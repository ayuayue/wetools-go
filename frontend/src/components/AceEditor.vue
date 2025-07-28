<template>
  <div class="ace-editor-container" :class="{ 'dark-theme': isDarkTheme, 'fullscreen': isFullscreen }">
    <div class="editor-header" v-if="showHeader">
      <div class="language-info">{{ language }}</div>
      <div class="editor-actions">
        <el-button 
          class="action-btn" 
          size="small" 
          @click="toggleFullscreen"
          v-if="showFullscreenButton"
        >
          <i :class="isFullscreen ? 'fas fa-compress' : 'fas fa-expand'"></i> 
          {{ isFullscreen ? '退出全屏' : '全屏' }}
        </el-button>
        <el-button 
          class="action-btn" 
          size="small" 
          @click="formatCode" 
          v-if="showFormatButton"
        >
          <i class="fas fa-magic"></i> 格式化
        </el-button>
        <el-button 
          class="action-btn" 
          type="success" 
          size="small" 
          @click="copyCode"
        >
          <i class="fas fa-copy"></i> 复制
        </el-button>
      </div>
    </div>
    <div ref="editorContainer" class="editor-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, computed, nextTick } from 'vue'
import { ElButton } from 'element-plus'
import ace from 'ace-builds'
import 'ace-builds/src-noconflict/mode-javascript'
import 'ace-builds/src-noconflict/mode-json'
import 'ace-builds/src-noconflict/mode-html'
import 'ace-builds/src-noconflict/mode-xml'
import 'ace-builds/src-noconflict/mode-css'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/theme-dracula'
import 'ace-builds/src-noconflict/ext-language_tools'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'text'
  },
  readonly: {
    type: Boolean,
    default: false
  },
  showHeader: {
    type: Boolean,
    default: true
  },
  showLineNumbers: {
    type: Boolean,
    default: true
  },
  showFormatButton: {
    type: Boolean,
    default: true
  },
  showFullscreenButton: {
    type: Boolean,
    default: true
  },
  theme: {
    type: String,
    default: 'light'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorContainer = ref(null)
let editor = null

// 计算是否为暗色主题
const isDarkTheme = computed(() => {
  return props.theme === 'dark' || document.body.classList.contains('dark-theme')
})

// 全屏状态
const isFullscreen = ref(false)

// 语言映射
const languageModes = {
  javascript: 'javascript',
  js: 'javascript',
  json: 'json',
  html: 'html',
  xml: 'xml',
  css: 'css',
  text: 'text'
}

// 获取语言模式
const getLanguageMode = () => {
  return languageModes[props.language] || 'text'
}

// 初始化编辑器
const initEditor = () => {
  if (!editorContainer.value) return

  editor = ace.edit(editorContainer.value, {
    maxLines: 30,
    theme: isDarkTheme.value ? 'ace/theme/dracula' : 'ace/theme/chrome',
    mode: `ace/mode/${getLanguageMode()}`,
    value: props.modelValue,
    readOnly: props.readonly,
    showLineNumbers: props.showLineNumbers,
    fontSize: 14,
    fontFamily: 'Consolas, Monaco, "Courier New", monospace',
    tabSize: 2,
    useSoftTabs: true,
    wrap: true, // 自动换行
    showPrintMargin: false,
    highlightActiveLine: true,
    highlightSelectedWord: true,
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
    enableSnippets: true
  })

  // 监听内容变化
  editor.session.on('change', () => {
    const content = editor.getValue()
    emit('update:modelValue', content)
  })
}

// 更新编辑器内容
const updateEditorContent = (newContent) => {
  if (editor) {
    editor.setValue(newContent, -1) // -1 表示将光标移到开头
  }
}

// 格式化代码
const formatCode = () => {
  if (!editor) return

  let formatted = props.modelValue
  try {
    if (props.language === 'json') {
      formatted = JSON.stringify(JSON.parse(props.modelValue), null, 2)
    } else if (props.language === 'html' || props.language === 'xml') {
      // 简单的格式化，实际项目中可能需要更复杂的格式化库
      formatted = props.modelValue.replace(/></g, '>\n<')
    } else if (props.language === 'css') {
      // CSS格式化
      formatted = props.modelValue.replace(/;/g, ';\n').replace(/}/g, '}\n')
    }
    updateEditorContent(formatted)
    emit('update:modelValue', formatted)
  } catch (e) {
    console.error('格式化错误:', e)
  }
}

// 复制代码
const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(props.modelValue)
  } catch (err) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = props.modelValue
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}

// 监听modelValue变化
watch(() => props.modelValue, (newVal) => {
  if (editor && newVal !== editor.getValue()) {
    updateEditorContent(newVal)
  }
})

// 监听主题变化
watch(isDarkTheme, (newVal) => {
  if (editor) {
    editor.setTheme(newVal ? 'ace/theme/dracula' : 'ace/theme/chrome')
  }
})

// 监听语言变化
watch(() => props.language, () => {
  if (editor) {
    editor.session.setMode(`ace/mode/${getLanguageMode()}`)
  }
})

onMounted(() => {
  initEditor()
})

// 全屏切换
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  
  if (editorContainer.value) {
    if (isFullscreen.value) {
      // 进入全屏
      editorContainer.value.classList.add('fullscreen-editor')
      document.body.classList.add('editor-fullscreen-active')
    } else {
      // 退出全屏
      editorContainer.value.classList.remove('fullscreen-editor')
      document.body.classList.remove('editor-fullscreen-active')
    }
    
    // 调整编辑器大小
    if (editor) {
      nextTick(() => {
        editor.resize()
      })
    }
  }
}

onBeforeUnmount(() => {
  // 确保退出全屏状态
  if (isFullscreen.value) {
    document.body.classList.remove('editor-fullscreen-active')
  }
  
  if (editor) {
    editor.destroy()
  }
})
</script>

<style scoped>
.ace-editor-container {
  border-radius: 4px;
  overflow: hidden;
  background: #f8fafc;
  margin: 10px 0;
  border: 1px solid #e2e8f0;
  position: relative;
}

.ace-editor-container.dark-theme {
  background: #2d3748;
  border: 1px solid #4a5568;
}

.ace-editor-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  margin: 0;
  border-radius: 0;
}

.editor-fullscreen-active {
  overflow: hidden;
}

.fullscreen-editor {
  height: calc(100vh - 40px) !important;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f1f5f9;
  padding: 8px 12px;
  font-size: 12px;
  color: #64748b;
}

.ace-editor-container.dark-theme .editor-header {
  background: #4a5568;
  color: #e2e8f0;
}

.language-info {
  font-weight: 500;
  text-transform: uppercase;
}

.editor-actions {
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

.ace-editor-container.dark-theme .action-btn {
  background: #5a6fd8;
}

.ace-editor-container.dark-theme .action-btn:hover {
  background: #4a5fc8;
}

.editor-container {
  min-height: 550px;
  height: 100%;
  overflow: auto;
  position: relative;
}

.editor-container :deep(.ace_editor) {
  font-family: 'Consolas', 'Monaco', monospace !important;
  font-size: 14px !important;
  line-height: 1.5 !important;
}

.editor-container :deep(.ace_content) {
  font-family: 'Consolas', 'Monaco', monospace !important;
}
</style>