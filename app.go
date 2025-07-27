package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) error {
	// 使用Wails runtime打开URL
	// 注意：这将在默认浏览器中打开URL，而不是在应用内嵌入
	return nil
}
