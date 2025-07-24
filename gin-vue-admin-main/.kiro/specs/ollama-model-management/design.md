# Ollama模型管理系统设计文档

## 概述

本文档描述了Ollama模型管理系统的技术设计，该系统基于Gin-Vue-Admin框架，提供完整的本地AI模型管理和部署能力。

## 紧急修复设计方案

### 问题诊断和修复策略

基于当前系统状态分析，主要问题集中在以下几个方面：

1. **WebSocket连接失效**：需要重新设计WebSocket连接管理机制
2. **前端组件交互问题**：需要修复组件间的数据传递和事件处理
3. **后端API响应问题**：需要确保API接口正常工作
4. **状态管理混乱**：需要统一状态管理和数据流

### 修复优先级

**P0 (立即修复)**：
- WebSocket连接和实时状态更新
- 模型操作按钮功能恢复
- 服务状态检测修复
- 基础错误处理

**P1 (短期修复)**：
- 批量操作功能
- 搜索和筛选功能
- 进度显示和用户反馈

**P2 (中期完善)**：
- 高级功能（域名绑定、日志查看等）
- 性能优化
- 用户体验改进

## 紧急修复技术方案

### 1. WebSocket连接修复设计

#### 问题分析
- 当前WebSocket连接不稳定，导致实时状态更新失效
- 连接断开后无法自动重连
- 消息处理机制存在问题

#### 修复方案
```typescript
// 简化的WebSocket管理器
class SimpleWebSocketManager {
  private socket: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000
  
  connect() {
    try {
      this.socket = new WebSocket(this.getWebSocketUrl())
      this.setupEventHandlers()
    } catch (error) {
      this.handleConnectionError(error)
    }
  }
  
  private setupEventHandlers() {
    if (!this.socket) return
    
    this.socket.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
    }
    
    this.socket.onmessage = (event) => {
      this.handleMessage(JSON.parse(event.data))
    }
    
    this.socket.onclose = () => {
      this.handleReconnect()
    }
    
    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }
  
  private handleReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      setTimeout(() => this.connect(), this.reconnectDelay)
    }
  }
}
```

### 2. 前端组件修复设计

#### 问题分析
- 组件间props传递不正确
- 事件处理函数未正确绑定
- 状态更新不及时

#### 修复方案
```vue
<!-- 主组件修复 -->
<template>
  <div class="ollama-model-management">
    <!-- 服务状态显示 -->
    <AppStatus 
      :status="serviceStatus" 
      :is-connected="isWebSocketConnected"
      @refresh="handleRefreshStatus"
    />
    
    <!-- 模型表格 -->
    <ModelTable 
      :models="modelList"
      :loading="tableLoading"
      @start-model="handleStartModel"
      @stop-model="handleStopModel"
      @delete-model="handleDeleteModel"
      @batch-operation="handleBatchOperation"
    />
  </div>
</template>

<script setup lang="ts">
// 使用简化的状态管理
const serviceStatus = ref('unknown')
const isWebSocketConnected = ref(false)
const modelList = ref([])
const tableLoading = ref(false)

// 修复后的事件处理函数
const handleStartModel = async (model: OllamaModelInfo) => {
  try {
    tableLoading.value = true
    await startOllamaModel(model.id)
    ElMessage.success(`模型 ${model.name} 启动成功`)
    await refreshModelList()
  } catch (error) {
    ElMessage.error(`启动模型失败: ${error.message}`)
  } finally {
    tableLoading.value = false
  }
}
</script>
```

### 3. 后端API修复设计

#### 问题分析
- API响应格式不统一
- 错误处理不完善
- 状态更新不及时

