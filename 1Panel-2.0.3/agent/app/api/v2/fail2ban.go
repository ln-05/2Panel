package v2

import (
	"os"

	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/gin-gonic/gin"
)

// @Tags Fail2ban
// @Summary Load fail2ban base info
// @Success 200 {object} dto.Fail2BanBaseInfo
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/base [get]
func (b *BaseApi) LoadFail2BanBaseInfo(c *gin.Context) {
	data, err := fail2banService.LoadBaseInfo()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.SuccessWithData(c, data)
}

// @Tags Fail2ban
// @Summary Page fail2ban ip list
// @Accept json
// @Param request body dto.Fail2BanSearch true "request"
// @Success 200 {Array} string
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/search [post]
func (b *BaseApi) SearchFail2Ban(c *gin.Context) {
	var req dto.Fail2BanSearch
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	list, err := fail2banService.Search(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.SuccessWithData(c, list)
}

// @Tags Fail2ban
// @Summary Operate fail2ban
// @Accept json
// @Param request body dto.Operate true "request"
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/operate [post]
// @x-panel-log {"bodyKeys":["operation"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"[operation] Fail2ban","formatEN":"[operation] Fail2ban"}
func (b *BaseApi) OperateFail2Ban(c *gin.Context) {
	var req dto.Operate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := fail2banService.Operate(req.Operation); err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.Success(c)
}

// @Tags Fail2ban
// @Summary Operate sshd of fail2ban
// @Accept json
// @Param request body dto.Operate true "request"
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/operate/sshd [post]
func (b *BaseApi) OperateSSHD(c *gin.Context) {
	var req dto.Fail2BanSet
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := fail2banService.OperateSSHD(req); err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.Success(c)
}

// @Tags Fail2ban
// @Summary Update fail2ban conf
// @Accept json
// @Param request body dto.Fail2BanUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/update [post]
// @x-panel-log {"bodyKeys":["key","value"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改 Fail2ban 配置 [key] => [value]","formatEN":"update fail2ban conf [key] => [value]"}
func (b *BaseApi) UpdateFail2BanConf(c *gin.Context) {
	var req dto.Fail2BanUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := fail2banService.UpdateConf(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Fail2ban
// @Summary Load fail2ban conf
// @Accept json
// @Success 200 {string} file
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/load/conf [get]
func (b *BaseApi) LoadFail2BanConf(c *gin.Context) {
	path := "/etc/fail2ban/jail.local"
	file, err := os.ReadFile(path)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.SuccessWithData(c, string(file))
}

// @Tags Fail2ban
// @Summary Update fail2ban conf by file
// @Accept json
// @Param request body dto.UpdateByFile true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /toolbox/fail2ban/update/byconf [post]
func (b *BaseApi) UpdateFail2BanConfByFile(c *gin.Context) {
	var req dto.UpdateByFile
	if err := helper.CheckBind(&req, c); err != nil {
		return
	}
	if err := fail2banService.UpdateConfByFile(req); err != nil {
		helper.InternalServer(c, err)
		return
	}

	helper.Success(c)
}
