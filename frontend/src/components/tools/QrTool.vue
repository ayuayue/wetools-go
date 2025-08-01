<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>二维码生成器</h2>
      <p>二维码生成与解析工具</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="tool-tabs">
        <el-tab-pane label="生成二维码" name="generate">
          <el-row :gutter="20">
            <el-col :span="24">
              <div class="input-section">
                <el-input
                  v-model="qrData"
                  type="textarea"
                  :rows="6"
                  placeholder="请输入要生成二维码的内容，例如：https://example.com"
                  resize="vertical"
                />
              </div>
            </el-col>
          </el-row>
          
          <div class="qr-options">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="option-item">
                  <el-form-item label="尺寸">
                    <el-slider
                      v-model="qrSize"
                      :min="100"
                      :max="500"
                      :step="10"
                      show-input
                    />
                  </el-form-item>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="option-item">
                  <el-form-item label="纠错级别">
                    <el-select v-model="qrErrorCorrection">
                      <el-option label="低 (7%)" value="L"></el-option>
                      <el-option label="中 (15%)" value="M"></el-option>
                      <el-option label="高 (25%)" value="Q"></el-option>
                      <el-option label="最高 (30%)" value="H"></el-option>
                    </el-select>
                  </el-form-item>
                </div>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="option-item">
                  <el-form-item label="前景色">
                    <el-color-picker v-model="qrForegroundColor" />
                  </el-form-item>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="option-item">
                  <el-form-item label="背景色">
                    <el-color-picker v-model="qrBackgroundColor" />
                  </el-form-item>
                </div>
              </el-col>
            </el-row>
          </div>
          
          <div class="button-group">
            <el-button type="primary" @click="generateQRCode">
              <i class="fas fa-qrcode"></i> 生成二维码
            </el-button>
            <el-button @click="clearQRData">
              <i class="fas fa-trash"></i> 清空
            </el-button>
            <el-button type="success" @click="downloadQRCode">
              <i class="fas fa-download"></i> 下载二维码
            </el-button>
          </div>
          
          <div class="qr-result" v-if="qrCodeUrl">
            <div class="qr-image-container">
              <img :src="qrCodeUrl" :alt="qrData" />
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="解析二维码" name="decode">
          <div class="decode-section">
            <div class="file-upload-area">
              <el-upload
                drag
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleQRImageUpload"
                @paste="handlePaste"
              >
                <i class="fas fa-cloud-upload-alt"></i>
                <div class="el-upload__text">
                  将二维码图片拖到此处，或<em>点击上传</em>
                </div>
                <div class="el-upload__tip">
                  支持 JPG、PNG 格式的二维码图片，可直接粘贴图片
                </div>
              </el-upload>
            </div>
            
            <div class="image-preview" v-if="qrImagePreview">
              <h3>图片预览</h3>
              <img :src="qrImagePreview" alt="二维码图片预览" />
            </div>
            
            <div class="button-group">
              <el-button type="primary" @click="decodeQRCode" :disabled="!qrImageFile">
                <i class="fas fa-search"></i> 解析二维码
              </el-button>
              <el-button @click="clearDecodeData">
                <i class="fas fa-trash"></i> 清空
              </el-button>
            </div>
            
            <div class="decode-result" v-if="decodedData">
              <h3>解析结果</h3>
              <el-input
                v-model="decodedData"
                type="textarea"
                :rows="4"
                readonly
              />
              <div class="button-group">
                <el-button type="success" @click="copyDecodedData">
                  <i class="fas fa-copy"></i> 复制内容
                </el-button>
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
import { ref, onMounted, onUnmounted } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElFormItem, ElInput, ElSlider, ElSelect, ElOption, ElColorPicker, ElButton, ElAlert, ElUpload } from 'element-plus'
import QRCode from 'qrcode'
import jsQR from 'jsqr'

// 数据模型
const activeTab = ref('generate')
const qrData = ref('')
const qrSize = ref(200)
const qrErrorCorrection = ref('M')
const qrForegroundColor = ref('#000000')
const qrBackgroundColor = ref('#ffffff')
const qrCodeUrl = ref('')
const qrImageFile = ref(null)
const qrImagePreview = ref('')
const decodedData = ref('')
const validationResult = ref(null)

// 生成二维码
const generateQRCode = async () => {
  try {
    if (!qrData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要生成二维码的内容'
      }
      return
    }
    
    // 使用本地生成二维码
    const options = {
      width: qrSize.value,
      margin: 2,
      color: {
        dark: qrForegroundColor.value,
        light: qrBackgroundColor.value
      },
      errorCorrectionLevel: qrErrorCorrection.value
    }
    
    qrCodeUrl.value = await QRCode.toDataURL(qrData.value, options)
    
    validationResult.value = {
      type: 'success',
      message: '二维码生成成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `二维码生成失败: ${error.message}`
    }
  }
}

// 清空二维码数据
const clearQRData = () => {
  qrData.value = ''
  qrCodeUrl.value = ''
  validationResult.value = null
}

