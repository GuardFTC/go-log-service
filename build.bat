@echo off
echo 开始编译 logging-mon-service...

:: 创建输出目录
if not exist "..\..\bin" mkdir "..\..\bin"

:: 设置项目名称
set PROJECT_NAME=logging-mon-service

:: 编译Linux版本 (amd64)
echo 编译Linux版本...
set GOOS=linux
set GOARCH=amd64
go build -o "..\..\bin\%PROJECT_NAME%-linux-amd64" .
if %errorlevel% neq 0 (
    echo Linux编译失败!
    exit /b 1
)

:: 编译Windows版本 (amd64)
echo 编译Windows版本...
set GOOS=windows
set GOARCH=amd64
go build -o "..\..\bin\%PROJECT_NAME%-windows.exe" .
if %errorlevel% neq 0 (
    echo Windows编译失败!
    exit /b 1
)

:: 编译Linux ARM64版本
echo 编译Linux ARM64版本...
set GOOS=linux
set GOARCH=arm64
go build -o "..\..\bin\%PROJECT_NAME%-linux-arm64" .
if %errorlevel% neq 0 (
    echo Linux ARM64编译失败!
    exit /b 1
)

echo.
echo 编译完成! 可执行文件已生成到 ..\..\bin\ 目录:
echo - %PROJECT_NAME%-linux-amd64 (Linux x64)
echo - %PROJECT_NAME%-linux-arm64 (Linux ARM64)
echo - %PROJECT_NAME%-windows.exe (Windows x64)