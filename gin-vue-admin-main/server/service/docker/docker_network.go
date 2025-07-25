package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerNetworkService struct{}

// GetNetworkList 获取网络列表
func (d *DockerNetworkService) GetNetworkList(filter dockerReq.NetworkFilter) ([]dockerRes.NetworkInfo, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return nil, 0, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 构建过滤器
	filterArgs := filters.NewArgs()
	if filter.Name != "" {
		filterArgs.Add("name", filter.Name)
	}
	if filter.Driver != "" {
		filterArgs.Add("driver", filter.Driver)
	}

	// 调用Docker API获取网络列表
	networks, err := global.GVA_DOCKER.NetworkList(ctx, types.NetworkListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		global.GVA_LOG.Error("Failed to get network list", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get network list: %v", err)
	}

	// 转换为响应模型
	networkInfos := make([]dockerRes.NetworkInfo, 0, len(networks))
	for _, dockerNetwork := range networks {
		networkInfo := d.convertToNetworkInfo(dockerNetwork)
		
		// 应用名称过滤（如果API过滤不够精确）
		if filter.Name != "" && !strings.Contains(strings.ToLower(networkInfo.Name), strings.ToLower(filter.Name)) {
			continue
		}
		
		networkInfos = append(networkInfos, networkInfo)
	}

	// 应用分页
	total := int64(len(networkInfos))
	start := (filter.Page - 1) * filter.PageSize
	end := start + filter.PageSize

	if start > len(networkInfos) {
		return []dockerRes.NetworkInfo{}, total, nil
	}
	if end > len(networkInfos) {
		end = len(networkInfos)
	}

	return networkInfos[start:end], total, nil
}

// GetNetworkDetail 获取网络详细信息
func (d *DockerNetworkService) GetNetworkDetail(networkID string) (*dockerRes.NetworkDetail, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用Docker API获取网络详细信息
	networkResource, err := global.GVA_DOCKER.NetworkInspect(ctx, networkID, types.NetworkInspectOptions{})
	if err != nil {
		global.GVA_LOG.Error("Failed to get network detail", zap.String("networkID", networkID), zap.Error(err))
		return nil, fmt.Errorf("failed to get network detail: %v", err)
	}

	// 转换为详细响应模型
	networkDetail := d.convertToNetworkDetail(networkResource)
	return &networkDetail, nil
}

// CreateNetwork 创建网络
func (d *DockerNetworkService) CreateNetwork(createReq dockerReq.NetworkCreateRequest) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return "", fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 构建网络创建选项
	createOptions := types.NetworkCreate{
		Driver:     createReq.Driver,
		EnableIPv6: createReq.EnableIPv6,
		Internal:   createReq.Internal,
		Attachable: createReq.Attachable,
		Labels:     createReq.Labels,
	}

	// 配置IPAM
	if createReq.Subnet != "" || createReq.Gateway != "" {
		createOptions.IPAM = &network.IPAM{
			Config: []network.IPAMConfig{{
				Subnet:  createReq.Subnet,
				Gateway: createReq.Gateway,
			}},
		}
	}

	// 创建网络
	response, err := global.GVA_DOCKER.NetworkCreate(ctx, createReq.Name, createOptions)
	if err != nil {
		global.GVA_LOG.Error("Failed to create network", 
			zap.String("name", createReq.Name), 
			zap.String("driver", createReq.Driver), 
			zap.Error(err))
		return "", fmt.Errorf("failed to create network: %v", err)
	}

	global.GVA_LOG.Info("Network created successfully", 
		zap.String("name", createReq.Name), 
		zap.String("id", response.ID))
	
	return response.ID, nil
}

// RemoveNetwork 删除网络
func (d *DockerNetworkService) RemoveNetwork(networkID string) error {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 删除网络
	err := global.GVA_DOCKER.NetworkRemove(ctx, networkID)
	if err != nil {
		global.GVA_LOG.Error("Failed to remove network", zap.String("networkID", networkID), zap.Error(err))
		return fmt.Errorf("failed to remove network: %v", err)
	}

	global.GVA_LOG.Info("Network removed successfully", zap.String("networkID", networkID))
	return nil
}

// PruneNetworks 清理未使用的网络
func (d *DockerNetworkService) PruneNetworks() (int64, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return 0, 0, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 清理网络
	pruneReport, err := global.GVA_DOCKER.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		global.GVA_LOG.Error("Failed to prune networks", zap.Error(err))
		return 0, 0, fmt.Errorf("failed to prune networks: %v", err)
	}

	deletedCount := int64(len(pruneReport.NetworksDeleted))
	spaceReclaimed := int64(0) // 网络清理不涉及磁盘空间

	global.GVA_LOG.Info("Networks pruned successfully", 
		zap.Int64("deletedCount", deletedCount))

	return deletedCount, spaceReclaimed, nil
}

// convertToNetworkInfo 将Docker API的NetworkResource转换为NetworkInfo响应模型
func (d *DockerNetworkService) convertToNetworkInfo(dockerNetwork types.NetworkResource) dockerRes.NetworkInfo {
	// 转换IPAM配置
	var ipamConfig *dockerRes.IPAMConfig
	if dockerNetwork.IPAM.Config != nil && len(dockerNetwork.IPAM.Config) > 0 {
		config := dockerNetwork.IPAM.Config[0]
		ipamConfig = &dockerRes.IPAMConfig{
			Subnet:  config.Subnet,
			Gateway: config.Gateway,
		}
	}

	return dockerRes.NetworkInfo{
		ID:         dockerNetwork.ID,
		Name:       dockerNetwork.Name,
		Driver:     dockerNetwork.Driver,
		Scope:      dockerNetwork.Scope,
		EnableIPv6: dockerNetwork.EnableIPv6,
		Internal:   dockerNetwork.Internal,
		Attachable: dockerNetwork.Attachable,
		Created:    time.Now().Unix(), // Docker API中的Created字段需要解析
		Labels:     dockerNetwork.Labels,
		IPAM:       ipamConfig,
	}
}

// convertToNetworkDetail 将Docker API的NetworkResource转换为NetworkDetail响应模型
func (d *DockerNetworkService) convertToNetworkDetail(dockerNetwork types.NetworkResource) dockerRes.NetworkDetail {
	// 转换基本信息
	networkInfo := d.convertToNetworkInfo(dockerNetwork)

	// 转换容器信息
	containers := make(map[string]dockerRes.NetworkContainer)
	for containerID, container := range dockerNetwork.Containers {
		containers[containerID] = dockerRes.NetworkContainer{
			Name:        container.Name,
			EndpointID:  container.EndpointID,
			MacAddress:  container.MacAddress,
			IPv4Address: container.IPv4Address,
			IPv6Address: container.IPv6Address,
		}
	}

	return dockerRes.NetworkDetail{
		NetworkInfo: networkInfo,
		Containers:  containers,
	}
}


