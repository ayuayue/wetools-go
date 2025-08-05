<template>
  <div class="tool-container">
    <ToolHeader title="图片转SVG" description="将图片转换为SVG格式" />
    
    <div class="tool-content">
      <div class="upload-section">
        <div 
          class="upload-area" 
          @dragover.prevent="handleDragOver"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <i class="fas fa-cloud-upload-alt upload-icon"></i>
          <p class="upload-text">点击或拖拽上传图片</p>
          <p class="upload-hint">支持 JPG, PNG, GIF 等常见图片格式</p>
          <input 
            ref="fileInput" 
            type="file" 
            accept="image/*" 
            @change="handleFileSelect" 
            style="display: none;"
          />
        </div>
        
        <div v-if="previewUrl" class="preview-section">
          <h3>预览</h3>
          <img :src="previewUrl" alt="预览图片" class="preview-image" />
        </div>
      </div>
      
      <!-- 添加颜色模式选择 -->
      <div class="options-section">
        <div class="form-group">
          <label class="form-label">转换模式:</label>
          <div class="radio-group">
            <label class="radio-item">
              <input 
                type="radio" 
                v-model="colorMode" 
                value="mono" 
                class="radio-input"
              />
              <span class="radio-label">单色</span>
            </label>
            <label class="radio-item">
              <input 
                type="radio" 
                v-model="colorMode" 
                value="color" 
                class="radio-input"
              />
              <span class="radio-label">彩色</span>
            </label>
          </div>
        </div>
        
        <!-- 彩色模式参数 -->
        <div v-if="colorMode === 'color'" class="form-group">
          <label class="form-label">颜色层数 (2-10):</label>
          <input 
            type="number" 
            v-model.number="colorSteps" 
            min="2" 
            max="10" 
            class="form-input"
          />
        </div>
      </div>
      
      <div class="controls-section">
        <button 
          class="btn btn-primary" 
          @click="convertToSvg" 
          :disabled="!selectedFile || isConverting"
        >
          <i class="fas fa-exchange-alt"></i>
          {{ isConverting ? '转换中...' : '转换为SVG' }}
        </button>
        
        <button 
          class="btn btn-secondary" 
          @click="clearAll" 
          :disabled="!selectedFile && !previewUrl"
        >
          <i class="fas fa-trash-alt"></i>
          清空
        </button>
      </div>
      
      <div v-if="svgResult" class="result-section">
        <h3>转换结果</h3>
        <div class="svg-preview" v-html="svgResult"></div>
        <div class="result-actions">
          <button class="btn btn-success" @click="copySvg">
            <i class="fas fa-copy"></i>
            复制SVG代码
          </button>
          <button class="btn btn-info" @click="downloadSvg">
            <i class="fas fa-download"></i>
            下载SVG文件
          </button>
        </div>
      </div>
      
      <div v-if="error" class="error-section">
        <div class="alert alert-danger">
          <i class="fas fa-exclamation-triangle"></i>
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ToolHeader from '../ToolHeader.vue'
import potrace from 'potrace'
import { Buffer } from 'buffer'

