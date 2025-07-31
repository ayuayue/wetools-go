<template>
  <el-container class="main-layout-container">
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
    </div>
    
    <!-- 菜单收起按钮 -->
    <div class="sidebar-toggle-wrapper" :style="{ left: isSidebarCollapsed ? '64px' : sidebarWidth + 'px' }">
      <el-button 
        class="sidebar-toggle" 
        :icon="isSidebarCollapsed ? 'ArrowRight' : 'ArrowLeft'" 
        circle 
        @click="toggleSidebar"
      />
    </div>

    <!-- Main Content and Footer -->
    <el-container direction="vertical" class="content-container">
      <el-main class="main-content">
        <ToolHeader :breadcrumb-items="breadcrumbItems" />
        
        <!-- 标签页 -->
        <el-tabs 
          v-model="activeTab" 
          type="card" 
          @tab-remove="removeTab"
          @tab-click="switchTab"
          class="tool-tabs"
        >
          <el-tab-pane
            v-for="tab in openTabs"
            :key="tab.id"
            :name="tab.id"
          >
            <template #label>
              <div class="tab-label" @contextmenu.prevent="showContextMenu(tab.id, $event)">
                <span>{{ tab.title }}</span>
                <el-icon 
                  v-if="tab.pinned" 
                  class="pin-icon pinned" 
                  @click.stop="togglePinTab(tab.id)"
                >
                  <Lock />
                </el-icon>
                <el-icon 
                  v-else 
                  class="pin-icon" 
                  @click.stop="togglePinTab(tab.id)"
                >
                  <Unlock />
                </el-icon>
                <el-icon 
                  v-if="!tab.pinned" 
                  class="close-icon" 
                  @click.stop="removeTab(tab.id)"
                >
                  <Close />
                </el-icon>
              </div>
            </template>
            <component 
              :is="tab.component" 
              @tool-result="handleToolResult"
            />
          </el-tab-pane>
        </el-tabs>
        
        <!-- 右键菜单 -->
        <div 
          v-show="contextMenuVisible" 
          class="context-menu" 
          :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }"
          @blur="contextMenuVisible = false"
          tabindex="-1"
        >
          <ul>
            <li @click="closeAllTabs">关闭所有标签页</li>
            <li @click="closeOtherTabs">关闭其他标签页</li>
            <li @click="closeLeftTabs">关闭左侧标签页</li>
            <li @click="closeRightTabs">关闭右侧标签页</li>
          </ul>
        </div>
      </el-main>
      
      <!-- Footer -->
      <el-footer class="footer">
        <p>© 2023 WeTools - 开发者工具箱 | 支持JSON、XML、HTML、Base64、二维码、加密解密等工具</p>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, shallowRef, onMounted, reactive, onBeforeMount, provide, computed, h } from 'vue'
import { ElContainer, ElAside, ElMain, ElFooter, ElButton, ElTabs, ElTabPane, ElIcon, ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'
import { Close, Lock, Unlock } from '@element-plus/icons-vue'
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
import WebsiteTool from './tools/WebsiteTool.vue'
import ProxyTool from './tools/ProxyTool.vue'
import SettingsTool from './tools/SettingsTool.vue'
import HtmlRunner from './tools/HtmlRunner.vue'
import JsRunner from './tools/JsRunner.vue'
import JwtTool from './tools/JwtTool.vue'
import AboutTool from './tools/AboutTool.vue'
import { menuItems } from '../config/menuConfig'

// 主题和菜单状态
const isDarkTheme = ref(false)
const isSidebarCollapsed = ref(false)
const sidebarWidth = ref(180)

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
  website: WebsiteTool,
  proxy: ProxyTool,
  settings: SettingsTool,
  htmlrunner: HtmlRunner,
  jsrunner: JsRunner,
  jwt: JwtTool,
  about: AboutTool
}

// 标签页管理
const openTabs = ref([])
const activeTab = ref('')

// 右键菜单相关
const contextMenuVisible = ref(false)
const contextMenuTabId = ref('')
const contextMenuX = ref(0)
const contextMenuY = ref(0)

// 面包屑数据
const breadcrumbItems = ref([
  { title: '首页', path: '/' }
])

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
    sidebarWidth.value = Math.max(150, Math.min(400, newWidth))
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
  
  // 检查是否是彩蛋标签页，如果是则移除
  const easterEggTabIndex = openTabs.value.findIndex(tab => tab.id === 'easter-egg')
  if (easterEggTabIndex !== -1) {
    openTabs.value.splice(easterEggTabIndex, 1)
  }
  
  // 检查标签页是否已存在
  const existingTab = openTabs.value.find(tab => tab.id === tool.id)
  
  if (!existingTab) {
    // 添加新标签页
    const newTab = {
      id: tool.id,
      title: tool.title,
      description: tool.description || '',
      icon: tool.icon,
      component: component,
      tool: tool,
      pinned: false // 默认不固定
    }
    openTabs.value.push(newTab)
  }
  
  // 激活标签页
  activeTab.value = tool.id
  
  // 更新当前工具
  currentTool.value = {
    title: tool.title,
    description: tool.description || '',
    icon: tool.icon,
    component: component
  }
  
  // 更新面包屑
  updateBreadcrumb(tool)
}

