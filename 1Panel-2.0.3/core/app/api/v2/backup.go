package v2

import (
	"fmt"

	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/gin-gonic/gin"
)

// @Tags Backup Account
// @Summary Create backup account
// @Accept json
// @Param request body dto.BackupOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups [post]
// @x-panel-log {"bodyKeys":["type"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建备份账号 [type]","formatEN":"create backup account [type]"}
func (b *BaseApi) CreateBackup(c *gin.Context) {
	var req dto.BackupOperate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := backupService.Create(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Backup Account
// @Summary Refresh token
// @Accept json
// @Param request body dto.OperateByName true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups/refresh/token [post]
func (b *BaseApi) RefreshToken(c *gin.Context) {
	var req dto.OperateByName
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := backupService.RefreshToken(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Backup Account
// @Summary List buckets
// @Accept json
// @Param request body dto.ForBuckets true "request"
// @Success 200 {array} string
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups/buckets [post]
func (b *BaseApi) ListBuckets(c *gin.Context) {
	var req dto.ForBuckets
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	buckets, err := backupService.GetBuckets(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, buckets)
}

// @Tags Backup Account
// @Summary Load backup account base info
// @Accept json
// @Success 200 {object} dto.BackupClientInfo
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups/client/:clientType [get]
func (b *BaseApi) LoadBackupClientInfo(c *gin.Context) {
	clientType, ok := c.Params.Get("clientType")
	if !ok {
		helper.BadRequest(c, fmt.Errorf("error %s in path", "clientType"))
		return
	}
	data, err := backupService.LoadBackupClientInfo(clientType)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// @Tags Backup Account
// @Summary Delete backup account
// @Accept json
// @Param request body dto.OperateByName true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups/del [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"删除备份账号 [name]","formatEN":"delete backup account [name]"}
func (b *BaseApi) DeleteBackup(c *gin.Context) {
	var req dto.OperateByName
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := backupService.Delete(req.Name); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Backup Account
// @Summary Update backup account
// @Accept json
// @Param request body dto.BackupOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/backups/update [post]
// @x-panel-log {"bodyKeys":["type"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新备份账号 [types]","formatEN":"update backup account [types]"}
func (b *BaseApi) UpdateBackup(c *gin.Context) {
	var req dto.BackupOperate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := backupService.Update(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}
