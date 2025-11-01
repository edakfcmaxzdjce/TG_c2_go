# Telegram Bot 配置指南

## 📱 完整配置步骤

### 步骤 1: 创建 Telegram Bot

1. **打开 Telegram**，搜索 `@BotFather`
2. **发送** `/start` 开始对话
3. **发送** `/newbot` 创建新机器人
4. **输入机器人名称**（显示名称）：例如 `My C2 Bot`
5. **输入机器人用户名**（必须以bot结尾）：例如 `myc2bot_bot`

**成功后会收到类似消息：**
```
Done! Congratulations on your new bot. You will find it at t.me/myc2bot_bot.
You can now add a description, about section and profile picture for your bot.

Use this token to access the HTTP API:
1234567890:ABCdefGhiJklMnoPqrsTuvWxyz-example

Keep your token secure and store it safely, it can be used by anyone to control your bot.
```

### 步骤 2: 获取 Chat ID

**方法 A: 创建群组**
1. 创建一个新的 Telegram 群组
2. 将您的 Bot 添加到群组中
3. 给 Bot 管理员权限
4. 在群组中发送任意消息
5. 访问 URL：`https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
6. 在返回的 JSON 中找到 `"chat":{"id":-1001234567890}`

**方法 B: 使用论坛群组（推荐）**
1. 创建一个论坛群组（SuperGroup with Topics）
2. 添加 Bot 为管理员，给予创建主题权限
3. 获取群组 ID

### 步骤 3: 使用配置生成器

运行配置生成器工具：

```bash
cd tools
go run config_generator.go
```

**示例输入：**
```
Bot Token: 1234567890:ABCdefGhiJklMnoPqrsTuvWxyz-example
Chat ID: -1001234567890
Bot Name: myc2bot_bot
```

### 步骤 4: 更新配置文件

将生成的配置复制到 `config/api.go` 文件中：

```go
var (
    // Telegram Bot配置
    UPLOAD_FILE_URL      = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    CHAT_ID              = "-1001234567890"
    UPDATE_URL           = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    SENDMSG_URL          = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    CREATEFORUMTOPIC_URL = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    SEND_PHOTO_URL       = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    BOT_NAME             = "myc2bot_bot"
    GET_FILE_PATH_URL    = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2JvdDEyMzQ1Njc4OTA6QUJDZGVmR2hpSmtsT..."
    DOWNLOAD_FILE_URL    = "aHR0cHM6Ly9hcGkudGVsZWdyYW0ub3JnL2ZpbGUvYm90MTIzNDU2Nzg5MDpBQkNkZWZH..."
    // ...
)
```

## 🔧 手动配置方法

如果您想手动配置，可以使用在线 Base64 编码工具：

### 需要编码的 URLs：

假设您的 Bot Token 是 `1234567890:ABCdefGhiJklMnoPqrsTuvWxyz`

```
UPDATE_URL:           https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/getUpdates
SENDMSG_URL:          https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/sendMessage
UPLOAD_FILE_URL:      https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/sendDocument
SEND_PHOTO_URL:       https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/sendPhoto
GET_FILE_PATH_URL:    https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/getFile
CREATEFORUMTOPIC_URL: https://api.telegram.org/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz/createForumTopic
DOWNLOAD_FILE_URL:    https://api.telegram.org/file/bot1234567890:ABCdefGhiJklMnoPqrsTuvWxyz
```

## 🛡️ 安全注意事项

1. **保护 Bot Token**: 这是最重要的凭证，不要泄露给他人
2. **使用私人群组**: 确保只有您能够访问控制群组
3. **定期轮换**: 如有必要，可以通过 BotFather 重新生成 Token
4. **监控访问**: 定期检查 Bot 的使用情况

## 🧪 测试配置

配置完成后，可以运行程序测试：

```bash
go run main.go
```

程序启动后应该能够：
1. 连接到 Telegram
2. 创建或找到现有主题
3. 发送连接确认消息

## ❗ 常见问题

**Q: Bot Token 格式不对？**
A: Token 格式应为 `数字:字母数字字符串`，例如 `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`

**Q: Chat ID 获取失败？**
A: 确保 Bot 已添加到群组并有管理员权限，然后在群组发送消息后访问 getUpdates API

**Q: 程序运行但无法收到消息？**
A: 检查 Chat ID 是否正确，群组 ID 通常以 `-100` 开头

**Q: Base64 编码错误？**
A: 使用提供的工具生成，或使用在线 Base64 编码器，确保没有多余的空格或换行符
