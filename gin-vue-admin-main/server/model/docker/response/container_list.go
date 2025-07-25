package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"

// ContainerListResponse 容器列表响应
type ContainerListResponse struct {
	response.PageResult
	List []ContainerInfo `json:"list"`
}

// ContainerDetailResponse 容器详情响应
type ContainerDetailResponse struct {
	response.Response
	Data ContainerDetail `json:"data"`
}

// ContainerLogsResponse 容器日志响应
type ContainerLogsResponse struct {
	response.Response
	Data string `json:"data"` // 日志内容
}