#### 修复方案
```go
// 统一的API响应格式
type APIResponse struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
}

// 修复后的模型启动接口
func (o *OllamaModelApi) StartOllamaModel(c *gin.Context) {
    var req request.OllamaModelStart
    if err := c.ShouldBindJSON(&req); err != nil {
        response.FailWithMessage(err.Error(), c)
        return
    }
    
    // 验证模型是否存在
    model, err := ollamaModelService.GetOllamaModel(req.ID)
    if err != nil {
        response.FailWithMessage("模型不存在", c)
        return
    }
    
    // 启动模型
    if err := ollamaModelService.StartModel(model); err != nil {
        response.FailWithMessage("启动模型失败: "+err.Error(), c)
        return
    }
    
    // 返回成功响应
    response.OkWithData(gin.H{
        "id": model.ID,
        "name": model.Name,
        "status": "starting",
    }, c)
}
```

### 4. 状态管理修复设计

#### 问题分析
- 状态管理分散，难以维护
- 数据流不清晰
- 状态同步问题

#### 修复方案
```typescript
// 统一的状态管理
export const useOllamaStore = defineStore('ollama', {
  state: () => ({
    // 服务状态
    serviceStatus: 'unknown' as ServiceStatus,
    isConnected: false,
    
    // 模型数据
    models: [] as OllamaModelInfo[],
    selectedModels: [] as OllamaModelInfo[],
    
    // UI状态
    loading: false,
    tableLoading: false,
    
    // 错误状态
    lastError: null as string | null,
  }),
  
  actions: {
    // 统一的错误处理
    handleError(error: any, context: string) {
      console.error(`[${context}]`, error)
      this.lastError = error.message || '未知错误'
      ElMessage.error(this.lastError)
    },
    
    // 统一的模型操作
    async startModel(modelId: number) {
      try {
        this.tableLoading = true
        const result = await startOllamaModel(modelId)
        
        // 更新本地状态
        const model = this.models.find(m => m.id === modelId)
        if (model) {
          model.status = 'starting'
        }
        
        return result
      } catch (error) {
        this.handleError(error, 'startModel')
        throw error
      } finally {
        this.tableLoading = false
      }
    }
  }
})
```

## 架构设计

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                    前端层 (Vue 3 + TypeScript)              │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   模型管理   │  │   GPU监控    │  │  MCP服务器   │         │
│  │   /model    │  │    /gpu     │  │    /mcp     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   添加模型   │  │   域名绑定   │  │   终端操作   │         │
│  │  /model/add │  │/model/domain│  │/model/terminal│       │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│                    API层 (RESTful)                         │
├─────────────────────────────────────────────────────────────┤
│                   后端层 (Gin + GORM)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   API层     │  │  Service层   │  │   Model层    │         │
│  │ Controller  │  │  Business    │  │   Database   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
├─────────────────────────────────────────────────────────────┤
│                   基础设施层                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Ollama    │  │   Docker     │  │   MySQL     │         │
│  │   Service   │  │  Container   │  │  Database   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

## 组件设计

### 前端组件架构

```
frontend/src/views/ai/
├── model/           # Ollama模型管理
│   ├── index.vue    # 主页面
│   ├── add/         # 添加模型
│   │   ├── index.vue
│   │   └── components/
│   ├── domain/      # 域名绑定
│   │   ├── index.vue
│   │   └── components/
│   ├── terminal/    # 终端操作
│   │   ├── index.vue
│   │   └── components/
│   ├── del/         # 删除模型
│   │   └── index.vue
│   └── conn/        # 连接信息
│       └── index.vue
├── gpu/             # GPU监控
│   └── index.vue
└── mcp/             # MCP服务器管理
    └── index.vue
```

### 数据模型设计

#### 后端数据模型

```go
// OllamaModel 模型实体
type OllamaModel struct {
    global.GVA_MODEL
    Name         string `json:"name"`         // 模型名称
    Size         string `json:"size"`         // 模型大小
    From         string `json:"from"`         // 来源
    Status       string `json:"status"`       // 状态
    Message      string `json:"message"`      // 消息
    LogFileExist bool   `json:"logFileExist"` // 日志文件存在
}

// OllamaBindDomain 域名绑定实体
type OllamaBindDomain struct {
    global.GVA_MODEL
    Domain       string `json:"domain"`       // 域名
    AppInstallID uint   `json:"appInstallID"` // 应用安装ID
    SSLID        uint   `json:"sslID"`        // SSL证书ID
    WebsiteID    uint   `json:"websiteID"`    // 网站ID
    IPList       string `json:"ipList"`       // IP白名单
}
```

