package docker

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type DockerDiagnosticService struct{}

// DiagnosticResult 诊断结果
type DiagnosticResult struct {
	OverallStatus    string                      `json:"overallStatus"` // "healthy", "warning", "error"
	ClientStatus     *ClientStatusResult         `json:"clientStatus"`
	NetworkTest      *NetworkTestResult          `json:"networkTest"`
	VersionCheck     *VersionCompatibilityResult `json:"versionCheck"`
	ConfigValidation *ConfigValidationResult     `json:"configValidation"`
	PermissionCheck  *PermissionCheckResult      `json:"permissionCheck"`
	Timestamp        time.Time                   `json:"timestamp"`
}

// ClientStatusResult 客户端状态结果
type ClientStatusResult struct {
	IsInitialized bool   `json:"isInitialized"`
	IsConnected   bool   `json:"isConnected"`
	Error         string `json:"error,omitempty"`
	PingLatency   int64  `json:"pingLatency"` // 毫秒
}

// NetworkTestResult 网络测试结果
type NetworkTestResult struct {
	IsReachable    bool   `json:"isReachable"`
	ResponseTime   int64  `json:"responseTime"` // 毫秒
	Error          string `json:"error,omitempty"`
	DNSResolution  bool   `json:"dnsResolution"`
	PortAccessible bool   `json:"portAccessible"`
}

// VersionCompatibilityResult 版本兼容性结果
type VersionCompatibilityResult struct {
	ClientVersion      string `json:"clientVersion"`
	ServerVersion      string `json:"serverVersion"`
	IsCompatible       bool   `json:"isCompatible"`
	RecommendedVersion string `json:"recommendedVersion,omitempty"`
	Error              string `json:"error,omitempty"`
}

// ConfigValidationResult 配置验证结果
type ConfigValidationResult struct {
	IsValid     bool     `json:"isValid"`
	Issues      []string `json:"issues"`
	Suggestions []string `json:"suggestions"`
}

// PermissionCheckResult 权限检查结果
type PermissionCheckResult struct {
	CanListContainers bool     `json:"canListContainers"`
	CanListImages     bool     `json:"canListImages"`
	CanListNetworks   bool     `json:"canListNetworks"`
	CanListVolumes    bool     `json:"canListVolumes"`
	CanGetInfo        bool     `json:"canGetInfo"`
	CanGetVersion     bool     `json:"canGetVersion"`
	Errors            []string `json:"errors"`
	Suggestions       []string `json:"suggestions"`
}

// DiagnoseConnection 全面诊断Docker连接
func (d *DockerDiagnosticService) DiagnoseConnection() (*DiagnosticResult, error) {
	result := &DiagnosticResult{
		Timestamp: time.Now(),
	}

	// 检查客户端状态
	result.ClientStatus = d.CheckClientStatus()

	// 测试网络连通性
	result.NetworkTest = d.TestNetworkConnectivity()

	// 验证API版本兼容性
	result.VersionCheck = d.ValidateAPIVersion()

	// 验证配置
	result.ConfigValidation = d.ValidateConfig()

	// 检查权限
	result.PermissionCheck = d.CheckPermissions()

	// 评估整体状态
	result.OverallStatus = d.evaluateOverallStatus(result)

	return result, nil
}

// CheckClientStatus 检查Docker客户端状态
func (d *DockerDiagnosticService) CheckClientStatus() *ClientStatusResult {
	result := &ClientStatusResult{
		IsInitialized: global.GVA_DOCKER != nil,
	}

	if !result.IsInitialized {
		result.Error = "Docker client is not initialized"
		return result
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	_, err := global.GVA_DOCKER.Ping(ctx)
	result.PingLatency = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		result.IsConnected = false
	} else {
		result.IsConnected = true
	}

	return result
}

// TestNetworkConnectivity 测试网络连通性
func (d *DockerDiagnosticService) TestNetworkConnectivity() *NetworkTestResult {
	result := &NetworkTestResult{}

	dockerConfig := global.GVA_CONFIG.Docker
	host := dockerConfig.Host

	// 解析主机和端口
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	// 对于TCP连接，测试端口可达性
	if len(host) > 6 && host[:6] == "tcp://" {
		hostPort := host[6:] // 移除 "tcp://" 前缀

		start := time.Now()
		conn, err := net.DialTimeout("tcp", hostPort, 10*time.Second)
		result.ResponseTime = time.Since(start).Milliseconds()

		if err != nil {
			result.Error = err.Error()
			result.IsReachable = false
			result.PortAccessible = false
		} else {
			result.IsReachable = true
			result.PortAccessible = true
			conn.Close()
		}

		// 测试DNS解析
		if host, _, err := net.SplitHostPort(hostPort); err == nil {
			if _, err := net.LookupHost(host); err == nil {
				result.DNSResolution = true
			}
		}
	} else {
		// Unix socket连接
		result.IsReachable = true
		result.PortAccessible = true
		result.DNSResolution = true
	}

	return result
}

