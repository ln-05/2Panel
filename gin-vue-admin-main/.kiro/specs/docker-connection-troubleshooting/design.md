# 设计文档

## 概述

Docker连接故障排除系统是一个综合性的诊断和修复工具，旨在解决用户在使用Docker编排界面时遇到的"获取容器失败"等连接问题。该系统通过多层次的检测机制、详细的错误分析和智能的修复建议，帮助用户快速定位和解决Docker连接问题。

## 架构

### 系统架构
```
Docker故障排除系统
├── 连接诊断层 (Connection Diagnosis Layer)
│   ├── Docker客户端状态检查
│   ├── 网络连通性测试
│   ├── API版本兼容性检查
│   └── TLS/证书验证
├── 配置验证层 (Configuration Validation Layer)
│   ├── 配置文件解析
│   ├── 参数有效性验证
│   ├── 环境变量检查
│   └── 权限验证
├── 故障分析层 (Fault Analysis Layer)
│   ├── 错误分类器
│   ├── 根因分析引擎
│   ├── 解决方案匹配器
│   └── 修复建议生成器
└── 监控和日志层 (Monitoring & Logging Layer)
    ├── 实时状态监控
    ├── 连接状态变化通知
    ├── 详细日志记录
    └── 性能指标收集
```

### 技术栈
- **后端**: Go + Gin框架
- **Docker集成**: Docker SDK for Go
- **网络诊断**: Go标准库 net/http, net
- **日志系统**: Zap日志库
- **配置管理**: Viper配置库
- **前端**: Vue 3 + Element Plus

## 组件和接口

### 1. Docker连接诊断服务 (DockerDiagnosticService)

#### 核心方法
```go
type DockerDiagnosticService struct{}

// 全面诊断Docker连接
func (d *DockerDiagnosticService) DiagnoseConnection() (*DiagnosticResult, error)

// 检查Docker客户端状态
func (d *DockerDiagnosticService) CheckClientStatus() *ClientStatusResult

// 测试网络连通性
func (d *DockerDiagnosticService) TestNetworkConnectivity(host string, port int) *NetworkTestResult

// 验证API版本兼容性
func (d *DockerDiagnosticService) ValidateAPIVersion() *VersionCompatibilityResult

// 检查TLS配置
func (d *DockerDiagnosticService) ValidateTLSConfig() *TLSValidationResult

// 生成修复建议
func (d *DockerDiagnosticService) GenerateFixSuggestions(diagnosticResult *DiagnosticResult) []FixSuggestion
```

#### 数据模型
```go
// 诊断结果
type DiagnosticResult struct {
    OverallStatus    string                    `json:"overallStatus"`    // "healthy", "warning", "error"
    ClientStatus     *ClientStatusResult       `json:"clientStatus"`
    NetworkTest      *NetworkTestResult        `json:"networkTest"`
    VersionCheck     *VersionCompatibilityResult `json:"versionCheck"`
    TLSValidation    *TLSValidationResult      `json:"tlsValidation"`
    ConfigValidation *ConfigValidationResult   `json:"configValidation"`
    Timestamp        time.Time                 `json:"timestamp"`
}

// 客户端状态结果
type ClientStatusResult struct {
    IsInitialized bool   `json:"isInitialized"`
    IsConnected   bool   `json:"isConnected"`
    Error         string `json:"error,omitempty"`
    PingLatency   int64  `json:"pingLatency"` // 毫秒
}

// 网络测试结果
type NetworkTestResult struct {
    IsReachable     bool   `json:"isReachable"`
    ResponseTime    int64  `json:"responseTime"` // 毫秒
    Error           string `json:"error,omitempty"`
    DNSResolution   bool   `json:"dnsResolution"`
    PortAccessible  bool   `json:"portAccessible"`
}

// 版本兼容性结果
type VersionCompatibilityResult struct {
    ClientVersion    string `json:"clientVersion"`
    ServerVersion    string `json:"serverVersion"`
    IsCompatible     bool   `json:"isCompatible"`
    RecommendedVersion string `json:"recommendedVersion,omitempty"`
}

// TLS验证结果
type TLSValidationResult struct {
    IsTLSEnabled     bool   `json:"isTlsEnabled"`
    CertificateValid bool   `json:"certificateValid"`
    Error            string `json:"error,omitempty"`
}

// 修复建议
type FixSuggestion struct {
    Category    string `json:"category"`    // "configuration", "network", "service", "permission"
    Priority    string `json:"priority"`    // "high", "medium", "low"
    Title       string `json:"title"`
    Description string `json:"description"`
    Steps       []string `json:"steps"`
    AutoFixable bool   `json:"autoFixable"`
}
```

### 2. 配置验证服务 (DockerConfigValidator)

