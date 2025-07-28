<template>
  <el-card class="code-block-container" shadow="never" :data-theme="theme" :class="{ 'fullscreen': isFullscreen }">
    <div class="code-block-header" :data-theme="theme" v-if="showHeader">
      <div class="language-info">{{ displayLanguage }}</div>
      <div class="actions">
        <el-button 
          class="action-btn" 
          size="small" 
          @click="toggleFullscreen"
        >
          <i :class="isFullscreen ? 'fas fa-compress' : 'fas fa-expand'"></i> 
          {{ isFullscreen ? '退出全屏' : '全屏' }}
        </el-button>
        <el-button 
          class="action-btn" 
          size="small" 
          @click="formatCode"
          v-if="canFormat"
        >
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button 
          class="action-btn" 
          type="success" 
          size="small" 
          @click="copyCode"
        >
          <i class="fas fa-copy"></i> {{ copied ? '已复制' : '复制' }}
        </el-button>
      </div>
    </div>
    <pre 
      class="code-block" 
      :class="{ 'line-numbers': showLineNumbers }"
      ref="codeBlockRef"
    ><code 
      :class="`language-${prismLanguage}`" 
      ref="codeRef"
    >{{ code }}</code></pre>
  </el-card>
</template>

<script setup>
import { ref, onMounted, onUpdated, onBeforeUnmount, computed, watch } from 'vue'
import { ElCard, ElButton } from 'element-plus'
import Prism from 'prismjs'
import 'prismjs/themes/prism.css'
import 'prismjs/themes/prism-tomorrow.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers.css'
import 'prismjs/plugins/line-numbers/prism-line-numbers'

const props = defineProps({
  code: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'text'
  },
  showLineNumbers: {
    type: Boolean,
    default: true
  },
  showHeader: {
    type: Boolean,
    default: true
  },
  autoDetect: {
    type: Boolean,
    default: true
  },
  theme: {
    type: String,
    default: 'light'
  }
})

const emit = defineEmits(['update:code'])

const codeBlockRef = ref(null)
const codeRef = ref(null)
const copied = ref(false)
const isFullscreen = ref(false)

// 语言映射
const languageMap = {
  'json': 'json',
  'xml': 'xml',
  'html': 'html',
  'javascript': 'javascript',
  'js': 'javascript',
  'css': 'css',
  'text': 'text',
  'plaintext': 'text',
  'plain': 'text'
}

// 计算Prism语言类名
const prismLanguage = computed(() => {
  if (props.autoDetect && props.language === 'auto') {
    return detectLanguage(props.code)
  }
  return languageMap[props.language] || 'text'
})

// 显示语言名称
const displayLanguage = computed(() => {
  const lang = prismLanguage.value
  const langNames = {
    'json': 'JSON',
    'xml': 'XML',
    'html': 'HTML',
    'javascript': 'JavaScript',
    'css': 'CSS',
    'text': 'Text'
  }
  return langNames[lang] || lang
})

// 检查是否可以格式化
const canFormat = computed(() => {
  return ['json', 'html', 'xml', 'css', 'javascript', 'js'].includes(prismLanguage.value)
})

// 语言检测函数
const detectLanguage = (code) => {
  if (!code) return 'text'
  
  // JSON检测
  if (code.trim().startsWith('{') || code.trim().startsWith('[')) {
    try {
      JSON.parse(code)
      return 'json'
    } catch (e) {
      // 不是有效的JSON
    }
  }
  
  // XML检测
  if (code.includes('<?xml') || (code.includes('<') && code.includes('>'))) {
    return 'xml'
  }
  
  // HTML检测
  if (code.includes('<html') || code.includes('<div') || code.includes('<p>')) {
    return 'html'
  }
  
  return 'text'
}

// 高亮代码
const highlightCode = () => {
  if (codeRef.value) {
    Prism.highlightElement(codeRef.value)
    
    // 如果启用行号，初始化行号插件
    if (props.showLineNumbers && codeBlockRef.value) {
      Prism.plugins.lineNumbers.resize(codeBlockRef.value)
    }
  }
}

// 复制代码
const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(props.code)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = props.code
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  }
}

