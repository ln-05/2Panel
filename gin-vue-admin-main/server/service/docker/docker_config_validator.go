package docker

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"go.uber.org/zap"
)

type DockerConfigValidator struct{}

// ValidateDockerConfig 验证完整的Docker配置
func (v *DockerConfigValidator) ValidateDockerConfig(config *dockerModel.DockerConfigRequest) *dockerModel.ValidationResponse {
	var errors []dockerModel.ConfigValidationError

	// 验证镜像加速器
	if err := v.ValidateRegistryMirrors(config.RegistryMirrors); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "registryMirrors",
			Message: err.Error(),
			Code:    "INVALID_REGISTRY_MIRRORS",
		})
	}

	// 验证不安全仓库
	if err := v.ValidateInsecureRegistries(config.InsecureRegistries); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "insecureRegistries",
			Message: err.Error(),
			Code:    "INVALID_INSECURE_REGISTRIES",
		})
	}

	// 验证私有仓库
	if config.PrivateRegistry != nil {
		if err := v.ValidatePrivateRegistry(config.PrivateRegistry); err != nil {
			errors = append(errors, dockerModel.ConfigValidationError{
				Field:   "privateRegistry",
				Message: err.Error(),
				Code:    "INVALID_PRIVATE_REGISTRY",
			})
		}
	}

	// 验证存储配置
	if err := v.ValidateStorageConfig(config.StorageDriver, config.StorageOpts); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "storageDriver",
			Message: err.Error(),
			Code:    "INVALID_STORAGE_CONFIG",
		})
	}

	// 验证日志配置
	if err := v.ValidateLogConfig(config.LogDriver, config.LogOpts); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "logDriver",
			Message: err.Error(),
			Code:    "INVALID_LOG_CONFIG",
		})
	}

	// 验证网络配置
	if err := v.ValidateNetworkConfig(config); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "networkConfig",
			Message: err.Error(),
			Code:    "INVALID_NETWORK_CONFIG",
		})
	}

	// 验证Cgroup驱动
	if err := v.ValidateCgroupDriver(config.CgroupDriver); err != nil {
		errors = append(errors, dockerModel.ConfigValidationError{
			Field:   "cgroupDriver",
			Message: err.Error(),
			Code:    "INVALID_CGROUP_DRIVER",
		})
	}

	// 验证Socket路径
	if config.SocketPath != "" {
		if err := v.ValidateSocketPath(config.SocketPath); err != nil {
			errors = append(errors, dockerModel.ConfigValidationError{
				Field:   "socketPath",
				Message: err.Error(),
				Code:    "INVALID_SOCKET_PATH",
			})
		}
	}

	// 验证数据根目录
	if config.DataRoot != "" {
		if err := v.ValidateDataRoot(config.DataRoot); err != nil {
			errors = append(errors, dockerModel.ConfigValidationError{
				Field:   "dataRoot",
				Message: err.Error(),
				Code:    "INVALID_DATA_ROOT",
			})
		}
	}

	// 验证执行根目录
	if config.ExecRoot != "" {
		if err := v.ValidateExecRoot(config.ExecRoot); err != nil {
			errors = append(errors, dockerModel.ConfigValidationError{
				Field:   "execRoot",
				Message: err.Error(),
				Code:    "INVALID_EXEC_ROOT",
			})
		}
	}

	response := &dockerModel.ValidationResponse{
		Valid:  len(errors) == 0,
		Errors: errors,
	}

	if len(errors) > 0 {
		global.GVA_LOG.Warn("Docker配置验证失败", zap.Int("errorCount", len(errors)))
	} else {
		global.GVA_LOG.Info("Docker配置验证通过")
	}

	return response
}

// ValidateRegistryMirrors 验证镜像加速器URL
func (v *DockerConfigValidator) ValidateRegistryMirrors(mirrors []string) error {
	for i, mirror := range mirrors {
		if mirror == "" {
			return fmt.Errorf("镜像加速器URL不能为空 (索引: %d)", i)
		}

		// 验证URL格式
		parsedURL, err := url.Parse(mirror)
		if err != nil {
			return fmt.Errorf("镜像加速器URL格式无效: %s (索引: %d)", mirror, i)
		}

		// 检查协议
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return fmt.Errorf("镜像加速器URL必须使用http或https协议: %s (索引: %d)", mirror, i)
		}

		// 检查主机名
		if parsedURL.Host == "" {
			return fmt.Errorf("镜像加速器URL缺少主机名: %s (索引: %d)", mirror, i)
		}
	}

	return nil
}

