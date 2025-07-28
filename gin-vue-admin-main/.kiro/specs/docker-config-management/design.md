# 设计文档

## 概述

Docker配置管理系统是一个用于管理Docker daemon配置的Web界面，允许用户通过图形界面查看、修改和应用Docker的各种配置参数。该系统提供了安全的配置管理机制，包括配置验证、备份恢复和服务重启等功能。

## 架构

### 系统架构
```
Docker配置管理系统
├── 配置读取层 (Configuration Reader)
│   ├── Docker配置文件解析
│   ├── Docker API配置获取
│   ├── 系统环境变量读取
│   └── 默认配置处理
├── 配置验证层 (Configuration Validator)
│   ├── 参数格式验证
│   ├── 配置兼容性检查
│   ├── 网络配置验证
│   └── 存储配置验证
├── 配置管理层 (Configuration Manager)
│   ├── 配置文件写入
│   ├── 配置备份管理
│   ├── 配置版本控制
│   └── 配置回滚机制
└── 服务控制层 (Service Controller)
    ├── Docker服务重启
    ├── 服务状态监控
    ├── 配置应用验证
    └── 错误恢复处理
```

### 技术栈
- **后端**: Go + Gin框架
- **Docker集成**: Docker SDK for Go
- **配置管理**: JSON/YAML配置文件处理
- **系统服务**: 系统服务管理API
- **前端**: Vue 3 + Element Plus

## 组件和接口

### 1. Docker配置服务 (DockerConfigService)

#### 核心方法
```go
type DockerConfigService struct{}

// 获取Docker配置
func (d *DockerConfigService) GetDockerConfig() (*DockerConfigResponse, error)

// 更新Docker配置
func (d *DockerConfigService) UpdateDockerConfig(config *DockerConfigRequest) error

// 验证Docker配置
func (d *DockerConfigService) ValidateDockerConfig(config *DockerConfigRequest) error

// 备份Docker配置
func (d *DockerConfigService) BackupDockerConfig() (string, error)

// 恢复Docker配置
func (d *DockerConfigService) RestoreDockerConfig(backupId string) error

// 重启Docker服务
func (d *DockerConfigService) RestartDockerService() error

// 获取Docker服务状态
func (d *DockerConfigService) GetDockerServiceStatus() (*ServiceStatusResponse, error)
```

#### 数据模型
```go
// Docker配置请求
type DockerConfigRequest struct {
    RegistryMirrors    []string          `json:"registryMirrors"`    // 镜像加速器
    InsecureRegistries []string          `json:"insecureRegistries"` // 不安全仓库
    PrivateRegistry    *PrivateRegistry  `json:"privateRegistry"`    // 私有仓库
    StorageDriver      string            `json:"storageDriver"`      // 存储驱动
    StorageOpts        map[string]string `json:"storageOpts"`        // 存储选项
    LogDriver          string            `json:"logDriver"`          // 日志驱动
    LogOpts            map[string]string `json:"logOpts"`            // 日志选项
    EnableIPv6         bool              `json:"enableIPv6"`         // 启用IPv6
    EnableIPForward    bool              `json:"enableIPForward"`    // 启用IP转发
    EnableIptables     bool              `json:"enableIptables"`     // 启用iptables
    LiveRestore        bool              `json:"liveRestore"`        // 实时恢复
    CgroupDriver       string            `json:"cgroupDriver"`       // Cgroup驱动
    SocketPath         string            `json:"socketPath"`         // Socket路径
    DataRoot           string            `json:"dataRoot"`           // 数据根目录
    ExecRoot           string            `json:"execRoot"`           // 执行根目录
}

// Docker配置响应
type DockerConfigResponse struct {
    Config          *DockerConfigRequest `json:"config"`
    ConfigPath      string               `json:"configPath"`
    IsDefault       bool                 `json:"isDefault"`
    LastModified    time.Time            `json:"lastModified"`
    ServiceStatus   string               `json:"serviceStatus"`
    Version         string               `json:"version"`
    BackupAvailable bool                 `json:"backupAvailable"`
}

// 私有仓库配置
type PrivateRegistry struct {
    URL      string `json:"url"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}

// 服务状态响应
type ServiceStatusResponse struct {
    Status      string    `json:"status"`      // running, stopped, error
    Uptime      string    `json:"uptime"`      // 运行时间
    Version     string    `json:"version"`     // Docker版本
    LastRestart time.Time `json:"lastRestart"` // 最后重启时间
    ErrorMsg    string    `json:"errorMsg"`    // 错误信息
}
```

### 2. 配置验证器 (DockerConfigValidator)

```go
type DockerConfigValidator struct{}

// 验证镜像加速器URL
func (v *DockerConfigValidator) ValidateRegistryMirrors(mirrors []string) error

// 验证私有仓库配置
func (v *DockerConfigValidator) ValidatePrivateRegistry(registry *PrivateRegistry) error

// 验证存储配置
func (v *DockerConfigValidator) ValidateStorageConfig(driver string, opts map[string]string) error

// 验证网络配置
func (v *DockerConfigValidator) ValidateNetworkConfig(config *DockerConfigRequest) error

// 验证Socket路径
func (v *DockerConfigValidator) ValidateSocketPath(path string) error

