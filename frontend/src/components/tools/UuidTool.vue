<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>UUID 生成器</h2>
      <p>UUID生成器，支持多种UUID版本</p>
    </div>
    
    <div class="tool-content">
      <div class="uuid-version-section">
        <h3>选择UUID版本</h3>
        <div class="version-buttons">
          <el-button 
            v-for="version in uuidVersions" 
            :key="version.value"
            :type="selectedVersion === version.value ? 'primary' : ''"
            @click="selectVersion(version.value)"
          >
            {{ version.name }}
          </el-button>
        </div>
        
        <div class="version-description">
          <el-alert
            :title="getVersionDescription(selectedVersion)"
            type="info"
            :closable="false"
          />
        </div>
      </div>
      
      <div class="generate-section">
        <el-row :gutter="20">
          <el-col :span="12">
            <div class="count-section">
              <el-form-item label="生成数量">
                <el-input-number
                  v-model="count"
                  :min="1"
                  :max="100"
                  controls-position="right"
                />
              </el-form-item>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="format-section">
              <el-form-item label="显示格式">
                <el-select v-model="format" placeholder="选择格式">
                  <el-option label="标准格式" value="standard"></el-option>
                  <el-option label="无连字符" value="no-dashes"></el-option>
                  <el-option label="大写" value="uppercase"></el-option>
                </el-select>
              </el-form-item>
            </div>
          </el-col>
        </el-row>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="generateUUIDs">
          <i class="fas fa-id-card"></i> 生成UUID
        </el-button>
        <el-button @click="clearResults">
          <i class="fas fa-trash"></i> 清空结果
        </el-button>
        <el-button type="success" @click="copyAllResults">
          <i class="fas fa-copy"></i> 复制所有结果
        </el-button>
      </div>
      
      <div class="results-section" v-if="results.length > 0">
        <h3>生成结果</h3>
        <div class="results-list">
          <div
            v-for="(result, index) in results"
            :key="index"
            class="result-item"
          >
            <el-input
              v-model="results[index]"
              readonly
            />
            <el-button
              type="success"
              icon="fas fa-copy"
              @click="copySingleResult(index)"
              circle
            />
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
import { ref } from 'vue'
import { ElCard, ElRow, ElCol, ElFormItem, ElInputNumber, ElSelect, ElOption, ElButton, ElAlert } from 'element-plus'
import { v1 as uuidv1, v4 as uuidv4, v5 as uuidv5 } from 'uuid'

// 数据模型
const selectedVersion = ref('v4')
const count = ref(1)
const format = ref('standard')
const results = ref([])
const validationResult = ref(null)

// UUID版本列表
const uuidVersions = ref([
  { value: 'v1', name: 'UUID v1' },
  { value: 'v4', name: 'UUID v4' },
  { value: 'v5', name: 'UUID v5' }
])

// 选择UUID版本
const selectVersion = (version) => {
  selectedVersion.value = version
}

// 获取版本描述
const getVersionDescription = (version) => {
  const descriptions = {
    v1: '基于时间和节点ID的UUID，包含时间戳信息',
    v4: '基于随机数的UUID，最常用的版本',
    v5: '基于命名空间和名称的UUID，使用SHA-1哈希'
  }
  return descriptions[version] || '未知版本'
}

// 生成UUID
const generateUUIDs = () => {
  try {
    const generatedResults = []
    
    for (let i = 0; i < count.value; i++) {
      let uuid = ''
      
      switch (selectedVersion.value) {
        case 'v1':
          uuid = generateUUIDv1()
          break
        case 'v4':
          uuid = generateUUIDv4()
          break
        case 'v5':
          uuid = generateUUIDv5()
          break
        default:
          uuid = generateUUIDv4()
      }
      
      // 格式化UUID
      switch (format.value) {
        case 'no-dashes':
          uuid = uuid.replace(/-/g, '')
          break
        case 'uppercase':
          uuid = uuid.toUpperCase()
          break
        default:
          // 标准格式，无需处理
          break
      }
      
      generatedResults.push(uuid)
    }
    
    results.value = generatedResults
    validationResult.value = {
      type: 'success',
      message: `成功生成${count.value}个UUID！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `UUID生成失败: ${error.message}`
    }
  }
}

// 生成UUID v1
const generateUUIDv1 = () => {
  return uuidv1()
}

// 生成UUID v4
const generateUUIDv4 = () => {
  return uuidv4()
}

// 生成UUID v5
const generateUUIDv5 = () => {
  // 使用一个固定的命名空间UUID和示例名称
  // 在实际应用中，用户可能需要提供自己的命名空间和名称
  const namespace = '6ba7b810-9dad-11d1-80b4-00c04fd430c8' // 示例命名空间
  const name = 'example' // 示例名称
  return uuidv5(name, namespace)
}

// 清空结果
const clearResults = () => {
  results.value = []
  validationResult.value = null
}

// 复制单个结果
const copySingleResult = async (index) => {
  try {
    await navigator.clipboard.writeText(results.value[index])
    validationResult.value = {
      type: 'success',
      message: 'UUID已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 复制所有结果
const copyAllResults = async () => {
  if (results.value.length === 0) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    const allResults = results.value.join('\n')
    await navigator.clipboard.writeText(allResults)
    validationResult.value = {
      type: 'success',
      message: '所有UUID已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 初始化示例数据
count.value = 5
</script>

<style scoped>
.tool-container {
  margin-bottom: 2rem;
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
}

.uuid-version-section {
  width: 100%;
}

.uuid-version-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.version-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.version-buttons .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
}

.version-description {
  width: 100%;
}

.generate-section {
  width: 100%;
}

.count-section,
.format-section {
  width: 100%;
}

.button-group {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.button-group .el-button {
  flex: 1;
  min-width: 120px;
  max-width: 160px;
}

.results-section {
  width: 100%;
}

.results-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.result-item {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.result-item .el-input {
  flex: 1;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .version-buttons {
    flex-direction: column;
  }
  
  .version-buttons .el-button {
    width: 100%;
    max-width: none;
  }
  
  .button-group {
    flex-direction: column;
    align-items: center;
  }
  
  .button-group .el-button {
    width: 100%;
    max-width: none;
  }
  
  .result-item {
    flex-direction: column;
  }
  
  .result-item .el-button {
    align-self: flex-end;
  }
}
</style>