#### 前端类型定义

```typescript
// 模型信息接口
interface OllamaModelInfo {
  id: number
  name: string
  size: string
  from: string
  logFileExist: boolean
  status: ModelStatus
  message: string
  createdAt: string
}

// 模型状态枚举
enum ModelStatus {
  STOPPED = 'stopped',
  RUNNING = 'running',
  DOWNLOADING = 'downloading',
  ERROR = 'error',
  UNAVAILABLE = 'unavailable'
}

// 域名绑定接口
interface OllamaBindDomain {
  id: number
  domain: string
  appInstallID: number
  sslID?: number
  websiteID?: number
  ipList?: string
}
```

## 接口设计

### RESTful API 设计

#### 模型管理接口

```
GET    /api/v1/ollamaModel/search          # 搜索模型
POST   /api/v1/ollamaModel/create          # 创建模型
POST   /api/v1/ollamaModel/close           # 停止模型
POST   /api/v1/ollamaModel/recreate        # 重新创建模型
DELETE /api/v1/ollamaModel/deleteOllamaModel # 删除模型
POST   /api/v1/ollamaModel/sync            # 同步模型
GET    /api/v1/ollamaModel/detail          # 获取详情
```

#### 域名绑定接口

```
POST   /api/v1/ollamaModel/bindDomain      # 绑定域名
GET    /api/v1/ollamaModel/getBindDomain   # 获取绑定信息
PUT    /api/v1/ollamaModel/updateBindDomain # 更新绑定
```

### WebSocket 实时通信

```
ws://localhost:8888/ws/ollama/status       # 模型状态推送
ws://localhost:8888/ws/ollama/progress     # 下载进度推送
ws://localhost:8888/ws/ollama/logs         # 日志实时推送
```

## 用户界面设计

### 主页面布局

```vue
<template>
  <div v-loading="loading">
    <!-- 面包屑导航 -->
    <RouterButton />
    
    <!-- 主要内容区域 -->
    <LayoutContent title="Ollama模型管理">
      <!-- 应用状态组件 -->
      <template #app>
        <AppStatus :status="ollamaStatus" />
      </template>
      
      <!-- 提示信息 -->
      <template #prompt>
        <el-alert 
          type="info" 
          :title="$t('ollama.tips.title')"
          :description="$t('ollama.tips.description')"
          show-icon
        />
      </template>
      
      <!-- 左侧工具栏 -->
      <template #leftToolBar>
        <el-button-group>
          <el-button type="primary" @click="showAddDialog">
            <el-icon><Plus /></el-icon>
            {{ $t('ollama.actions.add') }}
          </el-button>
          <el-button @click="syncModels">
            <el-icon><Refresh /></el-icon>
            {{ $t('ollama.actions.sync') }}
          </el-button>
          <el-button @click="showTerminal">
            <el-icon><Monitor /></el-icon>
            {{ $t('ollama.actions.terminal') }}
          </el-button>
        </el-button-group>
      </template>
      
      <!-- 右侧工具栏 -->
      <template #rightToolBar>
        <el-input
          v-model="searchQuery"
          :placeholder="$t('ollama.search.placeholder')"
          style="width: 200px"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button @click="refreshData">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </template>
      
      <!-- 主要内容 -->
      <template #main>
        <ModelTable 
          :data="modelList"
          :loading="tableLoading"
          @edit="handleEdit"
          @delete="handleDelete"
          @start="handleStart"
          @stop="handleStop"
          @detail="handleDetail"
        />
      </template>
    </LayoutContent>
    
    <!-- 对话框组件 -->
    <AddModelDialog v-model="addDialogVisible" @success="refreshData" />
    <DomainBindDialog v-model="domainDialogVisible" :model="selectedModel" />
    <TerminalDialog v-model="terminalVisible" />
  </div>
</template>
```

### 组件功能设计

