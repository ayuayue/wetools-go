<template>
  <el-breadcrumb separator="/" class="tool-breadcrumb">
    <el-breadcrumb-item 
      v-for="(item, index) in breadcrumbItems" 
      :key="index"
      :to="item.path"
    >
      {{ item.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
import { ElBreadcrumb, ElBreadcrumbItem } from 'element-plus'

// 定义props
const props = defineProps({
  breadcrumbItems: {
    type: Array,
    default: () => []
  }
})
</script>

<style scoped>
.tool-breadcrumb {
  margin-bottom: 0.5rem;
  padding: 0.25rem 0.5rem;
  background: white;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

/* 暗色主题 */
.dark-theme .tool-breadcrumb {
  background: #2d2d2d;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}
</style>