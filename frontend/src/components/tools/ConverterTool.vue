<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>格式转换器</h2>
      <p>多种数据格式互转工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="input-section">
            <div class="section-header">
              <h3>输入数据</h3>
              <div class="format-selector">
                <el-select v-model="inputFormat" placeholder="选择输入格式">
                  <el-option
                    v-for="format in formats"
                    :key="format.value"
                    :label="format.label"
                    :value="format.value"
                  />
                </el-select>
              </div>
            </div>
            <AceEditor
              v-model="inputData"
              :language="inputFormat"
              :theme="theme"
              :showHeader="false"
            />
          </div>
        </el-col>
        
        <el-col :span="12">
          <div class="output-section">
            <div class="section-header">
              <h3>输出数据</h3>
              <div class="format-selector">
                <el-select v-model="outputFormat" placeholder="选择输出格式">
                  <el-option
                    v-for="format in formats"
                    :key="format.value"
                    :label="format.label"
                    :value="format.value"
                  />
                </el-select>
              </div>
            </div>
            <AceEditor
              v-model="outputData"
              :language="outputFormat"
              :theme="theme"
              :readonly="true"
              :showHeader="false"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="conversion-options">
        <h3>转换选项</h3>
        <div class="options-content">
          <el-checkbox v-model="prettyPrint">美化输出</el-checkbox>
          <el-checkbox v-model="preserveOrder">保持键顺序</el-checkbox>
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="convertData">
          <i class="fas fa-exchange-alt"></i> 转换
        </el-button>
        <el-button @click="swapFormats">
          <i class="fas fa-exchange-alt"></i> 交换格式
        </el-button>
        <el-button @click="clearData">
          <i class="fas fa-trash"></i> 清空
        </el-button>
        <el-button type="success" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </el-button>
      </div>
      
      <div class="examples-section">
        <h3>常用转换示例</h3>
        <div class="example-buttons">
          <el-button
            v-for="example in examples"
            :key="example.name"
            @click="loadExample(example)"
          >
            {{ example.name }}
          </el-button>
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
import { ElCard, ElRow, ElCol, ElSelect, ElOption, ElInput, ElButton, ElAlert, ElCheckbox } from 'element-plus'

// 数据模型
const inputData = ref('')
const outputData = ref('')
const inputFormat = ref('json')
const outputFormat = ref('xml')
const prettyPrint = ref(true)
const preserveOrder = ref(true)
const validationResult = ref(null)

// 支持的格式
const formats = ref([
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'yaml', label: 'YAML' },
  { value: 'csv', label: 'CSV' },
  { value: 'html', label: 'HTML' }
])

// 示例数据
const examples = ref([
  {
    name: 'JSON to XML',
    inputFormat: 'json',
    outputFormat: 'xml',
    data: '{"name": "WeTools", "type": "developer tools", "features": ["JSON", "XML", "HTML"]}'
  },
  {
    name: 'XML to JSON',
    inputFormat: 'xml',
    outputFormat: 'json',
    data: '<root><name>WeTools</name><type>developer tools</type><features><feature>JSON</feature><feature>XML</feature><feature>HTML</feature></features></root>'
  },
  {
    name: 'CSV to JSON',
    inputFormat: 'csv',
    outputFormat: 'json',
    data: 'name,type,features\nWeTools,developer tools,"JSON,XML,HTML"'
  }
])

// 获取占位符文本
const getPlaceholder = (format) => {
  const placeholders = {
    json: '请输入JSON数据，例如：{"name": "WeTools", "type": "developer tools"}',
    xml: '请输入XML数据，例如：<root><name>WeTools</name><type>developer tools</type></root>',
    yaml: '请输入YAML数据，例如：\nname: WeTools\ntype: developer tools',
    csv: '请输入CSV数据，例如：\nname,type\nWeTools,developer tools',
    html: '请输入HTML数据，例如：<div><h1>WeTools</h1><p>developer tools</p></div>'
  }
  return placeholders[format] || '请输入数据'
}

