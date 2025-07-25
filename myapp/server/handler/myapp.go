package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"myapp/global"
	"myapp/models"
	pb "myapp/proto"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Server gRPC服务器结构体
type Server struct {
	pb.UnimplementedMyappServer
}

// DatabaseCreate 创建数据库连接配置
func (s *Server) DatabaseCreate(ctx context.Context, req *pb.DatabaseCreateRequest) (*pb.DatabaseCreateResponse, error) {
	// 检查名称是否重复
	var count int64
	if err := global.DB.Model(&models.SysDatabases{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	if count > 0 {
		return nil, fmt.Errorf("数据库连接名称已存在")
	}
	database := &models.SysDatabases{
		Name:        req.Name,        // 数据库连接名称
		Type:        req.Type,        // 数据库类型(mysql/postgresql等)
		Host:        req.Host,        // 主机地址
		Port:        req.Port,        // 端口号
		Username:    req.Username,    // 用户名
		Password:    req.Password,    // 密码
		Database:    req.Database,    // 数据库名
		Status:      "active",        // 状态设为激活
		Description: req.Description, // 描述
		CreatedAt:   time.Now(),      // 创建时间
		UpdatedAt:   time.Now(),      // 更新时间
	}
	if err := global.DB.Create(database).Error; err != nil {
		return nil, fmt.Errorf("保存失败")
	}
	return &pb.DatabaseCreateResponse{
		Success:    true,
		Message:    "数据库连接创建成功",
		DatabaseId: database.Id,
	}, nil
}

// DatabaseTest 测试数据库连接
func (s *Server) DatabaseTest(ctx context.Context, req *pb.DatabaseTestRequest) (*pb.DatabaseTestResponse, error) {
	var database models.SysDatabases
	//根据id去查存不存在
	if err := global.DB.First(&database, req.DatabaseId).Error; err != nil {
		return nil, fmt.Errorf("找不到数据库配置")
	}

	//简单的连接测试（只支持MySQL）
	if database.Type == "mysql" {
		// 构建连接字符串
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			database.Username, database.Password, database.Host, database.Port, database.Database)

		// 尝试连接
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("连接失败")
		}
		defer db.Close()

		// 测试连接
		if err := db.Ping(); err != nil {
			return &pb.DatabaseTestResponse{
				Success: false,
				Message: "连接测试失败: " + err.Error(),
			}, nil
		}

		return &pb.DatabaseTestResponse{
			Success: true,
			Message: "MySQL连接测试成功",
		}, nil
	}

	// 其他数据库类型暂时不支持
	return &pb.DatabaseTestResponse{
		Success: false,
		Message: "暂时只支持MySQL数据库",
	}, nil
}

// DatabaseList 获取数据库连接列表
func (s *Server) DatabaseList(ctx context.Context, req *pb.DatabaseListRequest) (*pb.DatabaseListResponse, error) {
	var databases []models.SysDatabases
	if err := global.DB.Find(&databases).Error; err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	var dbInfos []*pb.DatabaseInfo
	for _, db := range databases {
		dbInfo := &pb.DatabaseInfo{
			Id:          db.Id,
			Name:        db.Name,
			Type:        db.Type,
			Host:        db.Host,
			Port:        db.Port,
			Username:    db.Username,
			Database:    db.Database,
			Status:      db.Status,
			Description: db.Description,
			CreatedAt:   db.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		dbInfos = append(dbInfos, dbInfo)
	}
	return &pb.DatabaseListResponse{
		Success:   true,
		Message:   "获取数据库列表成功",
		Databases: dbInfos,
	}, nil
}

// DatabaseUpdate 更新数据库连接配置
func (s *Server) DatabaseUpdate(ctx context.Context, req *pb.DatabaseUpdateRequest) (*pb.DatabaseUpdateResponse, error) {
	var database models.SysDatabases
	//根据id去查存不存在
	if err := global.DB.First(&database, req.Id).Error; err != nil {
		return nil, fmt.Errorf("找不到数据库配置")
	}
	// 检查名称是否重复（排除当前记录）
	var count int64
	if err := global.DB.Model(&models.SysDatabases{}).Where("name = ? AND id != ?", req.Name, req.Id).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("查询失败")
	}
	if count > 0 {
		return nil, fmt.Errorf("数据库连接名称已存在")
	}
	//更新数据
	update := &models.SysDatabases{
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		Description: req.Description,
	}

	//保存到数据库 - 使用WHERE条件指定要更新的记录
	if err := global.DB.Model(&database).Updates(update).Error; err != nil {
		return nil, fmt.Errorf("更新失败")
	}
	return &pb.DatabaseUpdateResponse{
		Success: true,
		Message: "数据库连接更新成功",
	}, nil
}

