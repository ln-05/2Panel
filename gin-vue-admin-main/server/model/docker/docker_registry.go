package docker

import (
	"time"
	"gorm.io/gorm"
)

// DockerRegistry Docker仓库模型
type DockerRegistry struct {
	ID          uint           `json:"id" gorm:"primarykey"`                                    // 主键ID
	CreatedAt   time.Time      `json:"createdAt"`                                               // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`                                               // 更新时间
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`                                          // 删除时间
	Name        string         `json:"name" gorm:"column:name;type:varchar(100);not null;uniqueIndex"` // 仓库名称
	DownloadUrl string         `json:"downloadUrl" gorm:"column:download_url;type:varchar(500);not null"` // 下载地址
	Protocol    string         `json:"protocol" gorm:"column:protocol;type:varchar(10);not null;default:'https'"` // 协议
	Status      string         `json:"status" gorm:"column:status;type:varchar(20);not null;default:'active'"` // 状态
	Username    string         `json:"username" gorm:"column:username;type:varchar(100)"`       // 用户名
	Password    string         `json:"password" gorm:"column:password;type:varchar(200)"`       // 密码（加密存储）
	Description string         `json:"description" gorm:"column:description;type:text"`         // 描述
	IsDefault   bool           `json:"isDefault" gorm:"column:is_default;default:false"`        // 是否为默认仓库
	LastTestTime *time.Time    `json:"lastTestTime" gorm:"column:last_test_time"`               // 最后测试时间
	TestResult  string         `json:"testResult" gorm:"column:test_result;type:varchar(500)"`  // 测试结果
}

// TableName 设置表名
func (DockerRegistry) TableName() string {
	return "docker_registries"
}