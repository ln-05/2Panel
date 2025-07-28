package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"go.uber.org/zap"
)

type DockerConfigService struct {
	fileManager *DockerConfigFileManager
	validator   *DockerConfigValidator
	controller  *DockerServiceController
}

// NewDockerConfigService 创建Docker配置服务实例
func NewDockerConfigService() *DockerConfigService {
	return &DockerConfigService{
		fileManager: &DockerConfigFileManager{},
		validator:   &DockerConfigValidator{},
		controller:  &DockerServiceController{},
	}
}

// GetDockerConfig 获取当前Docker配置
func (s *DockerConfigService) GetDockerConfig() (*dockerModel.DockerConfigResponse, error) {
	global.GVA_LOG.Info("开始获取Docker配置")

	// 读取配置文件
	config, err := s.fileManager.ReadConfigFile()
	if err != nil {
		global.GVA_LOG.Error("读取Docker配置文件失败", zap.Error(err))
		return nil, fmt.Errorf("读取Docker配置文件失败: %v", err)
	}

	// 获取配置文件信息
	configPath := s.fileManager.GetConfigFilePath()
	var lastModified time.Time
	var isDefault bool = config == nil

	if !isDefault {
		if stat, err := os.Stat(configPath); err == nil {
			lastModified = stat.ModTime()
		}
	}

	// 获取服务状态
	serviceStatus := "unknown"
	version := ""
	if statusResp, err := s.controller.GetServiceStatus(); err == nil {
		serviceStatus = statusResp.Status
		version = statusResp.Version
	}

	// 检查是否有备份可用
	backupAvailable := s.hasBackupAvailable()

	response := &dockerModel.DockerConfigResponse{
		Config:          config,
		ConfigPath:      configPath,
		IsDefault:       isDefault,
		LastModified:    lastModified,
		ServiceStatus:   serviceStatus,
		Version:         version,
		BackupAvailable: backupAvailable,
	}

	global.GVA_LOG.Info("获取Docker配置成功", zap.String("configPath", configPath))
	return response, nil
}

// UpdateDockerConfig 更新Docker配置
func (s *DockerConfigService) UpdateDockerConfig(config dockerModel.DockerConfigRequest) error {
	global.GVA_LOG.Info("开始更新Docker配置")

	// 验证配置
	if err := s.validateConfig(&config); err != nil {
		global.GVA_LOG.Error("Docker配置验证失败", zap.Error(err))
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 创建备份
	if err := s.createBackupBeforeUpdate(); err != nil {
		global.GVA_LOG.Warn("创建配置备份失败", zap.Error(err))
		// 备份失败不阻止配置更新，但记录警告
	}

	// 写入配置文件
	if err := s.fileManager.WriteConfigFile(&config); err != nil {
		global.GVA_LOG.Error("写入Docker配置文件失败", zap.Error(err))
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	global.GVA_LOG.Info("Docker配置更新成功")
	return nil
}

// ValidateConfig 验证Docker配置
func (s *DockerConfigService) ValidateConfig(config dockerModel.DockerConfigRequest) (*dockerModel.ValidationResponse, error) {
	global.GVA_LOG.Info("开始验证Docker配置")

	response := &dockerModel.ValidationResponse{
		Valid:  true,
		Errors: []dockerModel.ConfigValidationError{},
	}

	// 验证镜像加速器
	if err := s.validator.ValidateRegistryMirrors(config.RegistryMirrors); err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
			Field:   "registryMirrors",
			Message: err.Error(),
			Code:    "INVALID_REGISTRY_MIRROR",
		})
	}

	// 验证私有仓库
	if config.PrivateRegistry != nil {
		if err := s.validator.ValidatePrivateRegistry(config.PrivateRegistry); err != nil {
			response.Valid = false
			response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
				Field:   "privateRegistry",
				Message: err.Error(),
				Code:    "INVALID_PRIVATE_REGISTRY",
			})
		}
	}

	// 验证存储配置
	if err := s.validator.ValidateStorageConfig(config.StorageDriver, config.StorageOpts); err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
			Field:   "storageConfig",
			Message: err.Error(),
			Code:    "INVALID_STORAGE_CONFIG",
		})
	}

	// 验证网络配置
	if err := s.validator.ValidateNetworkConfig(&config); err != nil {
		response.Valid = false
		response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
			Field:   "networkConfig",
			Message: err.Error(),
			Code:    "INVALID_NETWORK_CONFIG",
		})
	}

	// 验证Socket路径
	if config.SocketPath != "" {
		if err := s.validator.ValidateSocketPath(config.SocketPath); err != nil {
			response.Valid = false
			response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
				Field:   "socketPath",
				Message: err.Error(),
				Code:    "INVALID_SOCKET_PATH",
			})
		}
	}

	// 验证数据目录
	if config.DataRoot != "" {
		if err := s.validator.ValidateDataRoot(config.DataRoot); err != nil {
			response.Valid = false
			response.Errors = append(response.Errors, dockerModel.ConfigValidationError{
				Field:   "dataRoot",
				Message: err.Error(),
				Code:    "INVALID_DATA_ROOT",
			})
		}
	}

	global.GVA_LOG.Info("Docker配置验证完成", zap.Bool("valid", response.Valid), zap.Int("errors", len(response.Errors)))
	return response, nil
}