// ValidateInsecureRegistries 验证不安全仓库
func (v *DockerConfigValidator) ValidateInsecureRegistries(registries []string) error {
	for i, registry := range registries {
		if registry == "" {
			return fmt.Errorf("不安全仓库地址不能为空 (索引: %d)", i)
		}

		// 验证格式：可以是域名、IP地址或域名:端口、IP:端口
		if !v.isValidRegistryAddress(registry) {
			return fmt.Errorf("不安全仓库地址格式无效: %s (索引: %d)", registry, i)
		}
	}

	return nil
}

// ValidatePrivateRegistry 验证私有仓库配置
func (v *DockerConfigValidator) ValidatePrivateRegistry(registry *dockerModel.PrivateRegistry) error {
	if registry.URL == "" {
		return fmt.Errorf("私有仓库URL不能为空")
	}

	// 验证URL格式
	parsedURL, err := url.Parse(registry.URL)
	if err != nil {
		return fmt.Errorf("私有仓库URL格式无效: %s", registry.URL)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("私有仓库URL必须使用http或https协议: %s", registry.URL)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("私有仓库URL缺少主机名: %s", registry.URL)
	}

	// 如果提供了用户名，密码也应该提供
	if registry.Username != "" && registry.Password == "" {
		return fmt.Errorf("提供用户名时必须同时提供密码")
	}

	// 验证邮箱格式（如果提供）
	if registry.Email != "" && !v.isValidEmail(registry.Email) {
		return fmt.Errorf("邮箱格式无效: %s", registry.Email)
	}

	return nil
}

// ValidateStorageConfig 验证存储配置
func (v *DockerConfigValidator) ValidateStorageConfig(driver string, opts map[string]string) error {
	if driver == "" {
		return nil // 空驱动使用默认值
	}

	// 支持的存储驱动列表
	supportedDrivers := []string{
		"overlay2", "aufs", "devicemapper", "btrfs", "zfs", "vfs",
	}

	found := false
	for _, supported := range supportedDrivers {
		if driver == supported {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("不支持的存储驱动: %s，支持的驱动: %s", driver, strings.Join(supportedDrivers, ", "))
	}

	// 验证特定驱动的选项
	switch driver {
	case "devicemapper":
		return v.validateDeviceMapperOpts(opts)
	case "btrfs":
		return v.validateBtrfsOpts(opts)
	case "zfs":
		return v.validateZfsOpts(opts)
	}

	return nil
}

// ValidateLogConfig 验证日志配置
func (v *DockerConfigValidator) ValidateLogConfig(driver string, opts map[string]string) error {
	if driver == "" {
		return nil // 空驱动使用默认值
	}

	// 支持的日志驱动列表
	supportedDrivers := []string{
		"json-file", "syslog", "journald", "gelf", "fluentd", "awslogs", "splunk", "none",
	}

	found := false
	for _, supported := range supportedDrivers {
		if driver == supported {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("不支持的日志驱动: %s，支持的驱动: %s", driver, strings.Join(supportedDrivers, ", "))
	}

	// 验证特定驱动的选项
	switch driver {
	case "json-file":
		return v.validateJsonFileLogOpts(opts)
	case "syslog":
		return v.validateSyslogOpts(opts)
	}

	return nil
}

// ValidateNetworkConfig 验证网络配置
func (v *DockerConfigValidator) ValidateNetworkConfig(config *dockerModel.DockerConfigRequest) error {
	// IPv6和IP转发的组合验证
	if config.EnableIPv6 && !config.EnableIPForward {
		global.GVA_LOG.Warn("启用IPv6时建议同时启用IP转发")
	}

	// 检查iptables依赖
	if !config.EnableIptables && config.EnableIPForward {
		return fmt.Errorf("启用IP转发时通常需要启用iptables")
	}

	return nil
}

// ValidateCgroupDriver 验证Cgroup驱动
func (v *DockerConfigValidator) ValidateCgroupDriver(driver string) error {
	if driver == "" {
		return nil // 空驱动使用默认值
	}

	supportedDrivers := []string{"cgroupfs", "systemd"}
	
	found := false
	for _, supported := range supportedDrivers {
		if driver == supported {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("不支持的Cgroup驱动: %s，支持的驱动: %s", driver, strings.Join(supportedDrivers, ", "))
	}

	return nil
}

// ValidateSocketPath 验证Socket路径
func (v *DockerConfigValidator) ValidateSocketPath(path string) error {
	if path == "" {
		return fmt.Errorf("Socket路径不能为空")
	}

	// 检查路径格式
	if !filepath.IsAbs(path) {
		return fmt.Errorf("Socket路径必须是绝对路径: %s", path)
	}

	// 检查目录是否存在
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("Socket目录不存在: %s", dir)
	}

	return nil
}

// ValidateDataRoot 验证数据根目录
func (v *DockerConfigValidator) ValidateDataRoot(path string) error {
	if path == "" {
		return fmt.Errorf("数据根目录不能为空")
	}

	// 检查路径格式
	if !filepath.IsAbs(path) {
		return fmt.Errorf("数据根目录必须是绝对路径: %s", path)
	}

	// 检查目录权限
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("数据根路径不是目录: %s", path)
		}
		// 检查写权限
		testFile := filepath.Join(path, ".docker_test")
		if file, err := os.Create(testFile); err != nil {
			return fmt.Errorf("数据根目录不可写: %s", path)
		} else {
			file.Close()
			os.Remove(testFile)
		}
	}

	return nil
}

