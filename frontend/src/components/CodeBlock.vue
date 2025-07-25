<template>
  <div class="code-block-container">
    <div class="code-block-header" v-if="showHeader">
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
  background: #2d3748;
  margin: 10px 0;
}

.code-block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #4a5568;
  padding: 8px 12px;
  font-size: 12px;
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
</style>
</style>