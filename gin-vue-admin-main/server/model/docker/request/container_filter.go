package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// ContainerFilter 容器过滤请求结构
type ContainerFilter struct {
	request.PageInfo
	Status string `json:"status" form:"status"` // 容器状态过滤 (running, exited, paused, etc.)
	Name   string `json:"name" form:"name"`     // 容器名称过滤
}

// LogOptions 容器日志选项
type LogOptions struct {
	Tail   string `json:"tail" form:"tail"`     // 显示最后N行日志
	Since  string `json:"since" form:"since"`   // 显示指定时间之后的日志
	Follow bool   `json:"follow" form:"follow"` // 是否跟踪日志流
}