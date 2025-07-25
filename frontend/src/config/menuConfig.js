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
        icon: 'fas fa-file-code',
        description: '在线JSON格式化、校验、压缩工具，支持JSON转XML、CSV等格式'
      },
      {
        id: 'xml',
        title: 'XML工具',
        type: 'xml',
        icon: 'fas fa-file-code',
        description: 'XML格式化、校验、压缩工具，支持XML与其他格式互转'
      },
      {
        id: 'html',
        title: 'HTML工具',
        type: 'html',
        icon: 'fab fa-html5',
        description: 'HTML格式化、清理、压缩工具，支持HTML实体编解码'
      },
      {
        id: 'base64',
        title: 'Base64编码',
        type: 'base64',
        icon: 'fas fa-file-code',
        description: 'Base64编码解码工具，支持文件Base64转换'
      },
      {
        id: 'url',
        title: 'URL编解码',
        type: 'url',
        icon: 'fas fa-link',
        description: 'URL编码解码工具，支持URL参数解析'
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
        icon: 'fas fa-hashtag',
        description: 'MD5、SHA1、SHA256等哈希算法计算工具'
      },
      {
        id: 'encrypt',
        title: '加密解密',
        type: 'encrypt',
        icon: 'fas fa-key',
        description: 'AES、DES等对称加密算法工具'
      },
      {
        id: 'jwt',
        title: 'JWT编解码',
        type: 'jwt',
        icon: 'fas fa-key',
        description: 'JWT Token编解码工具，支持Header、Payload、Signature查看'
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
        icon: 'fas fa-random',
        description: '自定义长度和字符集的随机字符串生成器'
      },
      {
        id: 'uuid',
        title: 'UUID生成',
        type: 'uuid',
        icon: 'fas fa-id-card',
        description: 'UUID生成器，支持多种UUID版本'
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
        icon: 'fas fa-exchange-alt',
        description: '多种数据格式互转工具'
      },
      {
        id: 'qr',
        title: '二维码生成',
        type: 'qr',
        icon: 'fas fa-qrcode',
        description: '二维码生成与解析工具'
      },
      {
        id: 'timestamp',
        title: '时间戳转换',
        type: 'timestamp',
        icon: 'fas fa-clock',
        description: '时间戳与标准时间互转工具，支持多种时间格式'
      },
      {
        id: 'text',
        title: '文本处理',
        type: 'text',
        icon: 'fas fa-font',
        description: '文本格式化、大小写转换、繁简转换等工具'
      }
    ]
  }
]