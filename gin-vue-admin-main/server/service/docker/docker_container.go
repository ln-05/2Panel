package docker

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerContainerService struct{}

// GetContainerList 获取容器列表
func (d *DockerContainerService) GetContainerList(filter request.ContainerFilter) ([]response.ContainerInfo, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, 0, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 构建过滤器
	filterArgs := filters.NewArgs()

	// 状态过滤
	if filter.Status != "" {
		filterArgs.Add("status", filter.Status)
	}

	// 名称过滤
	if filter.Name != "" {
		filterArgs.Add("name", filter.Name)
	}

	// 设置列表选项
	options := types.ContainerListOptions{
		All:     true, // 显示所有容器（包括停止的）
		Filters: filterArgs,
	}

	// 调用Docker API获取容器列表
	containers, err := global.GVA_DOCKER.ContainerList(ctx, options)
	if err != nil {
		global.GVA_LOG.Error("Failed to get container list", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get container list: %v", err)
	}

	// 转换为响应模型
	containerInfos := make([]response.ContainerInfo, 0, len(containers))
	for _, dockerContainer := range containers {
		containerInfo := dockerModel.ConvertToContainerInfo(dockerContainer)
		containerInfos = append(containerInfos, containerInfo)
	}

	// 实现分页逻辑
	total := int64(len(containerInfos))

	// 计算分页
	if filter.Page > 0 && filter.PageSize > 0 {
		start := (filter.Page - 1) * filter.PageSize
		end := start + filter.PageSize

		if start >= len(containerInfos) {
			return []response.ContainerInfo{}, total, nil
		}

		if end > len(containerInfos) {
			end = len(containerInfos)
		}

		containerInfos = containerInfos[start:end]
	}

	return containerInfos, total, nil
}

// GetContainerDetail 获取容器详细信息
func (d *DockerContainerService) GetContainerDetail(containerID string) (*response.ContainerDetail, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	// 验证容器ID
	if containerID == "" {
		return nil, fmt.Errorf("container ID cannot be empty")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用Docker API获取容器详细信息
	containerJSON, err := global.GVA_DOCKER.ContainerInspect(ctx, containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return nil, fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to get container detail", zap.String("containerID", containerID), zap.Error(err))
		return nil, fmt.Errorf("failed to get container detail: %v", err)
	}

	// 转换为响应模型
	containerDetail := dockerModel.ConvertToContainerDetail(containerJSON)

	return &containerDetail, nil
}

// GetContainerLogs 获取容器日志
func (d *DockerContainerService) GetContainerLogs(containerID string, options request.LogOptions) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return "", fmt.Errorf("Docker client is not available")
	}

	// 验证容器ID
	if containerID == "" {
		return "", fmt.Errorf("container ID cannot be empty")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 构建日志选项
	logOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     options.Follow,
		Timestamps: true,
	}

	// 设置tail选项
	if options.Tail != "" {
		logOptions.Tail = options.Tail
	} else {
		logOptions.Tail = "100" // 默认显示最后100行
	}

	// 设置since选项
	if options.Since != "" {
		logOptions.Since = options.Since
	}

	// 调用Docker API获取容器日志
	logReader, err := global.GVA_DOCKER.ContainerLogs(ctx, containerID, logOptions)
	if err != nil {
		if client.IsErrNotFound(err) {
			return "", fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to get container logs", zap.String("containerID", containerID), zap.Error(err))
		return "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logReader.Close()

	// 读取日志内容
	logBytes, err := io.ReadAll(logReader)
	if err != nil {
		global.GVA_LOG.Error("Failed to read container logs", zap.String("containerID", containerID), zap.Error(err))
		return "", fmt.Errorf("failed to read container logs: %v", err)
	}

	// Docker日志格式包含8字节的头部信息，需要处理
	logContent := d.processDockerLogs(logBytes)

	return logContent, nil
}

// processDockerLogs 处理Docker日志格式
func (d *DockerContainerService) processDockerLogs(logBytes []byte) string {
	if len(logBytes) == 0 {
		return ""
	}

	var result strings.Builder
	i := 0

	for i < len(logBytes) {
		// Docker日志每行都有8字节的头部
		if i+8 > len(logBytes) {
			break
		}

		// 跳过前4字节（stream type和padding）
		// 第5-8字节是消息长度（big-endian）
		if i+4 < len(logBytes) {
			length := int(logBytes[i+4])<<24 | int(logBytes[i+5])<<16 | int(logBytes[i+6])<<8 | int(logBytes[i+7])

			// 跳过8字节头部
			i += 8

			// 读取实际的日志内容
			if i+length <= len(logBytes) {
				result.Write(logBytes[i : i+length])
				i += length
			} else {
				// 如果长度不匹配，直接添加剩余内容
				result.Write(logBytes[i:])
				break
			}
		} else {
			break
		}
	}

	// 如果处理失败，返回原始内容（去掉可能的控制字符）
	if result.Len() == 0 {
		return string(logBytes)
	}

	return result.String()
}

// IsDockerAvailable 检查Docker是否可用
func (d *DockerContainerService) IsDockerAvailable() bool {
	if global.GVA_DOCKER == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := global.GVA_DOCKER.Ping(ctx)
	return err == nil
}

// GetDockerInfo 获取Docker系统信息
func (d *DockerContainerService) GetDockerInfo() (*types.Info, error) {
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := global.GVA_DOCKER.Info(ctx)
	if err != nil {
		global.GVA_LOG.Error("Failed to get Docker info", zap.Error(err))
		return nil, fmt.Errorf("failed to get Docker info: %v", err)
	}

	return &info, nil
}

// StartContainer 启动容器
func (d *DockerContainerService) StartContainer(containerID string) error {
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := global.GVA_DOCKER.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			return fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to start container", zap.String("containerID", containerID), zap.Error(err))
		return fmt.Errorf("failed to start container: %v", err)
	}

	global.GVA_LOG.Info("Container started successfully", zap.String("containerID", containerID))
	return nil
}

// StopContainer 停止容器
func (d *DockerContainerService) StopContainer(containerID string, timeout *int) error {
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 设置停止超时时间
	var stopTimeout *time.Duration
	if timeout != nil {
		duration := time.Duration(*timeout) * time.Second
		stopTimeout = &duration
	}

	err := global.GVA_DOCKER.ContainerStop(ctx, containerID, stopTimeout)
	if err != nil {
		if client.IsErrNotFound(err) {
			return fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to stop container", zap.String("containerID", containerID), zap.Error(err))
		return fmt.Errorf("failed to stop container: %v", err)
	}

	global.GVA_LOG.Info("Container stopped successfully", zap.String("containerID", containerID))
	return nil
}

// RestartContainer 重启容器
func (d *DockerContainerService) RestartContainer(containerID string, timeout *int) error {
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 设置重启超时时间
	var restartTimeout *time.Duration
	if timeout != nil {
		duration := time.Duration(*timeout) * time.Second
		restartTimeout = &duration
	}

	err := global.GVA_DOCKER.ContainerRestart(ctx, containerID, restartTimeout)
	if err != nil {
		if client.IsErrNotFound(err) {
			return fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to restart container", zap.String("containerID", containerID), zap.Error(err))
		return fmt.Errorf("failed to restart container: %v", err)
	}

	global.GVA_LOG.Info("Container restarted successfully", zap.String("containerID", containerID))
	return nil
}

// RemoveContainer 删除容器
func (d *DockerContainerService) RemoveContainer(containerID string, force bool) error {
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := global.GVA_DOCKER.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: force,
	})
	if err != nil {
		if client.IsErrNotFound(err) {
			return fmt.Errorf("container not found")
		}
		global.GVA_LOG.Error("Failed to remove container", zap.String("containerID", containerID), zap.Error(err))
		return fmt.Errorf("failed to remove container: %v", err)
	}

	global.GVA_LOG.Info("Container removed successfully", zap.String("containerID", containerID))
	return nil
}

// BatchOperateByOrchestrationLabel 对同一label分组的容器批量操作
func (d *DockerContainerService) BatchOperateByOrchestrationLabel(label string, op string, timeout *int, force bool) (successIDs []string, failed map[string]string) {
	if global.GVA_DOCKER == nil {
		return nil, map[string]string{"_global": "Docker client is not available"}
	}
	ctx := context.Background()
	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, map[string]string{"_global": err.Error()}
	}
	failed = make(map[string]string)
	for _, ctn := range containers {
		var orchestrationName string
		
		// 优先检查自定义的orchestration标签
		if labelValue := ctn.Labels["orchestration"]; labelValue != "" {
			orchestrationName = labelValue
		} else if composeProject := ctn.Labels["com.docker.compose.project"]; composeProject != "" {
			// Docker Compose项目标签
			orchestrationName = composeProject
		} else if panelProject := ctn.Labels["com.1panel.compose.project"]; panelProject != "" {
			// 1Panel特有的项目标签
			orchestrationName = panelProject
		} else if panelApp := ctn.Labels["1panel.app"]; panelApp != "" {
			// 1Panel应用标签
			orchestrationName = panelApp
		}
		
		if orchestrationName != label {
			continue
		}
		var opErr error
		switch op {
		case "start":
			opErr = d.StartContainer(ctn.ID)
		case "stop":
			opErr = d.StopContainer(ctn.ID, timeout)
		case "restart":
			opErr = d.RestartContainer(ctn.ID, timeout)
		case "delete":
			opErr = d.RemoveContainer(ctn.ID, force)
		default:
			failed[ctn.ID] = "unsupported operation"
			continue
		}
		if opErr != nil {
			failed[ctn.ID] = opErr.Error()
		} else {
			successIDs = append(successIDs, ctn.ID)
		}
	}
	return
}

