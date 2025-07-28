package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerOverviewService struct{}

// GetOverviewStats 获取Docker概览统计信息
func (d *DockerOverviewService) GetOverviewStats() (*response.OverviewStats, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 检查Docker连接
	_, err := global.GVA_DOCKER.Ping(ctx)
	if err != nil {
		global.GVA_LOG.Error("Docker ping failed", zap.Error(err))
		return nil, fmt.Errorf("Docker服务连接失败: %w", err)
	}

	stats := &response.OverviewStats{}

	// 获取容器统计
	containerStats, err := d.getContainerStats(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get container stats", zap.Error(err))
		// 不阻塞其他统计，使用默认值
		stats.Containers = response.ContainerStats{}
	} else {
		stats.Containers = *containerStats
	}

	// 获取镜像统计
	imageStats, err := d.getImageStats(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get image stats", zap.Error(err))
		stats.Images = response.ImageStats{}
	} else {
		stats.Images = *imageStats
	}

	// 获取网络统计
	networkStats, err := d.getNetworkStats(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get network stats", zap.Error(err))
		stats.Networks = response.NetworkStats{}
	} else {
		stats.Networks = *networkStats
	}

	// 获取存储卷统计
	volumeStats, err := d.getVolumeStats(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get volume stats", zap.Error(err))
		stats.Volumes = response.VolumeStats{}
	} else {
		stats.Volumes = *volumeStats
	}

	// 获取系统信息
	systemStats, err := d.getSystemStats(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get system stats", zap.Error(err))
		stats.System = response.SystemStats{}
	} else {
		stats.System = *systemStats
	}

	global.GVA_LOG.Info("Docker overview stats collected successfully")
	return stats, nil
}

// getContainerStats 获取容器统计信息
func (d *DockerOverviewService) getContainerStats(ctx context.Context) (*response.ContainerStats, error) {
	// 获取所有容器（包括停止的）
	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	stats := &response.ContainerStats{
		Total: len(containers),
	}

	// 统计不同状态的容器
	for _, container := range containers {
		switch container.State {
		case "running":
			stats.Running++
		case "exited", "stopped", "dead":
			stats.Stopped++
		case "paused":
			stats.Paused++
		case "restarting":
			// 重启中的容器暂时归类为运行中
			stats.Running++
		case "created":
			// 已创建但未启动的容器归类为停止
			stats.Stopped++
		default:
			// 其他未知状态归类为停止
			global.GVA_LOG.Warn("Unknown container state", zap.String("state", container.State), zap.String("containerID", container.ID))
			stats.Stopped++
		}
	}

	global.GVA_LOG.Debug("Container stats collected",
		zap.Int("total", stats.Total),
		zap.Int("running", stats.Running),
		zap.Int("stopped", stats.Stopped),
		zap.Int("paused", stats.Paused))

	return stats, nil
}

