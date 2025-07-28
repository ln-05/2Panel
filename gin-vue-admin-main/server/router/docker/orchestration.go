package docker

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/gin-gonic/gin"
)

// DockerRouter 定义Docker路由组
type DockerRouter struct{}

// dockerApi 实例化Docker API
var dockerApi = &api.DockerContainerApi{}

// InitDockerOrchestrationRouter 初始化编排路由
func (d *DockerRouter) InitDockerOrchestrationRouter(Router *gin.RouterGroup) {
	orchestrationRouter := Router.Group("orchestration")
	{
		orchestrationRouter.GET("/list", dockerApi.GetOrchestrationList)
	}
}
