# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Wails desktop application project that combines Go (backend) with Vue.js (frontend) to create cross-platform desktop applications. The project is named "WeTools" and uses Wails v2 framework.

## Code Architecture

### Backend (Go)
- `main.go`: Entry point that initializes the Wails application, embeds frontend assets, and configures the window
- `app.go`: Contains the App struct with Go methods that can be called from the frontend
- The Go backend uses Wails runtime to communicate with the Vue frontend

### Frontend (Vue.js)
- Single page application with Vue 3 Composition API
- `App.vue`: Main application component that uses the MainLayout component
- Component structure:
  - `MainLayout.vue`: Main layout component containing header, sidebar, and main content
  - `Sidebar.vue`: Collapsible sidebar menu with category support
  - `ToolHeader.vue`: Dynamic tool header component
  - Tools components in `components/tools/`:
    - `JsonTool.vue`: JSON formatting tool
    - `XmlTool.vue`: XML formatting tool
    - `HtmlTool.vue`: HTML formatting tool
    - `Base64Tool.vue`: Base64 encoding/decoding tool
    - `UrlTool.vue`: URL encoding/decoding tool
    - `HashTool.vue`: Hash calculation tool
    - `EncryptTool.vue`: Encryption/decryption tool
    - `RandomTool.vue`: Random string generator
    - `UuidTool.vue`: UUID generator
    - `ConverterTool.vue`: Format converter
    - `QrTool.vue`: QR code generator
- Configuration:
  - `config/menuConfig.js`: Menu structure and configuration
- Assets are located in the `assets` directory (images, fonts)

### Build System
- Wails handles the build process, embedding the frontend dist into the Go binary
- Frontend uses Vite as the build tool with npm for package management

## Common Development Commands

### Development
- `wails dev`: Run in live development mode with hot reload for frontend changes
- `wails build`: Build a redistributable production package

### Frontend Development
- `npm install`: Install frontend dependencies (run from frontend directory)
- `npm run dev`: Run Vite development server (run from frontend directory)
- `npm run build`: Build frontend assets (run from frontend directory)

## Project Configuration
- `wails.json`: Main Wails project configuration
- `go.mod`: Go module dependencies
- `frontend/package.json`: Frontend dependencies and scripts

## Communication Between Frontend and Backend
- Go methods in `app.go` are exposed to frontend via Wails bindings
- Frontend calls Go functions using the generated bindings in `wailsjs/go/`
- The `Greet` function in App struct is an example of backend functionality exposed to frontend

## Component Architecture
- MainLayout: Root layout component
- Sidebar: Collapsible menu with category support
- ToolHeader: Dynamic header for tools
- Tools: Individual tool components (JSON, XML, etc.)
- Menu configuration is centralized in config/menuConfig.js
- All tools support copy to clipboard functionality
- All tools have responsive design for mobile devices