// RestartDockerService 重启Docker服务
func (s *DockerConfigService) RestartDockerService() (*dockerModel.ServiceOperationResponse, error) {
	global.GVA_LOG.Info("开始重启Docker服务")

	response := &dockerModel.ServiceOperationResponse{
		Operation: "restart",
		Timestamp: time.Now(),
	}

	if err := s.controller.RestartService(); err != nil {
		response.Success = false
		response.Message = fmt.Sprintf("重启Docker服务失败: %v", err)
		global.GVA_LOG.Error("重启Docker服务失败", zap.Error(err))
		return response, err
	}

	response.Success = true
	response.Message = "Docker服务重启成功"
	global.GVA_LOG.Info("Docker服务重启成功")
	return response, nil
}

// GetServiceStatus 获取Docker服务状态
func (s *DockerConfigService) GetServiceStatus() (*dockerModel.ServiceStatusResponse, error) {
	return s.controller.GetServiceStatus()
}

// StartDockerService 启动Docker服务
func (s *DockerConfigService) StartDockerService() (*dockerModel.ServiceOperationResponse, error) {
	global.GVA_LOG.Info("开始启动Docker服务")

	response := &dockerModel.ServiceOperationResponse{
		Operation: "start",
		Timestamp: time.Now(),
	}

	if err := s.controller.StartService(); err != nil {
		response.Success = false
		response.Message = fmt.Sprintf("启动Docker服务失败: %v", err)
		global.GVA_LOG.Error("启动Docker服务失败", zap.Error(err))
		return response, err
	}

	response.Success = true
	response.Message = "Docker服务启动成功"
	global.GVA_LOG.Info("Docker服务启动成功")
	return response, nil
}

// StopDockerService 停止Docker服务
func (s *DockerConfigService) StopDockerService() (*dockerModel.ServiceOperationResponse, error) {
	global.GVA_LOG.Info("开始停止Docker服务")

	response := &dockerModel.ServiceOperationResponse{
		Operation: "stop",
		Timestamp: time.Now(),
	}

	if err := s.controller.StopService(); err != nil {
		response.Success = false
		response.Message = fmt.Sprintf("停止Docker服务失败: %v", err)
		global.GVA_LOG.Error("停止Docker服务失败", zap.Error(err))
		return response, err
	}

	response.Success = true
	response.Message = "Docker服务停止成功"
	global.GVA_LOG.Info("Docker服务停止成功")
	return response, nil
}

// 私有方法

// validateConfig 内部配置验证
func (s *DockerConfigService) validateConfig(config *dockerModel.DockerConfigRequest) error {
	validationResp, err := s.ValidateConfig(*config)
	if err != nil {
		return err
	}

	if !validationResp.Valid {
		return fmt.Errorf("配置验证失败，共有%d个错误", len(validationResp.Errors))
	}

	return nil
}

