<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入内容</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入要编码/解码的内容"
      ></textarea>
    </div>

    <!-- 文件上传区域 -->
    <div class="tool-section">
      <h3><i class="fas fa-file-upload"></i> 文件上传</h3>
      <div class="file-upload-area" @dragover.prevent @drop.prevent="handleFileDrop">
        <input 
          type="file" 
          ref="fileInput" 
          @change="handleFileSelect" 
          style="display: none;"
        />
        <div class="upload-content" @click="triggerFileInput">
          <i class="fas fa-cloud-upload-alt upload-icon"></i>
          <p>点击选择文件或拖拽文件到此处</p>
          <p class="file-info" v-if="selectedFile">{{ selectedFile.name }} ({{ formatFileSize(selectedFile.size) }})</p>
        </div>
      </div>
    </div>

    <div class="button-group">
      <button @click="encodeBase64">
        <i class="fas fa-lock"></i> Base64编码
      </button>
      <button @click="decodeBase64">
        <i class="fas fa-unlock"></i> Base64解码
      </button>
      <button @click="encodeFileToBase64">
        <i class="fas fa-file-image"></i> 文件转Base64
      </button>
      <button @click="decodeBase64ToFile">
        <i class="fas fa-file-download"></i> Base64转文件
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
      <div v-if="outputImage" class="image-preview">
        <img :src="outputImage" alt="预览图片" />
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
      <p>Base64是一种基于64个可打印字符来表示二进制数据的表示方法。本工具可以帮助您：</p>
      <ul class="description-list">
        <li>将文本或二进制数据编码为Base64格式</li>
        <li>将Base64编码的数据解码为原始内容</li>
        <li>支持文件转Base64和Base64转文件</li>
        <li>支持图片预览功能</li>
        <li>支持中文、英文、特殊字符等各类内容</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const inputData = ref('')
const outputData = ref('')
const selectedFile = ref(null)
const outputImage = ref('')
const fileInput = ref(null)

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// Base64编码
const encodeBase64 = () => {
  try {
    outputData.value = btoa(unescape(encodeURIComponent(inputData.value)))
    // 如果是图片Base64，显示预览
    if (outputData.value.startsWith('data:image')) {
      outputImage.value = outputData.value
    } else {
      outputImage.value = ''
    }
  } catch (e) {
    outputData.value = '编码错误: ' + e.message
    outputImage.value = ''
  }
}

// Base64解码
const decodeBase64 = () => {
  try {
    outputData.value = decodeURIComponent(escape(atob(inputData.value)))
    outputImage.value = ''
  } catch (e) {
    outputData.value = '解码错误: ' + e.message
    outputImage.value = ''
  }
}

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value.click()
}

// 处理文件选择
const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    selectedFile.value = file
  }
}

// 处理文件拖拽
const handleFileDrop = (event) => {
  const file = event.dataTransfer.files[0]
  if (file) {
    selectedFile.value = file
    fileInput.value.files = event.dataTransfer.files
  }
}

// 文件转Base64
const encodeFileToBase64 = () => {
  if (!selectedFile.value) {
    outputData.value = '请先选择文件'
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    outputData.value = e.target.result
    // 如果是图片文件，显示预览
    if (selectedFile.value.type.startsWith('image/')) {
      outputImage.value = e.target.result
    } else {
      outputImage.value = ''
    }
  }
  reader.onerror = () => {
    outputData.value = '文件读取错误'
    outputImage.value = ''
  }
  reader.readAsDataURL(selectedFile.value)
}

// Base64转文件
const decodeBase64ToFile = () => {
  try {
    // 创建下载链接
    const link = document.createElement('a')
    link.href = inputData.value
    link.download = selectedFile.value ? selectedFile.value.name : 'download'
    link.click()
  } catch (e) {
    outputData.value = '文件下载错误: ' + e.message
  }
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
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
  inputData.value = 'WeTools 开发者工具箱'
  encodeBase64()
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

/* 文件上传区域样式 */
.file-upload-area {
  border: 2px dashed #cbd5e1;
  border-radius: 4px;
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
  background: #f8fafc;
}

.file-upload-area:hover {
  border-color: #667eea;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.upload-icon {
  font-size: 2rem;
  color: #94a3b8;
}

.file-info {
  font-size: 0.9rem;
  color: #64748b;
  margin-top: 0.5rem;
}

/* 图片预览样式 */
.image-preview {
  text-align: center;
  margin-bottom: 1rem;
}

.image-preview img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
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
  
  .file-upload-area {
    padding: 1rem;
  }
}
</style>