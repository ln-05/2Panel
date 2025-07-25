// 自动生成模板Image
package cion

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// image表 结构体  Image
type Image struct {
	global.GVA_MODEL
	ImageId      string    `json:"imageId" form:"imageId" gorm:"comment:镜像ID;column:image_id;size:64;"`             //镜像ID
	Repository   string    `json:"repository" form:"repository" gorm:"comment:仓库名;column:repository;size:255;"`     //仓库名
	Tag          string    `json:"tag" form:"tag" gorm:"comment:标签;column:tag;size:100;"`                           //标签
	Size         string    `json:"size" form:"size" gorm:"comment:大小;column:size;size:50;"`                         //大小
	CreatedTime  time.Time `json:"createdTime" form:"createdTime" gorm:"comment:创建时间;column:created_time;"`         //创建时间
	Architecture string    `json:"architecture" form:"architecture" gorm:"comment:架构;column:architecture;size:50;"` //架构
	Digest       string    `json:"digest" form:"digest" gorm:"comment:摘要;column:digest;size:255;"`                  //摘要
	CreatedBy    int       `json:"createdBy" form:"createdBy" gorm:"comment:创建者;column:created_by;size:19;"`        //创建者
	UpdatedBy    int       `json:"updatedBy" form:"updatedBy" gorm:"comment:更新者;column:updated_by;size:19;"`        //更新者
	DeletedBy    int       `json:"deletedBy" form:"deletedBy" gorm:"comment:删除者;column:deleted_by;size:19;"`        //删除者

}

// TableName image表 Image自定义表名 image
func (Image) TableName() string {
	return "image"
}
