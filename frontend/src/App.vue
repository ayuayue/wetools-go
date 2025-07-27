<template>
  <el-container class="app-container" :class="{ 'dark-theme': isDarkTheme }">
    <el-header class="app-header">
      <div class="header-content">
        <h1><i class="fas fa-toolbox"></i> WeTools 开发者工具箱</h1>
      </div>
      <div class="header-actions">
        <el-button 
          class="theme-toggle" 
          :icon="isDarkTheme ? 'Sunny' : 'Moon'" 
          circle 
          @click="toggleTheme"
        />
      </div>
    </el-header>
    
    <MainLayout />
  </el-container>
</template>

<script setup>
import { ref, onBeforeMount } from 'vue'
import { ElContainer, ElHeader, ElButton } from 'element-plus'
import MainLayout from './components/MainLayout.vue'
import { toggleTheme as switchTheme, initTheme } from './config/elementPlusConfig'

// 主题状态
const isDarkTheme = ref(false)

// 切换主题
const toggleTheme = () => {
  switchTheme()
  isDarkTheme.value = !isDarkTheme.value
}

// 初始化时加载保存的设置
onBeforeMount(() => {
  // 初始化主题
  initTheme()
  
  // 加载主题设置
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDarkTheme.value = savedTheme === 'dark'
  }
})
</script>

<style>
/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f7fa;
  color: #333;
}

.app-container {
  min-height: 100vh;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.75rem 1rem;
  box-shadow: 0 1px 5px rgba(0,0,0,0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: none;
}

.header-content h1 {
  font-size: 1.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
  margin: 0;
}

.header-content p {
  margin-top: 0.25rem;
  font-size: 0.75rem;
  opacity: 0.9;
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
}

.theme-toggle {
  background: rgba(255, 255, 255, 0.2) !important;
  border: none;
  color: white !important;
}

.theme-toggle:hover {
  background: rgba(255, 255, 255, 0.3) !important;
}

/* 暗色主题 */
.app-container.dark-theme .app-header {
  background: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
}

.app-container.dark-theme .theme-toggle {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #e0e0e0 !important;
}

.app-container.dark-theme .theme-toggle:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

/* 从原index.html中提取的样式 */
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}

@media (max-width: 768px) {
  .app-header {
    flex-direction: column;
    text-align: center;
    padding: 1rem;
  }
  
  .header-content {
    margin-bottom: 1rem;
  }
  
  .header-actions {
    align-self: flex-end;
  }
}
</style>