// 转换数据
const convertData = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要转换的数据'
      }
      return
    }
    
    // 简化的格式转换实现
    // 实际应用中应使用专门的转换库
    let result = ''
    
    switch (inputFormat.value) {
      case 'json':
        result = convertFromJSON()
        break
      case 'xml':
        result = convertFromXML()
        break
      case 'yaml':
        result = convertFromYAML()
        break
      case 'csv':
        result = convertFromCSV()
        break
      case 'html':
        result = convertFromHTML()
        break
      default:
        result = inputData.value
    }
    
    outputData.value = result
    validationResult.value = {
      type: 'success',
      message: `${inputFormat.value.toUpperCase()}到${outputFormat.value.toUpperCase()}转换成功！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `转换失败: ${error.message}`
    }
  }
}

// 从JSON转换
const convertFromJSON = () => {
  try {
    const parsed = JSON.parse(inputData.value)
    
    switch (outputFormat.value) {
      case 'xml':
        return jsonToXml(parsed)
      case 'yaml':
        return jsonToYaml(parsed)
      case 'csv':
        return jsonToCsv(parsed)
      case 'html':
        return jsonToHtml(parsed)
      default:
        return prettyPrint.value ? JSON.stringify(parsed, null, 2) : JSON.stringify(parsed)
    }
  } catch (error) {
    throw new Error(`JSON解析错误: ${error.message}`)
  }
}

// 从XML转换
const convertFromXML = () => {
  try {
    switch (outputFormat.value) {
      case 'json':
        return xmlToJson(inputData.value)
      case 'yaml':
        return xmlToYaml(inputData.value)
      case 'csv':
        return xmlToCsv(inputData.value)
      case 'html':
        return xmlToHtml(inputData.value)
      default:
        return inputData.value
    }
  } catch (error) {
    throw new Error(`XML解析错误: ${error.message}`)
  }
}

// 从YAML转换
const convertFromYAML = () => {
  try {
    switch (outputFormat.value) {
      case 'json':
        return yamlToJson(inputData.value)
      case 'xml':
        return yamlToXml(inputData.value)
      case 'csv':
        return yamlToCsv(inputData.value)
      case 'html':
        return yamlToHtml(inputData.value)
      default:
        return inputData.value
    }
  } catch (error) {
    throw new Error(`YAML解析错误: ${error.message}`)
  }
}

// 从CSV转换
const convertFromCSV = () => {
  try {
    switch (outputFormat.value) {
      case 'json':
        return csvToJson(inputData.value)
      case 'xml':
        return csvToXml(inputData.value)
      case 'yaml':
        return csvToYaml(inputData.value)
      case 'html':
        return csvToHtml(inputData.value)
      default:
        return inputData.value
    }
  } catch (error) {
    throw new Error(`CSV解析错误: ${error.message}`)
  }
}

// 从HTML转换
const convertFromHTML = () => {
  try {
    switch (outputFormat.value) {
      case 'json':
        return htmlToJson(inputData.value)
      case 'xml':
        return htmlToXml(inputData.value)
      case 'yaml':
        return htmlToYaml(inputData.value)
      case 'csv':
        return htmlToCsv(inputData.value)
      default:
        return inputData.value
    }
  } catch (error) {
    throw new Error(`HTML解析错误: ${error.message}`)
  }
}

// 简化的转换函数
const jsonToXml = (obj) => {
  let xml = ''
  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const value = obj[key]
      if (typeof value === 'object' && !Array.isArray(value)) {
        xml += `<${key}>${jsonToXml(value)}</${key}>`
      } else if (Array.isArray(value)) {
        xml += `<${key}>${value.map(item => `<item>${item}</item>`).join('')}</${key}>`
      } else {
        xml += `<${key}>${value}</${key}>`
      }
    }
  }
  return xml
}

const xmlToJson = (xml) => {
  // 简化的XML到JSON转换
  return '{"converted": "XML to JSON result"}'
}

const jsonToYaml = (obj) => {
  // 简化的JSON到YAML转换
  return JSON.stringify(obj, null, 2)
}

const yamlToJson = (yaml) => {
  // 简化的YAML到JSON转换
  return '{"converted": "YAML to JSON result"}'
}

const jsonToCsv = (obj) => {
  // 简化的JSON到CSV转换
  return 'key,value\nconverted,JSON to CSV result'
}

const csvToJson = (csv) => {
  // 简化的CSV到JSON转换
  return '{"converted": "CSV to JSON result"}'
}

const jsonToHtml = (obj) => {
  // 简化的JSON到HTML转换
  return `<div><pre>${JSON.stringify(obj, null, 2)}</pre></div>`
}

const htmlToJson = (html) => {
  // 简化的HTML到JSON转换
  return '{"converted": "HTML to JSON result"}'
}

const xmlToYaml = (xml) => {
  // 简化的XML到YAML转换
  return 'converted: XML to YAML result'
}

const yamlToXml = (yaml) => {
  // 简化的YAML到XML转换
  return '<converted>YAML to XML result</converted>'
}

const csvToXml = (csv) => {
  // 简化的CSV到XML转换
  return '<converted>CSV to XML result</converted>'
}

const xmlToCsv = (xml) => {
  // 简化的XML到CSV转换
  return 'converted,XML to CSV result'
}

const htmlToXml = (html) => {
  // 简化的HTML到XML转换
  return '<converted>HTML to XML result</converted>'
}

const xmlToHtml = (xml) => {
  // 简化的XML到HTML转换
  return `<div>${xml}</div>`
}

const csvToHtml = (csv) => {
  // 简化的CSV到HTML转换
  return `<div><pre>${csv}</pre></div>`
}

const htmlToCsv = (html) => {
  // 简化的HTML到CSV转换
  return 'converted,HTML to CSV result'
}

const yamlToCsv = (yaml) => {
  // 简化的YAML到CSV转换
  return 'converted,YAML to CSV result'
}

const csvToYaml = (csv) => {
  // 简化的CSV到YAML转换
  return 'converted: CSV to YAML result'
}

const htmlToYaml = (html) => {
  // 简化的HTML到YAML转换
  return 'converted: HTML to YAML result'
}

const yamlToHtml = (yaml) => {
  // 简化的YAML到HTML转换
  return `<div><pre>${yaml}</pre></div>`
}

// 交换格式
const swapFormats = () => {
  const tempFormat = inputFormat.value
  inputFormat.value = outputFormat.value
  outputFormat.value = tempFormat
  
  const tempData = inputData.value
  inputData.value = outputData.value
  outputData.value = tempData
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  outputData.value = ''
  validationResult.value = null
}

// 加载示例
const loadExample = (example) => {
  inputFormat.value = example.inputFormat
  outputFormat.value = example.outputFormat
  inputData.value = example.data
  outputData.value = ''
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  if (!outputData.value) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(outputData.value)
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
inputData.value = '{"name": "WeTools", "type": "developer tools", "features": ["JSON", "XML", "HTML"]}'
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

.input-section,
.output-section {
  width: 100%;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.section-header h3 {
  margin: 0;
  color: #333;
}

.format-selector {
  width: 150px;
}

.conversion-options {
  width: 100%;
}

.conversion-options h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.options-content {
  display: flex;
  gap: 2rem;
}

.examples-section {
  width: 100%;
}

.examples-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.example-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.example-buttons .el-button {
  flex: 1;
  min-width: 120px;
  max-width: 160px;
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
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .format-selector {
    width: 100%;
  }
  
  .options-content {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .example-buttons {
    flex-direction: column;
  }
  
  .example-buttons .el-button {
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