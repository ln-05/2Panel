import service from '@/utils/request'

// 获取Docker配置
export const getDockerConfig = () => {
  return service({
    url: '/docker/config',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取Docker配置失败:', error)
    throw error
  })
}

// 更新Docker配置
export const updateDockerConfig = (data) => {
  return service({
    url: '/docker/config',
    method: 'put',
    data,
    timeout: 15000, // 15秒超时
  }).catch(error => {
    console.error('更新Docker配置失败:', error)
    throw error
  })
}

// 验证Docker配置
export const validateDockerConfig = (data) => {
  return service({
    url: '/docker/config/validate',
    method: 'post',
    data,
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('验证Docker配置失败:', error)
    throw error
  })
}

// 备份Docker配置
export const backupDockerConfig = (description) => {
  return service({
    url: '/docker/config/backup',
    method: 'post',
    params: { description },
    timeout: 15000, // 15秒超时
  }).catch(error => {
    console.error('备份Docker配置失败:', error)
    throw error
  })
}

// 恢复Docker配置
export const restoreDockerConfig = (data) => {
  return service({
    url: '/docker/config/restore',
    method: 'post',
    data,
    timeout: 20000, // 20秒超时
  }).catch(error => {
    console.error('恢复Docker配置失败:', error)
    throw error
  })
}

// 获取备份列表
export const getBackupList = () => {
  return service({
    url: '/docker/config/backups',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取备份列表失败:', error)
    throw error
  })
}

// 删除备份
export const deleteBackup = (backupId) => {
  return service({
    url: `/docker/config/backups/${backupId}`,
    method: 'delete',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('删除备份失败:', error)
    throw error
  })
}

// 清理过期备份
export const cleanupOldBackups = (maxDays = 30, maxCount = 10) => {
  return service({
    url: '/docker/config/backups/cleanup',
    method: 'post',
    params: { maxDays, maxCount },
    timeout: 15000, // 15秒超时
  }).catch(error => {
    console.error('清理过期备份失败:', error)
    throw error
  })
}

// 重启Docker服务
export const restartDockerService = () => {
  return service({
    url: '/docker/service/restart',
    method: 'post',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    console.error('重启Docker服务失败:', error)
    throw error
  })
}

// 启动Docker服务
export const startDockerService = () => {
  return service({
    url: '/docker/service/start',
    method: 'post',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    console.error('启动Docker服务失败:', error)
    throw error
  })
}

// 停止Docker服务
export const stopDockerService = () => {
  return service({
    url: '/docker/service/stop',
    method: 'post',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    console.error('停止Docker服务失败:', error)
    throw error
  })
}

// 获取Docker服务状态
export const getDockerServiceStatus = () => {
  return service({
    url: '/docker/service/status',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取Docker服务状态失败:', error)
    throw error
  })
}

// 检查Docker服务健康状态
export const checkDockerServiceHealth = () => {
  return service({
    url: '/docker/service/health',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('检查Docker服务健康状态失败:', error)
    throw error
  })
}