// createBackupBeforeUpdate 更新前创建备份
func (s *DockerConfigService) createBackupBeforeUpdate() error {
	// 读取当前配置
	config, err := s.fileManager.ReadConfigFile()
	if err != nil {
		return fmt.Errorf("读取当前配置失败: %v", err)
	}

	// 如果配置文件不存在或为空，不需要备份
	if config == nil {
		global.GVA_LOG.Info("当前配置为空，跳过备份")
		return nil
	}

	// 创建备份
	configPath := s.fileManager.GetConfigFilePath()
	backupPath := s.generateBackupPath(configPath)

	// 确保备份目录存在
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 复制配置文件到备份位置
	if err := s.copyFile(configPath, backupPath); err != nil {
		return fmt.Errorf("创建配置备份失败: %v", err)
	}

	global.GVA_LOG.Info("配置备份创建成功", zap.String("backupPath", backupPath))
	return nil
}

// generateBackupPath 生成备份路径
func (s *DockerConfigService) generateBackupPath(configPath string) string {
	dir := filepath.Dir(configPath)
	filename := filepath.Base(configPath)
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("%s.backup.%s", filename, timestamp)
	return filepath.Join(dir, "backups", backupFilename)
}

// hasBackupAvailable 检查是否有备份可用
func (s *DockerConfigService) hasBackupAvailable() bool {
	configPath := s.fileManager.GetConfigFilePath()
	backupDir := filepath.Join(filepath.Dir(configPath), "backups")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return false
	}

	// 检查备份目录是否有文件
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return false
	}

	return len(entries) > 0
}

// BackupDockerConfig 创建Docker配置备份
func (s *DockerConfigService) BackupDockerConfig(description string) (*dockerModel.BackupResponse, error) {
	global.GVA_LOG.Info("开始创建Docker配置备份", zap.String("description", description))

	// 读取当前配置
	config, err := s.fileManager.ReadConfigFile()
	if err != nil {
		global.GVA_LOG.Error("读取当前配置失败", zap.Error(err))
		return nil, fmt.Errorf("读取当前配置失败: %v", err)
	}

	// 如果配置为空，创建默认配置的备份
	var configData []byte
	if config != nil {
		configData, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("序列化配置失败: %v", err)
		}
	} else {
		configData = []byte("{}")
	}

	// 生成备份ID和路径
	configPath := s.fileManager.GetConfigFilePath()
	backupId := s.generateBackupId()
	backupPath := s.generateBackupPathWithId(configPath, backupId)

	// 确保备份目录存在
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		global.GVA_LOG.Error("创建备份目录失败", zap.Error(err))
		return nil, fmt.Errorf("创建备份目录失败: %v", err)
	}

	// 创建备份文件
	if err := s.copyFile(configPath, backupPath); err != nil {
		global.GVA_LOG.Error("创建备份文件失败", zap.Error(err))
		return nil, fmt.Errorf("创建备份文件失败: %v", err)
	}

	// 创建备份元数据
	metadata := s.createBackupMetadata(backupId, description, configData)
	metadataPath := backupPath + ".meta"
	if err := s.saveBackupMetadata(metadataPath, metadata); err != nil {
		global.GVA_LOG.Warn("保存备份元数据失败", zap.Error(err))
		// 元数据保存失败不影响备份创建
	}

	response := &dockerModel.BackupResponse{
		BackupId:    backupId,
		BackupPath:  backupPath,
		CreatedAt:   time.Now(),
		Description: description,
	}

	global.GVA_LOG.Info("Docker配置备份创建成功",
		zap.String("backupId", backupId),
		zap.String("backupPath", backupPath))
	return response, nil
}

// RestoreDockerConfig 从备份恢复Docker配置
func (s *DockerConfigService) RestoreDockerConfig(restoreReq dockerModel.RestoreRequest) error {
	global.GVA_LOG.Info("开始恢复Docker配置", zap.String("backupId", restoreReq.BackupId))

	// 查找备份文件
	backupPath, err := s.findBackupPath(restoreReq.BackupId)
	if err != nil {
		global.GVA_LOG.Error("查找备份文件失败", zap.Error(err))
		return fmt.Errorf("查找备份文件失败: %v", err)
	}

	// 验证备份文件
	if err := s.validateBackupFile(backupPath); err != nil {
		global.GVA_LOG.Error("备份文件验证失败", zap.Error(err))
		return fmt.Errorf("备份文件验证失败: %v", err)
	}

	// 创建当前配置的备份（恢复前备份）
	if err := s.createBackupBeforeRestore(); err != nil {
		global.GVA_LOG.Warn("恢复前创建备份失败", zap.Error(err))
		// 不阻止恢复操作
	}

	// 恢复配置文件
	if err := s.fileManager.RestoreConfigFile(backupPath); err != nil {
		global.GVA_LOG.Error("恢复配置文件失败", zap.Error(err))
		return fmt.Errorf("恢复配置文件失败: %v", err)
	}

	global.GVA_LOG.Info("Docker配置恢复成功", zap.String("backupId", restoreReq.BackupId))
	return nil
}

