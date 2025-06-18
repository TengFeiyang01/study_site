import React from 'react'
import { Routes, Route, Link, useLocation } from 'react-router-dom'
import { ThemeProvider, useTheme } from './contexts/ThemeContext'
import QuestionList from './pages/QuestionList'
import QuestionForm from './pages/QuestionForm'
import QuestionEdit from './pages/QuestionEdit'
import CodingDashboard from './pages/CodingDashboard'

// 内部组件，可以使用主题
function AppContent() {
  const location = useLocation()
  const { theme, toggleTheme, isDark } = useTheme()
  
  // 判断当前在哪个模块
  const isCodingModule = location.pathname.startsWith('/coding')

  return (
    <div className="App">
      <header className="header">
        <div className="container">
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h1>{isCodingModule ? '刷题' : '八股复习'}</h1>
            <div style={{ display: 'flex', alignItems: 'center', gap: '20px' }}>
              <button 
                className="theme-toggle-btn"
                onClick={toggleTheme}
                title={`切换到${isDark ? '白天' : '黑夜'}模式`}
              >
                {isDark ? '☀️' : '🌙'}
              </button>
              <nav className="nav">
                <Link 
                  to="/" 
                  className={!isCodingModule ? 'active' : ''}
                >
                  八股复习
                </Link>
                <Link 
                  to="/coding" 
                  className={isCodingModule ? 'active' : ''}
                >
                  刷题
                </Link>
              </nav>
            </div>
          </div>
          
          {/* 八股复习模块的子导航 */}
          {!isCodingModule && (
            <nav className="sub-nav">
              <Link 
                to="/" 
                className={location.pathname === '/' ? 'active' : ''}
              >
                题目列表
              </Link>
              <Link 
                to="/create" 
                className={location.pathname === '/create' ? 'active' : ''}
              >
                添加题目
              </Link>
            </nav>
          )}
        </div>
      </header>

      <div className="container">
        <Routes>
          <Route path="/" element={<QuestionList />} />
          <Route path="/create" element={<QuestionForm />} />
          <Route path="/edit/:id" element={<QuestionEdit />} />
          <Route path="/category/:category" element={<QuestionList />} />
          <Route path="/coding" element={<CodingDashboard />} />
        </Routes>
      </div>
    </div>
  )
}

// 包装组件，提供主题上下文
function App() {
  return (
    <ThemeProvider>
      <AppContent />
    </ThemeProvider>
  )
}

export default App
