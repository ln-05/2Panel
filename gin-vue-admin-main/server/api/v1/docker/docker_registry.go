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

type DockerRegistryApi struct{}

// GetRegistryList 获取仓库列表
// @Tags Docker仓库管理
// @Summary 获取仓库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dockerReq.RegistryFilter true "分页参数"
// @Success 200 {object} response.Response{data=dockerRes.RegistryListResponse,msg=string} "获取成功"
// @Router /docker/registries [get]
func (d *DockerRegistryApi) GetRegistryList(c *gin.Context) {
	var pageInfo dockerReq.RegistryFilter
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

	// 调用服务层获取仓库列表
	registries, total, err := dockerRegistryService.GetRegistryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取仓库列表失败", zap.Error(err))
		response.FailWithMessage("获取仓库列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(dockerRes.RegistryListResponse{
		List:  registries,
		Total: total,
	}, "获取成功", c)
}

// GetRegistryDetail 获取仓库详细信息
// @Tags Docker仓库管理
// @Summary 获取仓库详细信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "仓库ID"
// @Success 200 {object} response.Response{data=dockerRes.RegistryDetail,msg=string} "获取成功"
// @Router /docker/registries/{id} [get]
func (d *DockerRegistryApi) GetRegistryDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的仓库ID", c)
		return
	}

	// 调用服务层获取仓库详细信息
	registryDetail, err := dockerRegistryService.GetRegistryDetail(uint(id))
	if err != nil {
		global.GVA_LOG.Error("获取仓库详细信息失败", zap.Uint64("id", id), zap.Error(err))
		response.FailWithMessage("获取仓库详细信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(registryDetail, c)
}

// CreateRegistry 创建仓库
// @Tags Docker仓库管理
// @Summary 创建仓库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.RegistryCreateRequest true "创建仓库参数"
// @Success 200 {object} response.Response{data=dockerRes.RegistryInfo,msg=string} "创建成功"
// @Router /docker/registries [post]
func (d *DockerRegistryApi) CreateRegistry(c *gin.Context) {
	var createReq dockerReq.RegistryCreateRequest
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层创建仓库
	registryInfo, err := dockerRegistryService.CreateRegistry(createReq)
	if err != nil {
		global.GVA_LOG.Error("创建仓库失败", zap.String("name", createReq.Name), zap.Error(err))
		response.FailWithMessage("创建仓库失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(registryInfo, "创建成功", c)
}

// UpdateRegistry 更新仓库
// @Tags Docker仓库管理
// @Summary 更新仓库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.RegistryUpdateRequest true "更新仓库参数"
// @Success 200 {object} response.Response{data=dockerRes.RegistryInfo,msg=string} "更新成功"
// @Router /docker/registries [put]
func (d *DockerRegistryApi) UpdateRegistry(c *gin.Context) {
	var updateReq dockerReq.RegistryUpdateRequest
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层更新仓库
	registryInfo, err := dockerRegistryService.UpdateRegistry(updateReq)
	if err != nil {
		global.GVA_LOG.Error("更新仓库失败", zap.Uint("id", updateReq.ID), zap.Error(err))
		response.FailWithMessage("更新仓库失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(registryInfo, "更新成功", c)
}

// DeleteRegistry 删除仓库
// @Tags Docker仓库管理
// @Summary 删除仓库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "仓库ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /docker/registries/{id} [delete]
func (d *DockerRegistryApi) DeleteRegistry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的仓库ID", c)
		return
	}

	// 调用服务层删除仓库
	err = dockerRegistryService.DeleteRegistry(uint(id))
	if err != nil {
		global.GVA_LOG.Error("删除仓库失败", zap.Uint64("id", id), zap.Error(err))
		response.FailWithMessage("删除仓库失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// TestRegistry 测试仓库连接
// @Tags Docker仓库管理
// @Summary 测试仓库连接
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "仓库ID"
// @Success 200 {object} response.Response{data=dockerRes.RegistryTestResponse,msg=string} "测试完成"
// @Router /docker/registries/{id}/test [post]
func (d *DockerRegistryApi) TestRegistry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的仓库ID", c)
		return
	}

	// 调用服务层测试仓库连接
	testResult, err := dockerRegistryService.TestRegistry(uint(id))
	if err != nil {
		global.GVA_LOG.Error("测试仓库连接失败", zap.Uint64("id", id), zap.Error(err))
		response.FailWithMessage("测试仓库连接失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(testResult, "测试完成", c)
}

// SetDefaultRegistry 设置默认仓库
// @Tags Docker仓库管理
// @Summary 设置默认仓库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path int true "仓库ID"
// @Success 200 {object} response.Response{msg=string} "设置成功"
// @Router /docker/registries/{id}/default [post]
func (d *DockerRegistryApi) SetDefaultRegistry(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的仓库ID", c)
		return
	}

	// 调用服务层设置默认仓库
	err = dockerRegistryService.SetDefaultRegistry(uint(id))
	if err != nil {
		global.GVA_LOG.Error("设置默认仓库失败", zap.Uint64("id", id), zap.Error(err))
		response.FailWithMessage("设置默认仓库失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("设置成功", c)
}

// SyncFrom1Panel 从1Panel同步仓库数据
// @Tags Docker仓库管理
// @Summary 从1Panel同步仓库数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "同步成功"
// @Router /docker/registries/sync [post]
func (d *DockerRegistryApi) SyncFrom1Panel(c *gin.Context) {
	// 调用服务层同步数据
	err := registrySyncService.SyncFrom1Panel()
	if err != nil {
		global.GVA_LOG.Error("从1Panel同步仓库数据失败", zap.Error(err))
		response.FailWithMessage("同步失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("同步成功", c)
}