// GetOrchestrationList 获取编排列表
func (d *DockerContainerService) GetOrchestrationList(page, pageSize int, search, statusFilter string) (interface{}, error) {
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		global.GVA_LOG.Error("Failed to get container list for orchestration", zap.Error(err))
		return nil, fmt.Errorf("failed to get container list: %v", err)
	}

	// 分组 - 支持多种编排标签
	orchestrationMap := make(map[string][]types.Container)
	for _, ctn := range containers {
		var orchestrationName string
		
		// 优先检查自定义的orchestration标签
		if label := ctn.Labels["orchestration"]; label != "" {
			orchestrationName = label
		} else if composeProject := ctn.Labels["com.docker.compose.project"]; composeProject != "" {
			// Docker Compose项目标签
			orchestrationName = composeProject
		} else if panelProject := ctn.Labels["com.1panel.compose.project"]; panelProject != "" {
			// 1Panel特有的项目标签
			orchestrationName = panelProject
		} else if panelApp := ctn.Labels["1panel.app"]; panelApp != "" {
			// 1Panel应用标签
			orchestrationName = panelApp
		}
		
		if orchestrationName != "" {
			orchestrationMap[orchestrationName] = append(orchestrationMap[orchestrationName], ctn)
		}
	}

	// 组装结果
	type OrchestrationListItem struct {
		Name           string    `json:"name"`
		Source         string    `json:"source"`
		Dir            string    `json:"dir"`
		Status         string    `json:"status"`
		ContainerCount int       `json:"containerCount"`
		CreatedAt      time.Time `json:"createdAt"`
	}
	
	var list []OrchestrationListItem
	for name, group := range orchestrationMap {
		if search != "" && !strings.Contains(name, search) {
			continue
		}
		var running, stopped int
		var earliest time.Time
		for i, ctn := range group {
			if ctn.State == "running" {
				running++
			} else {
				stopped++
			}
			created := time.Unix(ctn.Created, 0)
			if i == 0 || created.Before(earliest) {
				earliest = created
			}
		}
		status := "mixed"
		if running == len(group) {
			status = "running"
		} else if stopped == len(group) {
			status = "stopped"
		}
		if statusFilter != "" && status != statusFilter {
			continue
		}
		list = append(list, OrchestrationListItem{
			Name:           name,
			Source:         "应用商店", // 可根据实际情况调整
			Dir:            "-",    // 如有目录信息可补充
			Status:         status,
			ContainerCount: len(group),
			CreatedAt:      earliest,
		})
	}
	
	// 排序
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	total := len(list)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paged := list[start:end]
	
	return map[string]interface{}{
		"list":  paged,
		"total": total,
	}, nil
}

