package response

import "time"

// RegistryInfo 仓库基本信息
type RegistryInfo struct {
	ID          uint      `json:"id"`          // 仓库ID
	Name        string    `json:"name"`        // 仓库名称
	DownloadUrl string    `json:"downloadUrl"` // 下载地址
	Protocol    string    `json:"protocol"`    // 协议 (http/https)
	Status      string    `json:"status"`      // 状态 (active/inactive)
	Username    string    `json:"username"`    // 用户名
	Description string    `json:"description"` // 描述
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`   // 更新时间
}

// RegistryDetail 仓库详细信息
type RegistryDetail struct {
	RegistryInfo
	IsDefault    bool       `json:"isDefault"`    // 是否为默认仓库
	LastTestTime *time.Time `json:"lastTestTime"` // 最后测试时间
	TestResult   string     `json:"testResult"`   // 测试结果
}

// RegistryListResponse 仓库列表响应
type RegistryListResponse struct {
	List  []RegistryInfo `json:"list"`  // 仓库列表
	Total int64          `json:"total"` // 总数
}

// RegistryTestResponse 仓库测试响应
type RegistryTestResponse struct {
	Success bool   `json:"success"` // 测试是否成功
	Message string `json:"message"` // 测试结果消息
}