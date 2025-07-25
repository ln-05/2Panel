package docker

import (
	"database/sql"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"go.uber.org/zap"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type RegistrySyncService struct{}

// SyncFrom1Panel 从1Panel同步仓库数据
func (r *RegistrySyncService) SyncFrom1Panel() error {
	// 1Panel通常使用SQLite数据库，路径可能在 /opt/1panel/db/1Panel.db
	// 或者使用MySQL，需要根据实际情况调整
	
	// 尝试连接1Panel数据库
	db, err := r.connect1PanelDB()
	if err != nil {
		global.GVA_LOG.Error("连接1Panel数据库失败", zap.Error(err))
		return err
	}
	defer db.Close()

	// 查询1Panel中的仓库数据
	registries, err := r.query1PanelRegistries(db)
	if err != nil {
		global.GVA_LOG.Error("查询1Panel仓库数据失败", zap.Error(err))
		return err
	}

	// 同步到我们的数据库
	for _, registry := range registries {
		err := r.syncRegistry(registry)
		if err != nil {
			global.GVA_LOG.Error("同步仓库失败", zap.String("name", registry.Name), zap.Error(err))
		} else {
			global.GVA_LOG.Info("同步仓库成功", zap.String("name", registry.Name))
		}
	}

	return nil
}

// connect1PanelDB 连接1Panel数据库
func (r *RegistrySyncService) connect1PanelDB() (*sql.DB, error) {
	// 1Panel可能使用SQLite或MySQL
	// 这里需要根据实际的1Panel配置来调整
	
	// 尝试SQLite连接（1Panel默认使用SQLite）
	dbPath := "/opt/1panel/db/1Panel.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err == nil {
		if err = db.Ping(); err == nil {
			return db, nil
		}
	}

	// 如果SQLite连接失败，尝试MySQL连接
	// 这里使用当前系统的MySQL配置
	mysqlConfig := global.GVA_CONFIG.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Path,
		mysqlConfig.Port,
		"1panel", // 假设1Panel使用名为1panel的数据库
		mysqlConfig.Config,
	)
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	
	if err = db.Ping(); err != nil {
		return nil, err
	}
	
	return db, nil
}

// query1PanelRegistries 查询1Panel中的仓库数据
func (r *RegistrySyncService) query1PanelRegistries(db *sql.DB) ([]dockerModel.DockerRegistry, error) {
	// 1Panel中仓库数据的表名可能是 docker_registries 或类似的名称
	// 这里需要根据实际的1Panel表结构来调整
	
	queries := []string{
		// 尝试不同可能的表名和字段名
		"SELECT name, download_url, protocol, username, password, description FROM docker_registries",
		"SELECT name, url, protocol, username, password, description FROM registries", 
		"SELECT name, registry_url, protocol, username, password, description FROM docker_registry",
		"SELECT name, address, protocol, username, password, description FROM container_registries",
	}

	var registries []dockerModel.DockerRegistry
	
	for _, query := range queries {
		rows, err := db.Query(query)
		if err != nil {
			// 如果查询失败，尝试下一个查询
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var registry dockerModel.DockerRegistry
			var username, password, description sql.NullString
			
			err := rows.Scan(
				&registry.Name,
				&registry.DownloadUrl,
				&registry.Protocol,
				&username,
				&password,
				&description,
			)
			if err != nil {
				global.GVA_LOG.Error("扫描仓库数据失败", zap.Error(err))
				continue
			}

			// 处理可能为NULL的字段
			if username.Valid {
				registry.Username = username.String
			}
			if password.Valid {
				registry.Password = password.String
			}
			if description.Valid {
				registry.Description = description.String
			}

			registry.Status = "active"
			registries = append(registries, registry)
		}

		// 如果成功查询到数据，就不再尝试其他查询
		if len(registries) > 0 {
			break
		}
	}

	return registries, nil
}

// syncRegistry 同步单个仓库到我们的数据库
func (r *RegistrySyncService) syncRegistry(registry dockerModel.DockerRegistry) error {
	// 检查是否已存在同名仓库
	var existingRegistry dockerModel.DockerRegistry
	err := global.GVA_DB.Where("name = ?", registry.Name).First(&existingRegistry).Error
	
	if err != nil {
		// 不存在，创建新仓库
		return global.GVA_DB.Create(&registry).Error
	} else {
		// 已存在，更新仓库信息
		existingRegistry.DownloadUrl = registry.DownloadUrl
		existingRegistry.Protocol = registry.Protocol
		existingRegistry.Username = registry.Username
		existingRegistry.Password = registry.Password
		existingRegistry.Description = registry.Description
		existingRegistry.Status = registry.Status
		
		return global.GVA_DB.Save(&existingRegistry).Error
	}
}