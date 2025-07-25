<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 输入文本</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入要处理的文本"
      ></textarea>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-sliders-h"></i> 文本处理选项</h3>
      <div class="config-grid">
        <div class="config-group">
          <label>
            <input 
              v-model="config.toUpperCase" 
              type="checkbox" 
            /> 转大写
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.toLowerCase" 
              type="checkbox" 
            /> 转小写
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.removeSpaces" 
              type="checkbox" 
            /> 去除空格
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.removeLineBreaks" 
              type="checkbox" 
            /> 去除换行
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.removeInvisibleChars" 
              type="checkbox" 
            /> 去除不可见字符
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.removeLineNumbers" 
              type="checkbox" 
            /> 去除行号
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.deduplicateLines" 
              type="checkbox" 
            /> 行去重
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.chineseToEnglishPunctuation" 
              type="checkbox" 
            /> 中文标点转英文
          </label>
        </div>
        <div class="config-group">
          <label>
            <input 
              v-model="config.englishToChinesePunctuation" 
              type="checkbox" 
            /> 英文标点转中文
          </label>
        </div>
      </div>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-magic"></i> 高级处理</h3>
      <div class="advanced-group">
        <div class="regex-group">
          <label>正则替换:</label>
          <div class="regex-inputs">
            <input 
              v-model="regexPattern" 
              placeholder="正则表达式" 
            />
            <input 
              v-model="regexReplacement" 
              placeholder="替换内容" 
            />
            <button @click="applyRegexReplace">
              <i class="fas fa-exchange-alt"></i> 替换
            </button>
          </div>
        </div>
        <div class="save-group">
          <button @click="saveToFile">
            <i class="fas fa-save"></i> 保存到文件
          </button>
        </div>
      </div>
    </div>

    <div class="button-group">
      <button @click="processText">
        <i class="fas fa-magic"></i> 处理文本
      </button>
      <button @click="countStats">
        <i class="fas fa-chart-bar"></i> 统计信息
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <div class="result-header">
        <h3><i class="fas fa-file-alt"></i> 处理结果</h3>
        <div class="stats" v-if="stats">
          字数: {{ stats.characters }} | 行数: {{ stats.lines }} | 单词数: {{ stats.words }}
        </div>
      </div>
      <CodeBlock 
        :code="outputData" 
        language="text"
        :show-line-numbers="true"
      />
      <div class="result-footer">
        <button class="copy-btn" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </button>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>文本处理工具支持多种文本处理功能：</p>
      <ul class="description-list">
        <li>大小写转换、去除空格和换行符</li>
        <li>去除不可见字符和行号</li>
        <li>文本去重和标点符号转换</li>
        <li>正则表达式替换功能</li>
        <li>文本统计信息（字数、行数、单词数）</li>
        <li>支持保存处理结果到文件</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import CodeBlock from '../CodeBlock.vue'

const inputData = ref('')
const outputData = ref('')
const stats = ref(null)

const config = reactive({
  toUpperCase: false,
  toLowerCase: false,
  removeSpaces: false,
  removeLineBreaks: false,
  removeInvisibleChars: false,
  removeLineNumbers: false,
  deduplicateLines: false,
  chineseToEnglishPunctuation: false,
  englishToChinesePunctuation: false
})

const regexPattern = ref('')
const regexReplacement = ref('')

// 处理文本
const processText = () => {
  let result = inputData.value
  
  // 应用各种处理选项
  if (config.toUpperCase) {
    result = result.toUpperCase()
  }
  
  if (config.toLowerCase) {
    result = result.toLowerCase()
  }
  
  if (config.removeSpaces) {
    result = result.replace(/\s+/g, '')
  }
  
  if (config.removeLineBreaks) {
    result = result.replace(/\n/g, '')
  }
  
  if (config.removeInvisibleChars) {
    result = result.replace(/[\u200B-\u200D\uFEFF]/g, '')
  }
  
  if (config.removeLineNumbers) {
    result = result.split('\n').map(line => {
      return line.replace(/^\s*\d+[\.\-\)\s]*/, '')
    }).join('\n')
  }
  
  if (config.deduplicateLines) {
    const lines = result.split('\n')
    const uniqueLines = [...new Set(lines)]
    result = uniqueLines.join('\n')
  }
  
  if (config.chineseToEnglishPunctuation) {
    const punctuationMap = {
      '，': ',',
      '。': '.',
      '！': '!',
      '？': '?',
      '；': ';',
      '：': ':',
      '“': '"',
      '”': '"',
      '‘': "'",
      '’': "'",
      '（': '(',
      '）': ')',
      '【': '[',
      '】': ']',
      '《': '<',
      '》': '>'
    }
    for (const [chinese, english] of Object.entries(punctuationMap)) {
      result = result.replace(new RegExp(chinese, 'g'), english)
    }
  }
  
  if (config.englishToChinesePunctuation) {
    const punctuationMap = {
      ',': '，',
      '.': '。',
      '!': '！',
      '?': '？',
      ';': '；',
      ':': '：',
      '"': '“',
      "'": '‘',
      '(': '（',
      ')': '）',
      '[': '【',
      ']': '】',
      '<': '《',
      '>': '》'
    }
    for (const [english, chinese] of Object.entries(punctuationMap)) {
      result = result.replace(new RegExp('\\' + english, 'g'), chinese)
    }
  }
  
  outputData.value = result
  countStats(result)
}

// 正则替换
const applyRegexReplace = () => {
  try {
    const regex = new RegExp(regexPattern.value, 'g')
    outputData.value = outputData.value.replace(regex, regexReplacement.value)
    countStats(outputData.value)
  } catch (e) {
    outputData.value = '正则表达式错误: ' + e.message
  }
}

// 统计信息
const countStats = (text = outputData.value) => {
  if (!text) return
  
  stats.value = {
    characters: text.length,
    lines: text.split('\n').length,
    words: text.trim().split(/\s+/).filter(word => word.length > 0).length
  }
}

// 保存到文件
const saveToFile = () => {
  if (!outputData.value) return
  
  const blob = new Blob([outputData.value], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'text_result.txt'
  link.click()
  URL.revokeObjectURL(url)
}

// 复制结果
const copyResult = async () => {
  try {
    await navigator.clipboard.writeText(outputData.value)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = outputData.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
  stats.value = null
  
  // 重置配置
  Object.keys(config).forEach(key => {
    config[key] = false
  })
  
  regexPattern.value = ''
  regexReplacement.value = ''
}

// 初始化示例数据
const init = () => {
  inputData.value = `1. 这是第一行文本，包含中文标点符号！
2. This is the second line, with English punctuation!
3. 第三行文本用于测试去重功能。
4. This is the fourth line for testing deduplication.
5. 重复的行用于测试去重功能。
6. Duplicate line for testing deduplication.`
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

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.result-header h3 {
  margin-bottom: 0;
}

.stats {
  font-size: 0.9rem;
  color: #64748b;
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

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.config-group label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.config-group input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.advanced-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.regex-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.regex-group label {
  font-weight: 500;
  color: #64748b;
}

.regex-inputs {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.regex-inputs input {
  flex: 1;
  min-width: 150px;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
}

.regex-inputs button {
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

.regex-inputs button:hover {
  background: #5a6fd8;
}

.save-group button {
  padding: 0.5rem 1rem;
  background: #4ade80;
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

.save-group button:hover {
  background: #22c55e;
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
  .button-group,
  .regex-inputs {
    flex-direction: column;
  }
  
  button,
  .regex-inputs input,
  .regex-inputs button {
    width: 100%;
    justify-content: center;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
  }
  
  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
</style>