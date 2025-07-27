<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>HTML 格式化工具</h2>
      <p>在线HTML格式化、清理、压缩工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="10"
              placeholder='请输入HTML代码，例如：&lt;div&gt;&lt;h1&gt;Hello World&lt;/h1&gt;&lt;p&gt;This is a paragraph.&lt;/p&gt;&lt;/div&gt;'
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="button-group">
        <el-button type="primary" @click="formatHtml">
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button @click="compressHtml">
          <i class="fas fa-compress"></i> 压缩
        </el-button>
        <el-button @click="cleanHtml">
          <i class="fas fa-broom"></i> 清理
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
              placeholder="格式化后的HTML结果将显示在这里"
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

// 格式化HTML
const formatHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 简单的HTML格式化函数
    const formatted = formatHtmlString(inputData.value)
    outputData.value = formatted
    validationResult.value = {
      type: 'success',
      message: 'HTML格式化成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML格式错误: ${error.message}`
    }
  }
}

// 压缩HTML
const compressHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除多余的空白字符
    const compressed = inputData.value
      .replace(/>\s+</g, '><')
      .replace(/\s+/g, ' ')
      .trim()
    outputData.value = compressed
    validationResult.value = {
      type: 'success',
      message: 'HTML压缩成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML处理错误: ${error.message}`
    }
  }
}

// 清理HTML
const cleanHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除注释和多余的空白
    const cleaned = inputData.value
      .replace(/<!--[\s\S]*?-->/g, '')
      .replace(/>\s+</g, '><')
      .replace(/\s+/g, ' ')
      .trim()
    outputData.value = cleaned
    validationResult.value = {
      type: 'success',
      message: 'HTML清理成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML处理错误: ${error.message}`
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

// 简单的HTML格式化函数
const formatHtmlString = (html) => {
  let formatted = ''
  let indent = ''
  const tab = '  ' // 2个空格缩进
  let inTag = false
  let inScript = false
  let inStyle = false
  
  for (let i = 0; i < html.length; i++) {
    const char = html.charAt(i)
    
    // 检查是否进入script或style标签
    if (html.substring(i, i + 7).toLowerCase() === '<script') {
      inScript = true
    } else if (html.substring(i, i + 8).toLowerCase() === '</script') {
      inScript = false
    } else if (html.substring(i, i + 6).toLowerCase() === '<style') {
      inStyle = true
    } else if (html.substring(i, i + 7).toLowerCase() === '</style') {
      inStyle = false
    }
    
    if (char === '<' && !inScript && !inStyle) {
      // 开始标签
      if (html.charAt(i + 1) === '/') {
        // 结束标签
        indent = indent.slice(0, -tab.length)
        formatted += '\n' + indent + char
      } else {
        // 开始标签
        formatted += '\n' + indent + char
        // 不对自闭合标签增加缩进
        if (html.substring(i).indexOf('>') > html.substring(i).indexOf('/')) {
          // 不是自闭合标签
          indent += tab
        }
      }
      inTag = true
    } else if (char === '>' && inTag && !inScript && !inStyle) {
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
inputData.value = '<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello World</h1><p>This is a paragraph.</p></body></html>'
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