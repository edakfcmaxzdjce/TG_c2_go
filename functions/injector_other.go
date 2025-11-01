//go:build !windows

package functions

import (
	"fmt"

	"TG_c2_go/telegram"
)

// Injector 代码注入功能
type Injector struct {
	client *telegram.TelegramClient
}

// NewInjector 创建新的注入器
func NewInjector() *Injector {
	return &Injector{
		client: telegram.NewTelegramClient(),
	}
}

// InjectShellcode 注入shellcode (非Windows平台不支持)
func (inj *Injector) InjectShellcode(fileID string, tgClient *telegram.TelegramClient) error {
	msg := "Shellcode注入功能仅支持Windows系统"
	inj.client.SendMessage(msg)
	return fmt.Errorf(msg)
}
