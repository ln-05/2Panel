// 自动生成模板OllamaModel
package AI

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ollamaModel表 结构体  OllamaModel
type OllamaModel struct {
	global.GVA_MODEL
	Name         string `json:"name" form:"name" gorm:"comment:模型名称;column:name;size:255;not null"`                         //模型名称
	Size         string `json:"size" form:"size" gorm:"comment:模型大小;column:size;size:255;"`                                 //模型大小
	From         string `json:"from" form:"from" gorm:"comment:来源;column:from;size:255;"`                                   //来源
	Status       string `json:"status" form:"status" gorm:"comment:状态;column:status;size:255;default:stopped"`              //状态
	Message      string `json:"message" form:"message" gorm:"comment:消息;column:message;type:text"`                          //消息
	LogFileExist bool   `json:"logFileExist" form:"logFileExist" gorm:"comment:日志文件存在;column:log_file_exist;default:false"` //日志文件存在
	CreatedBy    uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint   `gorm:"column:deleted_by;comment:删除者"`
}

// OllamaModelInfo DTO for API responses
type OllamaModelInfo struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Size         string    `json:"size"`
	From         string    `json:"from"`
	LogFileExist bool      `json:"logFileExist"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
}

// OllamaBindDomain 域名绑定结构体
type OllamaBindDomain struct {
	global.GVA_MODEL
	Domain       string `json:"domain" form:"domain" gorm:"comment:域名;column:domain;size:255;not null" validate:"required"`
	AppInstallID uint   `json:"appInstallID" form:"appInstallID" gorm:"comment:应用安装ID;column:app_install_id;not null" validate:"required"`
	SSLID        uint   `json:"sslID" form:"sslID" gorm:"comment:SSL证书ID;column:ssl_id"`
	WebsiteID    uint   `json:"websiteID" form:"websiteID" gorm:"comment:网站ID;column:website_id"`
	IPList       string `json:"ipList" form:"ipList" gorm:"comment:IP白名单;column:ip_list;type:text"`
	CreatedBy    uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName ollamaModel表 OllamaModel自定义表名 ollama_model
func (OllamaModel) TableName() string {
	return "ollama_model"
}

// TableName 域名绑定表 OllamaBindDomain自定义表名 ollama_bind_domain
func (OllamaBindDomain) TableName() string {
	return "ollama_bind_domain"
}

// ToInfo 转换为DTO
func (o *OllamaModel) ToInfo() OllamaModelInfo {
	return OllamaModelInfo{
		ID:           o.ID,
		Name:         o.Name,
		Size:         o.Size,
		From:         o.From,
		LogFileExist: o.LogFileExist,
		Status:       o.Status,
		Message:      o.Message,
		CreatedAt:    o.CreatedAt,
	}
}
