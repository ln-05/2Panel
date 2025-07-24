
// 自动生成模板McpServer
package AI
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// mcpServer表 结构体  McpServer
type McpServer struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:服务器名称;column:name;size:255;"`  //服务器名称
  DockerCompose  *string `json:"dockerCompose" form:"dockerCompose" gorm:"comment:Docker Compose配置;column:docker_compose;"`  //配置
  Command  *string `json:"command" form:"command" gorm:"comment:启动命令;column:command;"`  //启动命令
  ContainerName  *string `json:"containerName" form:"containerName" gorm:"comment:容器名称;column:container_name;size:255;"`  //容器名称
  Message  *string `json:"message" form:"message" gorm:"comment:消息;column:message;"`  //消息
  Port  *int `json:"port" form:"port" gorm:"comment:端口号;column:port;size:10;"`  //端口号
  Status  *string `json:"status" form:"status" gorm:"comment:状态;column:status;size:255;"`  //状态
  Env  *string `json:"env" form:"env" gorm:"comment:环境变量;column:env;"`  //环境变量
  BaseUrl  *string `json:"baseUrl" form:"baseUrl" gorm:"comment:基础URL;column:base_url;size:255;"`  //基础URL
  SsePath  *string `json:"ssePath" form:"ssePath" gorm:"comment:SSE路径;column:sse_path;size:255;"`  //SSE路径
  WebsiteId  *int `json:"websiteId" form:"websiteId" gorm:"comment:网站ID;column:website_id;size:10;"`  //网站ID
  Dir  *string `json:"dir" form:"dir" gorm:"comment:目录;column:dir;size:255;"`  //目录
  HostIp  *string `json:"hostIp" form:"hostIp" gorm:"comment:主机IP;column:host_ip;size:255;"`  //主机IP
    CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
    UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
    DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName mcpServer表 McpServer自定义表名 mcp_server
func (McpServer) TableName() string {
    return "mcp_server"
}





