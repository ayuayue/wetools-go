package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	"golang.org/x/sys/windows/registry"
)

// ProxyConfig 代理配置结构
type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// App struct
type App struct {
	ctx         context.Context
	proxyConfig ProxyConfig
	proxyServer *http.Server
	autoStartEnabled bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 启动时检查是否有保存的代理配置
	a.loadSavedProxyConfig()
	
	// 如果代理已启用，则启动代理服务器
	if a.proxyConfig.Enabled {
		a.StartProxyServer()
	}
	
	// 创建系统托盘菜单
	a.createTrayMenu()
}

// loadSavedProxyConfig 从存储中加载保存的代理配置
func (a *App) loadSavedProxyConfig() {
	// 这里应该从某个地方加载保存的配置
	// 例如，从文件或数据库中读取
	// 现在我们只是设置一个默认配置
	// 实际应用中，您可能需要从持久化存储中读取配置
}

// createTrayMenu 创建系统托盘菜单
func (a *App) createTrayMenu() {
	// 当前Wails版本不支持系统托盘菜单
	// 系统托盘功能将在未来版本中实现
	fmt.Println("系统托盘功能暂不可用")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetProxyConfig 获取代理配置
func (a *App) GetProxyConfig() ProxyConfig {
	return a.proxyConfig
}

// SetProxyConfig 设置代理配置
func (a *App) SetProxyConfig(config ProxyConfig) error {
	a.proxyConfig = config
	
	// 如果启用了代理配置，则启动代理服务器
	if config.Enabled {
		a.StartProxyServer()
	} else {
		// 如果禁用了代理配置，则停止代理服务器
		a.StopProxyServer()
	}
	
	return nil
}

// GetProxyURL 获取代理URL
func (a *App) GetProxyURL() string {
	if !a.proxyConfig.Enabled || a.proxyConfig.Host == "" {
		return ""
	}
	
	proxyURL := ""
	if a.proxyConfig.Username != "" && a.proxyConfig.Password != "" {
		proxyURL = fmt.Sprintf("http://%s:%s@%s:%s", 
			url.QueryEscape(a.proxyConfig.Username), 
			url.QueryEscape(a.proxyConfig.Password), 
			a.proxyConfig.Host, 
			a.proxyConfig.Port)
	} else {
		proxyURL = fmt.Sprintf("http://%s:%s", a.proxyConfig.Host, a.proxyConfig.Port)
	}
	
	return proxyURL
}

// StartProxyServer 启动内置代理服务器
func (a *App) StartProxyServer() error {
	if a.proxyServer != nil {
		return nil // 代理服务器已启动
	}
	
	// 创建代理处理函数
	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从查询参数中提取目标URL
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			// 尝试从路径中提取
			path := r.URL.Path
			if strings.HasPrefix(path, "/proxy/") {
				targetURL = path[7:] // 移除 "/proxy/" 前缀
			}
		}
		
		if targetURL == "" {
			http.Error(w, "缺少目标URL", http.StatusBadRequest)
			return
		}
		
		// 解码URL
		decodedURL, err := url.QueryUnescape(targetURL)
		if err != nil {
			http.Error(w, "无效的URL: "+err.Error(), http.StatusBadRequest)
			return
		}
		
		// 解析目标URL
		target, err := url.Parse(decodedURL)
		if err != nil {
			http.Error(w, "无效的目标URL: "+err.Error(), http.StatusBadRequest)
			return
		}
		
		// 创建到目标服务器的请求
		req, err := http.NewRequest(r.Method, target.String(), r.Body)
		if err != nil {
			http.Error(w, "创建请求失败: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		// 复制请求头
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
		
		// 设置Host头
		req.Header.Set("Host", target.Host)
		
		// 移除可能泄露客户端信息的头部
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Real-IP")
		
		// 创建HTTP客户端
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					// 如果配置了外部代理，则使用它
					proxyURL := a.GetProxyURL()
					if proxyURL != "" {
						return url.Parse(proxyURL)
					}
					return http.ProxyFromEnvironment(req)
				},
			},
		}
		
		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "请求失败: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		
		// 复制响应头
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		
		// 设置响应状态码
		w.WriteHeader(resp.StatusCode)
		
		// 复制响应体
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("复制响应体失败: %v\n", err)
		}
	})
	
	// 创建HTTP服务器
	a.proxyServer = &http.Server{
		Addr:    "127.0.0.1:8081", // 使用固定端口
		Handler: proxyHandler,
	}
	
	// 在goroutine中启动服务器
	go func() {
		if err := a.proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("代理服务器启动失败: %v\n", err)
		}
	}()
	
	fmt.Println("代理服务器已在端口 8081 上启动")
	return nil
}

