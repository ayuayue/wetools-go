<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>HTML运行器</h2>
      <p>在线编写和运行HTML代码</p>
    </div>
    
    <div class="tool-content">
      <!-- 两栏布局 -->
      <div class="two-column-layout">
        <!-- 左侧代码编辑器 -->
        <div class="column code-column">
          <div class="editor-section">
            <div class="editor-header">
              <h3>HTML代码编辑器</h3>
              <div class="editor-actions">
                <el-button type="primary" @click="runCode">
                  <i class="fas fa-play"></i> 运行
                </el-button>
                <el-button @click="clearCode">
                  <i class="fas fa-trash"></i> 清空
                </el-button>
                <el-button @click="formatCode">
                  <i class="fas fa-indent"></i> 格式化
                </el-button>
              </div>
            </div>
            <AceEditor
              ref="editorRef"
              v-model="htmlCode"
              language="html"
              :theme="theme"
              :showFullscreenButton="true"
            />
          </div>
        </div>
        
        <!-- 右侧结果预览 -->
        <div class="column preview-column">
          <div class="preview-section">
            <div class="preview-header">
              <h3>结果预览</h3>
              <div class="preview-actions">
                <el-button @click="togglePreviewFullscreen">
                  <i :class="isPreviewFullscreen ? 'fas fa-compress' : 'fas fa-expand'"></i>
                  {{ isPreviewFullscreen ? '退出全屏' : '全屏' }}
                </el-button>
                <el-button @click="clearOutput" v-if="htmlCode">
                  <i class="fas fa-times"></i> 清空结果
                </el-button>
              </div>
            </div>
            <div class="preview-content" :class="{ 'fullscreen': isPreviewFullscreen }">
              <iframe 
                ref="outputFrame" 
                class="result-frame"
                sandbox="allow-scripts allow-same-origin"
                frameborder="0"
              ></iframe>
            </div>
          </div>
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
import { ref, onMounted, inject } from 'vue'
import { ElCard, ElButton, ElAlert } from 'element-plus'
import AceEditor from '../AceEditor.vue'

// 注入主题
const theme = inject('theme', ref('light'))

// 编辑器引用
const editorRef = ref(null)
const outputFrame = ref(null)

// 代码和输出状态
const htmlCode = ref(`<!DOCTYPE html>
<html>
<head>
    <title>HTML运行器示例</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .highlight {
            background-color: #ffeb3b;
            padding: 2px 4px;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>欢迎使用HTML运行器</h1>
        <p>这是一个<span class="highlight">HTML运行器</span>示例页面。</p>
        <p>您可以在这里编写和测试HTML代码。</p>
        <ul>
            <li>支持HTML标签</li>
            <li>支持内联CSS样式</li>
            <li>支持JavaScript脚本</li>
        </ul>
        <button onclick="alert('Hello, World!')">点击我</button>
    </div>
</body>
</html>`)

const validationResult = ref(null)
const isPreviewFullscreen = ref(false)

// 运行代码
const runCode = () => {
  try {
    if (!htmlCode.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入HTML代码'
      }
      return
    }
    
    // 在iframe中运行HTML代码
    if (outputFrame.value) {
      const iframe = outputFrame.value
      const iframeDoc = iframe.contentDocument || iframe.contentWindow.document
      iframeDoc.open()
      iframeDoc.write(htmlCode.value)
      iframeDoc.close()
      
      validationResult.value = {
        type: 'success',
        message: 'HTML代码运行成功'
      }
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `运行代码时出错: ${error.message}`
    }
  }
}

// 清空代码
const clearCode = () => {
  htmlCode.value = ''
  validationResult.value = null
}

// 格式化代码
const formatCode = () => {
  try {
    // 简单的HTML格式化
    let formatted = htmlCode.value
      .replace(/></g, '>\n<')
      .replace(/\n\s*\n/g, '\n')
      .trim()
    
    htmlCode.value = formatted
    validationResult.value = {
      type: 'success',
      message: '代码格式化完成'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `格式化代码时出错: ${error.message}`
    }
  }
}

// 切换预览全屏
const togglePreviewFullscreen = () => {
  isPreviewFullscreen.value = !isPreviewFullscreen.value
  if (isPreviewFullscreen.value) {
    document.body.classList.add('preview-fullscreen-active')
  } else {
    document.body.classList.remove('preview-fullscreen-active')
  }
}

// 清空输出
const clearOutput = () => {
  if (outputFrame.value) {
    const iframe = outputFrame.value
    const iframeDoc = iframe.contentDocument || iframe.contentWindow.document
    iframeDoc.open()
    iframeDoc.write('')
    iframeDoc.close()
  }
}

// 组件挂载后的处理
onMounted(() => {
  // 初始化运行一次示例代码
  runCode()
})
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
  height: calc(100vh - 200px);
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
  height: calc(100% - 100px);
}

/* 两栏布局 */
.two-column-layout {
  display: flex;
  gap: 1.5rem;
  flex: 1;
  min-height: 0;
}

.column {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.code-column {
  flex: 1;
}

.preview-column {
  flex: 1;
}

.editor-section,
.preview-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  height: 100%;
}

.editor-header,
.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.editor-header h3,
.preview-header h3 {
  margin: 0;
  color: #333;
}

.editor-actions,
.preview-actions {
  display: flex;
  gap: 0.5rem;
}

.preview-content {
  flex: 1;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  background: white;
  position: relative;
}

.preview-content.fullscreen {
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

.preview-fullscreen-active {
  overflow: hidden;
}

.result-frame {
  width: 100%;
  height: 100%;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .two-column-layout {
    flex-direction: column;
  }
  
  .editor-header,
  .preview-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .editor-actions,
  .preview-actions {
    justify-content: center;
  }
}
</style>