import service from '@/utils/request'

// 获取Docker连接诊断结果
export const getDockerDiagnosis = () => {
  return service({
    url: '/docker/diagnose',
    method: 'get',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    console.error('获取Docker诊断失败:', error)
    throw error
  })
}

// 快速连接测试
export const testDockerConnection = () => {
  return service({
    url: '/docker/test-connection',
    method: 'post',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('Docker连接测试失败:', error)
    throw error
  })
}

// 获取权限状态
export const getDockerPermissions = () => {
  return service({
    url: '/docker/permissions',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取Docker权限状态失败:', error)
    throw error
  })
}