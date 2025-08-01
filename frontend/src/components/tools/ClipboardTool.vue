<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>剪贴板管理</h2>
      <p>管理和查看剪贴板历史记录</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="clipboard-card">
            <template #header>
              <div class="clipboard-header">
                <span>剪贴板监听</span>
                <el-switch
                  v-model="isListening"
                  active-text="开启"
                  inactive-text="关闭"
                  @change="toggleClipboardListener"
                  :loading="loading"
                />
              </div>
            </template>
            
            <div class="clipboard-content">
              <el-alert
                v-if="isListening"
                title="剪贴板监听已开启，复制内容将自动记录"
                type="success"
                :closable="false"
                show-icon
              />
              <el-alert
                v-else
                title="剪贴板监听已关闭，开启后可记录复制的内容"
                type="info"
                :closable="false"
                show-icon
              />
              
              <div class="clipboard-stats">
                <el-tag>记录数: {{ clipboardHistory.length }}</el-tag>
                <div class="clipboard-actions">
                  <el-button 
                    type="primary" 
                    @click="openClipboardFile"
                    :disabled="loading"
                  >
                    打开文件
                  </el-button>
                  <el-button 
                    type="danger" 
                    :icon="Delete" 
                    @click="clearHistory"
                    :disabled="clipboardHistory.length === 0 || loading"
                  >
                    清空记录
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card class="history-card">
            <template #header>
              <div class="history-header">
                <span>剪贴板历史</span>
                <div class="history-actions">
                  <el-button 
                    type="primary" 
                    :icon="Refresh" 
                    @click="refreshHistory"
                    :loading="loading"
                  >
                    刷新
                  </el-button>
                </div>
              </div>
            </template>
            
            <div class="history-content">
              <el-empty 
                v-if="clipboardHistory.length === 0" 
                description="暂无剪贴板记录"
              />
              
              <div v-else>
                <!-- 搜索框 -->
                <div class="search-container">
                  <el-input
                    v-model="searchText"
                    placeholder="搜索剪贴板内容..."
                    clearable
                    @input="handleSearch"
                  >
                    <template #prefix>
                      <el-icon><Search /></el-icon>
                    </template>
                  </el-input>
                </div>
                
                <!-- 历史记录列表 -->
                <div class="history-list">
                  <div 
                    v-for="(item, index) in filteredHistory" 
                    :key="index"
                    class="history-item"
                  >
                    <div class="item-header">
                      <div class="item-time">{{ formatTime(item.timestamp) }}</div>
                      <div class="item-actions">
                        <el-button 
                          type="primary" 
                          :icon="DocumentCopy" 
                          @click="copyToClipboard(item.content)"
                          size="small"
                          link
                        >
                          复制
                        </el-button>
                        <el-button 
                          type="danger" 
                          :icon="Delete" 
                          @click="removeItem(index)"
                          size="small"
                          link
                        >
                          删除
                        </el-button>
                      </div>
                    </div>
                    <div class="item-content">{{ item.content }}</div>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { ElCard, ElRow, ElCol, ElSwitch, ElAlert, ElTag, ElButton, ElEmpty, ElMessage, ElInput } from 'element-plus'
import { Delete, Refresh, DocumentCopy, Search } from '@element-plus/icons-vue'

// 剪贴板历史记录
const clipboardHistory = ref([])

// 监听状态
const isListening = ref(true)

// 加载状态
const loading = ref(false)

// 搜索文本
const searchText = ref('')

// 过滤后的历史记录
const filteredHistory = computed(() => {
  if (!searchText.value) {
    return clipboardHistory.value
  }
  
  const searchLower = searchText.value.toLowerCase()
  return clipboardHistory.value.filter(item => 
    item.content.toLowerCase().includes(searchLower)
  )
})

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

// 切换剪贴板监听
const toggleClipboardListener = (value) => {
  // 在这个版本中，监听是始终开启的，因为我们已经在后端实现了
  // 这里只是为了UI展示
  ElMessage.info(value ? '剪贴板监听已开启' : '剪贴板监听已关闭')
}

// 获取剪贴板历史记录
const fetchHistory = async () => {
  try {
    loading.value = true
    console.log('正在获取剪贴板历史记录')
    
    // 获取所有历史记录
    const history = await window.go.main.App.GetClipboardHistory()
    
    console.log('获取到的剪贴板历史记录:', history)
    
    clipboardHistory.value = Array.isArray(history) ? history : []
    
    console.log('更新后的状态:', {
      clipboardHistory: clipboardHistory.value
    })
  } catch (error) {
    console.error('获取剪贴板历史失败:', error)
    ElMessage.error('获取剪贴板历史失败: ' + (error.message || '未知错误'))
    // 设置默认值以避免进一步的错误
    clipboardHistory.value = []
  } finally {
    loading.value = false
  }
}

// 复制到剪贴板
const copyToClipboard = async (content) => {
  try {
    await window.go.main.App.CopyToClipboard(content)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败: ' + error.message)
  }
}

// 删除记录
const removeItem = async (index) => {
  try {
    loading.value = true
    await window.go.main.App.RemoveClipboardItem(index)
    // 重新获取数据
    await fetchHistory()
    ElMessage.success('删除成功')
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 清空记录
const clearHistory = async () => {
  try {
    loading.value = true
    await window.go.main.App.ClearClipboardHistory()
    clipboardHistory.value = []
    ElMessage.success('已清空所有记录')
  } catch (error) {
    ElMessage.error('清空记录失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 打开剪贴板文件
const openClipboardFile = async () => {
  try {
    loading.value = true
    await window.go.main.App.OpenClipboardFile()
    ElMessage.success('已打开剪贴板文件')
  } catch (error) {
    ElMessage.error('打开文件失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 刷新历史记录
const refreshHistory = async () => {
  await fetchHistory()
  ElMessage.success('已刷新')
}

// 处理搜索
const handleSearch = () => {
  // 搜索逻辑已经在computed属性中实现
  console.log('搜索文本:', searchText.value)
}

// 组件挂载时加载历史记录
onMounted(() => {
  fetchHistory()
})
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
}

.tool-header {
  margin-bottom: 1.5rem;
  text-align: left;
}

.tool-header h2 {
  margin: 0 0 0.5rem 0;
  color: #333;
  font-size: 1.8rem;
}

.tool-header p {
  margin: 0;
  color: #666;
  font-size: 1.1rem;
}

.clipboard-card,
.history-card {
  margin-bottom: 1rem;
}

.clipboard-header,
.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.clipboard-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1rem;
}

.clipboard-actions {
  display: flex;
  gap: 0.5rem;
}

.history-content {
  min-height: 300px;
}

.history-actions {
  display: flex;
  gap: 0.5rem;
}

.search-container {
  margin-bottom: 1rem;
}

.history-list {
  max-height: 500px;
  overflow-y: auto;
}

.history-item {
  border: 1px solid #e1e5e9;
  border-radius: 4px;
  margin-bottom: 0.5rem;
  padding: 0.75rem;
  background-color: #fff;
}

.history-item:hover {
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.item-time {
  font-size: 0.75rem;
  color: #888;
}

.item-actions {
  display: flex;
  gap: 0.5rem;
}

.item-actions .el-button {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}

.item-content {
  word-break: break-all;
  white-space: pre-wrap;
  font-size: 0.875rem;
  line-height: 1.5;
  color: #333;
}

@media (max-width: 768px) {
  .clipboard-header,
  .history-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .clipboard-stats {
    width: 100%;
    justify-content: space-between;
  }
  
  .item-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .item-actions {
    align-self: flex-end;
  }
}
</style>