```go
type DockerConfigValidator struct{}

// 验证Docker配置
func (v *DockerConfigValidator) ValidateConfig(config *DockerConfig) *ConfigValidationResult

// 检查主机地址格式
func (v *DockerConfigValidator) ValidateHostFormat(host string) error

// 验证API版本
func (v *DockerConfigValidator) ValidateAPIVersion(version string) error

// 检查证书路径
func (v *DockerConfigValidator) ValidateCertPath(certPath string) error

// 验证超时设置
func (v *DockerConfigValidator) ValidateTimeout(timeout int) error
```

### 3. 网络诊断工具 (NetworkDiagnosticTool)

```go
type NetworkDiagnosticTool struct{}

// TCP连接测试
func (n *NetworkDiagnosticTool) TestTCPConnection(host string, port int, timeout time.Duration) error

// DNS解析测试
func (n *NetworkDiagnosticTool) TestDNSResolution(hostname string) error

// HTTP连接测试
func (n *NetworkDiagnosticTool) TestHTTPConnection(url string, timeout time.Duration) (*http.Response, error)

// 端口可达性测试
func (n *NetworkDiagnosticTool) TestPortReachability(host string, port int) bool
```

### 4. API端点设计

```go
// 获取Docker连接诊断结果
GET /api/v1/docker/diagnosis
Response: DiagnosticResult

// 执行快速连接测试
POST /api/v1/docker/test-connection
Response: ConnectionTestResult

// 获取修复建议
GET /api/v1/docker/fix-suggestions
Response: []FixSuggestion

// 应用自动修复
POST /api/v1/docker/auto-fix
Body: AutoFixRequest
Response: AutoFixResult

// 获取Docker配置验证结果
GET /api/v1/docker/config-validation
Response: ConfigValidationResult

// 实时连接状态监控 (WebSocket)
WS /api/v1/docker/status-monitor
```

### 5. 前端组件设计

#### DockerDiagnosticPanel.vue
```vue
<template>
  <div class="docker-diagnostic-panel">
    <!-- 总体状态显示 -->
    <el-card class="status-overview">
      <div class="status-indicator" :class="overallStatus">
        <el-icon><Connection /></el-icon>
        <span>{{ statusText }}</span>
      </div>
    </el-card>

    <!-- 详细诊断结果 -->
    <el-tabs v-model="activeTab">
      <el-tab-pane label="连接状态" name="connection">
        <ConnectionStatus :result="diagnosticResult.clientStatus" />
      </el-tab-pane>
      <el-tab-pane label="网络测试" name="network">
        <NetworkTest :result="diagnosticResult.networkTest" />
      </el-tab-pane>
      <el-tab-pane label="配置验证" name="config">
        <ConfigValidation :result="diagnosticResult.configValidation" />
      </el-tab-pane>
      <el-tab-pane label="修复建议" name="suggestions">
        <FixSuggestions :suggestions="fixSuggestions" @apply-fix="applyFix" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
```

## 数据模型

### 配置模型扩展
```go
type DockerConfig struct {
    Host       string `yaml:"host" json:"host"`
    Version    string `yaml:"version" json:"version"`
    TLSVerify  bool   `yaml:"tls-verify" json:"tlsVerify"`
    CertPath   string `yaml:"cert-path" json:"certPath"`
    Timeout    int    `yaml:"timeout" json:"timeout"`
    
    // 新增诊断相关配置
    EnableDiagnostic     bool `yaml:"enable-diagnostic" json:"enableDiagnostic"`
    DiagnosticInterval   int  `yaml:"diagnostic-interval" json:"diagnosticInterval"` // 秒
    AutoRetryCount       int  `yaml:"auto-retry-count" json:"autoRetryCount"`
    HealthCheckEndpoint  string `yaml:"health-check-endpoint" json:"healthCheckEndpoint"`
}
```

### 错误分类系统
```go
type DockerErrorCategory string

const (
    ErrorCategoryConnection    DockerErrorCategory = "connection"
    ErrorCategoryConfiguration DockerErrorCategory = "configuration"
    ErrorCategoryPermission    DockerErrorCategory = "permission"
    ErrorCategoryVersion       DockerErrorCategory = "version"
    ErrorCategoryNetwork       DockerErrorCategory = "network"
    ErrorCategoryService       DockerErrorCategory = "service"
)

type DockerError struct {
    Category    DockerErrorCategory `json:"category"`
    Code        string             `json:"code"`
    Message     string             `json:"message"`
    Details     string             `json:"details"`
    Suggestions []string           `json:"suggestions"`
}
```

## 错误处理

