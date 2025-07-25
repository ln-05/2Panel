# 设计文档

## 概述

本设计文档概述了在Gin-Vue-Admin系统中实现Docker容器管理API的方案。该功能将提供与Docker容器交互的RESTful端点，包括列出容器、检索容器详细信息和基本容器操作。实现将遵循系统现有的架构模式，并与当前的身份验证和授权机制集成。

## 架构

Docker容器API将按照现有的分层架构实现：

```
┌─────────────────┐
│   前端层        │ ← Vue.js容器管理组件
├─────────────────┤
│   API层         │ ← REST端点 (/api/v1/docker/*)
├─────────────────┤
│  服务层         │ ← 业务逻辑和Docker客户端操作
├─────────────────┤
│  Docker客户端   │ ← Docker Go SDK集成
├─────────────────┤
│  Docker守护进程 │ ← 本地或远程Docker守护进程
└─────────────────┘
```

### 关键设计决策

1. **Docker SDK集成**: 使用官方Docker Go SDK (`github.com/docker/docker/client`) 进行可靠的Docker守护进程通信
2. **基于权限的访问**: 使用新的"docker:manage"权限与现有的Casbin授权系统集成
3. **错误处理**: 为Docker守护进程连接问题实现全面的错误处理
4. **响应格式**: 使用通用响应结构遵循现有的响应模式
5. **中间件集成**: 使用现有的JWT身份验证和操作记录中间件

## 组件和接口

### 1. API层 (`server/api/v1/docker/`)

**文件: `docker_container.go`**
```go
type DockerApi struct{}

// API方法:
// - GetContainerList(c *gin.Context)    // GET /api/v1/docker/containers
// - GetContainerDetail(c *gin.Context)  // GET /api/v1/docker/containers/:id
// - GetContainerLogs(c *gin.Context)    // GET /api/v1/docker/containers/:id/logs
```

### 2. 服务层 (`server/service/docker/`)

**文件: `docker_container.go`**
```go
type DockerService struct {
    client *client.Client
}

// 服务方法:
// - GetContainerList(filter ContainerFilter) ([]ContainerInfo, error)
// - GetContainerDetail(containerID string) (*ContainerDetail, error)
// - GetContainerLogs(containerID string, options LogOptions) (string, error)
```

### 3. 模型层 (`server/model/docker/`)

**请求模型 (`server/model/docker/request/`):**
```go
type ContainerFilter struct {
    Status string `json:"status" form:"status"`
    Name   string `json:"name" form:"name"`
    Page   int    `json:"page" form:"page"`
    PageSize int  `json:"pageSize" form:"pageSize"`
}

type LogOptions struct {
    Tail   string `json:"tail" form:"tail"`
    Since  string `json:"since" form:"since"`
    Follow bool   `json:"follow" form:"follow"`
}
```

**响应模型 (`server/model/docker/response/`):**
```go
type ContainerInfo struct {
    ID      string            `json:"id"`
    Name    string            `json:"name"`
    Image   string            `json:"image"`
    Status  string            `json:"status"`
    State   string            `json:"state"`
    Created int64             `json:"created"`
    Ports   []PortMapping     `json:"ports"`
    Labels  map[string]string `json:"labels"`
}

type ContainerDetail struct {
    ContainerInfo
    Config      ContainerConfig `json:"config"`
    HostConfig  HostConfig      `json:"hostConfig"`
    NetworkSettings NetworkSettings `json:"networkSettings"`
    Mounts      []Mount         `json:"mounts"`
}
```

### 4. 路由层 (`server/router/docker/`)

**文件: `docker_container.go`**
```go
type DockerRouter struct{}

func (d *DockerRouter) InitDockerRouter(Router *gin.RouterGroup) {
    dockerRouter := Router.Group("docker").Use(middleware.OperationRecord())
    dockerRouterWithoutRecord := Router.Group("docker")
    
    // 带操作记录的路由
    dockerRouter.GET("containers", dockerApi.GetContainerList)
    dockerRouter.GET("containers/:id", dockerApi.GetContainerDetail)
    dockerRouter.GET("containers/:id/logs", dockerApi.GetContainerLogs)
}
```

## Data Models

### Container Information Structure

```go
type ContainerInfo struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Image       string            `json:"image"`
    ImageID     string            `json:"imageId"`
    Command     string            `json:"command"`
    Created     int64             `json:"created"`
    Status      string            `json:"status"`
    State       string            `json:"state"`
    Ports       []PortMapping     `json:"ports"`
    Labels      map[string]string `json:"labels"`
    SizeRw      int64             `json:"sizeRw,omitempty"`
    SizeRootFs  int64             `json:"sizeRootFs,omitempty"`
}

type PortMapping struct {
    PrivatePort int    `json:"privatePort"`
    PublicPort  int    `json:"publicPort,omitempty"`
    Type        string `json:"type"`
    IP          string `json:"ip,omitempty"`
}
```

