package AI

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/AI"
	AIReq "github.com/flipped-aurora/gin-vue-admin/server/model/AI/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OllamaModelService struct {
	resourceChecker *ResourceChecker
}

// NewOllamaModelService 创建OllamaModelService实例
func NewOllamaModelService() *OllamaModelService {
	return &OllamaModelService{
		resourceChecker: NewResourceChecker(),
	}
}

const (
	StatusStopped     = "stopped"
	StatusRunning     = "running"
	StatusDownloading = "downloading"
	StatusError       = "error"
	StatusUnavailable = "unavailable"
)

// CreateOllamaModel 创建ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) CreateOllamaModel(ctx context.Context, ollamaModel *AI.OllamaModel) (err error) {
	err = global.GVA_DB.Create(ollamaModel).Error
	return err
}

// DeleteOllamaModel 删除ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) DeleteOllamaModel(ctx context.Context, ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AI.OllamaModel{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&AI.OllamaModel{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteOllamaModelByIds 批量删除ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) DeleteOllamaModelByIds(ctx context.Context, IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AI.OllamaModel{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&AI.OllamaModel{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateOllamaModel 更新ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) UpdateOllamaModel(ctx context.Context, ollamaModel AI.OllamaModel) (err error) {
	err = global.GVA_DB.Model(&AI.OllamaModel{}).Where("id = ?", ollamaModel.ID).Updates(&ollamaModel).Error
	return err
}

// GetOllamaModel 根据ID获取ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) GetOllamaModel(ctx context.Context, ID string) (ollamaModel AI.OllamaModel, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&ollamaModel).Error
	return
}

// GetOllamaModelInfoList 分页获取ollamaModel表记录
// Author [yourname](https://github.com/yourname)
func (ollamaModelService *OllamaModelService) GetOllamaModelInfoList(ctx context.Context, info AIReq.OllamaModelSearch) (list []AI.OllamaModel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&AI.OllamaModel{})
	var ollamaModels []AI.OllamaModel
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&ollamaModels).Error
	return ollamaModels, total, err
}
func (ollamaModelService *OllamaModelService) GetOllamaModelPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// Search 搜索和分页查询模型
func (ollamaModelService *OllamaModelService) Search(ctx context.Context, info AIReq.OllamaModelSearch) (list []AI.OllamaModelInfo, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&AI.OllamaModel{})
	var ollamaModels []AI.OllamaModel

	// 条件搜索
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if len(info.Status) > 0 {
		db = db.Where("status IN ?", info.Status)
	}
	if len(info.From) > 0 {
		db = db.Where("from IN ?", info.From)
	}
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	// 大小范围筛选 (需要解析大小字符串)
	if info.MinSize != "" || info.MaxSize != "" {
		db = ollamaModelService.applySizeFilter(db, info.MinSize, info.MaxSize)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 排序
	if info.SortBy != "" {
		orderBy := info.SortBy
		if info.SortDesc {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}
		db = db.Order(orderBy)
	} else {
		db = db.Order("created_at DESC") // 默认按创建时间降序
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&ollamaModels).Error
	if err != nil {
		return
	}

	// 转换为DTO
	list = make([]AI.OllamaModelInfo, len(ollamaModels))
	for i, model := range ollamaModels {
		list[i] = model.ToInfo()
	}

	return list, total, err
}

