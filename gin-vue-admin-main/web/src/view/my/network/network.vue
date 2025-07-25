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
  getDockerNetworkList,
  getDockerNetworkDetail,
  createDockerNetwork,
  deleteDockerNetwork,
  pruneDockerNetworks
} from '@/api/dockerNetwork'

defineOptions({
  name: 'DockerNetworkManagement'
})

// 页面状态
const loading = ref(false)
const networkList = ref([])
const selectedNetworks = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框状态
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const networkDetail = ref({})

// 创建网络表单
const createForm = reactive({
  name: '',
  driver: 'bridge',
  subnet: '',
  gateway: '',
  enableIPv6: false,
  internal: false,
  attachable: false,
  labels: {}
})

const createRules = {
  name: [
    { required: true, message: '请输入网络名称', trigger: 'blur' }
  ]
}

// 表单引用
const createFormRef = ref()

// 获取网络列表
const fetchNetworkList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.name = searchKeyword.value
    }

    const response = await getDockerNetworkList(params)
    
    if (response.code === 0) {
      networkList.value = response.data.list || []
      total.value = response.data.total || 0
      // 调试信息：查看返回的数据结构
      console.log('Network list data:', response.data.list)
    } else {
      // 根据错误码提供更具体的错误信息
      let errorMessage = response.msg || '获取网络列表失败'
      
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
    console.error('获取网络列表失败:', error)
    
    // 处理网络错误
    let errorMessage = '获取网络列表失败'
    
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

// 刷新网络列表
const refreshNetworkList = () => {
  fetchNetworkList()
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchNetworkList()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchNetworkList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchNetworkList()
}

// 选择处理
const handleSelectionChange = (selection) => {
  selectedNetworks.value = selection
}

// 创建网络
const handleCreateNetwork = async () => {
  if (!createFormRef.value) return
  
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    const response = await createDockerNetwork(createForm)

    if (response.code === 0) {
      ElMessage.success('网络创建成功')
      showCreateDialog.value = false
      resetCreateForm()
      fetchNetworkList()
    } else {
      let errorMessage = response.msg || '创建网络失败'
      
      if (response.code === 7) {
        errorMessage = '请先登录系统'
      } else if (response.msg && response.msg.includes('Docker client is not available')) {
        errorMessage = 'Docker服务不可用，请检查Docker配置'
      } else if (response.msg && response.msg.includes('already exists')) {
        errorMessage = '网络名称已存在，请使用其他名称'
      }
      
      ElMessage.error(errorMessage)
    }
  } catch (error) {
    console.error('创建网络失败:', error)
    
    let errorMessage = '创建网络失败'
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

// 查看网络详情
const showNetworkDetail = async (row) => {
  try {
    const response = await getDockerNetworkDetail(row.id)
    
    if (response.code === 0) {
      networkDetail.value = response.data
      showDetailDialog.value = true
    } else {
      ElMessage.error(response.msg || '获取网络详情失败')
    }
  } catch (error) {
    console.error('获取网络详情失败:', error)
    ElMessage.error('获取网络详情失败')
  }
}

// 删除网络
const deleteNetwork = (row) => {
  ElMessageBox.confirm(
    `确定要删除网络 ${row.name} 吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteDockerNetwork(row.id)
      
      if (response.code === 0) {
        ElMessage.success('网络删除成功')
        fetchNetworkList()
      } else {
        ElMessage.error(response.msg || '删除网络失败')
      }
    } catch (error) {
      console.error('删除网络失败:', error)
      ElMessage.error('删除网络失败')
    }
  })
}

// 批量删除网络
const batchDeleteNetworks = () => {
  if (selectedNetworks.value.length === 0) {
    ElMessage.warning('请选择要删除的网络')
    return
  }

  ElMessageBox.confirm(
    `确定要删除选中的 ${selectedNetworks.value.length} 个网络吗？`,
    '确认批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const deletePromises = selectedNetworks.value.map(network => 
        deleteDockerNetwork(network.id)
      )
      
      await Promise.all(deletePromises)
      ElMessage.success('批量删除成功')
      selectedNetworks.value = []
      fetchNetworkList()
    } catch (error) {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  })
}

// 清理未使用的网络
const pruneNetworks = () => {
  ElMessageBox.confirm(
    '确定要清理所有未使用的网络吗？此操作不可恢复。',
    '确认清理',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await pruneDockerNetworks()
      
      if (response.code === 0) {
        const { deletedCount } = response.data
        ElMessage.success(`清理完成，删除了 ${deletedCount} 个网络`)
        fetchNetworkList()
      } else {
        ElMessage.error(response.msg || '清理网络失败')
      }
    } catch (error) {
      console.error('清理网络失败:', error)
      ElMessage.error('清理网络失败')
    }
  })
}

// 工具函数
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const formatLabels = (labels) => {
  if (!labels || Object.keys(labels).length === 0) return '-'
  return Object.entries(labels).map(([key, value]) => `${key}=${value}`).join(', ')
}

const getDriverTagType = (driver) => {
  switch (driver) {
    case 'bridge':
      return 'primary'
    case 'host':
      return 'success'
    case 'overlay':
      return 'warning'
    case 'macvlan':
      return 'info'
    case 'none':
      return 'danger'
    default:
      return 'info'
  }
}

// 重置表单
const resetCreateForm = () => {
  Object.assign(createForm, {
    name: '',
    driver: 'bridge',
    subnet: '',
    gateway: '',
    enableIPv6: false,
    internal: false,
    attachable: false,
    labels: {}
  })
}

// 页面加载时获取数据
onMounted(() => {
  fetchNetworkList()
})
</script>

<template>
  <div class="docker-network-management">
    <!-- 顶部操作按钮 -->
    <div class="operation-bar">
      <div class="left-buttons">
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          创建网络
        </el-button>
        <el-button @click="refreshNetworkList">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button 
          v-if="selectedNetworks.length > 0" 
          type="danger" 
          @click="batchDeleteNetworks"
        >
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
        <el-button @click="pruneNetworks">
          <el-icon><DeleteFilled /></el-icon>
          清理未使用网络
        </el-button>
      </div>
      <div class="right-search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索网络名称"
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

    <!-- 网络列表表格 -->
    <div class="table-container">
      <el-table
        v-loading="loading"
        :data="networkList"
        @selection-change="handleSelectionChange"
        style="width: 100%"
        height="calc(100vh - 200px)"
        border
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column label="名称" min-width="200">
          <template #default="{ row }">
            <div class="network-name">
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="驱动" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getDriverTagType(row.driver)">
              {{ row.driver }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="子网" width="150">
          <template #default="{ row }">
            <span>{{ row.ipam?.subnet || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="网关" width="150">
          <template #default="{ row }">
            <span>{{ row.ipam?.gateway || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column label="标签" min-width="200">
          <template #default="{ row }">
            <div v-if="row.labels && Object.keys(row.labels).length > 0" class="labels-container">
              <el-tag 
                v-for="(value, key) in row.labels" 
                :key="key" 
                size="small"
                class="label-tag"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="showNetworkDetail(row)"
            >
              详情
            </el-button>
            <el-button 
              type="danger" 
              link 
              size="small" 
              @click="deleteNetwork(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 创建网络对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建网络"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form 
        :model="createForm" 
        :rules="createRules" 
        ref="createFormRef" 
        label-width="100px"
      >
        <el-form-item label="网络名称" prop="name">
          <el-input
            v-model="createForm.name"
            placeholder="请输入网络名称"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="驱动类型" prop="driver">
          <el-select v-model="createForm.driver" style="width: 100%">
            <el-option label="bridge" value="bridge" />
            <el-option label="host" value="host" />
            <el-option label="overlay" value="overlay" />
            <el-option label="macvlan" value="macvlan" />
            <el-option label="none" value="none" />
          </el-select>
        </el-form-item>

        <el-form-item label="子网">
          <el-input
            v-model="createForm.subnet"
            placeholder="例如: 172.18.0.0/16"
            clearable
          />
        </el-form-item>

        <el-form-item label="网关">
          <el-input
            v-model="createForm.gateway"
            placeholder="例如: 172.18.0.1"
            clearable
          />
        </el-form-item>

        <el-form-item label="网络选项">
          <el-checkbox v-model="createForm.enableIPv6">启用 IPv6</el-checkbox>
          <el-checkbox v-model="createForm.internal">内部网络</el-checkbox>
          <el-checkbox v-model="createForm.attachable">可附加</el-checkbox>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateNetwork">
          创建
        </el-button>
      </template>
    </el-dialog>

    <!-- 网络详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="网络详情"
      width="800px"
    >
      <div v-if="networkDetail" class="network-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="网络ID">
            {{ networkDetail.id }}
          </el-descriptions-item>
          <el-descriptions-item label="网络名称">
            {{ networkDetail.name }}
          </el-descriptions-item>
          <el-descriptions-item label="驱动类型">
            <el-tag :type="networkDetail.driver === 'bridge' ? 'primary' : 'info'">
              {{ networkDetail.driver }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="范围">
            {{ networkDetail.scope }}
          </el-descriptions-item>
          <el-descriptions-item label="IPv6">
            <el-tag :type="networkDetail.enableIPv6 ? 'success' : 'info'">
              {{ networkDetail.enableIPv6 ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="内部网络">
            <el-tag :type="networkDetail.internal ? 'warning' : 'info'">
              {{ networkDetail.internal ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="可附加">
            <el-tag :type="networkDetail.attachable ? 'success' : 'info'">
              {{ networkDetail.attachable ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(networkDetail.created) }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="networkDetail.ipam" class="ipam-section">
          <h4>IPAM 配置</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="子网">
              {{ networkDetail.ipam.subnet || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="网关">
              {{ networkDetail.ipam.gateway || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <div v-if="networkDetail.containers && Object.keys(networkDetail.containers).length > 0" class="containers-section">
          <h4>连接的容器</h4>
          <el-table :data="Object.entries(networkDetail.containers)" border>
            <el-table-column label="容器ID" width="120">
              <template #default="{ row }">
                <span class="container-id">{{ row[0].substring(0, 12) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="容器名称" prop="1.name" />
            <el-table-column label="IPv4地址" prop="1.ipv4Address" />
            <el-table-column label="IPv6地址" prop="1.ipv6Address" />
            <el-table-column label="MAC地址" prop="1.macAddress" />
          </el-table>
        </div>

        <div v-if="networkDetail.labels && Object.keys(networkDetail.labels).length > 0" class="labels-section">
          <h4>标签</h4>
          <div class="labels-container">
            <el-tag 
              v-for="(value, key) in networkDetail.labels" 
              :key="key" 
              class="label-tag"
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
.docker-network-management {
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

.network-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #666;
  cursor: pointer;
}

.network-id:hover {
  color: #409eff;
}

.network-name {
  font-weight: 500;
  color: #333;
}

.container-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #666;
}

.pagination-container {
  padding: 20px;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #f0f0f0;
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

/* 网络详情样式 */
.network-detail {
  .ipam-section,
  .containers-section,
  .labels-section {
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
  
  .labels-container {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    
    .label-tag {
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
  .docker-network-management {
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

/* 描述列表样式 */
.el-descriptions {
  margin-bottom: 16px;
}

.el-descriptions__label {
  font-weight: 600;
  color: #333;
}

/* 复选框组样式 */
.el-form-item .el-checkbox {
  margin-right: 16px;
}

.el-form-item .el-checkbox:last-child {
  margin-right: 0;
}
</style>