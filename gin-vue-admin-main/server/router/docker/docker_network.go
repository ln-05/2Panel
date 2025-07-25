package docker

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerNetworkRouter struct{}

// InitDockerNetworkRouter 初始化Docker网络路由
func (d *DockerNetworkRouter) InitDockerNetworkRouter(Router *gin.RouterGroup) {
	dockerNetworkApi := api.DockerNetworkApi{}
	
	// 带操作记录的路由组 - 用于需要记录操作日志的API
	networkRouter := Router.Group("docker").Use(middleware.OperationRecord())
	// 不带操作记录的路由组 - 用于查询类API
	networkRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（网络操作）
	{
		networkRouter.POST("networks", dockerNetworkApi.CreateNetwork)           // 创建网络
		networkRouter.DELETE("networks/:id", dockerNetworkApi.RemoveNetwork)    // 删除网络
		networkRouter.POST("networks/prune", dockerNetworkApi.PruneNetworks)    // 清理未使用的网络
	}

	// 不需要记录操作的路由（查询类）
	{
		networkRouterWithoutRecord.GET("networks", dockerNetworkApi.GetNetworkList)       // 获取网络列表
		networkRouterWithoutRecord.GET("networks/:id", dockerNetworkApi.GetNetworkDetail) // 获取网络详情
	}
}
