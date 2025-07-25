import service from '@/utils/request'

// 获取存储卷列表
export const getDockerVolumeList = (params) => {
  return service({
    url: '/docker/volumes',
    method: 'get',
    params
  })
}

// 获取存储卷详细信息
export const getDockerVolumeDetail = (name) => {
  return service({
    url: `/docker/volumes/${name}`,
    method: 'get'
  })
}

// 创建存储卷
export const createDockerVolume = (data) => {
  return service({
    url: '/docker/volumes',
    method: 'post',
    data
  })
}

// 删除存储卷
export const deleteDockerVolume = (name, force = false) => {
  return service({
    url: `/docker/volumes/${name}`,
    method: 'delete',
    params: { force }
  })
}

// 清理未使用存储卷
export const pruneDockerVolumes = () => {
  return service({
    url: '/docker/volumes/prune',
    method: 'post'
  })
}