<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>XML 格式化工具</h2>
      <p>在线XML格式化、校验、压缩工具</p>
    </div>
    
    <div class="tool-content">
      <!-- 两栏布局 -->
      <div class="two-column-layout">
        <!-- 左侧输入区域 -->
        <div class="column input-column">
          <div class="input-section">
            <div class="editor-header">
              <h3>XML输入</h3>
              <div class="editor-actions">
                <!-- 占位元素，确保与右侧高度一致 -->
              </div>
            </div>
            <AceEditor
              v-model="inputData"
              language="xml"
              :theme="theme"
              :showHeader="false"
            />
          </div>
        </div>
        
        <!-- 右侧输出区域 -->
        <div class="column output-column">
          <div class="output-section">
            <div class="editor-header">
              <h3>XML输出</h3>
              <div class="editor-actions">
                <el-button type="success" @click="copyResult">
                  <i class="fas fa-copy"></i> 复制结果
                </el-button>
              </div>
            </div>
            <AceEditor
              v-model="outputData"
              language="xml"
              :theme="theme"
              :readonly="true"
              :showHeader="false"
            />
          </div>
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="formatXml">
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button @click="compressXml">
          <i class="fas fa-compress"></i> 压缩
        </el-button>
        <el-button @click="validateXml">
          <i class="fas fa-check-circle"></i> 校验
        </el-button>
        <el-button @click="clearData">
          <i class="fas fa-trash"></i> 清空
        </el-button>
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
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert } from 'element-plus'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 数据模型
const inputData = ref('')
const outputData = ref('')
const validationResult = ref(null)

// 格式化XML
const formatXml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 简单的XML格式化函数
    const formatted = formatXmlString(inputData.value)
    outputData.value = formatted
    validationResult.value = {
      type: 'success',
      message: 'XML格式化成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML格式错误: ${error.message}`
    }
  }
}

// 压缩XML
const compressXml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除多余的空白字符
    const compressed = inputData.value.replace(/>\s+</g, '><').trim()
    outputData.value = compressed
    validationResult.value = {
      type: 'success',
      message: 'XML压缩成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML处理错误: ${error.message}`
    }
  }
}

// 校验XML
const validateXml = () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'info',
        message: '请输入XML数据进行校验'
      }
      return
    }
    
    // 简单的XML校验
    const parser = new DOMParser()
    const doc = parser.parseFromString(inputData.value, 'text/xml')
    
    // 检查是否有解析错误
    const errorNode = doc.querySelector('parsererror')
    if (errorNode) {
      throw new Error('XML格式不正确')
    }
    
    validationResult.value = {
      type: 'success',
      message: 'XML格式正确！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `XML格式错误: ${error.message}`
    }
  }
}

// 清空数据
const clearData = () => {
  inputData.value = ''
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

// 改进的XML格式化函数
const formatXmlString = (xml) => {
  // 移除多余的空白字符但保留必要的结构
  let formatted = xml.trim()
  
  // 使用更智能的XML格式化方法
  const tab = '  '; // 2个空格缩进
  let result = '';
  let indentLevel = 0;
  
  // 预定义自闭合标签
  const selfClosingTags = ['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr'];
  
  // 使用正则表达式解析XML标签
  const tagRegex = /<[^>]+>/g;
  let lastIndex = 0;
  let match;
  
  while ((match = tagRegex.exec(formatted)) !== null) {
    // 添加标签前的文本内容
    const textBeforeTag = formatted.substring(lastIndex, match.index);
    if (textBeforeTag.trim()) {
      result += textBeforeTag.trim();
    }
    
    const tag = match[0];
    const tagName = tag.match(/<\/?([a-zA-Z]+)/)?.[1]?.toLowerCase() || '';
    const isEndTag = tag.startsWith('</');
    const isSelfClosing = tag.endsWith('/>') || selfClosingTags.includes(tagName);
    
    // 处理标签
    if (isEndTag) {
      // 结束标签
      indentLevel = Math.max(0, indentLevel - 1);
      result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
    } else if (isSelfClosing) {
      // 自闭合标签
      result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
    } else {
      // 开始标签
      result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
      indentLevel++;
    }
    
    lastIndex = match.index + tag.length;
  }
  
  // 添加剩余的文本
  const remainingText = formatted.substring(lastIndex);
  if (remainingText.trim()) {
    result += remainingText.trim();
  }
  
  // 清理多余的空白行，但保留必要的结构
  return result
    .replace(/\n\s*\n/g, '\n') // 移除多余空行
    .replace(/^\s*\n/g, '') // 移除开头的空行
    .trim();
}

// 初始化示例数据
inputData.value = '<root><name>WeTools</name><type>developer tools</type><features><feature>JSON</feature><feature>XML</feature><feature>HTML</feature></features></root>'
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
}

.tool-header {
  margin-bottom: 1rem;
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
}

.tool-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 两栏布局 */
.two-column-layout {
  display: flex;
  gap: 1.5rem;
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

.button-group {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.button-group .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
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