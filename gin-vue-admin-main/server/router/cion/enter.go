package cion

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ ImageRouter }

var imageApi = api.ApiGroupApp.CionApiGroup.ImageApi
