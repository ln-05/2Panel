package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"go.uber.org/zap"
)

type DockerConfigFileManager struct{}

// GetConfigFilePath 获取Docker配置文件路径
func (m *DockerConfigFileManager) GetConfigFilePath() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/docker/daemon.json"
	case "windows":
		return `C:\ProgramData\docker\config\daemon.json`
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), ".docker", "daemon.json")
	default:
		return "/etc/docker/daemon.json"
	}
}

// GetBackupDir 获取备份目录
func (m *DockerConfigFileManager) GetBackupDir() string {
	switch runtime.GOOS {
	case "linux":
		return "/etc/docker/backups"
	case "windows":
		return `C:\ProgramData\docker\backups`
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), ".docker", "backups")
	default:
		return "/etc/docker/backups"
	}
}

// ReadConfigFile 读取配置文件
func (m *DockerConfigFileManager) ReadConfigFile() (*dockerModel.DockerConfigRequest, error) {
	configPath := m.GetConfigFilePath()
	
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		global.GVA_LOG.Info("Docker配置文件不存在，返回默认配置", zap.String("path", configPath))
		return m.getDefaultConfig(), nil
	}

	// 读取文件内容
	data, err := os.ReadFile(configPath)
	if err != nil {
		global.GVA_LOG.Error("读取Docker配置文件失败", zap.String("path", configPath), zap.Error(err))
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		global.GVA_LOG.Error("解析Docker配置文件失败", zap.Error(err))
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 转换为我们的配置结构
	config := m.convertRawConfigToRequest(rawConfig)
	
	global.GVA_LOG.Info("成功读取Docker配置文件", zap.String("path", configPath))
	return config, nil
}

// WriteConfigFile 写入配置文件
func (m *DockerConfigFileManager) WriteConfigFile(config *dockerModel.DockerConfigRequest) error {
	configPath := m.GetConfigFilePath()
	
	// 确保配置目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 转换为Docker daemon.json格式
	daemonConfig := m.convertRequestToDaemonConfig(config)

	// 序列化为JSON
	data, err := json.MarshalIndent(daemonConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		global.GVA_LOG.Error("写入Docker配置文件失败", zap.String("path", configPath), zap.Error(err))
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	global.GVA_LOG.Info("成功写入Docker配置文件", zap.String("path", configPath))
	return nil
}

// BackupConfigFile 备份配置文件
func (m *DockerConfigFileManager) BackupConfigFile() (string, error) {
	configPath := m.GetConfigFilePath()
	backupDir := m.GetBackupDir()
	
	// 确保备份目录存在
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("daemon_backup_%s.json", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	// 检查原配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果原文件不存在，创建一个空的备份
		emptyConfig := m.getDefaultConfig()
		data, _ := json.MarshalIndent(m.convertRequestToDaemonConfig(emptyConfig), "", "  ")
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return "", fmt.Errorf("创建空备份文件失败: %v", err)
		}
	} else {
		// 复制原文件到备份位置
		data, err := os.ReadFile(configPath)
		if err != nil {
			return "", fmt.Errorf("读取原配置文件失败: %v", err)
		}

		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			return "", fmt.Errorf("写入备份文件失败: %v", err)
		}
	}

	global.GVA_LOG.Info("成功创建Docker配置备份", zap.String("backup", backupPath))
	return backupPath, nil
}

// RestoreConfigFile 恢复配置文件
func (m *DockerConfigFileManager) RestoreConfigFile(backupPath string) error {
	configPath := m.GetConfigFilePath()

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupPath)
	}

	// 读取备份文件
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %v", err)
	}

	// 验证备份文件格式
	var testConfig map[string]interface{}
	if err := json.Unmarshal(data, &testConfig); err != nil {
		return fmt.Errorf("备份文件格式无效: %v", err)
	}

	// 恢复配置文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("恢复配置文件失败: %v", err)
	}

	global.GVA_LOG.Info("成功恢复Docker配置", zap.String("from", backupPath), zap.String("to", configPath))
	return nil
}

