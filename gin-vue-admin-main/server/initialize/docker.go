package initialize

import (
	"context"
	"time"

	"github.com/docker/docker/client"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

func Docker() {
	dockerConfig := global.GVA_CONFIG.Docker
	
	// 设置默认值
	if dockerConfig.Host == "" {
		dockerConfig.Host = "unix:///var/run/docker.sock"
	}
	if dockerConfig.Version == "" {
		dockerConfig.Version = "1.41"
	}
	if dockerConfig.Timeout == 0 {
		dockerConfig.Timeout = 30
	}

	// 创建Docker客户端
	cli, err := client.NewClientWithOpts(
		client.WithHost(dockerConfig.Host),
		client.WithVersion(dockerConfig.Version),
		client.WithTimeout(time.Duration(dockerConfig.Timeout)*time.Second),
	)
	if err != nil {
		global.GVA_LOG.Error("Docker client initialization failed", zap.Error(err))
		return
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dockerConfig.Timeout)*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		global.GVA_LOG.Warn("Docker daemon is not accessible", zap.Error(err))
		// 不要返回错误，允许应用继续运行，但Docker功能将不可用
	} else {
		global.GVA_LOG.Info("Docker client connected successfully")
	}

	global.GVA_DOCKER = cli
}

// GetDockerClient 获取Docker客户端，如果未初始化则返回nil
func GetDockerClient() *client.Client {
	return global.GVA_DOCKER
}

// IsDockerAvailable 检查Docker是否可用
func IsDockerAvailable() bool {
	if global.GVA_DOCKER == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := global.GVA_DOCKER.Ping(ctx)
	return err == nil
}