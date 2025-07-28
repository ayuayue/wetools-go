<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>代理设置</h2>
      <p>配置网络代理以通过代理服务器访问网络资源</p>
    </div>
    
    <div class="tool-content">
      <ProxyConfig 
        :visible="proxyConfigVisible" 
        @update:visible="proxyConfigVisible = $event"
        @proxy-config-updated="handleProxyConfigUpdated"
      />
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElCard } from 'element-plus'
import ProxyConfig from '../ProxyConfig.vue'
import proxyManager from '../../utils/proxyManager'

// 代理配置对话框可见性
const proxyConfigVisible = ref(true)

// 处理代理配置更新
const handleProxyConfigUpdated = (config) => {
  console.log('代理配置已更新:', config)
  // 更新全局代理配置
  proxyManager.updateProxyConfig(config)
}
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
}

.tool-header {
  margin-bottom: 1.5rem;
}

.tool-header h2 {
  margin: 0 0 0.5rem 0;
  color: #333;
}

.tool-header p {
  margin: 0;
  color: #666;
}

.tool-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
</style>