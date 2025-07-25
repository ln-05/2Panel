<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Refresh, 
  Plus, 
  Delete, 
  DeleteFilled, 
  Search,
  InfoFilled,
  View
} from '@element-plus/icons-vue'
import {
  getDockerVolumeList,
  getDockerVolumeDetail,
  createDockerVolume,
  deleteDockerVolume,
  pruneDockerVolumes
} from '@/api/dockerVolume'

defineOptions({
  name: 'DockerVolumeManagement'
})

// 页面状态
const loading = ref(false)
const volumeList = ref([])
const selectedVolumes = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框状态
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const volumeDetail = ref({})

// 创建存储卷表单
const createForm = reactive({
  name: '',
  driver: 'local',
  driverOpts: {},
  labels: {}
})

const createRules = {
  name: [
    { required: true, message: '请输入存储卷名称', trigger: 'blur' }
  ]
}

// 表单引用
const createFormRef = ref()

// 获取存储卷列表
const fetchVolumeList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.name = searchKeyword.value
    }

    const response = await getDockerVolumeList(params)
    
    if (response.code === 0) {
      volumeList.value = response.data.list || []
      total.value = response.data.total || 0
    } else {
      let errorMessage = response.msg || '获取存储卷列表失败'
      
      if (response.code === 7) {
        errorMessage = '请先登录系统'
      } else if (response.msg && response.msg.includes('Docker client is not available')) {
        errorMessage = 'Docker服务不可用，请检查Docker配置和连接'
      } else if (response.msg && response.msg.includes('connection refused')) {
        errorMessage = '无法连接到Docker服务，请检查Docker守护进程是否运行'
      } else if (response.msg && response.msg.includes('timeout')) {
        errorMessage = 'Docker服务连接超时，请检查网络连接'
      }
      
      ElMessage.error(errorMessage)
    }
  } catch (error) {
    console.error('获取存储卷列表失败:', error)
    
    let errorMessage = '获取存储卷列表失败'
    
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        errorMessage = '认证失败，请重新登录'
      } else if (status === 403) {
        errorMessage = '权限不足，无法访问Docker功能'
      } else if (status === 500) {
        errorMessage = '服务器内部错误，请检查Docker服务配置'
      } else if (status === 502 || status === 503) {
        errorMessage = 'Docker服务不可用'
      }
    } else if (error.request) {
      errorMessage = '网络连接失败，请检查网络连接'
    }
    
    ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}

// 刷新存储卷列表
const refreshVolumeList = () => {
  fetchVolumeList()
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchVolumeList()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchVolumeList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchVolumeList()
}

// 选择处理
const handleSelectionChange = (selection) => {
  selectedVolumes.value = selection
}

// 创建存储卷
const handleCreateVolume = async () => {
  if (!createFormRef.value) return
  
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    const response = await createDockerVolume(createForm)

    if (response.code === 0) {
      ElMessage.success('存储卷创建成功')
      showCreateDialog.value = false
      resetCreateForm()
      fetchVolumeList()
    } else {
      let errorMessage = response.msg || '创建存储卷失败'
      
      if (response.code === 7) {
        errorMessage = '请先登录系统'
      } else if (response.msg && response.msg.includes('Docker client is not available')) {
        errorMessage = 'Docker服务不可用，请检查Docker配置'
      } else if (response.msg && response.msg.includes('already exists')) {
        errorMessage = '存储卷名称已存在，请使用其他名称'
      }
      
      ElMessage.error(errorMessage)
    }
  } catch (error) {
    console.error('创建存储卷失败:', error)
    
    let errorMessage = '创建存储卷失败'
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        errorMessage = '认证失败，请重新登录'
      } else if (status === 500) {
        errorMessage = '服务器错误，请检查Docker服务状态'
      }
    } else if (error.request) {
      errorMessage = '网络连接失败，请检查网络连接'
    }
    
    ElMessage.error(errorMessage)
  }
}

