<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入HTML</h3>
      <textarea 
        v-model="inputData" 
        placeholder='<div><h1>WeTools</h1><p>开发者工具箱</p></div>'
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="formatHtml">
        <i class="fas fa-magic"></i> 格式化
      </button>
      <button @click="compressHtml">
        <i class="fas fa-compress"></i> 压缩
      </button>
      <button @click="escapeHtml">
        <i class="fas fa-code"></i> HTML转义
      </button>
      <button @click="unescapeHtml">
        <i class="fas fa-code"></i> HTML反转义
      </button>
      <button class="copy-btn" @click="copyInput">
        <i class="fas fa-copy"></i> 复制输入
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <div class="result-header">
        <h3><i class="fas fa-file-alt"></i> 输出结果</h3>
      </div>
      <div class="result">{{ outputData }}</div>
      <div class="result-footer" v-if="outputData && !isValidationResult">
        <button class="copy-btn" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </button>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>HTML（超文本标记语言）是用于创建网页的标准标记语言。本工具可以帮助您：</p>
      <ul class="description-list">
        <li>格式化HTML代码，使其更易读</li>
        <li>压缩HTML代码，去除多余空格和换行</li>
        <li>HTML实体编码和解码</li>
        <li>清理HTML代码</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const inputData = ref('')
const outputData = ref('')

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// 格式化HTML
const formatHtml = () => {
  try {
    // 简单的HTML格式化
    let formatted = inputData.value
      .replace(/></g, ">\n<")
      .replace(/^\s+|\s+$/g, "")
    
    outputData.value = formatted
  } catch (e) {
    outputData.value = '处理错误: ' + e.message
  }
}

// 压缩HTML
const compressHtml = () => {
  try {
    // 简单的HTML压缩
    let compressed = inputData.value
      .replace(/\n/g, "")
      .replace(/\s+/g, " ")
      .replace(/>\s+</g, "><")
    
    outputData.value = compressed
  } catch (e) {
    outputData.value = '处理错误: ' + e.message
  }
}

// HTML转义
const escapeHtml = () => {
  try {
    const map = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#39;'
    }
    
    outputData.value = inputData.value.replace(/[&<>"']/g, (m) => map[m])
  } catch (e) {
    outputData.value = '处理错误: ' + e.message
  }
}

// HTML反转义
const unescapeHtml = () => {
  try {
    const map = {
      '&amp;': '&',
      '&lt;': '<',
      '&gt;': '>',
      '&quot;': '"',
      '&#39;': "'"
    }
    
    let unescaped = inputData.value
    for (let key in map) {
      unescaped = unescaped.replace(new RegExp(key, 'g'), map[key])
    }
    
    outputData.value = unescaped
  } catch (e) {
    outputData.value = '处理错误: ' + e.message
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
}

// 复制输入内容
const copyInput = async () => {
  try {
    await navigator.clipboard.writeText(inputData.value)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = inputData.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

// 复制结果
const copyResult = async () => {
  try {
    await navigator.clipboard.writeText(outputData.value)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = outputData.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

// 初始化数据
const init = () => {
  inputData.value = '<div><h1>WeTools</h1><p>开发者工具箱</p><ul><li>JSON工具</li><li>XML工具</li><li>HTML工具</li></ul></div>'
  formatHtml()
}

init()
</script>

<style scoped>
.tool-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.tool-section {
  margin-bottom: 1.5rem;
}

.tool-section:last-child {
  margin-bottom: 0;
}

.tool-section h3 {
  margin-bottom: 1rem;
  color: #444;
  display: flex;
  align-items: center;
  gap: 8px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.result-header h3 {
  margin-bottom: 0;
}

textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.2s;
  min-height: 150px;
}

textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

.button-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin: 1rem 0;
}

button {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 5px;
}

button:hover {
  background: #5a6fd8;
}

button.secondary {
  background: #f1f5f9;
  color: #333;
}

button.secondary:hover {
  background: #e2e8f0;
}

.copy-btn {
  background: #4ade80;
  font-size: 0.9rem;
  padding: 0.25rem 0.75rem;
}

.copy-btn:hover {
  background: #22c55e;
}

.result {
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 4px;
  padding: 1rem;
  min-height: 150px;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow: auto;
  max-height: 400px;
}

.result-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

.description-list {
  margin-top: 10px;
  padding-left: 20px;
}

@media (max-width: 768px) {
  .button-group {
    flex-direction: column;
  }
  
  button {
    width: 100%;
    justify-content: center;
  }
  
  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>