import service from '@/utils/request'

// 获取容器列表
export const getContainerList = (params) => {
  return service({
    url: '/docker/containers',
    method: 'get',
    params
  })
}

// 获取容器详细信息
export const getContainerDetail = (id) => {
  return service({
    url: `/docker/containers/${id}`,
    method: 'get'
  })
}

// 获取容器日志
export const getContainerLogs = (id, params) => {
  return service({
    url: `/docker/containers/${id}/logs`,
    method: 'get',
    params
  })
}

// 启动容器
export const startContainer = (id) => {
  return service({
    url: `/docker/containers/${id}/start`,
    method: 'post'
  })
}

// 停止容器
export const stopContainer = (id, timeout) => {
  return service({
    url: `/docker/containers/${id}/stop`,
    method: 'post',
    params: timeout ? { timeout } : {}
  })
}

// 重启容器
export const restartContainer = (id, timeout) => {
  return service({
    url: `/docker/containers/${id}/restart`,
    method: 'post',
    params: timeout ? { timeout } : {}
  })
}

// 删除容器
export const deleteContainer = (id, force = false) => {
  return service({
    url: `/docker/containers/${id}`,
    method: 'delete',
    params: { force }
  })
}

// 获取Docker系统信息
export const getDockerInfo = () => {
  return service({
    url: '/docker/info',
    method: 'get'
  })
}

// 检查Docker状态
export const checkDockerStatus = () => {
  return service({
    url: '/docker/status',
    method: 'get'
  })
} 