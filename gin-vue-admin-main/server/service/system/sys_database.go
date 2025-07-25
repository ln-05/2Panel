package system

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type DatabaseService struct{}

// CreateSysDatabase 创建数据库连接
func (databaseService *DatabaseService) CreateSysDatabase(database *system.SysDatabase) (err error) {
	// 检查名称是否重复
	var count int64
	if err = global.GVA_DB.Model(&system.SysDatabase{}).Where("name = ?", database.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("数据库连接名称已存在")
	}
	
	database.Status = "unknown"
	err = global.GVA_DB.Create(database).Error
	return err
}

// DeleteSysDatabase 删除数据库连接
func (databaseService *DatabaseService) DeleteSysDatabase(ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysDatabase{}, "id = ?", ID).Error
	return err
}

// DeleteSysDatabaseByIds 批量删除数据库连接
func (databaseService *DatabaseService) DeleteSysDatabaseByIds(IDs []int) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysDatabase{}, "id in ?", IDs).Error
	return err
}

// UpdateSysDatabase 更新数据库连接
func (databaseService *DatabaseService) UpdateSysDatabase(database system.SysDatabase) (err error) {
	// 检查名称是否重复（排除当前记录）
	var count int64
	if err = global.GVA_DB.Model(&system.SysDatabase{}).Where("name = ? AND id != ?", database.Name, database.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("数据库连接名称已存在")
	}
	
	err = global.GVA_DB.Model(&system.SysDatabase{}).Where("id = ?", database.ID).Updates(&database).Error
	return err
}

// GetSysDatabase 根据ID获取数据库连接
func (databaseService *DatabaseService) GetSysDatabase(ID string) (database system.SysDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&database).Error
	return
}

// GetSysDatabaseInfoList 分页获取数据库连接列表
func (databaseService *DatabaseService) GetSysDatabaseInfoList(info systemReq.SysDatabaseSearch) (list []system.SysDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	
	// 创建db
	db := global.GVA_DB.Model(&system.SysDatabase{})
	var databases []system.SysDatabase
	
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Type != "" {
		db = db.Where("type = ?", info.Type)
	}
	if info.Host != "" {
		db = db.Where("host LIKE ?", "%"+info.Host+"%")
	}
	
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	
	err = db.Find(&databases).Error
	return databases, total, err
}

// TestSysDatabase 测试数据库连接
func (databaseService *DatabaseService) TestSysDatabase(ID string) (response systemRes.SysDatabaseTestResponse, err error) {
	var database system.SysDatabase
	if err = global.GVA_DB.Where("id = ?", ID).First(&database).Error; err != nil {
		return response, fmt.Errorf("找不到数据库配置")
	}

	start := time.Now()
	
	// 目前只支持MySQL
	if database.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			database.Username, database.Password, database.Host, database.Port, database.Database)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			response.Success = false
			response.Message = "连接失败: " + err.Error()
			return response, nil
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			response.Success = false
			response.Message = "连接测试失败: " + err.Error()
			return response, nil
		}

		// 更新最后测试时间和状态
		now := time.Now()
		global.GVA_DB.Model(&database).Updates(map[string]interface{}{
			"status":       "online",
			"last_test_at": &now,
		})

		response.Success = true
		response.Message = "MySQL连接测试成功"
		response.ResponseTime = time.Since(start).Milliseconds()
		
		// 获取MySQL版本
		var version string
		if err := db.QueryRow("SELECT VERSION()").Scan(&version); err == nil {
			response.Version = version
		}
		
		return response, nil
	}

	response.Success = false
	response.Message = "暂时只支持MySQL数据库"
	return response, nil
}

// SyncSysDatabase 从远程服务器同步数据库配置
func (databaseService *DatabaseService) SyncSysDatabase(req systemReq.SysDatabaseSync) (response systemRes.SysDatabaseSyncResponse, err error) {
	// 从远程服务器获取数据库列表
	remoteDatabases, err := databaseService.fetchRemoteDatabases(req.ServerURL, req.APIKey)
	if err != nil {
		return response, fmt.Errorf("获取远程数据库列表失败: %v", err)
	}

	// 执行同步
	stats := databaseService.syncDatabases(remoteDatabases, req.Overwrite, req.SyncType)

	response.Success = true
	response.Message = "同步完成"
	response.TotalSynced = stats.total
	response.Created = stats.created
	response.Updated = stats.updated
	response.Skipped = stats.skipped
	response.Errors = stats.errors

	return response, nil
}

// 同步统计信息
type syncStats struct {
	total   int
	created int
	updated int
	skipped int
	errors  []string
}

// fetchRemoteDatabases 从远程服务器获取数据库列表
func (databaseService *DatabaseService) fetchRemoteDatabases(serverURL, apiKey string) ([]system.SysDatabase, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", serverURL+"/api/database/list", nil)
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("远程服务器返回错误状态码: %d", resp.StatusCode)
	}

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			List []system.SysDatabase `json:"list"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("远程服务器返回错误: %s", result.Msg)
	}

	return result.Data.List, nil
}

// syncDatabases 执行数据库同步
func (databaseService *DatabaseService) syncDatabases(remoteDatabases []system.SysDatabase, overwrite bool, syncType string) syncStats {
	stats := syncStats{}

	for _, remoteDB := range remoteDatabases {
		// 检查本地是否已存在
		var existingDB system.SysDatabase
		err := global.GVA_DB.Where("name = ? AND host = ? AND port = ?",
			remoteDB.Name, remoteDB.Host, remoteDB.Port).First(&existingDB).Error

		if err != nil {
			// 不存在，创建新记录
			if err == gorm.ErrRecordNotFound {
				newDB := system.SysDatabase{
					Name:        remoteDB.Name,
					Type:        remoteDB.Type,
					Host:        remoteDB.Host,
					Port:        remoteDB.Port,
					Username:    remoteDB.Username,
					Database:    remoteDB.Database,
					Status:      remoteDB.Status,
					Description: remoteDB.Description,
				}

				if err := global.GVA_DB.Create(&newDB).Error; err != nil {
					stats.errors = append(stats.errors, fmt.Sprintf("创建数据库 %s 失败: %v", remoteDB.Name, err))
				} else {
					stats.created++
				}
			} else {
				stats.errors = append(stats.errors, fmt.Sprintf("查询数据库 %s 失败: %v", remoteDB.Name, err))
			}
		} else {
			// 已存在
			if overwrite {
				// 更新现有记录
				updates := map[string]interface{}{
					"type":        remoteDB.Type,
					"username":    remoteDB.Username,
					"database":    remoteDB.Database,
					"status":      remoteDB.Status,
					"description": remoteDB.Description,
				}

				if err := global.GVA_DB.Model(&existingDB).Updates(updates).Error; err != nil {
					stats.errors = append(stats.errors, fmt.Sprintf("更新数据库 %s 失败: %v", remoteDB.Name, err))
				} else {
					stats.updated++
				}
			} else {
				// 跳过
				stats.skipped++
			}
		}

		stats.total++
	}

	return stats
}