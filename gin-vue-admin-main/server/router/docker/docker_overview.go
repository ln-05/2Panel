package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type DockerOverviewRouter struct{}

// InitDockerOverviewRouter 初始化Docker概览路由
func (d *DockerOverviewRouter) InitDockerOverviewRouter(Router *gin.RouterGroup) {
	dockerOverviewApi := v1.ApiGroupApp.DockerApiGroup.DockerOverviewApi

	// 不带操作记录的路由组 - 用于查询类API
	dockerRouterWithoutRecord := Router.Group("docker")

	// 概览相关路由（查询类，不需要记录操作）
	{
		dockerRouterWithoutRecord.GET("overview", dockerOverviewApi.GetOverviewStats)     // 获取Docker概览统计
		dockerRouterWithoutRecord.GET("config/summary", dockerOverviewApi.GetConfigSummary) // 获取配置摘要
		dockerRouterWithoutRecord.GET("disk-usage", dockerOverviewApi.GetDockerDiskUsage) // 获取磁盘使用情况
	}
}

// InitDockerOverviewPublicRouter 初始化Docker概览公开路由（用于测试）
func (d *DockerOverviewRouter) InitDockerOverviewPublicRouter(Router *gin.RouterGroup) {
	dockerOverviewApi := v1.ApiGroupApp.DockerApiGroup.DockerOverviewApi
	dockerRegistryApi := v1.ApiGroupApp.DockerApiGroup.DockerRegistryApi

	// 公开路由组 - 不需要权限验证（仅用于测试）
	dockerPublicRouter := Router.Group("docker")

	// 概览相关路由（公开访问，仅用于测试）
	{
		dockerPublicRouter.GET("overview", dockerOverviewApi.GetOverviewStats)     // 获取Docker概览统计
		dockerPublicRouter.GET("config/summary", dockerOverviewApi.GetConfigSummary) // 获取配置摘要
		dockerPublicRouter.GET("disk-usage", dockerOverviewApi.GetDockerDiskUsage) // 获取磁盘使用情况
		dockerPublicRouter.GET("registries-public", dockerRegistryApi.GetRegistryList)    // 获取仓库列表（公开访问，使用不同路径避免冲突）
	}
}