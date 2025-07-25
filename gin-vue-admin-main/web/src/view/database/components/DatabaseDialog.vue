<template>
  <el-dialog 
    :title="dialogTitle" 
    v-model="visible" 
    width="600px"
    @close="handleClose"
  >
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef"
      label-width="100px"
    >
      <el-form-item label="连接名称" prop="name">
        <el-input 
          v-model="form.name" 
          placeholder="请输入数据库连接名称"
        />
      </el-form-item>
      
      <el-form-item label="数据库类型" prop="type">
        <el-select 
          v-model="form.type" 
          placeholder="请选择数据库类型"
          style="width: 100%"
        >
          <el-option label="MySQL" value="mysql" />
          <el-option label="PostgreSQL" value="postgresql" />
          <el-option label="Redis" value="redis" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="主机地址" prop="host">
        <el-input 
          v-model="form.host" 
          placeholder="请输入主机地址，如：localhost"
        />
      </el-form-item>
      
      <el-form-item label="端口" prop="port">
        <el-input-number 
          v-model="form.port" 
          :min="1" 
          :max="65535"
          style="width: 100%"
        />
      </el-form-item>
      
      <el-form-item label="用户名" prop="username">
        <el-input 
          v-model="form.username" 
          placeholder="请输入数据库用户名"
        />
      </el-form-item>
      
      <el-form-item label="密码" prop="password">
        <el-input 
          type="password" 
          v-model="form.password" 
          placeholder="请输入数据库密码"
          show-password
        />
      </el-form-item>
      
      <el-form-item label="数据库名" prop="database">
        <el-input 
          v-model="form.database" 
          placeholder="请输入数据库名（可选）"
        />
      </el-form-item>
      
      <el-form-item label="字符集" prop="charset">
        <el-input 
          v-model="form.charset" 
          placeholder="请输入字符集（可选）"
        />
      </el-form-item>
      
      <el-form-item label="描述信息" prop="description">
        <el-input 
          type="textarea" 
          v-model="form.description" 
          placeholder="请输入描述信息（可选）"
          :rows="3"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button 
          @click="handleTestConnection" 
          :loading="testing"
          :disabled="!canTest"
        >
          测试连接
        </el-button>
        <el-button @click="handleClose">
          取消
        </el-button>
        <el-button 
          type="primary" 
          @click="handleSave" 
          :loading="saving"
          :disabled="!canSave"
        >
          保存
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createDatabase, updateDatabase, testDatabase } from '@/api/database'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  mode: {
    type: String,
    default: 'create' // create | edit
  },
  database: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:modelValue', 'save'])

// 响应式数据
const formRef = ref()
const testing = ref(false)
const saving = ref(false)
const connectionValid = ref(false)

