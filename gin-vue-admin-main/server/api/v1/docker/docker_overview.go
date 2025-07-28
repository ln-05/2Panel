package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerOverviewApi struct{}

// GetOverviewStats 获取Docker概览统计信息
// @Tags DockerOverview
// @Summary 获取Docker概览统计信息
// @Description 获取Docker环境的整体统计信息，包括容器、镜像、网络、存储卷等
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.OverviewStats} "获取成功"
// @Router /docker/overview [get]
func (d *DockerOverviewApi) GetOverviewStats(c *gin.Context) {
	stats, err := dockerOverviewService.GetOverviewStats()
	if err != nil {
		global.GVA_LOG.Error("获取Docker概览统计失败", zap.Error(err))
		response.FailWithMessage("获取Docker概览统计失败: "+err.Error(), c)
		return
	}

	response.OkWithData(stats, c)
}

// GetConfigSummary 获取Docker配置摘要信息
// @Tags DockerOverview
// @Summary 获取Docker配置摘要信息
// @Description 获取Docker配置的摘要信息，包括Socket路径、镜像加速器等
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.ConfigSummary} "获取成功"
// @Router /docker/config/summary [get]
func (d *DockerOverviewApi) GetConfigSummary(c *gin.Context) {
	summary, err := dockerOverviewService.GetConfigSummary()
	if err != nil {
		global.GVA_LOG.Error("获取Docker配置摘要失败", zap.Error(err))
		response.FailWithMessage("获取Docker配置摘要失败: "+err.Error(), c)
		return
	}

	response.OkWithData(summary, c)
}

// GetDockerDiskUsage 获取Docker磁盘使用情况
// @Tags DockerOverview
// @Summary 获取Docker磁盘使用情况
// @Description 获取Docker环境的磁盘使用详情，包括镜像、容器、存储卷等占用空间
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.DiskUsage} "获取成功"
// @Router /docker/disk-usage [get]
func (d *DockerOverviewApi) GetDockerDiskUsage(c *gin.Context) {
	usage, err := dockerOverviewService.GetDockerDiskUsage()
	if err != nil {
		global.GVA_LOG.Error("获取Docker磁盘使用情况失败", zap.Error(err))
		response.FailWithMessage("获取Docker磁盘使用情况失败: "+err.Error(), c)
		return
	}

	response.OkWithData(usage, c)
}