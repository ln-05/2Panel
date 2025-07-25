package response

// NetworkInfo 网络基本信息
type NetworkInfo struct {
	ID         string            `json:"id"`         // 网络ID
	Name       string            `json:"name"`       // 网络名称
	Driver     string            `json:"driver"`     // 网络驱动
	Scope      string            `json:"scope"`      // 网络范围
	EnableIPv6 bool              `json:"enableIPv6"` // 是否启用IPv6
	Internal   bool              `json:"internal"`   // 是否为内部网络
	Attachable bool              `json:"attachable"` // 是否可附加
	Created    int64             `json:"created"`    // 创建时间
	Labels     map[string]string `json:"labels"`     // 标签
	IPAM       *IPAMConfig       `json:"ipam"`       // IPAM配置
}

// NetworkDetail 网络详细信息
type NetworkDetail struct {
	NetworkInfo
	Containers map[string]NetworkContainer `json:"containers"` // 连接的容器
}

// IPAMConfig IPAM配置
type IPAMConfig struct {
	Subnet  string `json:"subnet"`  // 子网
	Gateway string `json:"gateway"` // 网关
}

// NetworkContainer 网络中的容器信息
type NetworkContainer struct {
	Name        string `json:"name"`        // 容器名称
	EndpointID  string `json:"endpointId"`  // 端点ID
	MacAddress  string `json:"macAddress"`  // MAC地址
	IPv4Address string `json:"ipv4Address"` // IPv4地址
	IPv6Address string `json:"ipv6Address"` // IPv6地址
}

// NetworkListResponse 网络列表响应
type NetworkListResponse struct {
	List  []NetworkInfo `json:"list"`  // 网络列表
	Total int64         `json:"total"` // 总数
}