<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>XML 格式化工具</h2>
      <p>在线XML格式化、校验、压缩工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="10"
              placeholder='请输入XML数据，例如：&lt;root&gt;&lt;name&gt;WeTools&lt;/name&gt;&lt;type&gt;developer tools&lt;/type&gt;&lt;/root&gt;'
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="button-group">
        <el-button type="primary" @click="formatXml">
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button @click="compressXml">
          <i class="fas fa-compress"></i> 压缩
        </el-button>
        <el-button @click="validateXml">
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
              placeholder="格式化后的XML结果将显示在这里"
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

// 格式化XML
const formatXml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 简单的XML格式化函数
    const formatted = formatXmlString(inputData.value)
    outputData.value = formatted
    validationResult.value = {
      type: 'success',
      message: 'XML格式化成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML格式错误: ${error.message}`
    }
  }
}

// 压缩XML
const compressXml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除多余的空白字符
    const compressed = inputData.value.replace(/>\s+</g, '><').trim()
    outputData.value = compressed
    validationResult.value = {
      type: 'success',
      message: 'XML压缩成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML处理错误: ${error.message}`
    }
  }
}

// 校验XML
const validateXml = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'info',
        message: '请输入XML数据进行校验'
      }
      return
    }
    
    // 简单的XML校验
    const parser = new DOMParser()
    const doc = parser.parseFromString(inputData.value, 'text/xml')
    
    // 检查是否有解析错误
    const errorNode = doc.querySelector('parsererror')
    if (errorNode) {
      throw new Error('XML格式不正确')
    }
    
    validationResult.value = {
      type: 'success',
      message: 'XML格式正确！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML格式错误: ${error.message}`
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

// 简单的XML格式化函数
const formatXmlString = (xml) => {
  let formatted = ''
  let indent = ''
  const tab = '  ' // 2个空格缩进
  let inTag = false
  
  for (let i = 0; i < xml.length; i++) {
    const char = xml.charAt(i)
    
    if (char === '<' && xml.charAt(i + 1) !== '/') {
      // 开始标签
      formatted += '\n' + indent + char
      indent += tab
      inTag = true
    } else if (char === '<' && xml.charAt(i + 1) === '/') {
      // 结束标签
      indent = indent.slice(0, -tab.length)
      formatted += '\n' + indent + char
      inTag = true
    } else if (char === '>' && inTag) {
      // 标签结束
      formatted += char
      inTag = false
    } else {
      formatted += char
    }
  }
  
  // 移除第一行的换行符并清理多余的空白
  return formatted.trim()
}

// 初始化示例数据
inputData.value = '<root><name>WeTools</name><type>developer tools</type><features><feature>JSON</feature><feature>XML</feature><feature>HTML</feature></features></root>'
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