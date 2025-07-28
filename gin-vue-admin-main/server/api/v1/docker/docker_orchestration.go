package docker

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

// GetOrchestrationList 获取编排列表（支持分页、搜索、状态过滤）
func (d *DockerContainerApi) GetOrchestrationList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	search := c.Query("name")
	statusFilter := c.Query("status")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 调用服务层获取编排列表
	result, err := dockerContainerService.GetOrchestrationList(page, pageSize, search, statusFilter)
	if err != nil {
		response.FailWithMessage("获取编排列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// OrchestrationDetail 获取单个编排详情
func (d *DockerContainerApi) OrchestrationDetail(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.FailWithMessage("编排名称不能为空", c)
		return
	}

	// 调用服务层获取编排详情
	group, err := dockerContainerService.GetOrchestrationDetail(name)
	if err != nil {
		if err.Error() == "orchestration not found" {
			response.FailWithMessage("未找到该编排", c)
			return
		}
		response.FailWithMessage("获取编排详情失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(group, "获取成功", c)
}

// CreateOrchestration 创建编排（仅组装label，不实际创建容器）
type CreateOrchestrationReq struct {
	Name string `json:"name" binding:"required"`
}

func (d *DockerContainerApi) CreateOrchestration(c *gin.Context) {
	var req CreateOrchestrationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	// 实际上，编排的创建应为创建一组带label的容器，这里仅返回成功
	response.OkWithMessage("编排创建接口（请用容器创建接口并加label实现）", c)
}

// EditOrchestration 编辑编排（仅演示，实际应批量更新label）
type EditOrchestrationReq struct {
	Name    string `json:"name" binding:"required"`
	NewName string `json:"newName" binding:"required"`
}

func (d *DockerContainerApi) EditOrchestration(c *gin.Context) {
	var req EditOrchestrationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}
	// 实际应批量更新所有容器的label
	response.OkWithMessage("编排编辑接口（请用容器编辑接口并批量更新label实现）", c)
}

// DeleteOrchestration 删除编排（批量删除该label下所有容器）
func (d *DockerContainerApi) DeleteOrchestration(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.FailWithMessage("编排名称不能为空", c)
		return
	}

	// 调用服务层删除编排
	failed, err := dockerContainerService.DeleteOrchestration(name)
	if err != nil {
		if len(failed) > 0 {
			response.FailWithDetailed(failed, "部分容器删除失败: "+err.Error(), c)
			return
		}
		response.FailWithMessage("删除编排失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("编排删除成功", c)
}

// OrchestrationStatus 获取编排下容器状态
func (d *DockerContainerApi) OrchestrationStatus(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.FailWithMessage("编排名称不能为空", c)
		return
	}

	// 调用服务层获取编排状态
	status, err := dockerContainerService.GetOrchestrationStatus(name)
	if err != nil {
		if err.Error() == "orchestration not found" {
			response.FailWithMessage("未找到该编排", c)
			return
		}
		response.FailWithMessage("获取编排状态失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(status, "获取成功", c)
}

// BatchOperateOrchestration 批量操作编排下所有容器
// @Tags Docker
// @Summary 批量操作编排下所有容器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param name path string true "编排名称(label)"
// @Param op path string true "操作类型(start/stop/restart/delete)"
// @Param timeout query int false "超时时间(秒)"
// @Param force query bool false "强制删除(仅delete有效)"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "操作结果"
// @Router /docker/orchestrations/{name}/{op} [post]
func (d *DockerContainerApi) BatchOperateOrchestration(c *gin.Context) {
	name := c.Param("name")
	op := c.Param("op")
	if name == "" || op == "" {
		response.FailWithMessage("编排名称和操作类型不能为空", c)
		return
	}
	var timeout *int
	if timeoutStr := c.Query("timeout"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = &t
		}
	}
	force := c.Query("force") == "true"
	successIDs, failed := dockerContainerService.BatchOperateByOrchestrationLabel(name, op, timeout, force)
	resp := map[string]interface{}{
		"successIDs": successIDs,
		"failed":     failed,
	}
	if len(failed) == 0 {
		response.OkWithDetailed(resp, "操作成功", c)
	} else if len(successIDs) == 0 {
		response.FailWithDetailed(resp, "全部失败", c)
	} else {
		response.OkWithDetailed(resp, "部分成功", c)
	}
}
