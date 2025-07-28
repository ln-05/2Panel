package docker

import (
	"time"
)

// DockerConfigRequest Docker配置请求
type DockerConfigRequest struct {
	RegistryMirrors    []string          `json:"registryMirrors"`    // 镜像加速器
	InsecureRegistries []string          `json:"insecureRegistries"` // 不安全仓库
	PrivateRegistry    *PrivateRegistry  `json:"privateRegistry"`    // 私有仓库
	StorageDriver      string            `json:"storageDriver"`      // 存储驱动
	StorageOpts        map[string]string `json:"storageOpts"`        // 存储选项
	LogDriver          string            `json:"logDriver"`          // 日志驱动
	LogOpts            map[string]string `json:"logOpts"`            // 日志选项
	EnableIPv6         bool              `json:"enableIPv6"`         // 启用IPv6
	EnableIPForward    bool              `json:"enableIPForward"`    // 启用IP转发
	EnableIptables     bool              `json:"enableIptables"`     // 启用iptables
	LiveRestore        bool              `json:"liveRestore"`        // 实时恢复
	CgroupDriver       string            `json:"cgroupDriver"`       // Cgroup驱动
	SocketPath         string            `json:"socketPath"`         // Socket路径
	DataRoot           string            `json:"dataRoot"`           // 数据根目录
	ExecRoot           string            `json:"execRoot"`           // 执行根目录
}

// DockerConfigResponse Docker配置响应
type DockerConfigResponse struct {
	Config          *DockerConfigRequest `json:"config"`
	ConfigPath      string               `json:"configPath"`
	IsDefault       bool                 `json:"isDefault"`
	LastModified    time.Time            `json:"lastModified"`
	ServiceStatus   string               `json:"serviceStatus"`
	Version         string               `json:"version"`
	BackupAvailable bool                 `json:"backupAvailable"`
}

// PrivateRegistry 私有仓库配置
type PrivateRegistry struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

// ServiceOperationError 服务操作错误
type ServiceOperationError struct {
	Operation string `json:"operation"`
	Message   string `json:"message"`
	Code      string `json:"code"`
	Retryable bool   `json:"retryable"`
}

// BackupResponse 备份响应
type BackupResponse struct {
	BackupId    string    `json:"backupId"`
	BackupPath  string    `json:"backupPath"`
	CreatedAt   time.Time `json:"createdAt"`
	Description string    `json:"description"`
}

// RestoreRequest 恢复请求
type RestoreRequest struct {
	BackupId string `json:"backupId"`
}

// BackupListResponse 备份列表响应
type BackupListResponse struct {
	Backups []BackupInfo `json:"backups"`
	Total   int          `json:"total"`
}

// BackupInfo 备份信息
type BackupInfo struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Size        int64     `json:"size"`
	ConfigHash  string    `json:"configHash"`
}

// ServiceOperationResponse 服务操作响应
type ServiceOperationResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}