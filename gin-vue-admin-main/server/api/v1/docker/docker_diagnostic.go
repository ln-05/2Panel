package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerDiagnosticApi struct{}

// DiagnoseConnection 诊断Docker连接
// @Tags DockerDiagnostic
// @Summary 诊断Docker连接
// @Description 全面诊断Docker连接状态，包括网络、权限、配置等
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=docker.DiagnosticResult} "诊断成功"
// @Router /docker/diagnose [get]
func (d *DockerDiagnosticApi) DiagnoseConnection(c *gin.Context) {
	result, err := dockerDiagnosticService.DiagnoseConnection()
	if err != nil {
		global.GVA_LOG.Error("Docker连接诊断失败", zap.Error(err))
		response.FailWithMessage("Docker连接诊断失败: "+err.Error(), c)
		return
	}

	response.OkWithData(result, c)
}

// TestConnection 快速连接测试
// @Tags DockerDiagnostic
// @Summary 快速连接测试
// @Description 快速测试Docker连接状态
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=docker.ClientStatusResult} "测试成功"
// @Router /docker/test-connection [post]
func (d *DockerDiagnosticApi) TestConnection(c *gin.Context) {
	result := dockerDiagnosticService.CheckClientStatus()
	
	if result.IsConnected {
		response.OkWithData(result, c)
	} else {
		response.FailWithDetailed(result, "Docker连接测试失败", c)
	}
}

// GetPermissionStatus 获取权限状态
// @Tags DockerDiagnostic
// @Summary 获取权限状态
// @Description 检查Docker API的各项权限状态
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=docker.PermissionCheckResult} "获取成功"
// @Router /docker/permissions [get]
func (d *DockerDiagnosticApi) GetPermissionStatus(c *gin.Context) {
	result := dockerDiagnosticService.CheckPermissions()
	response.OkWithData(result, c)
}