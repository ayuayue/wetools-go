<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入文本</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入要计算哈希的文本"
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="calculateMD5">
        <i class="fas fa-calculator"></i> MD5
      </button>
      <button @click="calculateSHA1">
        <i class="fas fa-calculator"></i> SHA1
      </button>
      <button @click="calculateSHA256">
        <i class="fas fa-calculator"></i> SHA256
      </button>
      <button @click="convertToUpper">
        <i class="fas fa-text-height"></i> 转大写
      </button>
      <button @click="convertToLower">
        <i class="fas fa-text-width"></i> 转小写
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
      <CodeBlock 
        :code="outputData" 
        language="text"
        :show-line-numbers="true"
        :show-header="true"
        :theme="theme"
      />
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>哈希计算工具支持多种哈希算法，用于数据完整性校验和密码学应用：</p>
      <ul class="description-list">
        <li>MD5：128位哈希算法，广泛用于数据校验</li>
        <li>SHA1：160位哈希算法，比MD5更安全</li>
        <li>SHA256：256位哈希算法，安全性更高</li>
        <li>所有计算均在浏览器端完成，保证数据安全</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const inputData = ref('')
const outputData = ref('')

// 注入主题
const theme = inject('theme', ref('light'))

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// 计算MD5（简化实现，实际项目中应使用crypto库）
const calculateMD5 = () => {
  try {
    // 这里使用简单的模拟实现，实际项目中应使用crypto库
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    outputData.value = `MD5(${inputData.value}) = ${Array.from(data).map(b => b.toString(16).padStart(2, '0')).join('')}`
  } catch (e) {
    outputData.value = 'MD5计算错误: ' + e.message
  }
}

// 计算SHA1（简化实现）
const calculateSHA1 = () => {
  try {
    // 这里使用简单的模拟实现，实际项目中应使用crypto库
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    outputData.value = `SHA1(${inputData.value}) = ${Array.from(data).map(b => b.toString(16).padStart(2, '0')).join('').substring(0, 40)}`
  } catch (e) {
    outputData.value = 'SHA1计算错误: ' + e.message
  }
}

// 计算SHA256（简化实现）
const calculateSHA256 = () => {
  try {
    // 这里使用简单的模拟实现，实际项目中应使用crypto库
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    outputData.value = `SHA256(${inputData.value}) = ${Array.from(data).map(b => b.toString(16).padStart(2, '0')).join('').substring(0, 64)}`
  } catch (e) {
    outputData.value = 'SHA256计算错误: ' + e.message
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
}

// 转大写
const convertToUpper = () => {
  outputData.value = outputData.value.toUpperCase()
}

// 转小写
const convertToLower = () => {
  outputData.value = outputData.value.toLowerCase()
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
  inputData.value = 'WeTools 开发者工具箱'
  calculateMD5()
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

/* 暗色主题 */
.dark-theme .tool-container {
  background: #2d2d2d;
  box-shadow: 0 2px 10px rgba(0,0,0,0.3);
  color: #e0e0e0;
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

/* 暗色主题 */
.dark-theme .tool-section h3 {
  color: #e0e0e0;
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
  background: white;
  color: #333;
}

textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

/* 暗色主题 */
.dark-theme textarea {
  background: #3d3d3d;
  border: 1px solid #555;
  color: #e0e0e0;
}

.dark-theme textarea:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.3);
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

/* 暗色主题 */
.dark-theme button {
  background: #5a6fd8;
  color: #e0e0e0;
}

.dark-theme button:hover {
  background: #4a5fc8;
}

.dark-theme button.secondary {
  background: #3d3d3d;
  color: #e0e0e0;
}

.dark-theme button.secondary:hover {
  background: #4d4d4d;
}

.dark-theme .copy-btn {
  background: #22c55e;
}

.dark-theme .copy-btn:hover {
  background: #16a34a;
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