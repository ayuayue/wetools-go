<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>网站集成工具</h2>
      <p>在应用内嵌入和管理常用网站</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="tool-tabs">
        <el-tab-pane label="网站列表" name="list">
          <div class="website-list-section">
            <div class="add-website-form">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="网站名称">
                    <el-input
                      v-model="newWebsite.name"
                      placeholder="请输入网站名称"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="网站地址">
                    <el-input
                      v-model="newWebsite.url"
                      placeholder="请输入网站地址，如 https://example.com"
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="4">
                  <div class="form-buttons">
                    <el-button type="primary" @click="addWebsite">
                      <i class="fas fa-plus"></i> 添加
                    </el-button>
                  </div>
                </el-col>
              </el-row>
            </div>
            
            <div class="website-list">
              <el-table
                :data="websites"
                style="width: 100%"
                row-key="id"
              >
                <el-table-column prop="name" label="网站名称" width="200" />
                <el-table-column prop="url" label="网站地址" />
                <el-table-column label="操作" width="200">
                  <template #default="scope">
                    <el-button
                      size="small"
                      @click="openWebsite(scope.row)"
                    >
                      <i class="fas fa-external-link-alt"></i> 打开
                    </el-button>
                    <el-button
                      size="small"
                      type="danger"
                      @click="removeWebsite(scope.row)"
                    >
                      <i class="fas fa-trash"></i> 删除
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="网站浏览" name="browse">
          <div class="website-browse-section">
            <div class="website-selector">
              <el-select
                v-model="selectedWebsite"
                placeholder="请选择要浏览的网站"
                @change="loadWebsite"
              >
                <el-option
                  v-for="website in websites"
                  :key="website.id"
                  :label="website.name"
                  :value="website.url"
                />
              </el-select>
              <el-button
                type="primary"
                @click="refreshWebsite"
                :disabled="!selectedWebsite"
              >
                <i class="fas fa-sync"></i> 刷新
              </el-button>
            </div>
            
            <div class="website-frame-container" v-if="selectedWebsite">
              <iframe
                :src="selectedWebsite"
                class="website-frame"
                frameborder="0"
                allowfullscreen
              ></iframe>
            </div>
            
            <div class="website-placeholder" v-else>
              <i class="fas fa-globe"></i>
              <p>请选择一个网站开始浏览</p>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      
      <div class="validation-result" v-if="validationResult">
        <el-alert
          :type="validationResult.type"
          :title="validationResult.message"
          show-icon
          :closable="false"
        />
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElFormItem, ElInput, ElButton, ElAlert, ElTable, ElTableColumn, ElSelect, ElOption } from 'element-plus'

// 数据模型
const activeTab = ref('list')
const newWebsite = ref({
  name: '',
  url: ''
})
const websites = ref([])
const selectedWebsite = ref('')
const validationResult = ref(null)

// 添加网站
const addWebsite = () => {
  try {
    // 验证输入
    if (!newWebsite.value.name.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入网站名称'
      }
      return
    }
    
    if (!newWebsite.value.url.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入网站地址'
      }
      return
    }
    
    // 验证URL格式
    try {
      new URL(newWebsite.value.url)
    } catch (e) {
      validationResult.value = {
        type: 'warning',
        message: '请输入有效的网站地址，如 https://example.com'
      }
      return
    }
    
    // 检查是否已存在
    const exists = websites.value.some(
      site => site.url === newWebsite.value.url
    )
    if (exists) {
      validationResult.value = {
        type: 'warning',
        message: '该网站已存在'
      }
      return
    }
    
    // 添加网站
    const website = {
      id: Date.now().toString(),
      name: newWebsite.value.name.trim(),
      url: newWebsite.value.url.trim()
    }
    
    websites.value.push(website)
    
    // 保存到本地存储
    saveWebsites()
    
    // 清空表单
    newWebsite.value = {
      name: '',
      url: ''
    }
    
    validationResult.value = {
      type: 'success',
      message: '网站添加成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `添加失败: ${error.message}`
    }
  }
}

// 删除网站
const removeWebsite = (website) => {
  try {
    websites.value = websites.value.filter(site => site.id !== website.id)
    
    // 如果删除的是当前选中的网站，清空选中
    if (selectedWebsite.value === website.url) {
      selectedWebsite.value = ''
    }
    
    // 保存到本地存储
    saveWebsites()
    
    validationResult.value = {
      type: 'success',
      message: '网站删除成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `删除失败: ${error.message}`
    }
  }
}

// 打开网站（在新标签页中打开）
const openWebsite = (website) => {
  try {
    window.open(website.url, '_blank')
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `打开网站失败: ${error.message}`
    }
  }
}

// 加载网站（在iframe中显示）
const loadWebsite = (url) => {
  selectedWebsite.value = url
}

// 刷新网站
const refreshWebsite = () => {
  if (selectedWebsite.value) {
    // 通过重新设置src来刷新iframe
    const currentUrl = selectedWebsite.value
    selectedWebsite.value = ''
    setTimeout(() => {
      selectedWebsite.value = currentUrl
    }, 100)
  }
}

// 保存网站列表到本地存储
const saveWebsites = () => {
  try {
    localStorage.setItem('integratedWebsites', JSON.stringify(websites.value))
  } catch (error) {
    console.error('保存网站列表失败:', error)
  }
}

// 从本地存储加载网站列表
const loadWebsites = () => {
  try {
    const savedWebsites = localStorage.getItem('integratedWebsites')
    if (savedWebsites) {
      websites.value = JSON.parse(savedWebsites)
    }
  } catch (error) {
    console.error('加载网站列表失败:', error)
  }
}

// 初始化示例数据
const initExampleData = () => {
  if (websites.value.length === 0) {
    // 添加示例网站
    websites.value = [
      {
        id: '1',
        name: '通义千问',
        url: 'https://chat.qwen.ai/'
      }
    ]
    saveWebsites()
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadWebsites()
  initExampleData()
})

// 组件卸载前保存数据
onUnmounted(() => {
  saveWebsites()
})
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

.tool-tabs {
  width: 100%;
}

.website-list-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.add-website-form {
  width: 100%;
  padding: 1rem;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.form-buttons {
  display: flex;
  align-items: flex-end;
  height: 100%;
}

.website-list {
  width: 100%;
}

.website-browse-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.website-selector {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.website-selector .el-select {
  flex: 1;
}

.website-frame-container {
  width: 100%;
  height: 600px;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
}

.website-frame {
  width: 100%;
  height: 100%;
}

.website-placeholder {
  width: 100%;
  height: 600px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  border: 1px dashed #ddd;
  border-radius: 4px;
  color: #999;
}

.website-placeholder i {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .website-selector {
    flex-direction: column;
    align-items: stretch;
  }
  
  .website-frame-container {
    height: 400px;
  }
  
  .website-placeholder {
    height: 400px;
  }
}
</style>