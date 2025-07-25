import service from '@/utils/request'

// 获取网络列表
export const getDockerNetworkList = (params) => {
  return service({
    url: '/docker/networks',
    method: 'get',
    params
  })
}

// 获取网络详细信息
export const getDockerNetworkDetail = (id) => {
  return service({
    url: `/docker/networks/${id}`,
    method: 'get'
  })
}

// 创建网络
export const createDockerNetwork = (data) => {
  return service({
    url: '/docker/networks',
    method: 'post',
    data
  })
}

// 删除网络
export const deleteDockerNetwork = (id) => {
  return service({
    url: `/docker/networks/${id}`,
    method: 'delete'
  })
}

// 清理未使用网络
export const pruneDockerNetworks = () => {
  return service({
    url: '/docker/networks/prune',
    method: 'post'
  })
}
