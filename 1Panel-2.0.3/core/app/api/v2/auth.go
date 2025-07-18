package v2

import (
	"encoding/base64"
	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/utils/captcha"
	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

// @Tags Auth
// @Summary User login
// @Accept json
// @Param EntranceCode header string true "安全入口 base64 加密串"
// @Param request body dto.Login true "request"
// @Success 200 {object} dto.UserLoginInfo
// @Router /core/auth/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var req dto.Login
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if !req.IgnoreCaptcha {
		if errMsg := captcha.VerifyCode(req.CaptchaID, req.Captcha); errMsg != "" {
			helper.BadAuth(c, errMsg, nil)
			return
		}
	}
	entranceItem := c.Request.Header.Get("EntranceCode")
	var entrance []byte
	if len(entranceItem) != 0 {
		entrance, _ = base64.StdEncoding.DecodeString(entranceItem)
	}
	if len(entrance) == 0 {
		cookieValue, err := c.Cookie("SecurityEntrance")
		if err == nil {
			entrance, _ = base64.StdEncoding.DecodeString(cookieValue)
		}
	}

	user, msgKey, err := authService.Login(c, req, string(entrance))
	go saveLoginLogs(c, err)
	if msgKey == "ErrAuth" || msgKey == "ErrEntrance" {
		helper.BadAuth(c, msgKey, err)
		return
	}
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, user)
}

// @Tags Auth
// @Summary User login with mfa
// @Accept json
// @Param request body dto.MFALogin true "request"
// @Success 200 {object} dto.UserLoginInfo
// @Router /core/auth/mfalogin [post]
// @Header 200 {string} EntranceCode "安全入口"
func (b *BaseApi) MFALogin(c *gin.Context) {
	var req dto.MFALogin
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	entranceItem := c.Request.Header.Get("EntranceCode")
	var entrance []byte
	if len(entranceItem) != 0 {
		entrance, _ = base64.StdEncoding.DecodeString(entranceItem)
	}

	user, msgKey, err := authService.MFALogin(c, req, string(entrance))
	if msgKey == "ErrAuth" {
		helper.BadAuth(c, msgKey, err)
		return
	}
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, user)
}

// @Tags Auth
// @Summary User logout
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /core/auth/logout [post]
func (b *BaseApi) LogOut(c *gin.Context) {
	if err := authService.LogOut(c); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Auth
// @Summary Load captcha
// @Success 200 {object} dto.CaptchaResponse
// @Router /core/auth/captcha [get]
func (b *BaseApi) Captcha(c *gin.Context) {
	captcha, err := captcha.CreateCaptcha()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, captcha)
}

func (b *BaseApi) GetResponsePage(c *gin.Context) {
	pageCode, err := authService.GetResponsePage()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, pageCode)
}

// @Tags Auth
// @Summary Get Setting For Login
// @Success 200 {object} dto.SystemSetting
// @Router /core/auth/setting [get]
func (b *BaseApi) GetLoginSetting(c *gin.Context) {
	settingInfo, err := settingService.GetSettingInfo()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	res := &dto.LoginSetting{
		IsDemo:    global.CONF.Base.IsDemo,
		IsIntl:    global.CONF.Base.IsIntl,
		Language:  settingInfo.Language,
		MenuTabs:  settingInfo.MenuTabs,
		PanelName: settingInfo.PanelName,
		Theme:     settingInfo.Theme,
	}
	helper.SuccessWithData(c, res)
}

func saveLoginLogs(c *gin.Context, err error) {
	var logs model.LoginLog
	if err != nil {
		logs.Status = constant.StatusFailed
		logs.Message = err.Error()
	} else {
		logs.Status = constant.StatusSuccess
	}
	logs.IP = c.ClientIP()
	logs.Agent = c.GetHeader("User-Agent")
	_ = logService.CreateLoginLog(logs)
}
