package AI

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type McpServerRouter struct {}

// InitMcpServerRouter 初始化 mcpServer表 路由信息
func (s *McpServerRouter) InitMcpServerRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	mcpServerRouter := Router.Group("mcpServer").Use(middleware.OperationRecord())
	mcpServerRouterWithoutRecord := Router.Group("mcpServer")
	mcpServerRouterWithoutAuth := PublicRouter.Group("mcpServer")
	{
		mcpServerRouter.POST("createMcpServer", mcpServerApi.CreateMcpServer)   // 新建mcpServer表
		mcpServerRouter.DELETE("deleteMcpServer", mcpServerApi.DeleteMcpServer) // 删除mcpServer表
		mcpServerRouter.DELETE("deleteMcpServerByIds", mcpServerApi.DeleteMcpServerByIds) // 批量删除mcpServer表
		mcpServerRouter.PUT("updateMcpServer", mcpServerApi.UpdateMcpServer)    // 更新mcpServer表
	}
	{
		mcpServerRouterWithoutRecord.GET("findMcpServer", mcpServerApi.FindMcpServer)        // 根据ID获取mcpServer表
		mcpServerRouterWithoutRecord.GET("getMcpServerList", mcpServerApi.GetMcpServerList)  // 获取mcpServer表列表
	}
	{
	    mcpServerRouterWithoutAuth.GET("getMcpServerPublic", mcpServerApi.GetMcpServerPublic)  // mcpServer表开放接口
	}
}
