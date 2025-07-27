<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>随机字符串生成器</h2>
      <p>自定义长度和字符集的随机字符串生成器</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="length-section">
            <el-form-item label="字符串长度">
              <el-input-number
                v-model="length"
                :min="1"
                :max="1000"
                controls-position="right"
              />
            </el-form-item>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="count-section">
            <el-form-item label="生成数量">
              <el-input-number
                v-model="count"
                :min="1"
                :max="100"
                controls-position="right"
              />
            </el-form-item>
          </div>
        </el-col>
      </el-row>
      
      <div class="charset-section">
        <h3>字符集选择</h3>
        <div class="charset-options">
          <el-checkbox v-model="charsets.uppercase">大写字母 (A-Z)</el-checkbox>
          <el-checkbox v-model="charsets.lowercase">小写字母 (a-z)</el-checkbox>
          <el-checkbox v-model="charsets.numbers">数字 (0-9)</el-checkbox>
          <el-checkbox v-model="charsets.symbols">符号 (!@#$%^&*)</el-checkbox>
        </div>
        
        <div class="custom-charset">
          <el-input
            v-model="customCharset"
            placeholder="自定义字符集 (可选)"
          />
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="generateRandomStrings">
          <i class="fas fa-random"></i> 生成随机字符串
        </el-button>
        <el-button @click="clearResults">
          <i class="fas fa-trash"></i> 清空结果
        </el-button>
        <el-button type="success" @click="copyAllResults">
          <i class="fas fa-copy"></i> 复制所有结果
        </el-button>
      </div>
      
      <div class="results-section" v-if="results.length > 0">
        <h3>生成结果</h3>
        <div class="results-list">
          <div
            v-for="(result, index) in results"
            :key="index"
            class="result-item"
          >
            <el-input
              v-model="results[index]"
              readonly
            />
            <el-button
              type="success"
              @click="copySingleResult(index)"
              circle
            >
              <i class="fas fa-copy"></i>
            </el-button>
          </div>
        </div>
      </div>
      
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
import { ElCard, ElRow, ElCol, ElFormItem, ElInputNumber, ElCheckbox, ElInput, ElButton, ElAlert } from 'element-plus'

// 数据模型
const length = ref(16)
const count = ref(1)
const charsets = ref({
  uppercase: true,
  lowercase: true,
  numbers: true,
  symbols: false
})
const customCharset = ref('')
const results = ref([])
const validationResult = ref(null)

// 生成随机字符串
const generateRandomStrings = async () => {
  try {
    // 检查是否选择了字符集
    const hasDefaultCharset = charsets.value.uppercase || charsets.value.lowercase || 
                             charsets.value.numbers || charsets.value.symbols
    const hasCustomCharset = customCharset.value.trim() !== ''
    
    if (!hasDefaultCharset && !hasCustomCharset) {
      validationResult.value = {
        type: 'warning',
        message: '请至少选择一个字符集或输入自定义字符集'
      }
      return
    }
    
    // 构建字符集
    let charset = ''
    if (hasCustomCharset) {
      charset = customCharset.value
    } else {
      if (charsets.value.uppercase) charset += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
      if (charsets.value.lowercase) charset += 'abcdefghijklmnopqrstuvwxyz'
      if (charsets.value.numbers) charset += '0123456789'
      if (charsets.value.symbols) charset += '!@#$%^&*()_+-=[]{}|;:,.<>?'
    }
    
    // 生成随机字符串
    const generatedResults = []
    for (let i = 0; i < count.value; i++) {
      let result = ''
      // 使用Web Crypto API生成更安全的随机数
      const randomBytes = new Uint8Array(length.value)
      if (window.crypto && window.crypto.getRandomValues) {
        // 使用Web Crypto API
        window.crypto.getRandomValues(randomBytes)
        for (let j = 0; j < length.value; j++) {
          const randomIndex = randomBytes[j] % charset.length
          result += charset.charAt(randomIndex)
        }
      } else {
        // 降级到Math.random()
        for (let j = 0; j < length.value; j++) {
          const randomIndex = Math.floor(Math.random() * charset.length)
          result += charset.charAt(randomIndex)
        }
      }
      generatedResults.push(result)
    }
    
    results.value = generatedResults
    validationResult.value = {
      type: 'success',
      message: `成功生成${count.value}个随机字符串！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `生成失败: ${error.message}`
    }
  }
}

// 清空结果
const clearResults = () => {
  results.value = []
  validationResult.value = null
}

// 复制单个结果
const copySingleResult = async (index) => {
  try {
    await navigator.clipboard.writeText(results.value[index])
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

// 复制所有结果
const copyAllResults = async () => {
  if (results.value.length === 0) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    const allResults = results.value.join('\n')
    await navigator.clipboard.writeText(allResults)
    validationResult.value = {
      type: 'success',
      message: '所有结果已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 初始化示例数据
length.value = 16
count.value = 5
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

.length-section,
.count-section {
  width: 100%;
}

.charset-section {
  width: 100%;
}

.charset-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.charset-options {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.custom-charset {
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

.results-section {
  width: 100%;
}

.results-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.result-item {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.result-item .el-input {
  flex: 1;
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
  
  .result-item {
    flex-direction: column;
  }
  
  .result-item .el-button {
    align-self: flex-end;
  }
}
</style>