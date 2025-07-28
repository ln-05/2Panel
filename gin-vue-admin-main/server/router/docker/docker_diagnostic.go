package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type DockerDiagnosticRouter struct{}

// InitDockerDiagnosticRouter 初始化Docker诊断路由
func (d *DockerDiagnosticRouter) InitDockerDiagnosticRouter(Router *gin.RouterGroup) {
	dockerDiagnosticApi := v1.ApiGroupApp.DockerApiGroup.DockerDiagnosticApi

	// 不带操作记录的路由组 - 用于查询类API
	dockerRouterWithoutRecord := Router.Group("docker")

	// 诊断相关路由（查询类，不需要记录操作）
	{
		dockerRouterWithoutRecord.GET("diagnose", dockerDiagnosticApi.DiagnoseConnection)     // 全面诊断
		dockerRouterWithoutRecord.POST("test-connection", dockerDiagnosticApi.TestConnection) // 快速连接测试
		dockerRouterWithoutRecord.GET("permissions", dockerDiagnosticApi.GetPermissionStatus) // 权限状态
	}
}

// InitDockerDiagnosticPublicRouter 初始化Docker诊断公开路由（用于测试）
func (d *DockerDiagnosticRouter) InitDockerDiagnosticPublicRouter(Router *gin.RouterGroup) {
	dockerDiagnosticApi := v1.ApiGroupApp.DockerApiGroup.DockerDiagnosticApi

	// 公开路由组 - 不需要权限验证（仅用于测试）
	dockerPublicRouter := Router.Group("docker")

	// 诊断相关路由（公开访问，仅用于测试）
	{
		dockerPublicRouter.GET("diagnose", dockerDiagnosticApi.DiagnoseConnection)     // 全面诊断
		dockerPublicRouter.POST("test-connection", dockerDiagnosticApi.TestConnection) // 快速连接测试
		dockerPublicRouter.GET("permissions", dockerDiagnosticApi.GetPermissionStatus) // 权限状态
	}
}