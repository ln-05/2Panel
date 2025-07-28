package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"go.uber.org/zap"
)

type DockerServiceController struct{}

// RestartService 重启Docker服务
func (c *DockerServiceController) RestartService() error {
	global.GVA_LOG.Info("开始重启Docker服务")

	// 根据操作系统选择重启命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "restart", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Restart-Service", "docker")
	case "darwin":
		// macOS通常使用Docker Desktop，需要特殊处理
		return c.restartDockerDesktop()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 执行重启命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("重启Docker服务失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("重启Docker服务失败: %v, 输出: %s", err, string(output))
	}

	global.GVA_LOG.Info("Docker服务重启命令执行成功")

	// 等待服务就绪
	if err := c.WaitForServiceReady(60 * time.Second); err != nil {
		return fmt.Errorf("Docker服务重启后未能正常启动: %v", err)
	}

	global.GVA_LOG.Info("Docker服务重启成功")
	return nil
}

// StopService 停止Docker服务
func (c *DockerServiceController) StopService() error {
	global.GVA_LOG.Info("开始停止Docker服务")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "stop", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Stop-Service", "docker")
	case "darwin":
		return fmt.Errorf("macOS上的Docker Desktop不支持通过命令行停止")
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("停止Docker服务失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("停止Docker服务失败: %v, 输出: %s", err, string(output))
	}

	global.GVA_LOG.Info("Docker服务停止成功")
	return nil
}

// StartService 启动Docker服务
func (c *DockerServiceController) StartService() error {
	global.GVA_LOG.Info("开始启动Docker服务")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "start", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Start-Service", "docker")
	case "darwin":
		return fmt.Errorf("macOS上的Docker Desktop不支持通过命令行启动")
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("启动Docker服务失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("启动Docker服务失败: %v, 输出: %s", err, string(output))
	}

	// 等待服务就绪
	if err := c.WaitForServiceReady(30 * time.Second); err != nil {
		return fmt.Errorf("Docker服务启动后未能正常运行: %v", err)
	}

	global.GVA_LOG.Info("Docker服务启动成功")
	return nil
}

// GetServiceStatus 获取Docker服务状态
func (c *DockerServiceController) GetServiceStatus() (*dockerModel.ServiceStatusResponse, error) {
	status := &dockerModel.ServiceStatusResponse{
		Status:   "unknown",
		Version:  "",
		ErrorMsg: "",
	}

	// 检查Docker守护进程是否可访问
	if global.GVA_DOCKER != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 尝试ping Docker守护进程
		if _, err := global.GVA_DOCKER.Ping(ctx); err == nil {
			status.Status = "running"
			
			// 获取Docker版本信息
			if version, err := global.GVA_DOCKER.ServerVersion(ctx); err == nil {
				status.Version = version.Version
			}

			// 获取系统信息
			if _, err := global.GVA_DOCKER.Info(ctx); err == nil {
				// 计算运行时间（如果可用）
				// Docker Info中没有直接的启动时间，这里使用一个近似值
				status.Uptime = "运行中"
			}
		} else {
			status.Status = "stopped"
			status.ErrorMsg = err.Error()
		}
	} else {
		status.Status = "error"
		status.ErrorMsg = "Docker客户端未初始化"
	}

	// 获取系统服务状态作为补充
	systemStatus, err := c.getSystemServiceStatus()
	if err == nil && systemStatus != "" {
		// 如果系统服务状态与Docker API状态不一致，优先使用系统服务状态
		if status.Status == "stopped" && systemStatus == "running" {
			status.Status = "starting"
			status.ErrorMsg = "Docker守护进程正在启动"
		}
	}

	global.GVA_LOG.Debug("获取Docker服务状态", zap.String("status", status.Status))
	return status, nil
}

// CheckServiceHealth 检查Docker服务健康状态
func (c *DockerServiceController) CheckServiceHealth() error {
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker客户端未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查Docker守护进程连接
	if _, err := global.GVA_DOCKER.Ping(ctx); err != nil {
		return fmt.Errorf("Docker守护进程不可访问: %v", err)
	}

	// 检查Docker版本
	if _, err := global.GVA_DOCKER.ServerVersion(ctx); err != nil {
		return fmt.Errorf("无法获取Docker版本: %v", err)
	}

	// 检查Docker系统信息
	if _, err := global.GVA_DOCKER.Info(ctx); err != nil {
		return fmt.Errorf("无法获取Docker系统信息: %v", err)
	}

	return nil
}

// WaitForServiceReady 等待Docker服务就绪
func (c *DockerServiceController) WaitForServiceReady(timeout time.Duration) error {
	global.GVA_LOG.Info("等待Docker服务就绪", zap.Duration("timeout", timeout))

	deadline := time.Now().Add(timeout)
	checkInterval := 2 * time.Second

	for time.Now().Before(deadline) {
		if err := c.CheckServiceHealth(); err == nil {
			global.GVA_LOG.Info("Docker服务已就绪")
			return nil
		}

		global.GVA_LOG.Debug("Docker服务尚未就绪，继续等待")
		time.Sleep(checkInterval)
	}

	return fmt.Errorf("等待Docker服务就绪超时")
}

