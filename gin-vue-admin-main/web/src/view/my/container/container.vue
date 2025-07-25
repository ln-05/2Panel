<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  getContainerList, 
  getContainerDetail, 
  getContainerLogs,
  startContainer, 
  stopContainer, 
  restartContainer,
  deleteContainer,
  checkDockerStatus 
} from '@/api/container'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Refresh, 
  Plus, 
  Delete, 
  DeleteFilled, 
  Search,
  InfoFilled,
  VideoPlay,
  VideoPause
} from '@element-plus/icons-vue'

const activeTab = ref('all')
const tableData = ref([])
const selectedContainers = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const search = ref('')
const dockerStatus = ref(true)
const checkingStatus = ref(false)

// 根据activeTab过滤状态
const statusFilter = computed(() => {
  switch (activeTab.value) {
    case 'running':
      return 'running'
    case 'stopped':
      return 'exited'
    default:
      return ''
  }
})

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      status: statusFilter.value,
      name: search.value || ''
    }
    const res = await getContainerList(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    } else {
      // 根据不同的错误码显示不同的错误信息
      if (res.msg && res.msg.includes('权限不足')) {
        ElMessage.error('权限不足：请联系管理员分配Docker容器管理权限')
      } else if (res.msg && res.msg.includes('Docker')) {
        ElMessage.error('Docker服务不可用：请检查Docker是否正常运行')
      } else {
        ElMessage.error(res.msg || '获取容器列表失败')
      }
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取容器列表失败:', error)
    // 检查是否是网络错误或权限错误
    if (error.response && error.response.status === 403) {
      ElMessage.error('权限不足：请联系管理员分配Docker容器管理权限')
    } else if (error.response && error.response.status === 503) {
      ElMessage.error('Docker服务不可用：请检查Docker是否正常运行')
    } else {
      ElMessage.error('获取容器列表失败：请检查网络连接')
    }
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 检查Docker状态
const checkDocker = async () => {
  try {
    const res = await checkDockerStatus()
    dockerStatus.value = res.code === 0 && res.data === true
    if (!dockerStatus.value) {
      ElMessage.warning('Docker守护进程不可用')
    }
  } catch (error) {
    dockerStatus.value = false
    ElMessage.warning('无法连接到Docker守护进程')
  }
}

onMounted(() => {
  checkDocker()
  fetchData()
})

const handlePageChange = (val: number) => {
  page.value = val
  fetchData()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  page.value = 1
  fetchData()
}

const handleTabChange = () => {
  page.value = 1
  fetchData()
}

const handleSearch = () => {
  page.value = 1
  fetchData()
}

// 启动容器
const handleStart = async (row: any) => {
  try {
    const res = await startContainer(row.id)
    if (res.code === 0) {
      ElMessage.success('容器启动成功')
      fetchData()
    } else {
      ElMessage.error(res.msg || '启动容器失败')
    }
  } catch (error) {
    console.error('启动容器失败:', error)
    ElMessage.error('启动容器失败')
  }
}

// 停止容器
const handleStop = async (row: any) => {
  try {
    const res = await stopContainer(row.id)
    if (res.code === 0) {
      ElMessage.success('容器停止成功')
      fetchData()
    } else {
      ElMessage.error(res.msg || '停止容器失败')
    }
  } catch (error) {
    console.error('停止容器失败:', error)
    ElMessage.error('停止容器失败')
  }
}

// 重启容器
const handleRestart = async (row: any) => {
  try {
    const res = await restartContainer(row.id)
    if (res.code === 0) {
      ElMessage.success('容器重启成功')
      fetchData()
    } else {
      ElMessage.error(res.msg || '重启容器失败')
    }
  } catch (error) {
    console.error('重启容器失败:', error)
    ElMessage.error('重启容器失败')
  }
}

// 删除容器
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除容器 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    const res = await deleteContainer(row.id, false)
    if (res.code === 0) {
      ElMessage.success('容器删除成功')
      fetchData()
    } else {
      ElMessage.error(res.msg || '删除容器失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除容器失败:', error)
      ElMessage.error('删除容器失败')
    }
  }
}

