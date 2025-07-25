package handler

import (
	__ "api/proto"
	"api/request"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

func DatabaseCreate(c *gin.Context) {
	var req request.DatabaseCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseCreate(c, &__.DatabaseCreateRequest{
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "数据库创建失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据库创建成功",
		"data": create,
	})
	return

}

func DatabaseList(c *gin.Context) {
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseList(c, &__.DatabaseListRequest{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "展示数据库失败",
			"data": err.Error(),
		})
		return
	}

	// 转换数据格式以匹配前端期望
	var list []gin.H
	for _, db := range create.Databases {
		list = append(list, gin.H{
			"id":          db.Id,
			"name":        db.Name,
			"type":        db.Type,
			"host":        db.Host,
			"port":        db.Port,
			"username":    db.Username,
			"database":    db.Database,
			"status":      db.Status,
			"description": db.Description,
			"created_at":  db.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "展示数据库成功",
		"data": gin.H{
			"list":  list,
			"total": len(list),
		},
	})
	return
}

func DatabaseGet(c *gin.Context) {
	var req request.DatabaseGetRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseGet(c, &__.DatabaseGetRequest{
		Id: req.ID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "单条数据查询失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "单条数据查询成功",
		"data": create,
	})
	return
}

func DatabaseUpdate(c *gin.Context) {
	var req request.DatabaseUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseUpdate(c, &__.DatabaseUpdateRequest{
		Id:          req.ID,
		Name:        req.Name,
		Type:        req.Type,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		Database:    req.Database,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "数据库修改失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据库修改成功",
		"data": create,
	})
	return
}

func DatabaseDelete(c *gin.Context) {
	var req request.DatabaseDeleteRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseDelete(c, &__.DatabaseDeleteRequest{
		Id: req.ID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "数据库删除失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据库删除成功",
		"data": create,
	})
	return
}

func DatabaseTest(c *gin.Context) {
	var req request.DatabaseTestRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)
	create, err := c1.DatabaseTest(c, &__.DatabaseTestRequest{
		DatabaseId: req.DatabaseID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "连接测试失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "连接测试成功",
		"data": create,
	})
	return
}

func DatabaseSync(c *gin.Context) {
	var req request.DatabaseSyncRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "验证失败",
			"data": err.Error(),
		})
		return
	}

	conn, err := grpc.NewClient("127.0.0.1:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c1 := __.NewMyappClient(conn)

	response, err := c1.DatabaseSync(c, &__.DatabaseSyncRequest{
		ServerUrl: req.ServerURL,
		ApiKey:    req.APIKey,
		SyncType:  req.SyncType,
		Overwrite: req.Overwrite,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "数据库同步失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据库同步成功",
		"data": gin.H{
			"success":      response.Success,
			"message":      response.Message,
			"total_synced": response.TotalSynced,
			"created":      response.Created,
			"updated":      response.Updated,
			"skipped":      response.Skipped,
			"errors":       response.Errors,
		},
	})
	return
}
