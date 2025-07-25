package docker

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockerModel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	dockerReq "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
	dockerRes "github.com/flipped-aurora/gin-vue-admin/server/model/docker/response"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DockerRegistryService struct{}

// GetRegistryList 获取仓库列表
func (d *DockerRegistryService) GetRegistryList(filter dockerReq.RegistryFilter) ([]dockerRes.RegistryInfo, int64, error) {
	var registries []dockerModel.DockerRegistry
	var total int64

	db := global.GVA_DB.Model(&dockerModel.DockerRegistry{})

	// 应用过滤条件
	if filter.Name != "" {
		db = db.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Protocol != "" {
		db = db.Where("protocol = ?", filter.Protocol)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		global.GVA_LOG.Error("获取仓库总数失败", zap.Error(err))
		return nil, 0, err
	}

	// 应用分页
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		db = db.Offset(offset).Limit(filter.PageSize)
	}

	// 获取数据
	if err := db.Order("created_at DESC").Find(&registries).Error; err != nil {
		global.GVA_LOG.Error("获取仓库列表失败", zap.Error(err))
		return nil, 0, err
	}

	// 转换为响应模型
	registryInfos := make([]dockerRes.RegistryInfo, 0, len(registries))
	for _, registry := range registries {
		registryInfos = append(registryInfos, d.convertToRegistryInfo(registry))
	}

	return registryInfos, total, nil
}

// GetRegistryDetail 获取仓库详细信息
func (d *DockerRegistryService) GetRegistryDetail(id uint) (*dockerRes.RegistryDetail, error) {
	var registry dockerModel.DockerRegistry
	if err := global.GVA_DB.First(&registry, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("仓库不存在")
		}
		global.GVA_LOG.Error("获取仓库详情失败", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	registryDetail := d.convertToRegistryDetail(registry)
	return &registryDetail, nil
}

// CreateRegistry 创建仓库
func (d *DockerRegistryService) CreateRegistry(createReq dockerReq.RegistryCreateRequest) (*dockerRes.RegistryInfo, error) {
	// 检查名称是否已存在
	var existingRegistry dockerModel.DockerRegistry
	if err := global.GVA_DB.Where("name = ?", createReq.Name).First(&existingRegistry).Error; err == nil {
		return nil, fmt.Errorf("仓库名称已存在")
	}

	// 验证下载地址格式
	if !strings.HasPrefix(createReq.DownloadUrl, "http://") && !strings.HasPrefix(createReq.DownloadUrl, "https://") {
		return nil, fmt.Errorf("下载地址格式不正确，必须以http://或https://开头")
	}

	// 创建仓库记录
	registry := dockerModel.DockerRegistry{
		Name:        createReq.Name,
		DownloadUrl: createReq.DownloadUrl,
		Protocol:    createReq.Protocol,
		Username:    createReq.Username,
		Password:    createReq.Password, // 实际项目中应该加密存储
		Description: createReq.Description,
		Status:      "active",
	}

	if err := global.GVA_DB.Create(&registry).Error; err != nil {
		global.GVA_LOG.Error("创建仓库失败", zap.String("name", createReq.Name), zap.Error(err))
		return nil, fmt.Errorf("创建仓库失败: %v", err)
	}

	global.GVA_LOG.Info("仓库创建成功", zap.String("name", createReq.Name), zap.Uint("id", registry.ID))

	registryInfo := d.convertToRegistryInfo(registry)
	return &registryInfo, nil
}

// UpdateRegistry 更新仓库
func (d *DockerRegistryService) UpdateRegistry(updateReq dockerReq.RegistryUpdateRequest) (*dockerRes.RegistryInfo, error) {
	var registry dockerModel.DockerRegistry
	if err := global.GVA_DB.First(&registry, updateReq.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("仓库不存在")
		}
		return nil, err
	}

	// 检查名称是否与其他仓库冲突
	var existingRegistry dockerModel.DockerRegistry
	if err := global.GVA_DB.Where("name = ? AND id != ?", updateReq.Name, updateReq.ID).First(&existingRegistry).Error; err == nil {
		return nil, fmt.Errorf("仓库名称已存在")
	}

	// 验证下载地址格式
	if !strings.HasPrefix(updateReq.DownloadUrl, "http://") && !strings.HasPrefix(updateReq.DownloadUrl, "https://") {
		return nil, fmt.Errorf("下载地址格式不正确，必须以http://或https://开头")
	}

	// 更新仓库信息
	registry.Name = updateReq.Name
	registry.DownloadUrl = updateReq.DownloadUrl
	registry.Protocol = updateReq.Protocol
	registry.Username = updateReq.Username
	registry.Password = updateReq.Password // 实际项目中应该加密存储
	registry.Description = updateReq.Description

	if err := global.GVA_DB.Save(&registry).Error; err != nil {
		global.GVA_LOG.Error("更新仓库失败", zap.Uint("id", updateReq.ID), zap.Error(err))
		return nil, fmt.Errorf("更新仓库失败: %v", err)
	}

	global.GVA_LOG.Info("仓库更新成功", zap.Uint("id", updateReq.ID))

	registryInfo := d.convertToRegistryInfo(registry)
	return &registryInfo, nil
}