// 下载二维码
const downloadQRCode = () => {
  if (!qrCodeUrl.value) {
    validationResult.value = {
      type: 'warning',
      message: '请先生成二维码'
    }
    return
  }
  
  try {
    const link = document.createElement('a')
    link.href = qrCodeUrl.value
    link.download = 'qrcode.png'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    validationResult.value = {
      type: 'success',
      message: '二维码下载成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `二维码下载失败: ${error.message}`
    }
  }
}

// 处理二维码图片上传
const handleQRImageUpload = (file) => {
  qrImageFile.value = file.raw
  // 创建图片预览
  const reader = new FileReader()
  reader.onload = (e) => {
    qrImagePreview.value = e.target.result
  }
  reader.readAsDataURL(file.raw)
  decodedData.value = ''
  validationResult.value = null
}

// 处理粘贴事件
const handlePaste = (event) => {
  const items = (event.clipboardData || event.originalEvent.clipboardData).items
  for (let i = 0; i < items.length; i++) {
    if (items[i].type.indexOf('image') !== -1) {
      const blob = items[i].getAsFile()
      if (blob) {
        qrImageFile.value = blob
        // 创建图片预览
        const reader = new FileReader()
        reader.onload = (e) => {
          qrImagePreview.value = e.target.result
        }
        reader.readAsDataURL(blob)
        decodedData.value = ''
        validationResult.value = null
      }
      break
    }
  }
}

// 全局粘贴处理
const handleGlobalPaste = (event) => {
  // 只在解析二维码标签页处理粘贴
  if (activeTab.value === 'decode') {
    const items = (event.clipboardData || event.originalEvent.clipboardData).items
    for (let i = 0; i < items.length; i++) {
      if (items[i].type.indexOf('image') !== -1) {
        event.preventDefault() // 阻止默认粘贴行为
        const blob = items[i].getAsFile()
        if (blob) {
          qrImageFile.value = blob
          // 创建图片预览
          const reader = new FileReader()
          reader.onload = (e) => {
            qrImagePreview.value = e.target.result
          }
          reader.readAsDataURL(blob)
          decodedData.value = ''
          validationResult.value = {
            type: 'success',
            message: '图片已粘贴！'
          }
        }
        break
      }
    }
  }
}

// 解析二维码
const decodeQRCode = () => {
  if (!qrImageFile.value) {
    validationResult.value = {
      type: 'warning',
      message: '请先上传二维码图片'
    }
    return
  }
  
  try {
    // 使用FileReader读取图片文件
    const reader = new FileReader()
    reader.onload = (e) => {
      // 创建图片元素来获取尺寸
      const img = new Image()
      img.onload = () => {
        // 创建canvas来处理图片
        const canvas = document.createElement('canvas')
        const ctx = canvas.getContext('2d')
        canvas.width = img.width
        canvas.height = img.height
        ctx.drawImage(img, 0, 0, img.width, img.height)
        
        // 获取图片数据
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
        
        // 使用jsQR解析二维码
        const code = jsQR(imageData.data, imageData.width, imageData.height)
        
        if (code) {
          decodedData.value = code.data
          validationResult.value = {
            type: 'success',
            message: '二维码解析成功！'
          }
        } else {
          validationResult.value = {
            type: 'error',
            message: '未找到有效的二维码'
          }
        }
      }
      img.onerror = () => {
        validationResult.value = {
          type: 'error',
          message: '图片加载失败'
        }
      }
      img.src = e.target.result
    }
    reader.onerror = () => {
      validationResult.value = {
        type: 'error',
        message: '图片读取失败'
      }
    }
    reader.readAsDataURL(qrImageFile.value)
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `二维码解析失败: ${error.message}`
    }
  }
}

// 清空解析数据
const clearDecodeData = () => {
  qrImageFile.value = null
  qrImagePreview.value = ''
  decodedData.value = ''
  validationResult.value = null
}

// 复制解析结果
const copyDecodedData = async () => {
  if (!decodedData.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(decodedData.value)
    validationResult.value = {
      type: 'success',
      message: '解析结果已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 初始化示例数据
qrData.value = 'https://example.com'

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

.input-section {
  width: 100%;
}

.qr-options {
  width: 100%;
}

.option-item {
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

.qr-result {
  width: 100%;
  text-align: center;
}

.qr-result h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.qr-image-container {
  display: flex;
  justify-content: center;
  margin: 1rem 0;
}

.qr-image-container img {
  max-width: 100%;
  height: auto;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.decode-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.paste-area {
  width: 100%;
  margin-bottom: 1rem;
}

.file-upload-area {
  width: 100%;
}

.decode-result {
  width: 100%;
}

.decode-result h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.image-preview {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f5f7fa;
  border-radius: 4px;
  text-align: center;
}

.image-preview h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.image-preview img {
  max-width: 100%;
  max-height: 300px;
  border: 1px solid #ddd;
  border-radius: 4px;
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