// 表单数据
const form = reactive({
  id: null,
  name: '',
  type: 'mysql',
  host: 'localhost',
  port: 3306,
  username: '',
  password: '',
  database: '',
  charset: 'utf8mb4',
  description: ''
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入连接名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择数据库类型', trigger: 'change' }
  ],
  host: [
    { required: true, message: '请输入主机地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', min: 1, max: 65535, message: '端口号必须在1-65535之间', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const dialogTitle = computed(() => {
  return props.mode === 'create' ? '添加数据库连接' : '编辑数据库连接'
})

const canTest = computed(() => {
  // 编辑模式下放宽限制，创建模式下需要完整信息
  if (props.mode === 'edit') {
    return form.host && form.port
  }
  return form.host && form.port && form.username && form.password
})

const canSave = computed(() => {
  // 编辑模式下放宽限制，创建模式下需要完整信息
  if (props.mode === 'edit') {
    return form.name && form.host && form.port
  }
  return form.name && form.host && form.port && form.username && form.password
})

// 监听数据库类型变化，设置默认端口
watch(() => form.type, (newType) => {
  const defaultPorts = {
    mysql: 3306,
    postgresql: 5432,
    redis: 6379
  }
  form.port = defaultPorts[newType] || 3306
})

// 监听props.database变化，填充表单
watch(() => props.database, (newDatabase) => {
  console.log('编辑数据库 - 接收到的数据:', newDatabase)
  console.log('当前模式:', props.mode)
  if (newDatabase && Object.keys(newDatabase).length > 0) {
    Object.assign(form, {
      id: newDatabase.id || null,
      name: newDatabase.name || '',
      type: newDatabase.type || 'mysql',
      host: newDatabase.host || 'localhost',
      port: newDatabase.port || 3306,
      username: newDatabase.username || '',
      password: newDatabase.password || '',
      database: newDatabase.database || '',
      charset: newDatabase.charset || 'utf8mb4',
      description: newDatabase.description || ''
    })
    console.log('填充后的表单数据:', form)
    // 编辑模式下默认连接有效
    if (props.mode === 'edit') {
      connectionValid.value = true
    }
  }
}, { immediate: true })

// 方法
const handleTestConnection = async () => {
  try {
    // 验证必填字段
    if (!form.host || !form.port || !form.username || !form.password) {
      ElMessage.error('请填写完整的连接信息')
      return
    }
    
    testing.value = true
    
    // 如果是编辑模式且有ID，使用现有连接测试
    if (props.mode === 'edit' && form.id) {
      const response = await testDatabase({ database_id: form.id })
      if (response.code === 200) {
        ElMessage.success('连接测试成功')
        connectionValid.value = true
      } else {
        ElMessage.error(response.msg || '连接测试失败')
        connectionValid.value = false
      }
    } else {
      // 创建模式下，先创建数据库连接然后测试
      const tempData = {
        name: form.name || '临时测试连接',
        type: form.type,
        host: form.host,
        port: form.port,
        username: form.username,
        password: form.password,
        database: form.database,
        charset: form.charset,
        description: form.description || '临时测试连接'
      }
      
      // 先创建连接
      const createResponse = await createDatabase(tempData)
      if (createResponse.code === 200) {
        // 创建成功后测试连接
        const testResponse = await testDatabase({ database_id: createResponse.data.id })
        if (testResponse.code === 200) {
          ElMessage.success('连接测试成功')
          connectionValid.value = true
          // 保存临时创建的ID，用于后续保存时更新而不是重复创建
          form.id = createResponse.data.id
        } else {
          ElMessage.error(testResponse.msg || '连接测试失败')
          connectionValid.value = false
        }
      } else {
        ElMessage.error(createResponse.msg || '创建临时连接失败')
        connectionValid.value = false
      }
    }
  } catch (error) {
    console.error('测试连接失败:', error)
    ElMessage.error('连接测试失败')
    connectionValid.value = false
  } finally {
    testing.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value.validate()
    saving.value = true
    
    const formData = {
      name: form.name,
      type: form.type,
      host: form.host,
      port: form.port,
      username: form.username,
      password: form.password,
      database: form.database,
      charset: form.charset,
      description: form.description
    }
    
    console.log('保存数据 - 表单数据:', form)
    console.log('保存数据 - 发送数据:', formData)
    
    let response
    // 如果是创建模式但已经有ID（通过测试连接创建了临时记录），则更新
    if (props.mode === 'create' && form.id) {
      response = await updateDatabase({
        id: parseInt(form.id), // 确保ID是数字类型
        ...formData
      })
    } else if (props.mode === 'create') {
      // 纯创建模式，没有测试过连接
      response = await createDatabase(formData)
    } else {
      // 编辑模式
      response = await updateDatabase({
        id: parseInt(form.id), // 确保ID是数字类型
        ...formData
      })
    }
    
    if (response.code === 200) {
      ElMessage.success(response.msg || '保存成功')
      emit('save')
      handleClose()
    } else {
      ElMessage.error(response.msg || '保存失败')
    }
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  visible.value = false
  // 重置表单
  formRef.value?.resetFields()
  connectionValid.value = false
  // 重置表单数据
  Object.assign(form, {
    id: null,
    name: '',
    type: 'mysql',
    host: 'localhost',
    port: 3306,
    username: '',
    password: '',
    database: '',
    charset: 'utf8mb4',
    description: ''
  })
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>