import service from '@/utils/request'

// 获取仓库列表
export const getRegistryList = (data) => {
  return service({
    url: '/docker/registries',
    method: 'get',
    params: data
  })
}

// 获取仓库详情
export const getRegistryDetail = (id) => {
  return service({
    url: `/docker/registries/${id}`,
    method: 'get'
  })
}

// 创建仓库
export const createRegistry = (data) => {
  return service({
    url: '/docker/registries',
    method: 'post',
    data
  })
}

// 更新仓库
export const updateRegistry = (data) => {
  return service({
    url: '/docker/registries',
    method: 'put',
    data
  })
}

// 删除仓库
export const deleteRegistry = (id) => {
  return service({
    url: `/docker/registries/${id}`,
    method: 'delete'
  })
}

// 测试仓库连接
export const testRegistry = (id) => {
  return service({
    url: `/docker/registries/${id}/test`,
    method: 'post'
  })
}

// 设置默认仓库
export const setDefaultRegistry = (id) => {
  return service({
    url: `/docker/registries/${id}/default`,
    method: 'post'
  })
}