// Create 创建/下载新模型
func (ollamaModelService *OllamaModelService) Create(ctx context.Context, req AIReq.OllamaModelCreate, userID uint) error {
	// 检查模型是否已存在
	var existingModel AI.OllamaModel
	err := global.GVA_DB.Where("name = ?", req.Name).First(&existingModel).Error
	if err == nil {
		return fmt.Errorf("模型 %s 已存在", req.Name)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 检查系统资源
	resourceStatus, err := ollamaModelService.resourceChecker.CheckResources()
	if err != nil {
		global.GVA_LOG.Error("检查系统资源失败", zap.Error(err))
		return fmt.Errorf("检查系统资源失败: %v", err)
	}

	if !resourceStatus.IsResourceEnough {
		global.GVA_LOG.Warn("系统资源不足", zap.String("message", resourceStatus.Message))
		return fmt.Errorf("系统资源不足: %s", resourceStatus.Message)
	}

	// 预估模型大小并检查磁盘空间
	estimatedSize := ollamaModelService.resourceChecker.EstimateModelSize(req.Name)
	if resourceStatus.DiskSpaceGB < estimatedSize {
		suggestions := ollamaModelService.resourceChecker.GetCleanupSuggestions()
		return fmt.Errorf("磁盘空间不足，预估需要%.1fGB，当前可用%.1fGB。建议清理: %v",
			estimatedSize, resourceStatus.DiskSpaceGB, suggestions)
	}

	// 创建模型记录
	model := AI.OllamaModel{
		Name:      req.Name,
		From:      req.From,
		Status:    StatusDownloading,
		Message:   "开始下载模型...",
		CreatedBy: userID,
	}

	err = global.GVA_DB.Create(&model).Error
	if err != nil {
		return err
	}

	// 异步下载模型
	go ollamaModelService.downloadModel(ctx, &model)

	return nil
}

// downloadModel 下载模型的异步方法
func (ollamaModelService *OllamaModelService) downloadModel(ctx context.Context, model *AI.OllamaModel) {
	global.GVA_LOG.Info("开始下载模型", zap.String("name", model.Name))

	// 更新状态为下载中
	err := ollamaModelService.updateModelStatus(model.ID, StatusDownloading, "正在下载模型...")
	if err != nil {
		return
	}

	// 执行ollama pull命令
	cmd := exec.Command("ollama", "pull", model.Name)
	output, err := cmd.CombinedOutput()

	if err != nil {
		global.GVA_LOG.Error("模型下载失败", zap.String("name", model.Name), zap.Error(err))
		ollamaModelService.updateModelStatus(model.ID, StatusError, fmt.Sprintf("下载失败: %s", string(output)))
		return
	}

	// 获取模型信息
	modelInfo, err := ollamaModelService.getOllamaModelInfo(model.Name)
	if err != nil {
		global.GVA_LOG.Error("获取模型信息失败", zap.String("name", model.Name), zap.Error(err))
		ollamaModelService.updateModelStatus(model.ID, StatusError, "获取模型信息失败")
		return
	}

	// 更新模型信息
	updates := map[string]interface{}{
		"status":  StatusStopped,
		"message": "模型下载完成",
		"size":    modelInfo.Size,
	}

	err = global.GVA_DB.Model(&AI.OllamaModel{}).Where("id = ?", model.ID).Updates(updates).Error
	if err != nil {
		global.GVA_LOG.Error("更新模型信息失败", zap.Error(err))
	}

	global.GVA_LOG.Info("模型下载完成", zap.String("name", model.Name))
}

// Start 启动模型
func (ollamaModelService *OllamaModelService) Start(ctx context.Context, ID string, userID uint) error {
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", ID).First(&model).Error
	if err != nil {
		return err
	}

	// 执行启动命令
	cmd := exec.Command("ollama", "run", model.Name)
	err = cmd.Start()
	if err != nil {
		global.GVA_LOG.Error("启动模型失败", zap.String("name", model.Name), zap.Error(err))
		return fmt.Errorf("启动模型失败: %v", err)
	}

	// 更新状态
	err = ollamaModelService.updateModelStatus(model.ID, StatusRunning, "模型已启动")
	if err != nil {
		return err
	}

	global.GVA_LOG.Info("模型已启动", zap.String("name", model.Name))
	return nil
}

// Stop 停止模型
func (ollamaModelService *OllamaModelService) Stop(ctx context.Context, ID string, userID uint) error {
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", ID).First(&model).Error
	if err != nil {
		return err
	}

	// 执行停止命令
	cmd := exec.Command("ollama", "stop", model.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		global.GVA_LOG.Error("停止模型失败", zap.String("name", model.Name), zap.Error(err))
		return fmt.Errorf("停止模型失败: %s", string(output))
	}

	// 更新状态
	err = ollamaModelService.updateModelStatus(model.ID, StatusStopped, "模型已停止")
	if err != nil {
		return err
	}

	global.GVA_LOG.Info("模型已停止", zap.String("name", model.Name))
	return nil
}

// Close 停止模型 (兼容旧接口)
func (ollamaModelService *OllamaModelService) Close(ctx context.Context, ID string, userID uint) error {
	return ollamaModelService.Stop(ctx, ID, userID)
}

// Recreate 重新创建模型
func (ollamaModelService *OllamaModelService) Recreate(ctx context.Context, ID string, userID uint) error {
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", ID).First(&model).Error
	if err != nil {
		return err
	}

	// 先停止模型
	ollamaModelService.Close(ctx, ID, userID)

	// 删除本地模型
	cmd := exec.Command("ollama", "rm", model.Name)
	cmd.Run() // 忽略错误，可能模型不存在

	// 重新下载
	ollamaModelService.updateModelStatus(model.ID, StatusDownloading, "重新创建模型...")
	go ollamaModelService.downloadModel(ctx, &model)

	return nil
}

// Sync 同步本地和远程模型
func (ollamaModelService *OllamaModelService) Sync(ctx context.Context, req AIReq.OllamaModelSync) (map[string]int, error) {
	result := map[string]int{
		"added":   0,
		"updated": 0,
		"removed": 0,
	}

	// 获取本地ollama模型列表
	localModels, err := ollamaModelService.getLocalOllamaModels()
	if err != nil {
		return result, err
	}

	// 获取数据库中的模型列表
	var dbModels []AI.OllamaModel
	err = global.GVA_DB.Find(&dbModels).Error
	if err != nil {
		return result, err
	}

	// 创建映射便于查找
	dbModelMap := make(map[string]*AI.OllamaModel)
	for i := range dbModels {
		dbModelMap[dbModels[i].Name] = &dbModels[i]
	}

	localModelMap := make(map[string]bool)

	// 处理本地模型
	for _, localModel := range localModels {
		localModelMap[localModel.Name] = true

		if dbModel, exists := dbModelMap[localModel.Name]; exists {
			// 更新现有模型
			if dbModel.Size != localModel.Size || dbModel.Status == StatusUnavailable {
				updates := map[string]interface{}{
					"size":    localModel.Size,
					"status":  StatusStopped,
					"message": "模型已同步",
				}
				global.GVA_DB.Model(&AI.OllamaModel{}).Where("id = ?", dbModel.ID).Updates(updates)
				result["updated"]++
			}
		} else {
			// 添加新模型
			newModel := AI.OllamaModel{
				Name:    localModel.Name,
				Size:    localModel.Size,
				Status:  StatusStopped,
				Message: "通过同步添加",
			}
			global.GVA_DB.Create(&newModel)
			result["added"]++
		}
	}

	// 处理数据库中但本地不存在的模型
	for _, dbModel := range dbModels {
		if !localModelMap[dbModel.Name] {
			if req.Force {
				// 强制删除
				global.GVA_DB.Delete(&dbModel)
				result["removed"]++
			} else {
				// 标记为不可用
				global.GVA_DB.Model(&AI.OllamaModel{}).Where("id = ?", dbModel.ID).Updates(map[string]interface{}{
					"status":  StatusUnavailable,
					"message": "本地模型不存在",
				})
				result["updated"]++
			}
		}
	}

	return result, nil
}

// LoadDetail 加载模型详情
func (ollamaModelService *OllamaModelService) LoadDetail(ctx context.Context, ID string) (AI.OllamaModelInfo, error) {
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", ID).First(&model).Error
	if err != nil {
		return AI.OllamaModelInfo{}, err
	}

	// 检查日志文件是否存在
	logPath := filepath.Join("/var/log/ollama", model.Name+".log")
	if _, err := os.Stat(logPath); err == nil {
		model.LogFileExist = true
		global.GVA_DB.Model(&model).Update("log_file_exist", true)
	}

	return model.ToInfo(), nil
}

// updateModelStatus 更新模型状态
func (ollamaModelService *OllamaModelService) updateModelStatus(modelID uint, status, message string) error {
	updates := map[string]interface{}{
		"status":     status,
		"message":    message,
		"updated_at": time.Now(),
	}
	return global.GVA_DB.Model(&AI.OllamaModel{}).Where("id = ?", modelID).Updates(updates).Error
}

// getOllamaModelInfo 获取ollama模型信息
func (ollamaModelService *OllamaModelService) getOllamaModelInfo(modelName string) (*AI.OllamaModelInfo, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		} `json:"models"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	for _, model := range result.Models {
		if strings.HasPrefix(model.Name, modelName) {
			return &AI.OllamaModelInfo{
				Name: model.Name,
				Size: fmt.Sprintf("%.2f GB", float64(model.Size)/(1024*1024*1024)),
			}, nil
		}
	}

	return nil, fmt.Errorf("模型 %s 未找到", modelName)
}

// getLocalOllamaModels 获取本地ollama模型列表
func (ollamaModelService *OllamaModelService) getLocalOllamaModels() ([]AI.OllamaModelInfo, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		} `json:"models"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var models []AI.OllamaModelInfo
	for _, model := range result.Models {
		models = append(models, AI.OllamaModelInfo{
			Name: model.Name,
			Size: fmt.Sprintf("%.2f GB", float64(model.Size)/(1024*1024*1024)),
		})
	}

	return models, nil
}

