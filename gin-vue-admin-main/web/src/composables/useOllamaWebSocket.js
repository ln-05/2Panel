import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'

/**
 * Ollama WebSocket连接管理
 */
export function useOllamaWebSocket() {
  const socket = ref(null)
  const isConnected = ref(false)
  const connectionStatus = ref('disconnected') // disconnected, connecting, connected, error
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectTimer = ref(null)

  // 获取WebSocket URL
  const getWebSocketUrl = () => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${protocol}//${host}/ws/ollama`
  }

  // 连接WebSocket
  const connect = () => {
    if (socket.value && (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING)) {
      return
    }

    connectionStatus.value = 'connecting'
    
    try {
      socket.value = new WebSocket(getWebSocketUrl())

      socket.value.onopen = () => {
        console.log('Ollama WebSocket连接已建立')
        isConnected.value = true
        connectionStatus.value = 'connected'
        reconnectAttempts.value = 0
        
        // 连接成功提示
        ElMessage.success('实时连接已建立')
      }

      socket.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          handleMessage(data)
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }

      socket.value.onclose = (event) => {
        console.log('Ollama WebSocket连接已关闭', event.code, event.reason)
        isConnected.value = false
        connectionStatus.value = 'disconnected'
        
        // 自动重连
        if (event.code !== 1000 && reconnectAttempts.value < maxReconnectAttempts) {
          scheduleReconnect()
        } else if (reconnectAttempts.value >= maxReconnectAttempts) {
          connectionStatus.value = 'error'
          ElMessage.error('WebSocket连接失败，请刷新页面重试')
        }
      }

      socket.value.onerror = (error) => {
        console.error('Ollama WebSocket错误:', error)
        connectionStatus.value = 'error'
      }

    } catch (error) {
      console.error('创建WebSocket连接失败:', error)
      connectionStatus.value = 'error'
    }
  }

  // 安排重连
  const scheduleReconnect = () => {
    if (reconnectTimer.value) return
    
    reconnectAttempts.value++
    const delay = Math.min(3000 * Math.pow(2, reconnectAttempts.value - 1), 30000)
    
    console.log(`尝试重连 (${reconnectAttempts.value}/${maxReconnectAttempts}) 延迟: ${delay}ms`)
    
    reconnectTimer.value = setTimeout(() => {
      reconnectTimer.value = null
      connect()
    }, delay)
  }

  // 断开连接
  const disconnect = () => {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value)
      reconnectTimer.value = null
    }

    if (socket.value) {
      socket.value.close(1000, 'Manual disconnect')
      socket.value = null
    }

    isConnected.value = false
    connectionStatus.value = 'disconnected'
  }

  // 发送消息
  const send = (data) => {
    if (!socket.value || socket.value.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket未连接，无法发送消息')
      return false
    }

    try {
      const message = typeof data === 'string' ? data : JSON.stringify(data)
      socket.value.send(message)
      return true
    } catch (error) {
      console.error('发送WebSocket消息失败:', error)
      return false
    }
  }

  // 消息处理器
  const messageHandlers = new Map()

  // 注册消息处理器
  const onMessage = (type, handler) => {
    if (!messageHandlers.has(type)) {
      messageHandlers.set(type, [])
    }
    messageHandlers.get(type).push(handler)
  }

  // 移除消息处理器
  const offMessage = (type, handler) => {
    if (messageHandlers.has(type)) {
      const handlers = messageHandlers.get(type)
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  // 处理接收到的消息
  const handleMessage = (data) => {
    const { type } = data
    
    if (messageHandlers.has(type)) {
      messageHandlers.get(type).forEach(handler => {
        try {
          handler(data)
        } catch (error) {
          console.error(`处理消息类型 ${type} 时出错:`, error)
        }
      })
    }

    // 处理通用消息类型
    switch (type) {
      case 'model_status_update':
        console.log('模型状态更新:', data)
        break
      case 'download_progress':
        console.log('下载进度:', data)
        break
      case 'operation_result':
        if (data.success) {
          ElMessage.success(data.message || '操作成功')
        } else {
          ElMessage.error(data.message || '操作失败')
        }
        break
      case 'system_status':
        console.log('系统状态更新:', data)
        break
      default:
        console.log('未知消息类型:', type, data)
    }
  }

  // 生命周期管理
  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected,
    connectionStatus,
    connect,
    disconnect,
    send,
    onMessage,
    offMessage
  }
}