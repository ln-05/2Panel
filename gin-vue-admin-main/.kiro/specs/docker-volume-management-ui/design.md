# Design Document

## Overview

本设计文档描述了Docker存储卷管理前端页面的架构和实现方案。该页面将仿照现有的网络管理页面（network.vue），提供完整的存储卷管理功能，包括列表展示、创建、删除、详情查看和清理等操作。页面将与已有的后端API集成，提供一致的用户体验。

## Architecture

### 技术栈
- **前端框架**: Vue 3 Composition API
- **UI组件库**: Element Plus
- **状态管理**: Vue 3 Reactive API
- **HTTP客户端**: Axios (通过现有的service封装)
- **样式**: SCSS

### 文件结构
```
web/src/
├── api/
│   └── dockerVolume.js          # 存储卷API接口
├── view/
│   └── my/
│       └── volume/
│           └── volume.vue       # 存储卷管理页面
└── pathInfo.json               # 路由信息更新
```

## Components and Interfaces

### 主要组件结构

#### 1. 存储卷管理页面 (volume.vue)
- **组件名称**: DockerVolumeManagement
- **功能**: 存储卷的完整管理界面
- **主要区域**:
  - 操作栏 (operation-bar)
  - 数据表格 (table-container)
  - 创建对话框 (create-dialog)
  - 详情对话框 (detail-dialog)

#### 2. API接口层 (dockerVolume.js)
- **功能**: 封装所有存储卷相关的API调用
- **接口列表**:
  - `getDockerVolumeList(params)` - 获取存储卷列表
  - `getDockerVolumeDetail(name)` - 获取存储卷详情
  - `createDockerVolume(data)` - 创建存储卷
  - `deleteDockerVolume(name, force)` - 删除存储卷
  - `pruneDockerVolumes()` - 清理未使用存储卷

### 页面布局设计

#### 操作栏 (Operation Bar)
```
[创建存储卷] [刷新] [批量删除] [清理未使用存储卷]     [搜索框]
```

#### 数据表格列设计
| 列名 | 宽度 | 数据源 | 说明 |
|------|------|--------|------|
| 选择框 | 55px | - | 多选功能 |
| 名称 | 200px | name | 存储卷名称 |
| 驱动 | 120px | driver | 驱动类型，带标签样式 |
| 挂载点 | 250px | mountpoint | 挂载路径 |
| 范围 | 100px | scope | local/global |
| 标签 | 200px | labels | 键值对标签 |
| 创建时间 | 180px | createdAt | 格式化时间 |
| 操作 | 120px | - | 详情/删除按钮 |

## Data Models

### 前端数据模型

#### 存储卷列表项 (VolumeInfo)
```typescript
interface VolumeInfo {
  name: string;           // 存储卷名称
  driver: string;         // 存储卷驱动
  mountpoint: string;     // 挂载点
  scope: string;          // 范围
  createdAt: string;      // 创建时间
  labels: Record<string, string>; // 标签
  options: Record<string, string>; // 选项
}
```

#### 存储卷详情 (VolumeDetail)
```typescript
interface VolumeDetail extends VolumeInfo {
  usageData?: {
    size: number;         // 大小
    refCount: number;     // 引用计数
  };
}
```

#### 创建存储卷表单 (CreateVolumeForm)
```typescript
interface CreateVolumeForm {
  name: string;           // 存储卷名称 (必填)
  driver: string;         // 驱动类型 (默认local)
  driverOpts: Record<string, string>; // 驱动选项
  labels: Record<string, string>;     // 标签
}
```

## Error Handling

### 错误处理策略

#### 1. API错误处理
- **Docker服务不可用**: 显示"Docker服务不可用，请检查Docker配置和连接"
- **网络连接失败**: 显示"网络连接失败，请检查网络连接"
- **认证失败**: 显示"认证失败，请重新登录"
- **权限不足**: 显示"权限不足，无法访问Docker功能"
- **服务器错误**: 显示"服务器内部错误，请检查Docker服务配置"

#### 2. 用户操作错误处理
- **表单验证**: 实时验证必填字段
- **重复名称**: 显示"存储卷名称已存在，请使用其他名称"
- **删除失败**: 显示具体的删除失败原因
- **批量操作**: 显示操作进度和结果统计

## Testing Strategy

### 单元测试
- **API接口测试**: 测试所有API调用的正确性
- **组件功能测试**: 测试各个功能组件的行为
- **表单验证测试**: 测试表单验证逻辑
- **错误处理测试**: 测试各种错误场景的处理

### 集成测试
- **页面交互测试**: 测试用户操作流程
- **API集成测试**: 测试前后端接口集成
- **响应式测试**: 测试不同屏幕尺寸的适配

## UI/UX Design

### 视觉设计原则
- **一致性**: 与现有网络管理页面保持一致的视觉风格
- **清晰性**: 信息层次清晰，操作路径明确
- **响应式**: 适配不同设备屏幕尺寸
- **可访问性**: 支持键盘导航和屏幕阅读器

### 交互设计
- **即时反馈**: 所有操作都有即时的视觉反馈
- **确认机制**: 危险操作需要用户确认
- **批量操作**: 支持多选和批量处理
- **搜索过滤**: 实时搜索和过滤功能