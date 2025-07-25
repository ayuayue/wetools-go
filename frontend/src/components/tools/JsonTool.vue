<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入JSON</h3>
      <textarea 
        id="input" 
        v-model="inputData" 
        placeholder='{"name": "WeTools", "type": "developer tools"}'
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="formatJson">
        <i class="fas fa-magic"></i> 格式化
      </button>
      <button @click="compressJson">
        <i class="fas fa-compress"></i> 压缩
      </button>
      <button @click="validateJson">
        <i class="fas fa-check-circle"></i> 校验
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
      <div 
        id="output" 
        class="result" 
        ref="outputRef"
        :class="{ 'validation-result': isValidationResult }"
      >{{ outputData }}</div>
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
      <p>JSON（JavaScript Object Notation）是一种轻量级的数据交换格式。本工具可以帮助您：</p>
      <ul class="description-list">
        <li>格式化JSON数据，使其更易读</li>
        <li>压缩JSON数据，去除多余空格</li>
        <li>校验JSON语法是否正确</li>
        <li>支持JSON与其他格式互转</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const inputData = ref('')
const outputData = ref('')
const outputRef = ref(null)

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// 格式化JSON
const formatJson = () => {
  try {
    const jsonObj = JSON.parse(inputData.value)
    outputData.value = JSON.stringify(jsonObj, null, 2)
  } catch (e) {
    outputData.value = 'JSON格式错误: ' + e.message
  }
}

// 压缩JSON
const compressJson = () => {
  try {
    const jsonObj = JSON.parse(inputData.value)
    outputData.value = JSON.stringify(jsonObj)
  } catch (e) {
    outputData.value = 'JSON格式错误: ' + e.message
  }
}

// 校验JSON
const validateJson = () => {
  try {
    JSON.parse(inputData.value)
    outputData.value = '✓ JSON格式正确'
  } catch (e) {
    outputData.value = '✗ JSON格式错误: ' + e.message
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
    // 可以添加一个提示，但这里简化处理
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
onMounted(() => {
  inputData.value = '{"name": "WeTools", "type": "developer tools", "features": ["JSON", "XML", "HTML", "Base64", "Hash", "Encryption"]}'
  formatJson()
})
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

.result.validation-result {
  font-family: inherit;
  white-space: normal;
  text-align: center;
  font-size: 16px;
  font-weight: bold;
  padding: 2rem;
}

.result.validation-result:nth-child(1) {
  color: #10b981;
}

.result.validation-result:nth-child(2) {
  color: #ef4444;
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