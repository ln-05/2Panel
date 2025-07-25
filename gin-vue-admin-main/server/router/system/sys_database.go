package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type SysDatabaseRouter struct{}

func (s *SysDatabaseRouter) InitDatabaseRouter(Router *gin.RouterGroup) {
	databaseRouter := Router.Group("database")
	{
		// 代理所有数据库请求到 myapp
		databaseRouter.Any("/*path", s.proxyToMyapp)
	}
}

// proxyToMyapp 代理到 myapp 服务
func (s *SysDatabaseRouter) proxyToMyapp(c *gin.Context) {
	// myapp API 服务地址
	targetURL := "http://localhost:8889"

	// 解析目标URL
	target, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "代理配置错误",
		})
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(target)

	// 修改请求路径
	// 原始路径: /api/database/list
	// 目标路径: /api/database/list
	// 获取path参数，例如 /list
	path := c.Param("path")

	// 保存原始路径
	originalPath := c.Request.URL.Path

	// 设置新的路径：/api/database + path
	c.Request.URL.Path = "/api/database" + path

	// 执行代理
	proxy.ServeHTTP(c.Writer, c.Request)

	// 恢复原始路径
	c.Request.URL.Path = originalPath
}
