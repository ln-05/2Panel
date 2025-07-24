// Ollama服务配置
export const OLLAMA_CONFIG = {
  // API基础URL
  API_BASE_URL: 'http://localhost:11434',
  
  // WebSocket URL
  WS_BASE_URL: (() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${protocol}//${host}/ws/ollama`
  })(),
  
  // 默认配置
  DEFAULT_SETTINGS: {
    // 对话配置
    chat: {
      stream: false,
      temperature: 0.7,
      max_tokens: 2048,
      timeout: 30000 // 30秒超时
    },
    
    // WebSocket配置
    websocket: {
      reconnectAttempts: 5,
      reconnectDelay: 3000,
      heartbeatInterval: 30000
    },
    
    // 分页配置
    pagination: {
      defaultPageSize: 10,
      pageSizes: [10, 20, 50, 100]
    }
  },
  
  // 推荐模型列表
  RECOMMENDED_MODELS: [
    {
      name: 'llama2:7b',
      description: '通用对话模型，平衡性能和资源消耗',
      size: '~4GB',
      category: 'general'
    },
    {
      name: 'llama2:13b',
      description: '更强大的通用对话模型',
      size: '~7GB',
      category: 'general'
    },
    {
      name: 'codellama:7b',
      description: '专门用于代码生成和编程问题',
      size: '~4GB',
      category: 'code'
    },
    {
      name: 'codellama:13b',
      description: '更强大的代码生成模型',
      size: '~7GB',
      category: 'code'
    },
    {
      name: 'mistral:7b',
      description: '高质量的对话模型',
      size: '~4GB',
      category: 'general'
    },
    {
      name: 'neural-chat:7b',
      description: '优化的聊天模型',
      size: '~4GB',
      category: 'chat'
    },
    {
      name: 'starling-lm:7b',
      description: '基于强化学习优化的对话模型',
      size: '~4GB',
      category: 'chat'
    },
    {
      name: 'vicuna:7b',
      description: '基于LLaMA微调的对话模型',
      size: '~4GB',
      category: 'chat'
    },
    {
      name: 'orca-mini:3b',
      description: '轻量级模型，适合资源受限环境',
      size: '~2GB',
      category: 'lightweight'
    },
    {
      name: 'phi:2.7b',
      description: '微软开发的轻量级模型',
      size: '~1.5GB',
      category: 'lightweight'
    }
  ],
  
  // 模型状态映射
  STATUS_MAP: {
    running: {
      text: '运行中',
      type: 'success',
      color: '#67C23A'
    },
    stopped: {
      text: '已停止',
      type: 'info',
      color: '#909399'
    },
    downloading: {
      text: '下载中',
      type: 'warning',
      color: '#E6A23C'
    },
    error: {
      text: '错误',
      type: 'danger',
      color: '#F56C6C'
    },
    unavailable: {
      text: '不可用',
      type: 'danger',
      color: '#F56C6C'
    }
  },
  
  // 连接状态映射
  CONNECTION_STATUS_MAP: {
    connected: {
      text: '实时连接已建立',
      type: 'success',
      icon: 'Connection'
    },
    connecting: {
      text: '正在连接...',
      type: 'warning',
      icon: 'Loading'
    },
    disconnected: {
      text: '连接已断开',
      type: 'info',
      icon: 'Disconnect'
    },
    error: {
      text: '连接错误',
      type: 'danger',
      icon: 'Warning'
    }
  }
}

// 获取模型分类
export const getModelsByCategory = (category) => {
  return OLLAMA_CONFIG.RECOMMENDED_MODELS.filter(model => model.category === category)
}

// 获取模型信息
export const getModelInfo = (modelName) => {
  return OLLAMA_CONFIG.RECOMMENDED_MODELS.find(model => model.name === modelName)
}

// 格式化模型状态
export const formatModelStatus = (status) => {
  return OLLAMA_CONFIG.STATUS_MAP[status] || {
    text: status,
    type: 'info',
    color: '#909399'
  }
}

// 格式化连接状态
export const formatConnectionStatus = (status) => {
  return OLLAMA_CONFIG.CONNECTION_STATUS_MAP[status] || {
    text: '未知状态',
    type: 'info',
    icon: 'QuestionFilled'
  }
}