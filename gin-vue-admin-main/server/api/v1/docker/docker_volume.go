package docker

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"go.uber.org/zap"
)

type DockerVolumeApi struct{}

// GetVolumeList 获取存储卷列表
// @Tags Docker存储卷管理
// @Summary 获取存储卷列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dockerReq.VolumeFilter true "分页参数"
// @Success 200 {object} response.Response{data=dockerRes.VolumeListResponse,msg=string} "获取成功"
// @Router /docker/volumes [get]
func (d *DockerVolumeApi) GetVolumeList(c *gin.Context) {
	var pageInfo dockerReq.VolumeFilter
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

	// 调用服务层获取存储卷列表
	volumes, total, err := dockerVolumeService.GetVolumeList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取存储卷列表失败", zap.Error(err))
		response.FailWithMessage("获取存储卷列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(dockerRes.VolumeListResponse{
		List:  volumes,
		Total: total,
	}, "获取成功", c)
}

// GetVolumeDetail 获取存储卷详细信息
// @Tags Docker存储卷管理
// @Summary 获取存储卷详细信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param name path string true "存储卷名称"
// @Success 200 {object} response.Response{data=dockerRes.VolumeDetail,msg=string} "获取成功"
// @Router /docker/volumes/{name} [get]
func (d *DockerVolumeApi) GetVolumeDetail(c *gin.Context) {
	volumeName := c.Param("name")
	if volumeName == "" {
		response.FailWithMessage("存储卷名称不能为空", c)
		return
	}

	// 调用服务层获取存储卷详细信息
	volumeDetail, err := dockerVolumeService.GetVolumeDetail(volumeName)
	if err != nil {
		global.GVA_LOG.Error("获取存储卷详细信息失败", zap.String("volumeName", volumeName), zap.Error(err))
		response.FailWithMessage("获取存储卷详细信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(volumeDetail, c)
}

// CreateVolume 创建存储卷
// @Tags Docker存储卷管理
// @Summary 创建存储卷
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.VolumeCreateRequest true "创建存储卷参数"
// @Success 200 {object} response.Response{data=string,msg=string} "创建成功"
// @Router /docker/volumes [post]
func (d *DockerVolumeApi) CreateVolume(c *gin.Context) {
	var createReq dockerReq.VolumeCreateRequest
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认驱动
	if createReq.Driver == "" {
		createReq.Driver = "local"
	}

	// 调用服务层创建存储卷
	volumeName, err := dockerVolumeService.CreateVolume(createReq)
	if err != nil {
		global.GVA_LOG.Error("创建存储卷失败", zap.String("name", createReq.Name), zap.Error(err))
		response.FailWithMessage("创建存储卷失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"volumeName": volumeName,
	}, "创建成功", c)
}

// RemoveVolume 删除存储卷
// @Tags Docker存储卷管理
// @Summary 删除存储卷
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param name path string true "存储卷名称"
// @Param force query bool false "是否强制删除"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /docker/volumes/{name} [delete]
func (d *DockerVolumeApi) RemoveVolume(c *gin.Context) {
	volumeName := c.Param("name")
	if volumeName == "" {
		response.FailWithMessage("存储卷名称不能为空", c)
		return
	}

	// 解析force参数
	forceStr := c.Query("force")
	force := forceStr == "true"

	// 调用服务层删除存储卷
	err := dockerVolumeService.RemoveVolume(volumeName, force)
	if err != nil {
		global.GVA_LOG.Error("删除存储卷失败", zap.String("volumeName", volumeName), zap.Error(err))
		response.FailWithMessage("删除存储卷失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// PruneVolumes 清理未使用的存储卷
// @Tags Docker存储卷管理
// @Summary 清理未使用的存储卷
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "清理成功"
// @Router /docker/volumes/prune [post]
func (d *DockerVolumeApi) PruneVolumes(c *gin.Context) {
	// 调用服务层清理存储卷
	deletedCount, spaceReclaimed, err := dockerVolumeService.PruneVolumes()
	if err != nil {
		global.GVA_LOG.Error("清理存储卷失败", zap.Error(err))
		response.FailWithMessage("清理存储卷失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"deletedCount":   deletedCount,
		"spaceReclaimed": spaceReclaimed,
	}, "清理完成，删除了 "+strconv.FormatInt(deletedCount, 10)+" 个存储卷，释放了 "+strconv.FormatInt(spaceReclaimed, 10)+" 字节空间", c)
}