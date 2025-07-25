#!/bin/bash

# MyApp Docker部署脚本
# 使用方法: ./deploy.sh [start|stop|restart|logs|status]

set -e

PROJECT_NAME="myapp"
COMPOSE_FILE="docker-compose.yml"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker和Docker Compose
check_requirements() {
    log_info "检查系统要求..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    log_success "系统要求检查通过"
}

# 启动服务
start_services() {
    log_info "启动MyApp服务..."
    
    # 构建并启动服务
    docker-compose -f $COMPOSE_FILE up -d --build
    
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    if docker-compose -f $COMPOSE_FILE ps | grep -q "Up"; then
        log_success "服务启动成功！"
        show_status
        show_urls
    else
        log_error "服务启动失败，请检查日志"
        docker-compose -f $COMPOSE_FILE logs
        exit 1
    fi
}

# 停止服务
stop_services() {
    log_info "停止MyApp服务..."
    docker-compose -f $COMPOSE_FILE down
    log_success "服务已停止"
}

# 重启服务
restart_services() {
    log_info "重启MyApp服务..."
    stop_services
    start_services
}

# 显示日志
show_logs() {
    log_info "显示服务日志..."
    docker-compose -f $COMPOSE_FILE logs -f
}

# 显示状态
show_status() {
    log_info "服务状态:"
    docker-compose -f $COMPOSE_FILE ps
    
    echo ""
    log_info "容器资源使用情况:"
    docker stats --no-stream $(docker-compose -f $COMPOSE_FILE ps -q) 2>/dev/null || true
}

# 显示访问URL
show_urls() {
    echo ""
    log_success "=== 服务访问地址 ==="
    echo "🌐 Web界面: http://localhost"
    echo "📡 API服务: http://localhost/api/database/list"
    echo "🔍 健康检查: http://localhost/health"
    echo "📊 数据库: localhost:3306 (用户: myapp, 密码: myapp123)"
    echo ""
    log_info "如果是远程服务器，请将localhost替换为服务器IP地址"
}

# 清理资源
cleanup() {
    log_info "清理Docker资源..."
    docker-compose -f $COMPOSE_FILE down -v --rmi all
    docker system prune -f
    log_success "清理完成"
}

# 备份数据
backup_data() {
    log_info "备份数据库数据..."
    BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p $BACKUP_DIR
    
    docker-compose -f $COMPOSE_FILE exec mysql mysqldump -u myapp -pmyapp123 myapp > $BACKUP_DIR/myapp_backup.sql
    log_success "数据备份完成: $BACKUP_DIR/myapp_backup.sql"
}

# 主函数
main() {
    case "${1:-start}" in
        start)
            check_requirements
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            check_requirements
            restart_services
            ;;
        logs)
            show_logs
            ;;
        status)
            show_status
            ;;
        urls)
            show_urls
            ;;
        cleanup)
            cleanup
            ;;
        backup)
            backup_data
            ;;
        *)
            echo "使用方法: $0 {start|stop|restart|logs|status|urls|cleanup|backup}"
            echo ""
            echo "命令说明:"
            echo "  start   - 启动所有服务"
            echo "  stop    - 停止所有服务"
            echo "  restart - 重启所有服务"
            echo "  logs    - 查看服务日志"
            echo "  status  - 查看服务状态"
            echo "  urls    - 显示访问地址"
            echo "  cleanup - 清理所有资源"
            echo "  backup  - 备份数据库"
            exit 1
            ;;
    esac
}

main "$@"