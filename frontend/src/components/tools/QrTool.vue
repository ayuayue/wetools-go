<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-qrcode"></i> 二维码内容</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入要生成二维码的内容（URL、文本等）"
      ></textarea>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-sliders-h"></i> 配置选项</h3>
      <div class="config-group">
        <label>尺寸:</label>
        <input 
          v-model.number="size" 
          type="range" 
          min="100" 
          max="500" 
          step="50"
        />
        <span>{{ size }}px</span>
      </div>
      
      <div class="config-group">
        <label>纠错级别:</label>
        <select v-model="errorCorrectionLevel">
          <option value="L">L (7%)</option>
          <option value="M">M (15%)</option>
          <option value="Q">Q (25%)</option>
          <option value="H">H (30%)</option>
        </select>
      </div>
    </div>

    <div class="button-group">
      <button @click="generateQRCode">
        <i class="fas fa-qrcode"></i> 生成二维码
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <div class="result-header">
        <h3><i class="fas fa-image"></i> 二维码</h3>
      </div>
      <div class="qr-code-container">
        <img 
          v-if="qrCodeDataUrl" 
          :src="qrCodeDataUrl" 
          :style="{ width: size + 'px', height: size + 'px' }"
          alt="QR Code"
        />
        <div v-else class="qr-placeholder">
          <i class="fas fa-qrcode"></i>
          <p>生成二维码后将显示在这里</p>
        </div>
      </div>
      <div class="result-footer" v-if="qrCodeDataUrl">
        <button class="copy-btn" @click="downloadQRCode">
          <i class="fas fa-download"></i> 下载二维码
        </button>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>二维码生成器可以将文本、URL等信息转换为二维码图片：</p>
      <ul class="description-list">
        <li>支持自定义尺寸和纠错级别</li>
        <li>可生成网站链接、文本、联系方式等二维码</li>
        <li>支持下载二维码图片</li>
        <li>所有生成操作在浏览器端完成，保证隐私安全</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import QRCode from 'qrcode'

const inputData = ref('')
const size = ref(200)
const errorCorrectionLevel = ref('M')
const qrCodeDataUrl = ref('')

// 生成二维码
const generateQRCode = async () => {
  if (!inputData.value) {
    qrCodeDataUrl.value = ''
    return
  }
  
  try {
    const dataUrl = await QRCode.toDataURL(inputData.value, {
      width: size.value,
      height: size.value,
      errorCorrectionLevel: errorCorrectionLevel.value,
      margin: 2,
      color: {
        dark: '#000000',
        light: '#ffffff'
      }
    })
    
    qrCodeDataUrl.value = dataUrl
  } catch (e) {
    console.error('生成二维码错误:', e)
    qrCodeDataUrl.value = ''
  }
}

// 下载二维码
const downloadQRCode = () => {
  if (!qrCodeDataUrl.value) return
  
  const link = document.createElement('a')
  link.href = qrCodeDataUrl.value
  link.download = 'qrcode.png'
  link.click()
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  qrCodeDataUrl.value = ''
}

// 初始化数据
onMounted(() => {
  inputData.value = 'https://github.com/yourusername/wetools-go'
  generateQRCode()
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
  min-height: 100px;
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

.config-group {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 1rem;
}

.config-group label {
  min-width: 80px;
}

.config-group input[type="range"] {
  flex: 1;
}

.config-group select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
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

.qr-code-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 4px;
  padding: 1rem;
}

.qr-placeholder {
  text-align: center;
  color: #94a3b8;
}

.qr-placeholder i {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.qr-placeholder p {
  margin: 0;
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