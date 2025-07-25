package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/cion"
	"github.com/flipped-aurora/gin-vue-admin/server/router/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/my"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Cion    cion.RouterGroup
	My      my.RouterGroup
	Docker  docker.RouterGroup
}
