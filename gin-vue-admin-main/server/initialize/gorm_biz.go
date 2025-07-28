package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cion"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"gorm.io/gorm"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(cion.Image{})
	if err != nil {
		return err
	}
	
	// Docker模型迁移 - 移除DockerRegistry，因为我们直接从1Panel读取
	err = db.AutoMigrate(
		&docker.DockerOrchestration{},
		&docker.DockerOrchestrationService{},
	)
	if err != nil {
		return err
	}
	
	// 清理可能存在的docker_registries表数据
	cleanupDockerRegistriesTable(db)
	
	return nil
}

// cleanupDockerRegistriesTable 清理docker_registries表
func cleanupDockerRegistriesTable(db *gorm.DB) {
	// 检查表是否存在
	if db.Migrator().HasTable("docker_registries") {
		// 删除表中的所有数据
		db.Exec("DELETE FROM docker_registries")
		global.GVA_LOG.Info("已清理docker_registries表中的数据")
		
		// 可选：完全删除表
		// db.Migrator().DropTable("docker_registries")
		// global.GVA_LOG.Info("已删除docker_registries表")
	}
}
