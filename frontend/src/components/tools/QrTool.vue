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
            <h3>生成的二维码</h3>
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
              >
                <i class="fas fa-cloud-upload-alt"></i>
                <div class="el-upload__text">
                  将二维码图片拖到此处，或<em>点击上传</em>
                </div>
                <div class="el-upload__tip" slot="tip">
                  支持 JPG、PNG 格式的二维码图片
                </div>
              </el-upload>
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
import { ref } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElFormItem, ElInput, ElSlider, ElSelect, ElOption, ElColorPicker, ElButton, ElAlert, ElUpload } from 'element-plus'

// 数据模型
const activeTab = ref('generate')
const qrData = ref('')
const qrSize = ref(200)
const qrErrorCorrection = ref('M')
const qrForegroundColor = ref('#000000')
const qrBackgroundColor = ref('#ffffff')
const qrCodeUrl = ref('')
const qrImageFile = ref(null)
const decodedData = ref('')
const validationResult = ref(null)

// 生成二维码
const generateQRCode = () => {
  try {
    if (!qrData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要生成二维码的内容'
      }
      return
    }
    
    // 使用在线API生成二维码
    const apiUrl = `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(qrData.value)}&size=${qrSize.value}x${qrSize.value}&ecc=${qrErrorCorrection.value}&color=${qrForegroundColor.value.replace('#', '')}&bgcolor=${qrBackgroundColor.value.replace('#', '')}`
    qrCodeUrl.value = apiUrl
    
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
  decodedData.value = ''
  validationResult.value = null
}

// 解析二维码
const decodeQRCode = async () => {
  if (!qrImageFile.value) {
    validationResult.value = {
      type: 'warning',
      message: '请先上传二维码图片'
    }
    return
  }
  
  try {
    // 使用在线API解析二维码
    const formData = new FormData()
    formData.append('file', qrImageFile.value)
    
    const response = await fetch('https://api.qrserver.com/v1/read-qr-code/', {
      method: 'POST',
      body: formData
    })
    
    const result = await response.json()
    
    if (result && result[0] && result[0].symbol && result[0].symbol[0]) {
      if (result[0].symbol[0].data) {
        decodedData.value = result[0].symbol[0].data
        validationResult.value = {
          type: 'success',
          message: '二维码解析成功！'
        }
      } else {
        validationResult.value = {
          type: 'warning',
          message: '未在二维码中找到有效数据'
        }
      }
    } else {
      validationResult.value = {
        type: 'error',
        message: '二维码解析失败'
      }
    }
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