// Element Plus全局配置
import { ref, computed } from 'vue'

// 主题状态
const isDarkTheme = ref(false)

// 计算主题
const theme = computed(() => isDarkTheme.value ? 'dark' : 'light')

// 切换主题
const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value
  // 保存主题设置到本地存储
  localStorage.setItem('theme', isDarkTheme.value ? 'dark' : 'light')
  
  // 应用主题到body
  if (isDarkTheme.value) {
    document.body.classList.add('dark-theme')
  } else {
    document.body.classList.remove('dark-theme')
  }
}

// 初始化主题
const initTheme = () => {
  // 加载主题设置
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDarkTheme.value = savedTheme === 'dark'
    if (isDarkTheme.value) {
      document.body.classList.add('dark-theme')
    }
  }
}

// Element Plus配置
const elementPlusConfig = {
  // 组件默认属性
  defaultProps: {
    button: {
      size: 'default',
      type: 'default'
    },
    input: {
      clearable: true,
      placeholder: '请输入内容'
    },
    textarea: {
      clearable: true,
      placeholder: '请输入内容',
      autosize: { minRows: 3, maxRows: 10 }
    },
    card: {
      shadow: 'never'
    },
    table: {
      border: true,
      stripe: true
    },
    pagination: {
      layout: 'prev, pager, next, jumper, ->, total'
    }
  },
  
  // 国际化配置
  locale: {
    lang: 'zh-cn',
    // 可以在这里添加更多国际化配置
  }
}

export {
  isDarkTheme,
  theme,
  toggleTheme,
  initTheme,
  elementPlusConfig
}