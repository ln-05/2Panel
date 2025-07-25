package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/cion"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/my"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	CionApiGroup    cion.ApiGroup
	MyApiGroup      my.ApiGroup
	DockerApiGroup  docker.ApiGroup
}
