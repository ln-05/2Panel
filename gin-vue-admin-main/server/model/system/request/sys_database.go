package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// SysDatabaseSearch 数据库连接搜索条件
type SysDatabaseSearch struct {
	system.SysDatabase
	request.PageInfo
	StartCreatedAt *string `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *string `json:"endCreatedAt" form:"endCreatedAt"`
}

// SysDatabaseCreate 创建数据库连接请求
type SysDatabaseCreate struct {
	Name        string `json:"name" binding:"required" validate:"required"`
	Type        string `json:"type" binding:"required" validate:"required"`
	Host        string `json:"host" binding:"required" validate:"required"`
	Port        int    `json:"port" binding:"required" validate:"required"`
	Username    string `json:"username" binding:"required" validate:"required"`
	Password    string `json:"password" binding:"required" validate:"required"`
	Database    string `json:"database"`
	Description string `json:"description"`
}

// SysDatabaseUpdate 更新数据库连接请求
type SysDatabaseUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Host        string `json:"host" binding:"required"`
	Port        int    `json:"port" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Database    string `json:"database"`
	Description string `json:"description"`
}

// SysDatabaseTest 测试数据库连接请求
type SysDatabaseTest struct {
	ID uint `json:"id" binding:"required"`
}

// SysDatabaseSync 数据库同步请求
type SysDatabaseSync struct {
	ServerURL string `json:"server_url" binding:"required"`
	APIKey    string `json:"api_key"`
	SyncType  string `json:"sync_type"`
	Overwrite bool   `json:"overwrite"`
}