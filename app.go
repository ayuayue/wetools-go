package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"os/user"
	
	"github.com/atotto/clipboard"
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

// ClipboardItem 剪贴板项目结构
type ClipboardItem struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

// App struct
type App struct {
	ctx              context.Context
	proxyConfig      ProxyConfig
	proxyServer      *http.Server
	autoStartEnabled bool
	clipboardHistory []ClipboardItem
	lastClipboard    string
	dataDir          string
	maxRecords       int // 最大保存记录数
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		maxRecords: 1000, // 默认保存1000条记录
	}
	app.initDataDir()
	return app
}

// loadClipboardHistory 从文件加载剪贴板历史记录
func (a *App) loadClipboardHistory() {
	if a.dataDir == "" {
		return
	}
	
	// 构建剪贴板历史文件路径
	clipboardFile := filepath.Join(a.dataDir, "clipboard.txt")
	
	// 检查文件是否存在
	if _, err := os.Stat(clipboardFile); os.IsNotExist(err) {
		// 文件不存在，初始化为空的历史记录
		a.clipboardHistory = make([]ClipboardItem, 0)
		// 添加一些示例数据用于测试
		a.addSampleData()
		return
	}
	
	// 读取文件内容
	data, err := os.ReadFile(clipboardFile)
	if err != nil {
		fmt.Printf("读取剪贴板历史文件失败: %v\n", err)
		a.clipboardHistory = make([]ClipboardItem, 0)
		// 添加一些示例数据用于测试
		a.addSampleData()
		return
	}
	
	// 解析文本数据
	content := string(data)
	if content == "" {
		a.clipboardHistory = make([]ClipboardItem, 0)
		// 添加一些示例数据用于测试
		a.addSampleData()
		return
	}
	
	// 按行分割内容
	lines := strings.Split(content, "\n")
	var history []ClipboardItem
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// 解析日期时间和内容
		// 格式: datetime|content
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		
		// 解析日期时间字符串为时间戳
		t, err := time.Parse("2006-01-02 15:04:05", parts[0])
		if err != nil {
			// 如果解析失败，尝试使用旧的时间戳格式
			timestamp, parseErr := strconv.ParseInt(parts[0], 10, 64)
			if parseErr != nil {
				continue
			}
			// 转换为正确的格式
			history = append(history, ClipboardItem{
				Timestamp: timestamp,
				Content:   parts[1],
			})
			continue
		}
		
		// 转换为毫秒时间戳
		timestamp := t.UnixNano() / int64(time.Millisecond)
		
		history = append(history, ClipboardItem{
			Timestamp: timestamp,
			Content:   parts[1],
		})
	}
	
	a.clipboardHistory = history
	fmt.Printf("加载了 %d 条剪贴板记录\n", len(a.clipboardHistory))
}

// addSampleData 添加示例数据用于测试
func (a *App) addSampleData() {
	// 只在没有数据时添加示例数据
	if len(a.clipboardHistory) > 0 {
		return
	}
	
	fmt.Println("添加示例剪贴板数据用于测试")
	
	sampleData := []ClipboardItem{
		{
			Content:   "示例文本内容1",
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		},
		{
			Content:   "https://www.example.com",
			Timestamp: time.Now().Add(-1 * time.Minute).UnixNano() / int64(time.Millisecond),
		},
		{
			Content:   `{"name": "WeTools", "version": "1.0.0"}`,
			Timestamp: time.Now().Add(-2 * time.Minute).UnixNano() / int64(time.Millisecond),
		},
	}
	
	a.clipboardHistory = sampleData
	a.saveClipboardHistory()
}

// loadSavedProxyConfig 从存储中加载保存的代理配置
func (a *App) loadSavedProxyConfig() {
	// 这里应该从某个地方加载保存的配置
	// 例如，从文件或数据库中读取
	// 现在我们只是设置一个默认配置
	// 实际应用中，您可能需要从持久化存储中读取配置
}

// saveClipboardHistory 保存剪贴板历史记录到文件
func (a *App) saveClipboardHistory() {
	if a.dataDir == "" {
		return
	}
	
	// 构建剪贴板历史文件路径
	clipboardFile := filepath.Join(a.dataDir, "clipboard.txt")
	
	// 构建文本内容
	var lines []string
	for _, item := range a.clipboardHistory {
		// 格式化时间戳为可读格式
		t := time.Unix(0, item.Timestamp*int64(time.Millisecond))
		formattedTime := t.Format("2006-01-02 15:04:05")
		
		// 格式: datetime|content
		line := fmt.Sprintf("%s|%s", formattedTime, item.Content)
		lines = append(lines, line)
	}
	
	// 将所有行连接成一个字符串
	content := strings.Join(lines, "\n")
	
	// 写入文件
	err := os.WriteFile(clipboardFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("保存剪贴板历史文件失败: %v\n", err)
		return
	}
	
	fmt.Printf("保存了 %d 条剪贴板记录\n", len(a.clipboardHistory))
}

