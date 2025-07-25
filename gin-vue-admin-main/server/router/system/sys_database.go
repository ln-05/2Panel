package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysDatabaseRouter struct{}

func (s *SysDatabaseRouter) InitDatabaseRouter(Router *gin.RouterGroup) {
	databaseRouter := Router.Group("database").Use(middleware.OperationRecord())
	databaseRouterWithoutRecord := Router.Group("database")
	var databaseApi = v1.ApiGroupApp.SystemApiGroup.DatabaseApi
	{
		databaseRouter.POST("create", databaseApi.CreateSysDatabase)   // 新建数据库连接
		databaseRouter.POST("deleted", databaseApi.DeleteSysDatabase) // 删除数据库连接
		databaseRouter.POST("update", databaseApi.UpdateSysDatabase)  // 更新数据库连接
		databaseRouter.POST("test", databaseApi.TestSysDatabase)      // 测试数据库连接
		databaseRouter.POST("sync", databaseApi.SyncSysDatabase)      // 同步数据库配置
	}
	{
		databaseRouterWithoutRecord.GET("id", databaseApi.FindSysDatabase)        // 根据ID获取数据库连接
		databaseRouterWithoutRecord.GET("list", databaseApi.GetSysDatabaseList)   // 获取数据库连接列表
	}
}
