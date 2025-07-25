# MyApp Docker 部署指南

## 🚀 快速开始

### 1. 准备工作

确保你的服务器已安装：
- Docker (版本 20.10+)
- Docker Compose (版本 1.29+)

```bash
# 检查Docker版本
docker --version
docker-compose --version

# 如果未安装，在Ubuntu/Debian上安装：
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. 部署到远程服务器

#### 方法一：使用部署脚本（推荐）

```bash
# 1. 上传项目文件到服务器
scp -r myapp/ root@14.103.149.197:/opt/

# 2. 登录服务器
ssh root@14.103.149.197

# 3. 进入项目目录
cd /opt/myapp

# 4. 给部署脚本执行权限
chmod +x deploy.sh

# 5. 启动服务
./deploy.sh start
```

#### 方法二：手动部署

```bash
# 1. 登录服务器
ssh root@14.103.149.197

# 2. 进入项目目录
cd /opt/myapp

# 3. 构建并启动服务
docker-compose up -d --build

# 4. 查看服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f
```

### 3. 验证部署

部署成功后，你可以通过以下方式验证：

```bash
# 检查服务状态
curl http://14.103.149.197/health

# 测试API接口
curl http://14.103.149.197/api/database/list

# 查看Web界面
# 在浏览器中访问: http://14.103.149.197
```

## 📋 服务说明

### 服务组件

1. **MySQL数据库** (端口3306)
   - 用户: myapp
   - 密码: myapp123
   - 数据库: myapp

2. **MyApp gRPC服务** (端口7777)
   - 内部服务，处理业务逻辑

3. **MyApp API服务** (端口8889)
   - HTTP API接口服务

4. **Nginx反向代理** (端口80)
   - 对外提供Web服务和API代理

### 访问地址

- **Web界面**: http://14.103.149.197
- **API文档**: http://14.103.149.197/api/database/list
- **健康检查**: http://14.103.149.197/health

## 🔧 配置说明

### 环境变量

可以通过修改 `docker-compose.yml` 中的环境变量来配置服务：

```yaml
environment:
  - DB_HOST=mysql          # 数据库主机
  - DB_PORT=3306          # 数据库端口
  - DB_USER=myapp         # 数据库用户
  - DB_PASSWORD=myapp123  # 数据库密码
  - DB_NAME=myapp         # 数据库名称
```

### 端口配置

如果需要修改端口，编辑 `docker-compose.yml`：

```yaml
ports:
  - "80:80"     # Web端口
  - "3306:3306" # 数据库端口（可选暴露）
```

## 🛠️ 常用命令

### 部署脚本命令

```bash
./deploy.sh start    # 启动所有服务
./deploy.sh stop     # 停止所有服务
./deploy.sh restart  # 重启所有服务
./deploy.sh logs     # 查看服务日志
./deploy.sh status   # 查看服务状态
./deploy.sh urls     # 显示访问地址
./deploy.sh backup   # 备份数据库
./deploy.sh cleanup  # 清理所有资源
```

### Docker Compose命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f [service_name]

# 重新构建
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 进入容器
docker-compose exec myapp sh
docker-compose exec mysql mysql -u myapp -p
```

## 🔍 故障排除

### 常见问题

1. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs myapp
   
   # 检查端口占用
   netstat -tlnp | grep :80
   ```

2. **数据库连接失败**
   ```bash
   # 检查MySQL服务
   docker-compose logs mysql
   
   # 测试数据库连接
   docker-compose exec mysql mysql -u myapp -pmyapp123 -e "SELECT 1"
   ```

3. **API接口404错误**
   ```bash
   # 检查nginx配置
   docker-compose logs nginx
   
   # 重启nginx
   docker-compose restart nginx
   ```

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs myapp
docker-compose logs mysql
docker-compose logs nginx

# 实时查看日志
docker-compose logs -f
```

### 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看磁盘使用
docker system df

# 清理未使用的资源
docker system prune
```

## 🔒 安全配置

### 生产环境建议

1. **修改默认密码**
   ```yaml
   # 在docker-compose.yml中修改
   environment:
     MYSQL_ROOT_PASSWORD: your_secure_password
     MYSQL_PASSWORD: your_secure_password
   ```

2. **使用HTTPS**
   ```bash
   # 添加SSL证书到 ./ssl/ 目录
   # 修改nginx.conf启用HTTPS
   ```

3. **限制网络访问**
   ```yaml
   # 只暴露必要端口
   ports:
     - "80:80"
     # - "3306:3306"  # 不暴露数据库端口
   ```

## 📊 监控和维护

### 数据备份

```bash
# 自动备份
./deploy.sh backup

# 手动备份
docker-compose exec mysql mysqldump -u myapp -pmyapp123 myapp > backup_$(date +%Y%m%d).sql
```

### 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并部署
./deploy.sh restart
```

### 扩容配置

如果需要处理更多请求，可以：

1. **增加API服务实例**
   ```yaml
   myapp:
     deploy:
       replicas: 3
   ```

2. **配置负载均衡**
   ```nginx
   upstream myapp_api {
       server myapp_1:8889;
       server myapp_2:8889;
       server myapp_3:8889;
   }
   ```

## 📞 技术支持

如果遇到问题，请：

1. 查看日志文件
2. 检查服务状态
3. 验证网络连接
4. 确认配置正确

需要帮助时，请提供：
- 错误日志
- 服务状态
- 系统环境信息