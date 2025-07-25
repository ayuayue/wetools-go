<template>
  <div class="app-container" :class="{ 'dark-theme': isDarkTheme, 'sidebar-collapsed': isSidebarCollapsed }">
    <!-- Header -->
    <header class="header">
      <h1><i class="fas fa-toolbox"></i> WeTools 开发者工具箱</h1>
      <p>JSON、XML、HTML、Base64、二维码、加密解密、随机字符串等一站式开发工具</p>
      <!-- 主题切换按钮 -->
      <button class="theme-toggle" @click="toggleTheme">
        <i :class="isDarkTheme ? 'fas fa-sun' : 'fas fa-moon'"></i>
      </button>
    </header>

    <div class="container">
      <!-- Sidebar -->
      <Sidebar 
        :menu-items="menuData" 
        @menu-item-click="handleMenuItemClick" 
        :class="{ 'collapsed': isSidebarCollapsed }"
      />
      
      <!-- 菜单收起按钮 -->
      <button class="sidebar-toggle" @click="toggleSidebar">
        <i :class="isSidebarCollapsed ? 'fas fa-arrow-right' : 'fas fa-arrow-left'"></i>
      </button>

      <!-- Main Content -->
      <main class="main-content">
        <ToolHeader 
          :title="currentTool.title" 
          :description="currentTool.description" 
          :icon="currentTool.icon" 
        />
        
        <component 
          :is="currentTool.component" 
          @tool-result="handleToolResult"
        />
      </main>
    </div>

    <!-- Footer -->
    <footer class="footer">
      <p>© 2023 WeTools - 开发者工具箱 | 支持JSON、XML、HTML、Base64、二维码、加密解密等工具</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, shallowRef, onMounted, reactive, onBeforeMount, provide, computed } from 'vue'
import Sidebar from './Sidebar.vue'
import ToolHeader from './ToolHeader.vue'
import JsonTool from './tools/JsonTool.vue'
import XmlTool from './tools/XmlTool.vue'
import HtmlTool from './tools/HtmlTool.vue'
import Base64Tool from './tools/Base64Tool.vue'
import UrlTool from './tools/UrlTool.vue'
import HashTool from './tools/HashTool.vue'
import EncryptTool from './tools/EncryptTool.vue'
import RandomTool from './tools/RandomTool.vue'
import UuidTool from './tools/UuidTool.vue'
import ConverterTool from './tools/ConverterTool.vue'
import QrTool from './tools/QrTool.vue'
import JwtTool from './tools/JwtTool.vue'
import TimestampTool from './tools/TimestampTool.vue'
import TextTool from './tools/TextTool.vue'
import { menuItems } from '../config/menuConfig'

// 主题和菜单状态
const isDarkTheme = ref(false)
const isSidebarCollapsed = ref(false)

// 计算主题
const theme = computed(() => isDarkTheme.value ? 'dark' : 'light')

// 提供主题信息
provide('theme', theme)

// 创建响应式菜单数据
const menuData = reactive(JSON.parse(JSON.stringify(menuItems)))

// 工具组件映射
const toolComponents = {
  json: JsonTool,
  xml: XmlTool,
  html: HtmlTool,
  base64: Base64Tool,
  url: UrlTool,
  hash: HashTool,
  encrypt: EncryptTool,
  random: RandomTool,
  uuid: UuidTool,
  converter: ConverterTool,
  qr: QrTool,
  jwt: JwtTool,
  timestamp: TimestampTool,
  text: TextTool
}

// 当前选中的工具
const currentTool = ref({
  title: 'JSON格式化工具',
  description: '在线JSON格式化、校验、压缩工具，支持JSON转XML、CSV等格式',
  icon: 'fas fa-file-code',
  component: JsonTool
})

// 切换主题
const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value
  // 保存主题设置到本地存储
  localStorage.setItem('theme', isDarkTheme.value ? 'dark' : 'light')
  
  // 应用主题到body
  if (isDarkTheme.value) {
    document.body.classList.add('dark-theme')
  } else {
    document.body.classList.remove('dark-theme')
  }
}

// 切换菜单收起状态
const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
}

// 处理菜单项点击
const handleMenuItemClick = (tool) => {
  // 根据tool.type动态加载对应的组件
  const component = toolComponents[tool.type] || JsonTool
  currentTool.value = {
    title: tool.title,
    description: tool.description,
    icon: tool.icon,
    component: component
  }
}

// 处理工具结果
const handleToolResult = (result) => {
  // 可以在这里处理工具返回的结果
  console.log('Tool result:', result)
}

// 初始化时加载保存的设置
onBeforeMount(() => {
  // 加载主题设置
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDarkTheme.value = savedTheme === 'dark'
    if (isDarkTheme.value) {
      document.body.classList.add('dark-theme')
    }
  }
})

onMounted(() => {
  // 初始化时设置默认工具
  const defaultTool = menuData[0].items[0]
  handleMenuItemClick(defaultTool)
})
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* 亮色主题 */
.app-container {
  background-color: #f5f7fa;
  color: #333;
}

/* 暗色主题 */
.app-container.dark-theme {
  background-color: #1a1a1a;
  color: #e0e0e0;
}

.container {
  display: flex;
  flex: 1;
  min-height: calc(100vh - 80px);
}

.main-content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

/* 暗色主题下的主内容区域 */
.app-container.dark-theme .main-content {
  background-color: #2d2d2d;
}

.footer {
  text-align: center;
  padding: 1.5rem;
  color: #666;
  font-size: 0.9rem;
  border-top: 1px solid #e1e5e9;
  background: white;
}

/* 暗色主题下的页脚 */
.app-container.dark-theme .footer {
  background-color: #222;
  color: #aaa;
  border-top: 1px solid #444;
}

/* 主题切换按钮 */
.theme-toggle {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: background 0.3s;
}

.theme-toggle:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 暗色主题下的主题切换按钮 */
.app-container.dark-theme .theme-toggle {
  background: rgba(255, 255, 255, 0.1);
  color: #e0e0e0;
}

.app-container.dark-theme .theme-toggle:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* 菜单收起按钮 */
.sidebar-toggle {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: #667eea;
  border: none;
  border-radius: 0 5px 5px 0;
  width: 20px;
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: all 0.3s;
  z-index: 100;
}

.sidebar-toggle:hover {
  background: #5a6fd8;
}

/* 展开时的按钮位置 */
.app-container:not(.sidebar-collapsed) .sidebar-toggle {
  left: 180px;
}

/* 收起时的按钮位置 */
.app-container.sidebar-collapsed .sidebar-toggle {
  left: 45px;
}

/* 暗色主题下的菜单收起按钮 */
.app-container.dark-theme .sidebar-toggle {
  background: #444;
}

.app-container.dark-theme .sidebar-toggle:hover {
  background: #555;
}

@media (max-width: 768px) {
  .container {
    flex-direction: column;
  }
  
  .main-content {
    padding: 1rem;
  }
  
  .theme-toggle {
    position: static;
    margin: 0.5rem;
  }
  
  .sidebar-toggle {
    display: none;
  }
}
</style>