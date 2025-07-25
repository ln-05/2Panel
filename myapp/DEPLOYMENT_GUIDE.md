# 🚀 MyApp 远程服务器部署指南

## 📦 部署包内容

```
myapp/
├── Dockerfile              # Docker镜像构建文件
├── docker-compose.yml      # Docker Compose配置
├── nginx.conf              # Nginx配置文件
├── init.sql                # 数据库初始化脚本
├── deploy.sh               # Linux部署脚本
├── deploy.bat              # Windows部署脚本
├── README_DOCKER.md        # 详细文档
├── .dockerignore           # Docker忽略文件
├── server/                 # gRPC服务器代码
├── api/                    # HTTP API服务器代码
└── DEPLOYMENT_GUIDE.md     # 本文件
```

## 🎯 部署步骤

### 步骤1: 上传文件到服务器

```bash
# 方法1: 使用scp上传整个目录
scp -r myapp/ root@14.103.149.197:/opt/

# 方法2: 使用rsync同步
rsync -avz myapp/ root@14.103.149.197:/opt/myapp/

# 方法3: 使用FTP/SFTP工具
# 推荐使用FileZilla, WinSCP等工具上传
```

### 步骤2: 登录服务器并部署

```bash
# 1. SSH登录服务器
ssh root@14.103.149.197

# 2. 进入项目目录
cd /opt/myapp

# 3. 给脚本执行权限
chmod +x deploy.sh

# 4. 启动服务
./deploy.sh start
```

### 步骤3: 验证部署

```bash
# 检查服务状态
./deploy.sh status

# 查看访问地址
./deploy.sh urls

# 测试API接口
curl http://14.103.149.197/api/database/list
```

## 🔧 手动部署（如果脚本失败）

```bash
# 1. 安装Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 2. 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 3. 启动服务
docker-compose up -d --build

# 4. 查看状态
docker-compose ps
```

## 📋 服务端口说明

- **80**: Web界面和API访问端口（对外开放）
- **3306**: MySQL数据库端口（可选开放）
- **7777**: gRPC服务端口（内部使用）
- **8889**: HTTP API服务端口（内部使用）

## 🌐 访问地址

部署成功后，可以通过以下地址访问：

- **主页**: http://14.103.149.197
- **API测试**: http://14.103.149.197/api/database/list
- **健康检查**: http://14.103.149.197/health

## 🔍 故障排除

### 问题1: 端口被占用

```bash
# 检查端口占用
netstat -tlnp | grep :80

# 停止占用端口的服务
sudo systemctl stop nginx  # 如果是nginx
sudo systemctl stop apache2  # 如果是apache

# 或者修改docker-compose.yml中的端口映射
ports:
  - "8080:80"  # 改为8080端口
```

### 问题2: Docker服务未启动

```bash
# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 检查Docker状态
sudo systemctl status docker
```

### 问题3: 权限问题

```bash
# 将用户添加到docker组
sudo usermod -aG docker $USER

# 重新登录或执行
newgrp docker
```

### 问题4: 内存不足

```bash
# 检查系统资源
free -h
df -h

# 清理Docker资源
docker system prune -a
```

## 📊 监控和维护

### 查看日志

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

### 备份数据

```bash
# 使用脚本备份
./deploy.sh backup

# 手动备��
docker-compose exec mysql mysqldump -u myapp -pmyapp123 myapp > backup.sql
```

### 更新服务

```bash
# 重新构建并启动
docker-compose up -d --build

# 或使用脚本
./deploy.sh restart
```

## 🔒 安全建议

1. **修改默认密码**
   - 编辑 `docker-compose.yml` 中的数据库密码
   - 重新部署服务

2. **配置防火墙**
   ```bash
   # 只开放必要端口
   sudo ufw allow 80
   sudo ufw allow 443
   sudo ufw enable
   ```

3. **使用HTTPS**
   - 获取SSL证书
   - 修改nginx配置启用HTTPS

## 📞 获取帮助

如果遇到问题：

1. 查看日志文件定位问题
2. 检查服务状态和端口占用
3. 确认Docker和Docker Compose版本
4. 验证网络连接和防火墙设置

常用调试命令：
```bash
# 检查Docker状态
docker info

# 检查容器状态
docker ps -a

# 检查网络
docker network ls

# 检查镜像
docker images

# 进入容器调试
docker-compose exec myapp sh
```

## 🎉 部署完成

部署成功后，你就可以：

1. 通过Web界面管理数据库连接
2. 使用API接口进行数据库操作
3. 从其他服务器同步数据库配置

现在你的数据库同步功能应该可以正常工作了！