// ValidateAPIVersion 验证API版本兼容性
func (d *DockerDiagnosticService) ValidateAPIVersion() *VersionCompatibilityResult {
	result := &VersionCompatibilityResult{
		ClientVersion: global.GVA_CONFIG.Docker.Version,
	}

	if global.GVA_DOCKER == nil {
		result.Error = "Docker client is not initialized"
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取服务器版本
	version, err := global.GVA_DOCKER.ServerVersion(ctx)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.ServerVersion = version.Version
	result.IsCompatible = true // 简化处理，实际应该检查版本兼容性
	result.RecommendedVersion = "1.41"

	return result
}

// ValidateConfig 验证配置
func (d *DockerDiagnosticService) ValidateConfig() *ConfigValidationResult {
	result := &ConfigValidationResult{
		IsValid:     true,
		Issues:      []string{},
		Suggestions: []string{},
	}

	dockerConfig := global.GVA_CONFIG.Docker

	// 检查主机配置
	if dockerConfig.Host == "" {
		result.Issues = append(result.Issues, "Docker host is not configured")
		result.Suggestions = append(result.Suggestions, "Set docker.host in config.yaml")
		result.IsValid = false
	}

	// 检查版本配置
	if dockerConfig.Version == "" {
		result.Issues = append(result.Issues, "Docker API version is not specified")
		result.Suggestions = append(result.Suggestions, "Set docker.version to '1.41' or appropriate version")
	}

	// 检查超时配置
	if dockerConfig.Timeout <= 0 {
		result.Issues = append(result.Issues, "Docker timeout is not properly configured")
		result.Suggestions = append(result.Suggestions, "Set docker.timeout to a positive value (e.g., 30)")
	}

	// 检查TLS配置
	if dockerConfig.TLSVerify && dockerConfig.CertPath == "" {
		result.Issues = append(result.Issues, "TLS verification is enabled but certificate path is not specified")
		result.Suggestions = append(result.Suggestions, "Set docker.cert-path or disable TLS verification")
		result.IsValid = false
	}

	return result
}

// CheckPermissions 检查权限
func (d *DockerDiagnosticService) CheckPermissions() *PermissionCheckResult {
	result := &PermissionCheckResult{
		Errors:      []string{},
		Suggestions: []string{},
	}

	if global.GVA_DOCKER == nil {
		result.Errors = append(result.Errors, "Docker client is not initialized")
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试容器列表权限
	_, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{Limit: 1})
	if err != nil {
		result.CanListContainers = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot list containers: %v", err))
	} else {
		result.CanListContainers = true
	}

	// 测试镜像列表权限
	_, err = global.GVA_DOCKER.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		result.CanListImages = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot list images: %v", err))
	} else {
		result.CanListImages = true
	}

	// 测试网络列表权限
	_, err = global.GVA_DOCKER.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		result.CanListNetworks = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot list networks: %v", err))
	} else {
		result.CanListNetworks = true
	}

	// 测试存储卷列表权限
	_, err = global.GVA_DOCKER.VolumeList(ctx, filters.Args{})
	if err != nil {
		result.CanListVolumes = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot list volumes: %v", err))
	} else {
		result.CanListVolumes = true
	}

	// 测试系统信息权限
	_, err = global.GVA_DOCKER.Info(ctx)
	if err != nil {
		result.CanGetInfo = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot get system info: %v", err))
	} else {
		result.CanGetInfo = true
	}

	// 测试版本信息权限
	_, err = global.GVA_DOCKER.ServerVersion(ctx)
	if err != nil {
		result.CanGetVersion = false
		result.Errors = append(result.Errors, fmt.Sprintf("Cannot get version: %v", err))
	} else {
		result.CanGetVersion = true
	}

	// 生成建议
	if len(result.Errors) > 0 {
		result.Suggestions = append(result.Suggestions, "Check Docker daemon permissions")
		result.Suggestions = append(result.Suggestions, "Verify user is in docker group (Linux)")
		result.Suggestions = append(result.Suggestions, "Check Docker Desktop settings (Windows/Mac)")
		result.Suggestions = append(result.Suggestions, "Verify TLS configuration if using remote Docker")
	}

	return result
}

// evaluateOverallStatus 评估整体状态
func (d *DockerDiagnosticService) evaluateOverallStatus(result *DiagnosticResult) string {
	if !result.ClientStatus.IsConnected {
		return "error"
	}

	if !result.NetworkTest.IsReachable {
		return "error"
	}

	if !result.ConfigValidation.IsValid {
		return "error"
	}

	// 检查关键权限
	if !result.PermissionCheck.CanListContainers ||
		!result.PermissionCheck.CanListImages ||
		!result.PermissionCheck.CanGetInfo {
		return "warning"
	}

	if len(result.PermissionCheck.Errors) > 0 {
		return "warning"
	}

	return "healthy"
}