// GetBackupList 获取备份列表
func (s *DockerConfigService) GetBackupList() (*dockerModel.BackupListResponse, error) {
	global.GVA_LOG.Info("开始获取备份列表")

	// 获取配置路径
	configPath := s.fileManager.GetConfigFilePath()
	backupDir := filepath.Join(filepath.Dir(configPath), "backups")

	// 检查备份目录是否存在
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return &dockerModel.BackupListResponse{
			Backups: []dockerModel.BackupInfo{},
			Total:   0,
		}, nil
	}

	// 读取备份目录
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("读取备份目录失败: %v", err)
	}

	var backups []dockerModel.BackupInfo
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) == ".meta" {
			continue
		}

		backupPath := filepath.Join(backupDir, entry.Name())
		backupInfo, err := s.getBackupInfo(backupPath)
		if err != nil {
			global.GVA_LOG.Warn("获取备份信息失败", zap.String("path", backupPath), zap.Error(err))
			continue
		}

		backups = append(backups, *backupInfo)
	}

	response := &dockerModel.BackupListResponse{
		Backups: backups,
		Total:   len(backups),
	}

	global.GVA_LOG.Info("获取备份列表成功", zap.Int("count", len(backups)))
	return response, nil
}

// DeleteBackup 删除备份
func (s *DockerConfigService) DeleteBackup(backupId string) error {
	global.GVA_LOG.Info("开始删除备份", zap.String("backupId", backupId))

	// 查找备份文件
	backupPath, err := s.findBackupPath(backupId)
	if err != nil {
		return fmt.Errorf("查找备份文件失败: %v", err)
	}

	// 删除备份文件
	if err := os.Remove(backupPath); err != nil {
		global.GVA_LOG.Error("删除备份文件失败", zap.Error(err))
		return fmt.Errorf("删除备份文件失败: %v", err)
	}

	// 删除元数据文件
	metadataPath := backupPath + ".meta"
	if _, err := os.Stat(metadataPath); err == nil {
		if err := os.Remove(metadataPath); err != nil {
			global.GVA_LOG.Warn("删除备份元数据失败", zap.Error(err))
		}
	}

	global.GVA_LOG.Info("备份删除成功", zap.String("backupId", backupId))
	return nil
}

// CleanupOldBackups 清理过期备份
func (s *DockerConfigService) CleanupOldBackups(maxAge time.Duration, maxCount int) error {
	global.GVA_LOG.Info("开始清理过期备份",
		zap.Duration("maxAge", maxAge),
		zap.Int("maxCount", maxCount))

	backupList, err := s.GetBackupList()
	if err != nil {
		return fmt.Errorf("获取备份列表失败: %v", err)
	}

	var deletedCount int
	now := time.Now()

	// 按创建时间排序（最新的在前）
	backups := backupList.Backups
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].CreatedAt.Before(backups[j].CreatedAt) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}

	// 删除过期或超出数量限制的备份
	for i, backup := range backups {
		shouldDelete := false

		// 检查是否超过最大数量
		if maxCount > 0 && i >= maxCount {
			shouldDelete = true
		}

		// 检查是否超过最大年龄
		if maxAge > 0 && now.Sub(backup.CreatedAt) > maxAge {
			shouldDelete = true
		}

		if shouldDelete {
			if err := s.DeleteBackup(backup.Id); err != nil {
				global.GVA_LOG.Warn("删除过期备份失败",
					zap.String("backupId", backup.Id),
					zap.Error(err))
			} else {
				deletedCount++
			}
		}
	}

	global.GVA_LOG.Info("过期备份清理完成", zap.Int("deletedCount", deletedCount))
	return nil
}

