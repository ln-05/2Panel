package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
)

type DockerImageService struct{}

// GetImageList 获取镜像列表
func (d *DockerImageService) GetImageList(filter request.ImageFilter) ([]response.ImageInfo, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, 0, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 构建过滤器
	filterArgs := filters.NewArgs()
	
	// 名称过滤
	if filter.Name != "" {
		filterArgs.Add("reference", filter.Name)
	}
	
	// 悬空镜像过滤
	if filter.Dangling != nil {
		if *filter.Dangling {
			filterArgs.Add("dangling", "true")
		} else {
			filterArgs.Add("dangling", "false")
		}
	}

	// 设置列表选项
	options := types.ImageListOptions{
		All:     true,
		Filters: filterArgs,
	}

	// 调用Docker API获取镜像列表
	images, err := global.GVA_DOCKER.ImageList(ctx, options)
	if err != nil {
		global.GVA_LOG.Error("Failed to get image list", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get image list: %v", err)
	}

	// 转换为响应模型
	imageInfos := make([]response.ImageInfo, 0, len(images))
	for _, dockerImage := range images {
		imageInfo := d.convertToImageInfo(dockerImage)
		imageInfos = append(imageInfos, imageInfo)
	}

	// 实现分页逻辑
	total := int64(len(imageInfos))
	
	// 计算分页
	if filter.Page > 0 && filter.PageSize > 0 {
		start := (filter.Page - 1) * filter.PageSize
		end := start + filter.PageSize
		
		if start >= len(imageInfos) {
			return []response.ImageInfo{}, total, nil
		}
		
		if end > len(imageInfos) {
			end = len(imageInfos)
		}
		
		imageInfos = imageInfos[start:end]
	}

	return imageInfos, total, nil
}

// GetImageDetail 获取镜像详细信息
func (d *DockerImageService) GetImageDetail(imageID string) (*response.ImageDetail, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	// 验证镜像ID
	if imageID == "" {
		return nil, fmt.Errorf("image ID cannot be empty")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用Docker API获取镜像详细信息
	imageInspect, _, err := global.GVA_DOCKER.ImageInspectWithRaw(ctx, imageID)
	if err != nil {
		global.GVA_LOG.Error("Failed to get image detail", zap.String("imageID", imageID), zap.Error(err))
		return nil, fmt.Errorf("failed to get image detail: %v", err)
	}

	// 转换为响应模型
	imageDetail := d.convertToImageDetail(imageInspect)

	return &imageDetail, nil
}

// PullImage 拉取镜像
func (d *DockerImageService) PullImage(pullReq request.ImagePullRequest) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return "", fmt.Errorf("Docker client is not available")
	}

	// 构建完整的镜像名称
	imageName := pullReq.Image
	if pullReq.Tag != "" && !strings.Contains(imageName, ":") {
		imageName = fmt.Sprintf("%s:%s", imageName, pullReq.Tag)
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second) // 拉取镜像可能需要较长时间
	defer cancel()

	// 拉取镜像
	reader, err := global.GVA_DOCKER.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		global.GVA_LOG.Error("Failed to pull image", zap.String("image", imageName), zap.Error(err))
		return "", fmt.Errorf("failed to pull image: %v", err)
	}
	defer reader.Close()

	// 读取拉取日志
	pullLog, err := io.ReadAll(reader)
	if err != nil {
		global.GVA_LOG.Error("Failed to read pull log", zap.String("image", imageName), zap.Error(err))
		return "", fmt.Errorf("failed to read pull log: %v", err)
	}

	global.GVA_LOG.Info("Image pulled successfully", zap.String("image", imageName))
	return string(pullLog), nil
}

// RemoveImage 删除镜像
func (d *DockerImageService) RemoveImage(imageID string, force bool) error {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	if imageID == "" {
		return fmt.Errorf("image ID cannot be empty")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 删除镜像
	_, err := global.GVA_DOCKER.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         force,
		PruneChildren: true,
	})
	if err != nil {
		global.GVA_LOG.Error("Failed to remove image", zap.String("imageID", imageID), zap.Error(err))
		return fmt.Errorf("failed to remove image: %v", err)
	}

	global.GVA_LOG.Info("Image removed successfully", zap.String("imageID", imageID))
	return nil
}

