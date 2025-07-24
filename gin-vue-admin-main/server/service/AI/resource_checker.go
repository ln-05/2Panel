package AI

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/AI"
	"go.uber.org/zap"
)

// ResourceChecker 资源检查器
type ResourceChecker struct {
	// 资源阈值配置
	MinDiskSpaceGB   float64 // 最小磁盘空间 (GB)
	MaxMemoryUsage   float64 // 最大内存使用率 (%)
	MaxCPUUsage      float64 // 最大CPU使用率 (%)
	MaxConcurrentOps int     // 最大并发操作数
}

// ResourceStatus 资源状态
type ResourceStatus struct {
	DiskSpaceGB      float64 `json:"diskSpaceGB"`      // 可用磁盘空间 (GB)
	MemoryUsage      float64 `json:"memoryUsage"`      // 内存使用率 (%)
	CPUUsage         float64 `json:"cpuUsage"`         // CPU使用率 (%)
	NetworkStatus    string  `json:"networkStatus"`    // 网络状态
	ConcurrentOps    int     `json:"concurrentOps"`    // 当前并发操作数
	IsResourceEnough bool    `json:"isResourceEnough"` // 资源是否充足
	Message          string  `json:"message"`          // 状态消息
}

// NewResourceChecker 创建资源检查器
func NewResourceChecker() *ResourceChecker {
	return &ResourceChecker{
		MinDiskSpaceGB:   5.0,  // 最少需要5GB磁盘空间
		MaxMemoryUsage:   80.0, // 内存使用率不超过80%
		MaxCPUUsage:      90.0, // CPU使用率不超过90%
		MaxConcurrentOps: 3,    // 最多3个并发操作
	}
}

// CheckResources 检查系统资源
func (rc *ResourceChecker) CheckResources() (*ResourceStatus, error) {
	status := &ResourceStatus{
		IsResourceEnough: true,
		Message:          "资源充足",
	}

	// 检查磁盘空间
	diskSpace, err := rc.getAvailableDiskSpace()
	if err != nil {
		global.GVA_LOG.Error("检查磁盘空间失败", zap.Error(err))
		status.Message = "无法检查磁盘空间"
		status.IsResourceEnough = false
	} else {
		status.DiskSpaceGB = diskSpace
		if diskSpace < rc.MinDiskSpaceGB {
			status.IsResourceEnough = false
			status.Message = fmt.Sprintf("磁盘空间不足，需要至少%.1fGB，当前可用%.1fGB", rc.MinDiskSpaceGB, diskSpace)
		}
	}

	// 检查内存使用
	memoryUsage, err := rc.getMemoryUsage()
	if err != nil {
		global.GVA_LOG.Error("检查内存使用失败", zap.Error(err))
	} else {
		status.MemoryUsage = memoryUsage
		if memoryUsage > rc.MaxMemoryUsage {
			status.IsResourceEnough = false
			status.Message = fmt.Sprintf("内存使用率过高，当前%.1f%%，建议低于%.1f%%", memoryUsage, rc.MaxMemoryUsage)
		}
	}

	// 检查CPU使用
	cpuUsage, err := rc.getCPUUsage()
	if err != nil {
		global.GVA_LOG.Error("检查CPU使用失败", zap.Error(err))
	} else {
		status.CPUUsage = cpuUsage
		if cpuUsage > rc.MaxCPUUsage {
			status.IsResourceEnough = false
			status.Message = fmt.Sprintf("CPU使用率过高，当前%.1f%%，建议低于%.1f%%", cpuUsage, rc.MaxCPUUsage)
		}
	}

	// 检查网络连接
	networkStatus, err := rc.checkNetworkConnection()
	if err != nil {
		global.GVA_LOG.Error("检查网络连接失败", zap.Error(err))
		status.NetworkStatus = "unknown"
	} else {
		status.NetworkStatus = networkStatus
		if networkStatus != "connected" {
			status.IsResourceEnough = false
			status.Message = "网络连接异常"
		}
	}

	// 检查并发操作数
	concurrentOps := rc.getCurrentDownloadCount()
	status.ConcurrentOps = concurrentOps
	if concurrentOps >= rc.MaxConcurrentOps {
		status.IsResourceEnough = false
		status.Message = fmt.Sprintf("并发操作数过多，当前%d个，最大允许%d个", concurrentOps, rc.MaxConcurrentOps)
	}

	return status, nil
}

// getAvailableDiskSpace 获取可用磁盘空间 (GB)
func (rc *ResourceChecker) getAvailableDiskSpace() (float64, error) {
	cmd := exec.Command("df", "-BG", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("执行df命令失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("df命令输出格式异常")
	}

	// 解析df输出，获取可用空间
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0, fmt.Errorf("df命令输出字段不足")
	}

	// 第4个字段是可用空间，格式如 "123G"
	availableStr := strings.TrimSuffix(fields[3], "G")
	available, err := strconv.ParseFloat(availableStr, 64)
	if err != nil {
		return 0, fmt.Errorf("解析可用空间失败: %v", err)
	}

	return available, nil
}

