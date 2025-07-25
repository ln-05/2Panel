package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// SysDatabase 数据库连接配置表
type SysDatabase struct {
	global.GVA_MODEL
	Name        string    `json:"name" gorm:"comment:数据库名称;not null"`
	Type        string    `json:"type" gorm:"comment:数据库类型;not null"`
	Host        string    `json:"host" gorm:"comment:主机地址;not null"`
	Port        int       `json:"port" gorm:"comment:端口号;not null"`
	Username    string    `json:"username" gorm:"comment:用户名;not null"`
	Password    string    `json:"password" gorm:"comment:密码;not null"`
	Database    string    `json:"database" gorm:"comment:数据库名"`
	Description string    `json:"description" gorm:"comment:描述信息"`
	Status      string    `json:"status" gorm:"comment:连接状态;default:unknown"`
	LastTestAt  *time.Time `json:"lastTestAt" gorm:"comment:最后测试时间"`
}

// TableName 设置表名
func (SysDatabase) TableName() string {
	return "sys_databases"
}