// ValidateExecRoot 验证执行根目录
func (v *DockerConfigValidator) ValidateExecRoot(path string) error {
	if path == "" {
		return fmt.Errorf("执行根目录不能为空")
	}

	// 检查路径格式
	if !filepath.IsAbs(path) {
		return fmt.Errorf("执行根目录必须是绝对路径: %s", path)
	}

	return nil
}

// 辅助方法

// isValidRegistryAddress 验证仓库地址格式
func (v *DockerConfigValidator) isValidRegistryAddress(address string) bool {
	// 域名或IP地址的正则表达式
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	ipRegex := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
	
	// 分离主机和端口
	parts := strings.Split(address, ":")
	host := parts[0]
	
	// 验证主机部分
	if !domainRegex.MatchString(host) && !ipRegex.MatchString(host) {
		return false
	}
	
	// 如果有端口，验证端口范围
	if len(parts) == 2 {
		port := parts[1]
		portRegex := regexp.MustCompile(`^[1-9][0-9]{0,4}$`)
		if !portRegex.MatchString(port) {
			return false
		}
	}
	
	return len(parts) <= 2
}

// isValidEmail 验证邮箱格式
func (v *DockerConfigValidator) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validateDeviceMapperOpts 验证devicemapper存储驱动选项
func (v *DockerConfigValidator) validateDeviceMapperOpts(opts map[string]string) error {
	// 验证常见的devicemapper选项
	if basesize, ok := opts["dm.basesize"]; ok {
		if !regexp.MustCompile(`^\d+[KMGT]?$`).MatchString(basesize) {
			return fmt.Errorf("dm.basesize格式无效: %s", basesize)
		}
	}
	
	return nil
}

// validateBtrfsOpts 验证btrfs存储驱动选项
func (v *DockerConfigValidator) validateBtrfsOpts(opts map[string]string) error {
	// btrfs通常不需要特殊选项验证
	return nil
}

// validateZfsOpts 验证zfs存储驱动选项
func (v *DockerConfigValidator) validateZfsOpts(opts map[string]string) error {
	// 验证ZFS存储池名称
	if zfsname, ok := opts["zfs.fsname"]; ok {
		if zfsname == "" {
			return fmt.Errorf("zfs.fsname不能为空")
		}
	}
	
	return nil
}

// validateJsonFileLogOpts 验证json-file日志驱动选项
func (v *DockerConfigValidator) validateJsonFileLogOpts(opts map[string]string) error {
	// 验证最大文件大小
	if maxSize, ok := opts["max-size"]; ok {
		if !regexp.MustCompile(`^\d+[kmg]?$`).MatchString(maxSize) {
			return fmt.Errorf("max-size格式无效: %s", maxSize)
		}
	}
	
	// 验证最大文件数量
	if maxFile, ok := opts["max-file"]; ok {
		if !regexp.MustCompile(`^\d+$`).MatchString(maxFile) {
			return fmt.Errorf("max-file必须是数字: %s", maxFile)
		}
	}
	
	return nil
}

// validateSyslogOpts 验证syslog日志驱动选项
func (v *DockerConfigValidator) validateSyslogOpts(opts map[string]string) error {
	// 验证syslog地址
	if address, ok := opts["syslog-address"]; ok {
		if address != "" {
			if _, err := url.Parse(address); err != nil {
				return fmt.Errorf("syslog-address格式无效: %s", address)
			}
		}
	}
	
	return nil
}