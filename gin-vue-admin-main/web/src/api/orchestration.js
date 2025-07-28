import service from '@/utils/request'

// 获取编排列表
export const getOrchestrationList = (data) => {
  return service({
    url: '/orchestration/list', // 请根据实际后端路由调整
    method: 'get', // 或 'get'，根据后端实际接口
    data
  })
}

// 删除编排
export const deleteOrchestration = (id) => {
  return service({
    url: `/orchestration/delete`, // 请根据实际后端路由调整
    method: 'post',
    data: { id }
  })
}
