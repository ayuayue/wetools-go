// proxyManager.js
// 全局代理配置管理器

class ProxyManager {
  constructor() {
    this.proxyConfig = null
    this.init()
  }

  // 初始化代理管理器
  init() {
    this.loadProxyConfig()
    this.setupGlobalProxy()
  }

  // 加载代理配置
  loadProxyConfig() {
    try {
      const savedConfig = localStorage.getItem('proxyConfig')
      if (savedConfig) {
        this.proxyConfig = JSON.parse(savedConfig)
      }
    } catch (error) {
      console.error('加载代理配置失败:', error)
    }
  }

  // 设置全局代理
  setupGlobalProxy() {
    if (!this.proxyConfig || !this.proxyConfig.enabled) {
      return
    }

    // 注意：在浏览器环境中，我们无法直接设置全局代理
    // 这里只是保存配置供其他组件使用
    console.log('代理配置已加载:', this.proxyConfig)
  }

  // 获取代理URL
  getProxyURL() {
    if (!this.proxyConfig || !this.proxyConfig.enabled || !this.proxyConfig.host) {
      return null
    }

    let proxyURL = ''
    if (this.proxyConfig.username && this.proxyConfig.password) {
      proxyURL = `http://${encodeURIComponent(this.proxyConfig.username)}:${encodeURIComponent(this.proxyConfig.password)}@${this.proxyConfig.host}:${this.proxyConfig.port}`
    } else {
      proxyURL = `http://${this.proxyConfig.host}:${this.proxyConfig.port}`
    }

    return proxyURL
  }

  // 检查URL是否应该通过代理
  shouldProxy(url) {
    if (!this.proxyConfig || !this.proxyConfig.enabled) {
      return false
    }

    // 简单检查，实际应用中可能需要更复杂的规则
    try {
      const urlObj = new URL(url)
      // 可以添加不使用代理的域名列表
      const bypassList = ['localhost', '127.0.0.1']
      return !bypassList.includes(urlObj.hostname)
    } catch (error) {
      console.error('解析URL失败:', error)
      return false
    }
  }

  // 通过代理获取资源
  async fetchWithProxy(url, options = {}) {
    if (!this.shouldProxy(url)) {
      return fetch(url, options)
    }

    const proxyURL = this.getProxyURL()
    if (!proxyURL) {
      return fetch(url, options)
    }

    // 构造代理请求URL
    // 注意：这需要后端代理服务器支持
    const proxyRequestURL = `${proxyURL}/${url}`

    try {
      return await fetch(proxyRequestURL, options)
    } catch (error) {
      console.error('通过代理请求失败:', error)
      // 如果代理请求失败，尝试直接请求
      return fetch(url, options)
    }
  }

  // 更新代理配置
  updateProxyConfig(config) {
    this.proxyConfig = config
    try {
      localStorage.setItem('proxyConfig', JSON.stringify(config))
      this.setupGlobalProxy()
    } catch (error) {
      console.error('保存代理配置失败:', error)
    }
  }
}

// 创建单例实例
const proxyManager = new ProxyManager()

// 导出单例实例
export default proxyManager