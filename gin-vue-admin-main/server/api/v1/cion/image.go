package cion

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cion"
	cionReq "github.com/flipped-aurora/gin-vue-admin/server/model/cion/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct{}

// CreateImage 创建image表
// @Tags Image
// @Summary 创建image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cion.Image true "创建image表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /image/createImage [post]
func (imageApi *ImageApi) CreateImage(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var image cion.Image
	err := c.ShouldBindJSON(&image)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	image.CreatedBy = int(utils.GetUserID(c))
	err = imageService.CreateImage(ctx, &image)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteImage 删除image表
// @Tags Image
// @Summary 删除image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cion.Image true "删除image表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /image/deleteImage [delete]
func (imageApi *ImageApi) DeleteImage(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	err := imageService.DeleteImage(ctx, ID, userID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteImageByIds 批量删除image表
// @Tags Image
// @Summary 批量删除image表

// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /image/deleteImageByIds [delete]
func (imageApi *ImageApi) DeleteImageByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	err := imageService.DeleteImageByIds(ctx, IDs, userID)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateImage 更新image表
// @Tags Image
// @Summary 更新image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cion.Image true "更新image表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /image/updateImage [put]
func (imageApi *ImageApi) UpdateImage(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var image cion.Image
	err := c.ShouldBindJSON(&image)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	image.UpdatedBy = int(utils.GetUserID(c))
	err = imageService.UpdateImage(ctx, image)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindImage 用id查询image表
// @Tags Image
// @Summary 用id查询image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询image表"
// @Success 200 {object} response.Response{data=cion.Image,msg=string} "查询成功"
// @Router /image/findImage [get]
func (imageApi *ImageApi) FindImage(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reimage, err := imageService.GetImage(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reimage, c)
}

// GetImageList 分页获取image表列表
// @Tags Image
// @Summary 分页获取image表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cionReq.ImageSearch true "分页获取image表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /image/getImageList [get]
func (imageApi *ImageApi) GetImageList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo cionReq.ImageSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := imageService.GetImageInfoList(ctx, pageInfo)
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

// GetImagePublic 不需要鉴权的image表接口
// @Tags Image
// @Summary 不需要鉴权的image表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /image/getImagePublic [get]
func (imageApi *ImageApi) GetImagePublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	imageService.GetImagePublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的image表接口信息",
	}, "获取成功", c)
}
