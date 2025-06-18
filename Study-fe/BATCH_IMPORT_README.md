# 八股题批量导入使用指南

## 📋 概述

本项目提供了两种方式批量导入八股题：
1. **前端网页导入**：通过浏览器界面导入
2. **命令行脚本导入**：使用Node.js脚本批量导入

## 🌐 方式一：前端网页导入

### 使用步骤

1. **打开添加题目页面**
   ```
   http://localhost:3000/create
   ```

2. **切换到批量导入标签**
   - 点击页面顶部的 "📋 批量导入" 标签

3. **准备JSON数据**
   - 可以点击 "📄 填充示例数据" 查看格式
   - 或者直接粘贴你的JSON数据

4. **开始导入**
   - 点击 "🚀 开始批量导入" 按钮
   - 查看实时导入进度
   - 等待导入完成

### 特点
- ✅ 界面友好，操作简单
- ✅ 实时进度显示
- ✅ 详细错误提示
- ✅ 支持Markdown预览
- ❌ 单次导入数量建议不超过50题（避免浏览器超时）

## 💻 方式二：命令行脚本导入

### 环境要求
```bash
# 需要安装axios依赖
npm install axios
```

### 使用方法

```bash
# 基本用法
node batch_import.js <json文件路径>

# 指定服务器地址
node batch_import.js <json文件路径> <服务器地址>
```

### 使用示例

```bash
# 导入示例文件
node batch_import.js example_questions.json

# 导入自定义文件
node batch_import.js ./data/my_questions.json

# 指定服务器地址
node batch_import.js questions.json http://192.168.1.100:8080

# 查看帮助
node batch_import.js --help
```

### 特点
- ✅ 支持大批量导入（无数量限制）
- ✅ 彩色进度条显示
- ✅ 详细的统计报告
- ✅ 自动错误重试
- ✅ 可中断操作（Ctrl+C）

## 📋 JSON数据格式

### 基本格式

```json
[
  {
    "content": "题目内容",
    "answer": "题目答案（支持Markdown）",
    "category": "题目分类"
  }
]
```

### 字段说明

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `content` | string | ✅ | 题目内容，支持Markdown格式 |
| `answer` | string | ✅ | 题目答案，支持Markdown格式 |
| `category` | string | ✅ | 题目分类（如：JavaScript、React等）|

### 完整示例

```json
[
  {
    "content": "什么是JavaScript闭包？",
    "answer": "## 闭包定义\n\n闭包是指一个函数可以访问其**外部作用域**中的变量，即使外部函数已经执行完毕。\n\n### 代码示例\n\n```javascript\nfunction outer() {\n  let count = 0;\n  return function inner() {\n    count++;\n    console.log(count);\n  };\n}\n```",
    "category": "JavaScript"
  },
  {
    "content": "解释React中的Hook是什么？",
    "answer": "Hook是React 16.8新增的特性，让你在不编写class的情况下使用state以及其他的React特性。\n\n**常用Hook:**\n- useState\n- useEffect\n- useContext",
    "category": "React"
  }
]
```

## 🎨 Markdown格式支持

### 支持的语法

- **标题**：`# ## ###`
- **粗体**：`**文本**`
- **斜体**：`*文本*`
- **代码块**：````javascript`
- **行内代码**：`代码`
- **列表**：`- 项目` 或 `1. 项目`
- **链接**：`[文本](URL)`
- **表格**：支持GitHub风格表格
- **引用**：`> 引用内容`

### 代码高亮

支持多种编程语言的语法高亮：

```json
{
  "answer": "### JavaScript示例\n\n```javascript\nconst arr = [1, 2, 3];\nconsole.log(arr);\n```\n\n### Python示例\n\n```python\ndef hello():\n    print('Hello World')\n```"
}
```

## 🔧 导入最佳实践

### 数据准备

1. **检查JSON格式**
   ```bash
   # 使用在线JSON验证器检查格式
   # 或使用VS Code的JSON格式化功能
   ```

2. **控制批次大小**
   ```bash
   # 前端导入：建议每次 ≤ 50 题
   # 命令行导入：可支持大批量（500+ 题）
   ```

3. **备份数据**
   ```bash
   # 导入前建议备份原有数据
   # 可以先在测试环境验证
   ```

### 错误处理

1. **常见错误及解决方案**

   | 错误类型 | 原因 | 解决方案 |
   |----------|------|----------|
   | JSON格式错误 | 语法不正确 | 使用JSON验证器检查 |
   | 缺少必填字段 | content/answer/category为空 | 检查数据完整性 |
   | 服务器连接失败 | 后端未启动 | 确保后端服务正常运行 |
   | 超时错误 | 批量过大 | 减少单次导入数量 |

2. **重试策略**
   ```bash
   # 如果部分导入失败，可以：
   # 1. 查看错误日志
   # 2. 修复问题数据
   # 3. 重新导入失败的题目
   ```

## 📊 导入性能

### 性能对比

| 导入方式 | 适用场景 | 导入速度 | 推荐批量 |
|----------|----------|----------|----------|
| 前端网页 | 小批量、临时导入 | 中等 | ≤ 50题 |
| 命令行脚本 | 大批量、批处理 | 快 | 无限制 |

### 优化建议

1. **大批量导入**：使用命令行脚本
2. **小批量导入**：使用前端界面
3. **定期导入**：可以编写自动化脚本
4. **数据验证**：导入前先验证JSON格式

## 🔍 故障排除

### 前端导入问题

```bash
# 1. 检查控制台错误
F12 → Console → 查看错误信息

# 2. 检查网络请求
F12 → Network → 查看API请求状态

# 3. 验证JSON格式
使用在线JSON验证器
```

### 命令行脚本问题

```bash
# 1. 检查Node.js版本
node --version  # 需要 v12+

# 2. 安装依赖
npm install axios

# 3. 检查文件路径
ls -la example_questions.json

# 4. 测试服务器连接
curl http://localhost:8080/question/categories
```

## 📚 示例文件

项目中提供了 `example_questions.json` 示例文件，包含：

- ✅ JavaScript相关题目
- ✅ React相关题目  
- ✅ CSS相关题目
- ✅ 网络相关题目
- ✅ Vue相关题目

可以直接使用此文件测试导入功能：

```bash
node batch_import.js example_questions.json
```

## 🎯 总结

选择合适的导入方式：

- **< 50题** → 使用前端网页导入
- **> 50题** → 使用命令行脚本导入
- **偶尔导入** → 前端网页更方便
- **批量处理** → 命令行脚本更高效

无论使用哪种方式，都要注意：
1. 提前验证JSON格式
2. 确保后端服务运行正常
3. 合理控制导入批量
4. 及时处理错误信息 