package docker

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	dockerService "github.com/flipped-aurora/gin-vue-admin/server/service/docker"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerConfigApi struct{}

// GetDockerConfig 获取Docker配置
// @Tags Docker
// @Summary 获取Docker配置
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.DockerConfigResponse,msg=string} "获取Docker配置成功"
// @Router /docker/config [get]
func (api *DockerConfigApi) GetDockerConfig(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	config, err := service.GetDockerConfig()
	if err != nil {
		global.GVA_LOG.Error("获取Docker配置失败", zap.Error(err))
		response.FailWithMessage("获取Docker配置失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(config, "获取Docker配置成功", c)
}

// UpdateDockerConfig 更新Docker配置
// @Tags Docker
// @Summary 更新Docker配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerModel.DockerConfigRequest true "Docker配置"
// @Success 200 {object} response.Response{msg=string} "更新Docker配置成功"
// @Router /docker/config [put]
func (api *DockerConfigApi) UpdateDockerConfig(c *gin.Context) {
	var config dockerModel.DockerConfigRequest
	if err := c.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	service := dockerService.NewDockerConfigService()
	if err := service.UpdateDockerConfig(config); err != nil {
		global.GVA_LOG.Error("更新Docker配置失败", zap.Error(err))
		response.FailWithMessage("更新Docker配置失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("更新Docker配置成功", c)
}

// ValidateDockerConfig 验证Docker配置
// @Tags Docker
// @Summary 验证Docker配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerModel.DockerConfigRequest true "Docker配置"
// @Success 200 {object} response.Response{data=dockerModel.ValidationResponse,msg=string} "验证Docker配置成功"
// @Router /docker/config/validate [post]
func (api *DockerConfigApi) ValidateDockerConfig(c *gin.Context) {
	var config dockerModel.DockerConfigRequest
	if err := c.ShouldBindJSON(&config); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	service := dockerService.NewDockerConfigService()
	validationResult, err := service.ValidateConfig(config)
	if err != nil {
		global.GVA_LOG.Error("验证Docker配置失败", zap.Error(err))
		response.FailWithMessage("验证Docker配置失败: "+err.Error(), c)
		return
	}

	if validationResult.Valid {
		response.OkWithDetailed(validationResult, "Docker配置验证通过", c)
	} else {
		response.FailWithDetailed(validationResult, "Docker配置验证失败", c)
	}
}

// BackupDockerConfig 备份Docker配置
// @Tags Docker
// @Summary 备份Docker配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param description query string false "备份描述"
// @Success 200 {object} response.Response{data=dockerModel.BackupResponse,msg=string} "备份Docker配置成功"
// @Router /docker/config/backup [post]
func (api *DockerConfigApi) BackupDockerConfig(c *gin.Context) {
	description := c.DefaultQuery("description", "手动备份")

	service := dockerService.NewDockerConfigService()
	backup, err := service.BackupDockerConfig(description)
	if err != nil {
		global.GVA_LOG.Error("备份Docker配置失败", zap.Error(err))
		response.FailWithMessage("备份Docker配置失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(backup, "备份Docker配置成功", c)
}

// RestoreDockerConfig 恢复Docker配置
// @Tags Docker
// @Summary 恢复Docker配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dockerModel.RestoreRequest true "恢复请求"
// @Success 200 {object} response.Response{msg=string} "恢复Docker配置成功"
// @Router /docker/config/restore [post]
func (api *DockerConfigApi) RestoreDockerConfig(c *gin.Context) {
	var restoreReq dockerModel.RestoreRequest
	if err := c.ShouldBindJSON(&restoreReq); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	service := dockerService.NewDockerConfigService()
	if err := service.RestoreDockerConfig(restoreReq); err != nil {
		global.GVA_LOG.Error("恢复Docker配置失败", zap.Error(err))
		response.FailWithMessage("恢复Docker配置失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("恢复Docker配置成功", c)
}

// GetBackupList 获取备份列表
// @Tags Docker
// @Summary 获取Docker配置备份列表
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.BackupListResponse,msg=string} "获取备份列表成功"
// @Router /docker/config/backups [get]
func (api *DockerConfigApi) GetBackupList(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	backupList, err := service.GetBackupList()
	if err != nil {
		global.GVA_LOG.Error("获取备份列表失败", zap.Error(err))
		response.FailWithMessage("获取备份列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(backupList, "获取备份列表成功", c)
}

// DeleteBackup 删除备份
// @Tags Docker
// @Summary 删除Docker配置备份
// @Security ApiKeyAuth
// @Produce application/json
// @Param backupId path string true "备份ID"
// @Success 200 {object} response.Response{msg=string} "删除备份成功"
// @Router /docker/config/backups/{backupId} [delete]
func (api *DockerConfigApi) DeleteBackup(c *gin.Context) {
	backupId := c.Param("backupId")
	if backupId == "" {
		response.FailWithMessage("备份ID不能为空", c)
		return
	}

	service := dockerService.NewDockerConfigService()
	if err := service.DeleteBackup(backupId); err != nil {
		global.GVA_LOG.Error("删除备份失败", zap.Error(err))
		response.FailWithMessage("删除备份失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除备份成功", c)
}

// CleanupOldBackups 清理过期备份
// @Tags Docker
// @Summary 清理过期Docker配置备份
// @Security ApiKeyAuth
// @Produce application/json
// @Param maxDays query int false "最大保留天数" default(30)
// @Param maxCount query int false "最大保留数量" default(10)
// @Success 200 {object} response.Response{msg=string} "清理过期备份成功"
// @Router /docker/config/backups/cleanup [post]
func (api *DockerConfigApi) CleanupOldBackups(c *gin.Context) {
	maxDaysStr := c.DefaultQuery("maxDays", "30")
	maxCountStr := c.DefaultQuery("maxCount", "10")

	maxDays, err := strconv.Atoi(maxDaysStr)
	if err != nil {
		response.FailWithMessage("maxDays参数无效", c)
		return
	}

	maxCount, err := strconv.Atoi(maxCountStr)
	if err != nil {
		response.FailWithMessage("maxCount参数无效", c)
		return
	}

	maxAge := time.Duration(maxDays) * 24 * time.Hour
	service := dockerService.NewDockerConfigService()
	if err := service.CleanupOldBackups(maxAge, maxCount); err != nil {
		global.GVA_LOG.Error("清理过期备份失败", zap.Error(err))
		response.FailWithMessage("清理过期备份失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("清理过期备份成功", c)
}

// RestartDockerService 重启Docker服务
// @Tags Docker
// @Summary 重启Docker服务
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.ServiceOperationResponse,msg=string} "重启Docker服务成功"
// @Router /docker/service/restart [post]
func (api *DockerConfigApi) RestartDockerService(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	result, err := service.RestartDockerService()
	if err != nil {
		global.GVA_LOG.Error("重启Docker服务失败", zap.Error(err))
		response.FailWithDetailed(result, "重启Docker服务失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "重启Docker服务成功", c)
}

// StartDockerService 启动Docker服务
// @Tags Docker
// @Summary 启动Docker服务
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.ServiceOperationResponse,msg=string} "启动Docker服务成功"
// @Router /docker/service/start [post]
func (api *DockerConfigApi) StartDockerService(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	result, err := service.StartDockerService()
	if err != nil {
		global.GVA_LOG.Error("启动Docker服务失败", zap.Error(err))
		response.FailWithDetailed(result, "启动Docker服务失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "启动Docker服务成功", c)
}

// StopDockerService 停止Docker服务
// @Tags Docker
// @Summary 停止Docker服务
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.ServiceOperationResponse,msg=string} "停止Docker服务成功"
// @Router /docker/service/stop [post]
func (api *DockerConfigApi) StopDockerService(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	result, err := service.StopDockerService()
	if err != nil {
		global.GVA_LOG.Error("停止Docker服务失败", zap.Error(err))
		response.FailWithDetailed(result, "停止Docker服务失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(result, "停止Docker服务成功", c)
}

// GetDockerServiceStatus 获取Docker服务状态
// @Tags Docker
// @Summary 获取Docker服务状态
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=dockerModel.ServiceStatusResponse,msg=string} "获取Docker服务状态成功"
// @Router /docker/service/status [get]
func (api *DockerConfigApi) GetDockerServiceStatus(c *gin.Context) {
	service := dockerService.NewDockerConfigService()
	status, err := service.GetServiceStatus()
	if err != nil {
		global.GVA_LOG.Error("获取Docker服务状态失败", zap.Error(err))
		response.FailWithMessage("获取Docker服务状态失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(status, "获取Docker服务状态成功", c)
}

// CheckDockerServiceHealth 检查Docker服务健康状态
// @Tags Docker
// @Summary 检查Docker服务健康状态
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "Docker服务健康检查成功"
// @Router /docker/service/health [get]
func (api *DockerConfigApi) CheckDockerServiceHealth(c *gin.Context) {
	controller := &dockerService.DockerServiceController{}
	if err := controller.CheckServiceHealth(); err != nil {
		global.GVA_LOG.Error("Docker服务健康检查失败", zap.Error(err))
		response.FailWithMessage("Docker服务健康检查失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("Docker服务运行正常", c)
}