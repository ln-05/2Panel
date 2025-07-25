package router

import (
	"api/handler"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	api := r.Group("/api")
	{
		database := api.Group("/database")
		{
			database.POST("/create", handler.DatabaseCreate)  // 创建数据库连接
			database.GET("/list", handler.DatabaseList)       // 获取数据库连接列表
			database.GET("/id", handler.DatabaseGet)          // 获取单个数据库连接
			database.POST("/update", handler.DatabaseUpdate)  // 更新数据库连接
			database.POST("/deleted", handler.DatabaseDelete) // 删除数据库连接
			database.POST("/test", handler.DatabaseTest)      // 测试数据库连接
		database.POST("/sync", handler.DatabaseSync)      // 从服务器同步数据库
		}
	}
}