// BindDomain 绑定域名到AI服务
func (ollamaModelService *OllamaModelService) BindDomain(ctx context.Context, req AIReq.OllamaBindDomainRequest, userID uint) error {
	// 检查域名是否已绑定
	var existingBind AI.OllamaBindDomain
	err := global.GVA_DB.Where("domain = ?", req.Domain).First(&existingBind).Error
	if err == nil {
		return fmt.Errorf("域名 %s 已被绑定", req.Domain)
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// 创建域名绑定记录
	bind := AI.OllamaBindDomain{
		Domain:       req.Domain,
		AppInstallID: req.AppInstallID,
		SSLID:        req.SSLID,
		WebsiteID:    req.WebsiteID,
		IPList:       req.IPList,
		CreatedBy:    userID,
	}

	return global.GVA_DB.Create(&bind).Error
}

// GetBindDomain 获取绑定的域名信息
func (ollamaModelService *OllamaModelService) GetBindDomain(ctx context.Context, domain string) (AI.OllamaBindDomain, error) {
	var bind AI.OllamaBindDomain
	err := global.GVA_DB.Where("domain = ?", domain).First(&bind).Error
	return bind, err
}

// UpdateBindDomain 更新域名绑定
func (ollamaModelService *OllamaModelService) UpdateBindDomain(ctx context.Context, domain string, req AIReq.OllamaBindDomainRequest, userID uint) error {
	updates := map[string]interface{}{
		"app_install_id": req.AppInstallID,
		"ssl_id":         req.SSLID,
		"website_id":     req.WebsiteID,
		"ip_list":        req.IPList,
		"updated_by":     userID,
		"updated_at":     time.Now(),
	}

	return global.GVA_DB.Model(&AI.OllamaBindDomain{}).Where("domain = ?", domain).Updates(updates).Error
}

// GetLogs 获取模型日志
func (ollamaModelService *OllamaModelService) GetLogs(ctx context.Context, ID string, lines int) (string, error) {
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", ID).First(&model).Error
	if err != nil {
		return "", err
	}

	// 构建日志文件路径
	logPath := filepath.Join("/var/log/ollama", model.Name+".log")

	// 检查日志文件是否存在
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return "日志文件不存在", nil
	}

	// 使用tail命令获取最后几行日志
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), logPath)
	output, err := cmd.Output()
	if err != nil {
		global.GVA_LOG.Error("读取日志失败", zap.String("path", logPath), zap.Error(err))
		return "", fmt.Errorf("读取日志失败: %v", err)
	}

	return string(output), nil
}

