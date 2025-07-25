<template>
  <div class="sidebar">
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
          <i :class="category.icon"></i> {{ category.title }}
          <i 
            :class="[
              'collapse-icon', 
              'fas', 
              category.collapsed ? 'fa-chevron-down' : 'fa-chevron-up'
            ]"
          ></i>
        </h3>
      </div>
      
      <div 
        v-show="!category.collapsed" 
        class="category-items"
      >
        <a
          v-for="item in category.items"
          :key="item.id"
          href="#"
          :class="['sidebar-item', { active: isActive(item) }]"
          @click="handleItemClick(item, $event)"
        >
          <i :class="item.icon"></i> {{ item.title }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  menuItems: {
    type: Array,
    required: true
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
const handleItemClick = (item, event) => {
  event.preventDefault()
  
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
  width: 250px;
  background: white;
  border-right: 1px solid #e1e5e9;
  padding: 1rem 0;
  overflow-y: auto;
}

.sidebar-category {
  margin-bottom: 1rem;
}

.category-header {
  cursor: pointer;
  user-select: none;
}

.category-header h3 {
  padding: 0.5rem 1.5rem;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #666;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.collapse-icon {
  font-size: 0.7rem;
  transition: transform 0.2s;
}

.sidebar-item {
  display: block;
  padding: 0.75rem 1.5rem 0.75rem 2.5rem;
  color: #333;
  text-decoration: none;
  transition: all 0.2s;
  border-left: 3px solid transparent;
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
  transition: all 0.3s ease;
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