# 设计文档

## 概述

Docker编排管理界面是一个基于Vue 3和Element Plus的前端界面，配合Go Gin后端API，提供类似1Panel的容器编排管理功能。该系统将允许用户通过直观的Web界面管理Docker容器的完整生命周期，包括创建、启动、停止、重启、删除等操作。

## 架构

### 前端架构
- **框架**: Vue 3 + Composition API
- **UI组件库**: Element Plus
- **状态管理**: 基于Vue 3的响应式系统
- **路由**: Vue Router
- **HTTP客户端**: Axios
- **构建工具**: Vite

### 后端架构
- **框架**: Gin (Go)
- **Docker集成**: Docker SDK for Go
- **数据库**: MySQL/SQLite (通过GORM)
- **认证**: JWT
- **权限控制**: Casbin

### 系统集成
- 复用现有的Docker容器API服务
- 集成现有的认证和权限系统
- 扩展现有的Docker服务层功能

## 组件和接口

### 前端组件结构

#### 1. 主容器编排页面 (OrchestrationView.vue)
```
OrchestrationView/
├── components/
│   ├── ContainerTable.vue          # 容器列表表格
│   ├── ContainerOperations.vue     # 容器操作按钮组
│   ├── ContainerCreateDialog.vue   # 创建容器对话框
│   ├── ContainerDetailDialog.vue   # 容器详情对话框
│   ├── ContainerLogDialog.vue      # 容器日志对话框
│   └── SearchAndFilter.vue         # 搜索和过滤组件
└── OrchestrationView.vue           # 主页面组件
```

#### 2. 核心组件设计

**ContainerTable.vue**
- 显示容器列表（名称、来源、编辑链接、容器数量、应用数量、创建时间）
- 支持多选操作
- 实时状态更新
- 分页功能

**ContainerOperations.vue**
- 创建编排按钮
- 批量操作按钮（删除、启动、停止）
- 刷新按钮

**ContainerCreateDialog.vue**
- 容器配置表单
- 镜像选择
- 端口映射配置
- 卷挂载配置
- 环境变量配置

### 后端API接口

#### 1. 容器编排API (docker_orchestration.go)

```go
// 获取容器编排列表
GET /api/v1/docker/orchestrations
// 参数: page, pageSize, name, status

// 获取容器编排详情
GET /api/v1/docker/orchestrations/{id}

// 创建容器编排
POST /api/v1/docker/orchestrations
// Body: OrchestrationCreateRequest

// 更新容器编排
PUT /api/v1/docker/orchestrations/{id}
// Body: OrchestrationUpdateRequest

// 删除容器编排
DELETE /api/v1/docker/orchestrations/{id}

// 启动容器编排
POST /api/v1/docker/orchestrations/{id}/start

// 停止容器编排
POST /api/v1/docker/orchestrations/{id}/stop

// 重启容器编排
POST /api/v1/docker/orchestrations/{id}/restart

// 获取容器日志
GET /api/v1/docker/orchestrations/{id}/logs
// 参数: lines, follow, timestamps
```

#### 2. 扩展现有容器API

```go
// 批量操作容器
POST /api/v1/docker/containers/batch
// Body: BatchOperationRequest

// 获取容器实时状态
GET /api/v1/docker/containers/status
// WebSocket连接用于实时更新
```

### 数据模型

#### 1. 前端数据模型

```typescript
// 容器编排信息
interface OrchestrationInfo {
  id: number
  name: string
  source: string
  editLink: string
  containerCount: number
  applicationCount: number
  status: 'running' | 'stopped' | 'error'
  createdAt: string
  updatedAt: string
}

// 容器创建请求
interface ContainerCreateRequest {
  name: string
  image: string
  ports: PortMapping[]
  volumes: VolumeMount[]
  environment: Record<string, string>
  restartPolicy: string
  networkMode: string
}

// 端口映射
interface PortMapping {
  hostPort: number
  containerPort: number
  protocol: 'tcp' | 'udp'
}

// 卷挂载
interface VolumeMount {
  hostPath: string
  containerPath: string
  readOnly: boolean
}
```