// GetSystemResourceStatus 获取系统资源状态
func (ollamaModelService *OllamaModelService) GetSystemResourceStatus(ctx context.Context) (*ResourceStatus, error) {
	return ollamaModelService.resourceChecker.CheckResources()
}

// Chat 与模型对话
func (ollamaModelService *OllamaModelService) Chat(ctx context.Context, req AIReq.OllamaChatRequest) (AIReq.OllamaChatResponse, error) {
	// 获取模型信息
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", req.ModelID).First(&model).Error
	if err != nil {
		return AIReq.OllamaChatResponse{}, fmt.Errorf("模型不存在: %v", err)
	}

	// 检查模型状态
	if model.Status != StatusRunning {
		return AIReq.OllamaChatResponse{}, fmt.Errorf("模型 %s 未运行，当前状态: %s", model.Name, model.Status)
	}

	// 构建请求体
	chatRequest := map[string]interface{}{
		"model":  model.Name,
		"prompt": req.Message,
		"stream": req.Stream,
	}

	// 如果有上下文，添加到请求中
	if req.Context != "" {
		chatRequest["context"] = req.Context
	}

	// 调用Ollama API
	response, err := ollamaModelService.callOllamaGenerate(chatRequest)
	if err != nil {
		global.GVA_LOG.Error("调用Ollama API失败", zap.String("model", model.Name), zap.Error(err))
		return AIReq.OllamaChatResponse{}, fmt.Errorf("对话失败: %v", err)
	}

	return response, nil
}

// callOllamaGenerate 调用Ollama的generate API
func (ollamaModelService *OllamaModelService) callOllamaGenerate(request map[string]interface{}) (AIReq.OllamaChatResponse, error) {
	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return AIReq.OllamaChatResponse{}, err
	}

	// 发送HTTP请求到Ollama API
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return AIReq.OllamaChatResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AIReq.OllamaChatResponse{}, fmt.Errorf("Ollama API错误 %d: %s", resp.StatusCode, string(body))
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AIReq.OllamaChatResponse{}, err
	}

	// 解析响应
	var ollamaResponse struct {
		Response string `json:"response"`
		Context  string `json:"context"`
		Done     bool   `json:"done"`
	}

	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return AIReq.OllamaChatResponse{}, err
	}

	return AIReq.OllamaChatResponse{
		Response: ollamaResponse.Response,
		Context:  ollamaResponse.Context,
		Done:     ollamaResponse.Done,
	}, nil
}

