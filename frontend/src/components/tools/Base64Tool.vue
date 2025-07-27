<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>Base64 编码/解码工具</h2>
      <p>Base64编码解码工具，支持文本和文件转换</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="tool-tabs">
        <el-tab-pane label="文本转换" name="text">
          <el-row :gutter="20">
            <el-col :span="24">
              <div class="input-section">
                <el-input
                  v-model="textInput"
                  type="textarea"
                  :rows="8"
                  placeholder="请输入要编码或解码的文本"
                  resize="vertical"
                />
              </div>
            </el-col>
          </el-row>
          
          <div class="button-group">
            <el-button type="primary" @click="encodeText">
              <i class="fas fa-lock"></i> 编码为Base64
            </el-button>
            <el-button @click="decodeText">
              <i class="fas fa-unlock"></i> 解码Base64
            </el-button>
            <el-button @click="clearText">
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
                  v-model="textOutput"
                  type="textarea"
                  :rows="8"
                  placeholder="转换结果将显示在这里"
                  resize="vertical"
                />
              </div>
            </el-col>
          </el-row>
        </el-tab-pane>
        
        <el-tab-pane label="文件转换" name="file">
          <div class="file-conversion-section">
            <div class="file-upload-area">
              <el-upload
                drag
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleFileChange"
                @paste="handlePaste"
              >
                <i class="fas fa-cloud-upload-alt"></i>
                <div class="el-upload__text">
                  将文件拖到此处，或<em>点击上传</em>
                </div>
                <div class="el-upload__tip">
                  支持所有文件类型，可直接粘贴图片或文件
                </div>
              </el-upload>
            </div>
            
            <div class="file-info" v-if="fileInfo">
              <el-descriptions title="文件信息" :column="1" border>
                <el-descriptions-item label="文件名">{{ fileInfo.name }}</el-descriptions-item>
                <el-descriptions-item label="文件大小">{{ formatFileSize(fileInfo.size) }}</el-descriptions-item>
                <el-descriptions-item label="文件类型">{{ fileInfo.type }}</el-descriptions-item>
              </el-descriptions>
            </div>
            
            <div class="base64-input" v-if="activeTab === 'file'">
              <el-input
                v-model="base64Input"
                type="textarea"
                :rows="6"
                placeholder="粘贴Base64字符串以转换为文件"
                resize="vertical"
              />
              <div class="button-group">
                <el-button type="primary" @click="decodeBase64String">
                  <i class="fas fa-file-download"></i> Base64转文件
                </el-button>
              </div>
            </div>
            
            <div class="button-group">
              <el-button type="primary" @click="encodeFile" :disabled="!fileInfo">
                <i class="fas fa-lock"></i> 文件转Base64
              </el-button>
              <el-button @click="clearFile">
                <i class="fas fa-trash"></i> 清空
              </el-button>
              <el-button type="success" @click="downloadFile" :disabled="!decodedFileContent && !decodedFromBase64">
                <i class="fas fa-download"></i> 下载文件
              </el-button>
            </div>
            
            <div class="file-result" v-if="fileResult || decodedFileName">
              <el-input
                v-model="fileResult"
                type="textarea"
                :rows="6"
                placeholder="Base64编码结果将显示在这里"
                resize="vertical"
              />
              <div class="decoded-file-info" v-if="decodedFileName">
                <p>解码文件名: {{ decodedFileName }}</p>
                <p>解码文件大小: {{ formatFileSize(decodedFileSize) }}</p>
              </div>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      
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
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElInput, ElButton, ElAlert, ElUpload, ElDescriptions, ElDescriptionsItem } from 'element-plus'

// 数据模型
const activeTab = ref('text')
const textInput = ref('')
const textOutput = ref('')
const fileInfo = ref(null)
const fileResult = ref('')
const base64Input = ref('')
const decodedFileContent = ref(null)
const decodedFromBase64 = ref(null)
const decodedFileName = ref('')
const decodedFileSize = ref(0)
const validationResult = ref(null)

// 编码文本为Base64
const encodeText = () => {
  try {
    if (!textInput.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要编码的文本'
      }
      return
    }
    
    textOutput.value = btoa(unescape(encodeURIComponent(textInput.value)))
    validationResult.value = {
      type: 'success',
      message: '文本编码成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `编码失败: ${error.message}`
    }
  }
}

