<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>Base64 编码/解码工具</h2>
      <p>Base64编码解码工具，支持文本和文件转换</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="tool-tabs">
        <el-tab-pane label="文本转换" name="text">
          <!-- 两栏布局 -->
          <div class="two-column-layout">
            <!-- 左侧输入区域 -->
            <div class="column input-column">
              <div class="input-section">
                <div class="editor-header">
                  <h3>文本输入</h3>
                  <div class="editor-actions">
                    <!-- 占位元素，确保与右侧高度一致 -->
                  </div>
                </div>
                <AceEditor
                  v-model="textInput"
                  language="text"
                  :theme="theme"
                  :showHeader="false"
                />
              </div>
            </div>
            
            <!-- 右侧输出区域 -->
            <div class="column output-column">
              <div class="output-section">
                <div class="editor-header">
                  <h3>Base64输出</h3>
                  <div class="editor-actions">
                    <el-button type="success" @click="copyResult">
                      <i class="fas fa-copy"></i> 复制结果
                    </el-button>
                  </div>
                </div>
                <AceEditor
                  v-model="textOutput"
                  language="text"
                  :theme="theme"
                  :readonly="true"
                  :showHeader="false"
                />
              </div>
            </div>
          </div>
          
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
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="文件转换" name="file">
          <el-row :gutter="20">
            <el-col :span="12">
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
              
              <div class="file-preview" v-if="filePreviewUrl">
                <h3>文件预览</h3>
                <div class="preview-container">
                  <img v-if="isImageFile" :src="filePreviewUrl" alt="文件预览" />
                  <div v-else class="file-icon">
                    <i class="fas fa-file"></i>
                    <p>{{ fileInfo.name }}</p>
                  </div>
                </div>
              </div>
              
              <div class="button-group vertical">
                <el-button type="primary" @click="encodeFile" :disabled="!fileInfo" size="large">
                  <i class="fas fa-lock"></i> 文件转Base64 (含前缀)
                </el-button>
                <el-button @click="clearFile" size="large">
                  <i class="fas fa-trash"></i> 清空
                </el-button>
              </div>
            </el-col>
            
            <el-col :span="12">
              <div class="base64-input">
                <el-input
                  v-model="base64Input"
                  type="textarea"
                  :rows="8"
                  placeholder="粘贴Base64字符串以转换为文件（可包含或不包含data:*/*;base64,前缀）"
                  resize="vertical"
                />
              </div>
              
              <div class="button-group vertical">
                <el-button type="primary" @click="decodeBase64String" size="large">
                  <i class="fas fa-file-download"></i> Base64转文件
                </el-button>
              </div>
              
              <div class="decoded-file-section" v-if="decodedFileName">
                <div class="decoded-file-info">
                  <el-descriptions title="解码文件信息" :column="1" border>
                    <el-descriptions-item label="文件名">{{ decodedFileName }}</el-descriptions-item>
                    <el-descriptions-item label="文件大小">{{ formatFileSize(decodedFileSize) }}</el-descriptions-item>
                    <el-descriptions-item label="文件类型">{{ decodedFileType }}</el-descriptions-item>
                  </el-descriptions>
                </div>
                <div class="button-group vertical">
                  <el-button type="success" @click="downloadFile" size="large">
                    <i class="fas fa-download"></i> 下载文件
                  </el-button>
                </div>
              </div>
            </el-col>
          </el-row>
          
          <div class="file-result" v-if="fileResult">
            <h3>Base64编码结果</h3>
            <AceEditor
              v-model="fileResult"
              language="text"
              :theme="theme"
              :readonly="true"
              :showHeader="false"
            />
            <div class="button-group">
              <el-button type="success" @click="copyFileResult" size="small">
                <i class="fas fa-copy"></i> 复制Base64
              </el-button>
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
import { ref, onMounted, onUnmounted, inject } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElInput, ElButton, ElAlert, ElUpload, ElDescriptions, ElDescriptionsItem } from 'element-plus'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 数据模型
const activeTab = ref('text')
const textInput = ref('')
const textOutput = ref('')
const fileInfo = ref(null)
const fileResult = ref('')
const base64Input = ref('')
const decodedFromBase64 = ref(null)
const decodedFileName = ref('')
const decodedFileSize = ref(0)
const decodedFileType = ref('')
const filePreviewUrl = ref('')
const isImageFile = ref(false)
const validationResult = ref(null)
const uploadRef = ref(null)

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
  decodedFromBase64.value = null
  decodedFileName.value = ''
  decodedFileSize.value = 0
  decodedFileType.value = ''
  validationResult.value = null
  
  // 创建文件预览
  createFilePreview(file.raw)
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
        decodedFromBase64.value = null
        decodedFileName.value = ''
        decodedFileSize.value = 0
        decodedFileType.value = ''
        validationResult.value = null
        
        // 创建文件预览
        createFilePreview(blob)
        break
      }
    }
  }
}