// 查看存储卷详情
const showVolumeDetail = async (row) => {
  try {
    const response = await getDockerVolumeDetail(row.name)
    
    if (response.code === 0) {
      volumeDetail.value = response.data
      showDetailDialog.value = true
    } else {
      ElMessage.error(response.msg || '获取存储卷详情失败')
    }
  } catch (error) {
    console.error('获取存储卷详情失败:', error)
    ElMessage.error('获取存储卷详情失败')
  }
}

// 删除存储卷
const deleteVolume = (row) => {
  ElMessageBox.confirm(
    `确定要删除存储卷 ${row.name} 吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteDockerVolume(row.name)
      
      if (response.code === 0) {
        ElMessage.success('存储卷删除成功')
        fetchVolumeList()
      } else {
        ElMessage.error(response.msg || '删除存储卷失败')
      }
    } catch (error) {
      console.error('删除存储卷失败:', error)
      ElMessage.error('删除存储卷失败')
    }
  })
}

// 批量删除存储卷
const batchDeleteVolumes = () => {
  if (selectedVolumes.value.length === 0) {
    ElMessage.warning('请选择要删除的存储卷')
    return
  }

  ElMessageBox.confirm(
    `确定要删除选中的 ${selectedVolumes.value.length} 个存储卷吗？`,
    '确认批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const deletePromises = selectedVolumes.value.map(volume => 
        deleteDockerVolume(volume.name)
      )
      
      await Promise.all(deletePromises)
      ElMessage.success('批量删除成功')
      selectedVolumes.value = []
      fetchVolumeList()
    } catch (error) {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  })
}

// 清理未使用的存储卷
const pruneVolumes = () => {
  ElMessageBox.confirm(
    '确定要清理所有未使用的存储卷吗？此操作不可恢复。',
    '确认清理',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await pruneDockerVolumes()
      
      if (response.code === 0) {
        const { deletedCount, spaceReclaimed } = response.data
        ElMessage.success(`清理完成，删除了 ${deletedCount} 个存储卷，释放了 ${formatBytes(spaceReclaimed)} 空间`)
        fetchVolumeList()
      } else {
        ElMessage.error(response.msg || '清理存储卷失败')
      }
    } catch (error) {
      console.error('清理存储卷失败:', error)
      ElMessage.error('清理存储卷失败')
    }
  })
}

// 工具函数
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

const formatBytes = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getDriverTagType = (driver) => {
  switch (driver) {
    case 'local':
      return 'primary'
    case 'nfs':
      return 'success'
    case 'tmpfs':
      return 'warning'
    default:
      return 'info'
  }
}

// 重置表单
const resetCreateForm = () => {
  Object.assign(createForm, {
    name: '',
    driver: 'local',
    driverOpts: {},
    labels: {}
  })
}

// 页面加载时获取数据
onMounted(() => {
  fetchVolumeList()
})
</script>

<template>
  <div class="docker-volume-management">
    <!-- 顶部操作按钮 -->
    <div class="operation-bar">
      <div class="left-buttons">
        <el-button type="primary" @click="showCreateDialog = true">
          创建存储卷
        </el-button>
        <el-button @click="pruneVolumes">
          清理存储卷
        </el-button>
        <el-button 
          v-if="selectedVolumes.length > 0" 
          type="danger" 
          @click="batchDeleteVolumes"
        >
          删除
        </el-button>
      </div>
      <div class="right-search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索"
          style="width: 200px"
          @input="handleSearch"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button @click="refreshVolumeList" style="margin-left: 8px">
          <el-icon><Refresh /></el-icon>
        </el-button>
        <el-button style="margin-left: 8px">
          不刷新
        </el-button>
      </div>
    </div>

    <!-- 存储卷列表表格 -->
    <div class="table-container">
      <el-table
        v-loading="loading"
        :data="volumeList"
        @selection-change="handleSelectionChange"
        style="width: 100%"
        border
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column label="名称" min-width="200">
          <template #default="{ row }">
            <div class="volume-name">
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="存储卷目录" min-width="300">
          <template #default="{ row }">
            <span class="mountpoint">{{ row.mountpoint || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="挂载点" min-width="300">
          <template #default="{ row }">
            <span class="mountpoint">{{ row.mountpoint || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="驱动" width="100">
          <template #default="{ row }">
            <span>{{ row.driver || 'local' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="deleteVolume(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <div class="pagination-info">
          共 {{ total }} 条
        </div>
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 创建存储卷对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建存储卷"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form 
        :model="createForm" 
        :rules="createRules" 
        ref="createFormRef" 
        label-width="120px"
      >
        <el-form-item label="存储卷名称" prop="name">
          <el-input
            v-model="createForm.name"
            placeholder="请输入存储卷名称"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="驱动类型" prop="driver">
          <el-select v-model="createForm.driver" style="width: 100%">
            <el-option label="local" value="local" />
            <el-option label="nfs" value="nfs" />
            <el-option label="tmpfs" value="tmpfs" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateVolume">
          创建
        </el-button>
      </template>
    </el-dialog>

    <!-- 存储卷详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="存储卷详情"
      width="800px"
    >
      <div v-if="volumeDetail" class="volume-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="存储卷名称">
            {{ volumeDetail.name }}
          </el-descriptions-item>
          <el-descriptions-item label="驱动类型">
            <el-tag :type="getDriverTagType(volumeDetail.driver)">
              {{ volumeDetail.driver }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="挂载点">
            {{ volumeDetail.mountpoint }}
          </el-descriptions-item>
          <el-descriptions-item label="范围">
            {{ volumeDetail.scope }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(volumeDetail.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="volumeDetail.usageData" label="大小">
            {{ formatBytes(volumeDetail.usageData.size) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="volumeDetail.usageData" label="引用计数">
            {{ volumeDetail.usageData.refCount }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="volumeDetail.labels && Object.keys(volumeDetail.labels).length > 0" class="labels-section">
          <h4>标签</h4>
          <div class="labels-container">
            <el-tag 
              v-for="(value, key) in volumeDetail.labels" 
              :key="key" 
              class="label-tag"
            >
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>

        <div v-if="volumeDetail.options && Object.keys(volumeDetail.options).length > 0" class="options-section">
          <h4>选项</h4>
          <div class="options-container">
            <el-tag 
              v-for="(value, key) in volumeDetail.options" 
              :key="key" 
              class="option-tag"
              type="info"
            >
              {{ key }}={{ value }}
            </el-tag>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
.docker-volume-management {
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

.volume-name {
  font-weight: 500;
  color: #333;
}

.mountpoint {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #666;
}

.pagination-container {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #f0f0f0;
}

.pagination-info {
  color: #666;
  font-size: 14px;
}

/* 对话框样式 */
.el-dialog__header {
  padding: 20px 20px 10px;
  border-bottom: 1px solid #f0f0f0;
}

.el-dialog__body {
  padding: 20px;
}

.el-dialog__footer {
  padding: 10px 20px 20px;
  border-top: 1px solid #f0f0f0;
}

/* 存储卷详情样式 */
.volume-detail {
  .labels-section,
  .options-section {
    margin-top: 24px;
    
    h4 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: #333;
      border-bottom: 2px solid #409eff;
      padding-bottom: 8px;
    }
  }
  
  .labels-container,
  .options-container {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    
    .label-tag,
    .option-tag {
      margin: 0;
    }
  }
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
  .docker-volume-management {
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

/* 对话框表单样式 */
.el-form-item {
  margin-bottom: 20px;
}

.el-form-item__label {
  font-weight: 500;
  color: #333;
}

/* 按钮组样式 */
.el-button-group {
  display: inline-flex;
}

.el-button-group .el-button {
  margin-left: 0;
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

/* 描述列表样式 */
.el-descriptions {
  margin-bottom: 16px;
}

.el-descriptions__label {
  font-weight: 600;
  color: #333;
}
</style>