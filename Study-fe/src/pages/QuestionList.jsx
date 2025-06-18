import React, { useState, useEffect } from 'react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { oneDark, oneLight } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { questionAPI } from '../api'

function QuestionList() {
  const [questions, setQuestions] = useState([])
  const [categories, setCategories] = useState([])
  const [masteryStats, setMasteryStats] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('')
  const [expandedQuestion, setExpandedQuestion] = useState(null)
  const [showAnswer, setShowAnswer] = useState({})
  const [activeModule, setActiveModule] = useState('开始复习')
  const [editMode, setEditMode] = useState({})
  const { category } = useParams()
  const navigate = useNavigate()

  // 专项复习模块配置
  const specialModules = [
    { 
      name: '开始复习', 
      icon: '📚', 
      color: '#4285f4',
      description: '开始今天的八股复习'
    },
    { 
      name: '掌握度追踪', 
      icon: '📊', 
      color: '#34a853',
      description: '查看学习进度和掌握情况'
    }
  ]

  // 分类颜色映射
  const categoryColors = [
    '#ea4335', '#fbbc05', '#4285f4', '#ff6d01', '#9aa0a6', 
    '#34a853', '#ab47bc', '#ff5722', '#795548', '#607d8b'
  ]

  // 掌握度级别映射
  const masteryLevels = {
    0: { name: '未学习', color: '#9e9e9e', icon: '📝' },
    1: { name: '学习中', color: '#ff9800', icon: '📖' },
    2: { name: '已掌握', color: '#4caf50', icon: '✅' }
  }

  useEffect(() => {
    fetchData()
  }, [category])

  const fetchData = async () => {
    try {
      setLoading(true)
      setError('')
      
      // 并行获取题目、分类和掌握度统计
      const [questionsResponse, categoriesResponse, masteryResponse] = await Promise.all([
        category ? questionAPI.getByCategory(category) : questionAPI.getAll(),
        questionAPI.getCategories(),
        questionAPI.getMasteryStats()
      ])
      
      setQuestions(questionsResponse.data || [])
      setCategories(categoriesResponse.data || [])
      setMasteryStats(masteryResponse.data || {})
      
      if (category) {
        setSelectedCategory(category)
      }
      
    } catch (err) {
      setError(err.response?.data?.error || '获取数据失败')
      console.error('Fetch data error:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleMasteryUpdate = async (questionId, newLevel) => {
    try {
      await questionAPI.updateMasteryLevel(questionId, newLevel)
      
      // 更新本地状态
      setQuestions(prev => prev.map(q => 
        q.id === questionId ? { ...q, mastery_level: newLevel } : q
      ))
      
      // 重新获取统计数据
      const masteryResponse = await questionAPI.getMasteryStats()
      setMasteryStats(masteryResponse.data || {})
      
    } catch (err) {
      setError(err.response?.data?.error || '更新掌握度失败')
      console.error('Update mastery error:', err)
    }
  }

  const handleDelete = async (id) => {
    if (!window.confirm('确定要删除这道题目吗？')) {
      return
    }

    try {
      await questionAPI.delete(id)
      setQuestions(questions.filter(q => q.id !== id))
      // 重新获取统计数据
      const masteryResponse = await questionAPI.getMasteryStats()
      setMasteryStats(masteryResponse.data || {})
    } catch (err) {
      setError(err.response?.data?.error || '删除失败')
      console.error('Delete question error:', err)
    }
  }

  const handleCategoryChange = (newCategory) => {
    setSelectedCategory(newCategory)
    setActiveModule(newCategory || '开始复习')
    if (newCategory) {
      navigate(`/category/${newCategory}`)
    } else {
      navigate('/')
    }
  }

  const toggleQuestion = (questionId) => {
    if (expandedQuestion === questionId) {
      setExpandedQuestion(null)
      setShowAnswer(prev => ({ ...prev, [questionId]: false }))
    } else {
      setExpandedQuestion(questionId)
      setShowAnswer(prev => ({ ...prev, [questionId]: false }))
    }
  }

  const toggleAnswer = (questionId) => {
    setShowAnswer(prev => ({
      ...prev,
      [questionId]: !prev[questionId]
    }))
  }

  const toggleEditMode = (questionId) => {
    setEditMode(prev => ({
      ...prev,
      [questionId]: !prev[questionId]
    }))
  }

  // 过滤题目
  const filteredQuestions = questions.filter(question => {
    const matchesSearch = question.content?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         question.answer?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         question.category?.toLowerCase().includes(searchTerm.toLowerCase())
    
    const matchesCategory = !selectedCategory || question.category === selectedCategory
    
    return matchesSearch && matchesCategory
  })

  // 获取分类统计
  const getCategoryStats = (categoryName) => {
    return questions.filter(q => q.category === categoryName).length
  }

  // 复制状态管理
  const [copiedCode, setCopiedCode] = useState({})
  // 代码块展开状态管理
  const [expandedCode, setExpandedCode] = useState({})
  // 代码主题管理
  const [codeTheme, setCodeTheme] = useState('dark')

  // 复制代码函数
  const copyToClipboard = async (text, codeId) => {
    try {
      await navigator.clipboard.writeText(text)
      // 设置复制成功状态
      setCopiedCode(prev => ({ ...prev, [codeId]: true }));
      // 2秒后重置状态
      setTimeout(() => {
        setCopiedCode(prev => {
          const newState = { ...prev };
          delete newState[codeId]; // 删除状态而不是设置为false
          return newState;
        });
      }, 2000);
    } catch (err) {
      console.error('复制失败:', err);
      // 设置失败状态
      setCopiedCode(prev => ({ ...prev, [codeId]: 'error' }));
      // 2秒后重置状态
      setTimeout(() => {
        setCopiedCode(prev => {
          const newState = { ...prev };
          delete newState[codeId];
          return newState;
        });
      }, 2000);
    }
  }

  // 切换代码块展开状态
  const toggleCodeExpand = (codeId) => {
    setExpandedCode(prev => ({ ...prev, [codeId]: !prev[codeId] }))
  }

  // 移除复杂的语法高亮，保持代码原始显示
  const highlightCode = (code, language) => {
    return code // 直接返回原始代码，不进行HTML标签处理
  }

  // 自定义markdown组件
  const markdownComponents = {
    code({ node, inline, className, children, ...props }) {
      const match = /language-(\w+)/.exec(className || '')
      const language = match ? match[1] : 'code'
      const codeContent = String(children).replace(/\n$/, '')
      // 使用代码内容的hash作为稳定的ID，添加错误处理
      let codeId;
      try {
        codeId = `code-${btoa(encodeURIComponent(codeContent.substring(0, 30))).replace(/[^a-zA-Z0-9]/g, '').substring(0, 10)}`;
      } catch (e) {
        // 如果btoa失败，使用简单的hash
        codeId = `code-${codeContent.length}-${codeContent.substring(0, 10).replace(/[^a-zA-Z0-9]/g, '')}`;
      }
      
      if (!inline && match) {
        const isExpanded = expandedCode[codeId]
        const shouldCollapse = codeContent.split('\n').length > 8
        
        return (
          <div className="modern-code-block">
            <div className="code-header">
              <div className="code-info">
                {shouldCollapse && (
                  <button 
                    className="expand-btn"
                    onClick={() => toggleCodeExpand(codeId)}
                    title={isExpanded ? "收起代码" : "展开代码"}
                  >
                    {isExpanded ? '−' : '+'}
                  </button>
                )}
                <span className="language-tag">{language}</span>
              </div>
              <div className="code-actions">
                <button 
                  className="theme-btn"
                  onClick={() => setCodeTheme(codeTheme === 'dark' ? 'light' : 'dark')}
                  title="切换主题"
                >
                  {codeTheme === 'dark' ? '☀️' : '🌙'}
                </button>
                <button 
                  className="copy-button"
                  onClick={() => copyToClipboard(codeContent, codeId)}
                  title="复制代码"
                >
                  {copiedCode[codeId] === true ? '✓ 已复制' : 
                   copiedCode[codeId] === 'error' ? '✗ 失败' : '📄 复制'}
                </button>
              </div>
            </div>
            <div className={`code-content ${shouldCollapse && !isExpanded ? 'collapsed' : ''}`}>
              <SyntaxHighlighter
                language={language.toLowerCase()}
                style={codeTheme === 'dark' ? oneDark : oneLight}
                customStyle={{
                  margin: 0,
                  padding: '20px',
                  background: 'transparent',
                  fontSize: '14px',
                  lineHeight: '1.6'
                }}
                wrapLines={true}
                wrapLongLines={true}
              >
                {codeContent}
              </SyntaxHighlighter>
              {shouldCollapse && !isExpanded && (
                <div className="code-fade">
                  <button 
                    className="expand-overlay-btn"
                    onClick={() => toggleCodeExpand(codeId)}
                  >
                    点击展开完整代码 ({codeContent.split('\n').length} 行)
                  </button>
                </div>
              )}
            </div>
          </div>
        )
      } else {
        return (
          <code className="inline-code" {...props}>
            {children}
          </code>
        )
      }
    }
  }

  // 渲染主要内容
  const renderMainContent = () => {
    if (activeModule === '掌握度追踪') {
      return (
        <div className="mastery-tracking">
          <h2>📊 掌握度追踪</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <h3>总题目数</h3>
              <div className="stat-number">{masteryStats.total || 0}</div>
            </div>
            <div className="stat-card">
              <h3>已掌握</h3>
              <div className="stat-number" style={{ color: '#4caf50' }}>
                {masteryStats.mastered || 0}
              </div>
            </div>
            <div className="stat-card">
              <h3>学习中</h3>
              <div className="stat-number" style={{ color: '#ff9800' }}>
                {masteryStats.learning || 0}
              </div>
            </div>
            <div className="stat-card">
              <h3>未学习</h3>
              <div className="stat-number" style={{ color: '#9e9e9e' }}>
                {masteryStats.unlearned || 0}
              </div>
            </div>
          </div>
          
          {/* 进度条 */}
          <div className="mastery-progress" style={{ marginTop: '2rem' }}>
            <h3>学习进度</h3>
            <div className="progress-bar">
              <div 
                className="progress-fill" 
                style={{ 
                  width: `${masteryStats.total ? (masteryStats.mastered / masteryStats.total * 100) : 0}%`,
                  background: 'linear-gradient(90deg, #4caf50, #8bc34a)'
                }}
              ></div>
            </div>
            <p>{masteryStats.total ? Math.round(masteryStats.mastered / masteryStats.total * 100) : 0}% 已掌握</p>
          </div>
        </div>
      )
    }

    return (
      <div>
        {/* 现代化的顶部横幅 */}
        <div className="modern-header">
          <div className="header-content">
            <div className="header-left">
              <div className="header-icon">🎯</div>
              <div className="header-text">
                <h1>八股文题库</h1>
                <p>掌握核心知识点，提升面试成功率</p>
              </div>
            </div>
            <div className="header-actions">
              <Link to="/create" className="btn btn-modern">
                <span className="btn-icon">➕</span>
                添加题目
              </Link>
            </div>
          </div>
          
          {/* 统计卡片 */}
          <div className="stats-cards">
            <div className="stat-card-mini">
              <div className="stat-icon">📚</div>
              <div className="stat-info">
                <div className="stat-number">{questions.length}</div>
                <div className="stat-label">总题目</div>
              </div>
            </div>
            <div className="stat-card-mini">
              <div className="stat-icon">📝</div>
              <div className="stat-info">
                <div className="stat-number">{categories.length}</div>
                <div className="stat-label">分类数</div>
              </div>
            </div>
            <div className="stat-card-mini">
              <div className="stat-icon">🎉</div>
              <div className="stat-info">
                <div className="stat-number">{filteredQuestions.length}</div>
                <div className="stat-label">当前显示</div>
              </div>
            </div>
          </div>
        </div>

        <div className="search-bar">
          <div className="search-wrapper">
            <span className="search-icon">🔍</span>
            <input
              type="text"
              placeholder="搜索题目内容、答案或分类..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
            {searchTerm && (
              <button 
                className="search-clear"
                onClick={() => setSearchTerm('')}
              >
                ✕
              </button>
            )}
          </div>
        </div>

        {error && <div className="error">{error}</div>}

        <div className="questions-list">
          {filteredQuestions.length === 0 ? (
            <div className="no-questions">
              {loading ? '加载中...' : '暂无题目'}
            </div>
          ) : (
            filteredQuestions.map((question, index) => (
              <div key={question.id} className="question-item">
                <div className="question-main" onClick={() => toggleQuestion(question.id)}>
                  <div className="question-content">
                    <ReactMarkdown 
                      remarkPlugins={[remarkGfm]}
                      components={markdownComponents}
                    >
                      {`${index + 1}. ${question.content}`}
                    </ReactMarkdown>
                  </div>
                  
                  <div className="question-actions">
                    <button 
                      onClick={(e) => {
                        e.stopPropagation()
                        toggleAnswer(question.id)
                      }}
                      className="btn btn-sm btn-outline-primary action-btn"
                    >
                      {showAnswer[question.id] ? '隐藏答案' : '显示答案'}
                    </button>
                    
                    <button 
                      onClick={(e) => {
                        e.stopPropagation()
                        toggleEditMode(question.id)
                      }} 
                      className="btn btn-sm btn-secondary action-btn"
                    >
                      {editMode[question.id] ? '取消' : '编辑'}
                    </button>
                    
                    {editMode[question.id] && (
                      <>
                        <Link 
                          to={`/edit/${question.id}`} 
                          className="btn btn-sm btn-outline-primary action-btn"
                          onClick={(e) => e.stopPropagation()}
                        >
                          修改
                        </Link>
                        <button 
                          onClick={(e) => {
                            e.stopPropagation()
                            handleDelete(question.id)
                          }} 
                          className="btn btn-sm btn-danger action-btn"
                        >
                          删除
                        </button>
                      </>
                    )}
                  </div>
                </div>
                
                {showAnswer[question.id] && (
                  <div className="question-answer">
                    <ReactMarkdown 
                      remarkPlugins={[remarkGfm]}
                      components={markdownComponents}
                    >
                      {question.answer}
                    </ReactMarkdown>
                    
                    {/* 掌握度按钮移到答案末尾 */}
                    <div className="answer-actions">
                      <div className="mastery-level-inline">
                        <span className="mastery-label">掌握度：</span>
                        <button
                          onClick={(e) => {
                            e.stopPropagation()
                            const currentLevel = question.mastery_level || 0
                            const nextLevel = (currentLevel + 1) % 3
                            handleMasteryUpdate(question.id, nextLevel)
                          }}
                          className="mastery-current-btn"
                          style={{ 
                            backgroundColor: masteryLevels[question.mastery_level || 0].color,
                            color: 'white',
                            padding: '4px 10px',
                            borderRadius: '8px',
                            fontSize: '0.8rem',
                            border: 'none',
                            cursor: 'pointer',
                            transition: 'all 0.2s ease'
                          }}
                          title={`当前: ${masteryLevels[question.mastery_level || 0].name}，点击切换到下一个状态`}
                        >
                          {masteryLevels[question.mastery_level || 0].icon} {masteryLevels[question.mastery_level || 0].name}
                        </button>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="question-list-container">
      {/* 左侧边栏 */}
      <div className="sidebar">
        {/* 专项复习部分 */}
        <div className="sidebar-section">
          <h3 className="sidebar-title">专项复习</h3>
          <div className="sidebar-items">
            {specialModules.map((module) => (
              <div
                key={module.name}
                className={`sidebar-item ${activeModule === module.name ? 'active' : ''}`}
                onClick={() => {
                  setActiveModule(module.name)
                  setSelectedCategory('')
                  navigate('/')
                }}
              >
                <span className="sidebar-icon">{module.icon}</span>
                <span className="sidebar-text">{module.name}</span>
              </div>
            ))}
          </div>
        </div>

        {/* 按分类复习部分 */}
        <div className="sidebar-section">
          <h3 className="sidebar-title">按分类复习</h3>
          <div className="sidebar-items">
            <div
              className={`sidebar-item ${selectedCategory === '' ? 'active' : ''}`}
              onClick={() => handleCategoryChange('')}
            >
              <span className="sidebar-icon">📋</span>
              <span className="sidebar-text">全部题目</span>
              <span className="category-count">{questions.length}</span>
            </div>
            {categories.map((cat, index) => (
              <div
                key={cat}
                className={`sidebar-item ${selectedCategory === cat ? 'active' : ''}`}
                onClick={() => handleCategoryChange(cat)}
              >
                <span 
                  className="category-indicator"
                  style={{ backgroundColor: categoryColors[index % categoryColors.length] }}
                ></span>
                <span className="sidebar-text">{cat}</span>
                <span className="category-count">{getCategoryStats(cat)}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* 主要内容区域 */}
      <div className="main-content">
        {renderMainContent()}
      </div>
    </div>
  )
}

export default QuestionList
