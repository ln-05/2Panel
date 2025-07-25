package docker

import (
	"io"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerImageApi struct{}

// GetImageList 获取镜像列表
// @Tags Docker
// @Summary 获取Docker镜像列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dockerReq.ImageFilter true "分页参数"
// @Success 200 {object} response.Response{data=dockerRes.ImageListResponse,msg=string} "获取成功"
// @Router /docker/images [get]
func (d *DockerImageApi) GetImageList(c *gin.Context) {
	var pageInfo dockerReq.ImageFilter
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认分页参数
	if pageInfo.Page <= 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize <= 0 {
		pageInfo.PageSize = 10
	}

	// 调用服务层获取镜像列表
	images, total, err := dockerImageService.GetImageList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取镜像列表失败", zap.Error(err))
		response.FailWithMessage("获取镜像列表失败: "+err.Error(), c)
		return
	}

	// 构建响应
	responseData := dockerRes.ImageListResponse{
		List: images,
		PageResult: response.PageResult{
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
			Total:    total,
		},
	}

	response.OkWithDetailed(responseData, "获取成功", c)
}

// GetImageDetail 获取镜像详细信息
// @Tags Docker
// @Summary 获取Docker镜像详细信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "镜像ID"
// @Success 200 {object} response.Response{data=dockerRes.ImageDetail,msg=string} "获取成功"
// @Router /docker/images/{id} [get]
func (d *DockerImageApi) GetImageDetail(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		response.FailWithMessage("镜像ID不能为空", c)
		return
	}

	// 调用服务层获取镜像详细信息
	imageDetail, err := dockerImageService.GetImageDetail(imageID)
	if err != nil {
		global.GVA_LOG.Error("获取镜像详细信息失败", zap.String("imageID", imageID), zap.Error(err))
		if err.Error() == "image not found" {
			response.FailWithMessage("镜像不存在", c)
			return
		}
		response.FailWithMessage("获取镜像详细信息失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(*imageDetail, "获取成功", c)
}

// PullImage 拉取镜像
// @Tags Docker
// @Summary 拉取Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.ImagePullRequest true "拉取镜像参数"
// @Success 200 {object} response.Response{data=string,msg=string} "拉取成功"
// @Router /docker/images/pull [post]
func (d *DockerImageApi) PullImage(c *gin.Context) {
	var pullReq dockerReq.ImagePullRequest
	err := c.ShouldBindJSON(&pullReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层拉取镜像
	pullLog, err := dockerImageService.PullImage(pullReq)
	if err != nil {
		global.GVA_LOG.Error("拉取镜像失败", zap.String("image", pullReq.Image), zap.Error(err))
		response.FailWithMessage("拉取镜像失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(pullLog, "镜像拉取成功", c)
}

// RemoveImage 删除镜像
// @Tags Docker
// @Summary 删除Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "镜像ID"
// @Param force query bool false "是否强制删除"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /docker/images/{id} [delete]
func (d *DockerImageApi) RemoveImage(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		response.FailWithMessage("镜像ID不能为空", c)
		return
	}

	// 解析强制删除参数
	force := c.Query("force") == "true"

	// 调用服务层删除镜像
	err := dockerImageService.RemoveImage(imageID, force)
	if err != nil {
		global.GVA_LOG.Error("删除镜像失败", zap.String("imageID", imageID), zap.Error(err))
		if err.Error() == "image not found" {
			response.FailWithMessage("镜像不存在", c)
			return
		}
		response.FailWithMessage("删除镜像失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("镜像删除成功", c)
}

// TagImage 给镜像打标签
// @Tags Docker
// @Summary 给Docker镜像打标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.ImageTagRequest true "标签参数"
// @Success 200 {object} response.Response{msg=string} "标签成功"
// @Router /docker/images/tag [post]
func (d *DockerImageApi) TagImage(c *gin.Context) {
	var tagReq dockerReq.ImageTagRequest
	err := c.ShouldBindJSON(&tagReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层给镜像打标签
	err = dockerImageService.TagImage(tagReq)
	if err != nil {
		global.GVA_LOG.Error("镜像打标签失败",
			zap.String("source", tagReq.SourceImage),
			zap.String("target", tagReq.TargetImage),
			zap.Error(err))
		response.FailWithMessage("镜像打标签失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("镜像标签设置成功", c)
}

// PruneImages 清理未使用的镜像
// @Tags Docker
// @Summary 清理未使用的Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param dangling query bool false "是否只清理悬空镜像"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "清理成功"
// @Router /docker/images/prune [post]
func (d *DockerImageApi) PruneImages(c *gin.Context) {
	// 解析参数
	danglingStr := c.Query("dangling")
	dangling := danglingStr == "true"

	// 调用服务层清理镜像
	deletedCount, spaceReclaimed, err := dockerImageService.PruneImages(dangling)
	if err != nil {
		global.GVA_LOG.Error("清理镜像失败", zap.Error(err))
		response.FailWithMessage("清理镜像失败: "+err.Error(), c)
		return
	}

	// 构建响应数据
	responseData := map[string]interface{}{
		"deletedCount":   deletedCount,
		"spaceReclaimed": spaceReclaimed,
	}

	response.OkWithDetailed(responseData, "镜像清理成功", c)
}

// BuildImage 构建镜像
// @Tags Docker
// @Summary 构建Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.ImageBuildRequest true "构建镜像参数"
// @Success 200 {object} response.Response{data=string,msg=string} "构建成功"
// @Router /docker/images/build [post]
func (d *DockerImageApi) BuildImage(c *gin.Context) {
	var buildReq dockerReq.ImageBuildRequest
	err := c.ShouldBindJSON(&buildReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层构建镜像
	buildLog, err := dockerImageService.BuildImage(buildReq)
	if err != nil {
		global.GVA_LOG.Error("构建镜像失败", 
			zap.String("imageName", buildReq.ImageName), 
			zap.String("tag", buildReq.Tag), 
			zap.Error(err))
		response.FailWithMessage("构建镜像失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(buildLog, "镜像构建成功", c)
}

// ExportImage 导出镜像
// @Tags Docker
// @Summary 导出Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/octet-stream
// @Param data body dockerReq.ImageExportRequest true "导出镜像参数"
// @Success 200 {file} binary "导出成功"
// @Router /docker/images/export [post]
func (d *DockerImageApi) ExportImage(c *gin.Context) {
	var exportReq dockerReq.ImageExportRequest
	err := c.ShouldBindJSON(&exportReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层导出镜像
	reader, err := dockerImageService.ExportImage(exportReq)
	if err != nil {
		global.GVA_LOG.Error("导出镜像失败", zap.Strings("images", exportReq.Images), zap.Error(err))
		response.FailWithMessage("导出镜像失败: "+err.Error(), c)
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=images.tar")

	// 将镜像数据流式传输给客户端
	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		global.GVA_LOG.Error("传输导出数据失败", zap.Error(err))
		return
	}

	global.GVA_LOG.Info("镜像导出成功", zap.Strings("images", exportReq.Images))
}

// ImportImage 导入镜像
// @Tags Docker
// @Summary 导入Docker镜像
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.ImageImportRequest true "导入镜像参数"
// @Success 200 {object} response.Response{data=string,msg=string} "导入成功"
// @Router /docker/images/import [post]
func (d *DockerImageApi) ImportImage(c *gin.Context) {
	var importReq dockerReq.ImageImportRequest
	err := c.ShouldBindJSON(&importReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层导入镜像
	importLog, err := dockerImageService.ImportImage(importReq)
	if err != nil {
		global.GVA_LOG.Error("导入镜像失败", zap.String("source", importReq.Source), zap.Error(err))
		response.FailWithMessage("导入镜像失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(importLog, "镜像导入成功", c)
}