// CheckConfigFilePermissions 检查配置文件权限
func (m *DockerConfigFileManager) CheckConfigFilePermissions() error {
	configPath := m.GetConfigFilePath()
	configDir := filepath.Dir(configPath)

	// 检查目录权限
	if info, err := os.Stat(configDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("配置目录不存在: %s", configDir)
		}
		return fmt.Errorf("检查配置目录失败: %v", err)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("配置路径不是目录: %s", configDir)
		}
	}

	// 检查文件权限（如果文件存在）
	if info, err := os.Stat(configPath); err == nil {
		// 检查是否可读写
		if info.Mode().Perm()&0200 == 0 {
			return fmt.Errorf("配置文件不可写: %s", configPath)
		}
	}

	return nil
}

// GetConfigFileInfo 获取配置文件信息
func (m *DockerConfigFileManager) GetConfigFileInfo() (map[string]interface{}, error) {
	configPath := m.GetConfigFilePath()
	
	info := map[string]interface{}{
		"path":   configPath,
		"exists": false,
	}

	if stat, err := os.Stat(configPath); err == nil {
		info["exists"] = true
		info["size"] = stat.Size()
		info["modTime"] = stat.ModTime()
		info["mode"] = stat.Mode().String()
	}

	return info, nil
}

// getDefaultConfig 获取默认配置
func (m *DockerConfigFileManager) getDefaultConfig() *dockerModel.DockerConfigRequest {
	return &dockerModel.DockerConfigRequest{
		RegistryMirrors:    []string{},
		InsecureRegistries: []string{},
		PrivateRegistry:    nil,
		StorageDriver:      "",
		StorageOpts:        make(map[string]string),
		LogDriver:          "json-file",
		LogOpts:            make(map[string]string),
		EnableIPv6:         false,
		EnableIPForward:    true,
		EnableIptables:     true,
		LiveRestore:        false,
		CgroupDriver:       "systemd",
		SocketPath:         "",
		DataRoot:           "",
		ExecRoot:           "",
	}
}

// convertRawConfigToRequest 将原始配置转换为请求结构
func (m *DockerConfigFileManager) convertRawConfigToRequest(rawConfig map[string]interface{}) *dockerModel.DockerConfigRequest {
	config := m.getDefaultConfig()

	// 镜像加速器
	if mirrors, ok := rawConfig["registry-mirrors"].([]interface{}); ok {
		config.RegistryMirrors = make([]string, len(mirrors))
		for i, mirror := range mirrors {
			if str, ok := mirror.(string); ok {
				config.RegistryMirrors[i] = str
			}
		}
	}

	// 不安全仓库
	if insecure, ok := rawConfig["insecure-registries"].([]interface{}); ok {
		config.InsecureRegistries = make([]string, len(insecure))
		for i, reg := range insecure {
			if str, ok := reg.(string); ok {
				config.InsecureRegistries[i] = str
			}
		}
	}

	// 存储驱动
	if driver, ok := rawConfig["storage-driver"].(string); ok {
		config.StorageDriver = driver
	}

	// 存储选项
	if opts, ok := rawConfig["storage-opts"].(map[string]interface{}); ok {
		config.StorageOpts = make(map[string]string)
		for k, v := range opts {
			if str, ok := v.(string); ok {
				config.StorageOpts[k] = str
			}
		}
	}

	// 日志驱动
	if driver, ok := rawConfig["log-driver"].(string); ok {
		config.LogDriver = driver
	}

	// 日志选项
	if opts, ok := rawConfig["log-opts"].(map[string]interface{}); ok {
		config.LogOpts = make(map[string]string)
		for k, v := range opts {
			if str, ok := v.(string); ok {
				config.LogOpts[k] = str
			}
		}
	}

	// 布尔选项
	if ipv6, ok := rawConfig["ipv6"].(bool); ok {
		config.EnableIPv6 = ipv6
	}
	if ipForward, ok := rawConfig["ip-forward"].(bool); ok {
		config.EnableIPForward = ipForward
	}
	if iptables, ok := rawConfig["iptables"].(bool); ok {
		config.EnableIptables = iptables
	}
	if liveRestore, ok := rawConfig["live-restore"].(bool); ok {
		config.LiveRestore = liveRestore
	}

	// Cgroup驱动
	if driver, ok := rawConfig["exec-opts"].([]interface{}); ok {
		for _, opt := range driver {
			if str, ok := opt.(string); ok && len(str) > 20 && str[:20] == "native.cgroupdriver=" {
				config.CgroupDriver = str[20:]
				break
			}
		}
	}

	// 数据根目录
	if dataRoot, ok := rawConfig["data-root"].(string); ok {
		config.DataRoot = dataRoot
	}

	// 执行根目录
	if execRoot, ok := rawConfig["exec-root"].(string); ok {
		config.ExecRoot = execRoot
	}

	return config
}