// GetOrchestrationDetail 获取编排详情
func (d *DockerContainerService) GetOrchestrationDetail(name string) ([]types.Container, error) {
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		global.GVA_LOG.Error("Failed to get container list for orchestration detail", zap.String("name", name), zap.Error(err))
		return nil, fmt.Errorf("failed to get container list: %v", err)
	}

	var group []types.Container
	for _, ctn := range containers {
		var orchestrationName string
		
		// 优先检查自定义的orchestration标签
		if label := ctn.Labels["orchestration"]; label != "" {
			orchestrationName = label
		} else if composeProject := ctn.Labels["com.docker.compose.project"]; composeProject != "" {
			// Docker Compose项目标签
			orchestrationName = composeProject
		} else if panelProject := ctn.Labels["com.1panel.compose.project"]; panelProject != "" {
			// 1Panel特有的项目标签
			orchestrationName = panelProject
		} else if panelApp := ctn.Labels["1panel.app"]; panelApp != "" {
			// 1Panel应用标签
			orchestrationName = panelApp
		}
		
		if orchestrationName == name {
			group = append(group, ctn)
		}
	}

	if len(group) == 0 {
		return nil, fmt.Errorf("orchestration not found")
	}

	return group, nil
}

// GetOrchestrationStatus 获取编排状态
func (d *DockerContainerService) GetOrchestrationStatus(name string) (map[string]interface{}, error) {
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		global.GVA_LOG.Error("Failed to get container list for orchestration status", zap.String("name", name), zap.Error(err))
		return nil, fmt.Errorf("failed to get container list: %v", err)
	}

	var group []types.Container
	for _, ctn := range containers {
		var orchestrationName string
		
		// 优先检查自定义的orchestration标签
		if label := ctn.Labels["orchestration"]; label != "" {
			orchestrationName = label
		} else if composeProject := ctn.Labels["com.docker.compose.project"]; composeProject != "" {
			// Docker Compose项目标签
			orchestrationName = composeProject
		} else if panelProject := ctn.Labels["com.1panel.compose.project"]; panelProject != "" {
			// 1Panel特有的项目标签
			orchestrationName = panelProject
		} else if panelApp := ctn.Labels["1panel.app"]; panelApp != "" {
			// 1Panel应用标签
			orchestrationName = panelApp
		}
		
		if orchestrationName == name {
			group = append(group, ctn)
		}
	}

	if len(group) == 0 {
		return nil, fmt.Errorf("orchestration not found")
	}

	var running, stopped int
	for _, ctn := range group {
		if ctn.State == "running" {
			running++
		} else {
			stopped++
		}
	}

	status := "mixed"
	if running == len(group) {
		status = "running"
	} else if stopped == len(group) {
		status = "stopped"
	}

	return map[string]interface{}{
		"name":    name,
		"status":  status,
		"running": running,
		"stopped": stopped,
		"total":   len(group),
	}, nil
}

