
<template>
  <div class="docker-image-management">
    <!-- 顶部操作按钮 -->
    <div class="operation-bar">
      <div class="left-buttons">
        <el-button type="primary" @click="showPullDialog = true">
          <el-icon><Download /></el-icon>
          拉取镜像
        </el-button>
        <el-button @click="showBuildDialog = true">
          <el-icon><Tools /></el-icon>
          构建镜像
        </el-button>
        <el-button @click="refreshImageList">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button @click="checkDockerStatus" :loading="checkingStatus">
          <el-icon><InfoFilled /></el-icon>
          检查Docker状态
        </el-button>
        <el-button 
          v-if="selectedImages.length > 0" 
          type="danger" 
          @click="batchDeleteImages"
        >
          <el-icon><Delete /></el-icon>
          批量删除
        </el-button>
        <el-button @click="pruneImages">
          <el-icon><DeleteFilled /></el-icon>
          清理悬空镜像
        </el-button>
      </div>
      <div class="right-search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索镜像名称或标签"
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

    <!-- 镜像列表表格 -->
    <div class="table-container">
      <el-table
        v-loading="loading"
        :data="imageList"
        @selection-change="handleSelectionChange"
        style="width: 100%"
        height="calc(100vh - 200px)"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column label="ID" width="120">
          <template #default="{ row }">
            <el-tooltip :content="row.id" placement="top">
              <span class="image-id">{{ row.id.substring(0, 12) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column label="名称" min-width="200">
          <template #default="{ row }">
            <div class="image-name-cell">
              <div v-if="row.repoTags && row.repoTags.length > 0">
                <div v-for="tag in row.repoTags" :key="tag" class="tag-item">
                  <el-tag size="small" type="info">{{ tag }}</el-tag>
                </div>
              </div>
              <div v-else class="no-tag">
                <el-tag size="small" type="warning">&lt;none&gt;</el-tag>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="标签" width="120">
          <template #default="{ row }">
            <span v-if="row.repoTags && row.repoTags.length > 0">
              {{ getImageTag(row.repoTags[0]) }}
            </span>
            <span v-else class="text-gray-400">&lt;none&gt;</span>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="showTagDialog(row)"
            >
              标签
            </el-button>
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="exportImage(row)"
            >
              推送
            </el-button>
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="exportImage(row)"
            >
              导出
            </el-button>
            <el-button 
              type="danger" 
              link 
              size="small" 
              @click="deleteImage(row)"
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

    <!-- 拉取镜像对话框 -->
    <el-dialog
      v-model="showPullDialog"
      title="拉取镜像"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="pullForm" :rules="pullRules" ref="pullFormRef" label-width="80px">
        <el-form-item label="镜像名称" prop="image">
          <el-input
            v-model="pullForm.image"
            placeholder="例如: nginx:latest"
            clearable
          />
        </el-form-item>
        <el-form-item label="标签" prop="tag">
          <el-input
            v-model="pullForm.tag"
            placeholder="默认为 latest"
            clearable
          />
        </el-form-item>
      </el-form>
      
      <div v-if="pullLoading" class="pull-progress">
        <el-progress :percentage="pullProgress" />
        <div class="pull-log">
          <pre>{{ pullLog }}</pre>
        </div>
      </div>

      <template #footer>
        <el-button @click="showPullDialog = false" :disabled="pullLoading">取消</el-button>
        <el-button type="primary" @click="handlePullImage" :loading="pullLoading">
          {{ pullLoading ? '拉取中...' : '拉取' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 构建镜像对话框 -->
    <el-dialog
      v-model="showBuildDialog"
      title="构建镜像"
      width="800px"
      :close-on-click-modal="false"
    >
      <el-form :model="buildForm" :rules="buildRules" ref="buildFormRef" label-width="100px">
        <el-form-item label="镜像名称" prop="imageName">
          <el-input
            v-model="buildForm.imageName"
            placeholder="例如: my-app"
            clearable
          />
        </el-form-item>
        <el-form-item label="标签" prop="tag">
          <el-input
            v-model="buildForm.tag"
            placeholder="默认为 latest"
            clearable
          />
        </el-form-item>
        <el-form-item label="Dockerfile" prop="dockerfile">
          <el-input
            v-model="buildForm.dockerfile"
            type="textarea"
            :rows="10"
            placeholder="请输入 Dockerfile 内容"
          />
        </el-form-item>
      </el-form>

      <div v-if="buildLoading" class="build-progress">
        <el-progress :percentage="buildProgress" />
        <div class="build-log">
          <pre>{{ buildLog }}</pre>
        </div>
      </div>

      <template #footer>
        <el-button @click="showBuildDialog = false" :disabled="buildLoading">取消</el-button>
        <el-button type="primary" @click="handleBuildImage" :loading="buildLoading">
          {{ buildLoading ? '构建中...' : '构建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 标签镜像对话框 -->
    <el-dialog
      v-model="showTagImageDialog"
      title="镜像标签"
      width="500px"
    >
      <el-form :model="tagForm" :rules="tagRules" ref="tagFormRef" label-width="100px">
        <el-form-item label="源镜像">
          <el-input v-model="tagForm.sourceImage" disabled />
        </el-form-item>
        <el-form-item label="目标镜像" prop="targetImage">
          <el-input
            v-model="tagForm.targetImage"
            placeholder="例如: my-repo/my-app:v1.0"
            clearable
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showTagImageDialog = false">取消</el-button>
        <el-button type="primary" @click="handleTagImage" :loading="tagLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Download, 
  Tools, 
  Refresh, 
  Delete, 
  DeleteFilled, 
  Search,
  InfoFilled
} from '@element-plus/icons-vue'
import {
  getDockerImageList,
  pullDockerImage,
  deleteDockerImage,
  buildDockerImage,
  tagDockerImage,
  pruneDockerImages,
  exportDockerImage
} from '@/api/dockerImage'

// 需要添加Docker状态检查API
const checkDockerStatus = async () => {
  checkingStatus.value = true
  try {
    // 这里应该调用Docker状态检查API
    const response = await fetch('/api/docker/status')
    const result = await response.json()
    
    if (result.code === 0) {
      ElMessage.success('Docker服务运行正常')
    } else {
      ElMessage.error(`Docker服务异常: ${result.msg}`)
    }
  } catch (error) {
    console.error('检查Docker状态失败:', error)
    ElMessage.error('无法检查Docker状态，请检查网络连接')
  } finally {
    checkingStatus.value = false
  }
}

defineOptions({
  name: 'DockerImageManagement'
})

// 页面状态
const loading = ref(false)
const imageList = ref([])
const selectedImages = ref([])
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const checkingStatus = ref(false)

// 对话框状态
const showPullDialog = ref(false)
const showBuildDialog = ref(false)
const showTagImageDialog = ref(false)

// 拉取镜像相关
const pullLoading = ref(false)
const pullProgress = ref(0)
const pullLog = ref('')
const pullForm = reactive({
  image: '',
  tag: ''
})
const pullRules = {
  image: [
    { required: true, message: '请输入镜像名称', trigger: 'blur' }
  ]
}

// 构建镜像相关
const buildLoading = ref(false)
const buildProgress = ref(0)
const buildLog = ref('')
const buildForm = reactive({
  imageName: '',
  tag: '',
  dockerfile: ''
})
const buildRules = {
  imageName: [
    { required: true, message: '请输入镜像名称', trigger: 'blur' }
  ],
  dockerfile: [
    { required: true, message: '请输入Dockerfile内容', trigger: 'blur' }
  ]
}

// 标签镜像相关
const tagLoading = ref(false)
const tagForm = reactive({
  sourceImage: '',
  targetImage: ''
})
const tagRules = {
  targetImage: [
    { required: true, message: '请输入目标镜像名称', trigger: 'blur' }
  ]
}

// 表单引用
const pullFormRef = ref()
const buildFormRef = ref()
const tagFormRef = ref()

// 获取镜像列表
const fetchImageList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    
    if (searchKeyword.value) {
      params.name = searchKeyword.value
    }

    const response = await getDockerImageList(params)
    
    if (response.code === 0) {
      imageList.value = response.data.list || []
      total.value = response.data.total || 0
    } else {
      // 根据错误码提供更具体的错误信息
      let errorMessage = response.msg || '获取镜像列表失败'
      
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
    console.error('获取镜像列表失败:', error)
    
    // 处理网络错误
    let errorMessage = '获取镜像列表失败'
    
    if (error.response) {
      // 服务器响应了错误状态码
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
      // 请求已发出但没有收到响应
      errorMessage = '网络连接失败，请检查网络连接'
    }
    
    ElMessage.error(errorMessage)
  } finally {
    loading.value = false
  }
}

// 刷新镜像列表
const refreshImageList = () => {
  fetchImageList()
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchImageList()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchImageList()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchImageList()
}

// 选择处理
const handleSelectionChange = (selection) => {
  selectedImages.value = selection
}

// 拉取镜像
const handlePullImage = async () => {
  if (!pullFormRef.value) return
  
  const valid = await pullFormRef.value.validate().catch(() => false)
  if (!valid) return

  pullLoading.value = true
  pullProgress.value = 0
  pullLog.value = ''

  try {
    const response = await pullDockerImage({
      image: pullForm.image,
      tag: pullForm.tag
    })

    if (response.code === 0) {
      pullLog.value = response.data
      pullProgress.value = 100
      ElMessage.success('镜像拉取成功')
      
      // 延迟关闭对话框并刷新列表
      setTimeout(() => {
        showPullDialog.value = false
        resetPullForm()
        fetchImageList()
      }, 1000)
    } else {
      // 根据错误信息提供更具体的提示
      let errorMessage = response.msg || '拉取镜像失败'
      
      if (response.code === 7) {
        errorMessage = '请先登录系统'
      } else if (response.msg && response.msg.includes('Docker client is not available')) {
        errorMessage = 'Docker服务不可用，请检查Docker配置'
      } else if (response.msg && response.msg.includes('not found')) {
        errorMessage = '镜像不存在，请检查镜像名称和标签'
      } else if (response.msg && response.msg.includes('timeout')) {
        errorMessage = '拉取超时，请检查网络连接或稍后重试'
      } else if (response.msg && response.msg.includes('unauthorized')) {
        errorMessage = '无权限拉取该镜像，请检查认证信息'
      }
      
      ElMessage.error(errorMessage)
    }
  } catch (error) {
    console.error('拉取镜像失败:', error)
    
    let errorMessage = '拉取镜像失败'
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
  } finally {
    pullLoading.value = false
  }
}

// 构建镜像
const handleBuildImage = async () => {
  if (!buildFormRef.value) return
  
  const valid = await buildFormRef.value.validate().catch(() => false)
  if (!valid) return

  buildLoading.value = true
  buildProgress.value = 0
  buildLog.value = ''

  try {
    const response = await buildDockerImage({
      imageName: buildForm.imageName,
      tag: buildForm.tag,
      dockerfile: buildForm.dockerfile
    })

    if (response.code === 0) {
      buildLog.value = response.data
      buildProgress.value = 100
      ElMessage.success('镜像构建成功')
      
      // 延迟关闭对话框并刷新列表
      setTimeout(() => {
        showBuildDialog.value = false
        resetBuildForm()
        fetchImageList()
      }, 1000)
    } else {
      ElMessage.error(response.msg || '构建镜像失败')
    }
  } catch (error) {
    console.error('构建镜像失败:', error)
    ElMessage.error('构建镜像失败')
  } finally {
    buildLoading.value = false
  }
}

// 显示标签对话框
const showTagDialog = (row) => {
  tagForm.sourceImage = row.repoTags?.[0] || row.id
  tagForm.targetImage = ''
  showTagImageDialog.value = true
}

// 标签镜像
const handleTagImage = async () => {
  if (!tagFormRef.value) return
  
  const valid = await tagFormRef.value.validate().catch(() => false)
  if (!valid) return

  tagLoading.value = true

  try {
    const response = await tagDockerImage({
      sourceImage: tagForm.sourceImage,
      targetImage: tagForm.targetImage
    })

    if (response.code === 0) {
      ElMessage.success('镜像标签设置成功')
      showTagImageDialog.value = false
      resetTagForm()
      fetchImageList()
    } else {
      ElMessage.error(response.msg || '设置镜像标签失败')
    }
  } catch (error) {
    console.error('设置镜像标签失败:', error)
    ElMessage.error('设置镜像标签失败')
  } finally {
    tagLoading.value = false
  }
}

// 删除镜像
const deleteImage = (row) => {
  ElMessageBox.confirm(
    `确定要删除镜像 ${row.repoTags?.[0] || row.id.substring(0, 12)} 吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteDockerImage(row.id)
      
      if (response.code === 0) {
        ElMessage.success('镜像删除成功')
        fetchImageList()
      } else {
        ElMessage.error(response.msg || '删除镜像失败')
      }
    } catch (error) {
      console.error('删除镜像失败:', error)
      ElMessage.error('删除镜像失败')
    }
  })
}

// 批量删除镜像
const batchDeleteImages = () => {
  if (selectedImages.value.length === 0) {
    ElMessage.warning('请选择要删除的镜像')
    return
  }

  ElMessageBox.confirm(
    `确定要删除选中的 ${selectedImages.value.length} 个镜像吗？`,
    '确认批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const deletePromises = selectedImages.value.map(image => 
        deleteDockerImage(image.id)
      )
      
      await Promise.all(deletePromises)
      ElMessage.success('批量删除成功')
      selectedImages.value = []
      fetchImageList()
    } catch (error) {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  })
}

// 清理悬空镜像
const pruneImages = () => {
  ElMessageBox.confirm(
    '确定要清理所有悬空镜像吗？此操作不可恢复。',
    '确认清理',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await pruneDockerImages(true)
      
      if (response.code === 0) {
        const { deletedCount, spaceReclaimed } = response.data
        ElMessage.success(`清理完成，删除了 ${deletedCount} 个镜像，释放了 ${formatSize(spaceReclaimed)} 空间`)
        fetchImageList()
      } else {
        ElMessage.error(response.msg || '清理镜像失败')
      }
    } catch (error) {
      console.error('清理镜像失败:', error)
      ElMessage.error('清理镜像失败')
    }
  })
}

// 导出镜像
const exportImage = (row) => {
  ElMessageBox.confirm(
    `确定要导出镜像 ${row.repoTags?.[0] || row.id.substring(0, 12)} 吗？`,
    '确认导出',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(async () => {
    try {
      const response = await exportDockerImage({
        images: [row.id]
      })
      
      // 创建下载链接
      const blob = new Blob([response])
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${row.repoTags?.[0]?.replace(/[:/]/g, '_') || 'image'}.tar`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      ElMessage.success('镜像导出成功')
    } catch (error) {
      console.error('导出镜像失败:', error)
      ElMessage.error('导出镜像失败')
    }
  })
}

// 工具函数
const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

const getImageTag = (repoTag) => {
  if (!repoTag) return 'latest'
  const parts = repoTag.split(':')
  return parts.length > 1 ? parts[parts.length - 1] : 'latest'
}

// 重置表单
const resetPullForm = () => {
  pullForm.image = ''
  pullForm.tag = ''
  pullProgress.value = 0
  pullLog.value = ''
}

const resetBuildForm = () => {
  buildForm.imageName = ''
  buildForm.tag = ''
  buildForm.dockerfile = ''
  buildProgress.value = 0
  buildLog.value = ''
}

const resetTagForm = () => {
  tagForm.sourceImage = ''
  tagForm.targetImage = ''
}

// 页面加载时获取数据
onMounted(() => {
  fetchImageList()
})
</script>

<style scoped>
.docker-image-management {
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

.image-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #666;
  cursor: pointer;
}

.image-id:hover {
  color: #409eff;
}

.image-name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tag-item {
  margin-bottom: 4px;
}

.tag-item:last-child {
  margin-bottom: 0;
}

.no-tag {
  color: #999;
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

/* 进度显示样式 */
.pull-progress,
.build-progress {
  margin-top: 20px;
  padding: 16px;
  background-color: #f8f9fa;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.pull-log,
.build-log {
  margin-top: 12px;
  max-height: 200px;
  overflow-y: auto;
  background-color: #2d3748;
  color: #e2e8f0;
  padding: 12px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.pull-log pre,
.build-log pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
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
  .docker-image-management {
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

.el-textarea__inner {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.4;
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

/* 进度条样式 */
.el-progress {
  margin-bottom: 12px;
}

.el-progress__text {
  font-size: 12px;
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
</style>
