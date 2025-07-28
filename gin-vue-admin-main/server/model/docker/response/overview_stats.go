package response

// OverviewStats Docker概览统计信息
type OverviewStats struct {
	Containers ContainerStats `json:"containers"`
	Images     ImageStats     `json:"images"`
	Networks   NetworkStats   `json:"networks"`
	Volumes    VolumeStats    `json:"volumes"`
	System     SystemStats    `json:"system"`
}

// ContainerStats 容器统计信息
type ContainerStats struct {
	Total   int `json:"total"`   // 总数
	Running int `json:"running"` // 运行中
	Stopped int `json:"stopped"` // 已停止
	Paused  int `json:"paused"`  // 暂停
}

// ImageStats 镜像统计信息
type ImageStats struct {
	Total     int    `json:"total"`     // 总数
	Size      string `json:"size"`      // 总大小（格式化）
	SizeBytes int64  `json:"sizeBytes"` // 总大小（字节）
}

// NetworkStats 网络统计信息
type NetworkStats struct {
	Total  int `json:"total"`  // 总数
	Bridge int `json:"bridge"` // bridge网络数量
	Host   int `json:"host"`   // host网络数量
	None   int `json:"none"`   // none网络数量
}

// VolumeStats 存储卷统计信息
type VolumeStats struct {
	Total     int    `json:"total"`     // 总数
	Size      string `json:"size"`      // 总大小（格式化）
	SizeBytes int64  `json:"sizeBytes"` // 总大小（字节）
}

// SystemStats 系统统计信息
type SystemStats struct {
	Version         string `json:"version"`         // Docker版本
	StorageDriver   string `json:"storageDriver"`   // 存储驱动
	CgroupDriver    string `json:"cgroupDriver"`    // Cgroup驱动
	KernelVersion   string `json:"kernelVersion"`   // 内核版本
	OperatingSystem string `json:"operatingSystem"` // 操作系统
	Architecture    string `json:"architecture"`    // 架构
	CPUs            int    `json:"cpus"`            // CPU数量
	MemoryTotal     int64  `json:"memoryTotal"`     // 总内存
}

// ConfigSummary 配置摘要信息
type ConfigSummary struct {
	SocketPath      string   `json:"socketPath"`      // Socket路径
	RegistryMirrors []string `json:"registryMirrors"` // 镜像加速器
	StorageDriver   string   `json:"storageDriver"`   // 存储驱动
	CgroupDriver    string   `json:"cgroupDriver"`    // Cgroup驱动
	Version         string   `json:"version"`         // Docker版本
	DataRoot        string   `json:"dataRoot"`        // 数据根目录
}

// DiskUsage Docker磁盘使用情况
type DiskUsage struct {
	LayersSize    int64  `json:"layersSize"`    // 镜像层大小
	ImagesSize    int64  `json:"imagesSize"`    // 镜像大小
	ContainersSize int64  `json:"containersSize"` // 容器大小
	VolumesSize   int64  `json:"volumesSize"`   // 存储卷大小
	BuildCacheSize int64  `json:"buildCacheSize"` // 构建缓存大小
	TotalSize     int64  `json:"totalSize"`     // 总大小
	
	// 格式化后的大小
	LayersSizeFormatted    string `json:"layersSizeFormatted"`
	ImagesSizeFormatted    string `json:"imagesSizeFormatted"`
	ContainersSizeFormatted string `json:"containersSizeFormatted"`
	VolumesSizeFormatted   string `json:"volumesSizeFormatted"`
	BuildCacheSizeFormatted string `json:"buildCacheSizeFormatted"`
	TotalSizeFormatted     string `json:"totalSizeFormatted"`
}