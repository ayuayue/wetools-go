<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>关于 WeTools</h2>
      <p>开发者工具箱 - 提升开发效率的利器</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="about-section">
            <h3>应用介绍</h3>
            <p>WeTools 是一款专为开发者设计的桌面工具箱应用，集成了多种常用的开发工具，帮助开发者提高工作效率。</p>
            <p>本应用基于 Wails 框架开发，结合了 Go 语言的高性能后端和 Vue.js 的现代化前端界面，提供了流畅的用户体验。</p>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="about-section">
            <h3>功能特性</h3>
            <ul class="feature-list">
              <li><i class="fas fa-check-circle"></i> JSON/XML/HTML 格式化工具</li>
              <li><i class="fas fa-check-circle"></i> Base64/URL 编解码工具</li>
              <li><i class="fas fa-check-circle"></i> JWT 编解码和验证工具</li>
              <li><i class="fas fa-check-circle"></i> 哈希计算和加密解密工具</li>
              <li><i class="fas fa-check-circle"></i> 随机字符串和 UUID 生成器</li>
              <li><i class="fas fa-check-circle"></i> 格式转换器和二维码生成器</li>
              <li><i class="fas fa-check-circle"></i> 网站集成管理</li>
              <li><i class="fas fa-check-circle"></i> 代码运行器（HTML/JavaScript）</li>
            </ul>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="about-section">
            <h3>技术栈</h3>
            <ul class="tech-list">
              <li><strong>后端:</strong> Go + Wails v2</li>
              <li><strong>前端:</strong> Vue 3 + Element Plus</li>
              <li><strong>构建工具:</strong> Vite</li>
              <li><strong>包管理:</strong> npm</li>
            </ul>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="about-section">
            <h3>开源信息</h3>
            <ul class="info-list">
              <li><strong>版本:</strong> v1.0.0</li>
              <li><strong>许可证:</strong> MIT</li>
              <li><strong>仓库:</strong> <a href="https://github.com/your-username/wetools-go" target="_blank">GitHub</a></li>
              <li><strong>反馈:</strong> <a href="https://github.com/your-username/wetools-go/issues" target="_blank">问题反馈</a></li>
            </ul>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="about-section">
            <h3>使用说明</h3>
            <ol class="usage-list">
              <li>1. 从左侧菜单选择需要的工具</li>
              <li>2. 在工具页面进行相应操作</li>
              <li>3. 大部分工具都支持一键复制结果</li>
              <li>4. 网站集成工具可以添加常用网站并在应用内浏览</li>
            </ol>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="about-section">
            <h3>联系我们</h3>
            <p>如果您有任何问题或建议，欢迎通过以下方式联系我们：</p>
            <ul class="contact-list">
              <li><i class="fas fa-envelope"></i> 邮箱: <a href="mailto:support@wetools.com">support@wetools.com</a></li>
              <li><i class="fab fa-github"></i> GitHub: <a href="https://github.com/your-username/wetools-go" target="_blank">WeTools 项目</a></li>
            </ul>
          </div>
        </el-col>
      </el-row>
    </div>
  </el-card>
</template>

<script setup>
import { ElCard, ElRow, ElCol } from 'element-plus'
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

.about-section {
  padding: 1rem;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.about-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
  border-bottom: 1px solid #e1e5e9;
  padding-bottom: 0.5rem;
}

.about-section p {
  color: #555;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.feature-list,
.tech-list,
.info-list,
.usage-list,
.contact-list {
  list-style: none;
  padding: 0;
}

.feature-list li,
.tech-list li,
.info-list li,
.usage-list li,
.contact-list li {
  margin-bottom: 0.5rem;
  color: #555;
}

.feature-list li i {
  color: #67c23a;
  margin-right: 0.5rem;
}

.contact-list a,
.info-list a {
  color: #409eff;
  text-decoration: none;
}

.contact-list a:hover,
.info-list a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .tool-content {
    gap: 1rem;
  }
  
  .about-section {
    padding: 0.75rem;
  }
}
</style>