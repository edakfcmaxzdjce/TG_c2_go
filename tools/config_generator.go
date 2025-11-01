package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=================================")
	fmt.Println("Telegram Bot 配置生成器")
	fmt.Println("=================================")

	// 获取 Bot Token
	fmt.Print("请输入您的 Bot Token: ")
	botToken, _ := reader.ReadString('\n')
	botToken = strings.TrimSpace(botToken)

	// 获取 Chat ID
	fmt.Print("请输入您的 Chat ID: ")
	chatID, _ := reader.ReadString('\n')
	chatID = strings.TrimSpace(chatID)

	// 获取 Bot 用户名
	fmt.Print("请输入您的 Bot 用户名 (不含@): ")
	botName, _ := reader.ReadString('\n')
	botName = strings.TrimSpace(botName)

	// 构建 API URLs
	baseURL := "https://api.telegram.org/bot" + botToken
	fileURL := "https://api.telegram.org/file/bot" + botToken

	urls := map[string]string{
		"UPDATE_URL":           baseURL + "/getUpdates",
		"SENDMSG_URL":          baseURL + "/sendMessage",
		"UPLOAD_FILE_URL":      baseURL + "/sendDocument",
		"SEND_PHOTO_URL":       baseURL + "/sendPhoto",
		"GET_FILE_PATH_URL":    baseURL + "/getFile",
		"CREATEFORUMTOPIC_URL": baseURL + "/createForumTopic",
		"DOWNLOAD_FILE_URL":    fileURL,
	}

	fmt.Println("\n=================================")
	fmt.Println("生成的配置 (请复制到 api.go 中):")
	fmt.Println("=================================")

	// 生成 Go 代码格式的配置
	fmt.Printf("CHAT_ID              = \"%s\"\n", chatID)
	fmt.Printf("BOT_NAME             = \"%s\"\n", botName)

	for name, url := range urls {
		encoded := base64.StdEncoding.EncodeToString([]byte(url))
		fmt.Printf("%-20s = \"%s\"\n", name, encoded)
	}

	fmt.Println("\n=================================")
	fmt.Println("原始 URLs (用于验证):")
	fmt.Println("=================================")
	for name, url := range urls {
		fmt.Printf("%-20s: %s\n", name, url)
	}

	fmt.Println("\n配置完成！请将上述配置复制到 config/api.go 文件中。")
}