// getMemoryUsage 获取内存使用率 (%)
func (rc *ResourceChecker) getMemoryUsage() (float64, error) {
	cmd := exec.Command("free", "-m")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("执行free命令失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("free命令输出格式异常")
	}

	// 解析内存信息
	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return 0, fmt.Errorf("free命令输出字段不足")
	}

	total, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, fmt.Errorf("解析总内存失败: %v", err)
	}

	used, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return 0, fmt.Errorf("解析已用内存失败: %v", err)
	}

	usage := (used / total) * 100
	return usage, nil
}

// getCPUUsage 获取CPU使用率 (%)
func (rc *ResourceChecker) getCPUUsage() (float64, error) {
	// 使用top命令获取CPU使用率
	cmd := exec.Command("top", "-bn1")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("执行top命令失败: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "%Cpu(s)") {
			// 解析CPU使用率行，格式类似: %Cpu(s):  1.2 us,  0.8 sy,  0.0 ni, 97.9 id
			fields := strings.Fields(line)
			for i, field := range fields {
				if strings.Contains(field, "id") && i > 0 {
					// idle百分比
					idleStr := strings.TrimSuffix(fields[i-1], ",")
					idle, err := strconv.ParseFloat(idleStr, 64)
					if err == nil {
						return 100 - idle, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("无法解析CPU使用率")
}

// checkNetworkConnection 检查网络连接状态
func (rc *ResourceChecker) checkNetworkConnection() (string, error) {
	// 尝试ping ollama官方服务器
	cmd := exec.Command("ping", "-c", "1", "-W", "3", "ollama.ai")
	err := cmd.Run()
	if err != nil {
		// 尝试ping Google DNS
		cmd = exec.Command("ping", "-c", "1", "-W", "3", "8.8.8.8")
		err = cmd.Run()
		if err != nil {
			return "disconnected", nil
		}
	}

	return "connected", nil
}

// getCurrentDownloadCount 获取当前下载任务数量
func (rc *ResourceChecker) getCurrentDownloadCount() int {
	var count int64
	global.GVA_DB.Model(&AI.OllamaModel{}).Where("status = ?", StatusDownloading).Count(&count)
	return int(count)
}

// EstimateModelSize 预估模型大小
func (rc *ResourceChecker) EstimateModelSize(modelName string) float64 {
	// 常见模型大小映射表 (GB)
	modelSizes := map[string]float64{
		"llama2:7b":      3.8,
		"llama2:13b":     7.3,
		"llama2:70b":     39.0,
		"codellama:7b":   3.8,
		"codellama:13b":  7.3,
		"codellama:34b":  19.0,
		"mistral:7b":     4.1,
		"mixtral:8x7b":   26.0,
		"neural-chat:7b": 4.1,
		"starcode:7b":    4.3,
		"vicuna:7b":      3.8,
		"vicuna:13b":     7.3,
		"orca-mini:3b":   1.9,
		"orca-mini:7b":   3.8,
		"orca-mini:13b":  7.3,
	}

	// 精确匹配
	if size, exists := modelSizes[modelName]; exists {
		return size
	}

	// 模糊匹配
	for pattern, size := range modelSizes {
		if strings.Contains(modelName, strings.Split(pattern, ":")[0]) {
			if strings.Contains(modelName, "7b") {
				return size
			} else if strings.Contains(modelName, "13b") {
				return size * 1.9 // 13b模型大约是7b的1.9倍
			} else if strings.Contains(modelName, "70b") {
				return size * 10.3 // 70b模型大约是7b的10.3倍
			}
		}
	}

	// 默认预估大小
	if strings.Contains(modelName, "3b") {
		return 1.9
	} else if strings.Contains(modelName, "7b") {
		return 3.8
	} else if strings.Contains(modelName, "13b") {
		return 7.3
	} else if strings.Contains(modelName, "34b") {
		return 19.0
	} else if strings.Contains(modelName, "70b") {
		return 39.0
	}

	// 默认大小
	return 4.0
}

// GetCleanupSuggestions 获取清理建议
func (rc *ResourceChecker) GetCleanupSuggestions() []string {
	suggestions := []string{}

	// 检查临时文件
	suggestions = append(suggestions, "清理系统临时文件: sudo rm -rf /tmp/*")

	// 检查Docker缓存
	suggestions = append(suggestions, "清理Docker缓存: docker system prune -f")

	// 检查日志文件
	suggestions = append(suggestions, "清理旧日志文件: sudo journalctl --vacuum-time=7d")

	// 检查未使用的模型
	var unusedModels []AI.OllamaModel
	global.GVA_DB.Where("status = ? AND updated_at < ?", StatusStopped, time.Now().AddDate(0, 0, -30)).Find(&unusedModels)

	if len(unusedModels) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("删除%d个30天未使用的模型", len(unusedModels)))
	}

	return suggestions
}
