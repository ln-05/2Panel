import service from '@/utils/request'

// 搜索模型列表
export const searchOllamaModel = (data) => {
  return service({
    url: '/ollamaModel/search',
    method: 'get',
    params: data
  })
}

// 创建/下载模型
export const createOllamaModelAdvanced = (data) => {
  return service({
    url: '/ollamaModel/create',
    method: 'post',
    data
  })
}

// 启动模型
export const startOllamaModel = (ID) => {
  return service({
    url: '/ollamaModel/start',
    method: 'post',
    params: { ID }
  })
}

// 停止模型
export const stopOllamaModel = (ID) => {
  return service({
    url: '/ollamaModel/stop',
    method: 'post',
    params: { ID }
  })
}

// 删除模型
export const deleteOllamaModel = (ID) => {
  return service({
    url: '/ollamaModel/deleteOllamaModel',
    method: 'delete',
    params: { ID }
  })
}

// 重新创建模型
export const recreateOllamaModel = (ID) => {
  return service({
    url: '/ollamaModel/recreate',
    method: 'post',
    params: { ID }
  })
}

// 同步模型
export const syncOllamaModel = (data) => {
  return service({
    url: '/ollamaModel/sync',
    method: 'post',
    data
  })
}

// 获取模型详情
export const getOllamaModelDetail = (ID) => {
  return service({
    url: '/ollamaModel/detail',
    method: 'get',
    params: { ID }
  })
}

// 获取模型日志
export const getOllamaModelLogs = (ID, lines = 100) => {
  return service({
    url: '/ollamaModel/logs',
    method: 'get',
    params: { ID, lines }
  })
}

// 获取系统资源状态
export const getSystemResourceStatus = () => {
  return service({
    url: '/ollamaModel/systemResource',
    method: 'get'
  })
}

// 批量删除模型
export const deleteOllamaModelByIds = (IDs) => {
  return service({
    url: '/ollamaModel/deleteOllamaModelByIds',
    method: 'delete',
    params: { 'IDs[]': IDs }
  })
}

// 与模型对话


// 与模型对话
export const chatWithOllamaModel = (data) => {
  return service({
    url: '/ollamaModel/chat',
    method: 'post',
    data
  })
}