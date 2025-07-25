package response

// SysDatabaseTestResponse 数据库连接测试响应
type SysDatabaseTestResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ResponseTime int64  `json:"response_time"`
	Version      string `json:"version,omitempty"`
}

// SysDatabaseSyncResponse 数据库同步响应
type SysDatabaseSyncResponse struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	TotalSynced int      `json:"total_synced"`
	Created     int      `json:"created"`
	Updated     int      `json:"updated"`
	Skipped     int      `json:"skipped"`
	Errors      []string `json:"errors,omitempty"`
}