#### 1. ModelTable 组件

```vue
<template>
  <el-table :data="data" :loading="loading" row-key="id">
    <el-table-column type="selection" width="55" />
    <el-table-column prop="name" :label="$t('ollama.table.name')" />
    <el-table-column prop="size" :label="$t('ollama.table.size')" />
    <el-table-column prop="status" :label="$t('ollama.table.status')">
      <template #default="{ row }">
        <StatusBadge :status="row.status" />
      </template>
    </el-table-column>
    <el-table-column prop="createdAt" :label="$t('ollama.table.createdAt')" />
    <el-table-column :label="$t('common.actions')" width="300">
      <template #default="{ row }">
        <ActionButtons 
          :model="row"
          @start="$emit('start', row)"
          @stop="$emit('stop', row)"
          @edit="$emit('edit', row)"
          @delete="$emit('delete', row)"
          @detail="$emit('detail', row)"
        />
      </template>
    </el-table-column>
  </el-table>
</template>
```

#### 2. AddModelDialog 组件

```vue
<template>
  <el-dialog v-model="visible" :title="$t('ollama.dialog.add.title')">
    <el-form :model="form" :rules="rules" ref="formRef">
      <el-form-item :label="$t('ollama.form.name')" prop="name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item :label="$t('ollama.form.from')" prop="from">
        <el-select v-model="form.from">
          <el-option label="Ollama Hub" value="ollama" />
          <el-option label="Hugging Face" value="huggingface" />
        </el-select>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <el-button @click="visible = false">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleSubmit">
        {{ $t('common.confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>
```

#### 3. TerminalDialog 组件

```vue
<template>
  <el-dialog 
    v-model="visible" 
    :title="$t('ollama.terminal.title')"
    width="80%"
    top="5vh"
  >
    <div class="terminal-container">
      <div ref="terminalRef" class="terminal"></div>
    </div>
    
    <template #footer>
      <el-button @click="clearTerminal">{{ $t('ollama.terminal.clear') }}</el-button>
      <el-button @click="visible = false">{{ $t('common.close') }}</el-button>
    </template>
  </el-dialog>
</template>
```

## 状态管理设计

### Pinia Store 设计

```typescript
// stores/ollama.ts
export const useOllamaStore = defineStore('ollama', {
  state: () => ({
    models: [] as OllamaModelInfo[],
    loading: false,
    ollamaStatus: 'unknown' as 'running' | 'stopped' | 'unknown',
    downloadProgress: {} as Record<string, number>,
  }),
  
  getters: {
    runningModels: (state) => 
      state.models.filter(m => m.status === ModelStatus.RUNNING),
    downloadingModels: (state) => 
      state.models.filter(m => m.status === ModelStatus.DOWNLOADING),
  },
  
  actions: {
    async fetchModels() {
      this.loading = true
      try {
        const response = await searchOllamaModel({})
        this.models = response.data.list
      } finally {
        this.loading = false
      }
    },
    
    async createModel(data: CreateModelRequest) {
      await createOllamaModelAdvanced(data)
      await this.fetchModels()
    },
    
    async syncModels() {
      const result = await syncOllamaModel({ force: false })
      await this.fetchModels()
      return result
    },
    
    updateModelStatus(modelId: number, status: ModelStatus, message?: string) {
      const model = this.models.find(m => m.id === modelId)
      if (model) {
        model.status = status
        if (message) model.message = message
      }
    },
    
    updateDownloadProgress(modelName: string, progress: number) {
      this.downloadProgress[modelName] = progress
    }
  }
})
```

## 错误处理设计

### 错误处理策略

```typescript
// utils/errorHandler.ts
export class OllamaErrorHandler {
  static handle(error: any, context: string) {
    console.error(`[${context}] Error:`, error)
    
    if (error.response?.status === 404) {
      ElMessage.error(i18n.global.t('errors.modelNotFound'))
    } else if (error.response?.status === 409) {
      ElMessage.error(i18n.global.t('errors.modelExists'))
    } else if (error.code === 'NETWORK_ERROR') {
      ElMessage.error(i18n.global.t('errors.networkError'))
    } else {
      ElMessage.error(error.message || i18n.global.t('errors.unknown'))
    }
  }
}
```