// 查看日志
const handleLogs = async (row: any) => {
  try {
    const res = await getContainerLogs(row.id, { tail: '100' })
    if (res.code === 0) {
      // 格式化日志显示
      const logs = res.data || '暂无日志'
      const formattedLogs = logs.length > 5000 ? logs.substring(0, 5000) + '\n\n... (日志过长，仅显示前5000字符)' : logs
      
      ElMessageBox.alert(
        `<pre style="max-height: 400px; overflow-y: auto; white-space: pre-wrap; font-family: monospace; font-size: 12px; background: #f5f5f5; padding: 10px; border-radius: 4px;">${formattedLogs}</pre>`, 
        `容器 ${row.name || row.id?.substring(0, 12)} 的日志`, 
        {
          confirmButtonText: '确定',
          dangerouslyUseHTMLString: true,
          customClass: 'log-dialog'
        }
      )
    } else {
      ElMessage.error(res.msg || '获取日志失败')
    }
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error('获取日志失败')
  }
}

// 查看详情
const handleDetail = async (row: any) => {
  try {
    const res = await getContainerDetail(row.id)
    if (res.code === 0) {
      // 这里可以打开一个对话框显示详情
      console.log('容器详情:', res.data)
      ElMessage.info('详情已在控制台输出')
    } else {
      ElMessage.error(res.msg || '获取详情失败')
    }
  } catch (error) {
    console.error('获取详情失败:', error)
    ElMessage.error('获取详情失败')
  }
}

const handleAction = (type: string, row: any) => {
  switch (type) {
    case '启动':
      handleStart(row)
      break
    case '停止':
      handleStop(row)
      break
    case '重启':
      handleRestart(row)
      break
    case '删除':
      handleDelete(row)
      break
    case '日志':
      handleLogs(row)
      break
    case '详情':
      handleDetail(row)
      break
    default:
      console.log(type, row)
  }
}

// 格式化端口显示
const formatPorts = (ports: any[]) => {
  if (!ports || ports.length === 0) return '-'
  return ports.map(port => {
    if (port.publicPort) {
      return `${port.publicPort}:${port.privatePort}/${port.type}`
    }
    return `${port.privatePort}/${port.type}`
  }).join(', ')
}

// 格式化创建时间
const formatTime = (timestamp: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString()
}

// 判断容器是否运行中
const isRunning = (state: string) => {
  console.log('容器状态:', state) // 调试用
  return state === 'running'
}

// 获取资源使用率
const getResourceUsage = (row: any, type: string) => {
  // 目前后端API没有返回资源使用率数据，显示占位符
  // 后续可以通过调用容器详情API或者统计API来获取实际数据
  return '--'
}

// 获取容器IP地址
const getContainerIP = (row: any) => {
  // 目前后端API的基本信息中没有IP地址
  // 需要调用容器详情API才能获取网络信息
  // 这里先显示占位符，后续可以优化
  return '--'
}

// 获取网络信息
const getNetworkInfo = (row: any) => {
  // 目前后端API的基本信息中没有网络信息
  // 需要调用容器详情API才能获取网络信息
  // 这里先显示占位符，后续可以优化
  return '--'
}

// 创建容器
const handleCreateContainer = () => {
  ElMessage.info('创建容器功能暂未实现')
}

// 清理容器
const handleCleanContainers = () => {
  ElMessage.info('清理容器功能暂未实现')
}

// 选择处理
const handleSelectionChange = (selection: any[]) => {
  selectedContainers.value = selection
}

