<template>
  <div class="sidebar-wrapper">
    <el-menu
      :default-active="activeItem?.id"
      class="sidebar-menu"
      :collapse="isCollapsed"
      :unique-opened="false"
      :default-openeds="getDefaultOpeneds()"
      @select="handleMenuSelect"
    >
      <el-sub-menu 
        v-for="category in menuItems" 
        :key="category.id" 
        :index="category.id"
      >
        <template #title>
          <i :class="category.icon"></i>
          <span>{{ category.title }}</span>
        </template>
        <el-menu-item
          v-for="item in category.items"
          :key="item.id"
          :index="item.id"
        >
          <i :class="item.icon"></i>
          <template #title>
            <span>{{ item.title }}</span>
          </template>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMenuItem, ElMenu, ElSubMenu } from 'element-plus'

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

// 当前激活的菜单项
const activeItem = ref(null)

// 当菜单数据加载完成后设置默认激活项
onMounted(() => {
  if (props.menuItems && props.menuItems.length > 0) {
    const firstCategory = props.menuItems[0]
    if (firstCategory && firstCategory.items && firstCategory.items.length > 0) {
      activeItem.value = firstCategory.items[0]
    }
  }
})

// 获取默认展开的菜单项
const getDefaultOpeneds = () => {
  if (!props.menuItems) return []
  return props.menuItems.map(category => category.id)
}

// 处理菜单选择
const handleMenuSelect = (index) => {
  // 查找选中的菜单项
  for (const category of props.menuItems) {
    const selectedItem = category.items.find(item => item.id === index)
    if (selectedItem) {
      activeItem.value = selectedItem
      emit('menuItemClick', selectedItem)
      break
    }
  }
}
</script>

<style scoped>
.sidebar-wrapper {
  height: 100%;
  background: white;
  border-right: 1px solid #e1e5e9;
}

.sidebar-menu {
  border-right: none;
  height: 100%;
}

.sidebar-menu :deep(.el-sub-menu__title) {
  font-size: 0.875rem;
  height: 40px;
  line-height: 40px;
  display: flex;
  align-items: center;
}

.sidebar-menu :deep(.el-menu-item) {
  font-size: 0.875rem;
  height: 36px;
  line-height: 36px;
  display: flex;
  align-items: center;
}

.sidebar-menu :deep(.el-menu-item i),
.sidebar-menu :deep(.el-sub-menu__title i) {
  margin-right: 8px;
  width: 16px;
  text-align: center;
}

/* 暗色主题 */
.dark-theme .sidebar-wrapper {
  background: #2d2d2d;
  border-right: 1px solid #444;
}

.dark-theme .sidebar-menu :deep(.el-sub-menu__title) {
  color: #e0e0e0;
}

.dark-theme .sidebar-menu :deep(.el-sub-menu__title:hover) {
  background-color: #3d3d3d;
}

.dark-theme .sidebar-menu :deep(.el-menu-item) {
  color: #e0e0e0;
}

.dark-theme .sidebar-menu :deep(.el-menu-item:hover) {
  background-color: #3d3d3d;
}

.dark-theme .sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: #3d3d3d;
  color: #667eea;
}

/* 收起状态下的样式 */
.sidebar-menu:not(.el-menu--collapse) {
  width: 200px;
}
</style>