// 备份相关私有方法

// generateBackupId 生成备份ID
func (s *DockerConfigService) generateBackupId() string {
	return fmt.Sprintf("backup_%d", time.Now().Unix())
}

// generateBackupPathWithId 使用ID生成备份路径
func (s *DockerConfigService) generateBackupPathWithId(configPath, backupId string) string {
	dir := filepath.Dir(configPath)
	filename := filepath.Base(configPath)
	backupFilename := fmt.Sprintf("%s_%s.backup", filename, backupId)
	return filepath.Join(dir, "backups", backupFilename)
}

// createBackupMetadata 创建备份元数据
func (s *DockerConfigService) createBackupMetadata(backupId, description string, configData []byte) map[string]interface{} {
	return map[string]interface{}{
		"id":          backupId,
		"description": description,
		"createdAt":   time.Now(),
		"configHash":  s.calculateConfigHash(configData),
		"size":        len(configData),
	}
}

// saveBackupMetadata 保存备份元数据
func (s *DockerConfigService) saveBackupMetadata(metadataPath string, metadata map[string]interface{}) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(metadataPath, data, 0644)
}

// findBackupPath 查找备份文件路径
func (s *DockerConfigService) findBackupPath(backupId string) (string, error) {
	configPath := s.fileManager.GetConfigFilePath()
	backupDir := filepath.Join(filepath.Dir(configPath), "backups")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if strings.Contains(filename, backupId) && strings.HasSuffix(filename, ".backup") {
			return filepath.Join(backupDir, filename), nil
		}
	}

	return "", fmt.Errorf("备份文件未找到: %s", backupId)
}

// validateBackupFile 验证备份文件
func (s *DockerConfigService) validateBackupFile(backupPath string) error {
	// 检查文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("备份文件不存在: %s", backupPath)
	}

	// 检查文件是否可读
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("无法读取备份文件: %v", err)
	}

	// 验证JSON格式
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("备份文件格式无效: %v", err)
	}

	return nil
}

// createBackupBeforeRestore 恢复前创建备份
func (s *DockerConfigService) createBackupBeforeRestore() error {
	description := fmt.Sprintf("恢复前自动备份 - %s", time.Now().Format("2006-01-02 15:04:05"))
	_, err := s.BackupDockerConfig(description)
	return err
}

// getBackupInfo 获取备份信息
func (s *DockerConfigService) getBackupInfo(backupPath string) (*dockerModel.BackupInfo, error) {
	// 获取文件信息
	stat, err := os.Stat(backupPath)
	if err != nil {
		return nil, err
	}

	// 从文件名提取备份ID
	filename := filepath.Base(backupPath)
	backupId := s.extractBackupIdFromFilename(filename)

	// 尝试读取元数据
	metadataPath := backupPath + ".meta"
	var description string
	var configHash string

	if metadataData, err := os.ReadFile(metadataPath); err == nil {
		var metadata map[string]interface{}
		if err := json.Unmarshal(metadataData, &metadata); err == nil {
			if desc, ok := metadata["description"].(string); ok {
				description = desc
			}
			if hash, ok := metadata["configHash"].(string); ok {
				configHash = hash
			}
		}
	}

	return &dockerModel.BackupInfo{
		Id:          backupId,
		Description: description,
		CreatedAt:   stat.ModTime(),
		Size:        stat.Size(),
		ConfigHash:  configHash,
	}, nil
}

// extractBackupIdFromFilename 从文件名提取备份ID
func (s *DockerConfigService) extractBackupIdFromFilename(filename string) string {
	// 文件名格式: daemon.json_backup_1234567890.backup
	parts := strings.Split(filename, "_")
	if len(parts) >= 2 {
		idPart := parts[len(parts)-1]
		return strings.TrimSuffix(idPart, ".backup")
	}
	return filename
}

// calculateConfigHash 计算配置哈希
func (s *DockerConfigService) calculateConfigHash(data []byte) string {
	// 简单的哈希实现，实际项目中可以使用更复杂的哈希算法
	return fmt.Sprintf("%x", len(data))
}

// 复制文件
func (s *DockerConfigService) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
