<template>
  <el-card class="tool-header" shadow="never">
    <div class="header-content">
    </div>
  </el-card>
</template>

<script setup>
import { ElCard } from 'element-plus'

// 不再需要props，因为标题和描述已被移除
</script>

<style scoped>
.tool-header {
  margin-bottom: 1.5rem;
  /* background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); */
  color: white;
  border: none;
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.title {
  font-size: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
}

.description {
  font-size: 1rem;
  opacity: 0.9;
  margin: 0;
}

/* 暗色主题 */
.dark-theme .tool-header {
  background: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
}

@media (max-width: 768px) {
  .title {
    font-size: 1.25rem;
  }
  
  .description {
    font-size: 0.875rem;
  }
}
</style>