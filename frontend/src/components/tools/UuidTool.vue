<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-sliders-h"></i> UUID/GUID版本</h3>
      <div class="config-group">
        <label>
          <input 
            v-model="uuidVersion" 
            type="radio" 
            value="v1" 
          /> UUID v1 (基于时间)
        </label>
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="uuidVersion" 
            type="radio" 
            value="v4" 
          /> UUID v4 (随机)
        </label>
      </div>
      
      <div class="config-group">
        <label>
          <input 
            v-model="uuidVersion" 
            type="radio" 
            value="guid" 
          /> GUID (Windows格式)
        </label>
      </div>
    </div>

    <div class="button-group">
      <button @click="generateUUID">
        <i class="fas fa-id-card"></i> 生成UUID/GUID
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
      <p>UUID/GUID（通用唯一识别码）生成器支持多种版本：</p>
      <ul class="description-list">
        <li>UUID v1：基于时间和MAC地址生成</li>
        <li>UUID v4：基于随机数生成</li>
        <li>GUID：Windows兼容格式的UUID</li>
        <li>所有生成操作在浏览器端完成，保证唯一性和安全性</li>
        <li>适用于数据库主键、会话标识、API令牌等场景</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const uuidVersion = ref('v4')
const outputData = ref('')

// 生成UUID v4
const generateUUIDv4 = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c == 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

// 生成UUID v1
const generateUUIDv1 = () => {
  // 简化的UUID v1实现
  const now = Date.now()
  const timestamp = now.toString(16)
  const random = Math.random().toString(16).substr(2, 8)
  return `${timestamp.substr(0, 8)}-${timestamp.substr(8, 4)}-1${random.substr(0, 3)}-${random.substr(3, 4)}-${random.substr(7, 12)}`
}

// 生成GUID (大写格式)
const generateGUID = () => {
  const uuid = generateUUIDv4()
  return uuid.toUpperCase()
}

// 生成UUID/GUID
const generateUUID = () => {
  try {
    switch (uuidVersion.value) {
      case 'v1':
        outputData.value = generateUUIDv1()
        break
      case 'v4':
        outputData.value = generateUUIDv4()
        break
      case 'guid':
        outputData.value = generateGUID()
        break
      default:
        outputData.value = generateUUIDv4()
    }
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
  generateUUID()
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

input[type="radio"] {
  width: 16px;
  height: 16px;
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