package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerImageRouter struct{}

// InitDockerImageRouter 初始化Docker镜像路由
func (d *DockerImageRouter) InitDockerImageRouter(Router *gin.RouterGroup) {
	// 带操作记录的路由组 - 用于需要记录操作日志的API
	dockerRouter := Router.Group("docker").Use(middleware.OperationRecord())
	
	// 不带操作记录的路由组 - 用于查询类API
	dockerRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（镜像操作）
	{
		dockerRouter.POST("images/pull", dockerImageApi.PullImage)       // 拉取镜像
		dockerRouter.DELETE("images/:id", dockerImageApi.RemoveImage)    // 删除镜像
		dockerRouter.POST("images/tag", dockerImageApi.TagImage)         // 给镜像打标签
		dockerRouter.POST("images/prune", dockerImageApi.PruneImages)    // 清理未使用的镜像
		dockerRouter.POST("images/build", dockerImageApi.BuildImage)     // 构建镜像
		dockerRouter.POST("images/export", dockerImageApi.ExportImage)   // 导出镜像
		dockerRouter.POST("images/import", dockerImageApi.ImportImage)   // 导入镜像
	}

	// 不需要记录操作的路由（查询类）
	{
		dockerRouterWithoutRecord.GET("images", dockerImageApi.GetImageList)       // 获取镜像列表
		dockerRouterWithoutRecord.GET("images/:id", dockerImageApi.GetImageDetail) // 获取镜像详情
	}
}