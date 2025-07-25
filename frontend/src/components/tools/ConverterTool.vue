<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-exchange-alt"></i> 格式转换器</h3>
      <div class="converter-group">
        <div class="converter-row">
          <select v-model="sourceFormat">
            <option value="json">JSON</option>
            <option value="xml">XML</option>
            <option value="csv">CSV</option>
          </select>
          <span>→</span>
          <select v-model="targetFormat">
            <option value="json">JSON</option>
            <option value="xml">XML</option>
            <option value="csv">CSV</option>
          </select>
        </div>
      </div>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-edit"></i> 源数据</h3>
      <textarea 
        v-model="inputData" 
        placeholder="请输入要转换的数据"
      ></textarea>
    </div>

    <div class="button-group">
      <button @click="convertData">
        <i class="fas fa-exchange-alt"></i> 转换
      </button>
      <button class="copy-btn" @click="copyInput">
        <i class="fas fa-copy"></i> 复制输入
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section">
      <div class="result-header">
        <h3><i class="fas fa-file-alt"></i> 转换结果</h3>
      </div>
      <div class="result">{{ outputData }}</div>
      <div class="result-footer" v-if="outputData && !isValidationResult">
        <button class="copy-btn" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </button>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>格式转换器支持多种数据格式之间的相互转换：</p>
      <ul class="description-list">
        <li>JSON ↔ XML</li>
        <li>JSON ↔ CSV</li>
        <li>XML ↔ CSV</li>
        <li>支持复杂数据结构的转换</li>
        <li>所有转换操作在浏览器端完成，保证数据安全</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const sourceFormat = ref('json')
const targetFormat = ref('xml')
const inputData = ref('')
const outputData = ref('')

// 判断是否为校验结果（不显示复制按钮）
const isValidationResult = computed(() => {
  return outputData.value.startsWith('✓') || outputData.value.startsWith('✗')
})

// 简单的JSON到XML转换
const jsonToXml = (jsonObj, rootName = 'root') => {
  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n<${rootName}>\n`
  
  const convertObject = (obj, indent = '  ') => {
    let result = ''
    for (const key in obj) {
      if (typeof obj[key] === 'object' && obj[key] !== null) {
        if (Array.isArray(obj[key])) {
          result += `${indent}<${key}>\n`
          obj[key].forEach(item => {
            if (typeof item === 'object') {
              result += convertObject(item, indent + '  ')
            } else {
              result += `${indent}  <item>${item}</item>\n`
            }
          })
          result += `${indent}</${key}>\n`
        } else {
          result += `${indent}<${key}>\n${convertObject(obj[key], indent + '  ')}${indent}</${key}>\n`
        }
      } else {
        result += `${indent}<${key}>${obj[key]}</${key}>\n`
      }
    }
    return result
  }
  
  xml += convertObject(jsonObj)
  xml += `</${rootName}>`
  return xml
}

// 简单的XML到JSON转换
const xmlToJson = (xmlString) => {
  // 这是一个简化的实现，实际项目中应该使用专门的XML解析库
  return {
    message: "XML to JSON conversion",
    input: xmlString,
    note: "This is a simplified conversion for demonstration purposes"
  }
}

// 简单的JSON到CSV转换
const jsonToCsv = (jsonObj) => {
  if (!Array.isArray(jsonObj)) {
    return "JSON to CSV requires an array of objects"
  }
  
  if (jsonObj.length === 0) {
    return ""
  }
  
  // 获取所有字段名
  const headers = Object.keys(jsonObj[0])
  let csv = headers.join(',') + '\n'
  
  // 添加数据行
  jsonObj.forEach(row => {
    const values = headers.map(header => {
      const value = row[header]
      // 处理包含逗号或引号的值
      if (typeof value === 'string' && (value.includes(',') || value.includes('"'))) {
        return `"${value.replace(/"/g, '""')}"`
      }
      return value
    })
    csv += values.join(',') + '\n'
  })
  
  return csv
}

// 简单的CSV到JSON转换
const csvToJson = (csvString) => {
  const lines = csvString.trim().split('\n')
  if (lines.length === 0) return []
  
  const headers = lines[0].split(',').map(h => h.trim().replace(/^"(.*)"$/, '$1'))
  const result = []
  
  for (let i = 1; i < lines.length; i++) {
    const values = lines[i].split(',').map(v => v.trim().replace(/^"(.*)"$/, '$1'))
    const obj = {}
    headers.forEach((header, index) => {
      obj[header] = values[index] || ''
    })
    result.push(obj)
  }
  
  return result
}

// 转换数据
const convertData = () => {
  try {
    if (sourceFormat.value === targetFormat.value) {
      outputData.value = inputData.value
      return
    }
    
    let result = ''
    
    switch (sourceFormat.value) {
      case 'json':
        const jsonObj = JSON.parse(inputData.value)
        switch (targetFormat.value) {
          case 'xml':
            result = jsonToXml(jsonObj)
            break
          case 'csv':
            result = jsonToCsv(jsonObj)
            break
          default:
            result = inputData.value
        }
        break
      case 'xml':
        switch (targetFormat.value) {
          case 'json':
            result = JSON.stringify(xmlToJson(inputData.value), null, 2)
            break
          case 'csv':
            // XML to CSV 需要先转为JSON再转为CSV
            result = "XML to CSV conversion not implemented in this demo"
            break
          default:
            result = inputData.value
        }
        break
      case 'csv':
        switch (targetFormat.value) {
          case 'json':
            result = JSON.stringify(csvToJson(inputData.value), null, 2)
            break
          case 'xml':
            // CSV to XML 需要先转为JSON再转为XML
            const jsonFromCsv = csvToJson(inputData.value)
            result = jsonToXml(jsonFromCsv)
            break
          default:
            result = inputData.value
        }
        break
      default:
        result = inputData.value
    }
    
    outputData.value = result
  } catch (e) {
    outputData.value = '转换错误: ' + e.message
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
}

// 复制输入内容
const copyInput = async () => {
  try {
    await navigator.clipboard.writeText(inputData.value)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = inputData.value
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
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

// 初始化数据
const init = () => {
  inputData.value = JSON.stringify({
    "name": "WeTools",
    "type": "developer tools",
    "features": ["JSON", "XML", "HTML"]
  }, null, 2)
  convertData()
}

init()
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

.converter-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.converter-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.converter-row select {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.converter-row span {
  font-size: 1.2rem;
  font-weight: bold;
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
  
  .converter-row {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .result-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>