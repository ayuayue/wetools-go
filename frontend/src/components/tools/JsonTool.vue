<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>JSON 格式化工具</h2>
      <p>在线JSON格式化、校验、压缩工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="10"
              placeholder='请输入JSON数据，例如：{"name": "WeTools", "type": "developer tools"}'
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="button-group">
        <el-button type="primary" @click="formatJson">
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button @click="compressJson">
          <i class="fas fa-compress"></i> 压缩
        </el-button>
        <el-button @click="validateJson">
          <i class="fas fa-check-circle"></i> 校验
        </el-button>
        <el-button @click="clearData">
          <i class="fas fa-trash"></i> 清空
        </el-button>
        <el-button type="success" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </el-button>
      </div>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="output-section">
            <el-input
              v-model="outputData"
              type="textarea"
              :rows="10"
              placeholder="格式化后的JSON结果将显示在这里"
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
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
import { ref } from 'vue'
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert } from 'element-plus'

// 数据模型
const inputData = ref('')
const outputData = ref('')
const validationResult = ref(null)

// 格式化JSON
const formatJson = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    const parsed = JSON.parse(inputData.value)
    outputData.value = JSON.stringify(parsed, null, 2)
    validationResult.value = {
      type: 'success',
      message: 'JSON格式化成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `JSON格式错误: ${error.message}`
    }
  }
}

// 压缩JSON
const compressJson = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    const parsed = JSON.parse(inputData.value)
    outputData.value = JSON.stringify(parsed)
    validationResult.value = {
      type: 'success',
      message: 'JSON压缩成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `JSON格式错误: ${error.message}`
    }
  }
}

// 校验JSON
const validateJson = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'info',
        message: '请输入JSON数据进行校验'
      }
      return
    }
    
    JSON.parse(inputData.value)
    validationResult.value = {
      type: 'success',
      message: 'JSON格式正确！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `JSON格式错误: ${error.message}`
    }
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  if (!outputData.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(outputData.value)
    validationResult.value = {
      type: 'success',
      message: '结果已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 初始化示例数据
inputData.value = '{"name": "WeTools", "type": "developer tools", "features": ["JSON", "XML", "HTML", "Base64", "Hash", "Encryption"]}'
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

.input-section,
.output-section {
  width: 100%;
}

.button-group {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.button-group .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .button-group {
    flex-direction: column;
    align-items: center;
  }
  
  .button-group .el-button {
    width: 100%;
    max-width: none;
  }
}
</style>