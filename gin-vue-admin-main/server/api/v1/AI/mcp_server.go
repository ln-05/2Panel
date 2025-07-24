package AI

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/AI"
	AIReq "github.com/flipped-aurora/gin-vue-admin/server/model/AI/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type McpServerApi struct{}

// CreateMcpServer 创建mcpServer表
// @Tags McpServer
// @Summary 创建mcpServer表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.McpServer true "创建mcpServer表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /mcpServer/createMcpServer [post]
func (mcpServerApi *McpServerApi) CreateMcpServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var mcpServer AI.McpServer
	err := c.ShouldBindJSON(&mcpServer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	mcpServer.CreatedBy = utils.GetUserID(c)
	err = mcpServerService.CreateMcpServer(ctx, &mcpServer)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteMcpServer 删除mcpServer表
// @Tags McpServer
// @Summary 删除mcpServer表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.McpServer true "删除mcpServer表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /mcpServer/deleteMcpServer [delete]
func (mcpServerApi *McpServerApi) DeleteMcpServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := mcpServerService.DeleteMcpServer(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteMcpServerByIds 批量删除mcpServer表
// @Tags McpServer
// @Summary 批量删除mcpServer表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /mcpServer/deleteMcpServerByIds [delete]
func (mcpServerApi *McpServerApi) DeleteMcpServerByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := mcpServerService.DeleteMcpServerByIds(ctx, IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMcpServer 更新mcpServer表
// @Tags McpServer
// @Summary 更新mcpServer表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.McpServer true "更新mcpServer表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /mcpServer/updateMcpServer [put]
func (mcpServerApi *McpServerApi) UpdateMcpServer(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var mcpServer AI.McpServer
	err := c.ShouldBindJSON(&mcpServer)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	mcpServer.UpdatedBy = utils.GetUserID(c)
	err = mcpServerService.UpdateMcpServer(ctx, mcpServer)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindMcpServer 用id查询mcpServer表
// @Tags McpServer
// @Summary 用id查询mcpServer表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询mcpServer表"
// @Success 200 {object} response.Response{data=AI.McpServer,msg=string} "查询成功"
// @Router /mcpServer/findMcpServer [get]
func (mcpServerApi *McpServerApi) FindMcpServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	remcpServer, err := mcpServerService.GetMcpServer(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(remcpServer, c)
}

// GetMcpServerList 分页获取mcpServer表列表
// @Tags McpServer
// @Summary 分页获取mcpServer表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query AIReq.McpServerSearch true "分页获取mcpServer表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /mcpServer/getMcpServerList [get]
func (mcpServerApi *McpServerApi) GetMcpServerList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo AIReq.McpServerSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := mcpServerService.GetMcpServerInfoList(ctx, pageInfo)
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

// GetMcpServerPublic 不需要鉴权的mcpServer表接口
// @Tags McpServer
// @Summary 不需要鉴权的mcpServer表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /mcpServer/getMcpServerPublic [get]
func (mcpServerApi *McpServerApi) GetMcpServerPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	mcpServerService.GetMcpServerPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的mcpServer表接口信息",
	}, "获取成功", c)
}