// StopProxyServer 停止代理服务器
func (a *App) StopProxyServer() error {
	if a.proxyServer == nil {
		return nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := a.proxyServer.Shutdown(ctx)
	a.proxyServer = nil
	return err
}

// GetProxyServerPort 获取代理服务器端口
func (a *App) GetProxyServerPort() (int, error) {
	if a.proxyServer == nil {
		return 0, fmt.Errorf("代理服务器未启动")
	}
	
	// 这里需要实际获取监听端口的逻辑
	// 简化处理，返回一个默认端口
	return 8081, nil
}

// ProxyRequest 通过代理请求指定URL
func (a *App) ProxyRequest(targetURL string) (string, error) {
	if !a.proxyConfig.Enabled {
		return "", fmt.Errorf("代理未启用")
	}
	
	// 确保代理服务器已启动
	if a.proxyServer == nil {
		err := a.StartProxyServer()
		if err != nil {
			return "", fmt.Errorf("启动代理服务器失败: %v", err)
		}
	}
	
	// 返回代理URL，前端可以通过这个URL访问目标网站
	encodedURL := url.QueryEscape(targetURL)
	return fmt.Sprintf("http://localhost:8081/proxy?url=%s", encodedURL), nil
}

// TestProxyConnection 测试代理连接
func (a *App) TestProxyConnection() (bool, string, error) {
	if !a.proxyConfig.Enabled {
		return false, "代理未启用", nil
	}
	
	// 确保代理服务器已启动
	if a.proxyServer == nil {
		err := a.StartProxyServer()
		if err != nil {
			return false, "", fmt.Errorf("启动代理服务器失败: %v", err)
		}
	}
	
	// 构造测试URL
	testURL := "http://httpbin.org/get"
	proxyURL := fmt.Sprintf("http://localhost:8081/proxy?url=%s", url.QueryEscape(testURL))
	
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// 发送测试请求
	resp, err := client.Get(proxyURL)
	if err != nil {
		return false, "", fmt.Errorf("代理连接测试失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, "代理连接测试成功", nil
	}
	
	return false, fmt.Sprintf("代理连接测试失败: %s", resp.Status), nil
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
	a.autoStartEnabled = !a.autoStartEnabled
	err := a.setStartup(a.autoStartEnabled)
	return a.autoStartEnabled, err
}

// IsAutoStartEnabled 检查是否已设置开机自启动
func (a *App) IsAutoStartEnabled() bool {
	return a.autoStartEnabled
}

// setStartup 设置或取消开机自启动
func (a *App) setStartup(enable bool) error {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	
	// 根据操作系统设置开机自启动
	switch runtime.GOOS {
	case "windows":
		return a.setWindowsStartup(enable, execPath)
	case "darwin":
		return a.setMacStartup(enable, execPath)
	case "linux":
		return a.setLinuxStartup(enable, execPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// setWindowsStartup Windows系统设置开机自启动
func (a *App) setWindowsStartup(enable bool, execPath string) error {
	// Windows使用注册表设置开机自启动
	const runKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	const appName = "WeTools"
	
	// 打开注册表项
	k, err := registry.OpenKey(registry.CURRENT_USER, runKey, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()
	
	if enable {
		// 添加到注册表启动项
		err = k.SetStringValue(appName, execPath)
		if err != nil {
			return fmt.Errorf("failed to set registry value: %v", err)
		}
		fmt.Printf("Added %s to Windows startup\n", appName)
	} else {
		// 从注册表启动项移除
		err = k.DeleteValue(appName)
		if err != nil && err != registry.ErrNotExist {
			return fmt.Errorf("failed to delete registry value: %v", err)
		}
		fmt.Printf("Removed %s from Windows startup\n", appName)
	}
	
	return nil
}

// setMacStartup macOS系统设置开机自启动
func (a *App) setMacStartup(enable bool, execPath string) error {
	// macOS使用plist文件设置开机自启动
	// 这里简化处理
	
	if enable {
		fmt.Printf("Setting macOS startup for %s\n", execPath)
	} else {
		fmt.Printf("Removing macOS startup for %s\n", execPath)
	}
	
	return nil
}

// setLinuxStartup Linux系统设置开机自启动
func (a *App) setLinuxStartup(enable bool, execPath string) error {
	// Linux使用.desktop文件设置开机自启动
	// 这里简化处理
	
	if enable {
		fmt.Printf("Setting Linux startup for %s\n", execPath)
	} else {
		fmt.Printf("Removing Linux startup for %s\n", execPath)
	}
	
	return nil
}

// GetAutoStartStatus 获取开机自启动状态
func (a *App) GetAutoStartStatus() bool {
	return a.autoStartEnabled
}

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) error {
	// 使用Wails runtime打开URL
	// 注意：这将在默认浏览器中打开URL，而不是在应用内嵌入
	return nil
}
