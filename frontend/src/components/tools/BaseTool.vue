<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入内容"
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="processData">
        <i class="fas fa-magic"></i> 处理
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-file-alt"></i> 输出结果</h3>
      <CodeBlock 
        :code="outputData" 
        :language="outputLanguage"
        :show-line-numbers="true"
        :auto-detect="true"
        :show-header="true"
      />
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>工具说明内容</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const inputData = ref('')
const outputData = ref('')
const outputLanguage = ref('auto')

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// 处理数据
const processData = () => {
  // 具体工具实现
  outputData.value = inputData.value
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
}

// 复制结果
const copyResult = async () => {
  try {
    await navigator.clipboard.writeText(outputData.value)
    copyButtonText.value = '已复制'
    setTimeout(() => {
      copyButtonText.value = '复制'
    }, 2000)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = outputData.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    
    copyButtonText.value = '已复制'
    setTimeout(() => {
      copyButtonText.value = '复制'
    }, 2000)
  }
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

textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  resize: vertical;
  transition: border-color 0.2s;
  min-height: 150px;
}

textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
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

.result {
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  border-radius: 4px;
  padding: 1rem;
  min-height: 150px;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  line-height: 1.5;
  overflow: auto;
  max-height: 400px;
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