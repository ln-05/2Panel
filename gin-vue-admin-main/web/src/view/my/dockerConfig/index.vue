<template>
  <div class="docker-config-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">Docker 配置</h2>
        <p class="page-description">管理 Docker 守护进程配置和服务状态</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          @click="refreshConfig" 
          :loading="loading.refresh"
        >
          刷新
        </el-button>
        <el-button 
          type="success" 
          @click="backupConfig" 
          :loading="loading.backup"
        >
          创建备份
        </el-button>
      </div>
    </div>

    <!-- 配置表单 -->
    <div class="config-container">
      <el-form :model="configForm" label-width="120px" ref="configFormRef" class="config-form">
        
        <!-- 镜像加速器 -->
        <div class="form-item-wrapper">
          <el-form-item label="镜像加速器">
            <el-input
              v-model="mirrorInput"
              placeholder="https://docker.1panel.live"
              class="mirror-input"
            >
              <template #append>
                <el-button @click="addMirror" :disabled="!mirrorInput.trim()">
                  添加
                </el-button>
              </template>
            </el-input>
            <div class="mirror-tags" v-if="configForm.registryMirrors.length > 0">
              <el-tag
                v-for="(mirror, index) in configForm.registryMirrors"
                :key="index"
                closable
                @close="removeMirror(index)"
                class="mirror-tag"
              >
                {{ mirror }}
              </el-tag>
            </div>
            <div class="form-help">
              <span>配置 Docker 镜像加速器 URL，提升镜像拉取速度。建议使用国内镜像源</span>
              <el-link type="primary" @click="showCommonMirrors = !showCommonMirrors">
                快速添加
              </el-link>
            </div>
            <div v-if="showCommonMirrors" class="common-mirrors">
              <el-button
                v-for="mirror in commonMirrors"
                :key="mirror.value"
                size="small"
                @click="addCommonMirror(mirror.value)"
                :disabled="configForm.registryMirrors.includes(mirror.value)"
              >
                {{ mirror.label }}
              </el-button>
            </div>
          </el-form-item>
        </div>

        <!-- 私有仓库 -->
        <div class="form-item-wrapper">
          <el-form-item label="私有仓库">
            <el-input
              v-model="insecureInput"
              placeholder="私有仓库地址"
              class="mirror-input"
            >
              <template #append>
                <el-button @click="addInsecure" :disabled="!insecureInput.trim()">
                  添加
                </el-button>
              </template>
            </el-input>
            <div class="mirror-tags" v-if="configForm.insecureRegistries.length > 0">
              <el-tag
                v-for="(registry, index) in configForm.insecureRegistries"
                :key="index"
                closable
                @close="removeInsecure(index)"
                class="mirror-tag"
              >
                {{ registry }}
              </el-tag>
            </div>
            <div class="form-help">
              <span>配置不安全的私有仓库地址，允许通过 HTTP 协议访问</span>
            </div>
          </el-form-item>
        </div>

        <!-- IPv6 -->
        <div class="form-item-wrapper">
          <el-form-item label="IPv6">
            <el-switch 
              v-model="configForm.enableIPv6"
              active-text="启用"
              inactive-text="禁用"
            />
            <div class="form-help">
              <span>启用 IPv6 网络支持</span>
            </div>
          </el-form-item>
        </div>

        <!-- 日志切割 -->
        <div class="form-item-wrapper">
          <el-form-item label="日志切割">
            <el-switch 
              v-model="configForm.enableLogRotation"
              active-text="启用"
              inactive-text="禁用"
            />
            <div class="form-help">
              <span>启用 Docker 日志自动切割功能</span>
            </div>
          </el-form-item>
        </div>

        <!-- iptables -->
        <div class="form-item-wrapper">
          <el-form-item label="iptables">
            <el-switch 
              v-model="configForm.enableIptables"
              active-text="启用"
              inactive-text="禁用"
            />
            <div class="form-help">
              <span>Docker 是否操作 iptables 规则，禁用后需要手动配置网络规则</span>
            </div>
          </el-form-item>
        </div>

        <!-- Live restore -->
        <div class="form-item-wrapper">
          <el-form-item label="Live restore">
            <el-switch 
              v-model="configForm.liveRestore"
              active-text="启用"
              inactive-text="禁用"
            />
            <div class="form-help">
              <span>允许在 Docker 守护进程重启时保持容器运行状态</span>
            </div>
          </el-form-item>
        </div>

        <!-- cgroup driver -->
        <div class="form-item-wrapper">
          <el-form-item label="cgroup driver">
            <el-radio-group v-model="configForm.cgroupDriver">
              <el-radio value="cgroupfs">cgroupfs</el-radio>
              <el-radio value="systemd">systemd</el-radio>
            </el-radio-group>
            <div class="form-help">
              <span>Docker 使用的 cgroup 驱动程序 (Docker Daemon)，需要与容器运行时保持一致</span>
            </div>
          </el-form-item>
        </div>

        <!-- Socket 路径 -->
        <div class="form-item-wrapper">
          <el-form-item label="Socket 路径">
            <el-input
              v-model="configForm.socketPath"
              placeholder="unix:///var/run/docker.sock"
            />
            <div class="form-help">
              <span>Docker 守护进程 Socket 文件路径，用于客户端与守护进程通信</span>
            </div>
          </el-form-item>
        </div>

        <!-- 操作按钮 -->
        <div class="form-actions">
          <el-button 
            type="primary" 
            @click="saveConfig" 
            :loading="loading.save"
          >
            保存
          </el-button>
          <el-button @click="resetConfig">
            重置
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getDockerConfig,
  updateDockerConfig,
  validateDockerConfig,
  backupDockerConfig,
  getBackupList,
  restoreDockerConfig,
  deleteBackup,
  getDockerServiceStatus
} from '@/api/dockerConfig'