// EnableService 启用Docker服务（开机自启）
func (c *DockerServiceController) EnableService() error {
	global.GVA_LOG.Info("启用Docker服务开机自启")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "enable", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Set-Service", "-Name", "docker", "-StartupType", "Automatic")
	case "darwin":
		// Docker Desktop通常默认开机自启
		global.GVA_LOG.Info("macOS上的Docker Desktop通常默认开机自启")
		return nil
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("启用Docker服务开机自启失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("启用Docker服务开机自启失败: %v", err)
	}

	global.GVA_LOG.Info("Docker服务开机自启启用成功")
	return nil
}

// DisableService 禁用Docker服务（开机自启）
func (c *DockerServiceController) DisableService() error {
	global.GVA_LOG.Info("禁用Docker服务开机自启")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "disable", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Set-Service", "-Name", "docker", "-StartupType", "Manual")
	case "darwin":
		return fmt.Errorf("macOS上的Docker Desktop开机自启需要在Docker Desktop设置中配置")
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("禁用Docker服务开机自启失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("禁用Docker服务开机自启失败: %v", err)
	}

	global.GVA_LOG.Info("Docker服务开机自启禁用成功")
	return nil
}

// ReloadService 重新加载Docker服务配置
func (c *DockerServiceController) ReloadService() error {
	global.GVA_LOG.Info("重新加载Docker服务配置")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		// 先重新加载systemd配置，然后重启Docker
		reloadCmd := exec.Command("systemctl", "daemon-reload")
		if output, err := reloadCmd.CombinedOutput(); err != nil {
			global.GVA_LOG.Warn("重新加载systemd配置失败", zap.Error(err), zap.String("output", string(output)))
		}
		cmd = exec.Command("systemctl", "reload-or-restart", "docker")
	case "windows":
		// Windows上重新加载配置通常需要重启服务
		return c.RestartService()
	case "darwin":
		// macOS上的Docker Desktop需要重启
		return c.restartDockerDesktop()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("重新加载Docker服务配置失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("重新加载Docker服务配置失败: %v", err)
	}

	// 等待服务就绪
	if err := c.WaitForServiceReady(30 * time.Second); err != nil {
		return fmt.Errorf("Docker服务重新加载后未能正常运行: %v", err)
	}

	global.GVA_LOG.Info("Docker服务配置重新加载成功")
	return nil
}

// 私有方法

// getSystemServiceStatus 获取系统服务状态
func (c *DockerServiceController) getSystemServiceStatus() (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("systemctl", "is-active", "docker")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Get-Service", "-Name", "docker", "|", "Select-Object", "-ExpandProperty", "Status")
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	status := strings.TrimSpace(string(output))
	
	// 标准化状态值
	switch strings.ToLower(status) {
	case "active", "running":
		return "running", nil
	case "inactive", "stopped":
		return "stopped", nil
	case "failed", "error":
		return "error", nil
	default:
		return status, nil
	}
}

// restartDockerDesktop 重启Docker Desktop (macOS)
func (c *DockerServiceController) restartDockerDesktop() error {
	global.GVA_LOG.Info("重启Docker Desktop")

	// 停止Docker Desktop
	stopCmd := exec.Command("osascript", "-e", "quit app \"Docker\"")
	if output, err := stopCmd.CombinedOutput(); err != nil {
		global.GVA_LOG.Warn("停止Docker Desktop失败", zap.Error(err), zap.String("output", string(output)))
	}

	// 等待一段时间确保完全停止
	time.Sleep(5 * time.Second)

	// 启动Docker Desktop
	startCmd := exec.Command("open", "-a", "Docker")
	if output, err := startCmd.CombinedOutput(); err != nil {
		global.GVA_LOG.Error("启动Docker Desktop失败", zap.Error(err), zap.String("output", string(output)))
		return fmt.Errorf("启动Docker Desktop失败: %v", err)
	}

	// 等待Docker Desktop启动
	if err := c.WaitForServiceReady(120 * time.Second); err != nil {
		return fmt.Errorf("Docker Desktop重启后未能正常启动: %v", err)
	}

	global.GVA_LOG.Info("Docker Desktop重启成功")
	return nil
}

// GetServiceLogs 获取Docker服务日志
func (c *DockerServiceController) GetServiceLogs(lines int) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		if lines > 0 {
			cmd = exec.Command("journalctl", "-u", "docker", "-n", fmt.Sprintf("%d", lines), "--no-pager")
		} else {
			cmd = exec.Command("journalctl", "-u", "docker", "--no-pager")
		}
	case "windows":
		// Windows事件日志查询
		cmd = exec.Command("powershell", "-Command", 
			fmt.Sprintf("Get-EventLog -LogName Application -Source Docker -Newest %d | Format-Table -AutoSize", lines))
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取Docker服务日志失败: %v", err)
	}

	return string(output), nil
}