<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>哈希计算工具</h2>
      <p>MD5、SHA1、SHA256等哈希算法计算工具</p>
    </div>
    
    <div class="tool-content">
      <!-- 两栏布局 -->
      <div class="two-column-layout">
        <!-- 左侧输入区域 -->
        <div class="column input-column">
          <div class="input-section">
            <div class="editor-header">
              <h3>文本输入</h3>
            </div>
            <AceEditor
              v-model="inputData"
              language="text"
              :theme="theme"
              :showHeader="false"
            />
          </div>
          
          <el-row :gutter="20">
            <el-col :span="24">
              <div class="salt-section">
                <el-input
                  v-model="salt"
                  placeholder="请输入盐值（可选）"
                />
              </div>
            </el-col>
          </el-row>
          
          <div class="hash-algorithms">
            <h3>选择哈希算法</h3>
            <div class="algorithm-buttons">
              <el-button 
                v-for="algorithm in algorithms" 
                :key="algorithm.name"
                :type="selectedAlgorithm === algorithm.name ? 'primary' : ''"
                @click="selectAlgorithm(algorithm.name)"
              >
                {{ algorithm.name }}
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 右侧输出区域 -->
        <div class="column output-column">
          <div class="output-section">
            <div class="editor-header">
              <h3>哈希输出</h3>
              <div class="editor-actions">
                <el-button type="success" @click="copyResult">
                  <i class="fas fa-copy"></i> 复制结果
                </el-button>
              </div>
            </div>
            <AceEditor
              v-model="hashResult"
              language="text"
              :theme="theme"
              :readonly="true"
              :showHeader="false"
            />
          </div>
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="calculateHash">
          <i class="fas fa-calculator"></i> 计算哈希值
        </el-button>
        <el-button @click="calculateAllHashes">
          <i class="fas fa-calculator"></i> 计算所有哈希值
        </el-button>
        <el-button @click="clearData">
          <i class="fas fa-trash"></i> 清空
        </el-button>
        <el-button @click="isUppercase = !isUppercase" :type="isUppercase ? 'primary' : 'default'">
          {{ isUppercase ? '转小写' : '转大写' }}
        </el-button>
      </div>
      
      <div class="all-results-section" v-if="allHashResults.length > 0">
        <h3>所有哈希值计算结果</h3>
        <el-table :data="allHashResults" style="width: 100%" border>
          <el-table-column prop="algorithm" label="算法" width="120"></el-table-column>
          <el-table-column prop="hash" label="哈希值"></el-table-column>
        </el-table>
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
import { ref, inject } from 'vue'
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert, ElTable, ElTableColumn } from 'element-plus'
import CryptoJS from 'crypto-js'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 数据模型
const inputData = ref('')
const selectedAlgorithm = ref('MD5')
const hashResult = ref('')
const allHashResults = ref([])
const validationResult = ref(null)
const isUppercase = ref(true) // 默认大写
const salt = ref('') // 盐值

// 哈希算法列表
const algorithms = ref([
  { name: 'MD5', fullName: 'Message-Digest Algorithm 5' },
  { name: 'SHA1', fullName: 'Secure Hash Algorithm 1' },
  { name: 'SHA256', fullName: 'Secure Hash Algorithm 256' },
  { name: 'SHA512', fullName: 'Secure Hash Algorithm 512' }
])

// 选择算法
const selectAlgorithm = (algorithm) => {
  selectedAlgorithm.value = algorithm
}

// 转换哈希结果大小写
const convertCase = (hash) => {
  return isUppercase.value ? hash.toUpperCase() : hash.toLowerCase()
}

// 计算哈希值
const calculateHash = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    // 准备数据和盐值
    const dataToHash = salt.value ? inputData.value + salt.value : inputData.value
    
    let hashResultValue = ''
    
    switch (selectedAlgorithm.value) {
      case 'MD5':
        hashResultValue = CryptoJS.MD5(dataToHash).toString()
        break
      case 'SHA1':
        hashResultValue = CryptoJS.SHA1(dataToHash).toString()
        break
      case 'SHA256':
        hashResultValue = CryptoJS.SHA256(dataToHash).toString()
        break
      case 'SHA512':
        hashResultValue = CryptoJS.SHA512(dataToHash).toString()
        break
      default:
        throw new Error('不支持的哈希算法')
    }
    
    hashResult.value = convertCase(hashResultValue)
    allHashResults.value = []
    validationResult.value = {
      type: 'success',
      message: `${selectedAlgorithm.value}哈希值计算成功！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `哈希计算失败: ${error.message}`
    }
  }
}

// 计算所有哈希值
const calculateAllHashes = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    // 准备数据和盐值
    const dataToHash = salt.value ? inputData.value + salt.value : inputData.value
    
    const results = []
    
    // 计算MD5
    const md5Hash = CryptoJS.MD5(dataToHash).toString()
    results.push({ algorithm: 'MD5', hash: convertCase(md5Hash) })
    
    // 计算SHA1
    const sha1Hash = CryptoJS.SHA1(dataToHash).toString()
    results.push({ algorithm: 'SHA1', hash: convertCase(sha1Hash) })
    
    // 计算SHA256
    const sha256Hash = CryptoJS.SHA256(dataToHash).toString()
    results.push({ algorithm: 'SHA256', hash: convertCase(sha256Hash) })
    
    // 计算SHA512
    const sha512Hash = CryptoJS.SHA512(dataToHash).toString()
    results.push({ algorithm: 'SHA512', hash: convertCase(sha512Hash) })
    
    allHashResults.value = results
    hashResult.value = ''
    validationResult.value = {
      type: 'success',
      message: '所有哈希值计算成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `哈希计算失败: ${error.message}`
    }
  }
}


// 清空数据
const clearData = () => {
  inputData.value = ''
  hashResult.value = ''
  allHashResults.value = []
  salt.value = ''
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  let result = hashResult.value || (allHashResults.value.length > 0 ? 
    allHashResults.value.map(r => `${r.algorithm}: ${r.hash}`).join('\n') : '')
  
  // 如果需要转换大小写，则进行转换
  if (result && !hashResult.value && allHashResults.value.length > 0) {
    // 对于所有结果的复制，保持原有格式
  } else if (result && hashResult.value) {
    // 对于单个结果的复制，根据当前设置转换大小写
    result = convertCase(result)
  }
  
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

// 初始化示例数据
inputData.value = 'Hello, WeTools!'
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

.input-section {
  width: 100%;
}

.salt-section {
  width: 100%;
  margin-top: 1rem;
}

.hash-algorithms {
  width: 100%;
}

.hash-algorithms h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.algorithm-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.algorithm-buttons .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
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

.result-section,
.all-results-section {
  width: 100%;
}

.result-section h3,
.all-results-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .algorithm-buttons {
    flex-direction: column;
  }
  
  .algorithm-buttons .el-button {
    width: 100%;
    max-width: none;
  }
  
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