// getImageStats 获取镜像统计信息
func (d *DockerOverviewService) getImageStats(ctx context.Context) (*response.ImageStats, error) {
	// 获取所有镜像（不包括中间层镜像）
	images, err := global.GVA_DOCKER.ImageList(ctx, types.ImageListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	stats := &response.ImageStats{
		Total: len(images),
	}

	// 计算总大小
	var totalSize int64
	for _, image := range images {
		totalSize += image.Size
	}

	stats.SizeBytes = totalSize
	stats.Size = d.formatBytes(totalSize)

	global.GVA_LOG.Debug("Image stats collected",
		zap.Int("total", stats.Total),
		zap.Int64("sizeBytes", stats.SizeBytes),
		zap.String("size", stats.Size))

	return stats, nil
}

// getNetworkStats 获取网络统计信息
func (d *DockerOverviewService) getNetworkStats(ctx context.Context) (*response.NetworkStats, error) {
	// 获取所有网络
	networks, err := global.GVA_DOCKER.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	stats := &response.NetworkStats{
		Total: len(networks),
	}

	// 统计不同类型的网络
	for _, network := range networks {
		switch network.Driver {
		case "bridge":
			stats.Bridge++
		case "host":
			stats.Host++
		case "none":
			stats.None++
		default:
			// 其他类型的网络（如overlay、macvlan等）
			global.GVA_LOG.Debug("Other network driver found", zap.String("driver", network.Driver), zap.String("networkID", network.ID))
		}
	}

	global.GVA_LOG.Debug("Network stats collected",
		zap.Int("total", stats.Total),
		zap.Int("bridge", stats.Bridge),
		zap.Int("host", stats.Host),
		zap.Int("none", stats.None))

	return stats, nil
}

// getVolumeStats 获取存储卷统计信息
func (d *DockerOverviewService) getVolumeStats(ctx context.Context) (*response.VolumeStats, error) {
	// 获取所有存储卷
	volumeResponse, err := global.GVA_DOCKER.VolumeList(ctx, filters.Args{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	stats := &response.VolumeStats{
		Total: len(volumeResponse.Volumes),
	}

	// 尝试计算存储卷总大小
	var totalSize int64
	for _, volume := range volumeResponse.Volumes {
		if volume.Mountpoint != "" {
			size, err := d.getDirectorySize(volume.Mountpoint)
			if err != nil {
				global.GVA_LOG.Debug("Failed to get volume size", zap.String("volume", volume.Name), zap.Error(err))
				continue
			}
			totalSize += size
		}
	}

	stats.SizeBytes = totalSize
	if totalSize > 0 {
		stats.Size = d.formatBytes(totalSize)
	} else {
		stats.Size = "未知"
	}

	global.GVA_LOG.Debug("Volume stats collected",
		zap.Int("total", stats.Total),
		zap.Int64("sizeBytes", stats.SizeBytes),
		zap.String("size", stats.Size))

	return stats, nil
}

// getSystemStats 获取系统统计信息
func (d *DockerOverviewService) getSystemStats(ctx context.Context) (*response.SystemStats, error) {
	// 获取Docker系统信息
	info, err := global.GVA_DOCKER.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker info: %w", err)
	}

	// 获取Docker版本信息
	version, err := global.GVA_DOCKER.ServerVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker version: %w", err)
	}

	stats := &response.SystemStats{
		Version:         version.Version,
		StorageDriver:   info.Driver,
		CgroupDriver:    info.CgroupDriver,
		KernelVersion:   info.KernelVersion,
		OperatingSystem: info.OperatingSystem,
		Architecture:    info.Architecture,
		CPUs:            info.NCPU,
		MemoryTotal:     info.MemTotal,
	}

	global.GVA_LOG.Debug("System stats collected",
		zap.String("version", stats.Version),
		zap.String("storageDriver", stats.StorageDriver),
		zap.String("cgroupDriver", stats.CgroupDriver),
		zap.Int("cpus", stats.CPUs),
		zap.Int64("memoryTotal", stats.MemoryTotal))

	return stats, nil
}

// GetConfigSummary 获取配置摘要信息
func (d *DockerOverviewService) GetConfigSummary() (*response.ConfigSummary, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取Docker系统信息
	info, err := global.GVA_DOCKER.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker info: %w", err)
	}

	// 获取Docker版本信息
	version, err := global.GVA_DOCKER.ServerVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker version: %w", err)
	}

	// 处理镜像加速器配置
	var registryMirrors []string
	if info.RegistryConfig != nil && info.RegistryConfig.Mirrors != nil {
		registryMirrors = info.RegistryConfig.Mirrors
	}

	summary := &response.ConfigSummary{
		SocketPath:      "unix:///var/run/docker.sock", // 默认路径
		StorageDriver:   info.Driver,
		CgroupDriver:    info.CgroupDriver,
		Version:         version.Version,
		DataRoot:        info.DockerRootDir,
		RegistryMirrors: registryMirrors,
	}

	global.GVA_LOG.Debug("Config summary collected",
		zap.String("version", summary.Version),
		zap.String("storageDriver", summary.StorageDriver),
		zap.String("dataRoot", summary.DataRoot),
		zap.Strings("registryMirrors", summary.RegistryMirrors))

	return summary, nil
}

// getDirectorySize 获取目录大小
func (d *DockerOverviewService) getDirectorySize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			// 忽略权限错误，继续计算其他文件
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

// formatBytes 格式化字节数为可读格式
func (d *DockerOverviewService) formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), sizes[exp+1])
}

// GetDockerDiskUsage 获取Docker磁盘使用情况
func (d *DockerOverviewService) GetDockerDiskUsage() (*response.DiskUsage, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取Docker磁盘使用情况
	diskUsage, err := global.GVA_DOCKER.DiskUsage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage: %w", err)
	}

	// 创建响应结构，处理可能的API版本差异
	usage := &response.DiskUsage{
		LayersSize:     0,
		ImagesSize:     0,
		ContainersSize: 0,
		VolumesSize:    0,
		BuildCacheSize: 0,
	}

	// 尝试获取各种大小信息，如果字段不存在则使用默认值
	if diskUsage.LayersSize > 0 {
		usage.LayersSize = diskUsage.LayersSize
	}

	// 计算镜像总大小
	for _, image := range diskUsage.Images {
		usage.ImagesSize += image.Size
	}

	// 计算容器总大小
	for _, container := range diskUsage.Containers {
		usage.ContainersSize += container.SizeRw
	}

	// 计算存储卷总大小

	for _, volume := range diskUsage.Volumes {
		// Docker API的Volume结构体没有Size字段，需要通过其他方式计算
		// 这里我们尝试通过Mountpoint计算大小
		if volume.Mountpoint != "" {
			size, err := d.getDirectorySize(volume.Mountpoint)
			if err != nil {
				global.GVA_LOG.Debug("Failed to get volume size", zap.String("volume", volume.Name), zap.Error(err))
				continue
			}
			usage.VolumesSize += size
		}
	}

	// 计算总使用量
	usage.TotalSize = usage.LayersSize + usage.ImagesSize + usage.ContainersSize + usage.VolumesSize + usage.BuildCacheSize

	// 格式化大小
	usage.LayersSizeFormatted = d.formatBytes(usage.LayersSize)
	usage.ImagesSizeFormatted = d.formatBytes(usage.ImagesSize)
	usage.ContainersSizeFormatted = d.formatBytes(usage.ContainersSize)
	usage.VolumesSizeFormatted = d.formatBytes(usage.VolumesSize)
	usage.BuildCacheSizeFormatted = d.formatBytes(usage.BuildCacheSize)
	usage.TotalSizeFormatted = d.formatBytes(usage.TotalSize)

	global.GVA_LOG.Debug("Docker disk usage collected",
		zap.Int64("totalSize", usage.TotalSize),
		zap.String("totalSizeFormatted", usage.TotalSizeFormatted))

	return usage, nil
}
