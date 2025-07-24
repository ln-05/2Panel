# Ollama模型管理系统 - 问题解决指南

## 🔧 常见问题解决

### 1. 编译错误

#### 问题：导入模块找不到
```
Cannot resolve module '@/config/ollama'
```

**解决方案**：
```bash
# 确保配置文件存在
ls web/src/config/ollama.js

# 如果不存在，重新创建
# 参考 web/src/config/ollama.js 文件内容
```

#### 问题：API函数未定义
```
'chatWithOllamaModel' is not defined
```

**解决方案**：
```javascript
// 确保在组件中正确导入
import { chatWithOllamaModel } from '@/api/AI/ollamaModel'
```

### 2. 运行时错误

#### 问题：WebSocket连接失败
**症状**：连接状态显示"连接错误"

**解决方案**：
1. 检查后端服务是否运行
2. 确认WebSocket端点配置正确
3. 检查防火墙设置

#### 问题：API调用失败
**症状**：网络请求返回404或500错误

**解决方案**：
1. 确认后端路由配置正确
2. 检查API接口是否正确注册
3. 验证请求参数格式

### 3. 功能问题

#### 问题：模型下载失败
**解决方案**：
1. 检查Ollama服务状态
2. 确认网络连接正常
3. 验证磁盘空间充足

#### 问题：对话无响应
**解决方案**：
1. 确认模型状态为"运行中"
2. 检查模型是否正确启动
3. 查看浏览器控制台错误

## 🛠️ 调试步骤

### 1. 前端调试
```bash
# 启动开发服务器
cd web
npm run dev

# 检查控制台错误
# 打开浏览器开发者工具 (F12)
# 查看 Console 和 Network 标签页
```

### 2. 后端调试
```bash
# 启动后端服务
cd server
go run main.go

# 检查日志输出
# 查看终端输出的错误信息
```

### 3. Ollama服务调试
```bash
# 检查Ollama服务状态
ollama list

# 测试Ollama API
curl http://localhost:11434/api/tags

# 重启Ollama服务
# 根据你的系统重启Ollama服务
```

## 🔍 日志检查

### 前端日志
- 浏览器控制台 (F12 -> Console)
- 网络请求 (F12 -> Network)
- Vue DevTools (如果安装)

### 后端日志
- 终端输出
- 日志文件 (如果配置)
- 数据库日志

### Ollama日志
- Ollama服务日志
- 模型运行日志

## 📞 获取帮助

如果问题仍然存在：

1. **收集信息**：
   - 错误消息截图
   - 浏览器控制台日志
   - 后端服务日志
   - 系统环境信息

2. **检查文档**：
   - README.md
   - QUICK_START.md
   - DEPLOYMENT_CHECKLIST.md

3. **重新部署**：
   - 按照部署检查清单重新部署
   - 确保所有依赖正确安装

## ✅ 验证修复

修复问题后，请按以下步骤验证：

1. **基础功能测试**：
   - 访问 `/ollama-test` 页面
   - 测试API和WebSocket连接

2. **完整流程测试**：
   - 添加模型
   - 启动模型
   - 进行对话

3. **稳定性测试**：
   - 长时间运行测试
   - 多用户并发测试