// 解码Base64文本
const decodeText = () => {
  try {
    if (!textInput.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要解码的Base64字符串'
      }
      return
    }
    
    textOutput.value = decodeURIComponent(escape(atob(textInput.value)))
    validationResult.value = {
      type: 'success',
      message: '文本解码成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `解码失败: ${error.message}`
    }
  }
}

// 清空文本数据
const clearText = () => {
  textInput.value = ''
  textOutput.value = ''
  validationResult.value = null
}

// 处理文件上传
const handleFileChange = (file) => {
  fileInfo.value = {
    name: file.name,
    size: file.size,
    type: file.type,
    raw: file.raw
  }
  fileResult.value = ''
  decodedFileContent.value = null
  validationResult.value = null
}

// 处理粘贴事件
const handlePaste = (event) => {
  const items = (event.clipboardData || event.originalEvent.clipboardData).items
  for (let i = 0; i < items.length; i++) {
    if (items[i].type.indexOf('image') !== -1 || items[i].type.indexOf('file') !== -1) {
      const blob = items[i].getAsFile()
      if (blob) {
        fileInfo.value = {
          name: blob.name || 'pasted-file',
          size: blob.size,
          type: blob.type,
          raw: blob
        }
        fileResult.value = ''
        decodedFileContent.value = null
        validationResult.value = null
        break
      }
    }
  }
}

// 编码文件为Base64
const encodeFile = () => {
  if (!fileInfo.value) {
    validationResult.value = {
      type: 'warning',
      message: '请先选择一个文件'
    }
    return
  }
  
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const base64String = e.target.result.split(',')[1] // 移除data:*/*;base64,前缀
      fileResult.value = base64String
      validationResult.value = {
        type: 'success',
        message: '文件编码成功！'
      }
    } catch (error) {
      validationResult.value = {
        type: 'error',
        message: `文件编码失败: ${error.message}`
      }
    }
  }
  reader.onerror = () => {
    validationResult.value = {
      type: 'error',
      message: '文件读取失败'
    }
  }
  reader.readAsDataURL(fileInfo.value.raw)
}

// 解码Base64字符串为文件
const decodeBase64String = () => {
  if (!base64Input.value.trim()) {
    validationResult.value = {
      type: 'warning',
      message: '请输入Base64字符串'
    }
    return
  }
  
  try {
    // 尝试解析Base64字符串
    const binaryString = atob(base64Input.value)
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    
    decodedFromBase64.value = new Blob([bytes])
    decodedFileName.value = 'decoded-file.bin'
    decodedFileSize.value = bytes.length
    
    validationResult.value = {
      type: 'success',
      message: 'Base64解码成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `Base64解码失败: ${error.message}`
    }
  }
}

// 下载文件
const downloadFile = () => {
  let blob, filename
  
  // 确定要下载的内容
  if (decodedFromBase64.value) {
    // 从Base64解码的文件
    blob = decodedFromBase64.value
    filename = decodedFileName.value
  } else if (decodedFileContent.value && fileInfo.value) {
    // 从文件解码的内容
    try {
      const binaryString = atob(decodedFileContent.value)
      const bytes = new Uint8Array(binaryString.length)
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i)
      }
      blob = new Blob([bytes])
      filename = `decoded_${fileInfo.value.name}`
    } catch (error) {
      validationResult.value = {
        type: 'error',
        message: `文件解码失败: ${error.message}`
      }
      return
    }
  } else {
    validationResult.value = {
      type: 'warning',
      message: '没有可下载的内容'
    }
    return
  }
  
  try {
    // 创建下载链接
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    validationResult.value = {
      type: 'success',
      message: '文件下载成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `文件下载失败: ${error.message}`
    }
  }
}

// 清空文件数据
const clearFile = () => {
  fileInfo.value = null
  fileResult.value = ''
  base64Input.value = ''
  decodedFileContent.value = null
  decodedFromBase64.value = null
  decodedFileName.value = ''
  decodedFileSize.value = 0
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  const result = activeTab.value === 'text' ? textOutput.value : fileResult.value
  if (!result) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(result)
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

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 初始化示例数据
textInput.value = 'Hello, WeTools!'
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

.tool-tabs {
  width: 100%;
}

.input-section,
.output-section,
.file-result,
.base64-input {
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

.file-conversion-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.file-upload-area {
  width: 100%;
}

.file-info {
  width: 100%;
}

.decoded-file-info {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.decoded-file-info p {
  margin: 0.5rem 0;
  color: #666;
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