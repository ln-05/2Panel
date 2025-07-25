<template>
  <el-dialog
    title="确认删除"
    v-model="visible"
    width="400px"
    :before-close="handleClose"
  >
    <div class="delete-confirm-content">
      <div class="warning-text">
        <p>确定要删除数据库连接吗？</p>
        <p class="database-name">{{ database?.name }}</p>
        <p class="warning-note">此操作不可恢复，请谨慎操作！</p>
      </div>
    </div>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          取消
        </el-button>
        <el-button 
          type="danger" 
          @click="handleConfirm"
          :loading="deleting"
        >
          确认删除
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { deleteDatabase } from '@/api/database'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  database: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

// 响应式数据
const deleting = ref(false)

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 方法
const handleConfirm = async () => {
  if (!props.database?.id) {
    ElMessage.error('数据库ID不存在')
    return
  }
  
  try {
    deleting.value = true
    const response = await deleteDatabase({ id: props.database.id })
    
    if (response.code === 200) {
      ElMessage.success(response.msg || '删除成功')
      emit('confirm')
      handleClose()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

const handleClose = () => {
  visible.value = false
}
</script>

<style scoped>
.delete-confirm-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px 0;
}

.warning-icon {
  flex-shrink: 0;
}

.warning-text {
  flex: 1;
}

.warning-text p {
  margin: 0 0 8px 0;
  line-height: 1.5;
}

.database-name {
  font-weight: bold;
  color: #409eff;
  font-size: 16px;
}

.warning-note {
  color: #E6A23C;
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>