### 1. 连接错误处理策略
```go
// 连接错误处理器
type ConnectionErrorHandler struct{}

func (h *ConnectionErrorHandler) HandleError(err error) *DockerError {
    switch {
    case isConnectionRefused(err):
        return &DockerError{
            Category: ErrorCategoryConnection,
            Code:     "CONNECTION_REFUSED",
            Message:  "Docker守护进程连接被拒绝",
            Details:  err.Error(),
            Suggestions: []string{
                "检查Docker Desktop是否已启动",
                "验证Docker守护进程是否正在运行",
                "检查防火墙设置",
            },
        }
    case isTimeout(err):
        return &DockerError{
            Category: ErrorCategoryNetwork,
            Code:     "CONNECTION_TIMEOUT",
            Message:  "Docker连接超时",
            Details:  err.Error(),
            Suggestions: []string{
                "检查网络连接",
                "增加连接超时时间",
                "验证Docker主机地址是否正确",
            },
        }
    // ... 其他错误类型处理
    }
}
```

### 2. 自动重试机制
```go
type RetryConfig struct {
    MaxRetries    int           `json:"maxRetries"`
    InitialDelay  time.Duration `json:"initialDelay"`
    MaxDelay      time.Duration `json:"maxDelay"`
    BackoffFactor float64       `json:"backoffFactor"`
}

func (d *DockerDiagnosticService) ConnectWithRetry(config RetryConfig) error {
    var lastErr error
    delay := config.InitialDelay
    
    for i := 0; i <= config.MaxRetries; i++ {
        if err := d.testConnection(); err == nil {
            return nil
        } else {
            lastErr = err
            if i < config.MaxRetries {
                time.Sleep(delay)
                delay = time.Duration(float64(delay) * config.BackoffFactor)
                if delay > config.MaxDelay {
                    delay = config.MaxDelay
                }
            }
        }
    }
    return lastErr
}
```

## 测试策略

### 1. 单元测试
- Docker客户端状态检查测试
- 网络连通性测试模拟
- 配置验证逻辑测试
- 错误分类和处理测试

### 2. 集成测试
- 完整诊断流程测试
- API端点集成测试
- WebSocket连接监控测试
- 自动修复功能测试

### 3. 模拟测试场景
```go
// 测试场景定义
type TestScenario struct {
    Name        string
    Setup       func() error
    Teardown    func() error
    Expected    DiagnosticResult
}

var testScenarios = []TestScenario{
    {
        Name: "Docker未启动",
        Setup: func() error {
            // 模拟Docker服务停止
            return stopDockerService()
        },
        Expected: DiagnosticResult{
            OverallStatus: "error",
            ClientStatus: &ClientStatusResult{
                IsInitialized: true,
                IsConnected:   false,
                Error:        "connection refused",
            },
        },
    },
    // ... 其他测试场景
}
```

## 性能优化

### 1. 诊断缓存机制
```go
type DiagnosticCache struct {
    cache    map[string]*CachedResult
    mutex    sync.RWMutex
    ttl      time.Duration
}

type CachedResult struct {
    Result    *DiagnosticResult
    Timestamp time.Time
}

func (c *DiagnosticCache) Get(key string) (*DiagnosticResult, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    
    if cached, exists := c.cache[key]; exists {
        if time.Since(cached.Timestamp) < c.ttl {
            return cached.Result, true
        }
        delete(c.cache, key)
    }
    return nil, false
}
```

### 2. 并发诊断
```go
func (d *DockerDiagnosticService) DiagnoseConnectionConcurrent() (*DiagnosticResult, error) {
    var wg sync.WaitGroup
    result := &DiagnosticResult{}
    
    // 并发执行各项检查
    wg.Add(4)
    
    go func() {
        defer wg.Done()
        result.ClientStatus = d.CheckClientStatus()
    }()
    
    go func() {
        defer wg.Done()
        result.NetworkTest = d.TestNetworkConnectivity(config.Host, config.Port)
    }()
    
    go func() {
        defer wg.Done()
        result.VersionCheck = d.ValidateAPIVersion()
    }()
    
    go func() {
        defer wg.Done()
        result.TLSValidation = d.ValidateTLSConfig()
    }()
    
    wg.Wait()
    
    // 综合评估整体状态
    result.OverallStatus = d.evaluateOverallStatus(result)
    result.Timestamp = time.Now()
    
    return result, nil
}
```

## 安全考虑

### 1. 敏感信息保护
- Docker证书和密钥的安全存储
- 配置信息的脱敏显示
- 日志中敏感信息的过滤

### 2. 权限控制
- 诊断功能的访问权限控制
- 自动修复功能的权限验证
- 配置修改的审计日志

## 部署和监控

### 1. 健康检查端点
```go
func (d *DockerDiagnosticService) HealthCheck() map[string]interface{} {
    return map[string]interface{}{
        "docker_client_available": d.IsDockerAvailable(),
        "last_diagnostic_time":    d.getLastDiagnosticTime(),
        "diagnostic_cache_size":   d.getCacheSize(),
        "active_monitors":         d.getActiveMonitorCount(),
    }
}
```

### 2. 监控指标
- Docker连接成功率
- 诊断执行时间
- 错误类型分布
- 自动修复成功率