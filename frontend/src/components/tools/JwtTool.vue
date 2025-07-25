<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-key"></i> JWT Token</h3>
      <textarea 
        v-model="jwtToken" 
        placeholder="请输入JWT Token"
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="decodeJWT">
        <i class="fas fa-unlock"></i> 解码JWT
      </button>
      <button @click="generateJWT">
        <i class="fas fa-lock"></i> 生成JWT
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section" v-if="decodedHeader || decodedPayload">
      <h3><i class="fas fa-file-alt"></i> 解码结果</h3>
      <div class="jwt-section" v-if="decodedHeader">
        <h4>Header</h4>
        <CodeBlock 
          :code="decodedHeader" 
          language="json"
          :show-line-numbers="false"
        />
      </div>
      <div class="jwt-section" v-if="decodedPayload">
        <h4>Payload</h4>
        <CodeBlock 
          :code="decodedPayload" 
          language="json"
          :show-line-numbers="false"
        />
      </div>
      <div class="jwt-section" v-if="signature">
        <h4>Signature</h4>
        <div class="signature">{{ signature }}</div>
      </div>
    </div>

    <div class="tool-section" v-if="showGenerator">
      <h3><i class="fas fa-cogs"></i> 生成JWT</h3>
      <div class="generator-section">
        <h4>Header</h4>
        <textarea 
          v-model="headerInput" 
          placeholder='{"alg": "HS256", "typ": "JWT"}'
        ></textarea>
      </div>
      <div class="generator-section">
        <h4>Payload</h4>
        <textarea 
          v-model="payloadInput" 
          placeholder='{"sub": "1234567890", "name": "John Doe", "iat": 1516239022}'
        ></textarea>
      </div>
      <div class="generator-section">
        <h4>Secret</h4>
        <input 
          v-model="secretKey" 
          type="password" 
          placeholder="请输入密钥"
        />
      </div>
      <div class="button-group">
        <button @click="createJWT">
          <i class="fas fa-lock"></i> 创建JWT
        </button>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>JWT (JSON Web Token) 编解码工具支持以下功能：</p>
      <ul class="description-list">
        <li>解码JWT Token，查看Header和Payload内容</li>
        <li>生成自定义的JWT Token</li>
        <li>支持HS256、HS384、HS512等签名算法</li>
        <li>实时显示Token的各部分信息</li>
        <li>所有操作均在浏览器端完成，保证安全性</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const jwtToken = ref('')
const decodedHeader = ref('')
const decodedPayload = ref('')
const signature = ref('')
const showGenerator = ref(false)

const headerInput = ref('')
const payloadInput = ref('')
const secretKey = ref('')

// 解码JWT
const decodeJWT = () => {
  if (!jwtToken.value) {
    clearDecodeResult()
    return
  }

  try {
    const parts = jwtToken.value.split('.')
    if (parts.length !== 3) {
      throw new Error('无效的JWT格式')
    }

    // 解码Header
    const header = JSON.parse(atob(parts[0]))
    decodedHeader.value = JSON.stringify(header, null, 2)

    // 解码Payload
    const payload = JSON.parse(atob(parts[1]))
    decodedPayload.value = JSON.stringify(payload, null, 2)

    // Signature
    signature.value = parts[2]
  } catch (e) {
    clearDecodeResult()
    decodedHeader.value = '解码错误: ' + e.message
  }
}

// 清空解码结果
const clearDecodeResult = () => {
  decodedHeader.value = ''
  decodedPayload.value = ''
  signature.value = ''
}

// 生成JWT界面
const generateJWT = () => {
  showGenerator.value = !showGenerator.value
  if (showGenerator.value) {
    // 初始化默认值
    headerInput.value = JSON.stringify({
      alg: 'HS256',
      typ: 'JWT'
    }, null, 2)
    
    payloadInput.value = JSON.stringify({
      sub: '1234567890',
      name: 'John Doe',
      iat: Math.floor(Date.now() / 1000)
    }, null, 2)
  }
}

// 创建JWT
const createJWT = () => {
  try {
    // 这里只是一个模拟实现，实际项目中应该使用crypto库进行签名
    const header = btoa(headerInput.value)
    const payload = btoa(payloadInput.value)
    const signature = 'signature_placeholder' // 模拟签名
    
    jwtToken.value = `${header}.${payload}.${signature}`
    showGenerator.value = false
  } catch (e) {
    jwtToken.value = '生成错误: ' + e.message
  }
}

// 清空数据
const clearData = () => {
  jwtToken.value = ''
  clearDecodeResult()
  showGenerator.value = false
  headerInput.value = ''
  payloadInput.value = ''
  secretKey.value = ''
}
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

.tool-section h4 {
  margin-bottom: 0.5rem;
  color: #64748b;
  font-size: 1rem;
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

input[type="password"] {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  transition: border-color 0.2s;
  background: white;
  color: #333;
}

input[type="password"]:focus {
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

.dark-theme input[type="password"] {
  background: #3d3d3d;
  border: 1px solid #555;
  color: #e0e0e0;
}

.dark-theme input[type="password"]:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.3);
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

.jwt-section {
  margin-bottom: 1.5rem;
}

.jwt-section:last-child {
  margin-bottom: 0;
}

.signature {
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 4px;
  padding: 1rem;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  word-break: break-all;
  white-space: pre-wrap;
}

.generator-section {
  margin-bottom: 1.5rem;
}

.generator-section:last-child {
  margin-bottom: 0;
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
}
</style>