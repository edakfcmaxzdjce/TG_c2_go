@echo off
echo ====================
echo TG C2 Go Build
echo ====================

rem 设置项目信息
set PROJECT_NAME=TG_c2_go
set VERSION=1.0.0

echo 项目: %PROJECT_NAME%
echo 版本: %VERSION%
echo.

rem 选择构建模式
echo 请选择构建模式：
echo   1) 正常版本 (Windows隐藏窗口，无回显)
echo   2) Debug版本 (显示控制台和所有输出)
echo.
set /p BUILD_MODE="请输入选项 (1 或 2，默认 1): "
if "%BUILD_MODE%"=="" set BUILD_MODE=1
if "%BUILD_MODE%"=="2" (
    set BUILD_TYPE=debug
    set WINDOWS_LDFLAGS=-ldflags="-s -w"
    echo ✅ 选择: Debug版本 (显示控制台)
) else (
    set BUILD_TYPE=release
    rem Windows使用windowsgui隐藏控制台窗口
    set WINDOWS_LDFLAGS=-ldflags="-H windowsgui -s -w"
    echo ✅ 选择: 正常版本 (Windows隐藏窗口)
)
echo.

rem 检查 Go 版本
echo 检查 Go 环境...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 未安装或未在 PATH 中
    pause
    exit /b 1
)
echo ✅ Go 环境正常
echo.

rem 安装依赖
echo 安装依赖...
go mod tidy
if %errorlevel% neq 0 (
    echo ❌ 依赖安装失败
    pause
    exit /b 1
)
echo ✅ 依赖安装完成
echo.

rem 创建构建目录
set BUILD_DIR=build
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

echo 开始构建...
echo 使用以下优化标志减小二进制大小：
echo   -trimpath: 移除文件路径信息
echo   构建模式: %BUILD_TYPE%
echo.

rem Windows 64位
echo 构建 Windows 64位版本...
set GOOS=windows
set GOARCH=amd64
go build -trimpath %WINDOWS_LDFLAGS% -o %BUILD_DIR%/TG_c2_go_windows_amd64.exe
if %errorlevel% equ 0 (
    for %%A in (%BUILD_DIR%\TG_c2_go_windows_amd64.exe) do set SIZE=%%~zA
    set /a SIZE_MB=%SIZE%/1024/1024
    echo ✅ Windows 64位构建完成 (大小: %SIZE_MB%MB, 模式: %BUILD_TYPE%)
) else (
    echo ❌ Windows 64位构建失败
)

rem Windows 32位
echo 构建 Windows 32位版本...
set GOOS=windows
set GOARCH=386
go build -trimpath %WINDOWS_LDFLAGS% -o %BUILD_DIR%/TG_c2_go_windows_386.exe
if %errorlevel% equ 0 (
    for %%A in (%BUILD_DIR%\TG_c2_go_windows_386.exe) do set SIZE=%%~zA
    set /a SIZE_MB=%SIZE%/1024/1024
    echo ✅ Windows 32位构建完成 (大小: %SIZE_MB%MB, 模式: %BUILD_TYPE%)
) else (
    echo ❌ Windows 32位构建失败
)

rem Linux 64位
echo 构建 Linux 64位版本...
set GOOS=linux
set GOARCH=amd64
go build -trimpath -ldflags="-s -w" -o %BUILD_DIR%/TG_c2_go_linux_amd64
if %errorlevel% equ 0 (
    echo ✅ Linux 64位构建完成
) else (
    echo ❌ Linux 64位构建失败
)

echo.
echo 构建完成！生成的文件在 %BUILD_DIR% 目录中：
dir %BUILD_DIR%

echo.
echo ====================
echo 构建脚本执行完成
echo ====================
pause
