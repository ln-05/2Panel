package docker

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerContainerApi struct{}

// GetContainerList 获取容器列表
// @Tags Docker
// @Summary 获取Docker容器列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dockerReq.ContainerFilter true "分页参数"
// @Success 200 {object} response.Response{data=dockerRes.ContainerListResponse,msg=string} "获取成功"
// @Router /docker/containers [get]
func (d *DockerContainerApi) GetContainerList(c *gin.Context) {
	var pageInfo dockerReq.ContainerFilter
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置默认分页参数
	if pageInfo.Page <= 0 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize <= 0 {
		pageInfo.PageSize = 10
	}

	// 调用服务层获取容器列表
	containers, total, err := dockerContainerService.GetContainerList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取容器列表失败", zap.Error(err))
		response.FailWithMessage("获取容器列表失败: "+err.Error(), c)
		return
	}

	// 构建响应
	responseData := dockerRes.ContainerListResponse{
		List: containers,
		PageResult: response.PageResult{
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
			Total:    total,
		},
	}

	response.OkWithDetailed(responseData, "获取成功", c)
}

// GetContainerDetail 获取容器详细信息
// @Tags Docker
// @Summary 获取Docker容器详细信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Success 200 {object} response.Response{data=dockerRes.ContainerDetail,msg=string} "获取成功"
// @Router /docker/containers/{id} [get]
func (d *DockerContainerApi) GetContainerDetail(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	// 调用服务层获取容器详细信息
	containerDetail, err := dockerContainerService.GetContainerDetail(containerID)
	if err != nil {
		global.GVA_LOG.Error("获取容器详细信息失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("获取容器详细信息失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(*containerDetail, "获取成功", c)
}

// GetContainerLogs 获取容器日志
// @Tags Docker
// @Summary 获取Docker容器日志
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Param data query dockerReq.LogOptions false "日志选项"
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /docker/containers/{id}/logs [get]
func (d *DockerContainerApi) GetContainerLogs(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	var logOptions dockerReq.LogOptions
	err := c.ShouldBindQuery(&logOptions)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 调用服务层获取容器日志
	logs, err := dockerContainerService.GetContainerLogs(containerID, logOptions)
	if err != nil {
		global.GVA_LOG.Error("获取容器日志失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("获取容器日志失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(logs, "获取成功", c)
}

// StartContainer 启动容器
// @Tags Docker
// @Summary 启动Docker容器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Success 200 {object} response.Response{msg=string} "启动成功"
// @Router /docker/containers/{id}/start [post]
func (d *DockerContainerApi) StartContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	// 调用服务层启动容器
	err := dockerContainerService.StartContainer(containerID)
	if err != nil {
		global.GVA_LOG.Error("启动容器失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("启动容器失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("容器启动成功", c)
}

// StopContainer 停止容器
// @Tags Docker
// @Summary 停止Docker容器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Param timeout query int false "停止超时时间(秒)"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /docker/containers/{id}/stop [post]
func (d *DockerContainerApi) StopContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	// 解析超时参数
	var timeout *int
	if timeoutStr := c.Query("timeout"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = &t
		}
	}

	// 调用服务层停止容器
	err := dockerContainerService.StopContainer(containerID, timeout)
	if err != nil {
		global.GVA_LOG.Error("停止容器失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("停止容器失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("容器停止成功", c)
}

// RestartContainer 重启容器
// @Tags Docker
// @Summary 重启Docker容器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Param timeout query int false "重启超时时间(秒)"
// @Success 200 {object} response.Response{msg=string} "重启成功"
// @Router /docker/containers/{id}/restart [post]
func (d *DockerContainerApi) RestartContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	// 解析超时参数
	var timeout *int
	if timeoutStr := c.Query("timeout"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = &t
		}
	}

	// 调用服务层重启容器
	err := dockerContainerService.RestartContainer(containerID, timeout)
	if err != nil {
		global.GVA_LOG.Error("重启容器失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("重启容器失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("容器重启成功", c)
}

// RemoveContainer 删除容器
// @Tags Docker
// @Summary 删除Docker容器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path string true "容器ID"
// @Param force query bool false "是否强制删除"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /docker/containers/{id} [delete]
func (d *DockerContainerApi) RemoveContainer(c *gin.Context) {
	containerID := c.Param("id")
	if containerID == "" {
		response.FailWithMessage("容器ID不能为空", c)
		return
	}

	// 解析强制删除参数
	force := c.Query("force") == "true"

	// 调用服务层删除容器
	err := dockerContainerService.RemoveContainer(containerID, force)
	if err != nil {
		global.GVA_LOG.Error("删除容器失败", zap.String("containerID", containerID), zap.Error(err))
		if err.Error() == "container not found" {
			response.FailWithMessage("容器不存在", c)
			return
		}
		response.FailWithMessage("删除容器失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("容器删除成功", c)
}

// GetDockerInfo 获取Docker系统信息
// @Tags Docker
// @Summary 获取Docker系统信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=types.Info,msg=string} "获取成功"
// @Router /docker/info [get]
func (d *DockerContainerApi) GetDockerInfo(c *gin.Context) {
	// 调用服务层获取Docker信息
	info, err := dockerContainerService.GetDockerInfo()
	if err != nil {
		global.GVA_LOG.Error("获取Docker信息失败", zap.Error(err))
		response.FailWithMessage("获取Docker信息失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(info, "获取成功", c)
}

// CheckDockerStatus 检查Docker状态
// @Tags Docker
// @Summary 检查Docker守护进程状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=bool,msg=string} "检查成功"
// @Router /docker/status [get]
func (d *DockerContainerApi) CheckDockerStatus(c *gin.Context) {
	// 调用服务层检查Docker状态
	isAvailable := dockerContainerService.IsDockerAvailable()
	
	if isAvailable {
		response.OkWithDetailed(true, "Docker守护进程运行正常", c)
	} else {
		// 提供更详细的错误信息
		response.OkWithDetailed(false, "Docker守护进程不可用，请检查Docker Desktop是否启动", c)
	}
}