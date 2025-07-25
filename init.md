# wetools-go 项目初始化文档

## 项目概述

这是一个使用 Wails 框架构建的桌面应用程序项目。Wails 允许使用 Go 语言作为后端，Vue.js 作为前端来构建跨平台的桌面应用。

## 技术栈

- **后端**: Go 语言
- **前端**: Vue 3 (版本 3.5.17)
- **构建工具**: Vite (版本 3.2.11)
- **UI框架**: Vue.js
- **包管理**: npm


## 开发环境搭建

1. 安装 Go 语言环境 (推荐版本 1.21+)
2. 安装 Node.js 和 npm
3. 安装 Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
4. 验证安装: `wails doctor`

## 开发流程

1. 运行开发服务器: `wails dev`
2. 构建生产版本: `wails build`
3. 生成应用程序图标: `wails generate icon`

## 构建命令

- 开发模式: `wails dev`
- 生产构建: `wails build`
- 生成图标: `wails generate icon`
- 打包应用: `wails build -package`

## 配置文件说明

- `project.json`: Wails 项目配置文件
- [go.mod](file://d:\thinks\other\github\wetools-go\go.mod): Go 模块配置文件
- `package.json`: npm 包配置文件
- `templates/`: 存放 HTML 模板文件
- `wailsjs/`: 自动生成的前端绑定代码

## 代码规范

- 后端代码遵循 Go 语言标准规范
- 前端代码遵循 Vue 3 和 JavaScript 标准规范
- 使用清晰的变量和函数命名
- 添加必要的注释说明

## 调试方法

- 使用 `wails dev` 启动开发服务器进行实时调试
- 使用浏览器开发者工具调试前端代码
- 使用 GoLand 或 VSCode 等 IDE 调试 Go 后端代码
- 查看控制台输出和日志文件定位问题