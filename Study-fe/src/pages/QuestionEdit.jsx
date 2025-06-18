import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { questionAPI } from '../api'

function QuestionEdit() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    content: '',
    answer: '',
    category: ''
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [showPreview, setShowPreview] = useState(true)

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

  useEffect(() => {
    fetchQuestion()
  }, [id])

  const fetchQuestion = async () => {
    try {
      setLoading(true)
      setError('')
      
      // 由于后端没有单独获取单个题目的接口，我们通过获取所有题目来找到目标题目
      const response = await questionAPI.getAll()
      const questions = response.data || []
      const question = questions.find(q => q.id === parseInt(id))
      
      if (question) {
        setFormData({
          content: question.content || '',
          answer: question.answer || '',
          category: question.category || ''
        })
      } else {
        setError('题目不存在')
      }
    } catch (err) {
      setError(err.response?.data?.error || '获取题目失败')
      console.error('Fetch question error:', err)
    } finally {
      setLoading(false)
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
      setSaving(true)
      setError('')
      
      await questionAPI.update(id, {
        content: formData.content.trim(),
        answer: formData.answer.trim(),
        category: formData.category.trim()
      })
      
      setSuccess('题目更新成功！')
      
      // 2秒后跳转到列表页
      setTimeout(() => {
        navigate('/')
      }, 2000)
      
    } catch (err) {
      setError(err.response?.data?.error || '更新题目失败')
      console.error('Update question error:', err)
    } finally {
      setSaving(false)
    }
  }

  const handleReset = () => {
    fetchQuestion() // 重新获取原始数据
    setError('')
    setSuccess('')
  }

  if (loading) {
    return (
      <div className="edit-loading">
        <div className="loading-spinner"></div>
        <p>加载中...</p>
      </div>
    )
  }

  return (
    <div className="question-edit-container">
      {/* 顶部导航栏 */}
      <div className="edit-header">
        <div className="edit-header-content">
          <button 
            onClick={() => navigate('/')}
            className="back-btn"
          >
            ← 返回列表
          </button>
          <h1>✏️ 编辑题目</h1>
          <div className="edit-actions">
            <button 
              onClick={() => setShowPreview(!showPreview)}
              className="preview-toggle-btn"
            >
              {showPreview ? '👁️ 隐藏预览' : '👁️ 显示预览'}
            </button>
          </div>
        </div>
      </div>

      <div className="edit-content">
        {/* 左侧编辑区 */}
        <div className="edit-form-section">
          <div className="edit-card">
            {error && <div className="alert alert-error">❌ {error}</div>}
            {success && <div className="alert alert-success">✅ {success}</div>}
            
            <form onSubmit={handleSubmit} className="modern-form">
              <div className="form-row">
                <div className="form-field">
                  <label htmlFor="category">
                    <span className="field-icon">🏷️</span>
                    分类
                  </label>
                  <input
                    type="text"
                    id="category"
                    name="category"
                    value={formData.category}
                    onChange={handleChange}
                    className="modern-input"
                    placeholder="例如：JavaScript、React、算法等"
                    disabled={saving}
                  />
                </div>
              </div>

              <div className="form-field">
                <label htmlFor="content">
                  <span className="field-icon">❓</span>
                  题目内容
                </label>
                <textarea
                  id="content"
                  name="content"
                  value={formData.content}
                  onChange={handleChange}
                  className="modern-textarea"
                  rows="5"
                  placeholder="请输入题目内容... 支持 Markdown 语法"
                  disabled={saving}
                />
                <div className="field-hint">支持 Markdown 语法，可以使用代码块、列表等格式</div>
              </div>

              <div className="form-field">
                <label htmlFor="answer">
                  <span className="field-icon">💡</span>
                  答案内容
                </label>
                <textarea
                  id="answer"
                  name="answer"
                  value={formData.answer}
                  onChange={handleChange}
                  className="modern-textarea"
                  rows="8"
                  placeholder="请输入答案... 支持 Markdown 语法"
                  disabled={saving}
                />
                <div className="field-hint">详细回答题目，支持代码块、图片、链接等</div>
              </div>

              <div className="form-actions">
                <button 
                  type="submit" 
                  className="btn-primary modern-btn"
                  disabled={saving}
                >
                  {saving ? (
                    <>
                      <span className="loading-spinner small"></span>
                      保存中...
                    </>
                  ) : (
                    <>
                      💾 保存更改
                    </>
                  )}
                </button>
                
                <button 
                  type="button" 
                  onClick={handleReset}
                  className="btn-secondary modern-btn"
                  disabled={saving}
                >
                  🔄 重置
                </button>
              </div>
            </form>
          </div>
        </div>

        {/* 右侧预览区 */}
        {showPreview && (
          <div className="preview-section">
            <div className="preview-card">
              <div className="preview-header">
                <h3>📋 实时预览</h3>
              </div>
              
              <div className="preview-content">
                {formData.category && (
                  <div className="preview-category">
                    <span className="category-label">分类</span>
                    <span className="category-tag">
                      {formData.category}
                    </span>
                  </div>
                )}
                
                {formData.content && (
                  <div className="preview-question">
                    <h4>📝 题目</h4>
                    <div className="preview-markdown">
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
                  <div className="preview-answer">
                    <h4>💡 答案</h4>
                    <div className="preview-markdown">
                      <ReactMarkdown 
                        remarkPlugins={[remarkGfm]}
                        components={markdownComponents}
                      >
                        {formData.answer}
                      </ReactMarkdown>
                    </div>
                  </div>
                )}

                {!formData.content && !formData.answer && !formData.category && (
                  <div className="preview-empty">
                    <div className="empty-icon">📝</div>
                    <p>开始编辑以查看预览</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default QuestionEdit
