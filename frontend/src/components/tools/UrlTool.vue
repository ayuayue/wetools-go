<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>URL 编码/解码工具</h2>
      <p>URL编码解码工具，支持URL参数解析</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="8"
              placeholder="请输入要编码或解码的URL或文本"
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="button-group">
        <el-button type="primary" @click="encodeUrl">
          <i class="fas fa-lock"></i> URL编码
        </el-button>
        <el-button @click="decodeUrl">
          <i class="fas fa-unlock"></i> URL解码
        </el-button>
        <el-button @click="parseUrlParams">
          <i class="fas fa-search"></i> 解析URL参数
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
              :rows="8"
              placeholder="转换结果将显示在这里"
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="params-section" v-if="urlParams.length > 0">
        <h3>URL参数解析结果</h3>
        <el-table :data="urlParams" style="width: 100%" border>
          <el-table-column prop="key" label="参数名" width="180"></el-table-column>
          <el-table-column prop="value" label="参数值"></el-table-column>
        </el-table>
      </div>
      
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
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert, ElTable, ElTableColumn } from 'element-plus'

// 数据模型
const inputData = ref('')
const outputData = ref('')
const urlParams = ref([])
const validationResult = ref(null)

// URL编码
const encodeUrl = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要编码的文本'
      }
      return
    }
    
    outputData.value = encodeURIComponent(inputData.value)
    urlParams.value = []
    validationResult.value = {
      type: 'success',
      message: 'URL编码成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `编码失败: ${error.message}`
    }
  }
}

// URL解码
const decodeUrl = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要解码的URL编码字符串'
      }
      return
    }
    
    outputData.value = decodeURIComponent(inputData.value)
    urlParams.value = []
    validationResult.value = {
      type: 'success',
      message: 'URL解码成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `解码失败: ${error.message}`
    }
  }
}

// 解析URL参数
const parseUrlParams = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入包含参数的URL'
      }
      return
    }
    
    // 创建一个虚拟的a标签来解析URL
    const a = document.createElement('a')
    a.href = inputData.value
    
    // 获取查询参数
    const params = new URLSearchParams(a.search)
    const paramArray = []
    
    for (const [key, value] of params) {
      paramArray.push({ key, value })
    }
    
    urlParams.value = paramArray
    outputData.value = JSON.stringify(Object.fromEntries(params), null, 2)
    
    validationResult.value = {
      type: 'success',
      message: `成功解析到 ${paramArray.length} 个参数！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `URL参数解析失败: ${error.message}`
    }
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
  urlParams.value = []
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
inputData.value = 'https://example.com/search?q=hello world&category=tools&page=1'
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
  min-width: 120px;
  max-width: 160px;
}

.params-section {
  width: 100%;
}

.params-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
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