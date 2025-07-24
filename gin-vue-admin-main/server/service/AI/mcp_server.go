
package AI

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/AI"
    AIReq "github.com/flipped-aurora/gin-vue-admin/server/model/AI/request"
    "gorm.io/gorm"
)

type McpServerService struct {}
// CreateMcpServer 创建mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService) CreateMcpServer(ctx context.Context, mcpServer *AI.McpServer) (err error) {
	err = global.GVA_DB.Create(mcpServer).Error
	return err
}

// DeleteMcpServer 删除mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService)DeleteMcpServer(ctx context.Context, ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&AI.McpServer{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&AI.McpServer{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteMcpServerByIds 批量删除mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService)DeleteMcpServerByIds(ctx context.Context, IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&AI.McpServer{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&AI.McpServer{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateMcpServer 更新mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService)UpdateMcpServer(ctx context.Context, mcpServer AI.McpServer) (err error) {
	err = global.GVA_DB.Model(&AI.McpServer{}).Where("id = ?",mcpServer.ID).Updates(&mcpServer).Error
	return err
}

// GetMcpServer 根据ID获取mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService)GetMcpServer(ctx context.Context, ID string) (mcpServer AI.McpServer, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&mcpServer).Error
	return
}
// GetMcpServerInfoList 分页获取mcpServer表记录
// Author [yourname](https://github.com/yourname)
func (mcpServerService *McpServerService)GetMcpServerInfoList(ctx context.Context, info AIReq.McpServerSearch) (list []AI.McpServer, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&AI.McpServer{})
    var mcpServers []AI.McpServer
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

	err = db.Find(&mcpServers).Error
	return  mcpServers, total, err
}
func (mcpServerService *McpServerService)GetMcpServerPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
