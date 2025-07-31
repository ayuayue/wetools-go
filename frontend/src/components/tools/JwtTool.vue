<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>JWT 工具</h2>
      <p>在线JWT编码、解码和验证工具</p>
    </div>
    
    <div class="tool-content">
      <!-- 左右布局 -->
      <div class="two-column-layout">
        <!-- 左侧：编码部分 -->
        <div class="column encode-column">
          <div class="section">
            <div class="section-header">
              <h3>编码JWT</h3>
            </div>
            
            <!-- Header部分 -->
            <div class="input-section">
              <div class="editor-header">
                <h4>Header</h4>
                <div class="editor-actions">
                  <el-button type="success" @click="formatHeader" size="small">
                    <i class="fas fa-magic"></i> 格式化
                  </el-button>
                </div>
              </div>
              <AceEditor
                v-model="headerData"
                language="json"
                :theme="theme"
                :showHeader="false"
                :style="{ height: '100px' }"
              />
            </div>
            
            <!-- Payload部分 -->
            <div class="input-section">
              <div class="editor-header">
                <h4>Payload</h4>
                <div class="editor-actions">
                  <el-button type="success" @click="formatPayload" size="small">
                    <i class="fas fa-magic"></i> 格式化
                  </el-button>
                </div>
              </div>
              <AceEditor
                v-model="payloadData"
                language="json"
                :theme="theme"
                :showHeader="false"
                :style="{ height: '250px' }"
              />
            </div>
            
            <!-- Signature部分 -->
            <div class="input-section">
              <div class="editor-header">
                <h4>Signature</h4>
                <div class="editor-actions">
                  <el-select v-model="algorithm" placeholder="选择算法" size="small">
                    <el-option
                      v-for="alg in algorithms"
                      :key="alg.value"
                      :label="alg.label"
                      :value="alg.value"
                    />
                  </el-select>
                </div>
              </div>
              <el-input
                v-model="secretKey"
                type="password"
                placeholder="请输入密钥"
                show-password
              />
              <div class="signature-info">
                <p>HMACSHA256(</p>
                <p class="indent">base64UrlEncode(header) + "." +</p>
                <p class="indent">base64UrlEncode(payload),</p>
                <p class="indent">your-256-bit-secret</p>
                <p>)</p>
              </div>
            </div>
            
            <!-- JWT Token显示 -->
            <div class="input-section">
              <div class="editor-header">
                <h4>JWT Token</h4>
                <div class="editor-actions">
                  <el-button type="success" @click="copyToken" size="small">
                    <i class="fas fa-copy"></i> 复制Token
                  </el-button>
                  <el-button type="primary" @click="generateToken" size="small">
                    <i class="fas fa-sync"></i> 生成Token
                  </el-button>
                </div>
              </div>
              <AceEditor
                v-model="jwtToken"
                language="text"
                :theme="theme"
                :readonly="true"
                :showHeader="false"
                :showFormatButton="false"
                :style="{ height: '120px' }"
              />
            </div>
          </div>
        </div>
        
        <!-- 右侧：解码部分 -->
        <div class="column decode-column">
          <div class="section">
            <div class="section-header">
              <h3>解码JWT</h3>
            </div>
            
            <!-- Token输入 -->
            <div class="input-section">
              <div class="editor-header">
                <h4>输入Token</h4>
                <div class="editor-actions">
                  <el-button type="primary" @click="decodeToken" size="small">
                    <i class="fas fa-search"></i> 解码Token
                  </el-button>
                </div>
              </div>
              <AceEditor
                v-model="tokenToDecode"
                language="text"
                :theme="theme"
                :showHeader="false"
                :showFormatButton="false"
                :style="{ height: '150px' }"
              />
            </div>
            
            <!-- 解码结果显示 -->
            <div v-if="decodedToken" class="decoded-result">
              <div class="result-section">
                <div class="editor-header">
                  <h4>Header</h4>
                </div>
                <AceEditor
                  :model-value="decodedHeader"
                  language="json"
                  :theme="theme"
                  :readonly="true"
                  :show-header="false"
                  :show-format-button="false"
                  :show-fullscreen-button="false"
                  :style="{ height: '100px' }"
                />
              </div>
              <div class="result-section">
                <div class="editor-header">
                  <h4>Payload</h4>
                </div>
                <AceEditor
                  :model-value="decodedPayload"
                  language="json"
                  :theme="theme"
                  :readonly="true"
                  :show-header="false"
                  :show-format-button="false"
                  :show-fullscreen-button="false"
                  :style="{ height: '250px' }"
                />
              </div>
              <div class="result-section">
                <div class="editor-header">
                  <h4>Signature</h4>
                </div>
                <AceEditor
                  :model-value="decodedSignature"
                  language="text"
                  :theme="theme"
                  :readonly="true"
                  :show-header="false"
                  :show-format-button="false"
                  :show-fullscreen-button="false"
                  :style="{ height: '50px' }"
                />
              </div>
              <div v-if="verificationResult" class="verification-result">
                <el-alert
                  :type="verificationResult.isValid ? 'success' : 'error'"
                  :title="verificationResult.message"
                  show-icon
                  :closable="false"
                />
              </div>
            </div>
            
            <!-- 无解码结果显示提示 -->
            <div v-else class="no-decode-result">
              <el-empty description="请输入JWT Token进行解码" />
            </div>
          </div>
        </div>
      </div>
      
      <div class="button-group">
        <el-button @click="clearAll">
          <i class="fas fa-trash"></i> 清空
        </el-button>
      </div>
    </div>
  </el-card>
