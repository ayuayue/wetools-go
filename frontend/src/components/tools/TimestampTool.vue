<template>
  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-clock"></i> 当前时间</h3>
      <div class="current-time">
        <div class="time-item">
          <label>当前时间:</label>
          <span>{{ currentTime }}</span>
        </div>
        <div class="time-item">
          <label>时间戳(秒):</label>
          <span>{{ currentTimestamp }}</span>
        </div>
        <div class="time-item">
          <label>时间戳(毫秒):</label>
          <span>{{ currentTimestampMs }}</span>
        </div>
        <div class="time-item">
          <label>ISO 8601:</label>
          <span>{{ currentISO }}</span>
        </div>
      </div>
    </div>

    <div class="tool-section">
      <h3><i class="fas fa-exchange-alt"></i> 时间转换</h3>
      <div class="converter-group">
        <div class="converter-row">
          <label>输入时间:</label>
          <input 
            v-model="inputTime" 
            placeholder="输入时间戳或时间格式 (如: 2023-01-01 12:00:00)"
          />
        </div>
        <div class="converter-row">
          <label>输入格式:</label>
          <select v-model="inputFormat">
            <option value="auto">自动识别</option>
            <option value="timestamp">时间戳(秒)</option>
            <option value="timestamp_ms">时间戳(毫秒)</option>
            <option value="iso">ISO 8601</option>
            <option value="datetime">日期时间 (YYYY-MM-DD HH:mm:ss)</option>
            <option value="date">日期 (YYYY-MM-DD)</option>
          </select>
        </div>
      </div>
    </div>

    <div class="button-group">
      <button @click="convertTime">
        <i class="fas fa-sync"></i> 转换
      </button>
      <button @click="setCurrentTime">
        <i class="fas fa-clock"></i> 使用当前时间
      </button>
      <button class="secondary" @click="clearData">
        <i class="fas fa-trash"></i> 清空
      </button>
    </div>

    <div class="tool-section" v-if="convertedTime">
      <h3><i class="fas fa-file-alt"></i> 转换结果</h3>
      <div class="result-grid">
        <div class="result-item">
          <label>标准时间:</label>
          <span>{{ convertedTime.datetime }}</span>
        </div>
        <div class="result-item">
          <label>时间戳(秒):</label>
          <span>{{ convertedTime.timestamp }}</span>
        </div>
        <div class="result-item">
          <label>时间戳(毫秒):</label>
          <span>{{ convertedTime.timestamp_ms }}</span>
        </div>
        <div class="result-item">
          <label>ISO 8601:</label>
          <span>{{ convertedTime.iso }}</span>
        </div>
        <div class="result-item">
          <label>UTC时间:</label>
          <span>{{ convertedTime.utc }}</span>
        </div>
      </div>
    </div>
  </div>

  <div class="tool-container">
    <div class="tool-section">
      <h3><i class="fas fa-info-circle"></i> 工具说明</h3>
      <p>时间戳转换工具支持多种时间格式之间的相互转换：</p>
      <ul class="description-list">
        <li>实时显示当前时间及各种格式的时间戳</li>
        <li>支持时间戳(秒/毫秒)与标准时间互转</li>
        <li>支持ISO 8601格式时间转换</li>
        <li>支持自定义时间格式转换</li>
        <li>自动识别输入的时间格式</li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const currentTime = ref('')
const currentTimestamp = ref('')
const currentTimestampMs = ref('')
const currentISO = ref('')

const inputTime = ref('')
const inputFormat = ref('auto')
const convertedTime = ref(null)

let timer = null

// 更新当前时间
const updateCurrentTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString('zh-CN')
  currentTimestamp.value = Math.floor(now.getTime() / 1000)
  currentTimestampMs.value = now.getTime()
  currentISO.value = now.toISOString()
}

// 设置当前时间为输入值
const setCurrentTime = () => {
  inputTime.value = currentTimestampMs.value
  inputFormat.value = 'timestamp_ms'
  convertTime()
}

