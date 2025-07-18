package v2

import (
	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/gin-gonic/gin"
)

// @Tags Website Acme
// @Summary Page website acme accounts
// @Accept json
// @Param request body dto.PageInfo true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /websites/acme/search [post]
func (b *BaseApi) PageWebsiteAcmeAccount(c *gin.Context) {
	var req dto.PageInfo
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	total, accounts, err := websiteAcmeAccountService.Page(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, dto.PageResult{
		Total: total,
		Items: accounts,
	})
}

// @Tags Website Acme
// @Summary Create website acme account
// @Accept json
// @Param request body request.WebsiteAcmeAccountCreate true "request"
// @Success 200 {object} response.WebsiteAcmeAccountDTO
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /websites/acme [post]
// @x-panel-log {"bodyKeys":["email"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建网站 acme [email]","formatEN":"Create website acme [email]"}
func (b *BaseApi) CreateWebsiteAcmeAccount(c *gin.Context) {
	var req request.WebsiteAcmeAccountCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	res, err := websiteAcmeAccountService.Create(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, res)
}

// @Tags Website Acme
// @Summary Delete website acme account
// @Accept json
// @Param request body request.WebsiteResourceReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /websites/acme/del [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"website_acme_accounts","output_column":"email","output_value":"email"}],"formatZH":"删除网站 acme [email]","formatEN":"Delete website acme [email]"}
func (b *BaseApi) DeleteWebsiteAcmeAccount(c *gin.Context) {
	var req request.WebsiteResourceReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := websiteAcmeAccountService.Delete(req.ID); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

// @Tags Website Acme
// @Summary Update website acme account
// @Accept json
// @Param request body request.WebsiteAcmeAccountUpdate true "request"
// @Success 200 {object} response.WebsiteAcmeAccountDTO
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /websites/acme/update [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"website_acme_accounts","output_column":"email","output_value":"email"}],"formatZH":"更新 acme [email]","formatEN":"Update acme [email]"}
func (b *BaseApi) UpdateWebsiteAcmeAccount(c *gin.Context) {
	var req request.WebsiteAcmeAccountUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	res, err := websiteAcmeAccountService.Update(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, res)
}
