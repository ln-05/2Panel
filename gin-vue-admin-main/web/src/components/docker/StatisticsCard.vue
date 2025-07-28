<template>
  <div 
    class="statistics-card" 
    :class="[`${type}-card`, { 'clickable': clickable, 'loading': loading }]"
    @click="handleClick"
  >
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <el-icon class="loading-icon"><Loading /></el-icon>
    </div>
    
    <!-- 卡片内容 -->
    <div class="card-content">
      <!-- 图标区域 -->
      <div class="stat-icon">
        <el-icon>
          <component :is="iconComponent" />
        </el-icon>
      </div>
      
      <!-- 统计内容 -->
      <div class="stat-content">
        <div class="stat-title">{{ title }}</div>
        <div class="stat-value">{{ displayValue }}</div>
        <div class="stat-details" v-if="details && details.length > 0">
          <div 
            v-for="(detail, index) in details" 
            :key="index"
            class="stat-detail-item"
            :class="detail.type"
          >
            {{ detail.label }}: {{ detail.value }}
          </div>
        </div>
        <div class="stat-subtitle" v-if="subtitle">
          {{ subtitle }}
        </div>
      </div>
    </div>
    
    <!-- 趋势指示器 -->
    <div class="trend-indicator" v-if="trend">
      <el-icon :class="trend.type">
        <ArrowUp v-if="trend.type === 'up'" />
        <ArrowDown v-if="trend.type === 'down'" />
        <Minus v-if="trend.type === 'stable'" />
      </el-icon>
      <span class="trend-text">{{ trend.text }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { 
  Box, 
  Picture, 
  Connection, 
  FolderOpened,
  Monitor,
  Setting,
  Loading,
  ArrowUp,
  ArrowDown,
  Minus
} from '@element-plus/icons-vue'

const props = defineProps({
  // 卡片类型，影响样式
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['container', 'image', 'network', 'volume', 'system', 'config', 'default'].includes(value)
  },
  // 卡片标题
  title: {
    type: String,
    required: true
  },
  // 主要数值
  value: {
    type: [Number, String],
    default: 0
  },
  // 副标题
  subtitle: {
    type: String,
    default: ''
  },
  // 详细信息数组
  details: {
    type: Array,
    default: () => []
  },
  // 图标名称
  icon: {
    type: String,
    default: 'Box'
  },
  // 是否可点击
  clickable: {
    type: Boolean,
    default: true
  },
  // 加载状态
  loading: {
    type: Boolean,
    default: false
  },
  // 趋势信息
  trend: {
    type: Object,
    default: null
  },
  // 格式化函数
  formatter: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['click'])

// 图标组件映射
const iconMap = {
  Box,
  Picture,
  Connection,
  FolderOpened,
  Monitor,
  Setting
}

// 计算图标组件
const iconComponent = computed(() => {
  return iconMap[props.icon] || Box
})

// 计算显示值
const displayValue = computed(() => {
  if (props.formatter && typeof props.formatter === 'function') {
    return props.formatter(props.value)
  }
  
  if (typeof props.value === 'number') {
    // 数字格式化
    if (props.value >= 1000000) {
      return (props.value / 1000000).toFixed(1) + 'M'
    } else if (props.value >= 1000) {
      return (props.value / 1000).toFixed(1) + 'K'
    }
    return props.value.toString()
  }
  
  return props.value || '0'
})

// 点击处理
const handleClick = () => {
  if (props.clickable && !props.loading) {
    emit('click')
  }
}
</script>

<style scoped>
.statistics-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
  min-height: 120px;
}

.statistics-card.clickable {
  cursor: pointer;
}

.statistics-card.clickable:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.statistics-card.loading {
  pointer-events: none;
}

/* 加载状态 */
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.loading-icon {
  font-size: 24px;
  color: #3b82f6;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 卡片内容 */
.card-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  height: 100%;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 80px;
}

.stat-title {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
  font-weight: 500;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
  line-height: 1;
}

.stat-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-bottom: 4px;
}

.stat-detail-item {
  font-size: 12px;
  color: #9ca3af;
}

.stat-subtitle {
  font-size: 12px;
  color: #6b7280;
  margin-top: auto;
}

/* 趋势指示器 */
.trend-indicator {
  position: absolute;
  top: 16px;
  right: 16px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.trend-indicator .up {
  color: #10b981;
}

.trend-indicator .down {
  color: #ef4444;
}

.trend-indicator .stable {
  color: #6b7280;
}

.trend-text {
  color: inherit;
}

/* 不同类型卡片的样式 */
.container-card .stat-icon {
  background: #dbeafe;
  color: #3b82f6;
}

.container-card {
  border-left: 4px solid #3b82f6;
}

.container-card .stat-detail-item.running {
  color: #10b981;
}

.container-card .stat-detail-item.stopped {
  color: #6b7280;
}

.container-card .stat-detail-item.paused {
  color: #f59e0b;
}

.image-card .stat-icon {
  background: #d1fae5;
  color: #10b981;
}

.image-card {
  border-left: 4px solid #10b981;
}

.image-card .stat-detail-item.size {
  color: #3b82f6;
}

.network-card .stat-icon {
  background: #fef3c7;
  color: #f59e0b;
}

.network-card {
  border-left: 4px solid #f59e0b;
}

.network-card .stat-detail-item.bridge {
  color: #f59e0b;
}

.network-card .stat-detail-item.host {
  color: #8b5cf6;
}

.volume-card .stat-icon {
  background: #fce7f3;
  color: #ec4899;
}

.volume-card {
  border-left: 4px solid #ec4899;
}

.volume-card .stat-detail-item.size {
  color: #3b82f6;
}

.system-card .stat-icon {
  background: #e0e7ff;
  color: #6366f1;
}

.system-card {
  border-left: 4px solid #6366f1;
}

.config-card .stat-icon {
  background: #f3e8ff;
  color: #8b5cf6;
}

.config-card {
  border-left: 4px solid #8b5cf6;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .statistics-card {
    min-height: 100px;
    padding: 16px;
  }
  
  .card-content {
    gap: 12px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .trend-indicator {
    position: static;
    margin-top: 8px;
    justify-content: flex-start;
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .statistics-card {
    background: #1f2937;
    border-color: #374151;
  }
  
  .stat-title,
  .stat-subtitle {
    color: #9ca3af;
  }
  
  .stat-value {
    color: #f9fafb;
  }
  
  .stat-detail-item {
    color: #6b7280;
  }
  
  .loading-overlay {
    background: rgba(31, 41, 55, 0.8);
  }
}
</style>