package AI

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/AI"
	AIReq "github.com/flipped-aurora/gin-vue-admin/server/model/AI/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OllamaModelApi struct{}

// CreateOllamaModel 创建ollamaModel表
// @Tags OllamaModel
// @Summary 创建ollamaModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.OllamaModel true "创建ollamaModel表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ollamaModel/createOllamaModel [post]
func (ollamaModelApi *OllamaModelApi) CreateOllamaModel(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var ollamaModel AI.OllamaModel
	err := c.ShouldBindJSON(&ollamaModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ollamaModel.CreatedBy = utils.GetUserID(c)
	err = ollamaModelService.CreateOllamaModel(ctx, &ollamaModel)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteOllamaModel 删除ollamaModel表
// @Tags OllamaModel
// @Summary 删除ollamaModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.OllamaModel true "删除ollamaModel表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /ollamaModel/deleteOllamaModel [delete]
func (ollamaModelApi *OllamaModelApi) DeleteOllamaModel(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := ollamaModelService.DeleteOllamaModel(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteOllamaModelByIds 批量删除ollamaModel表
// @Tags OllamaModel
// @Summary 批量删除ollamaModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /ollamaModel/deleteOllamaModelByIds [delete]
func (ollamaModelApi *OllamaModelApi) DeleteOllamaModelByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := ollamaModelService.DeleteOllamaModelByIds(ctx, IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateOllamaModel 更新ollamaModel表
// @Tags OllamaModel
// @Summary 更新ollamaModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AI.OllamaModel true "更新ollamaModel表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ollamaModel/updateOllamaModel [put]
func (ollamaModelApi *OllamaModelApi) UpdateOllamaModel(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var ollamaModel AI.OllamaModel
	err := c.ShouldBindJSON(&ollamaModel)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ollamaModel.UpdatedBy = utils.GetUserID(c)
	err = ollamaModelService.UpdateOllamaModel(ctx, ollamaModel)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindOllamaModel 用id查询ollamaModel表
// @Tags OllamaModel
// @Summary 用id查询ollamaModel表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询ollamaModel表"
// @Success 200 {object} response.Response{data=AI.OllamaModel,msg=string} "查询成功"
// @Router /ollamaModel/findOllamaModel [get]
func (ollamaModelApi *OllamaModelApi) FindOllamaModel(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reollamaModel, err := ollamaModelService.GetOllamaModel(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reollamaModel, c)
}

// GetOllamaModelList 分页获取ollamaModel表列表
// @Tags OllamaModel
// @Summary 分页获取ollamaModel表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query AIReq.OllamaModelSearch true "分页获取ollamaModel表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /ollamaModel/getOllamaModelList [get]
func (ollamaModelApi *OllamaModelApi) GetOllamaModelList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo AIReq.OllamaModelSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := ollamaModelService.GetOllamaModelInfoList(ctx, pageInfo)
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

// GetOllamaModelPublic 不需要鉴权的ollamaModel表接口
// @Tags OllamaModel
// @Summary 不需要鉴权的ollamaModel表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ollamaModel/getOllamaModelPublic [get]
func (ollamaModelApi *OllamaModelApi) GetOllamaModelPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	ollamaModelService.GetOllamaModelPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的ollamaModel表接口信息",
	}, "获取成功", c)
}

// SearchOllamaModel 搜索和分页查询模型
// @Tags OllamaModel
// @Summary 搜索和分页查询模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query AIReq.OllamaModelSearch true "搜索条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "搜索成功"
// @Router /ollamaModel/search [get]
func (ollamaModelApi *OllamaModelApi) SearchOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	var searchInfo AIReq.OllamaModelSearch
	err := c.ShouldBindQuery(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := ollamaModelService.Search(ctx, searchInfo)
	if err != nil {
		global.GVA_LOG.Error("搜索失败!", zap.Error(err))
		response.FailWithMessage("搜索失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchInfo.Page,
		PageSize: searchInfo.PageSize,
	}, "搜索成功", c)
}

// CreateOllamaModelAdvanced 创建/下载新模型
// @Tags OllamaModel
// @Summary 创建/下载新模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AIReq.OllamaModelCreate true "创建模型参数"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ollamaModel/create [post]
func (ollamaModelApi *OllamaModelApi) CreateOllamaModelAdvanced(c *gin.Context) {
	ctx := c.Request.Context()

	var req AIReq.OllamaModelCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	err = ollamaModelService.Create(ctx, req, userID)
	if err != nil {
		global.GVA_LOG.Error("创建模型失败!", zap.Error(err))
		response.FailWithMessage("创建模型失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("模型创建任务已启动", c)
}

// StartOllamaModel 启动模型
// @Tags OllamaModel
// @Summary 启动模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Success 200 {object} response.Response{msg=string} "启动成功"
// @Router /ollamaModel/start [post]
func (ollamaModelApi *OllamaModelApi) StartOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	err := ollamaModelService.Start(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("启动模型失败!", zap.Error(err))
		response.FailWithMessage("启动模型失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("模型已启动", c)
}

// StopOllamaModel 停止模型
// @Tags OllamaModel
// @Summary 停止模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /ollamaModel/stop [post]
func (ollamaModelApi *OllamaModelApi) StopOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	err := ollamaModelService.Stop(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("停止模型失败!", zap.Error(err))
		response.FailWithMessage("停止模型失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("模型已停止", c)
}

// CloseOllamaModel 停止模型 (兼容旧接口)
// @Tags OllamaModel
// @Summary 停止模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /ollamaModel/close [post]
func (ollamaModelApi *OllamaModelApi) CloseOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	err := ollamaModelService.Close(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("停止模型失败!", zap.Error(err))
		response.FailWithMessage("停止模型失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("模型已停止", c)
}

// RecreateOllamaModel 重新创建模型
// @Tags OllamaModel
// @Summary 重新创建模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Success 200 {object} response.Response{msg=string} "重新创建成功"
// @Router /ollamaModel/recreate [post]
func (ollamaModelApi *OllamaModelApi) RecreateOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	err := ollamaModelService.Recreate(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("重新创建模型失败!", zap.Error(err))
		response.FailWithMessage("重新创建模型失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("模型重新创建任务已启动", c)
}

// SyncOllamaModel 同步本地和远程模型
// @Tags OllamaModel
// @Summary 同步本地和远程模型
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AIReq.OllamaModelSync true "同步参数"
// @Success 200 {object} response.Response{data=map[string]int,msg=string} "同步成功"
// @Router /ollamaModel/sync [post]
func (ollamaModelApi *OllamaModelApi) SyncOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	var req AIReq.OllamaModelSync
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := ollamaModelService.Sync(ctx, req)
	if err != nil {
		global.GVA_LOG.Error("同步模型失败!", zap.Error(err))
		response.FailWithMessage("同步模型失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "同步完成", c)
}

// LoadOllamaModelDetail 加载模型详情
// @Tags OllamaModel
// @Summary 加载模型详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Success 200 {object} response.Response{data=AI.OllamaModelInfo,msg=string} "获取成功"
// @Router /ollamaModel/detail [get]
func (ollamaModelApi *OllamaModelApi) LoadOllamaModelDetail(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	detail, err := ollamaModelService.LoadDetail(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("获取模型详情失败!", zap.Error(err))
		response.FailWithMessage("获取模型详情失败:"+err.Error(), c)
		return
	}

	response.OkWithData(detail, c)
}

// BindDomainToOllama 绑定域名到AI服务
// @Tags OllamaModel
// @Summary 绑定域名到AI服务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AIReq.OllamaBindDomainRequest true "域名绑定参数"
// @Success 200 {object} response.Response{msg=string} "绑定成功"
// @Router /ollamaModel/bindDomain [post]
func (ollamaModelApi *OllamaModelApi) BindDomainToOllama(c *gin.Context) {
	ctx := c.Request.Context()

	var req AIReq.OllamaBindDomainRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	err = ollamaModelService.BindDomain(ctx, req, userID)
	if err != nil {
		global.GVA_LOG.Error("绑定域名失败!", zap.Error(err))
		response.FailWithMessage("绑定域名失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("域名绑定成功", c)
}

// GetOllamaBindDomain 获取绑定的域名信息
// @Tags OllamaModel
// @Summary 获取绑定的域名信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param domain query string true "域名"
// @Success 200 {object} response.Response{data=AI.OllamaBindDomain,msg=string} "获取成功"
// @Router /ollamaModel/getBindDomain [get]
func (ollamaModelApi *OllamaModelApi) GetOllamaBindDomain(c *gin.Context) {
	ctx := c.Request.Context()

	domain := c.Query("domain")
	if domain == "" {
		response.FailWithMessage("域名不能为空", c)
		return
	}

	bind, err := ollamaModelService.GetBindDomain(ctx, domain)
	if err != nil {
		global.GVA_LOG.Error("获取域名绑定信息失败!", zap.Error(err))
		response.FailWithMessage("获取域名绑定信息失败:"+err.Error(), c)
		return
	}

	response.OkWithData(bind, c)
}

// UpdateOllamaBindDomain 更新域名绑定
// @Tags OllamaModel
// @Summary 更新域名绑定
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param domain query string true "域名"
// @Param data body AIReq.OllamaBindDomainRequest true "更新参数"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ollamaModel/updateBindDomain [put]
func (ollamaModelApi *OllamaModelApi) UpdateOllamaBindDomain(c *gin.Context) {
	ctx := c.Request.Context()

	domain := c.Query("domain")
	if domain == "" {
		response.FailWithMessage("域名不能为空", c)
		return
	}

	var req AIReq.OllamaBindDomainRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	err = ollamaModelService.UpdateBindDomain(ctx, domain, req, userID)
	if err != nil {
		global.GVA_LOG.Error("更新域名绑定失败!", zap.Error(err))
		response.FailWithMessage("更新域名绑定失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("域名绑定更新成功", c)
}

// GetOllamaModelLogs 获取模型日志
// @Tags OllamaModel
// @Summary 获取模型日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "模型ID"
// @Param lines query int false "日志行数"
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /ollamaModel/logs [get]
func (ollamaModelApi *OllamaModelApi) GetOllamaModelLogs(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("模型ID不能为空", c)
		return
	}

	lines := 100 // 默认100行
	if linesStr := c.Query("lines"); linesStr != "" {
		if parsedLines, err := strconv.Atoi(linesStr); err == nil && parsedLines > 0 {
			lines = parsedLines
		}
	}

	logs, err := ollamaModelService.GetLogs(ctx, ID, lines)
	if err != nil {
		global.GVA_LOG.Error("获取模型日志失败!", zap.Error(err))
		response.FailWithMessage("获取模型日志失败:"+err.Error(), c)
		return
	}

	response.OkWithData(logs, c)
}

// GetSystemResourceStatus 获取系统资源状态
// @Tags OllamaModel
// @Summary 获取系统资源状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=AI.ResourceStatus,msg=string} "获取成功"
// @Router /ollamaModel/systemResource [get]
func (ollamaModelApi *OllamaModelApi) GetSystemResourceStatus(c *gin.Context) {
	ctx := c.Request.Context()

	status, err := ollamaModelService.GetSystemResourceStatus(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取系统资源状态失败!", zap.Error(err))
		response.FailWithMessage("获取系统资源状态失败:"+err.Error(), c)
		return
	}

	response.OkWithData(status, c)
}

// ChatWithOllamaModel 与模型对话
// @Tags OllamaModel
// @Summary 与模型对话
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body AIReq.OllamaChatRequest true "对话参数"
// @Success 200 {object} response.Response{data=AIReq.OllamaChatResponse,msg=string} "对话成功"
// @Router /ollamaModel/chat [post]
func (ollamaModelApi *OllamaModelApi) ChatWithOllamaModel(c *gin.Context) {
	ctx := c.Request.Context()

	var req AIReq.OllamaChatRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	chatResponse, err := ollamaModelService.Chat(ctx, req)
	if err != nil {
		global.GVA_LOG.Error("对话失败!", zap.Error(err))
		response.FailWithMessage("对话失败:"+err.Error(), c)
		return
	}

	response.OkWithData(chatResponse, c)
}
