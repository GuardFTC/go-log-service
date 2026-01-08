@echo off
chcp 65001 >nul
REM Docker镜像构建脚本 (Windows版本)

REM 设置变量
set "IMAGE_NAME=logging-mon-service"
set "VERSION=%1"
set "REGISTRY=%2"

REM 如果没有指定版本，使用latest
if "%VERSION%"=="" set "VERSION=latest"

REM 构建完整镜像名称
if "%REGISTRY%"=="" (
    set "FULL_IMAGE_NAME=%IMAGE_NAME%:%VERSION%"
) else (
    set "FULL_IMAGE_NAME=%REGISTRY%/%IMAGE_NAME%:%VERSION%"
)

echo 开始构建Docker镜像...
echo 镜像名称: "%FULL_IMAGE_NAME%"
echo.

REM 构建镜像
docker build -t "%FULL_IMAGE_NAME%" .

if %ERRORLEVEL% equ 0 (
    echo.
    echo Docker镜像构建成功!
    echo 镜像名称: "%FULL_IMAGE_NAME%"
    echo.
    echo 使用方法:
    echo 1. 运行容器 (开发环境):
    echo    docker run -p 39801:39801 "%FULL_IMAGE_NAME%" -env dev
    echo.
    echo 2. 运行容器 (生产环境):
    echo    docker run -p 39801:39801 "%FULL_IMAGE_NAME%"
    echo.
    echo 3. 运行容器 (自定义端口):
    echo    docker run -p 8080:8080 -e SERVICE_PORT=8080 "%FULL_IMAGE_NAME%"
    echo.
    echo 4. 推送到镜像仓库:
    echo    docker push "%FULL_IMAGE_NAME%"
    echo.
    
    REM 显示镜像信息
    echo 镜像信息:
    docker images | findstr "%IMAGE_NAME%"
) else (
    echo.
    echo Docker镜像构建失败!
    exit /b 1
)