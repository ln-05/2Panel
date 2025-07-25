package response

// ContainerInfo 容器基本信息响应结构
type ContainerInfo struct {
	ID          string            `json:"id"`          // 容器ID
	Name        string            `json:"name"`        // 容器名称
	Image       string            `json:"image"`       // 镜像名称
	ImageID     string            `json:"imageId"`     // 镜像ID
	Command     string            `json:"command"`     // 启动命令
	Created     int64             `json:"created"`     // 创建时间戳
	Status      string            `json:"status"`      // 状态描述
	State       string            `json:"state"`       // 运行状态
	Ports       []PortMapping     `json:"ports"`       // 端口映射
	Labels      map[string]string `json:"labels"`      // 标签
	SizeRw      int64             `json:"sizeRw,omitempty"`      // 可写层大小
	SizeRootFs  int64             `json:"sizeRootFs,omitempty"`  // 根文件系统大小
}

// PortMapping 端口映射结构
type PortMapping struct {
	PrivatePort int    `json:"privatePort"`           // 容器内部端口
	PublicPort  int    `json:"publicPort,omitempty"`  // 主机端口
	Type        string `json:"type"`                  // 协议类型 (tcp/udp)
	IP          string `json:"ip,omitempty"`          // 绑定IP地址
}

// ContainerDetail 容器详细信息响应结构
type ContainerDetail struct {
	ContainerInfo
	Config          ContainerConfig `json:"config"`          // 容器配置
	HostConfig      HostConfig      `json:"hostConfig"`      // 主机配置
	NetworkSettings NetworkSettings `json:"networkSettings"` // 网络设置
	Mounts          []Mount         `json:"mounts"`          // 挂载点
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Hostname     string            `json:"hostname"`
	Domainname   string            `json:"domainname"`
	User         string            `json:"user"`
	Env          []string          `json:"env"`          // 环境变量
	Cmd          []string          `json:"cmd"`          // 命令
	Image        string            `json:"image"`        // 镜像
	WorkingDir   string            `json:"workingDir"`   // 工作目录
	Entrypoint   []string          `json:"entrypoint"`   // 入口点
	Labels       map[string]string `json:"labels"`       // 标签
	ExposedPorts map[string]struct{} `json:"exposedPorts"` // 暴露端口
}

// HostConfig 主机配置
type HostConfig struct {
	Binds           []string          `json:"binds"`           // 绑定挂载
	ContainerIDFile string            `json:"containerIdFile"` // 容器ID文件
	LogConfig       LogConfig         `json:"logConfig"`       // 日志配置
	NetworkMode     string            `json:"networkMode"`     // 网络模式
	PortBindings    map[string][]PortBinding `json:"portBindings"` // 端口绑定
	RestartPolicy   RestartPolicy     `json:"restartPolicy"`   // 重启策略
	AutoRemove      bool              `json:"autoRemove"`      // 自动删除
	VolumeDriver    string            `json:"volumeDriver"`    // 卷驱动
	VolumesFrom     []string          `json:"volumesFrom"`     // 从其他容器挂载卷
	CapAdd          []string          `json:"capAdd"`          // 添加的能力
	CapDrop         []string          `json:"capDrop"`         // 删除的能力
	DNS             []string          `json:"dns"`             // DNS服务器
	DNSOptions      []string          `json:"dnsOptions"`      // DNS选项
	DNSSearch       []string          `json:"dnsSearch"`       // DNS搜索域
	ExtraHosts      []string          `json:"extraHosts"`      // 额外主机
	GroupAdd        []string          `json:"groupAdd"`        // 添加的组
	IpcMode         string            `json:"ipcMode"`         // IPC模式
	Cgroup          string            `json:"cgroup"`          // Cgroup
	Links           []string          `json:"links"`           // 链接
	OomScoreAdj     int               `json:"oomScoreAdj"`     // OOM分数调整
	PidMode         string            `json:"pidMode"`         // PID模式
	Privileged      bool              `json:"privileged"`      // 特权模式
	PublishAllPorts bool              `json:"publishAllPorts"` // 发布所有端口
	ReadonlyRootfs  bool              `json:"readonlyRootfs"`  // 只读根文件系统
	SecurityOpt     []string          `json:"securityOpt"`     // 安全选项
	UTSMode         string            `json:"utsMode"`         // UTS模式
	UsernsMode      string            `json:"usernsMode"`      // 用户命名空间模式
	ShmSize         int64             `json:"shmSize"`         // 共享内存大小
	Sysctls         map[string]string `json:"sysctls"`         // 系统控制参数
	Runtime         string            `json:"runtime"`         // 运行时
}

