package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS
//go:embed build/appicon.png
var icon []byte
var version = "0.0.1"
var gaMeasurementID, gaSecretKey string

const appName = "WeTools"

// 创建系统托盘菜单
func createTrayMenu(app *App) *menu.TrayMenu {
	// 创建托盘菜单项
	trayMenu := menu.NewMenu()
	
	// 添加菜单项
	trayMenu.Append(menu.Text("显示主窗口", keys.CmdOrCtrl("M"), func(data *menu.CallbackData) {
		// 显示主窗口
		runtime.WindowShow(app.ctx)
	}))
	
	trayMenu.Append(menu.Text("开机自启动", nil, func(data *menu.CallbackData) {
		// 切换开机自启动设置
		app.ToggleAutoStart()
	}))
	
	trayMenu.Append(menu.Separator())
	
	trayMenu.Append(menu.Text("退出", keys.CmdOrCtrl("Q"), func(data *menu.CallbackData) {
		// 退出应用
		runtime.Quit(app.ctx)
	}))
	
	return &menu.TrayMenu{
		Label: appName,
		Image: "main", // 引用trayicons/main.png
		Menu: trayMenu,
	}
}

func main() {
	// Create an instance of the app structure
	app := NewApp()
	
	// Create application with options
	err := wails.Run(&options.App{
		Title:  appName,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		HideWindowOnClose: true, // 关闭窗口时隐藏而不是退出
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
