<template>
  <el-dialog
    v-model="visible"
    title="添加模型"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="模型名称" prop="name">
        <el-input
          v-model="form.name"
          placeholder="请输入模型名称，如: llama2:7b"
          clearable
        />
        <div class="form-tip">
          支持的格式：llama2:7b, codellama:13b, mistral:latest 等
        </div>
      </el-form-item>
      
      <el-form-item label="来源" prop="from">
        <el-select v-model="form.from" placeholder="请选择模型来源" style="width: 100%">
          <el-option label="Ollama官方" value="ollama" />
          <el-option label="Hugging Face" value="huggingface" />
          <el-option label="本地文件" value="local" />
        </el-select>
      </el-form-item>
      
      <el-form-item>
        <div class="model-suggestions">
          <div class="suggestions-title">推荐模型：</div>
          <div class="suggestions-list">
            <el-tooltip
              v-for="model in recommendedModels"
              :key="model.name"
              :content="`${model.description} (${model.size})`"
              placement="top"
            >
              <el-tag
                class="suggestion-tag"
                :type="model.category === 'lightweight' ? 'success' : model.category === 'code' ? 'warning' : 'primary'"
                @click="selectModel(model.name)"
              >
                {{ model.name }}
              </el-tag>
            </el-tooltip>
          </div>
        </div>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleSubmit"
          :loading="loading"
        >
          开始下载
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { createOllamaModelAdvanced } from '@/api/AI/ollamaModel'
import { OLLAMA_CONFIG } from '@/config/ollama'

// Props & Emits
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

// 响应式数据
const formRef = ref()
const loading = ref(false)

const form = reactive({
  name: '',
  from: 'ollama'
})

const rules = {
  name: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9][a-zA-Z0-9._-]*:[a-zA-Z0-9][a-zA-Z0-9._-]*$|^[a-zA-Z0-9][a-zA-Z0-9._-]*$/, 
      message: '模型名称格式不正确', 
      trigger: 'blur' 
    }
  ],
  from: [
    { required: true, message: '请选择模型来源', trigger: 'change' }
  ]
}

// 推荐模型列表
const recommendedModels = computed(() => OLLAMA_CONFIG.RECOMMENDED_MODELS)

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 方法
const selectModel = (modelName) => {
  form.name = modelName
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    const response = await createOllamaModelAdvanced({
      name: form.name,
      from: form.from
    })
    
    if (response.code === 0) {
      ElMessage.success('模型下载任务已启动')
      emit('success')
      handleClose()
    }
  } catch (error) {
    if (error.message) {
      ElMessage.error('创建模型失败: ' + error.message)
    }
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  visible.value = false
  formRef.value?.resetFields()
  form.name = ''
  form.from = 'ollama'
}
</script>

<style scoped>
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.model-suggestions {
  width: 100%;
}

.suggestions-title {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.suggestions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.suggestion-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.suggestion-tag:hover {
  background-color: #409eff;
  color: white;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>