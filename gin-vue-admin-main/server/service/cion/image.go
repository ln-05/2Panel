
package cion

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cion"
    cionReq "github.com/flipped-aurora/gin-vue-admin/server/model/cion/request"
    "gorm.io/gorm"
)

type ImageService struct {}
// CreateImage 创建image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService) CreateImage(ctx context.Context, image *cion.Image) (err error) {
	err = global.GVA_DB.Create(image).Error
	return err
}

// DeleteImage 删除image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService)DeleteImage(ctx context.Context, ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&cion.Image{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&cion.Image{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteImageByIds 批量删除image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService)DeleteImageByIds(ctx context.Context, IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&cion.Image{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&cion.Image{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateImage 更新image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService)UpdateImage(ctx context.Context, image cion.Image) (err error) {
	err = global.GVA_DB.Model(&cion.Image{}).Where("id = ?",image.ID).Updates(&image).Error
	return err
}

// GetImage 根据ID获取image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService)GetImage(ctx context.Context, ID string) (image cion.Image, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&image).Error
	return
}
// GetImageInfoList 分页获取image表记录
// Author [yourname](https://github.com/yourname)
func (imageService *ImageService)GetImageInfoList(ctx context.Context, info cionReq.ImageSearch) (list []cion.Image, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&cion.Image{})
    var images []cion.Image
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&images).Error
	return  images, total, err
}
func (imageService *ImageService)GetImagePublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
