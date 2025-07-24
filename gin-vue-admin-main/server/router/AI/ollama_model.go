package AI

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type OllamaModelRouter struct{}

// InitOllamaModelRouter 初始化 ollamaModel表 路由信息
func (s *OllamaModelRouter) InitOllamaModelRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	ollamaModelRouter := Router.Group("ollamaModel").Use(middleware.OperationRecord())
	ollamaModelRouterWithoutRecord := Router.Group("ollamaModel")
	ollamaModelRouterWithoutAuth := PublicRouter.Group("ollamaModel")

	// 需要记录操作日志的接口
	{
		ollamaModelRouter.POST("createOllamaModel", ollamaModelApi.CreateOllamaModel)             // 新建ollamaModel表
		ollamaModelRouter.POST("create", ollamaModelApi.CreateOllamaModelAdvanced)                // 创建/下载新模型
		ollamaModelRouter.POST("start", ollamaModelApi.StartOllamaModel)                          // 启动模型
		ollamaModelRouter.POST("stop", ollamaModelApi.StopOllamaModel)                            // 停止模型
		ollamaModelRouter.POST("close", ollamaModelApi.CloseOllamaModel)                          // 停止模型(兼容)
		ollamaModelRouter.POST("recreate", ollamaModelApi.RecreateOllamaModel)                    // 重新创建模型
		ollamaModelRouter.POST("sync", ollamaModelApi.SyncOllamaModel)                            // 同步模型
		ollamaModelRouter.DELETE("deleteOllamaModel", ollamaModelApi.DeleteOllamaModel)           // 删除ollamaModel表
		ollamaModelRouter.DELETE("deleteOllamaModelByIds", ollamaModelApi.DeleteOllamaModelByIds) // 批量删除ollamaModel表
		ollamaModelRouter.PUT("updateOllamaModel", ollamaModelApi.UpdateOllamaModel)              // 更新ollamaModel表

		// 域名绑定相关
		ollamaModelRouter.POST("bindDomain", ollamaModelApi.BindDomainToOllama)          // 绑定域名
		ollamaModelRouter.PUT("updateBindDomain", ollamaModelApi.UpdateOllamaBindDomain) // 更新域名绑定
	}

	// 不需要记录操作日志的接口
	{
		ollamaModelRouterWithoutRecord.GET("search", ollamaModelApi.SearchOllamaModel)               // 搜索模型
		ollamaModelRouterWithoutRecord.GET("detail", ollamaModelApi.LoadOllamaModelDetail)           // 获取模型详情
		ollamaModelRouterWithoutRecord.GET("findOllamaModel", ollamaModelApi.FindOllamaModel)        // 根据ID获取ollamaModel表
		ollamaModelRouterWithoutRecord.GET("getOllamaModelList", ollamaModelApi.GetOllamaModelList)  // 获取ollamaModel表列表
		ollamaModelRouterWithoutRecord.GET("getBindDomain", ollamaModelApi.GetOllamaBindDomain)      // 获取域名绑定信息
		ollamaModelRouterWithoutRecord.GET("logs", ollamaModelApi.GetOllamaModelLogs)                // 获取模型日志
		ollamaModelRouterWithoutRecord.GET("systemResource", ollamaModelApi.GetSystemResourceStatus) // 获取系统资源状态
		ollamaModelRouterWithoutRecord.POST("chat", ollamaModelApi.ChatWithOllamaModel)              // 与模型对话
	}

	// 公开接口
	{
		ollamaModelRouterWithoutAuth.GET("getOllamaModelPublic", ollamaModelApi.GetOllamaModelPublic) // ollamaModel表开放接口
	}
}
