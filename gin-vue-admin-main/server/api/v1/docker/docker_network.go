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

type DockerNetworkApi struct{}

// GetNetworkList 获取网络列表
// @Tags Docker网络管理
// @Summary 获取网络列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dockerReq.NetworkFilter true "分页参数"
// @Success 200 {object} response.Response{data=dockerRes.NetworkListResponse,msg=string} "获取成功"
// @Router /docker/networks [get]
func (d *DockerNetworkApi) GetNetworkList(c *gin.Context) {
	var pageInfo dockerReq.NetworkFilter
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

	// 调用服务层获取网络列表
	networks, total, err := dockerNetworkService.GetNetworkList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取网络列表失败", zap.Error(err))
		response.FailWithMessage("获取网络列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(dockerRes.NetworkListResponse{
		List:  networks,
		Total: total,
	}, "获取成功", c)
}

// GetNetworkDetail 获取网络详细信息
// @Tags Docker网络管理
// @Summary 获取网络详细信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "网络ID"
// @Success 200 {object} response.Response{data=dockerRes.NetworkDetail,msg=string} "获取成功"
// @Router /docker/networks/{id} [get]
func (d *DockerNetworkApi) GetNetworkDetail(c *gin.Context) {
	networkID := c.Param("id")
	if networkID == "" {
		response.FailWithMessage("网络ID不能为空", c)
		return
	}

	// 调用服务层获取网络详细信息
	networkDetail, err := dockerNetworkService.GetNetworkDetail(networkID)
	if err != nil {
		global.GVA_LOG.Error("获取网络详细信息失败", zap.String("networkID", networkID), zap.Error(err))
		response.FailWithMessage("获取网络详细信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(networkDetail, c)
}

// CreateNetwork 创建网络
// @Tags Docker网络管理
// @Summary 创建网络
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerReq.NetworkCreateRequest true "创建网络参数"
// @Success 200 {object} response.Response{data=string,msg=string} "创建成功"
// @Router /docker/networks [post]
func (d *DockerNetworkApi) CreateNetwork(c *gin.Context) {
	var createReq dockerReq.NetworkCreateRequest
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认驱动
	if createReq.Driver == "" {
		createReq.Driver = "bridge"
	}

	// 调用服务层创建网络
	networkID, err := dockerNetworkService.CreateNetwork(createReq)
	if err != nil {
		global.GVA_LOG.Error("创建网络失败", zap.String("name", createReq.Name), zap.Error(err))
		response.FailWithMessage("创建网络失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"networkId": networkID,
	}, "创建成功", c)
}

// RemoveNetwork 删除网络
// @Tags Docker网络管理
// @Summary 删除网络
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "网络ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /docker/networks/{id} [delete]
func (d *DockerNetworkApi) RemoveNetwork(c *gin.Context) {
	networkID := c.Param("id")
	if networkID == "" {
		response.FailWithMessage("网络ID不能为空", c)
		return
	}

	// 调用服务层删除网络
	err := dockerNetworkService.RemoveNetwork(networkID)
	if err != nil {
		global.GVA_LOG.Error("删除网络失败", zap.String("networkID", networkID), zap.Error(err))
		response.FailWithMessage("删除网络失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// PruneNetworks 清理未使用的网络
// @Tags Docker网络管理
// @Summary 清理未使用的网络
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "清理成功"
// @Router /docker/networks/prune [post]
func (d *DockerNetworkApi) PruneNetworks(c *gin.Context) {
	// 调用服务层清理网络
	deletedCount, spaceReclaimed, err := dockerNetworkService.PruneNetworks()
	if err != nil {
		global.GVA_LOG.Error("清理网络失败", zap.Error(err))
		response.FailWithMessage("清理网络失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{
		"deletedCount":   deletedCount,
		"spaceReclaimed": spaceReclaimed,
	}, "清理完成，删除了 "+strconv.FormatInt(deletedCount, 10)+" 个网络", c)
}
