<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>哈希计算工具</h2>
      <p>MD5、SHA1、SHA256等哈希算法计算工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="8"
              placeholder="请输入要计算哈希值的文本"
              resize="vertical"
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
        <el-button type="success" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </el-button>
      </div>
      
      <div class="result-section" v-if="hashResult">
        <h3>计算结果</h3>
        <el-input
          v-model="hashResult"
          readonly
          placeholder="哈希计算结果将显示在这里"
        />
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
import { ref } from 'vue'
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert, ElTable, ElTableColumn } from 'element-plus'

// 数据模型
const inputData = ref('')
const selectedAlgorithm = ref('MD5')
const hashResult = ref('')
const allHashResults = ref([])
const validationResult = ref(null)

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

// 计算哈希值
const calculateHash = async () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    // 使用Web Crypto API计算哈希值
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    let hashBuffer
    
    switch (selectedAlgorithm.value) {
      case 'MD5':
        // MD5需要使用第三方库或自定义实现，这里简化处理
        hashResult.value = await calculateMD5(inputData.value)
        break
      case 'SHA1':
        hashBuffer = await crypto.subtle.digest('SHA-1', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      case 'SHA256':
        hashBuffer = await crypto.subtle.digest('SHA-256', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      case 'SHA512':
        hashBuffer = await crypto.subtle.digest('SHA-512', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      default:
        throw new Error('不支持的哈希算法')
    }
    
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
const calculateAllHashes = async () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    const results = []
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    
    // 计算SHA1
    const sha1Buffer = await crypto.subtle.digest('SHA-1', data)
    const sha1Hash = Array.from(new Uint8Array(sha1Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA1', hash: sha1Hash })
    
    // 计算SHA256
    const sha256Buffer = await crypto.subtle.digest('SHA-256', data)
    const sha256Hash = Array.from(new Uint8Array(sha256Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA256', hash: sha256Hash })
    
    // 计算SHA512
    const sha512Buffer = await crypto.subtle.digest('SHA-512', data)
    const sha512Hash = Array.from(new Uint8Array(sha512Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA512', hash: sha512Hash })
    
    // 计算MD5 (简化实现)
    const md5Hash = await calculateMD5(inputData.value)
    results.push({ algorithm: 'MD5', hash: md5Hash })
    
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

// 简化的MD5计算函数
const calculateMD5 = async (str) => {
  // 使用crypto.subtle.digest计算MD5
  const encoder = new TextEncoder()
  const data = encoder.encode(str)
  const hashBuffer = await crypto.subtle.digest('MD5', data)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('').toUpperCase()
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  hashResult.value = ''
  allHashResults.value = []
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  const result = hashResult.value || (allHashResults.value.length > 0 ? 
    allHashResults.value.map(r => `${r.algorithm}: ${r.hash}`).join('\n') : '')
  
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