// 批量启动
const handleBatchStart = () => {
  if (selectedContainers.value.length === 0) {
    ElMessage.warning('请选择要启动的容器')
    return
  }

  ElMessageBox.confirm(
    `确定要启动选中的 ${selectedContainers.value.length} 个容器吗？`,
    '确认批量启动',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(async () => {
    try {
      const startPromises = selectedContainers.value.map(container => 
        startContainer(container.id)
      )
      
      await Promise.all(startPromises)
      ElMessage.success('批量启动成功')
      selectedContainers.value = []
      fetchData()
    } catch (error) {
      console.error('批量启动失败:', error)
      ElMessage.error('批量启动失败')
    }
  })
}

// 批量停止
const handleBatchStop = () => {
  if (selectedContainers.value.length === 0) {
    ElMessage.warning('请选择要停止的容器')
    return
  }

  ElMessageBox.confirm(
    `确定要停止选中的 ${selectedContainers.value.length} 个容器吗？`,
    '确认批量停止',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const stopPromises = selectedContainers.value.map(container => 
        stopContainer(container.id)
      )
      
      await Promise.all(stopPromises)
      ElMessage.success('批量停止成功')
      selectedContainers.value = []
      fetchData()
    } catch (error) {
      console.error('批量停止失败:', error)
      ElMessage.error('批量停止失败')
    }
  })
}

// 批量重启
const handleBatchRestart = () => {
  ElMessage.info('批量重启功能暂未实现')
}

// 批量强制停止
const handleBatchForceStop = () => {
  ElMessage.info('批量强制停止功能暂未实现')
}

// 批量暂停
const handleBatchPause = () => {
  ElMessage.info('批量暂停功能暂未实现')
}

// 批量删除
const handleBatchDelete = () => {
  if (selectedContainers.value.length === 0) {
    ElMessage.warning('请选择要删除的容器')
    return
  }

  ElMessageBox.confirm(
    `确定要删除选中的 ${selectedContainers.value.length} 个容器吗？`,
    '确认批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const deletePromises = selectedContainers.value.map(container => 
        deleteContainer(container.id, false)
      )
      
      await Promise.all(deletePromises)
      ElMessage.success('批量删除成功')
      selectedContainers.value = []
      fetchData()
    } catch (error) {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  })
}

// 批量恢复
const handleBatchResume = () => {
  ElMessage.info('批量恢复功能暂未实现')
}
</script>