// DatabaseDelete 删除数据库连接配置
func (s *Server) DatabaseDelete(ctx context.Context, req *pb.DatabaseDeleteRequest) (*pb.DatabaseDeleteResponse, error) {
	// 先检查记录是否存在
	var database models.SysDatabases
	if err := global.DB.First(&database, req.Id).Error; err != nil {
		return nil, fmt.Errorf("找不到要删除的数据库配置")
	}

	//从数据库中删除这个配置
	del := &models.SysDatabases{
		Id: database.Id,
	}
	if err := global.DB.Delete(&del).Error; err != nil {
		return nil, fmt.Errorf("删除失败")
	}
	return &pb.DatabaseDeleteResponse{
		Success: true,
		Message: "数据库连接删除成功",
	}, nil
}

// DatabaseGet 获取单个数据库连接配置
func (s *Server) DatabaseGet(ctx context.Context, req *pb.DatabaseGetRequest) (*pb.DatabaseGetResponse, error) {
	//根据id去查存不存在
	var database models.SysDatabases
	if err := global.DB.First(&database, req.Id).Error; err != nil {
		return nil, fmt.Errorf("找不到数据库配置")
	}
	//把数据转换成前端需要的格式
	dbInfo := &pb.DatabaseInfo{
		Id:          database.Id,
		Name:        database.Name,
		Type:        database.Type,
		Host:        database.Host,
		Port:        database.Port,
		Username:    database.Username,
		Database:    database.Database,
		Status:      database.Status,
		Description: database.Description,
		CreatedAt:   database.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return &pb.DatabaseGetResponse{
		Success:  true,
		Message:  "获取数据库配置成功",
		Database: dbInfo,
	}, nil
}

// DatabaseSync 从远程服务器同步数据库配置
func (s *Server) DatabaseSync(ctx context.Context, req *pb.DatabaseSyncRequest) (*pb.DatabaseSyncResponse, error) {
	// 1. 验证请求参数
	if req.ServerUrl == "" {
		return nil, fmt.Errorf("服务器URL不能为空")
	}

	// 2. 从远程服务器获取数据库列表
	remoteDatabases, err := s.fetchRemoteDatabases(req.ServerUrl, req.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("获取远程数据库列表失败: %v", err)
	}

	// 3. 执行同步
	stats := s.syncDatabases(remoteDatabases, req.Overwrite, req.SyncType)

	return &pb.DatabaseSyncResponse{
		Success:     true,
		Message:     "同步完成",
		TotalSynced: int32(stats.total),
		Created:     int32(stats.created),
		Updated:     int32(stats.updated),
		Skipped:     int32(stats.skipped),
		Errors:      stats.errors,
	}, nil
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
func (s *Server) fetchRemoteDatabases(serverURL, apiKey string) ([]*pb.RemoteDatabaseInfo, error) {
	// 这里实现HTTP请求到远程服务器
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", serverURL+"/api/database/list", nil)
	if err != nil {
		return nil, err
	}

	// 只有在提供了API密钥时才设置Authorization头
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
			List []struct {
				Id          uint64 `json:"id"`
				Name        string `json:"name"`
				Type        string `json:"type"`
				Host        string `json:"host"`
				Port        int32  `json:"port"`
				Username    string `json:"username"`
				Database    string `json:"database"`
				Status      string `json:"status"`
				Description string `json:"description"`
				CreatedAt   string `json:"created_at"`
				UpdatedAt   string `json:"updated_at"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("远程服务器返回错误: %s", result.Msg)
	}

	// 转换为pb.RemoteDatabaseInfo格式
	var remoteDatabases []*pb.RemoteDatabaseInfo
	for _, db := range result.Data.List {
		remoteDatabases = append(remoteDatabases, &pb.RemoteDatabaseInfo{
			Id:          db.Id,
			Name:        db.Name,
			Type:        db.Type,
			Host:        db.Host,
			Port:        db.Port,
			Username:    db.Username,
			Database:    db.Database,
			Status:      db.Status,
			Description: db.Description,
			CreatedAt:   db.CreatedAt,
			UpdatedAt:   db.UpdatedAt,
		})
	}

	return remoteDatabases, nil
}

// syncDatabases 执行数据库同步
func (s *Server) syncDatabases(remoteDatabases []*pb.RemoteDatabaseInfo, overwrite bool, syncType string) syncStats {
	stats := syncStats{}

	for _, remoteDB := range remoteDatabases {
		// 检查本地是否已存在
		var existingDB models.SysDatabases
		err := global.DB.Where("name = ? AND host = ? AND port = ?",
			remoteDB.Name, remoteDB.Host, remoteDB.Port).First(&existingDB).Error

		if err != nil {
			// 不存在，创建新记录
			if err == gorm.ErrRecordNotFound {
				newDB := models.SysDatabases{
					Name:        remoteDB.Name,
					Type:        remoteDB.Type,
					Host:        remoteDB.Host,
					Port:        remoteDB.Port,
					Username:    remoteDB.Username,
					Database:    remoteDB.Database,
					Status:      remoteDB.Status,
					Description: remoteDB.Description,
				}

				if err := global.DB.Create(&newDB).Error; err != nil {
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
					"updated_at":  time.Now(),
				}

				if err := global.DB.Model(&existingDB).Updates(updates).Error; err != nil {
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