// LogConfig 日志配置
type LogConfig struct {
	Type   string            `json:"type"`   // 日志类型
	Config map[string]string `json:"config"` // 日志配置
}

// PortBinding 端口绑定
type PortBinding struct {
	HostIP   string `json:"hostIp"`   // 主机IP
	HostPort string `json:"hostPort"` // 主机端口
}

// RestartPolicy 重启策略
type RestartPolicy struct {
	Name              string `json:"name"`              // 策略名称
	MaximumRetryCount int    `json:"maximumRetryCount"` // 最大重试次数
}

// NetworkSettings 网络设置
type NetworkSettings struct {
	Bridge                 string                 `json:"bridge"`
	SandboxID              string                 `json:"sandboxId"`
	HairpinMode            bool                   `json:"hairpinMode"`
	LinkLocalIPv6Address   string                 `json:"linkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int                    `json:"linkLocalIPv6PrefixLen"`
	Ports                  map[string][]PortBinding `json:"ports"`
	SandboxKey             string                 `json:"sandboxKey"`
	SecondaryIPAddresses   []string               `json:"secondaryIPAddresses"`
	SecondaryIPv6Addresses []string               `json:"secondaryIPv6Addresses"`
	EndpointID             string                 `json:"endpointId"`
	Gateway                string                 `json:"gateway"`
	GlobalIPv6Address      string                 `json:"globalIPv6Address"`
	GlobalIPv6PrefixLen    int                    `json:"globalIPv6PrefixLen"`
	IPAddress              string                 `json:"ipAddress"`
	IPPrefixLen            int                    `json:"ipPrefixLen"`
	IPv6Gateway            string                 `json:"ipv6Gateway"`
	MacAddress             string                 `json:"macAddress"`
	Networks               map[string]EndpointSettings `json:"networks"`
}

// EndpointSettings 端点设置
type EndpointSettings struct {
	IPAMConfig          *EndpointIPAMConfig `json:"ipamConfig"`
	Links               []string            `json:"links"`
	Aliases             []string            `json:"aliases"`
	NetworkID           string              `json:"networkId"`
	EndpointID          string              `json:"endpointId"`
	Gateway             string              `json:"gateway"`
	IPAddress           string              `json:"ipAddress"`
	IPPrefixLen         int                 `json:"ipPrefixLen"`
	IPv6Gateway         string              `json:"ipv6Gateway"`
	GlobalIPv6Address   string              `json:"globalIPv6Address"`
	GlobalIPv6PrefixLen int                 `json:"globalIPv6PrefixLen"`
	MacAddress          string              `json:"macAddress"`
	DriverOpts          map[string]string   `json:"driverOpts"`
}

// EndpointIPAMConfig IPAM配置
type EndpointIPAMConfig struct {
	IPv4Address string `json:"ipv4Address"`
	IPv6Address string `json:"ipv6Address"`
}

// Mount 挂载点
type Mount struct {
	Target        string      `json:"target"`        // 目标路径
	Source        string      `json:"source"`        // 源路径
	Type          string      `json:"type"`          // 挂载类型
	ReadOnly      bool        `json:"readOnly"`      // 是否只读
	Consistency   string      `json:"consistency"`   // 一致性
	BindOptions   *BindOptions `json:"bindOptions"`   // 绑定选项
	VolumeOptions *VolumeOptions `json:"volumeOptions"` // 卷选项
	TmpfsOptions  *TmpfsOptions `json:"tmpfsOptions"`  // Tmpfs选项
}

// BindOptions 绑定选项
type BindOptions struct {
	Propagation string `json:"propagation"` // 传播模式
}

// VolumeOptions 卷选项
type VolumeOptions struct {
	NoCopy       bool              `json:"noCopy"`       // 不复制
	Labels       map[string]string `json:"labels"`       // 标签
	DriverConfig *VolumeDriverConfig `json:"driverConfig"` // 驱动配置
}

// VolumeDriverConfig 卷驱动配置
type VolumeDriverConfig struct {
	Name    string            `json:"name"`    // 驱动名称
	Options map[string]string `json:"options"` // 驱动选项
}

// TmpfsOptions Tmpfs选项
type TmpfsOptions struct {
	SizeBytes int64 `json:"sizeBytes"` // 大小（字节）
	Mode      int   `json:"mode"`      // 模式
}