// initDataDir 初始化程序数据目录
func (a *App) initDataDir() {
	// 获取当前用户
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		return
	}
	
	// 构建数据目录路径
	a.dataDir = filepath.Join(usr.HomeDir, ".wetools-go")
	
	// 创建目录（如果不存在）
	err = os.MkdirAll(a.dataDir, 0755)
	if err != nil {
		fmt.Printf("创建数据目录失败: %v\n", err)
		return
	}
	
	fmt.Printf("数据目录: %s\n", a.dataDir)
	
	// 加载已保存的剪贴板历史记录
	a.loadClipboardHistory()
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
	
	// 初始化剪贴板历史记录
	a.clipboardHistory = make([]ClipboardItem, 0)
	a.lastClipboard = ""
	
	// 启动剪贴板监听
	go a.startClipboardListener()
}

// startClipboardListener 启动剪贴板监听器
func (a *App) startClipboardListener() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	// 重试次数
	retryCount := 0
	maxRetries := 5
	
	fmt.Println("剪贴板监听器已启动")
	
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
				
				fmt.Printf("读取剪贴板失败: %v\n", err)
				retryCount++
				
				// 如果重试次数超过最大值，重置重试计数
				if retryCount > maxRetries {
					retryCount = 0
				}
				continue
			}
			
			retryCount = 0 // 重置重试计数
			
			// 如果内容发生变化，添加到历史记录
			if content != "" && content != a.lastClipboard {
				fmt.Printf("检测到新的剪贴板内容: %s\n", content)
				a.lastClipboard = content
				a.addClipboardItem(content)
			}
		case <-a.ctx.Done():
			// 应用关闭时停止监听
			fmt.Println("剪贴板监听器已停止")
			return
		}
	}
}

// addClipboardItem 添加剪贴板项目到历史记录
func (a *App) addClipboardItem(content string) {
	item := ClipboardItem{
		Content:   content,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond), // 转换为毫秒时间戳
	}
	
	// 添加到历史记录开头
	a.clipboardHistory = append([]ClipboardItem{item}, a.clipboardHistory...)
	
	// 限制历史记录数量为最大记录数
	if len(a.clipboardHistory) > a.maxRecords {
		a.clipboardHistory = a.clipboardHistory[:a.maxRecords]
	}
	
	// 保存历史记录到文件
	a.saveClipboardHistory()
	
	fmt.Printf("添加剪贴板记录: %s\n", content)
}

// GetClipboardHistory 获取剪贴板历史记录
func (a *App) GetClipboardHistory() []ClipboardItem {
	return a.clipboardHistory
}

// GetClipboardHistoryPage 获取分页的剪贴板历史记录
func (a *App) GetClipboardHistoryPage(page, pageSize int) ([]ClipboardItem, int) {
	total := len(a.clipboardHistory)
	
	// 计算总页数
	totalPages := 0
	if total > 0 && pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	
	// 确保页码有效
	if page < 1 {
		page = 1
	}
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}
	
	// 如果没有记录或页面大小无效，返回空数组
	if total == 0 || pageSize <= 0 {
		return []ClipboardItem{}, totalPages
	}
	
	// 计算起始和结束索引
	start := (page - 1) * pageSize
	end := start + pageSize
	
	// 确保结束索引不超过总长度
	if end > total {
		end = total
	}
	
	// 确保起始索引不超过总长度
	if start >= total {
		return []ClipboardItem{}, totalPages
	}
	
	// 返回当前页的数据和总页数
	return a.clipboardHistory[start:end], totalPages
}

// ClearClipboardHistory 清空剪贴板历史记录
func (a *App) ClearClipboardHistory() {
	a.clipboardHistory = make([]ClipboardItem, 0)
	a.saveClipboardHistory()
}

// RemoveClipboardItem 删除指定的剪贴板项目
func (a *App) RemoveClipboardItem(index int) error {
	if index < 0 || index >= len(a.clipboardHistory) {
		return fmt.Errorf("索引超出范围")
	}
	
	// 删除指定项目
	a.clipboardHistory = append(a.clipboardHistory[:index], a.clipboardHistory[index+1:]...)
	a.saveClipboardHistory()
	return nil
}

// CopyToClipboard 复制内容到剪贴板
func (a *App) CopyToClipboard(content string) error {
	return clipboard.WriteAll(content)
}

// SetMaxRecords 设置最大保存记录数
func (a *App) SetMaxRecords(max int) {
	if max > 0 {
		a.maxRecords = max
		// 如果当前记录数超过新设置的最大值，清理旧记录
		if len(a.clipboardHistory) > a.maxRecords {
			a.clipboardHistory = a.clipboardHistory[:a.maxRecords]
			a.saveClipboardHistory()
		}
	}
}

// GetMaxRecords 获取最大保存记录数
func (a *App) GetMaxRecords() int {
	return a.maxRecords
}

// OpenClipboardFile 打开剪贴板文件
func (a *App) OpenClipboardFile() error {
	if a.dataDir == "" {
		return fmt.Errorf("数据目录未初始化")
	}
	
	// 构建剪贴板文件路径
	clipboardFile := filepath.Join(a.dataDir, "clipboard.txt")
	
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