// convertRequestToDaemonConfig 将请求结构转换为daemon.json格式
func (m *DockerConfigFileManager) convertRequestToDaemonConfig(config *dockerModel.DockerConfigRequest) map[string]interface{} {
	daemonConfig := make(map[string]interface{})

	// 镜像加速器
	if len(config.RegistryMirrors) > 0 {
		daemonConfig["registry-mirrors"] = config.RegistryMirrors
	}

	// 不安全仓库
	if len(config.InsecureRegistries) > 0 {
		daemonConfig["insecure-registries"] = config.InsecureRegistries
	}

	// 存储驱动
	if config.StorageDriver != "" {
		daemonConfig["storage-driver"] = config.StorageDriver
	}

	// 存储选项
	if len(config.StorageOpts) > 0 {
		daemonConfig["storage-opts"] = config.StorageOpts
	}

	// 日志驱动
	if config.LogDriver != "" {
		daemonConfig["log-driver"] = config.LogDriver
	}

	// 日志选项
	if len(config.LogOpts) > 0 {
		daemonConfig["log-opts"] = config.LogOpts
	}

	// 网络选项
	daemonConfig["ipv6"] = config.EnableIPv6
	daemonConfig["ip-forward"] = config.EnableIPForward
	daemonConfig["iptables"] = config.EnableIptables

	// 实时恢复
	daemonConfig["live-restore"] = config.LiveRestore

	// Cgroup驱动
	if config.CgroupDriver != "" {
		daemonConfig["exec-opts"] = []string{fmt.Sprintf("native.cgroupdriver=%s", config.CgroupDriver)}
	}

	// 数据根目录
	if config.DataRoot != "" {
		daemonConfig["data-root"] = config.DataRoot
	}

	// 执行根目录
	if config.ExecRoot != "" {
		daemonConfig["exec-root"] = config.ExecRoot
	}

	return daemonConfig
}

// ListBackups 列出所有备份文件
func (m *DockerConfigFileManager) ListBackups() ([]dockerModel.BackupInfo, error) {
	backupDir := m.GetBackupDir()
	
	// 检查备份目录是否存在
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []dockerModel.BackupInfo{}, nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}

	var backups []dockerModel.BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backup := dockerModel.BackupInfo{
			Id:          entry.Name(),
			Description: fmt.Sprintf("备份于 %s", info.ModTime().Format("2006-01-02 15:04:05")),
			CreatedAt:   info.ModTime(),
			Size:        info.Size(),
		}

		backups = append(backups, backup)
	}

	global.GVA_LOG.Info("列出Docker配置备份", zap.Int("count", len(backups)))
	return backups, nil
}

// DeleteBackup 删除备份文件
func (m *DockerConfigFileManager) DeleteBackup(backupId string) error {
	backupDir := m.GetBackupDir()
	backupPath := filepath.Join(backupDir, backupId)

	if err := os.Remove(backupPath); err != nil {
		return fmt.Errorf("删除备份文件失败: %v", err)
	}

	global.GVA_LOG.Info("成功删除Docker配置备份", zap.String("backup", backupId))
	return nil
}