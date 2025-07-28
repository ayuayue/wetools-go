<template>
  <div v-if="!isDialogMode">
    <el-form
      ref="proxyFormRef"
      :model="proxyForm"
      :rules="proxyRules"
      label-width="80px"
    >
      <el-form-item label="启用代理" prop="enabled">
        <el-switch v-model="proxyForm.enabled" />
      </el-form-item>
      
      <el-form-item label="代理地址" prop="host">
        <el-input
          v-model="proxyForm.host"
          placeholder="例如: 127.0.0.1"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="端口" prop="port">
        <el-input
          v-model="proxyForm.port"
          placeholder="例如: 8080"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="proxyForm.username"
          placeholder="可选"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="密码" prop="password">
        <el-input
          v-model="proxyForm.password"
          type="password"
          placeholder="可选"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" @click="saveProxyConfig">保存配置</el-button>
        <el-button 
          @click="testProxyConnection" 
          :loading="testing"
          :disabled="!proxyForm.enabled || !proxyForm.host || !proxyForm.port"
        >
          {{ testing ? '测试中...' : '测试连接' }}
        </el-button>
      </el-form-item>
    </el-form>
    
    <div v-if="testResult" class="test-result" :class="testResult.type">
      <div class="test-result-header">
        <i :class="getResultIcon(testResult.type)"></i>
        <span class="test-result-title">{{ getResultTitle(testResult.type) }}</span>
      </div>
      <div class="test-result-message">{{ testResult.message }}</div>
      <div class="test-result-time" v-if="testResult.time">
        测试时间: {{ testResult.time }}
      </div>
    </div>
    
    <div class="proxy-info" v-if="proxyForm.enabled && proxyForm.host && proxyForm.port">
      <h4>代理信息</h4>
      <div class="proxy-info-item">
        <span class="label">代理地址:</span>
        <span class="value">{{ proxyForm.host }}:{{ proxyForm.port }}</span>
      </div>
      <div class="proxy-info-item" v-if="proxyForm.username">
        <span class="label">用户名:</span>
        <span class="value">{{ proxyForm.username }}</span>
      </div>
      <div class="proxy-info-item">
        <span class="label">状态:</span>
        <span class="value" :class="proxyStatusClass">{{ proxyStatusText }}</span>
      </div>
    </div>
  </div>
  
  <el-dialog
    v-else
    v-model="dialogVisible"
    title="代理设置"
    width="500px"
    @close="handleClose"
  >
    <el-form
      ref="proxyFormRef"
      :model="proxyForm"
      :rules="proxyRules"
      label-width="80px"
    >
      <el-form-item label="启用代理" prop="enabled">
        <el-switch v-model="proxyForm.enabled" />
      </el-form-item>
      
      <el-form-item label="代理地址" prop="host">
        <el-input
          v-model="proxyForm.host"
          placeholder="例如: 127.0.0.1"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="端口" prop="port">
        <el-input
          v-model="proxyForm.port"
          placeholder="例如: 8080"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="proxyForm.username"
          placeholder="可选"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
      
      <el-form-item label="密码" prop="password">
        <el-input
          v-model="proxyForm.password"
          type="password"
          placeholder="可选"
          :disabled="!proxyForm.enabled"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="saveProxyConfig">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElDialog, ElForm, ElFormItem, ElInput, ElSwitch, ElButton, ElMessage, ElAlert } from 'element-plus'

// 定义 emits
const emit = defineEmits(['update:visible', 'proxy-config-updated'])

// 定义 props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  isDialogMode: {
    type: Boolean,
    default: false
  }
})

// 表单引用
const proxyFormRef = ref()

// 对话框可见性
const dialogVisible = ref(false)

// 表单数据
const proxyForm = ref({
  enabled: false,
  host: '',
  port: '',
  username: '',
  password: ''
})

// 测试结果
const testResult = ref(null)

// 测试状态
const testing = ref(false)

// 代理状态
const proxyStatus = ref('unknown')

// 表单验证规则
const proxyRules = {
  host: [
    { required: true, message: '请输入代理地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口', trigger: 'blur' }
  ]
}

// 计算属性
const proxyStatusText = computed(() => {
  switch (proxyStatus.value) {
    case 'connected':
      return '已连接'
    case 'disconnected':
      return '未连接'
    case 'testing':
      return '测试中'
    default:
      return '未知'
  }
})

const proxyStatusClass = computed(() => {
  return `status-${proxyStatus.value}`
})

// 方法
const getResultIcon = (type) => {
  const icons = {
    success: 'fas fa-check-circle',
    error: 'fas fa-times-circle',
    warning: 'fas fa-exclamation-triangle',
    info: 'fas fa-info-circle'
  }
  return icons[type] || icons.info
}