<template>
  <div class="docker-container-management">
    <!-- 顶部操作按钮 -->
    <div class="operation-bar">
      <div class="left-buttons">
        <el-button type="primary" @click="handleCreateContainer">
          <el-icon><Plus /></el-icon>
          创建容器
        </el-button>
        <el-button @click="handleCleanContainers">
          <el-icon><DeleteFilled /></el-icon>
          清理容器
        </el-button>
        <el-button @click="fetchData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="checkDocker" :loading="checkingStatus">
          <el-icon><InfoFilled /></el-icon>
          检查Docker状态
        </el-button>
        <el-button 
          v-if="selectedContainers.length > 0" 
          type="success" 
          @click="handleBatchStart"
        >
          <el-icon><VideoPlay /></el-icon>
          批量启动
        </el-button>
        <el-button 
          v-if="selectedContainers.length > 0" 
          type="warning" 
          @click="handleBatchStop"
        >
          <el-icon><VideoPause /></el-icon>
          批量停止
        </el-button>
        <el-button 
          v-if="selectedContainers.length > 0" 
          type="danger" 
          @click="handleBatchDelete"
        >
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
      </div>
      <div class="right-search">
        <el-input
          v-model="search"
          placeholder="搜索容器名称"
          style="width: 300px"
          @input="handleSearch"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>
    <!-- 容器列表表格 -->
    <div class="table-container">
      <el-table 
        v-loading="loading"
        :data="tableData" 
        @selection-change="handleSelectionChange"
        style="width: 100%"
        height="calc(100vh - 200px)"
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column label="名称" min-width="140">
          <template #default="scope">
            <span>{{ scope.row.name || scope.row.id?.substring(0, 12) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="镜像" min-width="160">
          <template #default="scope">
            <span>{{ scope.row.image || '-' }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="状态" min-width="90">
          <template #default="scope">
            <el-tag v-if="isRunning(scope.row.state)" class="status-tag running">运行中</el-tag>
            <el-tag v-else class="status-tag stopped">已停止</el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="资源使用率" min-width="120">
          <template #default="scope">
            <div>CPU: <span class="resource-num">{{ getResourceUsage(scope.row, 'cpu') }}</span></div>
            <div>内存: <span class="resource-num">{{ getResourceUsage(scope.row, 'memory') }}</span></div>
          </template>
        </el-table-column>
        
        <el-table-column label="IP地址" min-width="120">
          <template #default="scope">
            <span>{{ getContainerIP(scope.row) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="关联资源" min-width="120">
          <template #default="scope">
            <span>{{ getNetworkInfo(scope.row) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="端口" min-width="160">
          <template #default="scope">
            <span>{{ formatPorts(scope.row.ports) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="创建时间" min-width="120">
          <template #default="scope">
            <span>{{ formatTime(scope.row.created) }}</span>
          </template>
        </el-table-column>
      <el-table-column label="操作" min-width="240" fixed="right">
        <template #default="scope">
          <div style="display: flex; gap: 8px; flex-wrap: wrap;">
            <!-- 临时显示所有按钮，用于调试 -->
            <el-button 
              v-if="!isRunning(scope.row.state)" 
              link 
              type="success" 
              size="small" 
              @click="handleAction('启动', scope.row)"
            >
              启动
            </el-button>
            <el-button 
              v-if="isRunning(scope.row.state)" 
              link 
              type="warning" 
              size="small" 
              @click="handleAction('停止', scope.row)"
            >
              停止
            </el-button>
            <!-- 临时显示所有按钮进行测试 -->
            <el-button 
              link 
              type="success" 
              size="small" 
              @click="handleAction('启动', scope.row)"
            >
              测试启动
            </el-button>
            <el-button 
              link 
              type="warning" 
              size="small" 
              @click="handleAction('停止', scope.row)"
            >
              测试停止
            </el-button>
            <el-button 
              link 
              type="primary" 
              size="small" 
              @click="handleAction('重启', scope.row)"
            >
              重启
            </el-button>
            <el-button 
              link 
              type="primary" 
              size="small" 
              @click="handleAction('日志', scope.row)"
            >
              日志
            </el-button>
            <el-button 
              link 
              type="info" 
              size="small" 
              @click="handleAction('详情', scope.row)"
            >
              详情
            </el-button>
            <el-button 
              link 
              type="danger" 
              size="small" 
              @click="handleAction('删除', scope.row)"
            >
              删除
            </el-button>
            <!-- 调试信息 -->
            <span style="font-size: 10px; color: #999;">{{ scope.row.state }}</span>
          </div>
        </template>
      </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.docker-container-management {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.operation-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.left-buttons {
  display: flex;
  gap: 12px;
}

.left-buttons .el-button {
  border-radius: 6px;
}

.right-search {
  display: flex;
  align-items: center;
}

.table-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.el-table {
  border-radius: 8px;
}

.el-table th {
  background-color: #fafafa;
  color: #333;
  font-weight: 600;
}

.el-table td {
  border-bottom: 1px solid #f0f0f0;
}

.status-tag {
  border-radius: 16px;
  font-size: 13px;
  padding: 0 16px;
  border: none;
}

.status-tag.running {
  background: linear-gradient(90deg, #00e4ff 0%, #00bfff 100%);
  color: #fff;
}

.status-tag.stopped {
  background: #ff4d4f;
  color: #fff;
}

.resource-num {
  color: #00e4ff;
  font-weight: bold;
}

.pagination-container {
  padding: 20px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #f0f0f0;
}

/* 表格操作按钮样式 */
.el-table .el-button--small {
  padding: 4px 8px;
  font-size: 12px;
}

.el-table .el-button + .el-button {
  margin-left: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .docker-container-management {
    padding: 10px;
  }
  
  .operation-bar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .left-buttons {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .right-search {
    justify-content: center;
  }
  
  .right-search .el-input {
    width: 100% !important;
  }
  
  .table-container .el-table {
    font-size: 12px;
  }
  
  .el-table .el-button--small {
    padding: 2px 6px;
    font-size: 11px;
  }
}

/* 加载状态样式 */
.el-loading-mask {
  background-color: rgba(255, 255, 255, 0.8);
}

/* 标签样式优化 */
.el-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}

.el-tag--small {
  height: 20px;
  line-height: 18px;
  font-size: 11px;
}

/* 表格行悬停效果 */
.el-table tbody tr:hover > td {
  background-color: #f5f7fa !important;
}

/* 选择框样式 */
.el-table .el-checkbox {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 工具提示样式 */
.el-tooltip__popper {
  max-width: 300px;
}

/* 搜索框样式 */
.el-input__prefix {
  display: flex;
  align-items: center;
}

/* 分页样式 */
.el-pagination {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.el-pagination .el-pagination__total {
  margin-right: auto;
}

/* 空状态样式 */
.el-table__empty-block {
  padding: 60px 0;
}

.el-table__empty-text {
  color: #909399;
  font-size: 14px;
}

/* 消息提示样式优化 */
.el-message {
  min-width: 300px;
  border-radius: 6px;
}

/* 确认对话框样式 */
.el-message-box {
  border-radius: 8px;
}

.el-message-box__header {
  padding: 20px 20px 10px;
}

.el-message-box__content {
  padding: 10px 20px;
}

.el-message-box__btns {
  padding: 10px 20px 20px;
}
.header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
  .tabs {
    flex: 1;
    .el-tabs__item {
      color: #bfcbd9;
      &.is-active {
        color: #00e4ff;
        font-weight: bold;
      }
    }
  }
  .header-actions {
    display: flex;
    gap: 10px;
    margin-left: 24px;
    .main-btn {
      font-weight: bold;
      border-radius: 20px;
      padding: 0 22px;
      background: linear-gradient(90deg, #00e4ff 0%, #00bfff 100%);
      color: #fff;
      border: none;
      box-shadow: 0 2px 8px 0 rgba(0,228,255,0.10);
    }
    .minor-btn {
      border-radius: 18px;
      background: #232a34;
      color: #bfcbd9;
      border: 1px solid #232a34;
      transition: background 0.2s;
      &:hover {
        background: #232a34cc;
        color: #00e4ff;
      }
    }
  }
  .header-search {
    display: flex;
    align-items: center;
    gap: 8px;
    .search-input {
      width: 180px;
      border-radius: 16px;
      background: #232a34;
      color: #fff;
      border: 1px solid #232a34;
    }
    .icon-btn {
      background: #232a34;
      color: #bfcbd9;
      border: none;
      &:hover {
        color: #00e4ff;
      }
    }
  }
}
.container-table {
  background: #232a34;
  border-radius: 10px;
  box-shadow: 0 2px 8px 0 rgba(0,0,0,0.08);
  .el-table__header th {
    background: #232a34;
    color: #bfcbd9;
    font-weight: 500;
    border-bottom: 1.5px solid #222b36;
    box-shadow: 0 2px 4px 0 rgba(0,0,0,0.04);
  }
  .el-table__row {
    background: #232a34;
    color: #fff;
    transition: background 0.2s;
    &:hover {
      background: #26303d;
    }
  }
  .el-tag.status-tag {
    border-radius: 16px;
    font-size: 13px;
    padding: 0 16px;
    border: none;
    &.running {
      background: linear-gradient(90deg, #00e4ff 0%, #00bfff 100%);
      color: #fff;
    }
    &.stopped {
      background: #ff4d4f;
      color: #fff;
    }
  }
  .resource-num {
    color: #00e4ff;
    font-weight: bold;
  }
}
.table-footer {
  display: flex;
  align-items: center;
  margin-top: 18px;
  color: #bfcbd9;
  .table-total {
    font-size: 15px;
    margin-right: 18px;
  }
  .table-pagination {
    margin-left: auto;
    .el-pagination__total {
      color: #bfcbd9;
    }
  }
}
</style>

