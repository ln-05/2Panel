package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerVolumeService struct{}

// GetVolumeList 获取存储卷列表
func (d *DockerVolumeService) GetVolumeList(filter dockerReq.VolumeFilter) ([]dockerRes.VolumeInfo, int64, error) {
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

	// 调用Docker API获取存储卷列表
	volumeListResponse, err := global.GVA_DOCKER.VolumeList(ctx, filters.Args{})
	if err != nil {
		global.GVA_LOG.Error("Failed to get volume list", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get volume list: %v", err)
	}

	// 转换为响应模型
	volumeInfos := make([]dockerRes.VolumeInfo, 0, len(volumeListResponse.Volumes))
	for _, dockerVolume := range volumeListResponse.Volumes {
		volumeInfo := d.convertToVolumeInfo(dockerVolume)
		
		// 应用名称过滤（如果API过滤不够精确）
		if filter.Name != "" && !containsIgnoreCase(volumeInfo.Name, filter.Name) {
			continue
		}
		
		// 应用驱动过滤
		if filter.Driver != "" && volumeInfo.Driver != filter.Driver {
			continue
		}
		
		volumeInfos = append(volumeInfos, volumeInfo)
	}

	// 应用分页
	total := int64(len(volumeInfos))
	start := (filter.Page - 1) * filter.PageSize
	end := start + filter.PageSize

	if start > len(volumeInfos) {
		return []dockerRes.VolumeInfo{}, total, nil
	}
	if end > len(volumeInfos) {
		end = len(volumeInfos)
	}

	return volumeInfos[start:end], total, nil
}

// GetVolumeDetail 获取存储卷详细信息
func (d *DockerVolumeService) GetVolumeDetail(volumeName string) (*dockerRes.VolumeDetail, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用Docker API获取存储卷详细信息
	dockerVolume, err := global.GVA_DOCKER.VolumeInspect(ctx, volumeName)
	if err != nil {
		global.GVA_LOG.Error("Failed to get volume detail", zap.String("volumeName", volumeName), zap.Error(err))
		return nil, fmt.Errorf("failed to get volume detail: %v", err)
	}

	// 转换为详细响应模型
	volumeDetail := d.convertToVolumeDetail(dockerVolume)
	return &volumeDetail, nil
}

// CreateVolume 创建存储卷
func (d *DockerVolumeService) CreateVolume(createReq dockerReq.VolumeCreateRequest) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return "", fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 设置默认驱动
	if createReq.Driver == "" {
		createReq.Driver = "local"
	}

	// 构建存储卷创建选项
	createOptions := volume.VolumeCreateBody{
		Name:       createReq.Name,
		Driver:     createReq.Driver,
		DriverOpts: createReq.DriverOpts,
		Labels:     createReq.Labels,
	}

	// 创建存储卷
	dockerVolume, err := global.GVA_DOCKER.VolumeCreate(ctx, createOptions)
	if err != nil {
		global.GVA_LOG.Error("Failed to create volume", 
			zap.String("name", createReq.Name), 
			zap.String("driver", createReq.Driver), 
			zap.Error(err))
		return "", fmt.Errorf("failed to create volume: %v", err)
	}

	global.GVA_LOG.Info("Volume created successfully", 
		zap.String("name", createReq.Name), 
		zap.String("mountpoint", dockerVolume.Mountpoint))
	
	return dockerVolume.Name, nil
}

// RemoveVolume 删除存储卷
func (d *DockerVolumeService) RemoveVolume(volumeName string, force bool) error {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 删除存储卷
	err := global.GVA_DOCKER.VolumeRemove(ctx, volumeName, force)
	if err != nil {
		global.GVA_LOG.Error("Failed to remove volume", zap.String("volumeName", volumeName), zap.Error(err))
		return fmt.Errorf("failed to remove volume: %v", err)
	}

	global.GVA_LOG.Info("Volume removed successfully", zap.String("volumeName", volumeName))
	return nil
}

// PruneVolumes 清理未使用的存储卷
func (d *DockerVolumeService) PruneVolumes() (int64, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		global.GVA_LOG.Error("Docker client is not available")
		return 0, 0, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 清理存储卷
	pruneReport, err := global.GVA_DOCKER.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		global.GVA_LOG.Error("Failed to prune volumes", zap.Error(err))
		return 0, 0, fmt.Errorf("failed to prune volumes: %v", err)
	}

	deletedCount := int64(len(pruneReport.VolumesDeleted))
	spaceReclaimed := int64(pruneReport.SpaceReclaimed)

	global.GVA_LOG.Info("Volumes pruned successfully", 
		zap.Int64("deletedCount", deletedCount),
		zap.Int64("spaceReclaimed", spaceReclaimed))

	return deletedCount, spaceReclaimed, nil
}

// convertToVolumeInfo 将Docker API的Volume转换为VolumeInfo响应模型
func (d *DockerVolumeService) convertToVolumeInfo(dockerVolume *types.Volume) dockerRes.VolumeInfo {
	return dockerRes.VolumeInfo{
		Name:       dockerVolume.Name,
		Driver:     dockerVolume.Driver,
		Mountpoint: dockerVolume.Mountpoint,
		Scope:      dockerVolume.Scope,
		CreatedAt:  dockerVolume.CreatedAt,
		Labels:     dockerVolume.Labels,
		Options:    dockerVolume.Options,
	}
}

// convertToVolumeDetail 将Docker API的Volume转换为VolumeDetail响应模型
func (d *DockerVolumeService) convertToVolumeDetail(dockerVolume types.Volume) dockerRes.VolumeDetail {
	// 转换基本信息
	volumeInfo := d.convertToVolumeInfo(&dockerVolume)

	// 转换使用情况数据
	var usageData *dockerRes.VolumeUsageData
	if dockerVolume.UsageData != nil {
		usageData = &dockerRes.VolumeUsageData{
			Size:     dockerVolume.UsageData.Size,
			RefCount: dockerVolume.UsageData.RefCount,
		}
	}

	return dockerRes.VolumeDetail{
		VolumeInfo: volumeInfo,
		UsageData:  usageData,
	}
}

// containsIgnoreCase 不区分大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(substr) == 0 || 
		    (len(s) > 0 && len(substr) > 0 && 
		     strings.Contains(strings.ToLower(s), strings.ToLower(substr))))
}

