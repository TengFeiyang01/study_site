import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { questionAPI } from '../api'

function QuestionForm() {
  const [formData, setFormData] = useState({
    content: '',
    answer: '',
    category: ''
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [activeTab, setActiveTab] = useState('single') // 'single' 或 'batch'
  
  // 批量导入相关状态
  const [jsonData, setJsonData] = useState('')
  const [batchLoading, setBatchLoading] = useState(false)
  const [batchError, setBatchError] = useState('')
  const [batchSuccess, setBatchSuccess] = useState('')
  const [importProgress, setImportProgress] = useState({ current: 0, total: 0 })
  
  const navigate = useNavigate()

  // 自定义markdown组件
  const markdownComponents = {
    code({ node, inline, className, children, ...props }) {
      const match = /language-(\w+)/.exec(className || '')
      return !inline && match ? (
        <SyntaxHighlighter
          style={tomorrow}
          language={match[1]}
          PreTag="div"
          customStyle={{
            background: '#f8f9fa',
            border: '1px solid #e9ecef',
            borderRadius: '8px',
            fontSize: '0.9rem',
            lineHeight: '1.4'
          }}
          {...props}
        >
          {String(children).replace(/\n$/, '')}
        </SyntaxHighlighter>
      ) : (
        <code className={className} {...props}>
          {children}
        </code>
      )
    }
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
    // 清除错误信息
    if (error) setError('')
    if (success) setSuccess('')
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    
    // 表单验证
    if (!formData.content.trim()) {
      setError('请输入题目内容')
      return
    }
    if (!formData.answer.trim()) {
      setError('请输入答案')
      return
    }
    if (!formData.category.trim()) {
      setError('请输入分类')
      return
    }

    try {
      setLoading(true)
      setError('')
      
      await questionAPI.create({
        content: formData.content.trim(),
        answer: formData.answer.trim(),
        category: formData.category.trim()
      })
      
      setSuccess('题目添加成功！')
      
      // 2秒后跳转到列表页
      setTimeout(() => {
        navigate('/')
      }, 2000)
      
    } catch (err) {
      setError(err.response?.data?.error || '添加题目失败')
      console.error('Create question error:', err)
    } finally {
      setLoading(false)
    }
  }

  // JSON批量导入功能
  const handleBatchImport = async () => {
    if (!jsonData.trim()) {
      setBatchError('请输入JSON数据')
      return
    }

    try {
      setBatchLoading(true)
      setBatchError('')
      setBatchSuccess('')
      setImportProgress({ current: 0, total: 0 })
      
      // 解析JSON数据
      let questions
      try {
        questions = JSON.parse(jsonData.trim())
      } catch (parseError) {
        setBatchError('JSON格式错误，请检查格式是否正确')
        return
      }

      // 验证数据格式
      if (!Array.isArray(questions)) {
        setBatchError('JSON数据必须是数组格式')
        return
      }

      if (questions.length === 0) {
        setBatchError('题目列表不能为空')
        return
      }

      // 验证每个题目的必填字段
      for (let i = 0; i < questions.length; i++) {
        const q = questions[i]
        if (!q.content || !q.answer || !q.category) {
          setBatchError(`第${i + 1}道题目缺少必填字段 (content, answer, category)`)
          return
        }
      }

      setImportProgress({ current: 0, total: questions.length })

      // 批量导入题目
      let successCount = 0
      let failCount = 0
      const errors = []

      for (let i = 0; i < questions.length; i++) {
        const question = questions[i]
        try {
          await questionAPI.create({
            content: question.content.trim(),
            answer: question.answer.trim(),
            category: question.category.trim()
          })
          successCount++
          setImportProgress({ current: i + 1, total: questions.length })
          
          // 添加小延迟避免请求过快
          await new Promise(resolve => setTimeout(resolve, 100))
        } catch (err) {
          failCount++
          errors.push(`第${i + 1}道题目导入失败: ${err.response?.data?.error || err.message}`)
          console.error(`Import question ${i + 1} failed:`, err)
        }
      }

      // 显示导入结果
      if (failCount === 0) {
        setBatchSuccess(`🎉 批量导入成功！共导入${successCount}道题目`)
        setJsonData('') // 清空输入
        
        // 3秒后跳转到列表页
        setTimeout(() => {
          navigate('/')
        }, 3000)
      } else {
        setBatchSuccess(`部分导入成功：成功${successCount}道，失败${failCount}道`)
        if (errors.length > 0) {
          setBatchError(errors.slice(0, 3).join('\n') + (errors.length > 3 ? '\n...' : ''))
        }
      }

    } catch (err) {
      setBatchError('批量导入失败: ' + (err.message || '未知错误'))
      console.error('Batch import error:', err)
    } finally {
      setBatchLoading(false)
    }
  }

  const handleReset = () => {
    setFormData({
      content: '',
      answer: '',
      category: ''
    })
    setError('')
    setSuccess('')
  }

  // JSON示例数据
  const jsonExample = `[
  {
    "content": "什么是JavaScript闭包？",
    "answer": "## 闭包定义\\n\\n闭包是指一个函数可以访问其**外部作用域**中的变量，即使外部函数已经执行完毕。\\n\\n### 代码示例\\n\\n\`\`\`javascript\\nfunction outer() {\\n  let count = 0;\\n  return function inner() {\\n    count++;\\n    console.log(count);\\n  };\\n}\\n\\nconst counter = outer();\\ncounter(); // 1\\ncounter(); // 2\\n\`\`\`",
    "category": "JavaScript"
  },
  {
    "content": "解释React中的Hook是什么？",
    "answer": "Hook是React 16.8新增的特性，让你在不编写class的情况下使用state以及其他的React特性。\\n\\n**常用Hook:**\\n- useState\\n- useEffect\\n- useContext\\n- useReducer",
    "category": "React"
  }
]`

  return (
    <div>
      {/* 顶部标签切换 */}
      <div className="card">
        <div className="tab-buttons" style={{ marginBottom: '2rem' }}>
          <button 
            className={`tab-btn ${activeTab === 'single' ? 'active' : ''}`}
            onClick={() => setActiveTab('single')}
          >
            📝 单个添加
          </button>
          <button 
            className={`tab-btn ${activeTab === 'batch' ? 'active' : ''}`}
            onClick={() => setActiveTab('batch')}
          >
            📋 批量导入
          </button>
        </div>

        {/* 单个添加表单 */}
        {activeTab === 'single' && (
          <>
            <h2>添加新题目</h2>
            
            {error && <div className="error">{error}</div>}
            {success && <div className="success">{success}</div>}
            
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label htmlFor="category">分类 *</label>
                <input
                  type="text"
                  id="category"
                  name="category"
                  value={formData.category}
                  onChange={handleChange}
                  className="form-control"
                  placeholder="例如：JavaScript、React、算法等"
                  disabled={loading}
                />
              </div>

              <div className="form-group">
                <label htmlFor="content">题目内容 *</label>
                <textarea
                  id="content"
                  name="content"
                  value={formData.content}
                  onChange={handleChange}
                  className="form-control"
                  rows="4"
                  placeholder="请输入题目内容..."
                  disabled={loading}
                />
              </div>

              <div className="form-group">
                <label htmlFor="answer">答案 *</label>
                <textarea
                  id="answer"
                  name="answer"
                  value={formData.answer}
                  onChange={handleChange}
                  className="form-control"
                  rows="6"
                  placeholder="请输入答案（支持Markdown格式）..."
                  disabled={loading}
                />
              </div>

              <div style={{ display: 'flex', gap: '1rem', marginTop: '2rem' }}>
                <button 
                  type="submit" 
                  className="btn btn-primary"
                  disabled={loading}
                >
                  {loading ? '添加中...' : '添加题目'}
                </button>
                
                <button 
                  type="button" 
                  onClick={handleReset}
                  className="btn btn-secondary"
                  disabled={loading}
                >
                  重置
                </button>
                
                <button 
                  type="button" 
                  onClick={() => navigate('/')}
                  className="btn btn-secondary"
                  disabled={loading}
                >
                  返回列表
                </button>
              </div>
            </form>
          </>
        )}

        {/* 批量导入表单 */}
        {activeTab === 'batch' && (
          <>
            <h2>📋 批量导入题目</h2>
            
            {batchError && (
              <div className="error" style={{ whiteSpace: 'pre-line' }}>
                {batchError}
              </div>
            )}
            {batchSuccess && (
              <div className="success">{batchSuccess}</div>
            )}

            {/* 导入进度 */}
            {batchLoading && importProgress.total > 0 && (
              <div className="import-progress" style={{ 
                marginBottom: '1rem', 
                padding: '1rem', 
                background: '#f8f9fa', 
                borderRadius: '8px',
                border: '1px solid #e9ecef'
              }}>
                <div style={{ marginBottom: '0.5rem' }}>
                  导入进度: {importProgress.current} / {importProgress.total}
                </div>
                <div style={{ 
                  width: '100%', 
                  height: '8px', 
                  background: '#e9ecef', 
                  borderRadius: '4px',
                  overflow: 'hidden'
                }}>
                  <div style={{ 
                    width: `${(importProgress.current / importProgress.total) * 100}%`, 
                    height: '100%', 
                    background: '#28a745',
                    transition: 'width 0.3s ease'
                  }}></div>
                </div>
              </div>
            )}
            
            <div className="form-group">
              <label htmlFor="jsonData">
                JSON数据 * 
                <span style={{ fontSize: '0.9em', color: '#666', marginLeft: '0.5rem' }}>
                  (请输入题目数组格式的JSON数据)
                </span>
              </label>
              <textarea
                id="jsonData"
                value={jsonData}
                onChange={(e) => setJsonData(e.target.value)}
                className="form-control"
                rows="15"
                placeholder="请粘贴JSON格式的题目数据..."
                disabled={batchLoading}
                style={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
              />
            </div>

            <div style={{ display: 'flex', gap: '1rem', marginTop: '2rem' }}>
              <button 
                onClick={handleBatchImport}
                className="btn btn-primary"
                disabled={batchLoading}
              >
                {batchLoading ? '导入中...' : '🚀 开始批量导入'}
              </button>
              
              <button 
                onClick={() => setJsonData(jsonExample)}
                className="btn btn-secondary"
                disabled={batchLoading}
              >
                📄 填充示例数据
              </button>
              
              <button 
                onClick={() => {
                  setJsonData('')
                  setBatchError('')
                  setBatchSuccess('')
                  setImportProgress({ current: 0, total: 0 })
                }}
                className="btn btn-secondary"
                disabled={batchLoading}
              >
                🧹 清空数据
              </button>
              
              <button 
                type="button" 
                onClick={() => navigate('/')}
                className="btn btn-secondary"
                disabled={batchLoading}
              >
                返回列表
              </button>
            </div>

            {/* JSON格式说明 */}
            <div className="card" style={{ marginTop: '2rem', background: '#f8f9fa' }}>
              <h4>📖 JSON格式说明</h4>
              <p>请提供一个包含题目对象的数组，每个题目对象必须包含以下字段：</p>
              <ul>
                <li><code>content</code> (必填): 题目内容</li>
                <li><code>answer</code> (必填): 题目答案，支持Markdown格式</li>
                <li><code>category</code> (必填): 题目分类</li>
              </ul>
              
              <details style={{ marginTop: '1rem' }}>
                <summary style={{ cursor: 'pointer', fontWeight: 'bold' }}>
                  💡 点击查看示例格式
                </summary>
                <pre style={{ 
                  background: 'white', 
                  padding: '1rem', 
                  borderRadius: '4px',
                  border: '1px solid #ddd',
                  marginTop: '0.5rem',
                  fontSize: '0.85rem',
                  overflow: 'auto'
                }}>
                  {jsonExample}
                </pre>
              </details>
            </div>
          </>
        )}
      </div>

      {/* 预览区域 - 只在单个添加模式显示 */}
      {activeTab === 'single' && (formData.content || formData.answer || formData.category) && (
        <div className="card" style={{ marginTop: '2rem' }}>
          <h3>预览</h3>
          <div style={{ border: '1px solid #e9ecef', borderRadius: '8px', padding: '1rem' }}>
            {formData.category && (
              <div style={{ marginBottom: '1rem' }}>
                <strong>分类：</strong>
                <span style={{ 
                  background: '#e9ecef', 
                  padding: '0.2rem 0.5rem', 
                  borderRadius: '4px',
                  fontSize: '0.85rem',
                  marginLeft: '0.5rem'
                }}>
                  {formData.category}
                </span>
              </div>
            )}
            {formData.content && (
              <div style={{ marginBottom: '1rem' }}>
                <strong>题目：</strong>
                <div style={{ marginTop: '0.5rem' }}>
                  <ReactMarkdown 
                    remarkPlugins={[remarkGfm]}
                    components={markdownComponents}
                  >
                    {formData.content}
                  </ReactMarkdown>
                </div>
              </div>
            )}
            {formData.answer && (
              <div>
                <strong>答案：</strong>
                <div style={{ marginTop: '0.5rem' }}>
                  <ReactMarkdown 
                    remarkPlugins={[remarkGfm]}
                    components={markdownComponents}
                  >
                    {formData.answer}
                  </ReactMarkdown>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default QuestionForm
