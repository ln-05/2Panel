@echo off
REM MyApp Docker部署脚本 (Windows版本)
REM 使用方法: deploy.bat [start|stop|restart|logs|status]

setlocal enabledelayedexpansion

set PROJECT_NAME=myapp
set COMPOSE_FILE=docker-compose.yml

REM 检查参数
if "%1"=="" (
    set ACTION=start
) else (
    set ACTION=%1
)

echo [INFO] 执行操作: %ACTION%

REM 检查Docker和Docker Compose
:check_requirements
echo [INFO] 检查系统要求...

docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker未安装，请先安装Docker Desktop
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker Compose未安装，请先安装Docker Compose
    pause
    exit /b 1
)

echo [SUCCESS] 系统要求检查通过

REM 根据参数执行相应操作
if "%ACTION%"=="start" goto start_services
if "%ACTION%"=="stop" goto stop_services
if "%ACTION%"=="restart" goto restart_services
if "%ACTION%"=="logs" goto show_logs
if "%ACTION%"=="status" goto show_status
if "%ACTION%"=="urls" goto show_urls
if "%ACTION%"=="cleanup" goto cleanup
if "%ACTION%"=="backup" goto backup_data

echo 使用方法: %0 {start^|stop^|restart^|logs^|status^|urls^|cleanup^|backup}
echo.
echo 命令说明:
echo   start   - 启动所有服务
echo   stop    - 停止所有服务
echo   restart - 重启所有服务
echo   logs    - 查看服务日志
echo   status  - 查看服务状态
echo   urls    - 显示访问地址
echo   cleanup - 清理所有资源
echo   backup  - 备份数据库
pause
exit /b 1

:start_services
echo [INFO] 启动MyApp服务...

REM 构建并启动服务
docker-compose -f %COMPOSE_FILE% up -d --build
if errorlevel 1 (
    echo [ERROR] 服务启动失败
    pause
    exit /b 1
)

echo [INFO] 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
docker-compose -f %COMPOSE_FILE% ps | findstr "Up" >nul
if errorlevel 1 (
    echo [ERROR] 服务启动失败，请检查日志
    docker-compose -f %COMPOSE_FILE% logs
    pause
    exit /b 1
) else (
    echo [SUCCESS] 服务启动成功！
    call :show_status
    call :show_urls
)
goto end

:stop_services
echo [INFO] 停止MyApp服务...
docker-compose -f %COMPOSE_FILE% down
echo [SUCCESS] 服务已停止
goto end

:restart_services
echo [INFO] 重启MyApp服务...
call :stop_services
call :start_services
goto end

:show_logs
echo [INFO] 显示服务日志...
docker-compose -f %COMPOSE_FILE% logs -f
goto end

:show_status
echo [INFO] 服务状态:
docker-compose -f %COMPOSE_FILE% ps
echo.
echo [INFO] 容器资源使用情况:
for /f "tokens=*" %%i in ('docker-compose -f %COMPOSE_FILE% ps -q') do (
    docker stats --no-stream %%i 2>nul
)
goto :eof

:show_urls
echo.
echo [SUCCESS] === 服务访问地址 ===
echo 🌐 Web界面: http://localhost
echo 📡 API服务: http://localhost/api/database/list
echo 🔍 健康检查: http://localhost/health
echo 📊 数据库: localhost:3306 (用户: myapp, 密码: myapp123)
echo.
echo [INFO] 如果是远程服务器，请将localhost替换为服务器IP地址
goto :eof

:cleanup
echo [INFO] 清理Docker资源...
docker-compose -f %COMPOSE_FILE% down -v --rmi all
docker system prune -f
echo [SUCCESS] 清理完成
goto end

:backup_data
echo [INFO] 备份数据库数据...
set BACKUP_DIR=backups\%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%
set BACKUP_DIR=%BACKUP_DIR: =0%
mkdir %BACKUP_DIR% 2>nul
docker-compose -f %COMPOSE_FILE% exec mysql mysqldump -u myapp -pmyapp123 myapp > %BACKUP_DIR%\myapp_backup.sql
echo [SUCCESS] 数据备份完成: %BACKUP_DIR%\myapp_backup.sql
goto end

:end
echo.
echo [INFO] 操作完成
pause