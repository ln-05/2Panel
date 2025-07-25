package cion

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ ImageApi }

var imageService = service.ServiceGroupApp.CionServiceGroup.ImageService
