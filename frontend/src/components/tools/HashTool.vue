<template>
  <el-card class="tool-container" shadow="never">
    <div class="tool-header">
      <h2>哈希计算工具</h2>
      <p>MD5、SHA1、SHA256等哈希算法计算工具</p>
    </div>
    
    <div class="tool-content">
      <el-row :gutter="20">
        <el-col :span="24">
          <div class="input-section">
            <el-input
              v-model="inputData"
              type="textarea"
              :rows="8"
              placeholder="请输入要计算哈希值的文本"
              resize="vertical"
            />
          </div>
        </el-col>
      </el-row>
      
      <div class="hash-algorithms">
        <h3>选择哈希算法</h3>
        <div class="algorithm-buttons">
          <el-button 
            v-for="algorithm in algorithms" 
            :key="algorithm.name"
            :type="selectedAlgorithm === algorithm.name ? 'primary' : ''"
            @click="selectAlgorithm(algorithm.name)"
          >
            {{ algorithm.name }}
          </el-button>
        </div>
      </div>
      
      <div class="button-group">
        <el-button type="primary" @click="calculateHash">
          <i class="fas fa-calculator"></i> 计算哈希值
        </el-button>
        <el-button @click="calculateAllHashes">
          <i class="fas fa-calculator"></i> 计算所有哈希值
        </el-button>
        <el-button @click="clearData">
          <i class="fas fa-trash"></i> 清空
        </el-button>
        <el-button type="success" @click="copyResult">
          <i class="fas fa-copy"></i> 复制结果
        </el-button>
      </div>
      
      <div class="result-section" v-if="hashResult">
        <h3>计算结果</h3>
        <el-input
          v-model="hashResult"
          readonly
          placeholder="哈希计算结果将显示在这里"
        />
      </div>
      
      <div class="all-results-section" v-if="allHashResults.length > 0">
        <h3>所有哈希值计算结果</h3>
        <el-table :data="allHashResults" style="width: 100%" border>
          <el-table-column prop="algorithm" label="算法" width="120"></el-table-column>
          <el-table-column prop="hash" label="哈希值"></el-table-column>
        </el-table>
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
import { ElCard, ElRow, ElCol, ElInput, ElButton, ElAlert, ElTable, ElTableColumn } from 'element-plus'

// 数据模型
const inputData = ref('')
const selectedAlgorithm = ref('MD5')
const hashResult = ref('')
const allHashResults = ref([])
const validationResult = ref(null)

// 哈希算法列表
const algorithms = ref([
  { name: 'MD5', fullName: 'Message-Digest Algorithm 5' },
  { name: 'SHA1', fullName: 'Secure Hash Algorithm 1' },
  { name: 'SHA256', fullName: 'Secure Hash Algorithm 256' },
  { name: 'SHA512', fullName: 'Secure Hash Algorithm 512' }
])

// 选择算法
const selectAlgorithm = (algorithm) => {
  selectedAlgorithm.value = algorithm
}

// 计算哈希值
const calculateHash = async () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    // 使用Web Crypto API计算哈希值
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    let hashBuffer
    
    switch (selectedAlgorithm.value) {
      case 'MD5':
        // MD5需要使用第三方库或自定义实现，这里简化处理
        hashResult.value = calculateMD5(inputData.value)
        break
      case 'SHA1':
        hashBuffer = await crypto.subtle.digest('SHA-1', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      case 'SHA256':
        hashBuffer = await crypto.subtle.digest('SHA-256', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      case 'SHA512':
        hashBuffer = await crypto.subtle.digest('SHA-512', data)
        hashResult.value = Array.from(new Uint8Array(hashBuffer))
          .map(b => b.toString(16).padStart(2, '0'))
          .join('').toUpperCase()
        break
      default:
        throw new Error('不支持的哈希算法')
    }
    
    allHashResults.value = []
    validationResult.value = {
      type: 'success',
      message: `${selectedAlgorithm.value}哈希值计算成功！`
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `哈希计算失败: ${error.message}`
    }
  }
}

