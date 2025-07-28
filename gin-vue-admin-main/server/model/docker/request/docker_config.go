package request

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

// PrivateRegistry 私有仓库配置
type PrivateRegistry struct {
	URL      string `json:"url" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// ConfigValidationRequest 配置验证请求
type ConfigValidationRequest struct {
	Config *DockerConfigRequest `json:"config" binding:"required"`
}

// RestoreConfigRequest 恢复配置请求
type RestoreConfigRequest struct {
	BackupID string `json:"backupId" binding:"required"`
}