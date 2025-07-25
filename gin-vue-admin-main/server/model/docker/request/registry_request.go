package request

// RegistryFilter 仓库过滤器
type RegistryFilter struct {
	Page     int    `form:"page" json:"page"`         // 页码
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页大小
	Name     string `form:"name" json:"name"`         // 仓库名称过滤
	Protocol string `form:"protocol" json:"protocol"` // 协议过滤
}

// RegistryCreateRequest 创建仓库请求
type RegistryCreateRequest struct {
	Name        string `json:"name" binding:"required"`        // 仓库名称
	DownloadUrl string `json:"downloadUrl" binding:"required"` // 下载地址
	Protocol    string `json:"protocol" binding:"required"`    // 协议 (http/https)
	Username    string `json:"username"`                       // 用户名（可选）
	Password    string `json:"password"`                       // 密码（可选）
	Description string `json:"description"`                    // 描述（可选）
}

// RegistryUpdateRequest 更新仓库请求
type RegistryUpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`          // 仓库ID
	Name        string `json:"name" binding:"required"`        // 仓库名称
	DownloadUrl string `json:"downloadUrl" binding:"required"` // 下载地址
	Protocol    string `json:"protocol" binding:"required"`    // 协议 (http/https)
	Username    string `json:"username"`                       // 用户名（可选）
	Password    string `json:"password"`                       // 密码（可选）
	Description string `json:"description"`                    // 描述（可选）
}