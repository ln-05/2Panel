package system

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DatabaseApi struct{}

// CreateSysDatabase 创建数据库连接
// @Tags SysDatabase
// @Summary 创建数据库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SysDatabaseCreate true "创建数据库连接"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /database/create [post]
func (databaseApi *DatabaseApi) CreateSysDatabase(c *gin.Context) {
	var req systemReq.SysDatabaseCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	database := system.SysDatabase{
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		Description: req.Description,
	}

	err = service.ServiceGroupApp.SystemServiceGroup.DatabaseService.CreateSysDatabase(&database)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSysDatabase 删除数据库连接
// @Tags SysDatabase
// @Summary 删除数据库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除数据库连接"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /database/deleted [post]
func (databaseApi *DatabaseApi) DeleteSysDatabase(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.ServiceGroupApp.SystemServiceGroup.DatabaseService.DeleteSysDatabaseByIds(IDS.Ids)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateSysDatabase 更新数据库连接
// @Tags SysDatabase
// @Summary 更新数据库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SysDatabaseUpdate true "更新数据库连接"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /database/update [post]
func (databaseApi *DatabaseApi) UpdateSysDatabase(c *gin.Context) {
	var req systemReq.SysDatabaseUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	database := system.SysDatabase{
		GVA_MODEL:   global.GVA_MODEL{ID: req.ID},
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		Description: req.Description,
	}

	err = service.ServiceGroupApp.SystemServiceGroup.DatabaseService.UpdateSysDatabase(database)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysDatabase 用id查询数据库连接
// @Tags SysDatabase
// @Summary 用id查询数据库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysDatabaseSearch true "用id查询数据库连接"
// @Success 200 {object} response.Response{data=object{resysDatabase=system.SysDatabase},msg=string} "查询成功"
// @Router /database/id [get]
func (databaseApi *DatabaseApi) FindSysDatabase(c *gin.Context) {
	ID := c.Query("ID")
	resysDatabase, err := service.ServiceGroupApp.SystemServiceGroup.DatabaseService.GetSysDatabase(ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"resysDatabase": resysDatabase}, c)
}

// GetSysDatabaseList 分页获取数据库连接列表
// @Tags SysDatabase
// @Summary 分页获取数据库连接列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query systemReq.SysDatabaseSearch true "分页获取数据库连接列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /database/list [get]
func (databaseApi *DatabaseApi) GetSysDatabaseList(c *gin.Context) {
	var pageInfo systemReq.SysDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := service.ServiceGroupApp.SystemServiceGroup.DatabaseService.GetSysDatabaseInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// TestSysDatabase 测试数据库连接
// @Tags SysDatabase
// @Summary 测试数据库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SysDatabaseTest true "测试数据库连接"
// @Success 200 {object} response.Response{data=systemRes.SysDatabaseTestResponse} "测试完成"
// @Router /database/test [post]
func (databaseApi *DatabaseApi) TestSysDatabase(c *gin.Context) {
	var req systemReq.SysDatabaseTest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := service.ServiceGroupApp.SystemServiceGroup.DatabaseService.TestSysDatabase(strconv.Itoa(int(req.ID)))
	if err != nil {
		global.GVA_LOG.Error("测试失败!", zap.Error(err))
		response.FailWithMessage("测试失败:"+err.Error(), c)
		return
	}

	if result.Success {
		response.OkWithDetailed(result, "连接测试成功", c)
	} else {
		response.FailWithDetailed(result, "连接测试失败", c)
	}
}

// SyncSysDatabase 从远程服务器同步数据库配置
// @Tags SysDatabase
// @Summary 从远程服务器同步数据库配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SysDatabaseSync true "同步数据库配置"
// @Success 200 {object} response.Response{data=systemRes.SysDatabaseSyncResponse} "同步完成"
// @Router /database/sync [post]
func (databaseApi *DatabaseApi) SyncSysDatabase(c *gin.Context) {
	var req systemReq.SysDatabaseSync
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := service.ServiceGroupApp.SystemServiceGroup.DatabaseService.SyncSysDatabase(req)
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "同步完成", c)
}