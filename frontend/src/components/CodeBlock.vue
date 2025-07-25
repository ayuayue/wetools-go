<template>
  <div class="code-block-container" :data-theme="theme">
    <div class="code-block-header" :data-theme="theme" v-if="showHeader">
      <div class="language-info">{{ displayLanguage }}</div>
      <div class="actions">
        <button class="copy-btn" @click="copyCode">
          <i class="fas fa-copy"></i> {{ copied ? '已复制' : '复制' }}
        </button>
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
  </div>
</template>

<script setup>
import { ref, onMounted, onUpdated, computed, watch } from 'vue'
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

const codeBlockRef = ref(null)
const codeRef = ref(null)
const copied = ref(false)

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

// 在挂载时高亮代码
onMounted(() => {
  highlightCode()
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
}

/* 暗色主题 */
.code-block-container[data-theme="dark"] {
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

.language-info {
  font-weight: 500;
  text-transform: uppercase;
}

.actions {
  display: flex;
  gap: 8px;
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

.copy-btn:active {
  background: #16a34a;
}

.code-block {
  margin: 0;
  border-radius: 0;
  max-height: 400px;
  overflow: auto;
}

/* 覆盖Prism.js的默认样式 */
pre[class*="language-"] {
  margin: 0;
  border-radius: 0;
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
}

.code-block-container[data-theme="dark"] pre[class*="language-"] {
  /* 使用Prism暗色主题 */
  background: #2d3748 !important;
  color: #e2e8f0 !important;
}

.code-block-container[data-theme="dark"] .token.comment,
.code-block-container[data-theme="dark"] .token.prolog,
.code-block-container[data-theme="dark"] .token.doctype,
.code-block-container[data-theme="dark"] .token.cdata {
  color: #88c0d0 !important;
}

.code-block-container[data-theme="dark"] .token.punctuation {
  color: #eceff4 !important;
}

.code-block-container[data-theme="dark"] .token.property,
.code-block-container[data-theme="dark"] .token.tag,
.code-block-container[data-theme="dark"] .token.boolean,
.code-block-container[data-theme="dark"] .token.number,
.code-block-container[data-theme="dark"] .token.constant,
.code-block-container[data-theme="dark"] .token.symbol,
.code-block-container[data-theme="dark"] .token.deleted {
  color: #b48ead !important;
}

.code-block-container[data-theme="dark"] .token.selector,
.code-block-container[data-theme="dark"] .token.attr-name,
.code-block-container[data-theme="dark"] .token.string,
.code-block-container[data-theme="dark"] .token.char,
.code-block-container[data-theme="dark"] .token.builtin,
.code-block-container[data-theme="dark"] .token.inserted {
  color: #a3be8c !important;
}

.code-block-container[data-theme="dark"] .token.operator,
.code-block-container[data-theme="dark"] .token.entity,
.code-block-container[data-theme="dark"] .token.url,
.code-block-container[data-theme="dark"] .language-css .token.string,
.code-block-container[data-theme="dark"] .style .token.string {
  color: #81a1c1 !important;
}

.code-block-container[data-theme="dark"] .token.atrule,
.code-block-container[data-theme="dark"] .token.attr-value,
.code-block-container[data-theme="dark"] .token.keyword {
  color: #81a1c1 !important;
}

.code-block-container[data-theme="dark"] .token.function {
  color: #8fbcbb !important;
}

.code-block-container[data-theme="dark"] .token.regex,
.code-block-container[data-theme="dark"] .token.important,
.code-block-container[data-theme="dark"] .token.variable {
  color: #8fbcbb !important;
}
</style>