// 格式化代码
const formatCode = () => {
  let formatted = props.code
  try {
    if (prismLanguage.value === 'json') {
      formatted = JSON.stringify(JSON.parse(props.code), null, 2)
    } else if (prismLanguage.value === 'html' || prismLanguage.value === 'xml') {
      // 简单的HTML/XML格式化
      formatted = props.code.replace(/></g, '>\n<')
    } else if (prismLanguage.value === 'css') {
      // CSS格式化
      formatted = props.code.replace(/;/g, ';\n').replace(/}/g, '}\n')
    } else if (prismLanguage.value === 'javascript' || prismLanguage.value === 'js') {
      // JavaScript格式化（简单处理）
      formatted = props.code.replace(/{/g, '{\n').replace(/}/g, '\n}')
    }
    
    // 更新父组件的code值
    emit('update:code', formatted)
  } catch (e) {
    console.error('格式化错误:', e)
  }
}

// 全屏切换
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  if (isFullscreen.value) {
    document.body.classList.add('codeblock-fullscreen-active')
  } else {
    document.body.classList.remove('codeblock-fullscreen-active')
  }
}

// 在挂载时高亮代码
onMounted(() => {
  highlightCode()
})

// 在卸载前确保退出全屏状态
onBeforeUnmount(() => {
  if (isFullscreen.value) {
    document.body.classList.remove('codeblock-fullscreen-active')
  }
})

// 在更新时重新高亮代码
onUpdated(() => {
  highlightCode()
})

// 监听代码变化
watch(() => props.code, () => {
  highlightCode()
})

// 监听语言变化
watch(() => props.language, () => {
  highlightCode()
})
</script>

<style scoped>
.code-block-container {
  border-radius: 4px;
  overflow: hidden;
  background: #f8fafc;
  margin: 10px 0;
  border: 1px solid #e2e8f0;
  position: relative;
}

.code-block-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  margin: 0;
  border-radius: 0;
  border: none;
}

.codeblock-fullscreen-active {
  overflow: hidden;
}

/* 暗色主题 */
.code-block-container[data-theme="dark"] {
  background: #2d3748;
  border: 1px solid #4a5568;
}

/* 响应MainLayout的暗色主题 */
.dark-theme .code-block-container {
  background: #2d3748;
  border: 1px solid #4a5568;
}

.code-block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f1f5f9;
  padding: 8px 12px;
  font-size: 12px;
  color: #64748b;
}

/* 暗色主题 */
.code-block-header[data-theme="dark"] {
  background: #4a5568;
  color: #e2e8f0;
}

/* 响应MainLayout的暗色主题 */
.dark-theme .code-block-header {
  background: #4a5568;
  color: #e2e8f0;
}

.language-info {
  font-weight: 500;
  text-transform: uppercase;
}

.actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: background 0.2s;
}

.action-btn:hover {
  background: #5a6fd8;
}

.code-block-container[data-theme="dark"] .action-btn {
  background: #5a6fd8;
}

.code-block-container[data-theme="dark"] .action-btn:hover {
  background: #4a5fc8;
}

.copy-btn {
  background: #4ade80;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 4px 8px;
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: background 0.2s;
}

.copy-btn:hover {
  background: #22c55e;
}

/* 暗色主题下的复制按钮 */
.code-block-container[data-theme="dark"] .copy-btn {
  background: #22c55e;
}

.code-block-container[data-theme="dark"] .copy-btn:hover {
  background: #16a34a;
}

/* 响应MainLayout的暗色主题下的复制按钮 */
.dark-theme .code-block-container .copy-btn {
  background: #22c55e;
}

.dark-theme .code-block-container .copy-btn:hover {
  background: #16a34a;
}

.code-block {
  margin: 0;
  border-radius: 0;
  max-height: 800px; /* 增加高度到800px */
  min-height: 200px; /* 设置最小高度 */
  height: auto; /* 自动高度 */
  overflow: auto;
}

/* 覆盖Prism.js的默认样式 */
pre[class*="language-"] {
  margin: 0;
  border-radius: 0;
  text-shadow: none !important;
  box-shadow: none !important;
}

code[class*="language-"] {
  text-shadow: none !important;
  box-shadow: none !important;
}

pre.line-numbers {
  padding-left: 3.8em;
  counter-reset: linenumber;
}

pre.line-numbers > code {
  position: relative;
  white-space: pre;
  word-wrap: normal;
  word-break: normal;
  overflow-wrap: normal;
  tab-size: 4;
}

.line-numbers .line-numbers-rows {
  position: absolute;
  pointer-events: none;
  top: 0;
  font-size: 100%;
  left: -3.8em;
  width: 3em;
  letter-spacing: -1px;
  border-right: 1px solid #999;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.line-numbers-rows > span {
  display: block;
  counter-increment: linenumber;
}

