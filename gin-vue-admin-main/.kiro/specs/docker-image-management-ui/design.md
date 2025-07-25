# Docker镜像管理页面设计文档

## 概述

本设计文档描述了Docker镜像管理页面的技术架构、组件设计和实现方案。页面将基于Vue 3 + Element Plus构建，提供现代化的用户界面和完整的镜像管理功能。

## 架构设计

### 技术栈
- **前端框架**: Vue 3 (Composition API)
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **HTTP客户端**: Axios
- **路由**: Vue Router
- **构建工具**: Vite

### 页面结构
```
DockerImageManagement/
├── index.vue (主页面)
├── components/
│   ├── ImageTable.vue (镜像列表表格)
│   ├── ImageOperations.vue (操作按钮组)
│   ├── PullImageDialog.vue (拉取镜像对话框)
│   ├── BuildImageDialog.vue (构建镜像对话框)
│   ├── TagImageDialog.vue (标签镜像对话框)
│   └── BatchOperations.vue (批量操作组件)
└── composables/
    ├── useImageList.js (镜像列表逻辑)
    ├── useImageOperations.js (镜像操作逻辑)
    └── useImageDialogs.js (对话框状态管理)
```

## 组件设计

### 1. 主页面组件 (index.vue)

**功能职责**:
- 页面布局和整体状态管理
- 协调各子组件交互
- 处理页面级别的事件

**主要状态**:
```javascript
const state = {
  loading: false,
  selectedImages: [],
  searchKeyword: '',
  currentPage: 1,
  pageSize: 10,
  total: 0
}
```

### 2. 镜像列表表格 (ImageTable.vue)

**功能职责**:
- 展示镜像列表数据
- 处理排序、分页、选择
- 提供行级操作按钮

**表格列设计**:
- 选择框 (checkbox)
- 镜像ID (可点击查看详情)
- 名称/标签
- 大小 (格式化显示)
- 创建时间 (相对时间)
- 操作按钮 (标签、删除、导出)

### 3. 操作按钮组 (ImageOperations.vue)

**功能职责**:
- 提供顶部操作按钮
- 根据选择状态动态显示按钮
- 处理批量操作

**按钮设计**:
- 拉取镜像 (主要按钮)
- 构建镜像
- 刷新列表
- 清理悬空镜像
- 批量删除 (选择时显示)

### 4. 拉取镜像对话框 (PullImageDialog.vue)

**功能职责**:
- 收集拉取镜像参数
- 显示拉取进度
- 处理拉取结果

**表单字段**:
- 镜像名称 (必填，支持自动补全)
- 标签 (可选，默认latest)
- 进度显示区域

### 5. 构建镜像对话框 (BuildImageDialog.vue)

**功能职责**:
- 收集构建参数
- 提供Dockerfile编辑器
- 显示构建日志

**表单字段**:
- 镜像名称 (必填)
- 标签 (可选)
- Dockerfile内容 (代码编辑器)
- 构建参数 (键值对)
- 构建日志显示区域

## 数据流设计

### API接口映射
```javascript
// 镜像列表
GET /api/v1/docker/images
// 镜像详情
GET /api/v1/docker/images/:id
// 拉取镜像
POST /api/v1/docker/images/pull
// 删除镜像
DELETE /api/v1/docker/images/:id
// 构建镜像
POST /api/v1/docker/images/build
// 标签镜像
POST /api/v1/docker/images/tag
// 导出镜像
POST /api/v1/docker/images/export
// 清理镜像
POST /api/v1/docker/images/prune
```

### 状态管理
使用Pinia store管理全局状态：
```javascript
export const useImageStore = defineStore('image', {
  state: () => ({
    images: [],
    loading: false,
    selectedImages: [],
    filters: {
      name: '',
      dangling: null
    }
  }),
  actions: {
    async fetchImages(),
    async pullImage(),
    async deleteImage(),
    async buildImage(),
    async tagImage(),
    async exportImage(),
    async pruneImages()
  }
})
```

## UI/UX设计

### 布局设计
- **顶部**: 页面标题 + 操作按钮组
- **中间**: 搜索栏 + 过滤器
- **主体**: 镜像列表表格
- **底部**: 分页组件

### 交互设计
1. **加载状态**: 使用骨架屏和加载动画
2. **操作反馈**: Toast消息提示操作结果
3. **确认对话框**: 危险操作需要二次确认
4. **进度显示**: 长时间操作显示进度条

### 响应式设计
- **桌面端**: 完整表格显示，操作按钮横向排列
- **平板端**: 适当调整列宽，部分列可隐藏
- **移动端**: 卡片式布局，操作按钮下拉菜单

## 错误处理

### 错误类型
1. **网络错误**: 显示重试按钮
2. **权限错误**: 跳转登录页面
3. **业务错误**: 显示具体错误信息
4. **验证错误**: 表单字段高亮显示

### 错误处理策略
```javascript
const handleError = (error) => {
  if (error.response?.status === 401) {
    // 跳转登录
    router.push('/login')
  } else if (error.response?.status >= 500) {
    // 服务器错误
    ElMessage.error('服务器错误，请稍后重试')
  } else {
    // 业务错误
    ElMessage.error(error.response?.data?.msg || '操作失败')
  }
}
```

## 性能优化

### 优化策略
1. **虚拟滚动**: 大量数据时使用虚拟列表
2. **懒加载**: 镜像详情按需加载
3. **缓存策略**: 合理缓存镜像列表数据
4. **防抖搜索**: 搜索输入防抖处理

### 代码分割
```javascript
// 路由级别代码分割
const DockerImageManagement = () => import('./views/DockerImageManagement.vue')

// 组件级别懒加载
const BuildImageDialog = defineAsyncComponent(() => 
  import('./components/BuildImageDialog.vue')
)
```

## 测试策略

### 单元测试
- 组件渲染测试
- 用户交互测试
- API调用测试
- 状态管理测试

### 集成测试
- 页面完整流程测试
- API接口集成测试
- 错误处理测试

### E2E测试
- 关键用户路径测试
- 跨浏览器兼容性测试
- 响应式布局测试