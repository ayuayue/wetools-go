export const menuItems = [
  {
    id: 'encoding',
    title: '编码解码',
    icon: 'fas fa-code',
    collapsed: false,
    items: [
      {
        id: 'json',
        title: 'JSON工具',
        type: 'json',
        icon: 'fas fa-file-code'
      },
      {
        id: 'xml',
        title: 'XML工具',
        type: 'xml',
        icon: 'fas fa-file-code'
      },
      {
        id: 'html',
        title: 'HTML工具',
        type: 'html',
        icon: 'fab fa-html5'
      },
      {
        id: 'base64',
        title: 'Base64编码',
        type: 'base64',
        icon: 'fas fa-file-code'
      },
      {
        id: 'url',
        title: 'URL编解码',
        type: 'url',
        icon: 'fas fa-link'
      },
      {
        id: 'jwt',
        title: 'JWT工具',
        type: 'jwt',
        icon: 'fas fa-key'
      }
    ]
  },
  {
    id: 'encryption',
    title: '加密解密',
    icon: 'fas fa-lock',
    collapsed: false,
    items: [
      {
        id: 'hash',
        title: '哈希计算',
        type: 'hash',
        icon: 'fas fa-hashtag'
      },
      {
        id: 'encrypt',
        title: '加密解密',
        type: 'encrypt',
        icon: 'fas fa-key'
      }
    ]
  },
  {
    id: 'random',
    title: '随机工具',
    icon: 'fas fa-random',
    collapsed: false,
    items: [
      {
        id: 'random-string',
        title: '随机字符串',
        type: 'random',
        icon: 'fas fa-random'
      },
      {
        id: 'uuid',
        title: 'UUID生成',
        type: 'uuid',
        icon: 'fas fa-id-card'
      }
    ]
  },
  {
    id: 'conversion',
    title: '格式转换',
    icon: 'fas fa-exchange-alt',
    collapsed: false,
    items: [
      {
        id: 'converter',
        title: '格式转换器',
        type: 'converter',
        icon: 'fas fa-exchange-alt'
      },
      {
        id: 'qr',
        title: '二维码生成',
        type: 'qr',
        icon: 'fas fa-qrcode'
      }
    ]
  },
  {
    id: 'websites',
    title: '网站集成',
    icon: 'fas fa-globe',
    collapsed: false,
    items: [
      {
        id: 'website-manager',
        title: '网站管理',
        type: 'website',
        icon: 'fas fa-cog'
      }
    ]
  },
  {
    id: 'code-execution',
    title: '代码运行',
    icon: 'fas fa-code',
    collapsed: false,
    items: [
      {
        id: 'html-runner',
        title: 'HTML运行器',
        type: 'htmlrunner',
        icon: 'fab fa-html5'
      },
      {
        id: 'js-runner',
        title: 'JavaScript运行器',
        type: 'jsrunner',
        icon: 'fab fa-js'
      }
    ]
  },
  {
    id: 'system',
    title: '系统工具',
    icon: 'fas fa-tools',
    collapsed: false,
    items: [
      {
        id: 'clipboard-manager',
        title: '剪贴板管理',
        type: 'clipboard',
        icon: 'fas fa-clipboard'
      }
    ]
  },
  {
    id: 'settings',
    title: '系统设置',
    icon: 'fas fa-cog',
    collapsed: false,
    items: [
      {
        id: 'proxy-config',
        title: '代理设置',
        type: 'proxy',
        icon: 'fas fa-network-wired'
      },
      {
        id: 'app-settings',
        title: '应用设置',
        type: 'settings',
        icon: 'fas fa-sliders-h'
      },
      {
        id: 'about',
        title: '关于',
        type: 'about',
        icon: 'fas fa-info-circle'
      }
    ]
  }
]
