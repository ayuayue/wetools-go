<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>HTML 格式化工具</h2>
      <p>在线HTML格式化、清理、压缩工具</p>
    </div>
    
    <div class="tool-content">
      <!-- 两栏布局 -->
      <div class="two-column-layout">
        <!-- 左侧输入区域 -->
        <div class="column input-column">
          <div class="input-section">
            <div class="editor-header">
              <h3>HTML输入</h3>
              <div class="editor-actions">
                <!-- 占位元素，确保与右侧高度一致 -->
              </div>
            </div>
            <AceEditor
              v-model="inputData"
              language="html"
              :theme="theme"
              :showHeader="false"
            />
          </div>
        </div>
        
        <!-- 右侧输出区域 -->
        <div class="column output-column">
          <div class="output-section">
            <div class="editor-header">
              <h3>HTML输出</h3>
              <div class="editor-actions">
                <el-button type="success" @click="copyResult">
                  <i class="fas fa-copy"></i> 复制结果
                </el-button>
              </div>
            </div>
            <AceEditor
              v-model="outputData"
              language="html"
              :theme="theme"
              :readonly="true"
              :showHeader="false"
            />
          </div>
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="formatHtml">
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button @click="compressHtml">
          <i class="fas fa-compress"></i> 压缩
        </el-button>
        <el-button @click="cleanHtml">
          <i class="fas fa-broom"></i> 清理
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

// 格式化HTML
const formatHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 简单的HTML格式化函数
    const formatted = formatHtmlString(inputData.value)
    outputData.value = formatted
    validationResult.value = {
      type: 'success',
      message: 'HTML格式化成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML格式错误: ${error.message}`
    }
  }
}

// 压缩HTML
const compressHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除多余的空白字符
    const compressed = inputData.value
      .replace(/>\s+</g, '><')
      .replace(/\s+/g, ' ')
      .trim()
    outputData.value = compressed
    validationResult.value = {
      type: 'success',
      message: 'HTML压缩成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML处理错误: ${error.message}`
    }
  }
}

// 清理HTML
const cleanHtml = () => {
  try {
    if (!inputData.value.trim()) {
      outputData.value = ''
      validationResult.value = null
      return
    }
    
    // 移除注释和多余的空白
    const cleaned = inputData.value
      .replace(/<!--[\s\S]*?-->/g, '')
      .replace(/>\s+</g, '><')
      .replace(/\s+/g, ' ')
      .trim()
    outputData.value = cleaned
    validationResult.value = {
      type: 'success',
      message: 'HTML清理成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `HTML处理错误: ${error.message}`
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

// 改进的HTML格式化函数
const formatHtmlString = (html) => {
  // 移除多余的空白字符但保留必要的结构
  let formatted = html
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&')
  
  // 使用更智能的HTML格式化方法
  const tab = '  '; // 2个空格缩进
  let result = '';
  let indentLevel = 0;
  
  // 预定义自闭合标签和特殊标签
  const selfClosingTags = ['area', 'base', 'br', 'col', 'embed', 'hr', 'img', 'input', 'link', 'meta', 'param', 'source', 'track', 'wbr'];
  const inlineTags = ['a', 'abbr', 'acronym', 'b', 'bdo', 'big', 'br', 'button', 'cite', 'code', 'dfn', 'em', 'i', 'img', 'input', 'kbd', 'label', 'map', 'object', 'q', 'samp', 'script', 'select', 'small', 'span', 'strong', 'sub', 'sup', 'textarea', 'time', 'tt', 'var'];
  const formatTags = ['html', 'head', 'body', 'div', 'section', 'article', 'aside', 'nav', 'header', 'footer', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'p', 'blockquote', 'pre', 'ul', 'ol', 'li', 'dl', 'dt', 'dd', 'table', 'thead', 'tbody', 'tfoot', 'tr', 'td', 'th', 'form', 'fieldset', 'legend'];
  
  // 使用正则表达式解析HTML标签
  const tagRegex = /<[^>]+>/g;
  let lastIndex = 0;
  let match;
  
  while ((match = tagRegex.exec(formatted)) !== null) {
    // 添加标签前的文本内容
    const textBeforeTag = formatted.substring(lastIndex, match.index);
    if (textBeforeTag.trim()) {
      // 如果文本前后有换行符，保持缩进
      const lines = textBeforeTag.split('\n');
      result += lines.map((line, index) => {
        if (index === 0 && !line.trim()) return line;
        return line.trim() ? tab.repeat(indentLevel) + line.trim() : line;
      }).join('\n') + '\n';
    }
    
    const tag = match[0];
    const tagName = tag.match(/<\/?([a-zA-Z]+)/)?.[1]?.toLowerCase() || '';
    const isEndTag = tag.startsWith('</');
    const isSelfClosing = tag.endsWith('/>') || selfClosingTags.includes(tagName);
    const isInlineTag = inlineTags.includes(tagName);
    const shouldFormatTag = formatTags.includes(tagName);
    
    // 处理标签
    if (isEndTag) {
      // 结束标签
      if (shouldFormatTag && !isInlineTag) {
        indentLevel = Math.max(0, indentLevel - 1);
        result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
      } else {
        result += tag;
      }
    } else if (isSelfClosing) {
      // 自闭合标签
      if (shouldFormatTag) {
        result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
      } else {
        result += tag;
      }
    } else {
      // 开始标签
      if (shouldFormatTag && !isInlineTag) {
        result += (result.endsWith('\n') ? '' : '\n') + tab.repeat(indentLevel) + tag + '\n';
        indentLevel++;
      } else {
        result += tag;
      }
    }
    
    lastIndex = match.index + tag.length;
  }
  
  // 添加剩余的文本
  const remainingText = formatted.substring(lastIndex);
  if (remainingText.trim()) {
    result += tab.repeat(indentLevel) + remainingText.trim();
  }
  
  // 清理多余的空白行，但保留必要的结构
  return result
    .replace(/\n\s*\n/g, '\n') // 移除多余空行
    .replace(/^\s*\n/g, '') // 移除开头的空行
    .trim();
}

// 初始化示例数据
inputData.value = '<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello World</h1><p>This is a paragraph.</p></body></html>'
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