// 全局粘贴处理
const handleGlobalPaste = (event) => {
  // 只在文件转换标签页处理粘贴
  if (activeTab.value === 'file') {
    const items = (event.clipboardData || event.originalEvent.clipboardData).items
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1 || items[i].type.indexOf('file') !== -1) {
        event.preventDefault() // 阻止默认粘贴行为
        const blob = items[i].getAsFile()
        if (blob) {
          fileInfo.value = {
            name: blob.name || 'pasted-file',
            size: blob.size,
            type: blob.type,
            raw: blob
          }
          fileResult.value = ''
          decodedFromBase64.value = null
          decodedFileName.value = ''
          decodedFileSize.value = 0
          decodedFileType.value = ''
          validationResult.value = {
            type: 'success',
            message: '文件已粘贴！'
          }
          
          // 创建文件预览
          createFilePreview(blob)
          break
        }
      }
    }
  }
}

// 创建文件预览
const createFilePreview = (file) => {
  // 检查是否为图片文件
  if (file.type && file.type.startsWith('image/')) {
    isImageFile.value = true
    const reader = new FileReader()
    reader.onload = (e) => {
      filePreviewUrl.value = e.target.result
    }
    reader.readAsDataURL(file)
  } else {
    isImageFile.value = false
    filePreviewUrl.value = ''
  }
}

// 编码文件为Base64 (包含data URI前缀)
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
      // 保留完整的data URI前缀
      fileResult.value = e.target.result
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
    let base64String = base64Input.value.trim()
    let mimeType = 'application/octet-stream'
    
    // 检查是否包含data URI前缀
    if (base64String.startsWith('data:')) {
      const parts = base64String.split(',')
      mimeType = parts[0].substring(5) // 移除"data:"前缀
      if (mimeType.endsWith(';base64')) {
        mimeType = mimeType.substring(0, mimeType.length - 7) // 移除";base64"后缀
      }
      base64String = parts[1]
    }
    
    // 尝试解析Base64字符串
    const binaryString = atob(base64String)
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    
    decodedFromBase64.value = new Blob([bytes], { type: mimeType })
    
    // 根据MIME类型推断文件扩展名
    let extension = 'bin'
    if (mimeType.startsWith('image/')) {
      extension = mimeType.split('/')[1]
      if (extension === 'jpeg') extension = 'jpg'
    } else if (mimeType.startsWith('text/')) {
      extension = 'txt'
    } else if (mimeType === 'application/pdf') {
      extension = 'pdf'
    } else if (mimeType === 'application/json') {
      extension = 'json'
    }
    
    decodedFileName.value = `decoded-file.${extension}`
    decodedFileSize.value = bytes.length
    decodedFileType.value = mimeType
    
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
  if (!decodedFromBase64.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有可下载的内容'
    }
    return
  }
  
  try {
    // 创建下载链接
    const url = URL.createObjectURL(decodedFromBase64.value)
    const a = document.createElement('a')
    a.href = url
    a.download = decodedFileName.value
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
  decodedFromBase64.value = null
  decodedFileName.value = ''
  decodedFileSize.value = 0
  decodedFileType.value = ''
  filePreviewUrl.value = ''
  isImageFile.value = false
  validationResult.value = null
}

// 复制文本结果
const copyResult = async () => {
  if (!textOutput.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(textOutput.value)
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

// 复制文件Base64结果
const copyFileResult = async () => {
  if (!fileResult.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(fileResult.value)
    validationResult.value = {
      type: 'success',
      message: 'Base64结果已复制到剪贴板！'
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

// 挂载时添加全局粘贴监听器
onMounted(() => {
  document.addEventListener('paste', handleGlobalPaste)
})

// 卸载时移除全局粘贴监听器
onUnmounted(() => {
  document.removeEventListener('paste', handleGlobalPaste)
})
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

/* 两栏布局 */
.two-column-layout {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.column {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.input-column,
.output-column {
  flex: 1;
}

.input-section,
.output-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  height: 100%;
  flex: 1;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.editor-header h3 {
  margin: 0;
  color: #333;
}

.editor-actions {
  display: flex;
  gap: 0.5rem;
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
  margin: 1rem 0;
}

.button-group.vertical {
  flex-direction: column;
  align-items: stretch;
}

.button-group .el-button {
  flex: 1;
  min-width: 120px;
  max-width: 160px;
}

.button-group.vertical .el-button {
  width: 100%;
  max-width: none;
}

.file-upload-area {
  width: 100%;
}

.file-info {
  width: 100%;
  margin-top: 1rem;
}

.file-preview {
  width: 100%;
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.file-preview h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
}

.preview-container img {
  max-width: 100%;
  max-height: 200px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.file-icon {
  text-align: center;
  color: #666;
}

.file-icon i {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.decoded-file-section {
  width: 100%;
  margin-top: 1rem;
}

.decoded-file-info {
  margin-bottom: 1rem;
}

.file-result {
  width: 100%;
  margin-top: 1rem;
}

.file-result h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .two-column-layout {
    flex-direction: column;
  }
  
  .button-group {
    flex-direction: column;
    align-items: center;
  }
  
  .button-group .el-button {
    width: 100%;
    max-width: none;
  }
  
  .button-group.vertical .el-button {
    width: 100%;
    max-width: none;
  }
}
</style>