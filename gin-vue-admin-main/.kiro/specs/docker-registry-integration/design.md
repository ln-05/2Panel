# Docker仓库管理前后端对接设计文档

## 概述

本设计文档描述了Docker仓库管理功能的前后端对接架构。系统已经具备完整的后端API和前端页面，需要确保两者能够正确对接并提供完整的用户体验。

## 架构

### 系统架构图

```
前端 (Vue.js)
    ↓ HTTP请求
API网关 (Gin)
    ↓ 服务调用
业务服务层 (DockerRegistryService)
    ↓ 数据访问
数据库 (MySQL/PostgreSQL)
```

### 组件交互流程

1. 用户在前端页面进行操作
2. 前端通过axios发送HTTP请求到后端API
3. Gin路由器将请求分发到对应的API处理器
4. API处理器调用业务服务层处理业务逻辑
5. 服务层通过GORM操作数据库
6. 结果通过相同路径返回给前端

## 组件和接口

### 前端组件

#### 1. 仓库管理页面 (`registry.vue`)
- **位置**: `web/src/view/my/registry/registry.vue`
- **功能**: 提供仓库管理的用户界面
- **依赖**: `dockerRegistry.js` API模块

#### 2. API接口模块 (`dockerRegistry.js`)
- **位置**: `web/src/api/dockerRegistry.js`
- **功能**: 封装所有仓库相关的API调用
- **方法**:
  - `getRegistryList(params)` - 获取仓库列表
  - `getRegistryDetail(id)` - 获取仓库详情
  - `createRegistry(data)` - 创建仓库
  - `updateRegistry(data)` - 更新仓库
  - `deleteRegistry(id)` - 删除仓库
  - `testRegistry(id)` - 测试仓库连接
  - `setDefaultRegistry(id)` - 设置默认仓库

### 后端组件

#### 1. API控制器 (`DockerRegistryApi`)
- **位置**: `server/api/v1/docker/docker_registry.go`
- **功能**: 处理HTTP请求和响应
- **方法**: 对应前端API的所有方法

#### 2. 业务服务 (`DockerRegistryService`)
- **位置**: `server/service/docker/docker_registry.go`
- **功能**: 实现业务逻辑和数据处理

#### 3. 数据模型
- **位置**: `server/model/docker/`
- **文件**:
  - `docker_registry.go` - 数据库模型
  - `request/registry_request.go` - 请求模型
  - `response/registry_response.go` - 响应模型

#### 4. 路由配置
- **位置**: `server/router/docker/docker_registry.go`
- **功能**: 配置API路由映射

## 数据模型

### 数据库表结构 (`docker_registries`)

```sql
CREATE TABLE docker_registries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name VARCHAR(100) NOT NULL UNIQUE,
    download_url VARCHAR(500) NOT NULL,
    protocol VARCHAR(10) NOT NULL DEFAULT 'https',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    username VARCHAR(100),
    password VARCHAR(200),
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    last_test_time DATETIME,
    test_result VARCHAR(500),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE INDEX idx_name (name)
);
```

### API数据格式

#### 请求格式示例

```json
// 创建仓库请求
{
    "name": "Docker Hub",
    "downloadUrl": "https://registry-1.docker.io",
    "protocol": "https",
    "username": "myuser",
    "password": "mypass",
    "description": "官方Docker Hub仓库"
}
```

#### 响应格式示例

```json
// 仓库列表响应
{
    "code": 0,
    "data": {
        "list": [
            {
                "id": 1,
                "name": "Docker Hub",
                "downloadUrl": "https://registry-1.docker.io",
                "protocol": "https",
                "status": "active",
                "username": "myuser",
                "description": "官方Docker Hub仓库",
                "createdAt": "2025-01-20T10:00:00Z",
                "updatedAt": "2025-01-20T10:00:00Z"
            }
        ],
        "total": 1
    },
    "msg": "获取成功"
}
```

## 错误处理

### 前端错误处理策略

1. **网络错误**: 显示"网络连接失败"提示
2. **认证错误**: 跳转到登录页面
3. **权限错误**: 显示"权限不足"提示
4. **业务错误**: 显示具体的错误信息
5. **表单验证错误**: 在对应字段显示验证提示

### 后端错误响应格式

```json
{
    "code": 7,
    "data": null,
    "msg": "仓库名称已存在"
}
```

### 错误码定义

- `0`: 成功
- `7`: 业务逻辑错误
- `401`: 认证失败
- `403`: 权限不足
- `500`: 服务器内部错误

## 测试策略

### 单元测试

1. **前端组件测试**: 使用Vue Test Utils测试组件功能
2. **API接口测试**: 使用Jest测试API调用逻辑
3. **后端服务测试**: 使用Go testing包测试业务逻辑
4. **数据库操作测试**: 测试CRUD操作的正确性

### 集成测试

1. **API集成测试**: 测试前后端API调用的完整流程
2. **数据库集成测试**: 测试数据持久化的正确性
3. **用户界面测试**: 使用Cypress进行端到端测试

### 测试用例

1. **正常流程测试**: 测试所有功能的正常使用场景
2. **异常处理测试**: 测试各种错误情况的处理
3. **边界条件测试**: 测试输入验证和边界值处理
4. **性能测试**: 测试大量数据情况下的系统性能

## 部署考虑

### 环境配置

1. **开发环境**: 本地开发和测试
2. **测试环境**: 集成测试和用户验收测试
3. **生产环境**: 正式部署环境

### 配置管理

1. **数据库连接**: 通过环境变量配置
2. **API基础URL**: 前端通过环境变量配置
3. **认证配置**: JWT密钥和过期时间配置
4. **日志配置**: 日志级别和输出配置

### 监控和日志

1. **API调用监控**: 记录API调用次数和响应时间
2. **错误日志**: 记录所有错误信息用于问题排查
3. **业务日志**: 记录重要的业务操作
4. **性能监控**: 监控系统资源使用情况