const loading = reactive({
  refresh: false,
  backup: false,
  validate: false,
  save: false,
  backups: false
})

const configFormRef = ref()
const configForm = reactive({
  registryMirrors: [],
  insecureRegistries: [],
  storageDriver: 'overlay2',
  logDriver: 'json-file',
  cgroupDriver: 'systemd',
  dataRoot: '',
  socketPath: '',
  enableIPv6: false,
  enableIPForward: true,
  enableIptables: true,
  enableLogRotation: false,
  liveRestore: false
})

const serviceStatus = reactive({
  status: 'unknown',
  uptime: '',
  version: '',
  lastRestart: null,
  errorMsg: ''
})

const backupList = ref([])
const mirrorInput = ref('')
const insecureInput = ref('')
const showCommonMirrors = ref(false)

const commonMirrors = [
  { label: '1Panel镜像源', value: 'https://docker.1panel.live' },
  { label: '阿里云镜像源', value: 'https://registry.cn-hangzhou.aliyuncs.com' },
  { label: '腾讯云镜像源', value: 'https://mirror.ccs.tencentyun.com' },
  { label: '网易云镜像源', value: 'https://hub-mirror.c.163.com' },
  { label: 'Docker中国镜像源', value: 'https://registry.docker-cn.com' }
]

const getStatusClass = (status) => {
  return {
    'status-running': status === 'running',
    'status-stopped': status === 'stopped',
    'status-error': status === 'error'
  }
}

const getStatusCardClass = (status) => {
  return {
    'status-running': status === 'running',
    'status-stopped': status === 'stopped',
    'status-error': status === 'error',
    'status-unknown': status === 'unknown'
  }
}

const getStatusText = (status) => {
  const statusMap = {
    running: '运行中',
    stopped: '已停止',
    error: '错误',
    unknown: '未知'
  }
  return statusMap[status] || '未知'
}

const formatTime = (timeStr) => {
  if (!timeStr) return '未知'
  return new Date(timeStr).toLocaleString()
}

const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const loadConfig = async () => {
  try {
    loading.refresh = true
    const result = await getDockerConfig()
    if (result.code === 0) {
      Object.assign(configForm, result.data.config)
      ElMessage.success('配置加载成功')
    }
  } catch (error) {
    ElMessage.error('加载配置失败: ' + error.message)
  } finally {
    loading.refresh = false
  }
}

const loadServiceStatus = async () => {
  try {
    const result = await getDockerServiceStatus()
    if (result.code === 0) {
      Object.assign(serviceStatus, result.data)
    }
  } catch (error) {
    console.error('获取服务状态失败:', error)
  }
}

