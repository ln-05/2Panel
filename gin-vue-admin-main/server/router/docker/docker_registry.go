package docker

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerRegistryRouter struct{}

// InitDockerRegistryRouter 初始化Docker仓库路由
func (d *DockerRegistryRouter) InitDockerRegistryRouter(Router *gin.RouterGroup) {
	dockerRegistryApi := api.DockerRegistryApi{}
	
	// 带操作记录的路由组 - 用于需要记录操作日志的API
	registryRouter := Router.Group("docker").Use(middleware.OperationRecord())
	// 不带操作记录的路由组 - 用于查询类API
	registryRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（仓库操作）
	{
		registryRouter.POST("registries", dockerRegistryApi.CreateRegistry)                    // 创建仓库
		registryRouter.PUT("registries", dockerRegistryApi.UpdateRegistry)                    // 更新仓库
		registryRouter.DELETE("registries/:id", dockerRegistryApi.DeleteRegistry)             // 删除仓库
		registryRouter.POST("registries/:id/test", dockerRegistryApi.TestRegistry)            // 测试仓库连接
		registryRouter.POST("registries/:id/default", dockerRegistryApi.SetDefaultRegistry)   // 设置默认仓库
		registryRouter.POST("registries/sync", dockerRegistryApi.SyncFrom1Panel)              // 从1Panel同步仓库数据
	}

	// 不需要记录操作的路由（查询类）
	{
		registryRouterWithoutRecord.GET("registries", dockerRegistryApi.GetRegistryList)       // 获取仓库列表
		registryRouterWithoutRecord.GET("registries/:id", dockerRegistryApi.GetRegistryDetail) // 获取仓库详情
	}
}