package router

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterContainerRouter(rg *gin.RouterGroup) {
	group := rg.Group("/container")
	{
		group.GET("/list", v1.GetContainerList)
		group.POST("/create", v1.CreateContainer)
		group.DELETE("/delete", v1.DeleteContainer)
	}
}