const loadBackups = async () => {
  try {
    loading.backups = true
    const result = await getBackupList()
    if (result.code === 0) {
      backupList.value = result.data.backups
    }
  } catch (error) {
    ElMessage.error('加载备份列表失败: ' + error.message)
  } finally {
    loading.backups = false
  }
}

const refreshConfig = async () => {
  await Promise.all([loadConfig(), loadServiceStatus()])
}

const validateConfig = async () => {
  try {
    loading.validate = true
    const result = await validateDockerConfig(configForm)
    if (result.code === 0) {
      if (result.data.valid) {
        ElMessage.success('配置验证通过')
      } else {
        ElMessage.error('配置验证失败: ' + result.data.errors.map(e => e.message).join(', '))
      }
    }
  } catch (error) {
    ElMessage.error('验证配置失败: ' + error.message)
  } finally {
    loading.validate = false
  }
}

const saveConfig = async () => {
  try {
    loading.save = true
    const result = await updateDockerConfig(configForm)
    if (result.code === 0) {
      ElMessage.success('配置保存成功')
      await loadServiceStatus()
    }
  } catch (error) {
    ElMessage.error('保存配置失败: ' + error.message)
  } finally {
    loading.save = false
  }
}

const backupConfig = async () => {
  try {
    const { value: description } = await ElMessageBox.prompt('请输入备份描述', '创建备份', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValue: '手动备份'
    })
    
    loading.backup = true
    const result = await backupDockerConfig(description)
    if (result.code === 0) {
      ElMessage.success('备份创建成功')
      await loadBackups()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('创建备份失败: ' + error.message)
    }
  } finally {
    loading.backup = false
  }
}

const restoreBackup = async (backup) => {
  try {
    await ElMessageBox.confirm(
      `确定要恢复备份 "${backup.description}" 吗？这将覆盖当前配置。`,
      '确认恢复',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const result = await restoreDockerConfig({ backupId: backup.id })
    if (result.code === 0) {
      ElMessage.success('备份恢复成功')
      await refreshConfig()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('恢复备份失败: ' + error.message)
    }
  }
}

const deleteBackupItem = async (backup) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除备份 "${backup.description}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const result = await deleteBackup(backup.id)
    if (result.code === 0) {
      ElMessage.success('备份删除成功')
      await loadBackups()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除备份失败: ' + error.message)
    }
  }
}

// 镜像加速器管理
const addMirror = () => {
  const mirror = mirrorInput.value.trim()
  if (mirror && !configForm.registryMirrors.includes(mirror)) {
    configForm.registryMirrors.push(mirror)
    mirrorInput.value = ''
    ElMessage.success('镜像加速器添加成功')
  } else if (configForm.registryMirrors.includes(mirror)) {
    ElMessage.warning('该镜像加速器已存在')
  }
}

const removeMirror = (index) => {
  configForm.registryMirrors.splice(index, 1)
  ElMessage.success('镜像加速器删除成功')
}

const addCommonMirror = (mirrorUrl) => {
  if (!configForm.registryMirrors.includes(mirrorUrl)) {
    configForm.registryMirrors.push(mirrorUrl)
    ElMessage.success('镜像加速器添加成功')
  }
}

// 私有仓库管理
const addInsecure = () => {
  const registry = insecureInput.value.trim()
  if (registry && !configForm.insecureRegistries.includes(registry)) {
    configForm.insecureRegistries.push(registry)
    insecureInput.value = ''
    ElMessage.success('私有仓库添加成功')
  } else if (configForm.insecureRegistries.includes(registry)) {
    ElMessage.warning('该私有仓库已存在')
  }
}

const removeInsecure = (index) => {
  configForm.insecureRegistries.splice(index, 1)
  ElMessage.success('私有仓库删除成功')
}

const resetConfig = () => {
  configFormRef.value?.resetFields()
  mirrorInput.value = ''
  insecureInput.value = ''
  showCommonMirrors.value = false
}

onMounted(() => {
  refreshConfig()
  loadBackups()
})
</script>