// 计算所有哈希值
const calculateAllHashes = async () => {
  try {
    if (!inputData.value.trim()) {
      validationResult.value = {
        type: 'warning',
        message: '请输入要计算哈希值的文本'
      }
      return
    }
    
    const results = []
    const encoder = new TextEncoder()
    const data = encoder.encode(inputData.value)
    
    // 计算SHA1
    const sha1Buffer = await crypto.subtle.digest('SHA-1', data)
    const sha1Hash = Array.from(new Uint8Array(sha1Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA1', hash: sha1Hash })
    
    // 计算SHA256
    const sha256Buffer = await crypto.subtle.digest('SHA-256', data)
    const sha256Hash = Array.from(new Uint8Array(sha256Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA256', hash: sha256Hash })
    
    // 计算SHA512
    const sha512Buffer = await crypto.subtle.digest('SHA-512', data)
    const sha512Hash = Array.from(new Uint8Array(sha512Buffer))
      .map(b => b.toString(16).padStart(2, '0'))
      .join('').toUpperCase()
    results.push({ algorithm: 'SHA512', hash: sha512Hash })
    
    // 计算MD5 (简化实现)
    const md5Hash = calculateMD5(inputData.value)
    results.push({ algorithm: 'MD5', hash: md5Hash })
    
    allHashResults.value = results
    hashResult.value = ''
    validationResult.value = {
      type: 'success',
      message: '所有哈希值计算成功！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: `哈希计算失败: ${error.message}`
    }
  }
}

// 简化的MD5计算函数
const calculateMD5 = async (str) => {
  // 使用简单的MD5实现
  const md5 = (function() {
    // MD5实现代码
    function md5_RotateLeft(lValue, iShiftBits) {
      return (lValue << iShiftBits) | (lValue >>> (32 - iShiftBits));
    }

    function md5_AddUnsigned(lX, lY) {
      var lX4, lY4, lX8, lY8, lResult;
      lX8 = (lX & 0x80000000);
      lY8 = (lY & 0x80000000);
      lX4 = (lX & 0x40000000);
      lY4 = (lY & 0x40000000);
      lResult = (lX & 0x3FFFFFFF) + (lY & 0x3FFFFFFF);
      if (lX4 & lY4) {
        return (lResult ^ 0x80000000 ^ lX8 ^ lY8);
      }
      if (lX4 | lY4) {
        if (lResult & 0x40000000) {
          return (lResult ^ 0xC0000000 ^ lX8 ^ lY8);
        } else {
          return (lResult ^ 0x40000000 ^ lX8 ^ lY8);
        }
      } else {
        return (lResult ^ lX8 ^ lY8);
      }
    }

    function md5_F(x, y, z) { return (x & y) | ((~x) & z); }
    function md5_G(x, y, z) { return (x & z) | (y & (~z)); }
    function md5_H(x, y, z) { return (x ^ y ^ z); }
    function md5_I(x, y, z) { return (y ^ (x | (~z))); }

    function md5_FF(a, b, c, d, x, s, ac) {
      a = md5_AddUnsigned(a, md5_AddUnsigned(md5_AddUnsigned(md5_F(b, c, d), x), ac));
      return md5_AddUnsigned(md5_RotateLeft(a, s), b);
    }

    function md5_GG(a, b, c, d, x, s, ac) {
      a = md5_AddUnsigned(a, md5_AddUnsigned(md5_AddUnsigned(md5_G(b, c, d), x), ac));
      return md5_AddUnsigned(md5_RotateLeft(a, s), b);
    }

    function md5_HH(a, b, c, d, x, s, ac) {
      a = md5_AddUnsigned(a, md5_AddUnsigned(md5_AddUnsigned(md5_H(b, c, d), x), ac));
      return md5_AddUnsigned(md5_RotateLeft(a, s), b);
    }

    function md5_II(a, b, c, d, x, s, ac) {
      a = md5_AddUnsigned(a, md5_AddUnsigned(md5_AddUnsigned(md5_I(b, c, d), x), ac));
      return md5_AddUnsigned(md5_RotateLeft(a, s), b);
    }

    function md5_ConvertToWordArray(string) {
      var lWordCount;
      var lMessageLength = string.length;
      var lNumberOfWords_temp1 = lMessageLength + 8;
      var lNumberOfWords_temp2 = (lNumberOfWords_temp1 - (lNumberOfWords_temp1 % 64)) / 64;
      var lNumberOfWords = (lNumberOfWords_temp2 + 1) * 16;
      var lWordArray = Array(lNumberOfWords - 1);
      var lBytePosition = 0;
      var lByteCount = 0;
      while (lByteCount < lMessageLength) {
        lWordCount = (lByteCount - (lByteCount % 4)) / 4;
        lBytePosition = (lByteCount % 4) * 8;
        lWordArray[lWordCount] = (lWordArray[lWordCount] | (string.charCodeAt(lByteCount) << lBytePosition));
        lByteCount++;
      }
      lWordCount = (lByteCount - (lByteCount % 4)) / 4;
      lBytePosition = (lByteCount % 4) * 8;
      lWordArray[lWordCount] = lWordArray[lWordCount] | (0x80 << lBytePosition);
      lWordArray[lNumberOfWords - 2] = lMessageLength << 3;
      lWordArray[lNumberOfWords - 1] = lMessageLength >>> 29;
      return lWordArray;
    }

    function md5_WordToHex(lValue) {
      var WordToHexValue = "", WordToHexValue_temp = "", lByte, lCount;
      for (lCount = 0; lCount <= 3; lCount++) {
        lByte = (lValue >>> (lCount * 8)) & 255;
        WordToHexValue_temp = "0" + lByte.toString(16);
        WordToHexValue = WordToHexValue + WordToHexValue_temp.substr(WordToHexValue_temp.length - 2, 2);
      }
      return WordToHexValue;
    }

    function md5_Utf8Encode(string) {
      string = string.replace(/\r\n/g, "\n");
      var utftext = "";
      for (var n = 0; n < string.length; n++) {
        var c = string.charCodeAt(n);
        if (c < 128) {
          utftext += String.fromCharCode(c);
        } else if ((c > 127) && (c < 2048)) {
          utftext += String.fromCharCode((c >> 6) | 192);
          utftext += String.fromCharCode((c & 63) | 128);
        } else {
          utftext += String.fromCharCode((c >> 12) | 224);
          utftext += String.fromCharCode(((c >> 6) & 63) | 128);
          utftext += String.fromCharCode((c & 63) | 128);
        }
      }
      return utftext;
    }

    var x = Array();
    var k, AA, BB, CC, DD, a, b, c, d;
    var S11 = 7, S12 = 12, S13 = 17, S14 = 22;
    var S21 = 5, S22 = 9, S23 = 14, S24 = 20;
    var S31 = 4, S32 = 11, S33 = 16, S34 = 23;
    var S41 = 6, S42 = 10, S43 = 15, S44 = 21;

    string = md5_Utf8Encode(str);
    x = md5_ConvertToWordArray(string);
    a = 0x67452301; b = 0xEFCDAB89; c = 0x98BADCFE; d = 0x10325476;

    for (k = 0; k < x.length; k += 16) {
      AA = a; BB = b; CC = c; DD = d;
      a = md5_FF(a, b, c, d, x[k + 0], S11, 0xD76AA478);
      d = md5_FF(d, a, b, c, x[k + 1], S12, 0xE8C7B756);
      c = md5_FF(c, d, a, b, x[k + 2], S13, 0x242070DB);
      b = md5_FF(b, c, d, a, x[k + 3], S14, 0xC1BDCEEE);
      a = md5_FF(a, b, c, d, x[k + 4], S11, 0xF57C0FAF);
      d = md5_FF(d, a, b, c, x[k + 5], S12, 0x4787C62A);
      c = md5_FF(c, d, a, b, x[k + 6], S13, 0xA8304613);
      b = md5_FF(b, c, d, a, x[k + 7], S14, 0xFD469501);
      a = md5_FF(a, b, c, d, x[k + 8], S11, 0x698098D8);
      d = md5_FF(d, a, b, c, x[k + 9], S12, 0x8B44F7AF);
      c = md5_FF(c, d, a, b, x[k + 10], S13, 0xFFFF5BB1);
      b = md5_FF(b, c, d, a, x[k + 11], S14, 0x895CD7BE);
      a = md5_FF(a, b, c, d, x[k + 12], S11, 0x6B901122);
      d = md5_FF(d, a, b, c, x[k + 13], S12, 0xFD987193);
      c = md5_FF(c, d, a, b, x[k + 14], S13, 0xA679438E);
      b = md5_FF(b, c, d, a, x[k + 15], S14, 0x49B40821);
      a = md5_GG(a, b, c, d, x[k + 1], S21, 0xF61E2562);
      d = md5_GG(d, a, b, c, x[k + 6], S22, 0xC040B340);
      c = md5_GG(c, d, a, b, x[k + 11], S23, 0x265E5A51);
      b = md5_GG(b, c, d, a, x[k + 0], S24, 0xE9B6C7AA);
      a = md5_GG(a, b, c, d, x[k + 5], S21, 0xD62F105D);
      d = md5_GG(d, a, b, c, x[k + 10], S22, 0x2441453);
      c = md5_GG(c, d, a, b, x[k + 15], S23, 0xD8A1E681);
      b = md5_GG(b, c, d, a, x[k + 4], S24, 0xE7D3FBC8);
      a = md5_GG(a, b, c, d, x[k + 9], S21, 0x21E1CDE6);
      d = md5_GG(d, a, b, c, x[k + 14], S22, 0xC33707D6);
      c = md5_GG(c, d, a, b, x[k + 3], S23, 0xF4D50D87);
      b = md5_GG(b, c, d, a, x[k + 8], S24, 0x455A14ED);
      a = md5_GG(a, b, c, d, x[k + 13], S21, 0xA9E3E905);
      d = md5_GG(d, a, b, c, x[k + 2], S22, 0xFCEFA3F8);
      c = md5_GG(c, d, a, b, x[k + 7], S23, 0x676F02D9);
      b = md5_GG(b, c, d, a, x[k + 12], S24, 0x8D2A4C8A);
      a = md5_HH(a, b, c, d, x[k + 5], S31, 0xFFFA3942);
      d = md5_HH(d, a, b, c, x[k + 8], S32, 0x8771F681);
      c = md5_HH(c, d, a, b, x[k + 11], S33, 0x6D9D6122);
      b = md5_HH(b, c, d, a, x[k + 14], S34, 0xFDE5380C);
      a = md5_HH(a, b, c, d, x[k + 1], S31, 0xA4BEEA44);
      d = md5_HH(d, a, b, c, x[k + 4], S32, 0x4BDECFA9);
      c = md5_HH(c, d, a, b, x[k + 7], S33, 0xF6BB4B60);
      b = md5_HH(b, c, d, a, x[k + 10], S34, 0xBEBFBC70);
      a = md5_HH(a, b, c, d, x[k + 13], S31, 0x289B7EC6);
      d = md5_HH(d, a, b, c, x[k + 0], S32, 0xEAA127FA);
      c = md5_HH(c, d, a, b, x[k + 3], S33, 0xD4EF3085);
      b = md5_HH(b, c, d, a, x[k + 6], S34, 0x4881D05);
      a = md5_HH(a, b, c, d, x[k + 9], S31, 0xD9D4D039);
      d = md5_HH(d, a, b, c, x[k + 12], S32, 0xE6DB99E5);
      c = md5_HH(c, d, a, b, x[k + 15], S33, 0x1FA27CF8);
      b = md5_HH(b, c, d, a, x[k + 2], S34, 0xC4AC5665);
      a = md5_II(a, b, c, d, x[k + 0], S41, 0xF4292244);
      d = md5_II(d, a, b, c, x[k + 7], S42, 0x432AFF97);
      c = md5_II(c, d, a, b, x[k + 14], S43, 0xAB9423A7);
      b = md5_II(b, c, d, a, x[k + 5], S44, 0xFC93A039);
      a = md5_II(a, b, c, d, x[k + 12], S41, 0x655B59C3);
      d = md5_II(d, a, b, c, x[k + 3], S42, 0x8F0CCC92);
      c = md5_II(c, d, a, b, x[k + 10], S43, 0xFFEFF47D);
      b = md5_II(b, c, d, a, x[k + 1], S44, 0x85845DD1);
      a = md5_II(a, b, c, d, x[k + 8], S41, 0x6FA87E4F);
      d = md5_II(d, a, b, c, x[k + 15], S42, 0xFE2CE6E0);
      c = md5_II(c, d, a, b, x[k + 6], S43, 0xA3014314);
      b = md5_II(b, c, d, a, x[k + 13], S44, 0x4E0811A1);
      a = md5_II(a, b, c, d, x[k + 4], S41, 0xF7537E82);
      d = md5_II(d, a, b, c, x[k + 11], S42, 0xBD3AF235);
      c = md5_II(c, d, a, b, x[k + 2], S43, 0x2AD7D2BB);
      b = md5_II(b, c, d, a, x[k + 9], S44, 0xEB86D391);
      a = md5_AddUnsigned(a, AA);
      b = md5_AddUnsigned(b, BB);
      c = md5_AddUnsigned(c, CC);
      d = md5_AddUnsigned(d, DD);
    }
    return (md5_WordToHex(a) + md5_WordToHex(b) + md5_WordToHex(c) + md5_WordToHex(d)).toUpperCase();
  })();
  return md5;
}

// 清空数据
const clearData = () => {
  inputData.value = ''
  hashResult.value = ''
  allHashResults.value = []
  validationResult.value = null
}

// 复制结果
const copyResult = async () => {
  const result = hashResult.value || (allHashResults.value.length > 0 ? 
    allHashResults.value.map(r => `${r.algorithm}: ${r.hash}`).join('\n') : '')
  
  if (!result) {
    validationResult.value = {
      type: 'warning',
      message: '没有内容可复制'
    }
    return
  }
  
  try {
    await navigator.clipboard.writeText(result)
    validationResult.value = {
      type: 'success',
      message: '结果已复制到剪贴板！'
    }
  } catch (error) {
    validationResult.value = {
      type: 'error',
      message: '复制失败，请手动复制'
    }
  }
}

// 初始化示例数据
inputData.value = 'Hello, WeTools!'
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

.input-section {
  width: 100%;
}

.hash-algorithms {
  width: 100%;
}

.hash-algorithms h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.algorithm-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.algorithm-buttons .el-button {
  flex: 1;
  min-width: 100px;
  max-width: 150px;
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

.result-section,
.all-results-section {
  width: 100%;
}

.result-section h3,
.all-results-section h3 {
  margin: 0 0 1rem 0;
  color: #333;
}

.validation-result {
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .algorithm-buttons {
    flex-direction: column;
  }
  
  .algorithm-buttons .el-button {
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
}
</style>