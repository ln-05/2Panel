<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="80%"
    top="5vh"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="chat-container">
      <!-- 模型选择器 -->
      <div class="model-selector">
        <el-select
          v-model="selectedModelId"
          placeholder="选择要对话的模型"
          style="width: 300px"
          @change="handleModelChange"
        >
          <el-option
            v-for="model in runningModels"
            :key="model.id"
            :label="model.name"
            :value="model.id"
          >
            <span>{{ model.name }}</span>
            <el-tag size="small" type="success" style="margin-left: 10px">运行中</el-tag>
          </el-option>
        </el-select>
        
        <div class="model-info" v-if="currentModel">
          <el-tag type="info" size="small">{{ currentModel.size }}</el-tag>
          <span class="model-status">{{ currentModel.message }}</span>
        </div>
      </div>

      <!-- 聊天区域 -->
      <div class="chat-area" ref="chatAreaRef">
        <div class="messages-container">
          <div
            v-for="(message, index) in messages"
            :key="index"
            :class="['message', message.role]"
          >
            <div class="message-avatar">
              <el-avatar :size="32">
                {{ message.role === 'user' ? 'U' : 'AI' }}
              </el-avatar>
            </div>
            <div class="message-content">
              <div class="message-text" v-html="formatMessage(message.content)"></div>
              <div class="message-time">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
          
          <!-- 加载指示器 -->
          <div v-if="isLoading" class="message assistant">
            <div class="message-avatar">
              <el-avatar :size="32">AI</el-avatar>
            </div>
            <div class="message-content">
              <div class="typing-indicator">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-area">
        <div class="input-container">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="3"
            placeholder="输入您的问题..."
            :disabled="!selectedModelId || isLoading"
            @keydown.ctrl.enter="handleSend"
            resize="none"
          />
          <div class="input-actions">
            <div class="input-tips">
              <span>Ctrl + Enter 发送</span>
            </div>
            <div class="input-buttons">
              <el-button @click="handleClear" size="small">清空对话</el-button>
              <el-button 
                type="primary" 
                @click="handleSend"
                :disabled="!inputMessage.trim() || !selectedModelId || isLoading"
                :loading="isLoading"
                size="small"
              >
                发送
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { searchOllamaModel, chatWithOllamaModel } from '@/api/AI/ollamaModel'

// Props & Emits
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  model: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue'])

// 响应式数据
const chatAreaRef = ref()
const selectedModelId = ref('')
const currentModel = ref(null)
const runningModels = ref([])
const inputMessage = ref('')
const isLoading = ref(false)
const messages = ref([])
const conversationContext = ref('')

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const dialogTitle = computed(() => {
  if (currentModel.value) {
    return `AI对话 - ${currentModel.value.name}`
  }
  return 'AI对话'
})

// 方法
const loadRunningModels = async () => {
  try {
    const response = await searchOllamaModel({
      status: ['running'],
      page: 1,
      pageSize: 100
    })
    
    if (response.code === 0) {
      runningModels.value = response.data.list
      
      // 如果传入了特定模型且正在运行，自动选择
      if (props.model && props.model.status === 'running') {
        const targetModel = runningModels.value.find(m => m.id === props.model.id)
        if (targetModel) {
          selectedModelId.value = targetModel.id
          currentModel.value = targetModel
        }
      }
      
      // 如果没有选中模型但有运行中的模型，选择第一个
      if (!selectedModelId.value && runningModels.value.length > 0) {
        selectedModelId.value = runningModels.value[0].id
        currentModel.value = runningModels.value[0]
      }
    }
  } catch (error) {
    ElMessage.error('获取运行中的模型失败: ' + error.message)
  }
}

const handleModelChange = (modelId) => {
  currentModel.value = runningModels.value.find(m => m.id === modelId)
  // 切换模型时清空对话历史和上下文
  messages.value = []
  clearContext()
  addSystemMessage(`已切换到模型: ${currentModel.value.name}`)
}

