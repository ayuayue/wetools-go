<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>系统设置</h2>
      <p>配置应用系统设置</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="settings-tabs">
        <el-tab-pane label="常规设置" name="general">
          <div class="settings-section">
            <el-form
              ref="generalFormRef"
              :model="generalSettings"
              label-width="120px"
            >
              <el-form-item label="开机自启动">
                <el-switch
                  v-model="generalSettings.autoStart"
                  @change="toggleAutoStart"
                />
              </el-form-item>
              
              <el-form-item label="最小化到托盘">
                <el-switch
                  v-model="generalSettings.minimizeToTray"
                  disabled
                />
              </el-form-item>
              
              <el-form-item label="启动时最小化">
                <el-switch
                  v-model="generalSettings.startMinimized"
                  disabled
                />
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="外观设置" name="appearance">
          <div class="settings-section">
            <el-form
              ref="appearanceFormRef"
              :model="appearanceSettings"
              label-width="120px"
            >
              <el-form-item label="主题">
                <el-select
                  v-model="appearanceSettings.theme"
                  placeholder="选择主题"
                  @change="changeTheme"
                >
                  <el-option label="浅色主题" value="light" />
                  <el-option label="深色主题" value="dark" />
                  <el-option label="跟随系统" value="system" />
                </el-select>
              </el-form-item>
              
              <el-form-item label="语言">
                <el-select
                  v-model="appearanceSettings.language"
                  placeholder="选择语言"
                  disabled
                >
                  <el-option label="简体中文" value="zh-CN" />
                  <el-option label="English" value="en" />
                </el-select>
              </el-form-item>
            </el-form>
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
import { ref, onMounted } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElForm, ElFormItem, ElSwitch, ElSelect, ElOption, ElAlert } from 'element-plus'

// 数据模型
const activeTab = ref('general')
const generalFormRef = ref()
const appearanceFormRef = ref()

const generalSettings = ref({
  autoStart: false,
  minimizeToTray: true,
  startMinimized: false
})

const appearanceSettings = ref({
  theme: 'light',
  language: 'zh-CN'
})

const validationResult = ref(null)

// 切换开机自启动
const toggleAutoStart = async (value) => {
  try {
    const result = await window.go.main.App.ToggleAutoStart()
    generalSettings.value.autoStart = result
    
    validationResult.value = {
      type: 'success',
      message: result ? '已设置开机自启动' : '已取消开机自启动'
    }
  } catch (error) {
    // 回滚开关状态
    generalSettings.value.autoStart = !value
    
    validationResult.value = {
      type: 'error',
      message: `设置失败: ${error.message}`
    }
  }
}

// 更改主题
const changeTheme = (value) => {
  try {
    // 保存主题设置到localStorage
    localStorage.setItem('theme', value)
    
    // 应用主题
    if (value === 'dark') {
      document.body.classList.add('dark-theme')
    } else {
      document.body.classList.remove('dark-theme')
    }
    
    validationResult.value = {
      type: 'success',
      message: '主题设置已保存'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `主题设置失败: ${error.message}`
    }
  }
}

// 加载设置
const loadSettings = async () => {
  try {
    // 加载开机自启动状态
    const autoStartStatus = await window.go.main.App.GetAutoStartStatus()
    generalSettings.value.autoStart = autoStartStatus
    
    // 加载主题设置
    const savedTheme = localStorage.getItem('theme')
    if (savedTheme) {
      appearanceSettings.value.theme = savedTheme
    }
  } catch (error) {
    console.error('加载设置失败:', error)
  }
}

// 组件挂载时加载设置
onMounted(() => {
  loadSettings()
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

.settings-tabs {
  width: 100%;
}

.settings-section {
  padding: 1rem 0;
}

.settings-section .el-form {
  max-width: 500px;
}

.validation-result {
  margin-top: 1rem;
}
</style>