</template>

<script setup>
import { ref, inject } from 'vue'
import { ElCard, ElInput, ElButton, ElAlert, ElSelect, ElOption, ElMessage, ElEmpty } from 'element-plus'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 数据模型
const headerData = ref('{\n  "alg": "HS256",\n  "typ": "JWT"\n}')
const payloadData = ref('{\n  "sub": "1234567890",\n  "name": "John Doe",\n  "iat": 1516239022\n}')
const secretKey = ref('your-256-bit-secret')
const jwtToken = ref('')
const tokenToDecode = ref('')
const decodedToken = ref(false)
const decodedHeader = ref('')
const decodedPayload = ref('')
const decodedSignature = ref('')
const verificationResult = ref(null)

// 算法选项
const algorithm = ref('HS256')
const algorithms = [
  { label: 'HS256', value: 'HS256' },
  { label: 'HS384', value: 'HS384' },
  { label: 'HS512', value: 'HS512' },
  { label: 'RS256', value: 'RS256' },
  { label: 'RS384', value: 'RS384' },
  { label: 'RS512', value: 'RS512' }
]

// UTF-8解码函数
const utf8Decode = (str) => {
  try {
    // 先尝试使用 atob 解码 base64
    const binaryString = atob(str)
    // 然后正确解码 UTF-8
    const bytes = new Uint8Array(binaryString.length)
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i)
    }
    const decoder = new TextDecoder('utf-8')
    return decoder.decode(bytes)
  } catch (e) {
    // 如果上面的方法失败，尝试传统的解码方式
    try {
      return decodeURIComponent(escape(str))
    } catch (e2) {
      return str
    }
  }
}

// 格式化Header
const formatHeader = () => {
  try {
    const parsed = JSON.parse(headerData.value)
    headerData.value = JSON.stringify(parsed, null, 2)
  } catch (error) {
    ElMessage.error('Header格式错误: ' + error.message)
  }
}

// 格式化Payload
const formatPayload = () => {
  try {
    const parsed = JSON.parse(payloadData.value)
    payloadData.value = JSON.stringify(parsed, null, 2)
  } catch (error) {
    ElMessage.error('Payload格式错误: ' + error.message)
  }
}

// 生成JWT Token
const generateToken = () => {
  // 这里应该实现JWT生成逻辑
  // 为了简化，我们只是模拟生成一个token
  try {
    const header = btoa(headerData.value).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
    const payload = btoa(payloadData.value).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
    const signature = 'signature' // 实际应用中需要根据算法和密钥生成签名
    
    jwtToken.value = `${header}.${payload}.${signature}`
  } catch (error) {
    ElMessage.error('生成Token失败: ' + error.message)
  }
}

// 复制Token
const copyToken = async () => {
  if (!jwtToken.value) {
    ElMessage.warning('没有Token可复制')
    return
  }
  
  try {
    await navigator.clipboard.writeText(jwtToken.value)
    ElMessage.success('Token已复制到剪贴板！')
  } catch (error) {
    ElMessage.error('复制失败: ' + error.message)
  }
}

