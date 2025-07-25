package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"

// ImageInfo 镜像基本信息响应结构
type ImageInfo struct {
	ID          string   `json:"id"`          // 镜像ID
	RepoTags    []string `json:"repoTags"`    // 仓库标签
	RepoDigests []string `json:"repoDigests"` // 仓库摘要
	Size        int64    `json:"size"`        // 镜像大小
	VirtualSize int64    `json:"virtualSize"` // 虚拟大小
	Created     int64    `json:"created"`     // 创建时间戳
	Labels      map[string]string `json:"labels"` // 标签
	Containers  int64    `json:"containers"`  // 使用此镜像的容器数量
}

// ImageDetail 镜像详细信息响应结构
type ImageDetail struct {
	ImageInfo
	Architecture string            `json:"architecture"` // 架构
	Os           string            `json:"os"`           // 操作系统
	Config       ImageConfig       `json:"config"`       // 镜像配置
	RootFS       RootFS            `json:"rootFS"`       // 根文件系统
	History      []ImageHistory    `json:"history"`      // 历史记录
}

// ImageConfig 镜像配置
type ImageConfig struct {
	Hostname     string            `json:"hostname"`
	Domainname   string            `json:"domainname"`
	User         string            `json:"user"`
	Env          []string          `json:"env"`          // 环境变量
	Cmd          []string          `json:"cmd"`          // 默认命令
	Image        string            `json:"image"`        // 镜像
	WorkingDir   string            `json:"workingDir"`   // 工作目录
	Entrypoint   []string          `json:"entrypoint"`   // 入口点
	Labels       map[string]string `json:"labels"`       // 标签
	ExposedPorts map[string]struct{} `json:"exposedPorts"` // 暴露端口
}

// RootFS 根文件系统信息
type RootFS struct {
	Type   string   `json:"type"`   // 类型
	Layers []string `json:"layers"` // 层
}

// ImageHistory 镜像历史记录
type ImageHistory struct {
	ID        string `json:"id"`        // ID
	Created   int64  `json:"created"`   // 创建时间
	CreatedBy string `json:"createdBy"` // 创建命令
	Size      int64  `json:"size"`      // 大小
	Comment   string `json:"comment"`   // 注释
}

// ImageListResponse 镜像列表响应
type ImageListResponse struct {
	List []ImageInfo `json:"list"`
	response.PageResult
}

// ImageDetailResponse 镜像详情响应
type ImageDetailResponse struct {
	response.Response
	Data ImageDetail `json:"data"`
}

// ImageBuildResponse 镜像构建响应
type ImageBuildResponse struct {
	response.Response
	Data string `json:"data"` // 构建日志
}