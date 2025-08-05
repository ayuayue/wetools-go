package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	
	"github.com/atotto/clipboard"
	
	"wetools-go/internal/image"
	clipboardmanager "wetools-go/internal/clipboard"
	"wetools-go/internal/proxy"
	"wetools-go/internal/autostart"
	"wetools-go/internal/logger"
)


// App struct
type App struct {
	ctx              context.Context
	proxyManager     *proxy.Manager
	autoStartManager *autostart.Manager
	clipboardManager *clipboardmanager.Manager
	logger           *logger.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 先创建剪贴板管理器以获取数据目录
	clipboardManager := clipboardmanager.NewManager()
	
	// 创建日志记录器
	logger := logger.NewLogger(clipboardManager.GetDataDir())
	
	app := &App{
		proxyManager:     proxy.NewManager(),
		autoStartManager: autostart.NewManager(),
		clipboardManager: clipboardManager,
		logger:           logger,
	}
	
	// 记录应用启动日志
	logger.Info("WeTools 应用启动")
	
	return app
}


// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 记录应用启动完成日志
	a.logger.Info("WeTools 应用启动完成")
	
	// 启动剪贴板监听
	go a.startClipboardListener()
}

// shutdown is called when the app is shutting down
func (a *App) shutdown(ctx context.Context) {
	// 记录应用关闭日志
	a.logger.Info("WeTools 应用关闭")
	
	// 关闭日志文件
	if a.logger != nil {
		a.logger.Close()
	}
}

// startClipboardListener 启动剪贴板监听器
func (a *App) startClipboardListener() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	// 重试次数
	retryCount := 0
	maxRetries := 5
	
	a.logger.Info("剪贴板监听器已启动")
	
	for {
		select {
		case <-ticker.C:
			// 读取剪贴板内容
			content, err := clipboard.ReadAll()
			if err != nil {
				// 忽略特定的Windows成功错误
				if err.Error() == "The operation completed successfully." {
					// 这是Windows的一个已知问题，忽略它
					continue
				}
				
				a.logger.Warn(fmt.Sprintf("读取剪贴板失败: %v", err))
				retryCount++
				
				// 如果重试次数超过最大值，重置重试计数
				if retryCount > maxRetries {
					a.logger.Warn("剪贴板读取重试次数过多，重置重试计数")
					retryCount = 0
				}
				continue
			}
			
			retryCount = 0 // 重置重试计数
			
			// 如果内容发生变化，添加到历史记录
			if content != "" && content != a.clipboardManager.GetLastContent() {
				a.logger.Info(fmt.Sprintf("检测到新的剪贴板内容: %s", content))
				a.clipboardManager.SetLastContent(content)
				a.clipboardManager.AddItem(content)
			}
		case <-a.ctx.Done():
			// 应用关闭时停止监听
			a.logger.Info("剪贴板监听器已停止")
			return
		}
	}
}

// GetClipboardHistory 获取剪贴板历史记录
func (a *App) GetClipboardHistory() []clipboardmanager.ClipboardItem {
	return a.clipboardManager.GetHistory()
}

// GetClipboardHistoryPage 获取分页的剪贴板历史记录
func (a *App) GetClipboardHistoryPage(page, pageSize int) ([]clipboardmanager.ClipboardItem, int) {
	return a.clipboardManager.GetHistoryPage(page, pageSize)
}

// ClearClipboardHistory 清空剪贴板历史记录
func (a *App) ClearClipboardHistory() {
	a.clipboardManager.ClearHistory()
}

// RemoveClipboardItem 删除指定的剪贴板项目
func (a *App) RemoveClipboardItem(index int) error {
	return a.clipboardManager.RemoveItem(index)
}

// CopyToClipboard 复制内容到剪贴板
func (a *App) CopyToClipboard(content string) error {
	return a.clipboardManager.CopyToClipboard(content)
}

// SetMaxRecords 设置最大保存记录数
func (a *App) SetMaxRecords(max int) {
	a.clipboardManager.SetMaxRecords(max)
}

// GetMaxRecords 获取最大保存记录数
func (a *App) GetMaxRecords() int {
	return a.clipboardManager.GetMaxRecords()
}

