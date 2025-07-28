// 数据库模块路由配置
export const databaseRoutes = [
  {
    path: '/database',
    name: 'Database',
    component: 'view/layout/index.vue',
    meta: {
      title: '数据库管理',
      icon: 'database',
      hideMenu: false
    },
    children: [
      {
        path: 'list',
        name: 'DatabaseList',
        component: 'view/database/index.vue',
        meta: {
          title: '数据库列表',
          icon: 'list',
          hideMenu: false
        }
      }
    ]
  }
] 