#### 2. 后端数据模型

```go
// 容器编排模型
type DockerOrchestration struct {
    gorm.Model
    Name            string `json:"name" gorm:"not null;uniqueIndex"`
    Description     string `json:"description"`
    ComposeContent  string `json:"composeContent" gorm:"type:text"`
    Status          string `json:"status" gorm:"default:'stopped'"`
    ContainerCount  int    `json:"containerCount"`
    ApplicationCount int   `json:"applicationCount"`
    Source          string `json:"source" gorm:"default:'manual'"`
    EditLink        string `json:"editLink"`
}

// 容器编排请求模型
type OrchestrationCreateRequest struct {
    Name           string                 `json:"name" binding:"required"`
    Description    string                 `json:"description"`
    ComposeContent string                 `json:"composeContent"`
    Services       []ServiceConfig        `json:"services"`
}

type ServiceConfig struct {
    Name         string            `json:"name"`
    Image        string            `json:"image"`
    Ports        []PortMapping     `json:"ports"`
    Volumes      []VolumeMount     `json:"volumes"`
    Environment  map[string]string `json:"environment"`
    RestartPolicy string           `json:"restartPolicy"`
}
```

## 错误处理

### 前端错误处理
1. **网络错误**: 显示友好的错误消息，提供重试选项
2. **权限错误**: 重定向到登录页面或显示权限不足提示
3. **Docker服务错误**: 显示Docker服务状态和连接指导
4. **表单验证错误**: 实时验证和错误提示

### 后端错误处理
1. **Docker连接错误**: 统一的Docker客户端错误处理
2. **容器操作错误**: 详细的操作失败原因和建议
3. **权限验证错误**: 标准的HTTP状态码和错误消息
4. **数据验证错误**: 结构化的验证错误响应

## 测试策略

### 前端测试
1. **单元测试**: 使用Vitest测试组件逻辑
2. **组件测试**: 使用Vue Test Utils测试组件交互
3. **集成测试**: 测试API调用和数据流
4. **E2E测试**: 使用Cypress测试完整用户流程

### 后端测试
1. **单元测试**: 测试服务层逻辑
2. **API测试**: 测试HTTP端点
3. **集成测试**: 测试Docker操作
4. **Mock测试**: 模拟Docker客户端进行测试

### 测试场景
1. **容器生命周期管理**: 创建、启动、停止、删除
2. **批量操作**: 多容器同时操作
3. **实时状态更新**: WebSocket连接和状态同步
4. **错误处理**: 各种异常情况的处理
5. **权限控制**: 不同用户角色的访问控制

## 性能优化

### 前端优化
1. **虚拟滚动**: 大量容器列表的性能优化
2. **懒加载**: 按需加载容器详情和日志
3. **缓存策略**: 合理缓存容器状态和配置
4. **防抖节流**: 搜索和过滤操作的性能优化

### 后端优化
1. **连接池**: Docker客户端连接复用
2. **并发控制**: 批量操作的并发限制
3. **缓存机制**: 容器状态和镜像信息缓存
4. **分页优化**: 大量数据的分页查询优化

## 安全考虑

### 认证和授权
1. **JWT认证**: 复用现有的JWT认证机制
2. **权限控制**: 基于Casbin的细粒度权限控制
3. **操作审计**: 记录所有容器操作日志

### Docker安全
1. **权限限制**: 限制Docker操作权限
2. **资源限制**: 容器资源使用限制
3. **网络隔离**: 容器网络安全配置
4. **镜像安全**: 镜像来源验证和安全扫描

## 部署和监控

### 部署策略
1. **容器化部署**: 应用本身的容器化
2. **配置管理**: 环境配置的外部化
3. **健康检查**: 应用和Docker服务的健康监控

### 监控指标
1. **容器状态**: 实时监控容器运行状态
2. **资源使用**: CPU、内存、磁盘使用情况
3. **操作日志**: 用户操作和系统事件日志
4. **性能指标**: API响应时间和错误率