<template>
  <div class="app-container">
    <!-- Header -->
    <header class="header">
      <h1><i class="fas fa-toolbox"></i> WeTools 开发者工具箱</h1>
      <p>JSON、XML、HTML、Base64、二维码、加密解密、随机字符串等一站式开发工具</p>
    </header>

    <div class="container">
      <!-- Sidebar -->
      <Sidebar :menu-items="menuData" @menu-item-click="handleMenuItemClick" />

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
import { ref, shallowRef, onMounted, reactive } from 'vue'
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
  description: '在线JSON格式化、校验、压缩工具，支持JSON转XML、CSV等格式',
  icon: 'fas fa-file-code',
  component: JsonTool
})

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

.footer {
  text-align: center;
  padding: 1.5rem;
  color: #666;
  font-size: 0.9rem;
  border-top: 1px solid #e1e5e9;
  background: white;
}

@media (max-width: 768px) {
  .container {
    flex-direction: column;
  }
  
  .main-content {
    padding: 1rem;
  }
}
</style>