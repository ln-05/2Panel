package docker

import (
	"time"
	"gorm.io/gorm"
)

// DockerOrchestration Docker编排模型
type DockerOrchestration struct {
	ID               uint           `json:"id" gorm:"primarykey"`                                                    // 主键ID
	CreatedAt        time.Time      `json:"createdAt"`                                                               // 创建时间
	UpdatedAt        time.Time      `json:"updatedAt"`                                                               // 更新时间
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`                                                          // 删除时间
	Name             string         `json:"name" gorm:"column:name;type:varchar(100);not null;uniqueIndex"`         // 编排名称
	Description      string         `json:"description" gorm:"column:description;type:text"`                        // 描述
	ComposeContent   string         `json:"composeContent" gorm:"column:compose_content;type:longtext"`             // Docker Compose内容
	Status           string         `json:"status" gorm:"column:status;type:varchar(20);not null;default:'stopped'"` // 状态 (running/stopped/error)
	ContainerCount   int            `json:"containerCount" gorm:"column:container_count;default:0"`                 // 容器数量
	ApplicationCount int            `json:"applicationCount" gorm:"column:application_count;default:0"`             // 应用数量
	Source           string         `json:"source" gorm:"column:source;type:varchar(50);not null;default:'manual'"` // 来源 (manual/imported/1panel)
	EditLink         string         `json:"editLink" gorm:"column:edit_link;type:varchar(500)"`                     // 编辑链接
	WorkingDir       string         `json:"workingDir" gorm:"column:working_dir;type:varchar(500)"`                 // 工作目录
	EnvFile          string         `json:"envFile" gorm:"column:env_file;type:varchar(500)"`                       // 环境变量文件路径
	LastStartTime    *time.Time     `json:"lastStartTime" gorm:"column:last_start_time"`                            // 最后启动时间
	LastStopTime     *time.Time     `json:"lastStopTime" gorm:"column:last_stop_time"`                              // 最后停止时间
}

// TableName 设置表名
func (DockerOrchestration) TableName() string {
	return "docker_orchestrations"
}

// DockerOrchestrationService Docker编排服务配置
type DockerOrchestrationService struct {
	ID              uint   `json:"id" gorm:"primarykey"`                                                       // 主键ID
	OrchestrationID uint   `json:"orchestrationId" gorm:"column:orchestration_id;not null;index"`            // 编排ID
	ServiceName     string `json:"serviceName" gorm:"column:service_name;type:varchar(100);not null"`        // 服务名称
	Image           string `json:"image" gorm:"column:image;type:varchar(200);not null"`                     // 镜像名称
	ContainerName   string `json:"containerName" gorm:"column:container_name;type:varchar(100)"`             // 容器名称
	Ports           string `json:"ports" gorm:"column:ports;type:text"`                                      // 端口映射 (JSON格式)
	Volumes         string `json:"volumes" gorm:"column:volumes;type:text"`                                  // 卷挂载 (JSON格式)
	Environment     string `json:"environment" gorm:"column:environment;type:text"`                         // 环境变量 (JSON格式)
	RestartPolicy   string `json:"restartPolicy" gorm:"column:restart_policy;type:varchar(50);default:'no'"` // 重启策略
	NetworkMode     string `json:"networkMode" gorm:"column:network_mode;type:varchar(100)"`                // 网络模式
	DependsOn       string `json:"dependsOn" gorm:"column:depends_on;type:text"`                            // 依赖服务 (JSON格式)
	Command         string `json:"command" gorm:"column:command;type:text"`                                  // 启动命令
	Status          string `json:"status" gorm:"column:status;type:varchar(20);default:'stopped'"`          // 服务状态
	ContainerID     string `json:"containerId" gorm:"column:container_id;type:varchar(100)"`                // 容器ID
	CreatedAt       time.Time `json:"createdAt"`                                                             // 创建时间
	UpdatedAt       time.Time `json:"updatedAt"`                                                             // 更新时间
}

// TableName 设置表名
func (DockerOrchestrationService) TableName() string {
	return "docker_orchestration_services"
}

// 建立关联关系
func (d *DockerOrchestration) Services() []DockerOrchestrationService {
	var services []DockerOrchestrationService
	// 这里可以通过GORM的关联查询来获取服务列表
	return services
}