### Docker Client Configuration

```go
type DockerConfig struct {
    Host       string `yaml:"host" json:"host"`
    Version    string `yaml:"version" json:"version"`
    TLSVerify  bool   `yaml:"tls-verify" json:"tlsVerify"`
    CertPath   string `yaml:"cert-path" json:"certPath"`
    Timeout    int    `yaml:"timeout" json:"timeout"`
}
```

## Error Handling

### Error Types and Responses

1. **Docker Daemon Connection Errors**
   - HTTP Status: 503 Service Unavailable
   - Message: "Docker daemon is not accessible"

2. **Container Not Found Errors**
   - HTTP Status: 404 Not Found
   - Message: "Container not found"

3. **Permission Errors**
   - HTTP Status: 403 Forbidden
   - Message: "Insufficient permissions to access Docker resources"

4. **Invalid Request Errors**
   - HTTP Status: 400 Bad Request
   - Message: Specific validation error message

### Error Handling Strategy

```go
func handleDockerError(err error) (int, string) {
    if client.IsErrNotFound(err) {
        return 404, "Container not found"
    }
    if client.IsErrConnectionFailed(err) {
        return 503, "Docker daemon is not accessible"
    }
    return 500, "Internal server error"
}
```

## Testing Strategy

### Unit Tests

1. **Service Layer Tests** (`server/service/docker/docker_container_test.go`)
   - Mock Docker client for isolated testing
   - Test container listing with various filters
   - Test error handling scenarios
   - Test data transformation logic

2. **API Layer Tests** (`server/api/v1/docker/docker_container_test.go`)
   - Test HTTP request/response handling
   - Test authentication and authorization
   - Test input validation
   - Test error response formatting

### Integration Tests

1. **Docker Integration Tests**
   - Test with real Docker daemon (in CI environment)
   - Test container lifecycle operations
   - Test network connectivity scenarios

### Test Data Setup

```go
// Mock container data for testing
var mockContainers = []types.Container{
    {
        ID:      "abc123",
        Names:   []string{"/test-container"},
        Image:   "nginx:latest",
        Status:  "Up 2 hours",
        State:   "running",
        Created: time.Now().Unix(),
    },
}
```

## Security Considerations

### Authentication and Authorization

1. **JWT Authentication**: All Docker API endpoints require valid JWT tokens
2. **Permission-based Access**: New permission `docker:manage` will be created
3. **Role-based Access Control**: Only users with appropriate roles can access Docker APIs

### Docker Security

1. **Docker Socket Access**: Secure access to Docker socket with appropriate permissions
2. **Container Isolation**: Ensure API cannot access sensitive container information
3. **Resource Limits**: Implement rate limiting for Docker API calls

### Configuration Security

```yaml
# config.yaml addition
docker:
  host: "unix:///var/run/docker.sock"  # Docker daemon socket
  version: "1.41"                       # Docker API version
  timeout: 30                          # Connection timeout in seconds
  tls-verify: false                    # TLS verification for remote Docker
  cert-path: ""                        # Certificate path for TLS
```

## Performance Considerations

### Caching Strategy

1. **Container List Caching**: Cache container list for 30 seconds to reduce Docker API calls
2. **Container Detail Caching**: Cache individual container details for 60 seconds

### Pagination

1. **Client-side Pagination**: Implement pagination in the service layer to handle large container lists
2. **Filtering**: Support server-side filtering to reduce data transfer

### Connection Pooling

1. **Docker Client Reuse**: Maintain persistent Docker client connection
2. **Connection Health Checks**: Implement health checks for Docker daemon connectivity

## API Endpoints Specification

### 1. List Containers
- **Endpoint**: `GET /api/v1/docker/containers`
- **Query Parameters**: 
  - `status` (optional): Filter by container status
  - `name` (optional): Filter by container name
  - `page` (optional): Page number for pagination
  - `pageSize` (optional): Number of items per page
- **Response**: Paginated list of containers

### 2. Get Container Detail
- **Endpoint**: `GET /api/v1/docker/containers/:id`
- **Path Parameters**: 
  - `id`: Container ID or name
- **Response**: Detailed container information

### 3. Get Container Logs
- **Endpoint**: `GET /api/v1/docker/containers/:id/logs`
- **Path Parameters**: 
  - `id`: Container ID or name
- **Query Parameters**:
  - `tail` (optional): Number of lines to tail
  - `since` (optional): Show logs since timestamp
- **Response**: Container logs as text

## Dependencies

### New Go Dependencies

```go
// go.mod additions
require (
    github.com/docker/docker v24.0.0+incompatible
    github.com/docker/go-connections v0.4.0
    github.com/moby/term v0.5.0
)
```

### Configuration Updates

The system configuration will be extended to include Docker-specific settings in the main config structure.