// 验证数据目录
func (v *DockerConfigValidator) ValidateDataRoot(path string) error
```

### 3. 配置文件管理器 (DockerConfigFileManager)

```go
type DockerConfigFileManager struct{}

// 读取配置文件
func (m *DockerConfigFileManager) ReadConfigFile() (*DockerConfigRequest, error)

// 写入配置文件
func (m *DockerConfigFileManager) WriteConfigFile(config *DockerConfigRequest) error

// 备份配置文件
func (m *DockerConfigFileManager) BackupConfigFile() (string, error)

// 恢复配置文件
func (m *DockerConfigFileManager) RestoreConfigFile(backupPath string) error

// 获取配置文件路径
func (m *DockerConfigFileManager) GetConfigFilePath() string

// 检查配置文件权限
func (m *DockerConfigFileManager) CheckConfigFilePermissions() error
```

### 4. Docker服务控制器 (DockerServiceController)

```go
type DockerServiceController struct{}

// 重启Docker服务
func (c *DockerServiceController) RestartService() error

// 停止Docker服务
func (c *DockerServiceController) StopService() error

// 启动Docker服务
func (c *DockerServiceController) StartService() error

// 获取服务状态
func (c *DockerServiceController) GetServiceStatus() (*ServiceStatusResponse, error)

// 检查服务健康状态
func (c *DockerServiceController) CheckServiceHealth() error

// 等待服务就绪
func (c *DockerServiceController) WaitForServiceReady(timeout time.Duration) error
```

### 5. API端点设计

```go
// 获取Docker配置
GET /api/v1/docker/config
Response: DockerConfigResponse

// 更新Docker配置
PUT /api/v1/docker/config
Body: DockerConfigRequest
Response: SuccessResponse

// 验证Docker配置
POST /api/v1/docker/config/validate
Body: DockerConfigRequest
Response: ValidationResponse

// 备份Docker配置
POST /api/v1/docker/config/backup
Response: BackupResponse

// 恢复Docker配置
POST /api/v1/docker/config/restore
Body: RestoreRequest
Response: SuccessResponse

// 重启Docker服务
POST /api/v1/docker/service/restart
Response: ServiceOperationResponse

// 获取Docker服务状态
GET /api/v1/docker/service/status
Response: ServiceStatusResponse

// 获取配置备份列表
GET /api/v1/docker/config/backups
Response: BackupListResponse
```

### 6. 前端组件设计

#### DockerConfigPanel.vue
```vue
<template>
  <div class="docker-config-panel">
    <!-- 配置表单 -->
    <el-form ref="configForm" :model="dockerConfig" :rules="configRules">
      <!-- 镜像加速器配置 -->
      <el-form-item label="镜像加速" prop="registryMirrors">
        <el-input
          v-model="dockerConfig.registryMirrors"
          type="textarea"
          placeholder="https://docker.1panel.live"
          :rows="3"
        />
      </el-form-item>

      <!-- 私有仓库配置 -->
      <el-form-item label="私有仓库">
        <el-input v-model="dockerConfig.privateRegistry.url" placeholder="仓库地址" />
      </el-form-item>

      <!-- 网络配置 -->
      <el-form-item label="IPv6">
        <el-switch v-model="dockerConfig.enableIPv6" />
      </el-form-item>

      <!-- 其他配置选项 -->
      <el-form-item label="iptables">
        <el-switch v-model="dockerConfig.enableIptables" />
      </el-form-item>

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="saveConfig">保存配置</el-button>
        <el-button @click="resetConfig">重置</el-button>
        <el-button @click="backupConfig">备份配置</el-button>
      </el-form-item>
    </el-form>

    <!-- 服务状态显示 -->
    <el-card class="service-status">
      <template #header>
        <span>Docker服务状态</span>
        <el-button style="float: right" @click="restartService">重启服务</el-button>
      </template>
      <div>
        <el-tag :type="serviceStatusType">{{ serviceStatus }}</el-tag>
        <span style="margin-left: 10px">版本: {{ dockerVersion }}</span>
      </div>
    </el-card>
  </div>
</template>
```

## 错误处理

### 1. 配置验证错误
```go
type ConfigValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

type ValidationResponse struct {
    Valid  bool                     `json:"valid"`
    Errors []ConfigValidationError  `json:"errors"`
}
```

### 2. 服务操作错误
```go
type ServiceOperationError struct {
    Operation string `json:"operation"`
    Message   string `json:"message"`
    Code      string `json:"code"`
    Retryable bool   `json:"retryable"`
}
```

## 安全考虑

### 1. 权限控制
- 配置文件读写权限检查
- Docker服务操作权限验证
- 用户角色权限控制

### 2. 配置安全
- 敏感信息加密存储
- 配置文件备份加密
- 网络配置安全验证

### 3. 操作审计
- 配置更改日志记录
- 服务操作审计
- 用户操作追踪

## 部署和监控

### 1. 配置文件位置
- **Linux**: `/etc/docker/daemon.json`
- **Windows**: `C:\ProgramData\docker\config\daemon.json`
- **macOS**: `~/.docker/daemon.json`

### 2. 服务管理
- **systemd**: `systemctl restart docker`
- **Windows Service**: `Restart-Service docker`
- **macOS**: Docker Desktop重启

### 3. 监控指标
- 配置更改频率
- 服务重启次数
- 配置验证失败率
- 服务健康状态