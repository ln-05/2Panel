<template>
  <div class="database-container">
    <!-- 顶部标签页 -->
    <div class="database-tabs">
      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="全部" name="all">
          <template #label>
            <span class="tab-label">全部</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="MySQL" name="mysql">
          <template #label>
            <span class="tab-label">MySQL</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="PostgreSQL" name="postgresql">
          <template #label>
            <span class="tab-label">PostgreSQL</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="Redis" name="redis">
          <template #label>
            <span class="tab-label">Redis</span>
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 操作按钮和搜索栏 -->
    <div class="action-bar">
      <div class="action-buttons">
        <el-button type="primary" size="large" @click="showCreateDialog">
          创建数据库连接
        </el-button>
        <el-button type="success" size="large" @click="showSyncDialog">
          从服务器同步
        </el-button>
        <el-button size="large" @click="refreshList">
          刷新列表
        </el-button>
      </div>
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索数据库连接"
          size="large"
          clearable
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- 数据库列表 -->
    <div class="database-content">
      <div class="filter-section">
        <el-select 
          v-model="filterType" 
          placeholder="选择数据库类型" 
          size="large"
          @change="handleFilterChange"
          clearable
        >
          <el-option label="MySQL" value="mysql" />
          <el-option label="PostgreSQL" value="postgresql" />
          <el-option label="Redis" value="redis" />
        </el-select>
      </div>

      <div class="database-list">
        <h2 class="list-title">数据库连接列表</h2>
        
        <el-table 
          :data="filteredTableData" 
          style="width: 100%" 
          v-loading="loading"
          empty-text="暂无数据库连接"
        >
          <el-table-column prop="name" label="连接名称" sortable />
          <el-table-column prop="type" label="数据库类型" width="120">
            <template #default="scope">
              <el-tag :type="getTypeTagType(scope.row.type)">
                {{ getTypeLabel(scope.row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="host" label="主机地址" />
          <el-table-column prop="port" label="端口" width="80" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="database" label="数据库名" />
          <el-table-column prop="description" label="描述信息" show-overflow-tooltip />
          <el-table-column label="状态" width="80">
            <template #default="scope">
              <el-tag 
                :type="scope.row.status === 'active' ? 'success' : 'danger'"
                size="small"
              >
                {{ scope.row.status === 'active' ? '在线' : '离线' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="handleTest(scope.row)"
              >
                测试
              </el-button>
              <el-button 
                type="primary" 
                link 
                size="small"
                @click="handleEdit(scope.row)"
              >
                编辑
              </el-button>
              <el-button 
                type="danger" 
                link 
                size="small"
                @click="handleDelete(scope.row)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container" v-if="total > 0">
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
    </div>

    <!-- 添加/编辑对话框 -->
    <DatabaseDialog
      v-model="dialogVisible"
      :mode="dialogMode"
      :database="currentDatabase"
      @save="handleSave"
    />

    <!-- 删除确认对话框 -->
    <DeleteConfirmDialog
      v-model="deleteDialogVisible"
      :database="currentDatabase"
      @confirm="handleDeleteConfirm"
    />

    <!-- 同步对话框 -->
    <SyncDialog
      v-model:visible="syncDialogVisible"
      @success="handleSyncSuccess"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import DatabaseDialog from './components/DatabaseDialog.vue'
import DeleteConfirmDialog from './components/DeleteConfirmDialog.vue'
import SyncDialog from './components/SyncDialog.vue'
import { getDatabaseList, testDatabase } from '@/api/database'

// 响应式数据
const activeTab = ref('all')
const searchKeyword = ref('')
const filterType = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框状态
const dialogVisible = ref(false)
const dialogMode = ref('create') // create | edit
const deleteDialogVisible = ref(false)
const syncDialogVisible = ref(false)
const currentDatabase = ref({})

// 表格数据
const tableData = ref([])

// 计算属性
const filteredTableData = computed(() => {
  let data = tableData.value

  // 根据标签页过滤（如果不是'all'的话）
  if (activeTab.value && activeTab.value !== 'all') {
    data = data.filter(item => item.type === activeTab.value)
  }

  // 根据类型过滤器过滤
  if (filterType.value) {
    data = data.filter(item => item.type === filterType.value)
  }

  // 根据搜索关键词过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    data = data.filter(item => 
      item.name?.toLowerCase().includes(keyword) ||
      item.description?.toLowerCase().includes(keyword) ||
      item.host?.toLowerCase().includes(keyword)
    )
  }

  return data
})

// 方法
const loadDatabaseList = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchKeyword.value,
      type_filter: filterType.value
    }
    
    const response = await getDatabaseList(params)
    if (response.code === 200) {
      tableData.value = response.data?.list || []
      total.value = response.data?.total || 0
    } else {
      ElMessage.error(response.msg || '获取数据库列表失败')
    }
  } catch (error) {
    console.error('获取数据库列表失败:', error)
    ElMessage.error('获取数据库列表失败')
  } finally {
    loading.value = false
  }
}

