package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// ImageFilter 镜像过滤请求结构
type ImageFilter struct {
	request.PageInfo
	Name      string `json:"name" form:"name"`           // 镜像名称过滤
	Tag       string `json:"tag" form:"tag"`             // 标签过滤
	Dangling  *bool  `json:"dangling" form:"dangling"`   // 是否为悬空镜像
	Reference string `json:"reference" form:"reference"` // 引用过滤
}

// ImagePullRequest 拉取镜像请求
type ImagePullRequest struct {
	Image string `json:"image" binding:"required"` // 镜像名称，如 nginx:latest
	Tag   string `json:"tag"`                      // 标签（可选，如果image中没有包含）
}

// ImageBuildRequest 构建镜像请求
type ImageBuildRequest struct {
	Dockerfile string            `json:"dockerfile" binding:"required"` // Dockerfile内容
	ImageName  string            `json:"imageName" binding:"required"`  // 镜像名称
	Tag        string            `json:"tag"`                           // 标签，默认latest
	BuildArgs  map[string]string `json:"buildArgs"`                     // 构建参数
	Context    string            `json:"context"`                       // 构建上下文路径
}

// ImageTagRequest 镜像标签请求
type ImageTagRequest struct {
	SourceImage string `json:"sourceImage" binding:"required"` // 源镜像
	TargetImage string `json:"targetImage" binding:"required"` // 目标镜像名称
}

// ImageExportRequest 导出镜像请求
type ImageExportRequest struct {
	Images []string `json:"images" binding:"required"` // 要导出的镜像列表
}

// ImageImportRequest 导入镜像请求
type ImageImportRequest struct {
	Source string `json:"source" binding:"required"` // 导入源（文件路径或URL）
	Tag    string `json:"tag"`                       // 标签
}