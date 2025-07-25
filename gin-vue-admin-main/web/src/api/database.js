import service from '@/utils/request'

// 创建数据库连接
export const createDatabase = (data) => {
  return service({
    url: '/database/create',
    method: 'post',
    data
  })
}

// 获取数据库连接列表s
export const getDatabaseList = (params) => {
  return service({
    url: '/database/list',
    method: 'get',
    params
  })
}

// 获取单个数据库连接信息
export const getDatabaseById = (params) => {
  return service({
    url: '/database/id',
    method: 'get',
    params
  })
}

// 更新数据库连接
export const updateDatabase = (data) => {
  return service({
    url: '/database/update',
    method: 'post',
    data
  })
}

// 删除数据库连接
export const deleteDatabase = (data) => {
  return service({
    url: '/database/deleted',
    method: 'post',
    data
  })
}

// 测试数据库连接
export const testDatabase = (data) => {
  return service({
    url: '/database/test',
    method: 'post',
    data
  })
}

// 从服务器同步数据库
export const syncDatabase = (data) => {
  return service({
    url: '/database/sync',
    method: 'post',
    data
  })
}