export default {
  name: 'ImageToSvgTool',
  components: {
    ToolHeader
  },
  data() {
    return {
      selectedFile: null,
      previewUrl: null,
      svgResult: null,
      isConverting: false,
      error: null,
      colorMode: 'color', // 默认使用彩色模式
      colorSteps: 5 // 默认颜色层数
    }
  },
  methods: {
    triggerFileInput() {
      this.$refs.fileInput.click()
    },
    
    handleDragOver(event) {
      event.currentTarget.classList.add('drag-over')
    },
    
    handleDrop(event) {
      const dropArea = event.currentTarget
      dropArea.classList.remove('drag-over')
      
      const files = event.dataTransfer.files
      if (files.length > 0) {
        this.processFile(files[0])
      }
    },
    
    handleFileSelect(event) {
      const files = event.target.files
      if (files.length > 0) {
        this.processFile(files[0])
      }
    },
    
    processFile(file) {
      // 检查文件类型
      if (!file.type.startsWith('image/')) {
        this.error = '请选择有效的图片文件'
        return
      }
      
      // 检查文件大小 (限制为10MB)
      if (file.size > 10 * 1024 * 1024) {
        this.error = '文件大小不能超过10MB'
        return
      }
      
      this.selectedFile = file
      this.error = null
      
      // 创建预览
      if (this.previewUrl) {
        URL.revokeObjectURL(this.previewUrl)
      }
      this.previewUrl = URL.createObjectURL(file)
    },
    
    async convertToSvg() {
      if (!this.selectedFile) {
        this.error = '请先选择图片文件'
        return
      }
      
      this.isConverting = true
      this.error = null
      this.svgResult = null
      
      try {
        // 读取文件内容并转换为base64
        const arrayBuffer = await this.selectedFile.arrayBuffer()
        const base64Data = btoa(String.fromCharCode(...new Uint8Array(arrayBuffer)))
        
        // 调用Go后端方法进行转换
        const request = {
          imageData: base64Data,
          colorMode: this.colorMode,
          colorSteps: this.colorSteps
        }
        
        // 使用Wails绑定调用后端方法
        const result = await window.go.main.App.ImageToSvg(request)
        this.svgResult = result.Svg
      } catch (err) {
        this.error = '转换失败: ' + err.message
      } finally {
        this.isConverting = false
      }
    },
    
    copySvg() {
      if (!this.svgResult) return
      
      navigator.clipboard.writeText(this.svgResult)
        .then(() => {
          alert('SVG代码已复制到剪贴板')
        })
        .catch(err => {
          this.error = '复制失败: ' + err.message
        })
    },
    
    downloadSvg() {
      if (!this.svgResult) return
      
      const blob = new Blob([this.svgResult], { type: 'image/svg+xml' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'converted-image.svg'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    },
    
    clearAll() {
      this.selectedFile = null
      this.svgResult = null
      this.error = null
      if (this.previewUrl) {
        URL.revokeObjectURL(this.previewUrl)
        this.previewUrl = null
      }
      this.$refs.fileInput.value = ''
    }
  },
  
  beforeUnmount() {
    // 清理对象URL
    if (this.previewUrl) {
      URL.revokeObjectURL(this.previewUrl)
    }
  }
}
</script>

<style scoped>
.tool-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.upload-section {
  margin-bottom: 30px;
}

.upload-area {
  border: 2px dashed #ddd;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: #fafafa;
}

.upload-area:hover,
.upload-area.drag-over {
  border-color: #3498db;
  background-color: #f0f8ff;
}

.upload-icon {
  font-size: 48px;
  color: #3498db;
  margin-bottom: 15px;
}

.upload-text {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 10px;
  color: #333;
}

.upload-hint {
  color: #666;
  font-size: 14px;
}

.preview-section {
  margin-top: 20px;
}

.preview-section h3 {
  margin-bottom: 15px;
  color: #333;
}

.preview-image {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.controls-section {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #3498db;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2980b9;
}

.btn-secondary {
  background-color: #95a5a6;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #7f8c8d;
}

.result-section {
  margin-top: 30px;
}

.result-section h3 {
  margin-bottom: 15px;
  color: #333;
}

.svg-preview {
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 20px;
  background-color: #fafafa;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.result-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
  flex-wrap: wrap;
}

.btn-success {
  background-color: #27ae60;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background-color: #229954;
}

.btn-info {
  background-color: #3498db;
  color: white;
}

.btn-info:hover:not(:disabled) {
  background-color: #2980b9;
}

.error-section {
  margin-top: 20px;
}

.alert {
  padding: 15px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.alert-danger {
  background-color: #fadbd8;
  border-color: #f5b7b1;
  color: #c0392b;
}

@media (max-width: 768px) {
  .tool-container {
    padding: 15px;
  }
  
  .controls-section,
  .result-actions {
    flex-direction: column;
    align-items: center;
  }
  
  .btn {
    width: 100%;
    max-width: 300px;
    justify-content: center;
  }
}
</style>