package response

import (
	"time"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
)

// DockerConfigResponse Docker配置响应
type DockerConfigResponse struct {
	Config          *dockerReq.DockerConfigRequest `json:"config"`
	ConfigPath      string                         `json:"configPath"`
	IsDefault       bool                           `json:"isDefault"`
	LastModified    time.Time                      `json:"lastModified"`
	ServiceStatus   string                         `json:"serviceStatus"`
	Version         string                         `json:"version"`
	BackupAvailable bool                           `json:"backupAvailable"`
}

// ServiceStatusResponse 服务状态响应
type ServiceStatusResponse struct {
	Status      string    `json:"status"`      // running, stopped, error
	Uptime      string    `json:"uptime"`      // 运行时间
	Version     string    `json:"version"`     // Docker版本
	LastRestart time.Time `json:"lastRestart"` // 最后重启时间
	ErrorMsg    string    `json:"errorMsg"`    // 错误信息
}

// ConfigValidationError 配置验证错误
type ConfigValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ValidationResponse 验证响应
type ValidationResponse struct {
	Valid  bool                     `json:"valid"`
	Errors []ConfigValidationError  `json:"errors"`
}

// BackupInfo 备份信息
type BackupInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	Size        int64     `json:"size"`
	Description string    `json:"description"`
}

// BackupResponse 备份响应
type BackupResponse struct {
	BackupID string `json:"backupId"`
	Message  string `json:"message"`
}

// BackupListResponse 备份列表响应
type BackupListResponse struct {
	Backups []BackupInfo `json:"backups"`
	Total   int          `json:"total"`
}

// ServiceOperationResponse 服务操作响应
type ServiceOperationResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
}