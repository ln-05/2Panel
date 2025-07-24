package AI

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	OllamaModelApi
	McpServerApi
}

var (
	ollamaModelService = service.ServiceGroupApp.AIServiceGroup.OllamaModelService
	mcpServerService   = service.ServiceGroupApp.AIServiceGroup.McpServerService
)