.line-numbers-rows > span:before {
  content: counter(linenumber);
  color: #999;
  display: block;
  padding-right: 0.8em;
  text-align: right;
}

/* Prism主题切换样式 */
.code-block-container[data-theme="light"] pre[class*="language-"] {
  /* 使用Prism默认亮色主题 */
  text-shadow: none !important;
  box-shadow: none !important;
}

.code-block-container[data-theme="dark"] pre[class*="language-"] {
  /* 使用Prism暗色主题 */
  background: #2d3748 !important;
  color: #e2e8f0 !important;
  text-shadow: none !important;
  box-shadow: none !important;
}

/* 响应MainLayout的暗色主题 */
.dark-theme .code-block-container pre[class*="language-"] {
  background: #2d3748 !important;
  color: #e2e8f0 !important;
  text-shadow: none !important;
  box-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.comment,
.code-block-container[data-theme="dark"] .token.prolog,
.code-block-container[data-theme="dark"] .token.doctype,
.code-block-container[data-theme="dark"] .token.cdata {
  color: #88c0d0 !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.punctuation {
  color: #eceff4 !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.property,
.code-block-container[data-theme="dark"] .token.tag,
.code-block-container[data-theme="dark"] .token.boolean,
.code-block-container[data-theme="dark"] .token.number,
.code-block-container[data-theme="dark"] .token.constant,
.code-block-container[data-theme="dark"] .token.symbol,
.code-block-container[data-theme="dark"] .token.deleted {
  color: #b48ead !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.selector,
.code-block-container[data-theme="dark"] .token.attr-name,
.code-block-container[data-theme="dark"] .token.string,
.code-block-container[data-theme="dark"] .token.char,
.code-block-container[data-theme="dark"] .token.builtin,
.code-block-container[data-theme="dark"] .token.inserted {
  color: #a3be8c !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.operator,
.code-block-container[data-theme="dark"] .token.entity,
.code-block-container[data-theme="dark"] .token.url,
.code-block-container[data-theme="dark"] .language-css .token.string,
.code-block-container[data-theme="dark"] .style .token.string {
  color: #81a1c1 !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.atrule,
.code-block-container[data-theme="dark"] .token.attr-value,
.code-block-container[data-theme="dark"] .token.keyword {
  color: #81a1c1 !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.function {
  color: #8fbcbb !important;
  text-shadow: none !important;
}

.code-block-container[data-theme="dark"] .token.regex,
.code-block-container[data-theme="dark"] .token.important,
.code-block-container[data-theme="dark"] .token.variable {
  color: #8fbcbb !important;
  text-shadow: none !important;
}

/* 响应MainLayout的暗色主题 */
.dark-theme .code-block-container .token.comment,
.dark-theme .code-block-container .token.prolog,
.dark-theme .code-block-container .token.doctype,
.dark-theme .code-block-container .token.cdata {
  color: #88c0d0 !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.punctuation {
  color: #eceff4 !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.property,
.dark-theme .code-block-container .token.tag,
.dark-theme .code-block-container .token.boolean,
.dark-theme .code-block-container .token.number,
.dark-theme .code-block-container .token.constant,
.dark-theme .code-block-container .token.symbol,
.dark-theme .code-block-container .token.deleted {
  color: #b48ead !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.selector,
.dark-theme .code-block-container .token.attr-name,
.dark-theme .code-block-container .token.string,
.dark-theme .code-block-container .token.char,
.dark-theme .code-block-container .token.builtin,
.dark-theme .code-block-container .token.inserted {
  color: #a3be8c !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.operator,
.dark-theme .code-block-container .token.entity,
.dark-theme .code-block-container .token.url,
.dark-theme .code-block-container .language-css .token.string,
.dark-theme .code-block-container .style .token.string {
  color: #81a1c1 !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.atrule,
.dark-theme .code-block-container .token.attr-value,
.dark-theme .code-block-container .token.keyword {
  color: #81a1c1 !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.function {
  color: #8fbcbb !important;
  text-shadow: none !important;
}

.dark-theme .code-block-container .token.regex,
.dark-theme .code-block-container .token.important,
.dark-theme .code-block-container .token.variable {
  color: #8fbcbb !important;
  text-shadow: none !important;
}
</style>