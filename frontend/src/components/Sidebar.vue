<template>
  <div class="sidebar" :class="{ 'collapsed': isCollapsed }">
    <div 
      v-for="category in menuItems" 
      :key="category.id" 
      class="sidebar-category"
    >
      <div 
        class="category-header" 
        @click="toggleCategory(category)"
      >
        <h3>
          <i :class="category.icon"></i>
          <span v-if="!isCollapsed">{{ category.title }}</span>
          <i 
            v-if="!isCollapsed"
            :class="[
              'collapse-icon', 
              'fas', 
              category.collapsed ? 'fa-chevron-down' : 'fa-chevron-up'
            ]"
          ></i>
        </h3>
      </div>
      
      <transition name="slide">
      <div 
        v-show="!category.collapsed && !isCollapsed" 
        class="category-items"
      >
        <el-menu
          :default-active="activeItem?.id"
          class="sidebar-menu"
          :collapse="isCollapsed"
        >
          <el-menu-item
            v-for="item in category.items"
            :key="item.id"
            :index="item.id"
            :class="['sidebar-item', { active: isActive(item) }]"
            @click="handleItemClick(item)"
          >
            <i :class="item.icon"></i>
            <span>{{ item.title }}</span>
          </el-menu-item>
        </el-menu>
      </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMenuItem, ElMenu } from 'element-plus'

const props = defineProps({
  menuItems: {
    type: Array,
    required: true
  },
  isCollapsed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['menuItemClick'])

// 使用reactive创建响应式副本
const menuData = reactive(props.menuItems)

// 当前激活的菜单项
const activeItem = ref(props.menuItems[0]?.items[0] || null)

// 切换分类折叠状态
const toggleCategory = (category) => {
  // 直接修改category的collapsed属性
  category.collapsed = !category.collapsed
}

// 处理菜单项点击
const handleItemClick = (item) => {
  // 更新激活状态
  activeItem.value = item
  
  // 发送事件到父组件
  emit('menuItemClick', item)
}

// 检查是否为激活项
const isActive = (item) => {
  return activeItem.value && activeItem.value.id === item.id
}
</script>

<style scoped>
.sidebar {
  width: 100%;
  background: white;
  border-right: 1px solid #e1e5e9;
  padding: 0.5rem 0;
  overflow-y: auto;
  font-size: 0.875rem; /* 14px */
  transition: all 0.3s ease;
  min-height: 100%;
  height: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 暗色主题下的侧边栏 */
.dark-theme .sidebar {
  background: #2d2d2d;
  border-right: 1px solid #444;
}

/* 暗色主题下的菜单项 */
.dark-theme .sidebar-item {
  color: #e0e0e0;
}

.dark-theme .sidebar-item:hover {
  background-color: #3d3d3d;
}

.dark-theme .sidebar-item.active {
  background-color: #3d3d3d;
}

.dark-theme .category-header h3 {
  color: #aaa;
}

.sidebar-category {
  margin-bottom: 0.5rem;
  flex-shrink: 0;
}

/* 最后一个分类使用 flex-grow 填充剩余空间 */
.sidebar-category:last-child {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.sidebar-category:last-child .category-items {
  flex-grow: 1;
}

.category-header {
  cursor: pointer;
  user-select: none;
}

.category-header h3 {
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem; /* 14px */
  color: #666;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.25rem;
  margin: 0;
}

.category-header h3 > i:first-child {
  margin-right: 0.125rem;
  font-size: 0.75rem; /* 12px */
}

.collapse-icon {
  font-size: 0.5rem;
  transition: transform 0.2s;
  margin-left: auto;
}

.sidebar-item {
  display: flex;
  align-items: center;
  padding: 0.375rem 0.75rem 0.375rem 1.25rem;
  color: #333;
  text-decoration: none;
  transition: all 0.2s;
  border-left: 2px solid transparent;
  font-size: 0.875rem; /* 14px */
  gap: 0.25rem;
  margin: 0 0.25rem;
  border-radius: 2px;
}

.sidebar-item:hover {
  background-color: #f8f9fa;
  border-left-color: #667eea;
}

.sidebar-item.active {
  background-color: #eef2ff;
  border-left-color: #667eea;
  font-weight: 500;
}

.category-items {
  transition: all 0.2s ease;
  margin-top: 0.25rem;
  overflow: hidden;
}

/* Element Plus Menu 样式覆盖 */
.sidebar-menu {
  border-right: none !important;
  background: transparent !important;
}

.sidebar-menu :deep(.el-menu-item) {
  display: flex;
  align-items: center;
  padding: 0.375rem 0.75rem 0.375rem 1.25rem !important;
  color: #333;
  text-decoration: none;
  transition: all 0.2s;
  border-left: 2px solid transparent;
  font-size: 0.875rem; /* 14px */
  gap: 0.25rem;
  margin: 0 0.25rem;
  border-radius: 2px;
  height: auto !important;
  line-height: 1.5 !important;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background-color: #f8f9fa !important;
  border-left-color: #667eea;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: #eef2ff !important;
  border-left-color: #667eea;
  font-weight: 500;
}

/* 暗色主题下的菜单项样式 */
.dark-theme .sidebar-menu :deep(.el-menu-item) {
  color: #e0e0e0;
  background: transparent !important;
}

.dark-theme .sidebar-menu :deep(.el-menu-item:hover) {
  background-color: #3d3d3d !important;
}

.dark-theme .sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: #3d3d3d !important;
}

/* 收起的菜单样式 */
.sidebar.collapsed {
  overflow: hidden;
}

/* 菜单展开/收起动画 */
.slide-enter-active {
  transition: all 0.2s ease;
}

.slide-leave-active {
  transition: all 0.2s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.sidebar.collapsed .category-header h3 span {
  display: none;
}

.sidebar.collapsed .category-header {
  justify-content: center;
}

.sidebar.collapsed .category-header h3 {
  padding: 0.5rem;
  justify-content: center;
  overflow: hidden;
}

.sidebar.collapsed .category-header h3 i:first-child {
  margin-right: 0;
  font-size: 0.875rem; /* 14px */
}

@media (max-width: 768px) {
  .sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e1e5e9;
    max-height: 200px;
  }
}
</style>