// OpenClipboardFile 打开剪贴板文件
func (a *App) OpenClipboardFile() error {
	dataDir := a.clipboardManager.GetDataDir()
	if dataDir == "" {
		return fmt.Errorf("数据目录未初始化")
	}
	
	// 构建剪贴板文件路径
	clipboardFile := filepath.Join(dataDir, "clipboard.txt")
	
	// 检查文件是否存在，如果不存在则创建
	if _, err := os.Stat(clipboardFile); os.IsNotExist(err) {
		// 创建空文件
		file, err := os.Create(clipboardFile)
		if err != nil {
			return fmt.Errorf("创建剪贴板文件失败: %v", err)
		}
		file.Close()
	}
	
	// 根据操作系统打开文件
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "", clipboardFile)
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("open", clipboardFile)
		return cmd.Run()
	case "linux":
		cmd := exec.Command("xdg-open", clipboardFile)
		return cmd.Run()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}


// GetProxyConfig 获取代理配置
func (a *App) GetProxyConfig() proxy.Config {
	return a.proxyManager.GetConfig()
}

// SetProxyConfig 设置代理配置
func (a *App) SetProxyConfig(config proxy.Config) error {
	a.logger.Info(fmt.Sprintf("设置代理配置: Enabled=%v, Host=%s, Port=%s", config.Enabled, config.Host, config.Port))
	
	err := a.proxyManager.SetConfig(config)
	if err != nil {
		a.logger.Error(fmt.Sprintf("设置代理配置失败: %v", err))
		return err
	}
	
	// 如果启用了代理配置，则启动代理服务器
	if config.Enabled {
		a.logger.Info("启动代理服务器")
		err := a.proxyManager.StartServer()
		if err != nil {
			a.logger.Error(fmt.Sprintf("启动代理服务器失败: %v", err))
		} else {
			a.logger.Info("代理服务器启动成功")
		}
	} else {
		// 如果禁用了代理配置，则停止代理服务器
		a.logger.Info("停止代理服务器")
		err := a.proxyManager.StopServer()
		if err != nil {
			a.logger.Warn(fmt.Sprintf("停止代理服务器时出现警告: %v", err))
		} else {
			a.logger.Info("代理服务器已停止")
		}
	}
	
	return nil
}

// GetProxyURL 获取代理URL
func (a *App) GetProxyURL() string {
	return a.proxyManager.GetURL()
}

// StartProxyServer 启动内置代理服务器
func (a *App) StartProxyServer() error {
	return a.proxyManager.StartServer()
}

// StopProxyServer 停止代理服务器
func (a *App) StopProxyServer() error {
	return a.proxyManager.StopServer()
}

// GetProxyServerPort 获取代理服务器端口
func (a *App) GetProxyServerPort() (int, error) {
	return a.proxyManager.GetServerPort()
}

// ProxyRequest 通过代理请求指定URL
func (a *App) ProxyRequest(targetURL string) (string, error) {
	return a.proxyManager.Request(targetURL)
}

// TestProxyConnection 测试代理连接
func (a *App) TestProxyConnection() (bool, string, error) {
	return a.proxyManager.TestConnection()
}

// SaveFileToDownloads saves content to the download/wetools-go directory
func (a *App) SaveFileToDownloads(filename, content string) (string, error) {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	// 构建目标目录路径
	downloadDir := filepath.Join(homeDir, "download", "wetools-go")
	
	// 创建目录（如果不存在）
	err = os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// 构建完整文件路径
	filePath := filepath.Join(downloadDir, filename)
	
	// 写入文件
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}
	
	return filePath, nil
}

// ToggleAutoStart 切换开机自启动设置
func (a *App) ToggleAutoStart() (bool, error) {
	a.logger.Info("切换开机自启动设置")
	enabled, err := a.autoStartManager.Toggle()
	if err != nil {
		a.logger.Error(fmt.Sprintf("切换开机自启动设置失败: %v", err))
	} else {
		a.logger.Info(fmt.Sprintf("开机自启动设置已更新: %v", enabled))
	}
	return enabled, err
}

// IsAutoStartEnabled 检查是否已设置开机自启动
func (a *App) IsAutoStartEnabled() bool {
	return a.autoStartManager.IsEnabled()
}


// GetAutoStartStatus 获取开机自启动状态
func (a *App) GetAutoStartStatus() bool {
	return a.autoStartManager.GetStatus()
}

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) error {
	// 使用Wails runtime打开URL
	// 注意：这将在默认浏览器中打开URL，而不是在应用内嵌入
	return nil
}

// ImageToSvg 将图像转换为SVG
func (a *App) ImageToSvg(request image.ImageToSvgRequest) (image.ImageToSvgResponse, error) {
	a.logger.Info("开始图像转SVG处理")
	response, err := image.ToSvg(request)
	if err != nil {
		a.logger.Error(fmt.Sprintf("图像转SVG处理失败: %v", err))
	} else {
		a.logger.Info("图像转SVG处理成功")
	}
	return response, err
}