// 切换标签页固定状态
const togglePinTab = (tabId) => {
  const tab = openTabs.value.find(tab => tab.id === tabId)
  if (tab) {
    tab.pinned = !tab.pinned
  }
}

// 显示右键菜单
const showContextMenu = (tabId, event) => {
  event.preventDefault()
  contextMenuTabId.value = tabId
  contextMenuX.value = event.clientX
  contextMenuY.value = event.clientY
  contextMenuVisible.value = true
  
  // 确保菜单获得焦点，以便可以监听blur事件
  setTimeout(() => {
    const menu = document.querySelector('.context-menu')
    if (menu) {
      menu.focus()
    }
  }, 0)
}

// 关闭右侧标签页
const closeRightTabs = () => {
  const tabIndex = openTabs.value.findIndex(tab => tab.id === contextMenuTabId.value)
  if (tabIndex !== -1) {
    // 保留左侧和当前标签页，过滤掉右侧标签页（固定标签页除外）
    openTabs.value = openTabs.value.filter((tab, index) => {
      if (index <= tabIndex) return true
      if (tab.pinned) return true
      return false
    })
    
    // 如果当前激活的标签页被关闭了，激活最后一个标签页
    if (!openTabs.value.some(tab => tab.id === activeTab.value)) {
      activeTab.value = openTabs.value[openTabs.value.length - 1]?.id || ''
    }
    
    // 如果没有标签页了，添加默认标签页
    if (openTabs.value.length === 0) {
      const defaultTool = menuData[0].items[0]
      handleMenuItemClick(defaultTool)
    }
  }
}

// 关闭左侧标签页
const closeLeftTabs = () => {
  const tabIndex = openTabs.value.findIndex(tab => tab.id === contextMenuTabId.value)
  if (tabIndex !== -1) {
    // 保留右侧和当前标签页，过滤掉左侧标签页（固定标签页除外）
    openTabs.value = openTabs.value.filter((tab, index) => {
      if (index >= tabIndex) return true
      if (tab.pinned) return true
      return false
    })
    
    // 如果当前激活的标签页被关闭了，激活第一个标签页
    if (!openTabs.value.some(tab => tab.id === activeTab.value)) {
      activeTab.value = openTabs.value[0]?.id || ''
    }
    
    // 如果没有标签页了，添加默认标签页
    if (openTabs.value.length === 0) {
      const defaultTool = menuData[0].items[0]
      handleMenuItemClick(defaultTool)
    }
  }
}

// 关闭其他标签页
const closeOtherTabs = () => {
  // 只保留当前标签页和固定标签页
  openTabs.value = openTabs.value.filter(tab => {
    if (tab.id === contextMenuTabId.value) return true
    if (tab.pinned) return true
    return false
  })
  
  // 如果当前激活的标签页被关闭了，激活右键点击的标签页
  if (!openTabs.value.some(tab => tab.id === activeTab.value)) {
    activeTab.value = contextMenuTabId.value
  }
}

// 关闭所有标签页
const closeAllTabs = () => {
  // 只保留固定标签页
  openTabs.value = openTabs.value.filter(tab => tab.pinned)
  
  // 如果没有标签页了，添加默认标签页
  if (openTabs.value.length === 0) {
    const defaultTool = menuData[0].items[0]
    handleMenuItemClick(defaultTool)
  } else {
    // 激活第一个标签页
    activeTab.value = openTabs.value[0].id
  }
}

// 更新面包屑
const updateBreadcrumb = (tool) => {
  // 找到工具所属的分类
  const category = menuData.find(cat => cat.items.some(item => item.id === tool.id))
  
  breadcrumbItems.value = [
    { title: '首页', path: '/' },
    { title: category ? category.title : '工具', path: '#' },
    { title: tool.title, path: '#' }
  ]
}

// 移除标签页
const removeTab = (targetName) => {
  const tabs = openTabs.value
  let activeName = activeTab.value
  
  // 检查是否是固定标签页
  const tabToRemove = tabs.find(tab => tab.id === targetName)
  if (tabToRemove && tabToRemove.pinned) {
    // 固定的标签页不能被移除
    return
  }
  
  if (activeName === targetName) {
    tabs.forEach((tab, index) => {
      if (tab.id === targetName) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          activeName = nextTab.id
        }
      }
    })
  }
  
  activeTab.value = activeName
  openTabs.value = tabs.filter(tab => tab.id !== targetName)
  
  // 如果没有标签页了，显示默认工具而不是彩蛋
  if (openTabs.value.length === 0) {
    const defaultTool = menuData[0].items[0]
    handleMenuItemClick(defaultTool)
  }
}

