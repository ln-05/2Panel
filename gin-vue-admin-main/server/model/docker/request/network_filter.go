package request

// NetworkFilter 网络过滤器
type NetworkFilter struct {
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小
	Name     string `form:"name" json:"name"`         // 网络名称过滤
	Driver   string `form:"driver" json:"driver"`     // 驱动类型过滤
}

// NetworkCreateRequest 创建网络请求
type NetworkCreateRequest struct {
	Name       string            `json:"name" binding:"required"`       // 网络名称
	Driver     string            `json:"driver"`                        // 网络驱动，默认为bridge
	Subnet     string            `json:"subnet"`                        // 子网
	Gateway    string            `json:"gateway"`                       // 网关
	EnableIPv6 bool              `json:"enableIPv6"`                    // 是否启用IPv6
	Internal   bool              `json:"internal"`                      // 是否为内部网络
	Attachable bool              `json:"attachable"`                    // 是否可附加
	Labels     map[string]string `json:"labels"`                        // 标签
}