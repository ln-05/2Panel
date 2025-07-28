






<template>
  <div class="docker-overview-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">概览</h2>
        <p class="page-description">Docker环境整体状态概览</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          @click="refreshData" 
          :loading="loading.refresh"
          icon="Refresh"
        >
          刷新
        </el-button>
      </div>
    </div>

    <!-- 统计卡片区域 -->
    <div class="statistics-section">
      <el-row :gutter="20">
        <!-- 容器统计 -->
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon><Box /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ overviewData.containers.total }}</div>
                <div class="stat-title">容器</div>
                <div class="stat-details">
                  <span>运行中: {{ overviewData.containers.running }}</span>
                  <span>已停止: {{ overviewData.containers.stopped }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 镜像统计 -->
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon><Picture /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ overviewData.images.total }}</div>
                <div class="stat-title">镜像</div>
                <div class="stat-details">
                  <span>大小: {{ overviewData.images.size }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 网络统计 -->
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon><Connection /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ overviewData.networks.total }}</div>
                <div class="stat-title">网络</div>
                <div class="stat-details">
                  <span>Bridge: {{ overviewData.networks.bridge }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 存储卷统计 -->
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon">
                <el-icon><FolderOpened /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ overviewData.volumes.total }}</div>
                <div class="stat-title">存储卷</div>
                <div class="stat-details">
                  <span>大小: {{ overviewData.volumes.size }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 诊断信息面板 -->
    <div class="diagnostic-section" v-if="showDiagnostic">
      <el-card class="diagnostic-card">
        <template #header>
          <div class="card-header">
            <span class="card-title">连接诊断</span>
            <el-button 
              size="small" 
              @click="runDiagnostic" 
              :loading="loading.diagnostic"
              icon="Tools"
            >
              重新诊断
            </el-button>
          </div>
        </template>
        <div class="diagnostic-content">
          <!-- 整体状态 -->
          <div class="status-overview">
            <div class="status-indicator" :class="diagnosticData.overallStatus">
              <el-icon>
                <CircleCheck v-if="diagnosticData.overallStatus === 'healthy'" />
                <Warning v-else-if="diagnosticData.overallStatus === 'warning'" />
                <CircleClose v-else />
              </el-icon>
              <span>{{ getStatusText(diagnosticData.overallStatus) }}</span>
            </div>
          </div>

          <!-- 详细信息 -->
          <el-tabs v-model="activeTab" class="diagnostic-tabs">
            <el-tab-pane label="连接状态" name="connection">
              <div class="diagnostic-item">
                <div class="item-header">
                  <span class="item-title">客户端状态</span>
                  <el-tag :type="diagnosticData.clientStatus?.isConnected ? 'success' : 'danger'">
                    {{ diagnosticData.clientStatus?.isConnected ? '已连接' : '未连接' }}
                  </el-tag>
                </div>
                <div class="item-content" v-if="diagnosticData.clientStatus?.error">
                  <p class="error-message">{{ diagnosticData.clientStatus.error }}</p>
                </div>
                <div class="item-content" v-if="diagnosticData.clientStatus?.isConnected">
                  <p>延迟: {{ diagnosticData.clientStatus.pingLatency }}ms</p>
                </div>
              </div>
            </el-tab-pane>

            <el-tab-pane label="权限检查" name="permissions">
              <div class="permissions-grid">
                <div class="permission-item">
                  <span class="permission-label">容器列表</span>
                  <el-tag :type="permissionData.canListContainers ? 'success' : 'danger'">
                    {{ permissionData.canListContainers ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
                <div class="permission-item">
                  <span class="permission-label">镜像列表</span>
                  <el-tag :type="permissionData.canListImages ? 'success' : 'danger'">
                    {{ permissionData.canListImages ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
                <div class="permission-item">
                  <span class="permission-label">网络列表</span>
                  <el-tag :type="permissionData.canListNetworks ? 'success' : 'danger'">
                    {{ permissionData.canListNetworks ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
                <div class="permission-item">
                  <span class="permission-label">存储卷列表</span>
                  <el-tag :type="permissionData.canListVolumes ? 'success' : 'danger'">
                    {{ permissionData.canListVolumes ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
                <div class="permission-item">
                  <span class="permission-label">系统信息</span>
                  <el-tag :type="permissionData.canGetInfo ? 'success' : 'danger'">
                    {{ permissionData.canGetInfo ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
                <div class="permission-item">
                  <span class="permission-label">版本信息</span>
                  <el-tag :type="permissionData.canGetVersion ? 'success' : 'danger'">
                    {{ permissionData.canGetVersion ? '有权限' : '无权限' }}
                  </el-tag>
                </div>
              </div>
              
              <!-- 权限错误信息 -->
              <div v-if="permissionData.errors && permissionData.errors.length > 0" class="error-section">
                <h4>权限错误:</h4>
                <ul class="error-list">
                  <li v-for="error in permissionData.errors" :key="error" class="error-item">
                    {{ error }}
                  </li>
                </ul>
              </div>

              <!-- 建议 -->
              <div v-if="permissionData.suggestions && permissionData.suggestions.length > 0" class="suggestions-section">
                <h4>解决建议:</h4>
                <ul class="suggestions-list">
                  <li v-for="suggestion in permissionData.suggestions" :key="suggestion" class="suggestion-item">
                    {{ suggestion }}
                  </li>
                </ul>
              </div>
            </el-tab-pane>

            <el-tab-pane label="网络测试" name="network">
              <div class="diagnostic-item">
                <div class="item-header">
                  <span class="item-title">网络连通性</span>
                  <el-tag :type="diagnosticData.networkTest?.isReachable ? 'success' : 'danger'">
                    {{ diagnosticData.networkTest?.isReachable ? '可达' : '不可达' }}
                  </el-tag>
                </div>
                <div class="item-content">
                  <p v-if="diagnosticData.networkTest?.responseTime">
                    响应时间: {{ diagnosticData.networkTest.responseTime }}ms
                  </p>
                  <p v-if="diagnosticData.networkTest?.error" class="error-message">
                    {{ diagnosticData.networkTest.error }}
                  </p>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-card>
    </div>

    <!-- 配置信息面板 -->
    <div class="config-section">
      <el-card class="config-card">
        <template #header>
          <div class="card-header">
            <span class="card-title">配置信息</span>
          </div>
        </template>
        <div class="config-content">
          <div class="config-item">
            <span class="config-label">Socket 路径</span>
            <span class="config-value">{{ configData.socketPath || 'tcp://14.103.168.20:2376' }}</span>
          </div>
          <div class="config-item">
            <span class="config-label">Docker 版本</span>
            <span class="config-value">{{ configData.version || '未知' }}</span>
          </div>
          <div class="config-item">
            <span class="config-label">存储驱动</span>
            <span class="config-value">{{ configData.storageDriver || '未知' }}</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading.initial" class="loading-overlay">
      <el-loading-directive />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Box, 
  Picture, 
  Connection, 
  FolderOpened, 
  Refresh,
  CircleCheck,
  Warning,
  CircleClose,
  Tools
} from '@element-plus/icons-vue'
import { 
  getDockerOverview, 
  getDockerConfigSummary
} from '@/api/dockerOverview'
import { 
  getDockerDiagnosis,
  getDockerPermissions
} from '@/api/dockerDiagnostic'

// 响应式数据
const loading = reactive({
  initial: true,
  refresh: false,
  diagnostic: false
})

const overviewData = reactive({
  containers: {
    total: 0,
    running: 0,
    stopped: 0
  },
  images: {
    total: 0,
    size: '0 B'
  },
  networks: {
    total: 0,
    bridge: 0
  },
  volumes: {
    total: 0,
    size: '未知'
  }
})

const configData = reactive({
  socketPath: '',
  version: '',
  storageDriver: ''
})

// 诊断相关数据
const showDiagnostic = ref(false) // 初始不显示诊断面板，等数据加载完成后显示
const activeTab = ref('connection')

const diagnosticData = reactive({
  overallStatus: 'error',
  clientStatus: {
    isInitialized: false,
    isConnected: false,
    error: '',
    pingLatency: 0
  },
  networkTest: {
    isReachable: false,
    responseTime: 0,
    error: '',
    dnsResolution: false,
    portAccessible: false
  },
  versionCheck: {
    clientVersion: '',
    serverVersion: '',
    isCompatible: false,
    recommendedVersion: '',
    error: ''
  },
  configValidation: {
    isValid: false,
    issues: [],
    suggestions: []
  }
})

const permissionData = reactive({
  canListContainers: false,
  canListImages: false,
  canListNetworks: false,
  canListVolumes: false,
  canGetInfo: false,
  canGetVersion: false,
  errors: [],
  suggestions: []
})

// 真实数据加载
const loadData = async (showLoading = false) => {
  try {
    if (showLoading) {
      loading.refresh = true
    }

    // 调用真实的Docker API
    const [overviewResponse, configResponse] = await Promise.all([
      getDockerOverview(),
      getDockerConfigSummary()
    ])
    
    // 更新概览数据
    if (overviewResponse.code === 0) {
      Object.assign(overviewData, overviewResponse.data)
    } else {
      throw new Error(overviewResponse.msg || '获取概览数据失败')
    }
    
    // 更新配置数据
    if (configResponse.code === 0) {
      Object.assign(configData, configResponse.data)
    } else {
      console.warn('获取配置数据失败:', configResponse.msg)
      // 配置数据失败不影响主要功能，使用默认值
      Object.assign(configData, {
        socketPath: 'tcp://14.103.168.20:2376',
        version: '未知',
        storageDriver: '未知'
      })
    }

    if (showLoading) {
      ElMessage.success('数据刷新成功')
    }

  } catch (error) {
    console.error('加载数据失败:', error)
    
    // 显示错误信息
    const errorMsg = error.response?.data?.msg || error.message || '加载数据失败'
    if (showLoading) {
      ElMessage.error(errorMsg)
    }
    
    // 设置默认值，避免页面显示异常
    Object.assign(overviewData, {
      containers: { total: 0, running: 0, stopped: 0 },
      images: { total: 0, size: '0 B' },
      networks: { total: 0, bridge: 0 },
      volumes: { total: 0, size: '0 B' }
    })
    
    Object.assign(configData, {
      socketPath: 'tcp://14.103.168.20:2376',
      version: '未知',
      storageDriver: '未知'
    })
    
  } finally {
    loading.initial = false
    loading.refresh = false
  }
}

// 运行诊断
const runDiagnostic = async () => {
  try {
    loading.diagnostic = true
    
    // 获取诊断结果和权限状态
    const [diagnosisResponse, permissionsResponse] = await Promise.all([
      getDockerDiagnosis(),
      getDockerPermissions()
    ])
    
    // 更新诊断数据
    if (diagnosisResponse.code === 0) {
      Object.assign(diagnosticData, diagnosisResponse.data)
    }
    
    // 更新权限数据
    if (permissionsResponse.code === 0) {
      Object.assign(permissionData, permissionsResponse.data)
    }
    
    // 诊断完成后显示诊断面板
    showDiagnostic.value = true
    
    ElMessage.success('诊断完成')
  } catch (error) {
    console.error('诊断失败:', error)
    ElMessage.error('诊断失败: ' + (error.response?.data?.msg || error.message))
    // 即使诊断失败也显示面板，这样用户可以看到错误信息
    showDiagnostic.value = true
  } finally {
    loading.diagnostic = false
  }
}

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'healthy':
      return '连接正常'
    case 'warning':
      return '存在警告'
    case 'error':
      return '连接异常'
    default:
      return '未知状态'
  }
}

// 刷新数据
const refreshData = () => {
  loadData(true)
  runDiagnostic() // 同时运行诊断
}

// 生命周期
onMounted(() => {
  loadData()
  runDiagnostic() // 页面加载时运行诊断
})
</script>

<style scoped>
.docker-overview-page {
  padding: 20px;
  background: #f5f6fa;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 24px;
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
  align-items: center;
  gap: 12px;
}

.statistics-section {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
}

.stat-title {
  font-size: 14px;
  color: #6b7280;
  margin: 4px 0;
}

.stat-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-details span {
  font-size: 12px;
  color: #9ca3af;
}

.diagnostic-section {
  margin-bottom: 24px;
}

.diagnostic-card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.diagnostic-content {
  padding: 16px 0;
}

.status-overview {
  margin-bottom: 20px;
  text-align: center;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 16px;
}

.status-indicator.healthy {
  background: #f0f9ff;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

.status-indicator.warning {
  background: #fffbeb;
  color: #d97706;
  border: 1px solid #fed7aa;
}

.status-indicator.error {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.diagnostic-tabs {
  margin-top: 16px;
}

.diagnostic-item {
  margin-bottom: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.item-title {
  font-weight: 600;
  color: #374151;
}

.item-content {
  color: #6b7280;
  font-size: 14px;
}

.error-message {
  color: #dc2626;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: #fef2f2;
  padding: 8px 12px;
  border-radius: 4px;
  border-left: 4px solid #dc2626;
  margin: 8px 0;
}

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.permission-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
}

.permission-label {
  font-size: 14px;
  color: #374151;
  font-weight: 500;
}

.error-section, .suggestions-section {
  margin-top: 16px;
  padding: 16px;
  border-radius: 8px;
}

.error-section {
  background: #fef2f2;
  border-left: 4px solid #dc2626;
}

.suggestions-section {
  background: #f0f9ff;
  border-left: 4px solid #0369a1;
}

.error-section h4, .suggestions-section h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 600;
}

.error-section h4 {
  color: #dc2626;
}

.suggestions-section h4 {
  color: #0369a1;
}

.error-list, .suggestions-list {
  margin: 0;
  padding-left: 20px;
}

.error-item, .suggestion-item {
  font-size: 13px;
  margin-bottom: 4px;
}

.error-item {
  color: #991b1b;
}

.suggestion-item {
  color: #1e40af;
}

.config-section {
  margin-bottom: 24px;
}

.config-card {
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.config-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}

.config-item:last-child {
  border-bottom: none;
}

.config-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.config-value {
  font-size: 14px;
  color: #1f2937;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: #f8f9fa;
  padding: 4px 8px;
  border-radius: 4px;
}

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
  z-index: 1000;
}

@media (max-width: 768px) {
  .docker-overview-page {
    padding: 12px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
}
</style>