// 解码Token
const decodeToken = () => {
  // 重置之前的结果
  decodedToken.value = false
  decodedHeader.value = ''
  decodedPayload.value = ''
  decodedSignature.value = ''
  verificationResult.value = null
  
  if (!tokenToDecode.value) {
    ElMessage.warning('请输入要解码的Token')
    return
  }
  
  try {
    const parts = tokenToDecode.value.split('.')
    if (parts.length !== 3) {
      throw new Error('无效的JWT Token格式')
    }
    
    const [header, payload, signature] = parts
    
    // 解码Header
    try {
      const base64Header = header.replace(/-/g, '+').replace(/_/g, '/')
      // 补全base64字符串长度
      const paddedHeader = base64Header.padEnd(base64Header.length + (4 - base64Header.length % 4) % 4, '=')
      const decodedHeaderStr = utf8Decode(paddedHeader)
      decodedHeader.value = JSON.stringify(JSON.parse(decodedHeaderStr), null, 2)
    } catch (e) {
      throw new Error('Header解码失败: ' + e.message)
    }
    
    // 解码Payload
    try {
      const base64Payload = payload.replace(/-/g, '+').replace(/_/g, '/')
      // 补全base64字符串长度
      const paddedPayload = base64Payload.padEnd(base64Payload.length + (4 - base64Payload.length % 4) % 4, '=')
      const decodedPayloadStr = utf8Decode(paddedPayload)
      decodedPayload.value = JSON.stringify(JSON.parse(decodedPayloadStr), null, 2)
    } catch (e) {
      throw new Error('Payload解码失败: ' + e.message)
    }
    
    // Signature
    decodedSignature.value = signature
    
    decodedToken.value = true
    verificationResult.value = {
      isValid: true,
      message: 'Token解码成功'
    }
  } catch (error) {
    decodedToken.value = false
    verificationResult.value = {
      isValid: false,
      message: 'Token解码失败: ' + error.message
    }
    ElMessage.error('Token解码失败: ' + error.message)
  }
}

// 清空所有数据
const clearAll = () => {
  headerData.value = '{\n  "alg": "HS256",\n  "typ": "JWT"\n}'
  payloadData.value = '{\n  "sub": "1234567890",\n  "name": "John Doe",\n  "iat": 1516239022\n}'
  secretKey.value = 'your-256-bit-secret'
  jwtToken.value = ''
  tokenToDecode.value = ''
  decodedToken.value = false
  decodedHeader.value = ''
  decodedPayload.value = ''
  decodedSignature.value = ''
  verificationResult.value = null
}

// 初始化示例数据
headerData.value = '{\n  "alg": "HS256",\n  "typ": "JWT"\n}'
payloadData.value = '{\n  "sub": "1234567890",\n  "name": "张三",\n  "iat": 1516239022\n}'
tokenToDecode.value = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IuW8oOS4iSIsImlhdCI6MTUxNjIzOTAyMn0.RTJwB1E34VJ3LtaBIV6OWQO4szBY9O5Ywr7F3R8G2do'

// 初始化时自动解码示例Token
setTimeout(() => {
  decodeToken()
}, 0)
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
}

.tool-header {
  margin-bottom: 1rem;
  text-align: left;
}

.tool-header h2 {
  margin: 0 0 0.25rem 0;
  color: #333;
  font-size: 1.5rem;
}

.tool-header p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
  text-align: left;
}

.tool-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 左右布局 */
.two-column-layout {
  display: flex;
  gap: 1.5rem;
}

.column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #eee;
  text-align: left;
}

.input-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.editor-header h4 {
  margin: 0;
  color: #333;
  font-size: 1rem;
  text-align: left;
}

.editor-actions {
  display: flex;
  gap: 0.5rem;
}

.signature-info {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 0.5rem;
  font-family: monospace;
  font-size: 0.8rem;
  color: #666;
  margin-top: 0.5rem;
  text-align: left;
}

.signature-info .indent {
  margin-left: 1rem;
}

.result-section {
  margin-bottom: 1rem;
}

.result-section h4 {
  margin: 0 0 0.5rem 0;
  color: #333;
  text-align: left;
}

.verification-result {
  margin-top: 1rem;
}

.no-decode-result {
  text-align: center;
  padding: 2rem 0;
}

.button-group {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.button-group .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
}

/* 优化编辑器容器样式 */
.ace-editor-container {
  border-radius: 4px;
  overflow: hidden;
  background: #f8fafc;
  margin: 10px 0;
  border: 1px solid #e2e8f0;
  text-align: left;
}

.ace-editor-container :deep(.ace_content) {
  text-align: left;
}

/* 优化文本对齐 */
:deep(.el-textarea__inner) {
  text-align: left;
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
  
  .tool-header,
  .tool-header h2,
  .tool-header p,
  .section-header h3,
  .editor-header h4,
  .result-section h4 {
    text-align: center;
  }
}
</style>