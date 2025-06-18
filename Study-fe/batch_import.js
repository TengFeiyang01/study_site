#!/usr/bin/env node

/**
 * 八股题批量导入脚本
 * 
 * 使用方法：
 * node batch_import.js <json文件路径> [服务器地址]
 * 
 * 示例：
 * node batch_import.js example_questions.json
 * node batch_import.js ./data/questions.json http://localhost:8080
 */

const fs = require('fs')
const path = require('path')
const axios = require('axios')

// 默认配置
const DEFAULT_SERVER = 'http://localhost:8080'
const API_ENDPOINT = '/question/'

// 延迟函数
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

// 彩色输出
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  dim: '\x1b[2m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m'
}

function colorLog(message, color = 'reset') {
  console.log(colors[color] + message + colors.reset)
}

function showProgress(current, total, message = '') {
  const percentage = Math.round((current / total) * 100)
  const barLength = 30
  const filledLength = Math.round((current / total) * barLength)
  const bar = '█'.repeat(filledLength) + '░'.repeat(barLength - filledLength)
  
  process.stdout.write(`\r${colors.cyan}[${bar}] ${percentage}% (${current}/${total}) ${message}${colors.reset}`)
}

async function batchImport(jsonFilePath, serverUrl = DEFAULT_SERVER) {
  try {
    // 读取JSON文件
    colorLog('\n📁 读取JSON文件...', 'blue')
    
    if (!fs.existsSync(jsonFilePath)) {
      throw new Error(`文件不存在: ${jsonFilePath}`)
    }
    
    const fileContent = fs.readFileSync(jsonFilePath, 'utf8')
    let questions
    
    try {
      questions = JSON.parse(fileContent)
    } catch (parseError) {
      throw new Error(`JSON格式错误: ${parseError.message}`)
    }
    
    // 验证数据格式
    if (!Array.isArray(questions)) {
      throw new Error('JSON数据必须是数组格式')
    }
    
    if (questions.length === 0) {
      throw new Error('题目列表不能为空')
    }
    
    colorLog(`✅ 成功读取 ${questions.length} 道题目`, 'green')
    
    // 验证每个题目的必填字段
    colorLog('\n🔍 验证数据格式...', 'blue')
    for (let i = 0; i < questions.length; i++) {
      const q = questions[i]
      if (!q.content || !q.answer || !q.category) {
        throw new Error(`第${i + 1}道题目缺少必填字段 (content, answer, category)`)
      }
    }
    colorLog('✅ 数据格式验证通过', 'green')
    
    // 配置axios
    const api = axios.create({
      baseURL: serverUrl,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })
    
    // 测试服务器连接
    colorLog('\n🔗 测试服务器连接...', 'blue')
    try {
      await api.get('/question/categories')
      colorLog('✅ 服务器连接正常', 'green')
    } catch (err) {
      throw new Error(`无法连接到服务器 ${serverUrl}: ${err.message}`)
    }
    
    // 开始批量导入
    colorLog('\n🚀 开始批量导入...', 'blue')
    
    let successCount = 0
    let failCount = 0
    const errors = []
    
    for (let i = 0; i < questions.length; i++) {
      const question = questions[i]
      
      try {
        await api.post(API_ENDPOINT, {
          content: question.content.trim(),
          answer: question.answer.trim(),
          category: question.category.trim()
        })
        
        successCount++
        showProgress(i + 1, questions.length, `成功: ${successCount}, 失败: ${failCount}`)
        
        // 添加小延迟避免请求过快
        await delay(100)
        
      } catch (err) {
        failCount++
        const errorMessage = err.response?.data?.error || err.message
        errors.push(`第${i + 1}道题目导入失败: ${errorMessage}`)
        showProgress(i + 1, questions.length, `成功: ${successCount}, 失败: ${failCount}`)
      }
    }
    
    // 显示导入结果
    console.log('\n')
    colorLog('\n📊 导入结果统计:', 'bright')
    colorLog(`✅ 成功导入: ${successCount}道题目`, 'green')
    
    if (failCount > 0) {
      colorLog(`❌ 导入失败: ${failCount}道题目`, 'red')
      
      if (errors.length > 0) {
        colorLog('\n❌ 错误详情:', 'red')
        errors.slice(0, 5).forEach(error => {
          colorLog(`  ${error}`, 'red')
        })
        
        if (errors.length > 5) {
          colorLog(`  ... 还有 ${errors.length - 5} 个错误`, 'dim')
        }
      }
    }
    
    const successRate = Math.round((successCount / questions.length) * 100)
    
    if (successRate === 100) {
      colorLog(`\n🎉 批量导入完成！成功率: ${successRate}%`, 'green')
    } else if (successRate >= 80) {
      colorLog(`\n⚠️  批量导入基本完成，成功率: ${successRate}%`, 'yellow')
    } else {
      colorLog(`\n❌ 批量导入遇到较多问题，成功率: ${successRate}%`, 'red')
    }
    
  } catch (error) {
    colorLog(`\n❌ 导入失败: ${error.message}`, 'red')
    process.exit(1)
  }
}

// 显示使用帮助
function showHelp() {
  colorLog('\n📋 八股题批量导入工具', 'bright')
  colorLog('\n🔧 使用方法:', 'blue')
  colorLog('  node batch_import.js <json文件路径> [服务器地址]', 'cyan')
  
  colorLog('\n📝 示例:', 'blue')
  colorLog('  node batch_import.js example_questions.json', 'cyan')
  colorLog('  node batch_import.js ./data/questions.json http://localhost:8080', 'cyan')
  
  colorLog('\n📋 JSON格式要求:', 'blue')
  colorLog('  - 必须是题目对象的数组', 'dim')
  colorLog('  - 每个题目对象必须包含: content, answer, category', 'dim')
  colorLog('  - answer字段支持Markdown格式', 'dim')
  
  colorLog('\n📄 参考示例文件: example_questions.json', 'yellow')
}

// 主函数
async function main() {
  const args = process.argv.slice(2)
  
  if (args.length === 0 || args.includes('--help') || args.includes('-h')) {
    showHelp()
    return
  }
  
  const jsonFilePath = path.resolve(args[0])
  const serverUrl = args[1] || DEFAULT_SERVER
  
  colorLog('🎯 八股题批量导入工具', 'bright')
  colorLog(`📁 文件路径: ${jsonFilePath}`, 'dim')
  colorLog(`🌐 服务器地址: ${serverUrl}`, 'dim')
  
  await batchImport(jsonFilePath, serverUrl)
}

// 错误处理
process.on('unhandledRejection', (reason, promise) => {
  colorLog(`\n❌ 未处理的错误: ${reason}`, 'red')
  process.exit(1)
})

process.on('SIGINT', () => {
  colorLog('\n\n⏹️  用户中断操作', 'yellow')
  process.exit(0)
})

// 运行主函数
if (require.main === module) {
  main()
} 