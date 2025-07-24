<template>
  <div class="ollama-model-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>Ollama模型管理</h2>
        <div class="connection-status">
          <el-tag 
            :type="connectionStatusType" 
            size="small"
            :icon="connectionStatusIcon"
          >
            {{ connectionStatusText }}
          </el-tag>
        </div>
      </div>
      
      <div class="header-right">
        <el-button-group>
          <el-button 
            type="primary" 
            :icon="Plus" 
            @click="showAddDialog = true"
            :loading="loading"
          >
            添加模型
          </el-button>
          <el-button 
            :icon="Refresh" 
            @click="handleSync"
            :loading="syncLoading"
          >
            同步模型
          </el-button>
          <el-button 
            :icon="Monitor" 
            @click="showChatDialog = true"
          >
            AI对话
          </el-button>
        </el-button-group>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索模型名称..."
        style="width: 300px"
        clearable
        @input="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      
      <el-select
        v-model="statusFilter"
        placeholder="状态筛选"
        style="width: 150px; margin-left: 10px"
        clearable
        @change="handleSearch"
      >
        <el-option label="运行中" value="running" />
        <el-option label="已停止" value="stopped" />
        <el-option label="下载中" value="downloading" />
        <el-option label="错误" value="error" />
        <el-option label="不可用" value="unavailable" />
      </el-select>
    </div>

    <!-- 模型列表 -->
    <div class="model-list">
      <el-table 
        :data="modelList" 
        :loading="tableLoading"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="name" label="模型名称" min-width="200">
          <template #default="{ row }">
            <div class="model-name">
              <span>{{ row.name }}</span>
              <el-tag v-if="row.from" size="small" type="info" style="margin-left: 8px">
                {{ row.from }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="size" label="大小" width="120" />
        
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="message" label="消息" min-width="200" show-overflow-tooltip />
        
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                v-if="row.status === 'stopped'"
                type="success"
                size="small"
                :icon="VideoPlay"
                @click="handleStartModel(row)"
                :loading="row.loading"
              >
                启动
              </el-button>
              
              <el-button
                v-if="row.status === 'running'"
                type="warning"
                size="small"
                :icon="VideoPause"
                @click="handleStopModel(row)"
                :loading="row.loading"
              >
                停止
              </el-button>
              
              <el-button
                v-if="row.status === 'running'"
                type="primary"
                size="small"
                :icon="ChatDotRound"
                @click="handleChatWithModel(row)"
              >
                对话
              </el-button>
              
              <el-button
                type="info"
                size="small"
                :icon="Refresh"
                @click="handleRecreateModel(row)"
                :loading="row.loading"
              >
                重建
              </el-button>
              
              <el-button
                type="danger"
                size="small"
                :icon="Delete"
                @click="handleDeleteModel(row)"
                :loading="row.loading"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 批量操作 -->
      <div v-if="selectedModels.length > 0" class="batch-operations">
        <el-alert
          :title="`已选择 ${selectedModels.length} 个模型`"
          type="info"
          show-icon
          :closable="false"
        >
          <template #default>
            <el-button-group>
              <el-button size="small" @click="handleBatchStart">批量启动</el-button>
              <el-button size="small" @click="handleBatchStop">批量停止</el-button>
              <el-button size="small" type="danger" @click="handleBatchDelete">批量删除</el-button>
            </el-button-group>
          </template>
        </el-alert>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handlePageSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 添加模型对话框 -->
    <AddModelDialog 
      v-model="showAddDialog" 
      @success="refreshModelList"
    />

    <!-- AI对话对话框 -->
    <ChatDialog 
      v-model="showChatDialog"
      :model="selectedChatModel"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, Refresh, Monitor, Search, VideoPlay, VideoPause, 
  Delete, ChatDotRound, Connection, Disconnect 
} from '@element-plus/icons-vue'

// API imports
import {
  searchOllamaModel,
  startOllamaModel,
  stopOllamaModel,
  deleteOllamaModel,
  recreateOllamaModel,
  syncOllamaModel,
  deleteOllamaModelByIds
} from '@/api/AI/ollamaModel'

// Composables
import { useOllamaWebSocket } from '@/composables/useOllamaWebSocket'

// Config
import { formatModelStatus, formatConnectionStatus } from '@/config/ollama'

// Components
import AddModelDialog from './components/AddModelDialog.vue'
import ChatDialog from './components/ChatDialog.vue'

// 响应式数据
const loading = ref(false)
const tableLoading = ref(false)
const syncLoading = ref(false)
const modelList = ref([])
const selectedModels = ref([])
const selectedChatModel = ref(null)

// 搜索和筛选
const searchQuery = ref('')
const statusFilter = ref('')

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框状态
const showAddDialog = ref(false)
const showChatDialog = ref(false)

// WebSocket连接
const { isConnected, connectionStatus, onMessage } = useOllamaWebSocket()

// 计算属性
const connectionStatusType = computed(() => {
  return formatConnectionStatus(connectionStatus.value).type
})

const connectionStatusText = computed(() => {
  return formatConnectionStatus(connectionStatus.value).text
})

const connectionStatusIcon = computed(() => {
  return isConnected.value ? Connection : Disconnect
})

