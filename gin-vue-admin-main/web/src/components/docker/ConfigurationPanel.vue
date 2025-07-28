<template>
  <el-card class="configuration-panel" :class="{ 'loading': loading }">
    <template #header>
      <div class="panel-header">
        <div class="header-left">
          <el-icon class="header-icon"><Setting /></el-icon>
          <span class="panel-title">{{ title }}</span>
        </div>
        <div class="header-actions">
          <el-button 
            v-if="showConfigButton"
            type="text" 
            @click="handleConfigClick"
            class="config-button"
          >
            {{ configButtonText }}
          </el-button>
          <el-button 
            v-if="showRefreshButton"
            type="text" 
            @click="handleRefresh"
            :loading="loading"
            class="refresh-button"
          >
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </template>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-content">
      <el-skeleton :rows="3" animated />
    </div>
    
    <!-- 配置内容 -->
    <div v-else class="panel-content">
      <div 
        v-for="(item, index) in configItems" 
        :key="index"
        class="config-item"
        :class="{ 'clickable': item.clickable }"
        @click="handleItemClick(item)"
      >
        <div class="item-left">
          <el-icon v-if="item.icon" class="item-icon">
            <component :is="getIconComponent(item.icon)" />
          </el-icon>
          <div class="item-content">
            <span class="item-label">{{ item.label }}</span>
            <span v-if="item.description" class="item-description">{{ item.description }}</span>
          </div>
        </div>
        <div class="item-right">
          <span class="item-value" :class="item.valueType">
            {{ formatValue(item.value, item.formatter) }}
          </span>
          <el-icon v-if="item.clickable" class="arrow-icon"><ArrowRight /></el-icon>
        </div>
      </div>
      
      <!-- 空状态 -->
      <div v-if="configItems.length === 0" class="empty-state">
        <el-icon class="empty-icon"><DocumentRemove /></el-icon>
        <p class="empty-text">暂无配置信息</p>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { computed } from 'vue'
import { 
  Setting, 
  Refresh, 
  ArrowRight, 
  DocumentRemove,
  Link,
  FolderOpened,
  Monitor,
  Connection
} from '@element-plus/icons-vue'

const props = defineProps({
  // 面板标题
  title: {
    type: String,
    default: '配置信息'
  },
  // 配置项数组
  configItems: {
    type: Array,
    default: () => []
  },
  // 加载状态
  loading: {
    type: Boolean,
    default: false
  },
  // 是否显示配置按钮
  showConfigButton: {
    type: Boolean,
    default: true
  },
  // 配置按钮文本
  configButtonText: {
    type: String,
    default: '去配置'
  },
  // 是否显示刷新按钮
  showRefreshButton: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['config-click', 'refresh', 'item-click'])

// 图标组件映射
const iconMap = {
  Link,
  FolderOpened,
  Monitor,
  Connection,
  Setting
}

// 获取图标组件
const getIconComponent = (iconName) => {
  return iconMap[iconName] || Setting
}

// 格式化值
const formatValue = (value, formatter) => {
  if (formatter && typeof formatter === 'function') {
    return formatter(value)
  }
  
  if (value === null || value === undefined) {
    return '未设置'
  }
  
  if (Array.isArray(value)) {
    if (value.length === 0) {
      return '未配置'
    }
    if (value.length === 1) {
      return value[0]
    }
    return `${value[0]} 等 ${value.length} 项`
  }
  
  if (typeof value === 'boolean') {
    return value ? '启用' : '禁用'
  }
  
  return value.toString()
}

// 处理配置按钮点击
const handleConfigClick = () => {
  emit('config-click')
}

// 处理刷新按钮点击
const handleRefresh = () => {
  emit('refresh')
}

// 处理配置项点击
const handleItemClick = (item) => {
  if (item.clickable) {
    emit('item-click', item)
  }
}
</script>

<style scoped>
.configuration-panel {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.configuration-panel.loading {
  pointer-events: none;
}

.configuration-panel :deep(.el-card__header) {
  background: #fafbfc;
  border-bottom: 1px solid #e5e7eb;
  padding: 16px 20px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-icon {
  font-size: 18px;
  color: #3b82f6;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.config-button {
  color: #3b82f6;
  font-weight: 500;
  padding: 4px 8px;
}

.config-button:hover {
  background: #eff6ff;
}

.refresh-button {
  color: #6b7280;
  padding: 4px;
}

.refresh-button:hover {
  color: #3b82f6;
  background: #f3f4f6;
}

/* 面板内容 */
.panel-content {
  padding: 0;
}

.loading-content {
  padding: 20px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f3f4f6;
  transition: all 0.2s ease;
}

.config-item:last-child {
  border-bottom: none;
}

.config-item.clickable {
  cursor: pointer;
}

.config-item.clickable:hover {
  background: #f9fafb;
}

.item-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.item-icon {
  font-size: 16px;
  color: #6b7280;
  flex-shrink: 0;
}

.item-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.item-label {
  font-size: 14px;
  color: #374151;
  font-weight: 500;
}

.item-description {
  font-size: 12px;
  color: #9ca3af;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-value {
  font-size: 14px;
  color: #1f2937;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-value.success {
  background: #d1fae5;
  color: #065f46;
}

.item-value.warning {
  background: #fef3c7;
  color: #92400e;
}

.item-value.error {
  background: #fee2e2;
  color: #991b1b;
}

.item-value.info {
  background: #dbeafe;
  color: #1e40af;
}

.arrow-icon {
  font-size: 14px;
  color: #9ca3af;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 12px;
  color: #d1d5db;
}

.empty-text {
  font-size: 14px;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .configuration-panel :deep(.el-card__header) {
    padding: 12px 16px;
  }
  
  .panel-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .config-item {
    padding: 12px 16px;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .item-left {
    width: 100%;
  }
  
  .item-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .item-value {
    max-width: 100%;
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .configuration-panel {
    background: #1f2937;
    border-color: #374151;
  }
  
  .configuration-panel :deep(.el-card__header) {
    background: #111827;
    border-color: #374151;
  }
  
  .panel-title,
  .item-label {
    color: #f9fafb;
  }
  
  .item-description {
    color: #6b7280;
  }
  
  .item-value {
    background: #374151;
    color: #f9fafb;
  }
  
  .item-value.success {
    background: #064e3b;
    color: #6ee7b7;
  }
  
  .item-value.warning {
    background: #78350f;
    color: #fcd34d;
  }
  
  .item-value.error {
    background: #7f1d1d;
    color: #fca5a5;
  }
  
  .item-value.info {
    background: #1e3a8a;
    color: #93c5fd;
  }
  
  .config-item.clickable:hover {
    background: #374151;
  }
  
  .config-button:hover {
    background: #1e3a8a;
  }
  
  .refresh-button:hover {
    background: #374151;
  }
}
</style>