// 转换时间
const convertTime = () => {
  if (!inputTime.value) {
    convertedTime.value = null
    return
  }

  try {
    let date = null
    
    // 根据输入格式解析时间
    switch (inputFormat.value) {
      case 'timestamp':
        date = new Date(parseInt(inputTime.value) * 1000)
        break
      case 'timestamp_ms':
        date = new Date(parseInt(inputTime.value))
        break
      case 'iso':
        date = new Date(inputTime.value)
        break
      case 'datetime':
        // 处理 YYYY-MM-DD HH:mm:ss 格式
        date = new Date(inputTime.value.replace(' ', 'T'))
        break
      case 'date':
        // 处理 YYYY-MM-DD 格式
        date = new Date(inputTime.value)
        break
      case 'auto':
      default:
        // 自动识别格式
        if (/^\d{10}$/.test(inputTime.value)) {
          // 10位时间戳(秒)
          date = new Date(parseInt(inputTime.value) * 1000)
        } else if (/^\d{13}$/.test(inputTime.value)) {
          // 13位时间戳(毫秒)
          date = new Date(parseInt(inputTime.value))
        } else if (inputTime.value.includes('T') && inputTime.value.includes('Z')) {
          // ISO格式
          date = new Date(inputTime.value)
        } else if (inputTime.value.includes(' ')) {
          // YYYY-MM-DD HH:mm:ss 格式
          date = new Date(inputTime.value.replace(' ', 'T'))
        } else {
          // 默认尝试解析
          date = new Date(inputTime.value)
        }
        break
    }
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      throw new Error('无效的时间格式')
    }
    
    // 格式化结果
    convertedTime.value = {
      datetime: date.toLocaleString('zh-CN'),
      timestamp: Math.floor(date.getTime() / 1000),
      timestamp_ms: date.getTime(),
      iso: date.toISOString(),
      utc: date.toUTCString()
    }
  } catch (e) {
    convertedTime.value = {
      datetime: '转换错误: ' + e.message,
      timestamp: '',
      timestamp_ms: '',
      iso: '',
      utc: ''
    }
  }
}

// 清空数据
const clearData = () => {
  inputTime.value = ''
  convertedTime.value = null
}

// 初始化
onMounted(() => {
  updateCurrentTime()
  timer = setInterval(updateCurrentTime, 1000)
})

// 清理定时器
onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.tool-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.05);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

/* 暗色主题 */
.dark-theme .tool-container {
  background: #2d2d2d;
  box-shadow: 0 2px 10px rgba(0,0,0,0.3);
  color: #e0e0e0;
}

.tool-section {
  margin-bottom: 1.5rem;
}

.tool-section:last-child {
  margin-bottom: 0;
}

.tool-section h3 {
  margin-bottom: 1rem;
  color: #444;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 暗色主题 */
.dark-theme .tool-section h3 {
  color: #e0e0e0;
}

.current-time {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  background: #f8fafc;
  padding: 1rem;
  border-radius: 4px;
}

.time-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.time-item label {
  font-weight: 500;
  color: #64748b;
  font-size: 0.9rem;
}

.time-item span {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 1rem;
  color: #444;
}

.converter-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.converter-row {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.converter-row label {
  font-weight: 500;
  color: #64748b;
}

.converter-row input,
.converter-row select {
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 14px;
  transition: border-color 0.2s;
  background: white;
  color: #333;
}

.converter-row input:focus,
.converter-row select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.1);
}

/* 暗色主题 */
.dark-theme .converter-row input,
.dark-theme .converter-row select {
  background: #3d3d3d;
  border: 1px solid #555;
  color: #e0e0e0;
}

.dark-theme .converter-row input:focus,
.dark-theme .converter-row select:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.3);
}

.button-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin: 1rem 0;
}

button {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 5px;
}

button:hover {
  background: #5a6fd8;
}

button.secondary {
  background: #f1f5f9;
  color: #333;
}

button.secondary:hover {
  background: #e2e8f0;
}

/* 暗色主题 */
.dark-theme button {
  background: #5a6fd8;
  color: #e0e0e0;
}

.dark-theme button:hover {
  background: #4a5fc8;
}

.dark-theme button.secondary {
  background: #3d3d3d;
  color: #e0e0e0;
}

.dark-theme button.secondary:hover {
  background: #4d4d4d;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  background: #f8fafc;
  padding: 1rem;
  border-radius: 4px;
}

/* 暗色主题 */
.dark-theme .result-grid {
  background: #3d3d3d;
}

.result-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.result-item label {
  font-weight: 500;
  color: #64748b;
  font-size: 0.9rem;
}

.result-item span {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 1rem;
  color: #444;
  word-break: break-all;
}

/* 暗色主题 */
.dark-theme .result-item label {
  color: #a0aec0;
}

.dark-theme .result-item span {
  color: #e0e0e0;
}

.description-list {
  margin-top: 10px;
  padding-left: 20px;
}

@media (max-width: 768px) {
  .button-group {
    flex-direction: column;
  }
  
  button {
    width: 100%;
    justify-content: center;
  }
  
  .current-time,
  .result-grid {
    grid-template-columns: 1fr;
  }
}
</style>