// 显示彩蛋内容
const showEasterEgg = () => {
  // 创建一个特殊的彩蛋标签页
  const easterEggTab = {
    id: 'easter-egg',
    title: '🎉 惊喜彩蛋',
    description: '',
    icon: 'fas fa-gift',
    component: {
      render() {
        return h('div', { class: 'easter-egg-container' }, [
          h('div', { class: 'easter-egg-content' }, [
            h('h2', { class: 'easter-egg-title' }, '🎉 恭喜你发现彩蛋！'),
            h('p', { class: 'easter-egg-text' }, '感谢使用 WeTools 开发者工具箱！'),
            h('p', { class: 'easter-egg-text' }, '这个工具箱旨在帮助开发者提高工作效率。'),
            h('p', { class: 'easter-egg-text' }, '如果你喜欢这个工具，欢迎推荐给更多的朋友！'),
            h('div', { class: 'easter-egg-icon' }, '🛠️'),
            h('p', { class: 'easter-egg-text small' }, '继续探索，也许还有更多惊喜等着你...')
          ])
        ])
      }
    },
    tool: {},
    pinned: false
  }
  
  openTabs.value.push(easterEggTab)
  activeTab.value = 'easter-egg'
}

// 切换标签页
const switchTab = (tab) => {
  const tool = openTabs.value.find(t => t.id === tab.props.name)
  if (tool) {
    currentTool.value = {
      title: tool.title,
      description: tool.description || '',
      icon: tool.icon,
      component: tool.component
    }
    updateBreadcrumb(tool.tool)
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
  // 初始化时设置默认工具，不显示彩蛋
  const defaultTool = menuData[0].items[0]
  handleMenuItemClick(defaultTool)
})
</script>

<style scoped>
.main-layout-container {
  flex: 1;
  height: 100vh;
  overflow: hidden;
  display: flex;
  position: relative;
}

.content-container {
  flex: 1;
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
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
  padding: 1rem;
  overflow-y: auto;
  background-color: #f5f7fa;
  flex: 1;
}

.footer {
  text-align: center;
  padding: 0.5rem;
  color: #666;
  font-size: 0.75rem;
  border-top: 1px solid #e1e5e9;
  background: white;
  flex-shrink: 0;
  height: 35px; /* 设置固定高度 */
}

/* 标签页样式 */
.tool-tabs {
  margin-top: 1rem;
}

.tool-tabs :deep(.el-tabs__header) {
  margin-bottom: 1rem;
}

.tool-tabs :deep(.el-tabs__content) {
  min-height: 300px;
}

/* 标签页标签样式 */
.tab-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.pin-icon {
  font-size: 0.8rem;
  cursor: pointer;
  color: #999;
  transition: color 0.2s;
}

.pin-icon:hover {
  color: #409eff;
}

.pin-icon.pinned {
  color: #409eff;
}

.close-icon {
  font-size: 0.8rem;
  cursor: pointer;
  color: #999;
  transition: color 0.2s;
}

.close-icon:hover {
  color: #f56c6c;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: white;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  z-index: 3000;
  min-width: 160px;
}

.context-menu ul {
  list-style: none;
  padding: 6px 0;
  margin: 0;
}

.context-menu li {
  padding: 8px 16px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
}

.context-menu li:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

/* 暗色主题下的右键菜单 */
.dark-theme .context-menu {
  background: #2d2d2d;
  border: 1px solid #444;
}

.dark-theme .context-menu li {
  color: #aaa;
}

.dark-theme .context-menu li:hover {
  background-color: #3d3d3d;
  color: #409eff;
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
  z-index: 100;
  transition: left 0.3s ease;
  width: 0;
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

/* 彩蛋内容样式 */
.easter-egg-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding: 2rem;
}

.easter-egg-content {
  text-align: center;
  max-width: 600px;
  padding: 2rem;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4edf9 100%);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.easter-egg-title {
  color: #667eea;
  margin-bottom: 1.5rem;
  font-size: 2rem;
}

.easter-egg-text {
  color: #555;
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.easter-egg-text.small {
  font-size: 0.9rem;
  color: #888;
}

.easter-egg-icon {
  font-size: 4rem;
  margin: 1.5rem 0;
}

/* 暗色主题下的彩蛋样式 */
.dark-theme .easter-egg-content {
  background: linear-gradient(135deg, #2d3748 0%, #222831 100%);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.dark-theme .easter-egg-title {
  color: #667eea;
}

.dark-theme .easter-egg-text {
  color: #e0e0e0;
}

.dark-theme .easter-egg-text.small {
  color: #aaa;
}

@media (max-width: 768px) {
  .sidebar-wrapper {
    width: 100% !important;
    max-height: 150px;
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