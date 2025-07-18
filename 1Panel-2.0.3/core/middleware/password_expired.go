package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/core/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/core/app/repo"
	"github.com/1Panel-dev/1Panel/core/constant"
	"github.com/1Panel-dev/1Panel/core/utils/common"
	"github.com/gin-gonic/gin"
)

func PasswordExpired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/v2/core/auth") ||
			c.Request.URL.Path == "/api/v2/core/settings/expired/handle" ||
			c.Request.URL.Path == "/api/v2/core/settings/search" {
			c.Next()
			return
		}
		settingRepo := repo.NewISettingRepo()
		setting, err := settingRepo.Get(repo.WithByKey("ExpirationDays"))
		if err != nil {
			helper.ErrorWithDetail(c, http.StatusInternalServerError, "ErrPasswordExpired", err)
			return
		}
		expiredDays, _ := strconv.Atoi(setting.Value)
		if expiredDays == 0 {
			c.Next()
			return
		}

		extime, err := settingRepo.Get(repo.WithByKey("ExpirationTime"))
		if err != nil {
			helper.ErrorWithDetail(c, http.StatusInternalServerError, "ErrPasswordExpired", err)
			return
		}
		loc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
		expiredTime, err := time.ParseInLocation(constant.DateTimeLayout, extime.Value, loc)
		if err != nil {
			helper.ErrorWithDetail(c, 313, "ErrPasswordExpired", err)
			return
		}
		if time.Now().After(expiredTime) {
			helper.ErrorWithDetail(c, 313, "ErrPasswordExpired", err)
			return
		}
		c.Next()
	}
}