const handleSend = async () => {
  if (!inputMessage.value.trim() || !selectedModelId.value || isLoading.value) {
    return
  }

  const userMessage = inputMessage.value.trim()
  inputMessage.value = ''

  // 添加用户消息
  addMessage('user', userMessage)

  // 发送到AI模型
  await sendToModel(userMessage)
}

const sendToModel = async (message) => {
  isLoading.value = true
  
  try {
    // 调用真实的Ollama对话API
    const response = await chatWithOllamaModel({
      modelId: selectedModelId.value,
      message: message,
      stream: false,
      context: getLastContext()
    })
    
    if (response.code === 0) {
      addMessage('assistant', response.data.response)
      // 保存上下文用于下次对话
      saveContext(response.data.context)
    } else {
      throw new Error(response.msg || '对话失败')
    }
  } catch (error) {
    ElMessage.error('发送消息失败: ' + error.message)
    addMessage('assistant', '抱歉，我现在无法回应您的消息。请检查模型是否正常运行。')
  } finally {
    isLoading.value = false
  }
}



const addMessage = (role, content) => {
  messages.value.push({
    role,
    content,
    timestamp: new Date()
  })
  
  nextTick(() => {
    scrollToBottom()
  })
}

const addSystemMessage = (content) => {
  messages.value.push({
    role: 'system',
    content,
    timestamp: new Date()
  })
  
  nextTick(() => {
    scrollToBottom()
  })
}

const handleClear = () => {
  messages.value = []
  clearContext()
  if (currentModel.value) {
    addSystemMessage(`对话已清空，当前模型: ${currentModel.value.name}`)
  }
}

const scrollToBottom = () => {
  if (chatAreaRef.value) {
    chatAreaRef.value.scrollTop = chatAreaRef.value.scrollHeight
  }
}

const formatMessage = (content) => {
  // 简单的markdown格式化
  return content
    .replace(/\n/g, '<br>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/\*([^*]+)\*/g, '<em>$1</em>')
}

const formatTime = (timestamp) => {
  return timestamp.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleClose = () => {
  visible.value = false
}

// 上下文管理
const getLastContext = () => {
  return conversationContext.value
}

const saveContext = (context) => {
  if (context) {
    conversationContext.value = context
  }
}

const clearContext = () => {
  conversationContext.value = ''
}

// 监听对话框打开
watch(visible, (newVal) => {
  if (newVal) {
    loadRunningModels()
    if (runningModels.value.length === 0) {
      addSystemMessage('当前没有运行中的模型，请先启动一个模型后再进行对话。')
    }
  }
})

// 生命周期
onMounted(() => {
  if (visible.value) {
    loadRunningModels()
  }
})
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 70vh;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  border-bottom: 1px solid #ebeef5;
  background-color: #fafafa;
}

.model-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.model-status {
  font-size: 12px;
  color: #909399;
}

.chat-area {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: #f8f9fa;
}

.messages-container {
  max-width: 800px;
  margin: 0 auto;
}

.message {
  display: flex;
  margin-bottom: 20px;
  align-items: flex-start;
}

.message.user {
  flex-direction: row-reverse;
}

.message.user .message-content {
  background-color: #409eff;
  color: white;
  margin-right: 10px;
}

.message.assistant .message-content {
  background-color: white;
  border: 1px solid #ebeef5;
  margin-left: 10px;
}

.message.system .message-content {
  background-color: #f0f9ff;
  border: 1px solid #b3d8ff;
  color: #409eff;
  margin-left: 10px;
  font-style: italic;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  position: relative;
}

.message-text {
  line-height: 1.6;
  word-wrap: break-word;
}

.message-text :deep(code) {
  background-color: rgba(0, 0, 0, 0.1);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.message-time {
  font-size: 11px;
  opacity: 0.7;
  margin-top: 5px;
}

.typing-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #409eff;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-indicator span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-indicator span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.input-area {
  border-top: 1px solid #ebeef5;
  padding: 20px;
  background-color: white;
}

.input-container {
  max-width: 800px;
  margin: 0 auto;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.input-tips {
  font-size: 12px;
  color: #909399;
}

.input-buttons {
  display: flex;
  gap: 10px;
}
</style>