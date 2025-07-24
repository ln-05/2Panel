<template>
  <div style="padding: 20px;">
    <h1>Ollama模型管理测试页面</h1>
    
    <div style="margin: 20px 0;">
      <el-button type="primary" @click="testAPI">测试API连接</el-button>
      <el-button type="success" @click="testWebSocket">测试WebSocket连接</el-button>
      <el-button @click="goToMain">进入主页面</el-button>
    </div>
    
    <div v-if="apiResult" style="margin: 20px 0;">
      <h3>API测试结果：</h3>
      <pre>{{ JSON.stringify(apiResult, null, 2) }}</pre>
    </div>
    
    <div v-if="wsStatus" style="margin: 20px 0;">
      <h3>WebSocket状态：</h3>
      <el-tag :type="wsStatusType">{{ wsStatus }}</el-tag>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { searchOllamaModel } from '@/api/AI/ollamaModel'
import { useOllamaWebSocket } from '@/composables/useOllamaWebSocket'

const router = useRouter()
const apiResult = ref(null)

// WebSocket测试
const { isConnected, connectionStatus } = useOllamaWebSocket()

const wsStatus = computed(() => {
  return connectionStatus.value
})

const wsStatusType = computed(() => {
  switch (connectionStatus.value) {
    case 'connected': return 'success'
    case 'connecting': return 'warning'
    case 'disconnected': return 'info'
    case 'error': return 'danger'
    default: return 'info'
  }
})

const testAPI = async () => {
  try {
    const response = await searchOllamaModel({
      page: 1,
      pageSize: 10
    })
    apiResult.value = response
    ElMessage.success('API连接成功')
  } catch (error) {
    apiResult.value = { error: error.message }
    ElMessage.error('API连接失败: ' + error.message)
  }
}

const testWebSocket = () => {
  if (isConnected.value) {
    ElMessage.success('WebSocket已连接')
  } else {
    ElMessage.warning('WebSocket未连接，状态: ' + connectionStatus.value)
  }
}

const goToMain = () => {
  router.push('/ollama')
}
</script>