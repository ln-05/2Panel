<template>
  <div class="registry-container">
    <!-- 顶部操作栏 -->
    <div class="header-actions">
      <div class="left-buttons">
        <el-button type="primary" @click="handleCreate">
          添加仓库
        </el-button>
        <el-button type="success" @click="handleSync">
          从1Panel同步
        </el-button>
      </div>
      <div class="search-box">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #suffix>
            <el-icon class="search-icon" @click="handleSearch">
              <Search />
            </el-icon>
          </template>
        </el-input>
        <el-button @click="handleRefresh" class="refresh-btn">
          <el-icon><Refresh /></el-icon>
        </el-button>
        <el-button @click="handleBatchDelete" :disabled="!multipleSelection.length">
          批量删除
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      ref="multipleTable"
      :data="tableData"
      style="width: 100%"
      @selection-change="handleSelectionChange"
      v-loading="loading"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="名称" min-width="120" />
      <el-table-column prop="downloadUrl" label="下载地址" min-width="200" show-overflow-tooltip />
      <el-table-column prop="protocol" label="协议" width="80" />
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">
          <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" size="small">
            {{ scope.row.status === 'active' ? '正常' : '异常' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="160">
        <template #default="scope">
          {{ formatTime(scope.row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleDetail(scope.row)">
            详情
          </el-button>
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">
            编辑
          </el-button>
          <el-button link type="primary" size="small" @click="handleTest(scope.row)">
            测试
          </el-button>
          <el-button link type="danger" size="small" @click="handleDelete(scope.row)">
            删除
          </el-button>
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
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="仓库名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入仓库名称" />
        </el-form-item>
        <el-form-item label="下载地址" prop="downloadUrl">
          <el-input v-model="form.downloadUrl" placeholder="请输入下载地址，如：https://registry-1.docker.io" />
        </el-form-item>
        <el-form-item label="协议" prop="protocol">
          <el-select v-model="form.protocol" placeholder="请选择协议">
            <el-option label="https" value="https" />
            <el-option label="http" value="http" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="form.username" placeholder="请输入用户名（可选）" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" placeholder="请输入密码（可选）" show-password />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" placeholder="请输入描述（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="仓库详情"
      width="600px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="仓库名称">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="下载地址">{{ detailData.downloadUrl }}</el-descriptions-item>
        <el-descriptions-item label="协议">{{ detailData.protocol }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailData.status === 'active' ? 'success' : 'danger'" size="small">
            {{ detailData.status === 'active' ? '正常' : '异常' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailData.username || '无' }}</el-descriptions-item>
        <el-descriptions-item label="是否默认">
          <el-tag :type="detailData.isDefault ? 'success' : 'info'" size="small">
            {{ detailData.isDefault ? '是' : '否' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatTime(detailData.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ formatTime(detailData.updatedAt) }}</el-descriptions-item>
        <el-descriptions-item label="最后测试时间">
          {{ detailData.lastTestTime ? formatTime(detailData.lastTestTime) : '未测试' }}
        </el-descriptions-item>
        <el-descriptions-item label="测试结果">{{ detailData.testResult || '无' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ detailData.description || '无' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <span class="dialog-footer">
          <el-button v-if="!detailData.isDefault" type="warning" @click="handleSetDefault(detailData)">
            设为默认
          </el-button>
          <el-button type="primary" @click="handleTest(detailData)">
            测试连接
          </el-button>
          <el-button @click="detailVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import {
  getRegistryList,
  getRegistryDetail,
  createRegistry,
  updateRegistry,
  deleteRegistry,
  testRegistry,
  setDefaultRegistry,
  syncFrom1Panel
} from '@/api/dockerRegistry'

// 响应式数据
const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref([])
const multipleSelection = ref([])
const searchKeyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
const dialogVisible = ref(false)
const detailVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)

// 表单数据
const form = reactive({
  id: null,
  name: '',
  downloadUrl: '',
  protocol: 'https',
  username: '',
  password: '',
  description: ''
})

// 详情数据
const detailData = ref({})

// 表单引用
const formRef = ref()

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入仓库名称', trigger: 'blur' }
  ],
  downloadUrl: [
    { required: true, message: '请输入下载地址', trigger: 'blur' },
    { pattern: /^https?:\/\//, message: '下载地址必须以http://或https://开头', trigger: 'blur' }
  ],
  protocol: [
    { required: true, message: '请选择协议', trigger: 'change' }
  ]
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

// 获取列表数据
const getList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (searchKeyword.value) {
      params.name = searchKeyword.value
    }
    
    console.log('正在请求仓库列表，参数:', params)
    const res = await getRegistryList(params)
    console.log('仓库列表响应:', res)
    
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
      console.log('获取到仓库数据:', tableData.value)
    } else {
      console.error('API返回错误:', res)
      ElMessage.error(`获取仓库列表失败: ${res.msg || '未知错误'}`)
    }
  } catch (error) {
    console.error('获取仓库列表失败:', error)
    let errorMessage = '获取仓库列表失败'
    
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      console.error('HTTP错误响应:', status, data)
      
      if (status === 401) {
        errorMessage = '认证失败，请重新登录'
      } else if (status === 403) {
        errorMessage = '权限不足，无法访问仓库管理功能'
      } else if (status === 404) {
        errorMessage = 'API接口不存在，请检查后端服务'
      } else if (status === 500) {
        errorMessage = '服务器内部错误'
      } else {
        errorMessage = `请求失败 (${status}): ${data?.msg || '未知错误'}`
      }
    } else if (error.request) {
      console.error('网络请求失败:', error.request)
      errorMessage = '网络连接失败，请检查网络连接'
    } else {
      console.error('请求配置错误:', error.message)
      errorMessage = `请求配置错误: ${error.message}`
    }
    
    ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  page.value = 1
  getList()
}

// 刷新
const handleRefresh = () => {
  searchKeyword.value = ''
  page.value = 1
  getList()
}

// 分页相关
const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getList()
}

const handleCurrentChange = (val) => {
  page.value = val
  getList()
}

// 多选相关
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 批量删除
const handleBatchDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请选择要删除的仓库')
    return
  }
  
  try {
    await ElMessageBox.confirm('确定要删除选中的仓库吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    for (const item of multipleSelection.value) {
      await deleteRegistry(item.id)
    }
    
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 创建
const handleCreate = () => {
  dialogTitle.value = '添加仓库'
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  dialogTitle.value = '编辑仓库'
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    name: row.name,
    downloadUrl: row.downloadUrl,
    protocol: row.protocol,
    username: row.username,
    password: '', // 密码不回显
    description: row.description
  })
  dialogVisible.value = true
}

// 详情
const handleDetail = async (row) => {
  try {
    const res = await getRegistryDetail(row.id)
    if (res.code === 0) {
      detailData.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error('获取仓库详情失败:', error)
    ElMessage.error('获取仓库详情失败')
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除仓库 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteRegistry(row.id)
    ElMessage.success('删除成功')
    getList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 测试连接
const handleTest = async (row) => {
  try {
    const res = await testRegistry(row.id)
    if (res.code === 0) {
      const { success, message } = res.data
      if (success) {
        ElMessage.success(`测试成功: ${message}`)
      } else {
        ElMessage.error(`测试失败: ${message}`)
      }
      getList() // 刷新列表以更新状态
    }
  } catch (error) {
    console.error('测试连接失败:', error)
    ElMessage.error('测试连接失败')
  }
}

// 设为默认
const handleSetDefault = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要将 "${row.name}" 设为默认仓库吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await setDefaultRegistry(row.id)
    ElMessage.success('设置成功')
    getList()
    detailVisible.value = false
  } catch (error) {
    if (error !== 'cancel') {
      console.error('设置默认仓库失败:', error)
      ElMessage.error('设置失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitLoading.value = true
    
    const formData = { ...form }
    if (isEdit.value) {
      await updateRegistry(formData)
      ElMessage.success('更新成功')
    } else {
      delete formData.id
      await createRegistry(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    getList()
  } catch (error) {
    if (error !== false) { // 不是表单验证错误
      console.error('提交失败:', error)
      ElMessage.error('操作失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 从1Panel同步数据
const handleSync = async () => {
  try {
    await ElMessageBox.confirm('确定要从1Panel同步仓库数据吗？这将会导入1Panel中配置的仓库信息。', '确认同步', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    loading.value = true
    const res = await syncFrom1Panel()
    
    if (res.code === 0) {
      ElMessage.success('同步成功')
      // 刷新列表
      await getList()
    } else {
      ElMessage.error(`同步失败: ${res.msg || '未知错误'}`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步失败:', error)
      let errorMessage = '同步失败'
      
      if (error.response) {
        const status = error.response.status
        const data = error.response.data
        
        if (status === 404) {
          errorMessage = '未找到1Panel数据库或表结构不匹配'
        } else if (status === 500) {
          errorMessage = '服务器错误，请检查1Panel数据库连接'
        } else {
          errorMessage = `同步失败 (${status}): ${data?.msg || '未知错误'}`
        }
      } else if (error.request) {
        errorMessage = '网络连接失败，请检查网络连接'
      } else {
        errorMessage = `同步失败: ${error.message}`
      }
      
      ElMessage.error(errorMessage)
    }
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(form, {
    id: null,
    name: '',
    downloadUrl: '',
    protocol: 'https',
    username: '',
    password: '',
    description: ''
  })
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style scoped>
.registry-container {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-box .el-input {
  width: 250px;
}

.search-icon {
  cursor: pointer;
}

.refresh-btn {
  padding: 8px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>