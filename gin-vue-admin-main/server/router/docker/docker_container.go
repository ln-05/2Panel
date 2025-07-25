package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerContainerRouter struct{}

// InitDockerContainerRouter 初始化Docker容器路由
func (d *DockerContainerRouter) InitDockerContainerRouter(Router *gin.RouterGroup) {
	// 带操作记录的路由组 - 用于需要记录操作日志的API
	dockerRouter := Router.Group("docker").Use(middleware.OperationRecord())
	
	// 不带操作记录的路由组 - 用于查询类API
	dockerRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（容器操作）
	{
		dockerRouter.POST("containers/:id/start", dockerContainerApi.StartContainer)     // 启动容器
		dockerRouter.POST("containers/:id/stop", dockerContainerApi.StopContainer)       // 停止容器
		dockerRouter.POST("containers/:id/restart", dockerContainerApi.RestartContainer) // 重启容器
		dockerRouter.DELETE("containers/:id", dockerContainerApi.RemoveContainer)        // 删除容器
	}

	// 不需要记录操作的路由（查询类）
	{
		dockerRouterWithoutRecord.GET("containers", dockerContainerApi.GetContainerList)       // 获取容器列表
		dockerRouterWithoutRecord.GET("containers/:id", dockerContainerApi.GetContainerDetail) // 获取容器详情
		dockerRouterWithoutRecord.GET("containers/:id/logs", dockerContainerApi.GetContainerLogs) // 获取容器日志
		dockerRouterWithoutRecord.GET("info", dockerContainerApi.GetDockerInfo)                // 获取Docker信息
		dockerRouterWithoutRecord.GET("status", dockerContainerApi.CheckDockerStatus)          // 检查Docker状态
	}
}