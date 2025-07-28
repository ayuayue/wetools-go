<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>加密/解密工具</h2>
      <p>AES、DES等对称加密算法工具</p>
    </div>
    
    <div class="tool-content">
      <el-tabs v-model="activeTab" class="tool-tabs">
        <el-tab-pane label="AES加密" name="aes">
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
                  v-model="aesInput"
                  language="text"
                  :theme="theme"
                  :showHeader="false"
                />
              </div>
              
              <el-row :gutter="20">
                <el-col :span="12">
                  <div class="key-section">
                    <el-input
                      v-model="aesKey"
                      placeholder="请输入加密密钥"
                      show-password
                    />
                  </div>
                </el-col>
                <el-col :span="12">
                  <div class="iv-section">
                    <el-input
                      v-model="aesIv"
                      placeholder="请输入初始化向量(IV)"
                    />
                  </div>
                </el-col>
              </el-row>
            </div>
            
            <!-- 右侧输出区域 -->
            <div class="column output-column">
              <div class="output-section">
                <div class="editor-header">
                  <h3>加密输出</h3>
                  <div class="editor-actions">
                    <el-button type="success" @click="copyAESResult">
                      <i class="fas fa-copy"></i> 复制结果
                    </el-button>
                  </div>
                </div>
                <AceEditor
                  v-model="aesOutput"
                  language="text"
                  :theme="theme"
                  :readonly="true"
                  :showHeader="false"
                />
              </div>
            </div>
          </div>
          
          <div class="button-group">
            <el-button type="primary" @click="encryptAES">
              <i class="fas fa-lock"></i> AES加密
            </el-button>
            <el-button @click="decryptAES">
              <i class="fas fa-unlock"></i> AES解密
            </el-button>
            <el-button @click="generateAESKey">
              <i class="fas fa-key"></i> 生成密钥
            </el-button>
            <el-button @click="clearAES">
              <i class="fas fa-trash"></i> 清空
            </el-button>
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
import { ref, inject } from 'vue'
import { ElCard, ElTabs, ElTabPane, ElRow, ElCol, ElInput, ElButton, ElAlert } from 'element-plus'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 数据模型
const activeTab = ref('aes')
const aesInput = ref('')
const aesKey = ref('')
const aesIv = ref('')
const aesOutput = ref('')
const validationResult = ref(null)

// AES加密
const encryptAES = async () => {
  try {
    if (!aesInput.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要加密的文本'
      }
      return
    }
    
    if (!aesKey.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入加密密钥'
      }
      return
    }
    
    // 使用Web Crypto API进行AES加密
    const encoder = new TextEncoder()
    const data = encoder.encode(aesInput.value)
    
    // 生成密钥
    const keyData = encoder.encode(aesKey.value.padEnd(32, '0').substring(0, 32))
    const key = await crypto.subtle.importKey(
      'raw',
      keyData,
      { name: 'AES-CBC' },
      false,
      ['encrypt']
    )
    
    // 生成IV
    let iv
    if (aesIv.value.trim()) {
      const ivData = encoder.encode(aesIv.value.padEnd(16, '0').substring(0, 16))
      iv = ivData
    } else {
      iv = new Uint8Array(16)
      crypto.getRandomValues(iv)
    }
    
    // 执行加密
    const encrypted = await crypto.subtle.encrypt(
      { name: 'AES-CBC', iv: iv },
      key,
      data
    )
    
    // 将结果转换为Base64
    const encryptedArray = new Uint8Array(encrypted)
    aesOutput.value = btoa(String.fromCharCode(...iv, ...encryptedArray))
    
    validationResult.value = {
      type: 'success',
      message: 'AES加密成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `AES加密失败: ${error.message}`
    }
  }
}

// AES解密
const decryptAES = async () => {
  try {
    if (!aesInput.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要解密的文本'
      }
      return
    }
    
    if (!aesKey.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入解密密钥'
      }
      return
    }
    
    // 使用Web Crypto API进行AES解密
    const encoder = new TextEncoder()
    
    // 从Base64解码
    const encryptedData = atob(aesInput.value)
    const encryptedArray = new Uint8Array(encryptedData.length)
    for (let i = 0; i < encryptedData.length; i++) {
      encryptedArray[i] = encryptedData.charCodeAt(i)
    }
    
    // 提取IV和加密数据
    const iv = encryptedArray.slice(0, 16)
    const data = encryptedArray.slice(16)
    
    // 生成密钥
    const keyData = encoder.encode(aesKey.value.padEnd(32, '0').substring(0, 32))
    const key = await crypto.subtle.importKey(
      'raw',
      keyData,
      { name: 'AES-CBC' },
      false,
      ['decrypt']
    )
    
    // 执行解密
    const decrypted = await crypto.subtle.decrypt(
      { name: 'AES-CBC', iv: iv },
      key,
      data
    )
    
    // 将结果转换为文本
    const decoder = new TextDecoder()
    aesOutput.value = decoder.decode(decrypted)
    
    validationResult.value = {
      type: 'success',
      message: 'AES解密成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `AES解密失败: ${error.message}`
    }
  }
}

// 生成AES密钥
const generateAESKey = () => {
  try {
    // 生成随机密钥
    const keyLength = 32 // 256位
    const array = new Uint8Array(keyLength)
    crypto.getRandomValues(array)
    aesKey.value = Array.from(array)
      .map(b => b.toString(16).padStart(2, '0'))
      .join('')
    validationResult.value = {
      type: 'success',
      message: 'AES密钥生成成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `密钥生成失败: ${error.message}`
    }
  }
}

// 清空AES数据
const clearAES = () => {
  aesInput.value = ''
  aesKey.value = ''
  aesIv.value = ''
  aesOutput.value = ''
  validationResult.value = null
}


// 复制AES结果
const copyAESResult = async () => {
  if (!aesOutput.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(aesOutput.value)
    validationResult.value = {
      type: 'success',
      message: 'AES结果已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}


// 初始化示例数据
aesInput.value = 'Hello, WeTools!'
aesKey.value = 'my-secret-key-12345'
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
.key-section,
.iv-section {
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
}
</style>