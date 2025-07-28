<template>
  <el-dialog 
    title="从服务器同步" 
    v-model="visible" 
    width="500px"
    @close="handleClose"
  >
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef"
      label-width="100px"
    >
      <el-form-item label="服务器地址" prop="serverUrl">
        <el-input 
          v-model="form.serverUrl" 
          placeholder="https://remote-server.com"
        />
      </el-form-item>
      
      <el-form-item label="认证方式" prop="authType">
        <el-radio-group v-model="form.authType">
          <el-radio value="none">无需认证</el-radio>
          <el-radio value="apikey">API密钥</el-radio>
          <el-radio value="token">访问令牌</el-radio>
        </el-radio-group>
      </el-form-item>
      
      <el-form-item 
        v-if="form.authType === 'apikey'" 
        label="API密钥" 
        prop="apiKey"
      >
        <el-input 
          type="password" 
          v-model="form.apiKey" 
          placeholder="请输入远程服务器的API密钥"
          show-password
        />
      </el-form-item>
      
      <el-form-item 
        v-if="form.authType === 'token'" 
        label="访问令牌" 
        prop="accessToken"
      >
        <el-input 
          type="password" 
          v-model="form.accessToken" 
          placeholder="请输入访问令牌"
          show-password
        />
      </el-form-item>
      
      <el-form-item label="同步方式" prop="syncType">
        <el-radio-group v-model="form.syncType">
          <el-radio value="full">完全同步</el-radio>
          <el-radio value="incremental">增量同步</el-radio>
        </el-radio-group>
      </el-form-item>
      
      <el-form-item label="冲突处理">
        <el-checkbox v-model="form.overwrite">覆盖已存在的记录</el-checkbox>
      </el-form-item>
    </el-form>
    
    <!-- 同步进度 -->
    <div v-if="syncing" class="sync-progress">
      <el-progress :percentage="progress" :status="progressStatus" />
      <p class="progress-text">{{ progressText }}</p>
    </div>
    
    <!-- 同步结果 -->
    <div v-if="syncResult" class="sync-result">
      <el-alert
        :title="syncResult.success ? '同步成功' : '同步失败'"
        :type="syncResult.success ? 'success' : 'error'"
        :description="syncResult.message"
        show-icon
        :closable="false"
      />
      
      <div v-if="syncResult.success" class="sync-stats">
        <p>总计处理: {{ syncResult.total_synced }} 条记录</p>
        <p>新增: {{ syncResult.created }} 条</p>
        <p>更新: {{ syncResult.updated }} 条</p>
        <p>跳过: {{ syncResult.skipped }} 条</p>
        <div v-if="syncResult.errors && syncResult.errors.length > 0">
          <p>错误信息:</p>
          <ul>
            <li v-for="error in syncResult.errors" :key="error">{{ error }}</li>
          </ul>
        </div>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          {{ syncing ? '取消' : '关闭' }}
        </el-button>
        <el-button 
          type="primary" 
          @click="handleSync" 
          :loading="syncing"
          :disabled="syncing"
        >
          {{ syncing ? '同步中...' : '开始同步' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { syncDatabase } from '@/api/database'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:visible', 'success'])

// 响应式数据
const formRef = ref()
const syncing = ref(false)
const progress = ref(0)
const progressStatus = ref('')
const progressText = ref('')
const syncResult = ref(null)

// 表单数据
const form = reactive({
  serverUrl: '',
  authType: 'none',
  apiKey: '',
  accessToken: '',
  syncType: 'full',
  overwrite: false
})

// 表单验证规则
const rules = {
  serverUrl: [
    { required: true, message: '请输入服务器地址', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL地址', trigger: 'blur' }
  ],
  apiKey: [
    { required: true, message: '请输入API密钥', trigger: 'blur' }
  ],
  syncType: [
    { required: true, message: '请选择同步方式', trigger: 'change' }
  ]
}

// 计算属性
const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

// 方法
const handleSync = async () => {
  try {
    await formRef.value.validate()
    
    syncing.value = true
    progress.value = 0
    progressStatus.value = ''
    progressText.value = '正在连接远程服务器...'
    syncResult.value = null
    
    // 模拟进度更新
    const progressInterval = setInterval(() => {
      if (progress.value < 90) {
        progress.value += 10
        if (progress.value === 30) {
          progressText.value = '正在获取远程数据库列表...'
        } else if (progress.value === 60) {
          progressText.value = '正在同步数据库配置...'
        } else if (progress.value === 90) {
          progressText.value = '正在完成同步操作...'
        }
      }
    }, 500)
    
    // 根据认证方式设置API密钥
    let apiKey = ''
    if (form.authType === 'apikey') {
      apiKey = form.apiKey
    } else if (form.authType === 'token') {
      apiKey = form.accessToken
    }
    
    const response = await syncDatabase({
      server_url: form.serverUrl,
      api_key: apiKey,
      sync_type: form.syncType,
      overwrite: form.overwrite
    })
    
    clearInterval(progressInterval)
    progress.value = 100
    progressStatus.value = response.code === 0 ? 'success' : 'exception'
    progressText.value = response.code === 0 ? '同步完成' : '同步失败'
    
    if (response.code === 0) {
      syncResult.value = {
        success: true,
        message: response.msg,
        ...response.data
      }
      ElMessage.success('数据库同步成功')
      emit('success')
    } else {
      syncResult.value = {
        success: false,
        message: response.msg || '同步失败'
      }
      ElMessage.error(response.msg || '数据库同步失败')
    }
  } catch (error) {
    console.error('同步失败:', error)
    progress.value = 100
    progressStatus.value = 'exception'
    progressText.value = '同步失败'
    syncResult.value = {
      success: false,
      message: '网络错误或服务器异常'
    }
    ElMessage.error('数据库同步失败')
  } finally {
    syncing.value = false
  }
}

const handleClose = () => {
  if (syncing.value) {
    // 这里可以添加取消同步的逻辑
    syncing.value = false
  }
  visible.value = false
  // 重置表单和状态
  formRef.value?.resetFields()
  progress.value = 0
  progressStatus.value = ''
  progressText.value = ''
  syncResult.value = null
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.sync-progress {
  margin: 20px 0;
}

.progress-text {
  text-align: center;
  margin-top: 10px;
  color: #666;
}

.sync-result {
  margin: 20px 0;
}

.sync-stats {
  margin-top: 15px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.sync-stats p {
  margin: 5px 0;
  color: #606266;
}

.sync-stats ul {
  margin: 10px 0;
  padding-left: 20px;
}

.sync-stats li {
  color: #f56c6c;
  margin: 5px 0;
}
</style>