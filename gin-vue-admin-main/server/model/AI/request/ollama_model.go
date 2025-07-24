package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type OllamaModelSearch struct {
	Name           string      `json:"name" form:"name"`                       // 模型名称模糊搜索
	Status         []string    `json:"status" form:"status[]"`                 // 多状态筛选
	From           []string    `json:"from" form:"from[]"`                     // 多来源筛选
	MinSize        string      `json:"minSize" form:"minSize"`                 // 最小大小
	MaxSize        string      `json:"maxSize" form:"maxSize"`                 // 最大大小
	SortBy         string      `json:"sortBy" form:"sortBy"`                   // 排序字段 (name, size, createdAt, status)
	SortDesc       bool        `json:"sortDesc" form:"sortDesc"`               // 是否降序
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"` // 创建时间范围
	request.PageInfo
}

type OllamaModelCreate struct {
	Name string `json:"name" form:"name" validate:"required"`
	From string `json:"from" form:"from"`
}

type OllamaModelSync struct {
	Force bool `json:"force" form:"force"`
}

type OllamaBindDomainRequest struct {
	Domain       string `json:"domain" validate:"required"`
	AppInstallID uint   `json:"appInstallID" validate:"required"`
	SSLID        uint   `json:"sslID"`
	WebsiteID    uint   `json:"websiteID"`
	IPList       string `json:"ipList"`
}

type OllamaChatRequest struct {
	ModelID string `json:"modelId" validate:"required"`
	Message string `json:"message" validate:"required"`
	Stream  bool   `json:"stream"`
	Context string `json:"context"`
}

type OllamaChatResponse struct {
	Response string `json:"response"`
	Context  string `json:"context"`
	Done     bool   `json:"done"`
}