// TagImage 给镜像打标签
func (d *DockerImageService) TagImage(tagReq request.ImageTagRequest) error {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 给镜像打标签
	err := global.GVA_DOCKER.ImageTag(ctx, tagReq.SourceImage, tagReq.TargetImage)
	if err != nil {
		global.GVA_LOG.Error("Failed to tag image", 
			zap.String("source", tagReq.SourceImage), 
			zap.String("target", tagReq.TargetImage), 
			zap.Error(err))
		return fmt.Errorf("failed to tag image: %v", err)
	}

	global.GVA_LOG.Info("Image tagged successfully", 
		zap.String("source", tagReq.SourceImage), 
		zap.String("target", tagReq.TargetImage))
	return nil
}

// PruneImages 清理未使用的镜像
func (d *DockerImageService) PruneImages(dangling bool) (int64, int64, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return 0, 0, fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 构建过滤器
	filterArgs := filters.NewArgs()
	if dangling {
		filterArgs.Add("dangling", "true")
	}

	// 清理镜像
	pruneReport, err := global.GVA_DOCKER.ImagesPrune(ctx, filterArgs)
	if err != nil {
		global.GVA_LOG.Error("Failed to prune images", zap.Error(err))
		return 0, 0, fmt.Errorf("failed to prune images: %v", err)
	}

	deletedCount := int64(len(pruneReport.ImagesDeleted))
	spaceReclaimed := int64(pruneReport.SpaceReclaimed)

	global.GVA_LOG.Info("Images pruned successfully", 
		zap.Int64("deletedCount", deletedCount), 
		zap.Int64("spaceReclaimed", spaceReclaimed))

	return deletedCount, spaceReclaimed, nil
}

// convertToImageInfo 将Docker API的ImageSummary转换为ImageInfo响应模型
func (d *DockerImageService) convertToImageInfo(dockerImage types.ImageSummary) response.ImageInfo {
	return response.ImageInfo{
		ID:          dockerImage.ID,
		RepoTags:    dockerImage.RepoTags,
		RepoDigests: dockerImage.RepoDigests,
		Size:        dockerImage.Size,
		VirtualSize: dockerImage.VirtualSize,
		Created:     dockerImage.Created,
		Labels:      dockerImage.Labels,
		Containers:  dockerImage.Containers,
	}
}

// convertToImageDetail 将Docker API的ImageInspect转换为ImageDetail响应模型
func (d *DockerImageService) convertToImageDetail(dockerImage types.ImageInspect) response.ImageDetail {
	// 转换基本信息
	imageInfo := response.ImageInfo{
		ID:          dockerImage.ID,
		RepoTags:    dockerImage.RepoTags,
		RepoDigests: dockerImage.RepoDigests,
		Size:        dockerImage.Size,
		VirtualSize: dockerImage.VirtualSize,
		Created:     time.Now().Unix(), // 临时使用当前时间，实际应该解析dockerImage.Created字符串
		Labels:      dockerImage.Config.Labels,
		Containers:  -1, // ImageInspect中没有这个字段
	}

	// 转换配置信息
	config := response.ImageConfig{
		Hostname:     dockerImage.Config.Hostname,
		Domainname:   dockerImage.Config.Domainname,
		User:         dockerImage.Config.User,
		Env:          dockerImage.Config.Env,
		Cmd:          dockerImage.Config.Cmd,
		Image:        dockerImage.Config.Image,
		WorkingDir:   dockerImage.Config.WorkingDir,
		Entrypoint:   dockerImage.Config.Entrypoint,
		Labels:       dockerImage.Config.Labels,
		ExposedPorts: convertExposedPorts(dockerImage.Config.ExposedPorts),
	}

	// 转换根文件系统信息
	rootFS := response.RootFS{
		Type:   dockerImage.RootFS.Type,
		Layers: dockerImage.RootFS.Layers,
	}

	// 转换历史记录 - ImageInspect中没有History字段，使用空切片
	history := make([]response.ImageHistory, 0)

	return response.ImageDetail{
		ImageInfo:    imageInfo,
		Architecture: dockerImage.Architecture,
		Os:           dockerImage.Os,
		Config:       config,
		RootFS:       rootFS,
		History:      history,
	}
}

