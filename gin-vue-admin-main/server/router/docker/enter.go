package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	DockerContainerRouter
	DockerImageRouter
	DockerNetworkRouter
	DockerVolumeRouter
	DockerRegistryRouter
	DockerRouter
	DockerConfigRouter
	DockerOverviewRouter
	DockerDiagnosticRouter
}

// 适配 initialize/router.go 的调用，转发到 DockerRouter 的实现
func (g *RouterGroup) InitDockerOrchestrationRouter(Router *gin.RouterGroup) {
	g.DockerRouter.InitDockerOrchestrationRouter(Router)
}

var (
	dockerContainerApi = docker.DockerContainerApi{}
)

func InitDockerRouter(Router *gin.RouterGroup) {
	// 编排相关API
	Router.GET("/docker/orchestrations", dockerContainerApi.GetOrchestrationList)
	Router.GET("/docker/orchestrations/:name", dockerContainerApi.OrchestrationDetail)
	Router.POST("/docker/orchestrations", dockerContainerApi.CreateOrchestration)
	Router.PUT("/docker/orchestrations/:name", dockerContainerApi.EditOrchestration)
	Router.DELETE("/docker/orchestrations/:name", dockerContainerApi.DeleteOrchestration)
	Router.GET("/docker/orchestrations/:name/status", dockerContainerApi.OrchestrationStatus)
	// 批量操作接口保留
	Router.POST("/docker/orchestrations/:name/:op", dockerContainerApi.BatchOperateOrchestration)
}