// DeleteRegistry 删除仓库
func (d *DockerRegistryService) DeleteRegistry(id uint) error {
	var registry dockerModel.DockerRegistry
	if err := global.GVA_DB.First(&registry, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("仓库不存在")
		}
		return err
	}

	// 检查是否为默认仓库
	if registry.IsDefault {
		return fmt.Errorf("不能删除默认仓库")
	}

	if err := global.GVA_DB.Delete(&registry).Error; err != nil {
		global.GVA_LOG.Error("删除仓库失败", zap.Uint("id", id), zap.Error(err))
		return fmt.Errorf("删除仓库失败: %v", err)
	}

	global.GVA_LOG.Info("仓库删除成功", zap.Uint("id", id))
	return nil
}

// TestRegistry 测试仓库连接
func (d *DockerRegistryService) TestRegistry(id uint) (*dockerRes.RegistryTestResponse, error) {
	var registry dockerModel.DockerRegistry
	if err := global.GVA_DB.First(&registry, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("仓库不存在")
		}
		return nil, err
	}

	// 测试连接
	success, message := d.testRegistryConnection(registry)

	// 更新测试结果
	now := time.Now()
	registry.LastTestTime = &now
	registry.TestResult = message
	if success {
		registry.Status = "active"
	} else {
		registry.Status = "inactive"
	}

	global.GVA_DB.Save(&registry)

	return &dockerRes.RegistryTestResponse{
		Success: success,
		Message: message,
	}, nil
}

// SetDefaultRegistry 设置默认仓库
func (d *DockerRegistryService) SetDefaultRegistry(id uint) error {
	// 先取消所有默认仓库
	if err := global.GVA_DB.Model(&dockerModel.DockerRegistry{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		return err
	}

	// 设置新的默认仓库
	if err := global.GVA_DB.Model(&dockerModel.DockerRegistry{}).Where("id = ?", id).Update("is_default", true).Error; err != nil {
		return err
	}

	global.GVA_LOG.Info("设置默认仓库成功", zap.Uint("id", id))
	return nil
}

// testRegistryConnection 测试仓库连接
func (d *DockerRegistryService) testRegistryConnection(registry dockerModel.DockerRegistry) (bool, string) {
	// 构建测试URL
	testUrl := registry.DownloadUrl
	if !strings.HasSuffix(testUrl, "/") {
		testUrl += "/"
	}
	testUrl += "v2/"

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 跳过SSL验证，实际项目中应该根据需要配置
			},
		},
	}

	// 创建请求
	req, err := http.NewRequestWithContext(context.Background(), "GET", testUrl, nil)
	if err != nil {
		return false, fmt.Sprintf("创建请求失败: %v", err)
	}

	// 添加认证信息
	if registry.Username != "" && registry.Password != "" {
		req.SetBasicAuth(registry.Username, registry.Password)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Sprintf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == 200 || resp.StatusCode == 401 {
		// 200表示成功，401表示需要认证但仓库存在
		return true, "连接成功"
	}

	return false, fmt.Sprintf("连接失败，状态码: %d", resp.StatusCode)
}

// convertToRegistryInfo 转换为RegistryInfo
func (d *DockerRegistryService) convertToRegistryInfo(registry dockerModel.DockerRegistry) dockerRes.RegistryInfo {
	return dockerRes.RegistryInfo{
		ID:          registry.ID,
		Name:        registry.Name,
		DownloadUrl: registry.DownloadUrl,
		Protocol:    registry.Protocol,
		Status:      registry.Status,
		Username:    registry.Username,
		Description: registry.Description,
		CreatedAt:   registry.CreatedAt,
		UpdatedAt:   registry.UpdatedAt,
	}
}

// convertToRegistryDetail 转换为RegistryDetail
func (d *DockerRegistryService) convertToRegistryDetail(registry dockerModel.DockerRegistry) dockerRes.RegistryDetail {
	return dockerRes.RegistryDetail{
		RegistryInfo: d.convertToRegistryInfo(registry),
		IsDefault:    registry.IsDefault,
		LastTestTime: registry.LastTestTime,
		TestResult:   registry.TestResult,
	}
}