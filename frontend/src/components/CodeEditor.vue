<template>
  <div class="code-editor-container" :class="{ 'dark-theme': isDarkTheme }">
    <div class="editor-header" v-if="showHeader">
      <div class="language-info">{{ language }}</div>
      <div class="editor-actions">
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
import { ref, onMounted, onBeforeUnmount, watch, computed } from 'vue'
import { ElButton } from 'element-plus'
import { EditorView, basicSetup } from 'codemirror'
import { EditorState } from '@codemirror/state'
import { javascript } from '@codemirror/lang-javascript'
import { json } from '@codemirror/lang-json'
import { html } from '@codemirror/lang-html'
import { xml } from '@codemirror/lang-xml'
import { css } from '@codemirror/lang-css'
import { oneDark } from '@codemirror/theme-one-dark'
import { lineNumbers } from '@codemirror/view'

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
  theme: {
    type: String,
    default: 'light'
  }
})

const emit = defineEmits(['update:modelValue'])

const editorContainer = ref(null)
let editorView = null

// 计算是否为暗色主题
const isDarkTheme = computed(() => {
  return props.theme === 'dark' || document.body.classList.contains('dark-theme')
})

// 语言映射
const languageExtensions = {
  javascript: javascript(),
  js: javascript(),
  json: json(),
  html: html(),
  xml: xml(),
  css: css(),
  text: null
}

// 获取语言扩展
const getLanguageExtension = () => {
  return languageExtensions[props.language] || null
}

// 初始化编辑器
const initEditor = () => {
  if (!editorContainer.value) return

  const extensions = [
    basicSetup,
    EditorView.editable.of(!props.readonly),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        emit('update:modelValue', update.state.doc.toString())
      }
    })
  ]

  // 添加语言支持
  const languageExtension = getLanguageExtension()
  if (languageExtension) {
    extensions.push(languageExtension)
  }

  // 添加主题
  if (isDarkTheme.value) {
    extensions.push(oneDark)
  }

  // 添加行号
  if (props.showLineNumbers) {
    extensions.push(lineNumbers())
  }

  editorView = new EditorView({
    state: EditorState.create({
      doc: props.modelValue,
      extensions
    }),
    parent: editorContainer.value
  })
}

// 更新编辑器内容
const updateEditorContent = (newContent) => {
  if (editorView) {
    editorView.dispatch({
      changes: {
        from: 0,
        to: editorView.state.doc.length,
        insert: newContent
      }
    })
  }
}

// 格式化代码
const formatCode = () => {
  if (!editorView) return

  let formatted = props.modelValue
  try {
    if (props.language === 'json') {
      formatted = JSON.stringify(JSON.parse(props.modelValue), null, 2)
    } else if (props.language === 'html' || props.language === 'xml') {
      // 简单的格式化，实际项目中可能需要更复杂的格式化库
      formatted = props.modelValue.replace(/></g, '>\n<')
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
  if (editorView && newVal !== editorView.state.doc.toString()) {
    updateEditorContent(newVal)
  }
})

// 监听主题变化
watch(isDarkTheme, (newVal) => {
  if (editorView) {
    // 重新初始化编辑器以应用新主题
    editorView.destroy()
    initEditor()
  }
})

// 监听语言变化
watch(() => props.language, () => {
  if (editorView) {
    // 重新初始化编辑器以应用新语言
    editorView.destroy()
    initEditor()
  }
})

onMounted(() => {
  initEditor()
})

onBeforeUnmount(() => {
  if (editorView) {
    editorView.destroy()
  }
})
</script>

<style scoped>
.code-editor-container {
  border-radius: 4px;
  overflow: hidden;
  background: #f8fafc;
  margin: 10px 0;
  border: 1px solid #e2e8f0;
}

.code-editor-container.dark-theme {
  background: #2d3748;
  border: 1px solid #4a5568;
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

.code-editor-container.dark-theme .editor-header {
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

.code-editor-container.dark-theme .action-btn {
  background: #5a6fd8;
}

.code-editor-container.dark-theme .action-btn:hover {
  background: #4a5fc8;
}

.editor-container {
  min-height: 150px;
  max-height: 500px;
  overflow: auto;
}

.editor-container :deep(.cm-editor) {
  height: 100%;
}

.editor-container :deep(.cm-scroller) {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
}
</style>