// ChatStream 流式对话（WebSocket支持）
func (ollamaModelService *OllamaModelService) ChatStream(ctx context.Context, req AIReq.OllamaChatRequest, responseChannel chan<- AIReq.OllamaChatResponse) error {
	// 获取模型信息
	var model AI.OllamaModel
	err := global.GVA_DB.Where("id = ?", req.ModelID).First(&model).Error
	if err != nil {
		return fmt.Errorf("模型不存在: %v", err)
	}

	// 检查模型状态
	if model.Status != StatusRunning {
		return fmt.Errorf("模型 %s 未运行，当前状态: %s", model.Name, model.Status)
	}

	// 构建流式请求
	chatRequest := map[string]interface{}{
		"model":  model.Name,
		"prompt": req.Message,
		"stream": true,
	}

	if req.Context != "" {
		chatRequest["context"] = req.Context
	}

	// 调用流式API
	return ollamaModelService.callOllamaGenerateStream(chatRequest, responseChannel)
}

// callOllamaGenerateStream 调用Ollama的流式generate API
func (ollamaModelService *OllamaModelService) callOllamaGenerateStream(request map[string]interface{}, responseChannel chan<- AIReq.OllamaChatResponse) error {
	defer close(responseChannel)

	// 序列化请求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// 发送HTTP请求
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Ollama API错误 %d: %s", resp.StatusCode, string(body))
	}

	// 逐行读取流式响应
	scanner := json.NewDecoder(resp.Body)
	for scanner.More() {
		var streamResponse struct {
			Response string `json:"response"`
			Context  string `json:"context"`
			Done     bool   `json:"done"`
		}

		err := scanner.Decode(&streamResponse)
		if err != nil {
			global.GVA_LOG.Error("解析流式响应失败", zap.Error(err))
			continue
		}

		// 发送到响应通道
		select {
		case responseChannel <- AIReq.OllamaChatResponse{
			Response: streamResponse.Response,
			Context:  streamResponse.Context,
			Done:     streamResponse.Done,
		}:
		case <-time.After(5 * time.Second):
			global.GVA_LOG.Warn("发送响应到通道超时")
			return fmt.Errorf("响应通道超时")
		}

		// 如果完成，退出循环
		if streamResponse.Done {
			break
		}
	}

	return nil
}

// applySizeFilter 应用大小范围筛选
func (ollamaModelService *OllamaModelService) applySizeFilter(db *gorm.DB, minSize, maxSize string) *gorm.DB {
	if minSize != "" {
		minBytes := ollamaModelService.parseSize(minSize)
		if minBytes > 0 {
			// 这里需要将数据库中的size字段转换为字节数进行比较
			// 由于size字段存储的是字符串格式，我们需要在应用层进行筛选
			// 或者考虑在数据库中添加size_bytes字段存储字节数
		}
	}
	if maxSize != "" {
		maxBytes := ollamaModelService.parseSize(maxSize)
		if maxBytes > 0 {
			// 同上
		}
	}
	return db
}

// parseSize 解析大小字符串为字节数
func (ollamaModelService *OllamaModelService) parseSize(sizeStr string) int64 {
	if sizeStr == "" {
		return 0
	}

	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	var multiplier int64 = 1
	var numStr string

	if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(sizeStr, "GB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		numStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "B") {
		multiplier = 1
		numStr = strings.TrimSuffix(sizeStr, "B")
	} else {
		// 默认为字节
		numStr = sizeStr
	}

	numStr = strings.TrimSpace(numStr)
	if num, err := strconv.ParseFloat(numStr, 64); err == nil {
		return int64(num * float64(multiplier))
	}

	return 0
}

// sortModels 对模型列表进行排序 (应用层排序，用于复杂排序逻辑)
func (ollamaModelService *OllamaModelService) sortModels(models []AI.OllamaModelInfo, sortBy string, sortDesc bool) {
	if sortBy == "" {
		return
	}

	sort.Slice(models, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "name":
			less = models[i].Name < models[j].Name
		case "size":
			sizeI := ollamaModelService.parseSize(models[i].Size)
			sizeJ := ollamaModelService.parseSize(models[j].Size)
			less = sizeI < sizeJ
		case "createdAt":
			less = models[i].CreatedAt.Before(models[j].CreatedAt)
		case "status":
			less = models[i].Status < models[j].Status
		default:
			return false
		}

		if sortDesc {
			return !less
		}
		return less
	})
}
