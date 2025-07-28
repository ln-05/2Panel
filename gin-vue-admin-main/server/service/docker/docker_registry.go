package docker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerRegistryService struct{}

// GetRegistryList 获取仓库列表 - 从Docker配置和实际使用情况获取
func (d *DockerRegistryService) GetRegistryList(filter dockerReq.RegistryFilter) ([]dockerRes.RegistryInfo, int64, error) {
	global.GVA_LOG.Info("开始从Docker获取仓库列表")
	
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker客户端不可用")
		return nil, 0, fmt.Errorf("Docker client is not available")
	}

	// 获取Docker系统信息
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := global.GVA_DOCKER.Info(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取Docker系统信息失败", zap.Error(err))
		return nil, 0, fmt.Errorf("获取Docker系统信息失败: %v", err)
	}

	// 构建仓库列表
	var registries []dockerRes.RegistryInfo
	registryMap := make(map[string]bool) // 用于去重
	
	// 添加默认的Docker Hub
	dockerHub := dockerRes.RegistryInfo{
		ID:          1,
		Name:        "Docker Hub",
		DownloadUrl: "https://registry-1.docker.io",
		Protocol:    "https",
		Status:      "active",
		Description: "Docker官方仓库",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	registries = append(registries, dockerHub)
	registryMap["https://registry-1.docker.io"] = true

	// 从Docker信息中获取配置的仓库镜像
	if info.RegistryConfig != nil && len(info.RegistryConfig.Mirrors) > 0 {
		global.GVA_LOG.Info("找到Docker镜像仓库", zap.Int("count", len(info.RegistryConfig.Mirrors)))
		for _, mirror := range info.RegistryConfig.Mirrors {
			if !registryMap[mirror] {
				registry := dockerRes.RegistryInfo{
					ID:          uint(len(registries) + 1),
					Name:        d.getRegistryName(mirror),
					DownloadUrl: mirror,
					Protocol:    d.getProtocol(mirror),
					Status:      "active",
					Description: "Docker仓库镜像",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				registries = append(registries, registry)
				registryMap[mirror] = true
			}
		}
	}

	// 从Docker信息中获取不安全的仓库
	if info.RegistryConfig != nil && len(info.RegistryConfig.InsecureRegistryCIDRs) > 0 {
		global.GVA_LOG.Info("找到不安全的Docker仓库", zap.Int("count", len(info.RegistryConfig.InsecureRegistryCIDRs)))
		for _, insecure := range info.RegistryConfig.InsecureRegistryCIDRs {
			insecureUrl := insecure.String()
			if !registryMap[insecureUrl] {
				registry := dockerRes.RegistryInfo{
					ID:          uint(len(registries) + 1),
					Name:        d.getRegistryName(insecureUrl),
					DownloadUrl: insecureUrl,
					Protocol:    "http",
					Status:      "active",
					Description: "不安全的Docker仓库",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				registries = append(registries, registry)
				registryMap[insecureUrl] = true
			}
		}
	}

	// 从现有镜像中推断使用的仓库
	images, err := global.GVA_DOCKER.ImageList(ctx, types.ImageListOptions{})
	if err == nil {
		global.GVA_LOG.Info("分析现有镜像以推断仓库", zap.Int("image_count", len(images)))
		for _, image := range images {
			for _, repoTag := range image.RepoTags {
				if repoTag != "<none>:<none>" {
					registryUrl := d.extractRegistryFromImage(repoTag)
					if registryUrl != "" && !registryMap[registryUrl] {
						registry := dockerRes.RegistryInfo{
							ID:          uint(len(registries) + 1),
							Name:        d.getRegistryName(registryUrl),
							DownloadUrl: registryUrl,
							Protocol:    d.getProtocol(registryUrl),
							Status:      "detected",
							Description: "从镜像推断的仓库",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						}
						registries = append(registries, registry)
						registryMap[registryUrl] = true
					}
				}
			}
		}
	}

	global.GVA_LOG.Info("从Docker获取到仓库数据", zap.Int("count", len(registries)))

	// 应用过滤条件
	var filteredRegistries []dockerRes.RegistryInfo
	for _, registry := range registries {
		// 名称过滤
		if filter.Name != "" && !strings.Contains(registry.Name, filter.Name) {
			continue
		}
		// 协议过滤
		if filter.Protocol != "" && registry.Protocol != filter.Protocol {
			continue
		}
		filteredRegistries = append(filteredRegistries, registry)
	}

	total := int64(len(filteredRegistries))

	// 应用分页
	if filter.Page > 0 && filter.PageSize > 0 {
		start := (filter.Page - 1) * filter.PageSize
		end := start + filter.PageSize
		if start >= len(filteredRegistries) {
			filteredRegistries = []dockerRes.RegistryInfo{}
		} else {
			if end > len(filteredRegistries) {
				end = len(filteredRegistries)
			}
			filteredRegistries = filteredRegistries[start:end]
		}
	}

	global.GVA_LOG.Info("返回过滤后的仓库数据", zap.Int("filtered_count", len(filteredRegistries)), zap.Int64("total", total))
	return filteredRegistries, total, nil
}

// GetRegistryDetail 获取仓库详细信息
func (d *DockerRegistryService) GetRegistryDetail(id uint) (*dockerRes.RegistryDetail, error) {
	// 从Docker API获取仓库列表
	registries, _, err := d.GetRegistryList(dockerReq.RegistryFilter{})
	if err != nil {
		return nil, err
	}

	// 查找指定ID的仓库
	for _, registry := range registries {
		if registry.ID == id {
			return &dockerRes.RegistryDetail{
				RegistryInfo: registry,
				IsDefault:    registry.ID == 1, // Docker Hub为默认仓库
				LastTestTime: nil,
				TestResult:   "",
			}, nil
		}
	}

	return nil, fmt.Errorf("仓库不存在")
}

// CreateRegistry 创建仓库 - Docker仓库配置应该在Docker配置文件中管理
func (d *DockerRegistryService) CreateRegistry(createReq dockerReq.RegistryCreateRequest) (*dockerRes.RegistryInfo, error) {
	return nil, fmt.Errorf("Docker仓库配置应该在Docker daemon配置文件中管理")
}

// UpdateRegistry 更新仓库 - Docker仓库配置应该在Docker配置文件中管理
func (d *DockerRegistryService) UpdateRegistry(updateReq dockerReq.RegistryUpdateRequest) (*dockerRes.RegistryInfo, error) {
	return nil, fmt.Errorf("Docker仓库配置应该在Docker daemon配置文件中管理")
}

// DeleteRegistry 删除仓库 - Docker仓库配置应该在Docker配置文件中管理
func (d *DockerRegistryService) DeleteRegistry(id uint) error {
	return fmt.Errorf("Docker仓库配置应该在Docker daemon配置文件中管理")
}

// TestRegistry 测试仓库连接
func (d *DockerRegistryService) TestRegistry(id uint) (*dockerRes.RegistryTestResponse, error) {
	// 获取仓库详情
	registryDetail, err := d.GetRegistryDetail(id)
	if err != nil {
		return nil, err
	}

	// 测试连接
	success, message := d.testRegistryConnection(registryDetail.RegistryInfo)

	return &dockerRes.RegistryTestResponse{
		Success: success,
		Message: message,
	}, nil
}

// SetDefaultRegistry 设置默认仓库 - Docker仓库配置应该在Docker配置文件中管理
func (d *DockerRegistryService) SetDefaultRegistry(id uint) error {
	return fmt.Errorf("Docker仓库配置应该在Docker daemon配置文件中管理")
}

// testRegistryConnection 测试仓库连接
func (d *DockerRegistryService) testRegistryConnection(registry dockerRes.RegistryInfo) (bool, string) {
	// 构建测试URL
	testUrl := registry.DownloadUrl
	if !strings.HasSuffix(testUrl, "/") {
		testUrl += "/"
	}
	testUrl += "v2/"

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequestWithContext(context.Background(), "GET", testUrl, nil)
	if err != nil {
		return false, fmt.Sprintf("创建请求失败: %v", err)
	}

	// 添加认证信息
	if registry.Username != "" {
		// 注意：这里没有密码信息，因为我们从Docker API获取的信息不包含密码
		// 实际测试时可能需要用户手动输入认证信息
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == 200 || resp.StatusCode == 401 {
		// 200表示成功，401表示需要认证但仓库存在
		return true, "连接成功"
	}

	return false, fmt.Sprintf("连接失败，状态码: %d", resp.StatusCode)
}

// getRegistryName 根据URL生成仓库名称
func (d *DockerRegistryService) getRegistryName(url string) string {
	// 移除协议前缀
	name := strings.TrimPrefix(url, "https://")
	name = strings.TrimPrefix(name, "http://")
	
	// 移除路径
	if idx := strings.Index(name, "/"); idx != -1 {
		name = name[:idx]
	}
	
	// 根据常见的仓库地址生成友好名称
	switch {
	case strings.Contains(name, "docker.io"):
		return "Docker Hub"
	case strings.Contains(name, "registry.cn-hangzhou.aliyuncs.com"):
		return "阿里云容器镜像服务"
	case strings.Contains(name, "ccr.ccs.tencentyun.com"):
		return "腾讯云容器镜像服务"
	case strings.Contains(name, "registry.cn-beijing.aliyuncs.com"):
		return "阿里云容器镜像服务(北京)"
	case strings.Contains(name, "registry.cn-shenzhen.aliyuncs.com"):
		return "阿里云容器镜像服务(深圳)"
	case strings.Contains(name, "dockerhub.azk8s.cn"):
		return "Azure中国镜像"
	case strings.Contains(name, "reg-mirror.qiniu.com"):
		return "七牛云镜像"
	case strings.Contains(name, "hub-mirror.c.163.com"):
		return "网易云镜像"
	case strings.Contains(name, "mirror.baidubce.com"):
		return "百度云镜像"
	default:
		return name
	}
}

// getProtocol 从URL中提取协议
func (d *DockerRegistryService) getProtocol(url string) string {
	if strings.HasPrefix(url, "https://") {
		return "https"
	} else if strings.HasPrefix(url, "http://") {
		return "http"
	}
	// 默认假设是https
	return "https"
}

// extractRegistryFromImage 从镜像名称中提取仓库地址
func (d *DockerRegistryService) extractRegistryFromImage(imageTag string) string {
	// 分割镜像名称和标签
	parts := strings.Split(imageTag, ":")
	imageName := parts[0]
	
	// 如果镜像名称包含斜杠，可能包含仓库地址
	if strings.Contains(imageName, "/") {
		parts := strings.Split(imageName, "/")
		
		// 检查第一部分是否像仓库地址（包含点号或端口）
		firstPart := parts[0]
		if strings.Contains(firstPart, ".") || strings.Contains(firstPart, ":") {
			// 构建完整的仓库URL
			if strings.Contains(firstPart, ":") && !strings.Contains(firstPart, "://") {
				// 包含端口，假设是http
				return "http://" + firstPart
			} else if !strings.Contains(firstPart, "://") {
				// 不包含协议，假设是https
				return "https://" + firstPart
			}
			return firstPart
		}
	}
	
	// 如果没有明确的仓库地址，返回空字符串（表示使用默认的Docker Hub）
	return ""
}
