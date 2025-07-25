package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cion"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(cion.Image{})
	if err != nil {
		return err
	}
	return nil
}
