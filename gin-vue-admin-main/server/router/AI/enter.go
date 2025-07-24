package AI

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	OllamaModelRouter
	McpServerRouter
}

var (
	ollamaModelApi = api.ApiGroupApp.AIApiGroup.OllamaModelApi
	mcpServerApi   = api.ApiGroupApp.AIApiGroup.McpServerApi
)
