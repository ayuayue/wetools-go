<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-sliders-h"></i> 配置选项</h3>
      <div class="config-group">
        <label>长度:</label>
        <input 
          v-model.number="length" 
          type="number" 
          min="1" 
          max="1000" 
        />
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="includeUppercase" 
            type="checkbox" 
          /> 大写字母 (A-Z)
        </label>
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="includeLowercase" 
            type="checkbox" 
          /> 小写字母 (a-z)
        </label>
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="includeNumbers" 
            type="checkbox" 
          /> 数字 (0-9)
        </label>
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="includeSymbols" 
            type="checkbox" 
          /> 特殊符号 (!@#$%^&*)
        </label>
      </div>
    </div>

    <div class="button-group">
      <button @click="generateRandomString">
        <i class="fas fa-random"></i> 生成随机字符串
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <div class="result-header">
        <h3><i class="fas fa-file-alt"></i> 生成结果</h3>
      </div>
      <CodeBlock 
        :code="outputData" 
        language="text"
        :show-line-numbers="true"
        :show-header="true"
      />
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>随机字符串生成器可以根据自定义规则生成安全的随机字符串：</p>
      <ul class="description-list">
        <li>支持自定义字符串长度</li>
        <li>可选择包含大写字母、小写字母、数字、特殊符号</li>
        <li>适用于密码、验证码、API密钥等场景</li>
        <li>所有生成操作在浏览器端完成，保证安全性</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const length = ref(16)
const includeUppercase = ref(true)
const includeLowercase = ref(true)
const includeNumbers = ref(true)
const includeSymbols = ref(true)
const outputData = ref('')

// 生成随机字符串
const generateRandomString = () => {
  try {
    let charset = ''
    if (includeUppercase.value) charset += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    if (includeLowercase.value) charset += 'abcdefghijklmnopqrstuvwxyz'
    if (includeNumbers.value) charset += '0123456789'
    if (includeSymbols.value) charset += '!@#$%^&*()_+-=[]{}|;:,.<>?'
    
    if (!charset) {
      outputData.value = '请至少选择一种字符类型'
      return
    }
    
    let result = ''
    const values = new Uint32Array(length.value)
    crypto.getRandomValues(values)
    
    for (let i = 0; i < length.value; i++) {
      result += charset[values[i] % charset.length]
    }
    
    outputData.value = result
  } catch (e) {
    outputData.value = '生成错误: ' + e.message
  }
}

// 清空数据
const clearData = () => {
  outputData.value = ''
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
  generateRandomString()
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

.config-group {
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 10px;
}

.config-group label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
}

input[type="number"] {
  width: 80px;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  background: white;
  color: #333;
}

input[type="number"]:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

/* 暗色主题 */
.dark-theme input[type="number"] {
  background: #3d3d3d;
  border: 1px solid #555;
  color: #e0e0e0;
}

.dark-theme input[type="number"]:focus {
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
  min-height: 50px;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow: auto;
  max-height: 200px;
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
  
  .config-group {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>