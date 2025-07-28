package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerConfigRouter struct{}

// InitDockerConfigRouter 初始化Docker配置路由
func (d *DockerConfigRouter) InitDockerConfigRouter(Router *gin.RouterGroup) {
	dockerConfigApi := v1.ApiGroupApp.DockerApiGroup.DockerConfigApi

	// 带操作记录的路由组 - 用于需要记录操作日志的API
	dockerRouter := Router.Group("docker").Use(middleware.OperationRecord())

	// 不带操作记录的路由组 - 用于查询类API
	dockerRouterWithoutRecord := Router.Group("docker")

	// 需要记录操作的路由（配置修改操作）
	{
		dockerRouter.PUT("config", dockerConfigApi.UpdateDockerConfig)                 // 更新Docker配置
		dockerRouter.POST("config/backup", dockerConfigApi.BackupDockerConfig)         // 备份Docker配置
		dockerRouter.POST("config/restore", dockerConfigApi.RestoreDockerConfig)       // 恢复Docker配置
		dockerRouter.DELETE("config/backups/:backupId", dockerConfigApi.DeleteBackup)  // 删除备份
		dockerRouter.POST("config/backups/cleanup", dockerConfigApi.CleanupOldBackups) // 清理过期备份
		dockerRouter.POST("service/restart", dockerConfigApi.RestartDockerService)     // 重启Docker服务
		dockerRouter.POST("service/start", dockerConfigApi.StartDockerService)         // 启动Docker服务
		dockerRouter.POST("service/stop", dockerConfigApi.StopDockerService)           // 停止Docker服务
	}

	// 不需要记录操作的路由（查询类）
	{
		dockerRouterWithoutRecord.GET("config", dockerConfigApi.GetDockerConfig)                  // 获取Docker配置
		dockerRouterWithoutRecord.POST("config/validate", dockerConfigApi.ValidateDockerConfig)   // 验证Docker配置
		dockerRouterWithoutRecord.GET("config/backups", dockerConfigApi.GetBackupList)            // 获取备份列表
		dockerRouterWithoutRecord.GET("service/status", dockerConfigApi.GetDockerServiceStatus)   // 获取Docker服务状态
		dockerRouterWithoutRecord.GET("service/health", dockerConfigApi.CheckDockerServiceHealth) // 检查Docker服务健康状态
	}
}
