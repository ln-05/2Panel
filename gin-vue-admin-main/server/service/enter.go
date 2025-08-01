package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/cion"
	"github.com/flipped-aurora/gin-vue-admin/server/service/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/my"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	CionServiceGroup    cion.ServiceGroup
	MyServiceGroup      my.ServiceGroup
	DockerServiceGroup  docker.ServiceGroup
}