// DeleteOrchestration 删除编排（删除所有相关容器）
func (d *DockerContainerService) DeleteOrchestration(name string) ([]string, error) {
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	containers, err := global.GVA_DOCKER.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		global.GVA_LOG.Error("Failed to get container list for orchestration deletion", zap.String("name", name), zap.Error(err))
		return nil, fmt.Errorf("failed to get container list: %v", err)
	}

	var failed []string
	for _, ctn := range containers {
		var orchestrationName string
		
		// 优先检查自定义的orchestration标签
		if label := ctn.Labels["orchestration"]; label != "" {
			orchestrationName = label
		} else if composeProject := ctn.Labels["com.docker.compose.project"]; composeProject != "" {
			// Docker Compose项目标签
			orchestrationName = composeProject
		} else if panelProject := ctn.Labels["com.1panel.compose.project"]; panelProject != "" {
			// 1Panel特有的项目标签
			orchestrationName = panelProject
		} else if panelApp := ctn.Labels["1panel.app"]; panelApp != "" {
			// 1Panel应用标签
			orchestrationName = panelApp
		}
		
		if orchestrationName == name {
			err := global.GVA_DOCKER.ContainerRemove(ctx, ctn.ID, types.ContainerRemoveOptions{Force: true})
			if err != nil {
				global.GVA_LOG.Error("Failed to remove container in orchestration", zap.String("containerID", ctn.ID), zap.String("orchestration", name), zap.Error(err))
				failed = append(failed, ctn.ID)
			}
		}
	}

	if len(failed) > 0 {
		return failed, fmt.Errorf("some containers failed to delete")
	}

	global.GVA_LOG.Info("Orchestration deleted successfully", zap.String("name", name))
	return nil, nil
}