<style scoped>
/* 1Panel风格的现代化样式 */
.docker-config-page {
  padding: 20px;
  background: #f5f6fa;
  min-height: 100vh;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.page-description {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

/* 配置容器 */
.config-container {
  background: white;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 表单样式 */
.config-form {
  max-width: none;
}

.form-item-wrapper {
  margin-bottom: 24px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 24px;
}

.form-item-wrapper:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.mirror-input {
  margin-bottom: 12px;
}

.mirror-tags {
  margin-bottom: 8px;
}

.mirror-tag {
  margin-right: 8px;
  margin-bottom: 4px;
}

.form-help {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

.common-mirrors {
  margin-top: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 4px;
}

.common-mirrors .el-button {
  margin-right: 8px;
  margin-bottom: 4px;
}

/* 状态概览卡片 */
.status-overview {
  margin-bottom: 24px;
}

.status-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.status-card.status-running {
  border-left: 4px solid #10b981;
}

.status-card.status-stopped {
  border-left: 4px solid #ef4444;
}

.status-card.status-error {
  border-left: 4px solid #f59e0b;
}

.status-card.status-unknown {
  border-left: 4px solid #6b7280;
}

.status-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
}

.status-card.status-running .status-icon-wrapper {
  background: #d1fae5;
}

.status-card.status-stopped .status-icon-wrapper {
  background: #fee2e2;
}

.status-card.status-error .status-icon-wrapper {
  background: #fef3c7;
}

.status-icon {
  font-size: 24px;
}

.status-card.status-running .status-icon {
  color: #10b981;
}

.status-card.status-stopped .status-icon {
  color: #ef4444;
}

.status-card.status-error .status-icon {
  color: #f59e0b;
}

.status-card.status-unknown .status-icon {
  color: #6b7280;
}

.status-content {
  flex: 1;
}

.status-label {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 4px;
  font-weight: 500;
}

.status-value {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

/* 内容区域 */
.content-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 卡片样式 */
.config-card,
.backup-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  background: #fafbfc;
}

.card-title {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.card-title .el-icon {
  margin-right: 8px;
  font-size: 20px;
  color: #3b82f6;
}

.header-btn {
  border-radius: 6px;
  font-weight: 500;
}

.card-content {
  padding: 24px;
}

/* 表单样式 */
.config-form {
  max-width: none;
}

.form-section {
  margin-bottom: 32px;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e5e7eb;
}

.form-select,
.form-input {
  width: 100%;
}

.form-select :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #d1d5db;
  transition: all 0.2s ease;
}

.form-select :deep(.el-input__wrapper):hover {
  border-color: #3b82f6;
}

.form-input :deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #d1d5db;
  transition: all 0.2s ease;
}

.form-input :deep(.el-input__wrapper):hover {
  border-color: #3b82f6;
}

.checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.config-checkbox {
  font-weight: 500;
}

.config-checkbox :deep(.el-checkbox__label) {
  color: #374151;
  font-weight: 500;
}

/* 操作按钮 */
.form-actions {
  display: flex;
  gap: 12px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
  margin-top: 24px;
}

.form-actions .el-button {
  border-radius: 6px;
  font-weight: 500;
  padding: 8px 16px;
}

.action-button {
  border-radius: 8px;
  font-weight: 500;
  padding: 10px 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
}

.action-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 表格样式 */
.backup-table {
  border-radius: 8px;
  overflow: hidden;
}

.backup-table :deep(.el-table__header) {
  background: #f9fafb;
}

.backup-table :deep(.el-table__header th) {
  background: #f9fafb;
  color: #374151;
  font-weight: 600;
  border-bottom: 1px solid #e5e7eb;
}

.backup-table :deep(.el-table__body tr:hover) {
  background: #f9fafb;
}

.backup-table :deep(.el-table__body td) {
  border-bottom: 1px solid #f3f4f6;
}

.time-text,
.size-text {
  color: #6b7280;
  font-size: 13px;
}

.table-btn {
  border-radius: 6px;
  font-weight: 500;
  margin-right: 8px;
}

.table-btn:last-child {
  margin-right: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .docker-config-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .status-overview .el-col {
    margin-bottom: 16px;
  }
  
  .checkbox-group {
    flex-direction: column;
    gap: 12px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .action-button {
    width: 100%;
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .docker-config-page {
    background: #111827;
  }
  
  .status-card,
  .config-card,
  .backup-card {
    background: #1f2937;
    border-color: #374151;
  }
  
  .card-header {
    background: #111827;
    border-color: #374151;
  }
  
  .page-title,
  .card-title,
  .section-title,
  .status-value {
    color: #f9fafb;
  }
  
  .page-description,
  .status-label {
    color: #9ca3af;
  }
}
</style>