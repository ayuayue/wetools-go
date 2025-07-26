<template>
  <el-container class="main-layout-container" direction="horizontal">
    <!-- Sidebar -->
    <div class="sidebar-container">
      <el-aside 
        class="sidebar-wrapper" 
        :class="{ 'collapsed': isSidebarCollapsed }"
        :style="{ width: isSidebarCollapsed ? '64px' : sidebarWidth + 'px' }"
      >
        <Sidebar 
          :menu-items="menuData" 
          @menu-item-click="handleMenuItemClick" 
          :is-collapsed="isSidebarCollapsed"
        />
      </el-aside>
      
      <!-- 菜单拉伸手柄 -->
      <div 
        v-if="!isSidebarCollapsed"
        class="sidebar-resize-handle" 
        @mousedown="startResize"
      ></div>
      
      <!-- 菜单收起按钮 -->
      <div class="sidebar-toggle-wrapper">
        <el-button 
          class="sidebar-toggle" 
          :icon="isSidebarCollapsed ? 'ArrowRight' : 'ArrowLeft'" 
          circle 
          @click="toggleSidebar"
        />
      </div>
    </div>

    <!-- Main Content -->
    <el-main class="main-content">
      <ToolHeader />
      
      <component 
        :is="currentTool.component" 
        @tool-result="handleToolResult"
      />
    </el-main>
  </el-container>
  
  <!-- Footer -->
  <el-footer class="footer">
    <p>© 2023 WeTools - 开发者工具箱 | 支持JSON、XML、HTML、Base64、二维码、加密解密等工具</p>
  </el-footer>
</template>

<script setup>
import { ref, shallowRef, onMounted, reactive, onBeforeMount, provide, computed } from 'vue'
import { ElContainer, ElAside, ElMain, ElFooter, ElButton } from 'element-plus'
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
import { menuItems } from '../config/menuConfig'

// 主题和菜单状态
const isDarkTheme = ref(false)
const isSidebarCollapsed = ref(false)
const sidebarWidth = ref(220)

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
  qr: QrTool
}

// 当前选中的工具
const currentTool = ref({
  title: 'JSON格式化工具',
  description: '',
  icon: 'fas fa-file-code',
  component: JsonTool
})

// 切换菜单收起状态
const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
}

// 开始调整侧边栏宽度
const startResize = (e) => {
  e.preventDefault()
  const startX = e.clientX
  const startWidth = sidebarWidth.value
  
  const doDrag = (e) => {
    const newWidth = startWidth + (e.clientX - startX)
    sidebarWidth.value = Math.max(180, Math.min(400, newWidth))
  }
  
  const stopDrag = () => {
    document.removeEventListener('mousemove', doDrag)
    document.removeEventListener('mouseup', stopDrag)
  }
  
  document.addEventListener('mousemove', doDrag)
  document.addEventListener('mouseup', stopDrag)
}

// 处理菜单项点击
const handleMenuItemClick = (tool) => {
  // 根据tool.type动态加载对应的组件
  const component = toolComponents[tool.type] || JsonTool
  currentTool.value = {
    title: tool.title,
    description: tool.description || '',
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
.main-layout-container {
  flex: 1;
  height: calc(100vh - 120px);
  overflow: hidden;
}

.sidebar-container {
  position: relative;
  height: 100%;
  display: flex;
}

.sidebar-wrapper {
  transition: all 0.3s ease;
  background: white;
  border-right: 1px solid #e1e5e9;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
}

.main-content {
  padding: 2rem;
  overflow: hidden;
  background-color: #f5f7fa;
  height: 100%;
}

.footer {
  text-align: center;
  padding: 1.5rem;
  color: #666;
  font-size: 0.9rem;
  border-top: 1px solid #e1e5e9;
  background: white;
}

/* 暗色主题 */
.dark-theme .sidebar-wrapper {
  background: #2d2d2d;
  border-right: 1px solid #444;
}

.dark-theme .main-content {
  background-color: #2d2d2d;
}

.dark-theme .footer {
  background-color: #222;
  color: #aaa;
  border-top: 1px solid #444;
}

/* 菜单拉伸手柄 */
.sidebar-resize-handle {
  width: 4px;
  height: 100%;
  background: #667eea;
  cursor: col-resize;
  position: absolute;
  right: 0;
  top: 0;
  z-index: 100;
  opacity: 0;
  transition: opacity 0.2s;
}

.sidebar-resize-handle:hover {
  opacity: 1;
}

/* 菜单收起按钮包装器 */
.sidebar-toggle-wrapper {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: -10px;
  z-index: 100;
}

.sidebar-toggle {
  background: #667eea !important;
  border: none;
  color: white !important;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.sidebar-toggle:hover {
  background: #5a6fd8 !important;
}

/* 暗色主题下的菜单收起按钮 */
.dark-theme .sidebar-toggle {
  background: #444 !important;
}

.dark-theme .sidebar-toggle:hover {
  background: #555 !important;
}

@media (max-width: 768px) {
  .sidebar-wrapper {
    width: 100% !important;
    max-height: 200px;
    border-right: none;
    border-bottom: 1px solid #e1e5e9;
  }
  
  .main-content {
    padding: 1rem;
  }
  
  .sidebar-toggle {
    display: none;
  }
  
  .collapsed + .sidebar-toggle-wrapper .sidebar-toggle {
    display: none;
  }
}
</style>