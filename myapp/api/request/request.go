package request

type DatabaseCreateRequest struct {
	Name        string `json:"name" binding:"required"`     // 数据库连接名称（必填）
	Type        string `json:"type" binding:"required"`     // 数据库类型（必填）
	Host        string `json:"host" binding:"required"`     // 主机地址（必填）
	Port        int32  `json:"port"`                        // 端口号
	Username    string `json:"username" binding:"required"` // 用户名（必填）
	Password    string `json:"password" binding:"required"` // 密码（必填）
	Database    string `json:"database"`                    // 数据库名
	Charset     string `json:"charset"`                     // 字符集
	Description string `json:"description"`                 // 描述信息
	Config      string `json:"config"`                      // 额外配置信息（JSON格式）
	CreatedBy   uint64 `json:"created_by"`                  // 创建者ID
}

type DatabaseUpdateRequest struct {
	ID          uint64 `json:"id" binding:"required"` // 数据库ID（必填）
	Name        string `json:"name"`                  // 数据库连接名称
	Type        string `json:"type"`                  // 数据库类型
	Host        string `json:"host"`                  // 主机地址
	Port        int32  `json:"port"`                  // 端口号
	Username    string `json:"username"`              // 用户名
	Password    string `json:"password"`              // 密码
	Database    string `json:"database"`              // 数据库名
	Charset     string `json:"charset"`               // 字符集
	Description string `json:"description"`           // 描述信息
	Config      string `json:"config"`                // 额外配置信息（JSON格式）
	UpdatedBy   uint64 `json:"updated_by"`            // 更新者ID
}

type DatabaseListRequest struct {
	Page       int32  `json:"page" form:"page"`               // 页码
	PageSize   int32  `json:"page_size" form:"page_size"`     // 每页大小
	Search     string `json:"search" form:"search"`           // 搜索关键词
	TypeFilter string `json:"type_filter" form:"type_filter"` // 类型过滤
}

type DatabaseTestRequest struct {
	DatabaseID uint64 `json:"database_id" binding:"required"` // 数据库ID（必填）
}

type DatabaseDeleteRequest struct {
	ID uint64 `json:"id" binding:"required"` // 数据库ID（必填）
}

type DatabaseGetRequest struct {
	ID uint64 `json:"id" binding:"required"` // 数据库ID（必填）
}

// 数据库同步请求
type DatabaseSyncRequest struct {
	ServerURL string `json:"server_url" binding:"required"` // 远程服务器URL（必填）
	APIKey    string `json:"api_key"`                       // API密钥（可选）
	SyncType  string `json:"sync_type"`                     // 同步类型: full/incremental
	Overwrite bool   `json:"overwrite"`                     // 是否覆盖已存在的记录
}
