package response

import "time"

// OrchestrationInfo 编排基本信息
type OrchestrationInfo struct {
	ID               uint      `json:"id"`               // 编排ID
	Name             string    `json:"name"`             // 编排名称
	Description      string    `json:"description"`      // 描述
	Status           string    `json:"status"`           // 状态 (running/stopped/error)
	ContainerCount   int       `json:"containerCount"`   // 容器数量
	ApplicationCount int       `json:"applicationCount"` // 应用数量
	Source           string    `json:"source"`           // 来源 (manual/imported/1panel)
	EditLink         string    `json:"editLink"`         // 编辑链接
	CreatedAt        time.Time `json:"createdAt"`        // 创建时间
	UpdatedAt        time.Time `json:"updatedAt"`        // 更新时间
	LastStartTime    *time.Time `json:"lastStartTime"`   // 最后启动时间
	LastStopTime     *time.Time `json:"lastStopTime"`    // 最后停止时间
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ID            uint              `json:"id"`            // 服务ID
	ServiceName   string            `json:"serviceName"`   // 服务名称
	Image         string            `json:"image"`         // 镜像名称
	ContainerName string            `json:"containerName"` // 容器名称
	Status        string            `json:"status"`        // 服务状态
	ContainerID   string            `json:"containerId"`   // 容器ID
	Ports         []PortMappingInfo `json:"ports"`         // 端口映射
	Volumes       []VolumeMountInfo `json:"volumes"`       // 卷挂载
	Environment   map[string]string `json:"environment"`   // 环境变量
	RestartPolicy string            `json:"restartPolicy"` // 重启策略
	NetworkMode   string            `json:"networkMode"`   // 网络模式
	DependsOn     []string          `json:"dependsOn"`     // 依赖服务
	Command       []string          `json:"command"`       // 启动命令
}

// PortMappingInfo 端口映射信息
type PortMappingInfo struct {
	HostPort      int    `json:"hostPort"`      // 主机端口
	ContainerPort int    `json:"containerPort"` // 容器端口
	Protocol      string `json:"protocol"`      // 协议 (tcp/udp)
}

// VolumeMountInfo 卷挂载信息
type VolumeMountInfo struct {
	HostPath      string `json:"hostPath"`      // 主机路径
	ContainerPath string `json:"containerPath"` // 容器路径
	ReadOnly      bool   `json:"readOnly"`      // 是否只读
	Type          string `json:"type"`          // 挂载类型 (bind/volume/tmpfs)
}

// OrchestrationDetail 编排详细信息
type OrchestrationDetail struct {
	OrchestrationInfo
	ComposeContent string        `json:"composeContent"` // Docker Compose内容
	Services       []ServiceInfo `json:"services"`       // 服务列表
	WorkingDir     string        `json:"workingDir"`     // 工作目录
	EnvFile        string        `json:"envFile"`        // 环境变量文件路径
}

// OrchestrationListResponse 编排列表响应
type OrchestrationListResponse struct {
	List  []OrchestrationInfo `json:"list"`  // 编排列表
	Total int64               `json:"total"` // 总数
}

// OrchestrationOperationResponse 编排操作响应
type OrchestrationOperationResponse struct {
	Success bool   `json:"success"` // 操作是否成功
	Message string `json:"message"` // 操作结果消息
	Status  string `json:"status"`  // 操作后的状态
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	SuccessCount int                    `json:"successCount"` // 成功数量
	FailureCount int                    `json:"failureCount"` // 失败数量
	Results      []OperationResult      `json:"results"`      // 详细结果
}

// OperationResult 操作结果
type OperationResult struct {
	ID      uint   `json:"id"`      // 编排ID
	Name    string `json:"name"`    // 编排名称
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 结果消息
}

// OrchestrationLogResponse 编排日志响应
type OrchestrationLogResponse struct {
	ServiceName string    `json:"serviceName"` // 服务名称
	Logs        []LogLine `json:"logs"`        // 日志行
}

// LogLine 日志行
type LogLine struct {
	Timestamp time.Time `json:"timestamp"` // 时间戳
	Level     string    `json:"level"`     // 日志级别
	Message   string    `json:"message"`   // 日志消息
	Source    string    `json:"source"`    // 日志来源 (stdout/stderr)
}

// OrchestrationStatusResponse 编排状态响应
type OrchestrationStatusResponse struct {
	ID               uint                `json:"id"`               // 编排ID
	Name             string              `json:"name"`             // 编排名称
	Status           string              `json:"status"`           // 整体状态
	ContainerCount   int                 `json:"containerCount"`   // 容器数量
	ApplicationCount int                 `json:"applicationCount"` // 应用数量
	Services         []ServiceStatusInfo `json:"services"`         // 服务状态列表
	UpdatedAt        time.Time           `json:"updatedAt"`        // 更新时间
}

// ServiceStatusInfo 服务状态信息
type ServiceStatusInfo struct {
	ServiceName   string `json:"serviceName"`   // 服务名称
	ContainerName string `json:"containerName"` // 容器名称
	ContainerID   string `json:"containerId"`   // 容器ID
	Status        string `json:"status"`        // 服务状态
	Health        string `json:"health"`        // 健康状态
}