// BuildImage 构建镜像
func (d *DockerImageService) BuildImage(buildReq request.ImageBuildRequest) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return "", fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second) // 构建镜像可能需要很长时间
	defer cancel()

	// 构建镜像名称和标签
	imageName := buildReq.ImageName
	if buildReq.Tag != "" {
		imageName = fmt.Sprintf("%s:%s", buildReq.ImageName, buildReq.Tag)
	} else {
		imageName = fmt.Sprintf("%s:latest", buildReq.ImageName)
	}

	// 创建构建上下文（使用Dockerfile内容）
	dockerfileContent := buildReq.Dockerfile
	buildContext := strings.NewReader(dockerfileContent)

	// 转换构建参数
	buildArgs := make(map[string]*string)
	for key, value := range buildReq.BuildArgs {
		buildArgs[key] = &value
	}

	// 设置构建选项
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		BuildArgs:  buildArgs,
		Remove:     true, // 构建完成后删除中间容器
	}

	// 如果指定了构建上下文路径，需要创建tar包
	if buildReq.Context != "" {
		// 这里简化处理，实际应该创建包含Dockerfile和上下文文件的tar包
		global.GVA_LOG.Warn("Build context path specified but not fully implemented", zap.String("context", buildReq.Context))
	}

	// 构建镜像
	buildResponse, err := global.GVA_DOCKER.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		global.GVA_LOG.Error("Failed to build image", zap.String("imageName", imageName), zap.Error(err))
		return "", fmt.Errorf("failed to build image: %v", err)
	}
	defer buildResponse.Body.Close()

	// 读取构建日志
	buildLog, err := io.ReadAll(buildResponse.Body)
	if err != nil {
		global.GVA_LOG.Error("Failed to read build log", zap.String("imageName", imageName), zap.Error(err))
		return "", fmt.Errorf("failed to read build log: %v", err)
	}

	global.GVA_LOG.Info("Image built successfully", zap.String("imageName", imageName))
	return string(buildLog), nil
}

// ExportImage 导出镜像
func (d *DockerImageService) ExportImage(exportReq request.ImageExportRequest) (io.ReadCloser, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return nil, fmt.Errorf("Docker client is not available")
	}

	if len(exportReq.Images) == 0 {
		return nil, fmt.Errorf("no images specified for export")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// 导出镜像
	reader, err := global.GVA_DOCKER.ImageSave(ctx, exportReq.Images)
	if err != nil {
		global.GVA_LOG.Error("Failed to export images", zap.Strings("images", exportReq.Images), zap.Error(err))
		return nil, fmt.Errorf("failed to export images: %v", err)
	}

	global.GVA_LOG.Info("Images exported successfully", zap.Strings("images", exportReq.Images))
	return reader, nil
}

// ImportImage 导入镜像
func (d *DockerImageService) ImportImage(importReq request.ImageImportRequest) (string, error) {
	// 检查Docker客户端是否可用
	if global.GVA_DOCKER == nil {
		return "", fmt.Errorf("Docker client is not available")
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	// 构建导入选项
	importOptions := types.ImageImportOptions{}
	if importReq.Tag != "" {
		importOptions.Tag = importReq.Tag
	}

	// 导入镜像 - 这里简化处理，实际应该根据source类型处理不同的导入源
	var source io.Reader
	if strings.HasPrefix(importReq.Source, "http://") || strings.HasPrefix(importReq.Source, "https://") {
		// URL导入 - 实际实现中需要下载文件
		return "", fmt.Errorf("URL import not implemented yet")
	} else {
		// 文件路径导入 - 实际实现中需要读取文件
		return "", fmt.Errorf("file import not implemented yet - source: %s", importReq.Source)
	}

	// 执行导入
	importResponse, err := global.GVA_DOCKER.ImageImport(ctx, types.ImageImportSource{Source: source}, "", importOptions)
	if err != nil {
		global.GVA_LOG.Error("Failed to import image", zap.String("source", importReq.Source), zap.Error(err))
		return "", fmt.Errorf("failed to import image: %v", err)
	}
	defer importResponse.Close()

	// 读取导入日志
	importLog, err := io.ReadAll(importResponse)
	if err != nil {
		global.GVA_LOG.Error("Failed to read import log", zap.String("source", importReq.Source), zap.Error(err))
		return "", fmt.Errorf("failed to read import log: %v", err)
	}

	global.GVA_LOG.Info("Image imported successfully", zap.String("source", importReq.Source))
	return string(importLog), nil
}

// convertExposedPorts 转换暴露端口
func convertExposedPorts(exposedPorts nat.PortSet) map[string]struct{} {
	result := make(map[string]struct{})
	for port := range exposedPorts {
		result[string(port)] = struct{}{}
	}
	return result
}