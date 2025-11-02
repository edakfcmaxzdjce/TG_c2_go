#!/bin/bash

# TG C2 Go 构建脚本

echo "===================="
echo "TG C2 Go Build"
echo "===================="

# 设置项目信息
PROJECT_NAME="TG_c2_go"
VERSION="1.0.0"

echo "项目: $PROJECT_NAME"
echo "版本: $VERSION"
echo ""

# 选择构建模式
echo "请选择构建模式："
echo "  1) 正常版本 (Windows隐藏窗口，无回显)"
echo "  2) Debug版本 (显示控制台和所有输出)"
echo ""
read -p "请输入选项 (1 或 2，默认 1): " BUILD_MODE
BUILD_MODE=${BUILD_MODE:-1}

if [ "$BUILD_MODE" = "2" ]; then
    BUILD_TYPE="debug"
    WINDOWS_LDFLAGS='-ldflags="-s -w"'
    OTHER_LDFLAGS='-ldflags="-s -w"'
    echo "✅ 选择: Debug版本 (显示控制台)"
else
    BUILD_TYPE="release"
    # Windows使用windowsgui隐藏控制台窗口
    WINDOWS_LDFLAGS='-ldflags="-H windowsgui -s -w"'
    OTHER_LDFLAGS='-ldflags="-s -w"'
    echo "✅ 选择: 正常版本 (Windows隐藏窗口)"
fi
echo ""

# 检查 Go 版本
echo "检查 Go 环境..."
go version
if [ $? -ne 0 ]; then
    echo "❌ Go 未安装或未在 PATH 中"
    exit 1
fi
echo "✅ Go 环境正常"
echo ""

# 安装依赖
echo "安装依赖..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ 依赖安装失败"
    exit 1
fi
echo "✅ 依赖安装完成"
echo ""

# 创建构建目录
BUILD_DIR="build"
mkdir -p $BUILD_DIR

echo "开始构建..."
echo ""
echo "使用以下优化标志减小二进制大小："
echo "  -trimpath: 移除文件路径信息"
echo "  构建模式: $BUILD_TYPE"
echo ""

# Windows 64位
echo "构建 Windows 64位版本..."
eval "GOOS=windows GOARCH=amd64 go build -trimpath $WINDOWS_LDFLAGS -o $BUILD_DIR/TG_c2_go_windows_amd64.exe"
if [ $? -eq 0 ]; then
    SIZE=$(ls -lh $BUILD_DIR/TG_c2_go_windows_amd64.exe | awk '{print $5}')
    echo "✅ Windows 64位构建完成 (大小: $SIZE, 模式: $BUILD_TYPE)"
else
    echo "❌ Windows 64位构建失败"
fi

# Windows 32位
echo "构建 Windows 32位版本..."
eval "GOOS=windows GOARCH=386 go build -trimpath $WINDOWS_LDFLAGS -o $BUILD_DIR/TG_c2_go_windows_386.exe"
if [ $? -eq 0 ]; then
    SIZE=$(ls -lh $BUILD_DIR/TG_c2_go_windows_386.exe | awk '{print $5}')
    echo "✅ Windows 32位构建完成 (大小: $SIZE, 模式: $BUILD_TYPE)"
else
    echo "❌ Windows 32位构建失败"
fi

# Linux 64位
echo "构建 Linux 64位版本..."
eval "GOOS=linux GOARCH=amd64 go build -trimpath $OTHER_LDFLAGS -o $BUILD_DIR/TG_c2_go_linux_amd64"
if [ $? -eq 0 ]; then
    SIZE=$(ls -lh $BUILD_DIR/TG_c2_go_linux_amd64 | awk '{print $5}')
    echo "✅ Linux 64位构建完成 (大小: $SIZE, 模式: $BUILD_TYPE)"
else
    echo "❌ Linux 64位构建失败"
fi

# macOS 64位 (Intel)
echo "构建 macOS 64位版本..."
eval "GOOS=darwin GOARCH=amd64 go build -trimpath $OTHER_LDFLAGS -o $BUILD_DIR/TG_c2_go_darwin_amd64"
if [ $? -eq 0 ]; then
    SIZE=$(ls -lh $BUILD_DIR/TG_c2_go_darwin_amd64 | awk '{print $5}')
    echo "✅ macOS Intel 64位构建完成 (大小: $SIZE, 模式: $BUILD_TYPE)"
else
    echo "❌ macOS Intel 64位构建失败"
fi

# macOS ARM64 (Apple Silicon)
echo "构建 macOS ARM64版本..."
eval "GOOS=darwin GOARCH=arm64 go build -trimpath $OTHER_LDFLAGS -o $BUILD_DIR/TG_c2_go_darwin_arm64"
if [ $? -eq 0 ]; then
    SIZE=$(ls -lh $BUILD_DIR/TG_c2_go_darwin_arm64 | awk '{print $5}')
    echo "✅ macOS ARM64构建完成 (大小: $SIZE, 模式: $BUILD_TYPE)"
else
    echo "❌ macOS ARM64构建失败"
fi

echo ""
echo "构建完成！生成的文件在 $BUILD_DIR 目录中："
ls -la $BUILD_DIR/

echo ""
echo "===================="
echo "构建脚本执行完成"
echo "===================="
