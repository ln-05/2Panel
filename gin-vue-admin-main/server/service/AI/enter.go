package AI

type ServiceGroup struct {
	OllamaModelService *OllamaModelService
	McpServerService
}

// 初始化服务组
func init() {
	ServiceGroupApp.OllamaModelService = NewOllamaModelService()
}

var ServiceGroupApp = new(ServiceGroup)
