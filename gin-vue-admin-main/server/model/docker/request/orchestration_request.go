package request

// OrchestrationFilter 编排过滤器
type OrchestrationFilter struct {
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小
	Name     string `form:"name" json:"name"`         // 编排名称过滤
	Status   string `form:"status" json:"status"`     // 状态过滤 (running/stopped/error)
	Source   string `form:"source" json:"source"`     // 来源过滤 (manual/imported/1panel)
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	ServiceName   string            `json:"serviceName" binding:"required"` // 服务名称
	Image         string            `json:"image" binding:"required"`       // 镜像名称
	ContainerName string            `json:"containerName"`                  // 容器名称
	Ports         []PortMapping     `json:"ports"`                          // 端口映射
	Volumes       []VolumeMount     `json:"volumes"`                        // 卷挂载
	Environment   map[string]string `json:"environment"`                    // 环境变量
	RestartPolicy string            `json:"restartPolicy"`                  // 重启策略
	NetworkMode   string            `json:"networkMode"`                    // 网络模式
	DependsOn     []string          `json:"dependsOn"`                      // 依赖服务
	Command       []string          `json:"command"`                        // 启动命令
}

// PortMapping 端口映射
type PortMapping struct {
	HostPort      int    `json:"hostPort"`      // 主机端口
	ContainerPort int    `json:"containerPort"` // 容器端口
	Protocol      string `json:"protocol"`      // 协议 (tcp/udp)
}

// VolumeMount 卷挂载
type VolumeMount struct {
	HostPath      string `json:"hostPath"`      // 主机路径
	ContainerPath string `json:"containerPath"` // 容器路径
	ReadOnly      bool   `json:"readOnly"`      // 是否只读
	Type          string `json:"type"`          // 挂载类型 (bind/volume/tmpfs)
}

// OrchestrationCreateRequest 创建编排请求
type OrchestrationCreateRequest struct {
	Name           string          `json:"name" binding:"required"`    // 编排名称
	Description    string          `json:"description"`                // 描述
	ComposeContent string          `json:"composeContent"`             // Docker Compose内容
	Services       []ServiceConfig `json:"services"`                   // 服务配置列表
	WorkingDir     string          `json:"workingDir"`                 // 工作目录
	EnvFile        string          `json:"envFile"`                    // 环境变量文件路径
}

// OrchestrationUpdateRequest 更新编排请求
type OrchestrationUpdateRequest struct {
	ID             uint            `json:"id" binding:"required"`      // 编排ID
	Name           string          `json:"name" binding:"required"`    // 编排名称
	Description    string          `json:"description"`                // 描述
	ComposeContent string          `json:"composeContent"`             // Docker Compose内容
	Services       []ServiceConfig `json:"services"`                   // 服务配置列表
	WorkingDir     string          `json:"workingDir"`                 // 工作目录
	EnvFile        string          `json:"envFile"`                    // 环境变量文件路径
}

// OrchestrationOperationRequest 编排操作请求
type OrchestrationOperationRequest struct {
	ID        uint   `json:"id" binding:"required"`        // 编排ID
	Operation string `json:"operation" binding:"required"` // 操作类型 (start/stop/restart)
}

// BatchOrchestrationOperationRequest 批量编排操作请求
type BatchOrchestrationOperationRequest struct {
	IDs       []uint `json:"ids" binding:"required"`       // 编排ID列表
	Operation string `json:"operation" binding:"required"` // 操作类型 (start/stop/restart/delete)
}

// OrchestrationLogRequest 编排日志请求
type OrchestrationLogRequest struct {
	ID         uint   `form:"id" json:"id" binding:"required"`     // 编排ID
	ServiceName string `form:"serviceName" json:"serviceName"`     // 服务名称（可选，为空则获取所有服务日志）
	Lines      int    `form:"lines" json:"lines"`                  // 日志行数
	Follow     bool   `form:"follow" json:"follow"`                // 是否跟踪日志
	Timestamps bool   `form:"timestamps" json:"timestamps"`        // 是否显示时间戳
	Since      string `form:"since" json:"since"`                  // 开始时间
	Until      string `form:"until" json:"until"`                  // 结束时间
}