// 方法
const refreshModelList = async () => {
  tableLoading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    if (searchQuery.value) {
      params.name = searchQuery.value
    }
    
    if (statusFilter.value) {
      params.status = [statusFilter.value]
    }

    const response = await searchOllamaModel(params)
    if (response.code === 0) {
      modelList.value = response.data.list.map(item => ({
        ...item,
        loading: false
      }))
      total.value = response.data.total
    }
  } catch (error) {
    ElMessage.error('获取模型列表失败: ' + error.message)
  } finally {
    tableLoading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  refreshModelList()
}

const handleSync = async () => {
  syncLoading.value = true
  try {
    const response = await syncOllamaModel({ force: false })
    if (response.code === 0) {
      const result = response.data
      ElMessage.success(`同步完成: 新增${result.added}个，更新${result.updated}个，移除${result.removed}个`)
      await refreshModelList()
    }
  } catch (error) {
    ElMessage.error('同步失败: ' + error.message)
  } finally {
    syncLoading.value = false
  }
}

const handleStartModel = async (model) => {
  model.loading = true
  try {
    const response = await startOllamaModel(model.id)
    if (response.code === 0) {
      ElMessage.success(`模型 ${model.name} 启动成功`)
      model.status = 'running'
      model.message = '模型已启动'
    }
  } catch (error) {
    ElMessage.error(`启动模型失败: ${error.message}`)
  } finally {
    model.loading = false
  }
}

const handleStopModel = async (model) => {
  model.loading = true
  try {
    const response = await stopOllamaModel(model.id)
    if (response.code === 0) {
      ElMessage.success(`模型 ${model.name} 停止成功`)
      model.status = 'stopped'
      model.message = '模型已停止'
    }
  } catch (error) {
    ElMessage.error(`停止模型失败: ${error.message}`)
  } finally {
    model.loading = false
  }
}

const handleDeleteModel = async (model) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除模型 "${model.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    model.loading = true
    const response = await deleteOllamaModel(model.id)
    if (response.code === 0) {
      ElMessage.success(`模型 ${model.name} 删除成功`)
      await refreshModelList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`删除模型失败: ${error.message}`)
    }
  } finally {
    model.loading = false
  }
}

const handleRecreateModel = async (model) => {
  try {
    await ElMessageBox.confirm(
      `确定要重建模型 "${model.name}" 吗？这将重新下载模型文件。`,
      '确认重建',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    model.loading = true
    const response = await recreateOllamaModel(model.id)
    if (response.code === 0) {
      ElMessage.success(`模型 ${model.name} 重建任务已启动`)
      model.status = 'downloading'
      model.message = '重新创建模型...'
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`重建模型失败: ${error.message}`)
    }
  } finally {
    model.loading = false
  }
}

const handleChatWithModel = (model) => {
  selectedChatModel.value = model
  showChatDialog.value = true
}

const handleSelectionChange = (selection) => {
  selectedModels.value = selection
}

const handleBatchStart = async () => {
  const stoppedModels = selectedModels.value.filter(m => m.status === 'stopped')
  if (stoppedModels.length === 0) {
    ElMessage.warning('没有可启动的模型')
    return
  }

  for (const model of stoppedModels) {
    await handleStartModel(model)
  }
}

const handleBatchStop = async () => {
  const runningModels = selectedModels.value.filter(m => m.status === 'running')
  if (runningModels.length === 0) {
    ElMessage.warning('没有可停止的模型')
    return
  }

  for (const model of runningModels) {
    await handleStopModel(model)
  }
}

const handleBatchDelete = async () => {
  if (selectedModels.value.length === 0) {
    ElMessage.warning('请选择要删除的模型')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedModels.value.length} 个模型吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const ids = selectedModels.value.map(m => m.id.toString())
    const response = await deleteOllamaModelByIds(ids)
    if (response.code === 0) {
      ElMessage.success(`成功删除 ${selectedModels.value.length} 个模型`)
      selectedModels.value = []
      await refreshModelList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`批量删除失败: ${error.message}`)
    }
  }
}

const handlePageSizeChange = (newPageSize) => {
  pageSize.value = newPageSize
  currentPage.value = 1
  refreshModelList()
}

const handlePageChange = (newPage) => {
  currentPage.value = newPage
  refreshModelList()
}

// 工具函数
const getStatusType = (status) => {
  return formatModelStatus(status).type
}

const getStatusText = (status) => {
  return formatModelStatus(status).text
}

const formatTime = (timeStr) => {
  if (!timeStr) return '-'
  return new Date(timeStr).toLocaleString('zh-CN')
}

// WebSocket消息处理
onMessage('model_status_update', (data) => {
  const model = modelList.value.find(m => m.id === data.modelId)
  if (model) {
    model.status = data.status
    model.message = data.message || ''
  }
})

onMessage('download_progress', (data) => {
  const model = modelList.value.find(m => m.name === data.modelName)
  if (model) {
    model.message = `下载进度: ${data.progress}%`
  }
})

// 生命周期
onMounted(() => {
  refreshModelList()
})
</script>

<style scoped>
.ollama-model-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-left h2 {
  margin: 0;
  color: #303133;
}

.connection-status {
  display: flex;
  align-items: center;
}

.search-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.model-list {
  margin-bottom: 20px;
}

.model-name {
  display: flex;
  align-items: center;
}

.action-buttons {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.batch-operations {
  margin-top: 15px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>