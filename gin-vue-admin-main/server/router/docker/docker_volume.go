package docker

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerVolumeRouter struct{}

// InitDockerVolumeRouter 初始化Docker存储卷路由
func (d *DockerVolumeRouter) InitDockerVolumeRouter(Router *gin.RouterGroup) {
	dockerVolumeApi := api.DockerVolumeApi{}
	
	// 带操作记录的路由组 - 用于需要记录操作日志的API
	volumeRouter := Router.Group("docker").Use(middleware.OperationRecord())
	// 不带操作记录的路由组 - 用于查询类API
	volumeRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（存储卷操作）
	{
		volumeRouter.POST("volumes", dockerVolumeApi.CreateVolume)           // 创建存储卷
		volumeRouter.DELETE("volumes/:name", dockerVolumeApi.RemoveVolume)   // 删除存储卷
		volumeRouter.POST("volumes/prune", dockerVolumeApi.PruneVolumes)     // 清理未使用的存储卷
	}

	// 不需要记录操作的路由（查询类）
	{
		volumeRouterWithoutRecord.GET("volumes", dockerVolumeApi.GetVolumeList)       // 获取存储卷列表
		volumeRouterWithoutRecord.GET("volumes/:name", dockerVolumeApi.GetVolumeDetail) // 获取存储卷详情
	}
}