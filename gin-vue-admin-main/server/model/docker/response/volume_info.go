package response

// VolumeInfo 存储卷基本信息
type VolumeInfo struct {
	Name       string            `json:"name"`       // 存储卷名称
	Driver     string            `json:"driver"`     // 存储卷驱动
	Mountpoint string            `json:"mountpoint"` // 挂载点
	Scope      string            `json:"scope"`      // 范围
	CreatedAt  string            `json:"createdAt"`  // 创建时间
	Labels     map[string]string `json:"labels"`     // 标签
	Options    map[string]string `json:"options"`    // 选项
}

// VolumeDetail 存储卷详细信息
type VolumeDetail struct {
	VolumeInfo
	UsageData *VolumeUsageData `json:"usageData"` // 使用情况数据
}

// VolumeUsageData 存储卷使用情况
type VolumeUsageData struct {
	Size     int64 `json:"size"`     // 大小
	RefCount int64 `json:"refCount"` // 引用计数
}

// VolumeListResponse 存储卷列表响应
type VolumeListResponse struct {
	List  []VolumeInfo `json:"list"`  // 存储卷列表
	Total int64        `json:"total"` // 总数
}