const getResultTitle = (type) => {
  const titles = {
    success: '成功',
    error: '错误',
    warning: '警告',
    info: '信息'
  }
  return titles[type] || titles.info
}

// 监听可见性变化
onMounted(() => {
  dialogVisible.value = props.visible
  loadProxyConfig()
})

// 加载代理配置
const loadProxyConfig = async () => {
  try {
    // 从 localStorage 加载配置
    const savedConfig = localStorage.getItem('proxyConfig')
    if (savedConfig) {
      proxyForm.value = JSON.parse(savedConfig)
    }
  } catch (error) {
    console.error('加载代理配置失败:', error)
  }
}

// 保存代理配置
const saveProxyConfig = async () => {
  if (!proxyFormRef.value) return
  
  await proxyFormRef.value.validate((valid) => {
    if (valid) {
      try {
        // 保存到 localStorage
        localStorage.setItem('proxyConfig', JSON.stringify(proxyForm.value))
        
        // 通知父组件配置已更新
        emit('proxy-config-updated', proxyForm.value)
        
        // 如果是对话框模式，关闭对话框
        if (props.isDialogMode) {
          handleClose()
        }
        
        // 显示成功消息
        ElMessage.success('代理配置保存成功！')
      } catch (error) {
        ElMessage.error('保存代理配置失败: ' + error.message)
      }
    }
  })
}

// 测试代理连接
const testProxyConnection = async () => {
  if (!proxyForm.value.enabled) {
    testResult.value = {
      type: 'warning',
      message: '请先启用代理'
    }
    return
  }
  
  if (!proxyForm.value.host || !proxyForm.value.port) {
    testResult.value = {
      type: 'warning',
      message: '请填写代理地址和端口'
    }
    return
  }
  
  // 设置测试状态
  testing.value = true
  proxyStatus.value = 'testing'
  
  try {
    testResult.value = {
      type: 'info',
      message: '正在测试代理连接...',
      time: new Date().toLocaleString()
    }
    
    // 调用Wails后端方法进行代理测试
    const [success, message, error] = await window.go.main.App.TestProxyConnection()
    
    if (error) {
      testResult.value = {
        type: 'error',
        message: '代理连接测试失败: ' + error,
        time: new Date().toLocaleString()
      }
      proxyStatus.value = 'disconnected'
    } else if (success) {
      testResult.value = {
        type: 'success',
        message: message,
        time: new Date().toLocaleString()
      }
      proxyStatus.value = 'connected'
    } else {
      testResult.value = {
        type: 'error',
        message: message,
        time: new Date().toLocaleString()
      }
      proxyStatus.value = 'disconnected'
    }
  } catch (error) {
    console.error('代理测试错误:', error)
    testResult.value = {
      type: 'error',
      message: '代理连接测试失败: ' + (error.message || '未知错误'),
      time: new Date().toLocaleString()
    }
    proxyStatus.value = 'disconnected'
  } finally {
    testing.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  emit('update:visible', false)
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.test-result {
  margin-top: 1rem;
  padding: 1rem;
  border-radius: 4px;
  border-left: 4px solid #909399;
  background-color: #f4f4f5;
}

.test-result.success {
  border-left-color: #67c23a;
  background-color: #f0f9ff;
}

.test-result.error {
  border-left-color: #f56c6c;
  background-color: #fef0f0;
}

.test-result.warning {
  border-left-color: #e6a23c;
  background-color: #fdf6ec;
}

.test-result.info {
  border-left-color: #409eff;
  background-color: #f0f9ff;
}

.test-result-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  font-weight: bold;
}

.test-result-header i {
  font-size: 1.2rem;
}

.test-result.success .test-result-header i {
  color: #67c23a;
}

.test-result.error .test-result-header i {
  color: #f56c6c;
}

.test-result.warning .test-result-header i {
  color: #e6a23c;
}

.test-result.info .test-result-header i {
  color: #409eff;
}

.test-result-message {
  margin-bottom: 0.5rem;
}

.test-result-time {
  font-size: 0.8rem;
  color: #909399;
}

.proxy-info {
  margin-top: 1rem;
  padding: 1rem;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #f5f7fa;
}

.proxy-info h4 {
  margin: 0 0 1rem 0;
  color: #606266;
}

.proxy-info-item {
  display: flex;
  margin-bottom: 0.5rem;
}

.proxy-info-item:last-child {
  margin-bottom: 0;
}

.label {
  width: 80px;
  font-weight: bold;
  color: #606266;
}

.value {
  flex: 1;
  color: #606266;
}

.value.status-connected {
  color: #67c23a;
}

.value.status-disconnected {
  color: #f56c6c;
}

.value.status-testing {
  color: #409eff;
}

.value.status-unknown {
  color: #909399;
}
</style>