## 性能优化设计

### 1. 组件懒加载

```typescript
// router/ai.ts
const routes = [
  {
    path: '/ai/model',
    component: () => import('@/views/ai/model/index.vue'),
    children: [
      {
        path: 'add',
        component: () => import('@/views/ai/model/add/index.vue')
      },
      {
        path: 'domain',
        component: () => import('@/views/ai/model/domain/index.vue')
      }
    ]
  }
]
```

### 2. 数据缓存策略

```typescript
// composables/useOllamaCache.ts
export function useOllamaCache() {
  const cache = new Map<string, { data: any; timestamp: number }>()
  const CACHE_DURATION = 5 * 60 * 1000 // 5分钟
  
  const get = (key: string) => {
    const item = cache.get(key)
    if (item && Date.now() - item.timestamp < CACHE_DURATION) {
      return item.data
    }
    return null
  }
  
  const set = (key: string, data: any) => {
    cache.set(key, { data, timestamp: Date.now() })
  }
  
  return { get, set }
}
```

## 国际化设计

### 多语言支持

```typescript
// locales/zh-CN.ts
export default {
  ollama: {
    title: 'Ollama模型管理',
    actions: {
      add: '添加模型',
      sync: '同步模型',
      terminal: '终端',
      start: '启动',
      stop: '停止',
      delete: '删除',
      detail: '详情'
    },
    table: {
      name: '模型名称',
      size: '大小',
      status: '状态',
      createdAt: '创建时间'
    },
    status: {
      running: '运行中',
      stopped: '已停止',
      downloading: '下载中',
      error: '错误',
      unavailable: '不可用'
    }
  }
}
```

## 测试设计

### 单元测试

```typescript
// tests/components/ModelTable.spec.ts
describe('ModelTable', () => {
  it('should render model list correctly', () => {
    const wrapper = mount(ModelTable, {
      props: {
        data: mockModelList,
        loading: false
      }
    })
    
    expect(wrapper.find('.el-table').exists()).toBe(true)
    expect(wrapper.findAll('.el-table__row')).toHaveLength(mockModelList.length)
  })
  
  it('should emit events when actions are clicked', async () => {
    const wrapper = mount(ModelTable, {
      props: { data: mockModelList, loading: false }
    })
    
    await wrapper.find('.start-btn').trigger('click')
    expect(wrapper.emitted('start')).toBeTruthy()
  })
})
```

## 部署设计

### Docker 配置

```dockerfile
# Dockerfile
FROM node:18-alpine as builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 环境配置

```typescript
// config/index.ts
export const config = {
  development: {
    apiBaseUrl: 'http://localhost:8888/api/v1',
    wsBaseUrl: 'ws://localhost:8888/ws'
  },
  production: {
    apiBaseUrl: '/api/v1',
    wsBaseUrl: `ws://${location.host}/ws`
  }
}
```
#
# 安全和权限控制设计

### 权限管理策略

Ollama模型管理系统使用框架提供的Casbin RBAC权限管理系统，不实现自定义的权限管理模块。所有API接口通过框架的`middleware.CasbinHandler()`中间件进行权限验证，操作日志通过框架的`middleware.OperationRecord()`中间件记录。

权限控制主要包括以下几个方面：

1. **API访问控制**：通过Casbin策略控制不同角色对API的访问权限
2. **操作审计**：使用框架的OperationRecord中间件记录用户操作
3. **用户认证**：使用框架的JWT认证机制
4. **数据权限**：通过CreatedBy和UpdatedBy字段记录数据所有者

### 安全防护措施

1. **输入验证**：对所有用户输入进行严格验证，防止注入攻击
2. **XSS防护**：对输出内容进行转义处理
3. **CSRF防护**：使用框架提供的CSRF防护机制
4. **敏感操作确认**：对删除等危险操作进行二次确认