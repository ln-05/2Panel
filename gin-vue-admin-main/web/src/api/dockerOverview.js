import service from '@/utils/request'

// 获取Docker概览统计信息
export const getDockerOverview = () => {
  return service({
    url: '/docker/overview',
    method: 'get',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    // 统一错误处理
    console.error('获取Docker概览统计失败:', error)
    throw error
  })
}

// 获取Docker配置摘要信息
export const getDockerConfigSummary = () => {
  return service({
    url: '/docker/config/summary',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取Docker配置摘要失败:', error)
    throw error
  })
}

// 获取Docker磁盘使用情况
export const getDockerDiskUsage = () => {
  return service({
    url: '/docker/disk-usage',
    method: 'get',
    timeout: 30000, // 30秒超时
  }).catch(error => {
    console.error('获取Docker磁盘使用情况失败:', error)
    throw error
  })
}

// 批量获取所有概览数据
export const getAllOverviewData = async () => {
  try {
    const [overviewStats, configSummary, diskUsage] = await Promise.allSettled([
      getDockerOverview(),
      getDockerConfigSummary(),
      getDockerDiskUsage()
    ])

    return {
      overviewStats: overviewStats.status === 'fulfilled' ? overviewStats.value : null,
      configSummary: configSummary.status === 'fulfilled' ? configSummary.value : null,
      diskUsage: diskUsage.status === 'fulfilled' ? diskUsage.value : null,
      errors: {
        overviewStats: overviewStats.status === 'rejected' ? overviewStats.reason : null,
        configSummary: configSummary.status === 'rejected' ? configSummary.reason : null,
        diskUsage: diskUsage.status === 'rejected' ? diskUsage.reason : null
      }
    }
  } catch (error) {
    console.error('批量获取概览数据失败:', error)
    throw error
  }
}

// 带重试机制的API调用
export const getDockerOverviewWithRetry = async (maxRetries = 3, delay = 1000) => {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await getDockerOverview()
    } catch (error) {
      console.warn(`获取Docker概览统计失败，第${i + 1}次重试:`, error)
      
      if (i === maxRetries - 1) {
        throw error
      }
      
      // 等待后重试
      await new Promise(resolve => setTimeout(resolve, delay * (i + 1)))
    }
  }
}

// 检查Docker服务状态
export const checkDockerStatus = () => {
  return service({
    url: '/docker/status',
    method: 'get',
    timeout: 5000, // 5秒超时
  }).catch(error => {
    console.error('检查Docker状态失败:', error)
    throw error
  })
}

// 获取Docker系统信息
export const getDockerInfo = () => {
  return service({
    url: '/docker/info',
    method: 'get',
    timeout: 10000, // 10秒超时
  }).catch(error => {
    console.error('获取Docker系统信息失败:', error)
    throw error
  })
}