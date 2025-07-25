package cion

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ImageRouter struct {}

// InitImageRouter 初始化 image表 路由信息
func (s *ImageRouter) InitImageRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	imageRouter := Router.Group("image").Use(middleware.OperationRecord())
	imageRouterWithoutRecord := Router.Group("image")
	imageRouterWithoutAuth := PublicRouter.Group("image")
	{
		imageRouter.POST("createImage", imageApi.CreateImage)   // 新建image表
		imageRouter.DELETE("deleteImage", imageApi.DeleteImage) // 删除image表
		imageRouter.DELETE("deleteImageByIds", imageApi.DeleteImageByIds) // 批量删除image表
		imageRouter.PUT("updateImage", imageApi.UpdateImage)    // 更新image表
	}
	{
		imageRouterWithoutRecord.GET("findImage", imageApi.FindImage)        // 根据ID获取image表
		imageRouterWithoutRecord.GET("getImageList", imageApi.GetImageList)  // 获取image表列表
	}
	{
	    imageRouterWithoutAuth.GET("getImagePublic", imageApi.GetImagePublic)  // image表开放接口
	}
}
