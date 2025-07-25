import service from '@/utils/request'

// 获取镜像列表
export const getDockerImageList = (data) => {
  return service({
    url: '/docker/images',
    method: 'get',
    params: data
  })
}

// 获取镜像详情
export const getDockerImageDetail = (id) => {
  return service({
    url: `/docker/images/${id}`,
    method: 'get'
  })
}

// 拉取镜像
export const pullDockerImage = (data) => {
  return service({
    url: '/docker/images/pull',
    method: 'post',
    data
  })
}

// 删除镜像
export const deleteDockerImage = (id, force = false) => {
  return service({
    url: `/docker/images/${id}`,
    method: 'delete',
    params: { force }
  })
}

// 构建镜像
export const buildDockerImage = (data) => {
  return service({
    url: '/docker/images/build',
    method: 'post',
    data
  })
}

// 给镜像打标签
export const tagDockerImage = (data) => {
  return service({
    url: '/docker/images/tag',
    method: 'post',
    data
  })
}

// 清理未使用的镜像
export const pruneDockerImages = (dangling = true) => {
  return service({
    url: '/docker/images/prune',
    method: 'post',
    params: { dangling }
  })
}

// 导出镜像
export const exportDockerImage = (data) => {
  return service({
    url: '/docker/images/export',
    method: 'post',
    data,
    responseType: 'blob'
  })
}

// 导入镜像
export const importDockerImage = (data) => {
  return service({
    url: '/docker/images/import',
    method: 'post',
    data
  })
}