<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>JavaScript运行器</h2>
      <p>在线编写和运行JavaScript代码</p>
    </div>
    
    <div class="tool-content">
      <!-- 两栏布局 -->
      <div class="two-column-layout">
        <!-- 左侧代码编辑器 -->
        <div class="column code-column">
          <div class="editor-section">
            <div class="editor-header">
              <h3>JavaScript代码编辑器</h3>
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
              v-model="jsCode"
              language="javascript"
              :theme="theme"
              :showFullscreenButton="true"
            />
          </div>
        </div>
        
        <!-- 右侧结果输出 -->
        <div class="column output-column">
          <div class="output-section">
            <div class="output-header">
              <h3>运行结果</h3>
              <div class="output-actions">
                <el-button @click="toggleOutputFullscreen">
                  <i :class="isOutputFullscreen ? 'fas fa-compress' : 'fas fa-expand'"></i>
                  {{ isOutputFullscreen ? '退出全屏' : '全屏' }}
                </el-button>
                <el-button @click="clearOutput" v-if="outputContent">
                  <i class="fas fa-times"></i> 清空结果
                </el-button>
              </div>
            </div>
            <div class="output-content" :class="{ 'fullscreen': isOutputFullscreen }">
              <AceEditor
                v-model="outputContent"
                language="text"
                :theme="theme"
                :readonly="true"
                :showHeader="false"
                :showFullscreenButton="false"
                :showFormatButton="false"
                :showLineNumbers="false"
              />
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

// 代码和输出状态
const jsCode = ref(`// JavaScript运行器示例
console.log('欢迎使用JavaScript运行器！');

// 基本变量和操作
const name = 'WeTools';
let count = 0;

console.log(\`Hello, \${name}!\`);
console.log('计数器初始值:', count);

// 函数示例
function greet(person) {
    return \`你好, \${person}!\`;
}

console.log(greet('开发者'));

// 数组操作
const fruits = ['苹果', '香蕉', '橙子'];
console.log('水果列表:', fruits);

fruits.forEach((fruit, index) => {
    console.log(\`第\${index + 1}个水果: \${fruit}\`);
});

// 对象示例
const user = {
    name: '张三',
    age: 25,
    skills: ['JavaScript', 'Vue', 'Node.js'],
    greet() {
        return \`我是\${this.name}，今年\${this.age}岁\`;
    }
};

console.log(user.greet());
console.log('技能:', user.skills);

// 异步操作示例
async function fetchData() {
    console.log('开始获取数据...');
    
    // 模拟异步操作
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    console.log('数据获取完成！');
    return { message: '这是模拟的数据' };
}

fetchData().then(data => {
    console.log('获取到的数据:', data);
});

// 数学计算
const numbers = [1, 2, 3, 4, 5];
const sum = numbers.reduce((acc, num) => acc + num, 0);
const average = sum / numbers.length;

console.log(\`数字: [\${numbers.join(', ')}]\`);
console.log(\`总和: \${sum}\`);
console.log(\`平均值: \${average.toFixed(2)}\`);`)

const outputContent = ref('')
const validationResult = ref(null)
const isOutputFullscreen = ref(false)

// 运行代码
const runCode = () => {
  try {
    if (!jsCode.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入JavaScript代码'
      }
      return
    }
    
    // 捕获console输出
    const logs = []
    const originalLog = console.log
    
    console.log = (...args) => {
      const formattedArgs = args.map(arg => {
        // 对于对象和数组，使用格式化显示
        if (typeof arg === 'object' && arg !== null) {
          try {
            return JSON.stringify(arg, null, 2)
          } catch (e) {
            // 如果JSON.stringify失败，使用toString
            return String(arg)
          }
        }
        return String(arg)
      })
      // 保持原有console.log的行为，将所有参数在同一行显示
      logs.push(formattedArgs.join(' '))
      originalLog.apply(console, args)
    }
    
    try {
      // 执行JavaScript代码
      eval(jsCode.value)
      
      // 显示输出，每条console.log占一行
      outputContent.value = logs.join('\n') || '代码已执行（无输出）'
      
      validationResult.value = {
        type: 'success',
        message: 'JavaScript代码运行成功'
      }
    } catch (error) {
      outputContent.value = `错误: ${error.message}`
      validationResult.value = {
        type: 'error',
        message: `运行代码时出错: ${error.message}`
      }
    } finally {
      // 恢复原始console.log
      console.log = originalLog
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
  jsCode.value = ''
  validationResult.value = null
}

// 格式化代码
const formatCode = () => {
  // JavaScript代码格式化可以使用更专业的库
  validationResult.value = {
    type: 'info',
    message: 'JavaScript格式化功能待实现'
  }
}

// 切换输出全屏
const toggleOutputFullscreen = () => {
  isOutputFullscreen.value = !isOutputFullscreen.value
  if (isOutputFullscreen.value) {
    document.body.classList.add('output-fullscreen-active')
  } else {
    document.body.classList.remove('output-fullscreen-active')
  }
}

// 清空输出
const clearOutput = () => {
  outputContent.value = ''
}

// 组件挂载后的处理
onMounted(() => {
  // 可以在这里添加一些初始化逻辑
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

.output-column {
  flex: 1;
}

.editor-section,
.output-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  height: 100%;
}

.editor-header,
.output-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.editor-header h3,
.output-header h3 {
  margin: 0;
  color: #333;
}

.editor-actions,
.output-actions {
  display: flex;
  gap: 0.5rem;
}

.output-content {
  flex: 1;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  background: #f8f9fa;
  min-height: 150px;
  position: relative;
  height: 100%;
}

.output-content.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 9999;
  margin: 0;
  border-radius: 0;
  border: none;
  padding: 2rem;
}

.output-fullscreen-active {
  overflow: hidden;
}

.tool-container.dark-theme .output-content {
  background: #1e1e1e;
  color: #d4d4d4;
  border: 1px solid #444;
}

.output-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .two-column-layout {
    flex-direction: column;
  }
  
  .editor-header,
  .output-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .editor-actions,
  .output-actions {
    justify-content: center;
  }
}
</style>