const handleTabClick = (tab) => {
  activeTab.value = tab.props.name
  currentPage.value = 1
  loadDatabaseList()
}

const showCreateDialog = () => {
  dialogMode.value = 'create'
  currentDatabase.value = {}
  dialogVisible.value = true
}

const handleEdit = (database) => {
  console.log('编辑数据库:', database)
  console.log('设置模式为 edit')
  dialogMode.value = 'edit'
  currentDatabase.value = { ...database }
  console.log('设置当前数据库:', currentDatabase.value)
  dialogVisible.value = true
}

const handleDelete = (database) => {
  currentDatabase.value = database
  deleteDialogVisible.value = true
}

const handleTest = async (database) => {
  try {
    loading.value = true
    const response = await testDatabase({ database_id: database.id })
    if (response.code === 200) {
      ElMessage.success('连接测试成功')
      // 更新数据库状态
      const index = tableData.value.findIndex(item => item.id === database.id)
      if (index !== -1) {
        tableData.value[index].status = 'active'
      }
    } else {
      ElMessage.error(response.msg || '连接测试失败')
      // 更新数据库状态
      const index = tableData.value.findIndex(item => item.id === database.id)
      if (index !== -1) {
        tableData.value[index].status = 'inactive'
      }
    }
  } catch (error) {
    console.error('连接测试失败:', error)
    ElMessage.error('连接测试失败')
  } finally {
    loading.value = false
  }
}

const handleSave = () => {
  loadDatabaseList()
}

const handleDeleteConfirm = () => {
  loadDatabaseList()
}

const showSyncDialog = () => {
  syncDialogVisible.value = true
}

const handleSyncSuccess = () => {
  loadDatabaseList()
}

const refreshList = () => {
  currentPage.value = 1
  loadDatabaseList()
}

const handleSearch = () => {
  currentPage.value = 1
  // 这里可以添加防抖逻辑
  setTimeout(() => {
    loadDatabaseList()
  }, 300)
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadDatabaseList()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  loadDatabaseList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadDatabaseList()
}

// 工具方法
const getTypeTagType = (type) => {
  const typeMap = {
    mysql: 'success',
    postgresql: 'primary',
    redis: 'warning'
  }
  return typeMap[type] || 'info'
}

const getTypeLabel = (type) => {
  const typeMap = {
    mysql: 'MySQL',
    postgresql: 'PostgreSQL',
    redis: 'Redis'
  }
  return typeMap[type] || type
}

// 监听器
watch([searchKeyword, filterType], () => {
  // 搜索和过滤时重置页码
  currentPage.value = 1
})

// 生命周期
onMounted(() => {
  loadDatabaseList()
})
</script>

<style scoped>
.database-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.database-tabs {
  background: white;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.tab-label {
  font-weight: 500;
}

.database-status-bar {
  background: white;
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.status-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.version-info {
  border: 1px solid #409eff;
  border-radius: 4px;
  padding: 4px 8px;
  background-color: #f0f9ff;
}

.version-label {
  color: #409eff;
  font-size: 14px;
  font-weight: 500;
}

.status-actions {
  display: flex;
  gap: 8px;
}

.action-bar {
  background: white;
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.action-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-box {
  width: 300px;
}

.database-content {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.filter-section {
  margin-bottom: 20px;
}

.list-title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 20px;
  color: #333;
}

.sort-icon {
  margin-left: 4px;
  color: #409eff;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    gap: 16px;
  }
  
  .action-buttons {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .search-box {
    width: 100%;
  }
  
  .database-status-bar {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
}
</style> 