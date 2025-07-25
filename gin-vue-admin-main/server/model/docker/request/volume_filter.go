package request

// VolumeFilter 存储卷过滤器
type VolumeFilter struct {
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小
	Name     string `form:"name" json:"name"`         // 存储卷名称过滤
	Driver   string `form:"driver" json:"driver"`     // 驱动类型过滤
}

// VolumeCreateRequest 创建存储卷请求
type VolumeCreateRequest struct {
	Name       string            `json:"name" binding:"required"`       // 存储卷名称
	Driver     string            `json:"driver"`                        // 存储卷驱动，默认为local
	DriverOpts map[string]string `json:"driverOpts"`                    // 驱动选项
	Labels     map[string]string `json:"labels"`                        // 标签
}