import React, { useState, useEffect } from 'react'
import { codingProblemAPI } from '../api'

const CodingDashboard = () => {
  const [problems, setProblems] = useState([])
  const [dailyProblem, setDailyProblem] = useState(null)
  const [randomProblem, setRandomProblem] = useState(null)
  const [dailyHistory, setDailyHistory] = useState([])
  const [filteredProblems, setFilteredProblems] = useState([])
  const [selectedDifficulty, setSelectedDifficulty] = useState('all')
  const [selectedSource, setSelectedSource] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [activeTab, setActiveTab] = useState('all')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [stats, setStats] = useState({})
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage] = useState(10)

  // 获取所有题目的来源列表
  const uniqueSources = [...new Set(problems.map(p => p.source).filter(Boolean))]

  // 获取当前页的题目
  const currentProblems = activeTab === 'daily'
    ? dailyHistory.length > 0 ? dailyHistory : (dailyProblem ? [dailyProblem] : [])
    : filteredProblems.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)

  // 获取难度配置
  const getDifficultyConfig = (difficulty) => {
    const config = {
      '简单': { text: '简单', bg: 'from-green-500 to-emerald-400', textColor: 'text-green-600' },
      'Easy': { text: '简单', bg: 'from-green-500 to-emerald-400', textColor: 'text-green-600' },
      '中等': { text: '中等', bg: 'from-yellow-500 to-amber-400', textColor: 'text-yellow-600' },
      'Medium': { text: '中等', bg: 'from-yellow-500 to-amber-400', textColor: 'text-yellow-600' },
      '困难': { text: '困难', bg: 'from-red-500 to-rose-400', textColor: 'text-red-600' },
      'Hard': { text: '困难', bg: 'from-red-500 to-rose-400', textColor: 'text-red-600' }
    }
    return config[difficulty] || { text: '未知', bg: 'from-gray-500 to-gray-400', textColor: 'text-gray-600' }
  }

  // 获取每日一题历史
  const fetchDailyHistory = async () => {
    try {
      const response = await codingProblemAPI.getDailyHistory()
      if (response.data && response.data.problems) {
        setDailyHistory(response.data.problems)
      }
    } catch (err) {
      console.error('获取每日一题历史失败:', err)
    }
  }

  // 获取所有数据
  const fetchData = async () => {
    console.log('开始获取数据...')
    setLoading(true)
    setError('')
    try {
      // 获取每日一题
      console.log('正在获取每日一题...')
      let retryCount = 0;
      let dailyResponse = null;
      
      while (retryCount < 3) {
        try {
          dailyResponse = await codingProblemAPI.getDaily()
          if (dailyResponse?.data) {
            break;
          }
          retryCount++;
          await new Promise(resolve => setTimeout(resolve, 1000)); // 等待1秒后重试
        } catch (err) {
          console.error(`获取每日一题失败(第${retryCount + 1}次尝试):`, err)
          retryCount++;
          if (retryCount < 3) {
            await new Promise(resolve => setTimeout(resolve, 1000));
          }
        }
      }

      console.log('每日一题数据:', dailyResponse)
      if (dailyResponse?.data) {
        setDailyProblem(dailyResponse.data)
      } else {
        console.error('无法获取每日一题')
        setError('无法获取每日一题，请稍后重试')
      }

      // 获取随机一题
      console.log('正在获取随机一题...')
      const randomResponse = await codingProblemAPI.getRandom()
      console.log('随机一题数据:', randomResponse)
      if (randomResponse?.data) {
        setRandomProblem(randomResponse.data)
      }

      // 获取所有题目
      console.log('正在获取所有题目...')
      const problemsResponse = await codingProblemAPI.getAll({})
      console.log('题目数据:', problemsResponse)
      if (problemsResponse?.data) {
        setProblems(problemsResponse.data)
      }

      // 获取统计信息
      console.log('正在获取统计信息...')
      const statsResponse = await codingProblemAPI.getStats()
      console.log('统计信息:', statsResponse)
      if (statsResponse?.data) {
        setStats(statsResponse.data)
      }

      // 获取每日一题历史
      console.log('正在获取每日一题历史...')
      await fetchDailyHistory()
    } catch (err) {
      console.error('获取数据失败:', err)
      setError(`获取数据失败: ${err.message}`)
    } finally {
      setLoading(false)
    }
  }

  // 获取新的随机题目
  const getNewRandomProblem = async () => {
    try {
      const response = await codingProblemAPI.getRandom()
      if (response?.data) {
        setRandomProblem(response.data)
      }
    } catch (err) {
      console.error('获取随机题目失败:', err)
    }
  }

  // 处理题目点击
  const handleProblemClick = (problem) => {
    if (problem?.source_url) {
      window.open(problem.source_url, '_blank')
    }
  }

  // 过滤题目
  useEffect(() => {
    console.log('正在过滤题目...')
    console.log('当前过滤条件:', { selectedDifficulty, selectedSource, searchQuery, activeTab })
    let filtered = [...problems]

    if (activeTab === 'daily') {
      filtered = dailyHistory.length > 0 ? dailyHistory : (dailyProblem ? [dailyProblem] : [])
    }

    if (selectedDifficulty !== 'all') {
      const difficultyMap = {
        'easy': ['简单', 'Easy'],
        'medium': ['中等', 'Medium'],
        'hard': ['困难', 'Hard']
      }
      const targetDifficulties = difficultyMap[selectedDifficulty] || []
      filtered = filtered.filter(p => targetDifficulties.includes(p.difficulty))
    }

    if (selectedSource !== 'all') {
      filtered = filtered.filter(p => p.source?.toLowerCase() === selectedSource)
    }

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(p =>
        p.title?.toLowerCase().includes(query) ||
        p.description?.toLowerCase().includes(query) ||
        p.tags?.some(tag => tag.toLowerCase().includes(query))
      )
    }

    console.log('过滤后的题目数量:', filtered.length)
    setFilteredProblems(filtered)
    setCurrentPage(1)
  }, [problems, selectedDifficulty, selectedSource, searchQuery, activeTab, dailyProblem, dailyHistory])

  // 初始加载
  useEffect(() => {
    console.log('组件已挂载，开始获取数据...')
    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="text-6xl mb-8 animate-bounce">🚀</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-4">加载中...</h2>
          <p className="text-gray-600">正在获取最新题目数据</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="text-6xl mb-8">😢</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-4">出错了</h2>
          <p className="text-gray-600 mb-8">{error}</p>
          <button
            onClick={fetchData}
            className="px-6 py-3 bg-purple-600 text-white rounded-xl font-semibold hover:bg-purple-700 transition-colors duration-300"
          >
            重试
          </button>
        </div>
      </div>
    )
  }

  // 如果没有数据，显示空状态
  if (!problems.length && !dailyProblem && !randomProblem) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="text-6xl mb-8">📝</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-4">暂无题目数据</h2>
          <p className="text-gray-600 mb-8">系统中还没有任何题目</p>
          <button
            onClick={fetchData}
            className="px-6 py-3 bg-purple-600 text-white rounded-xl font-semibold hover:bg-purple-700 transition-colors duration-300"
          >
            刷新
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      <div className="space-y-8">
        {/* 顶部装饰 */}
        <div className="absolute top-0 left-0 w-full h-96 bg-gradient-to-br from-violet-600 via-purple-600 to-pink-600 opacity-10 -skew-y-3 transform origin-top-left"></div>
        
        <div className="relative z-10 max-w-7xl mx-auto px-4 py-12">
          {/* Hero Section */}
          <div className="text-center mb-16">
            <div className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-r from-violet-500 via-purple-500 to-pink-500 rounded-3xl shadow-xl mb-8 transform hover:rotate-12 transition-transform duration-300">
              <span className="text-4xl">⚡</span>
            </div>
            <h1 className="text-6xl md:text-7xl font-black mb-6">
              <span className="bg-gradient-to-r from-violet-600 via-purple-600 to-pink-600 bg-clip-text text-transparent">
                LeetCode
              </span>
              <br />
              <span className="text-gray-800">刷题神器</span>
            </h1>
            <p className="text-xl md:text-2xl text-gray-600 max-w-3xl mx-auto leading-relaxed">
              💡 精选优质算法题目 · 🚀 一键直达LeetCode · ✨ 提升编程思维
            </p>
          </div>

          {/* 统计卡片 */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
            {[
              { label: '题目总数', value: stats.total || 0, icon: '📊', gradient: 'from-blue-500 to-purple-500' },
              { label: '简单题目', value: stats.easy || 0, icon: '🟢', gradient: 'from-emerald-500 to-green-500' },
              { label: '中等题目', value: stats.medium || 0, icon: '🟡', gradient: 'from-amber-500 to-orange-500' },
              { label: '困难题目', value: stats.hard || 0, icon: '🔴', gradient: 'from-red-500 to-pink-500' }
            ].map((stat, index) => (
              <div key={index} className="group relative">
                <div className="absolute -inset-1 bg-gradient-to-r from-violet-600 to-pink-600 rounded-2xl opacity-20 group-hover:opacity-30 transition-opacity"></div>
                <div className="relative bg-white rounded-2xl p-6 transform group-hover:-translate-y-2 transition-all duration-300 shadow-lg group-hover:shadow-2xl">
                  <div className="flex items-center justify-between mb-4">
                    <div className={`w-12 h-12 bg-gradient-to-r ${stat.gradient} rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-300`}>
                      <span className="text-xl">{stat.icon}</span>
                    </div>
                    <div className="text-right">
                      <div className="text-3xl font-bold text-gray-800">{stat.value}</div>
                      <div className="text-sm text-gray-500">{stat.label}</div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* 每日一题和随机一题 */}
          <div className="flex flex-wrap gap-6">
            {/* 每日一题卡片 */}
            <div className="flex-1 min-w-[300px] bg-gradient-to-br from-violet-600 to-purple-600 rounded-2xl p-6 text-white">
              <div className="flex items-center gap-3 mb-4">
                <span className="text-2xl">⭐</span>
                <div>
                  <h3 className="text-xl font-bold">今日精选</h3>
                  <p className="text-sm opacity-80">每日一题，坚持成长</p>
                </div>
              </div>
              {dailyProblem ? (
                <div className="space-y-4">
                  <h4 className="text-lg font-semibold">{dailyProblem.title}</h4>
                  <div className="flex gap-2 flex-wrap">
                    <span className="px-3 py-1 bg-white/20 rounded-full text-sm">
                      {dailyProblem.difficulty}
                    </span>
                    {dailyProblem.tags?.map((tag, index) => (
                      <span key={index} className="px-3 py-1 bg-white/20 rounded-full text-sm">
                        {tag}
                      </span>
                    ))}
                  </div>
                  <button
                    onClick={() => handleProblemClick(dailyProblem)}
                    className="w-full py-3 bg-white text-purple-600 rounded-xl font-semibold hover:bg-purple-50 transition-colors duration-300"
                  >
                    开始刷题 🚀
                  </button>
                </div>
              ) : error ? (
                <div className="text-center py-8">
                  <p className="text-white/80 mb-4">{error}</p>
                  <button
                    onClick={fetchData}
                    className="px-6 py-2 bg-white/20 rounded-xl hover:bg-white/30 transition-colors duration-300"
                  >
                    重试
                  </button>
                </div>
              ) : (
                <div className="text-center py-8">
                  <div className="animate-spin text-4xl mb-4">🔄</div>
                  <p className="text-white/80">正在获取每日一题...</p>
                </div>
              )}
            </div>

            {/* 随机一题卡片 */}
            <div className="flex-1 min-w-[300px] bg-gradient-to-br from-pink-600 to-rose-600 rounded-2xl p-6 text-white">
              <div className="flex items-center gap-3 mb-4">
                <span className="text-2xl">🎲</span>
                <div>
                  <h3 className="text-xl font-bold">随机一题</h3>
                  <p className="text-sm opacity-80">挑战自己，突破瓶颈</p>
                </div>
              </div>
              {randomProblem ? (
                <div className="space-y-4">
                  <h4 className="text-lg font-semibold">{randomProblem.title}</h4>
                  <div className="flex gap-2">
                    <span className="px-3 py-1 bg-white/20 rounded-full text-sm">
                      {randomProblem.difficulty}
                    </span>
                    {randomProblem.tags?.map((tag, index) => (
                      <span key={index} className="px-3 py-1 bg-white/20 rounded-full text-sm">
                        {tag}
                      </span>
                    ))}
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleProblemClick(randomProblem)}
                      className="flex-1 py-3 bg-white text-rose-600 rounded-xl font-semibold hover:bg-rose-50 transition-colors duration-300"
                    >
                      开始刷题 🚀
                    </button>
                    <button
                      onClick={getNewRandomProblem}
                      className="px-4 py-3 bg-white/20 rounded-xl font-semibold hover:bg-white/30 transition-colors duration-300"
                    >
                      换一题 🔄
                    </button>
                  </div>
                </div>
              ) : (
                <div className="text-center py-8">
                  <p>暂无随机题目</p>
                </div>
              )}
            </div>
          </div>

          {/* 主要内容 */}
          <div className="bg-white rounded-3xl shadow-2xl overflow-hidden">
            {/* 标签页 */}
            <div className="bg-gradient-to-r from-gray-50 to-white p-8 border-b border-gray-100">
              <div className="flex flex-col lg:flex-row lg:items-center justify-between gap-8">
                <div className="flex gap-2 bg-gray-100 rounded-2xl p-2">
                  {[
                    { key: 'all', label: '全部题目', icon: '📚' },
                    { key: 'daily', label: '每日一题', icon: '⭐' }
                  ].map(tab => (
                    <button
                      key={tab.key}
                      onClick={() => {
                        setActiveTab(tab.key)
                        if (tab.key === 'daily') {
                          fetchDailyHistory()
                        }
                      }}
                      className={`px-6 py-3 rounded-xl font-semibold transition-all duration-300 flex items-center gap-3 ${
                        activeTab === tab.key
                          ? 'bg-white text-purple-600 shadow-lg transform scale-105'
                          : 'text-gray-600 hover:text-gray-800 hover:bg-white hover:bg-opacity-50'
                      }`}
                    >
                      <span className="text-lg">{tab.icon}</span>
                      {tab.label}
                    </button>
                  ))}
                </div>
              </div>

              {/* 搜索和筛选 */}
              {activeTab === 'all' && (
                <div className="flex flex-col lg:flex-row gap-4 mt-8">
                  <div className="flex-1 relative">
                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                      <svg className="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                      </svg>
                    </div>
                    <input
                      type="text"
                      placeholder="🔍 搜索题目、标签或描述..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="w-full pl-12 pr-6 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 transition-all duration-300 text-lg"
                    />
                  </div>

                  <select
                    value={selectedDifficulty}
                    onChange={(e) => setSelectedDifficulty(e.target.value)}
                    className="px-6 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 transition-all duration-300 font-medium text-gray-700 bg-white"
                  >
                    <option value="all">📊 全部难度</option>
                    <option value="easy">🟢 简单</option>
                    <option value="medium">🟡 中等</option>
                    <option value="hard">🔴 困难</option>
                  </select>

                  <select
                    value={selectedSource}
                    onChange={(e) => setSelectedSource(e.target.value)}
                    className="px-6 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 transition-all duration-300 font-medium text-gray-700 bg-white"
                  >
                    <option value="all">🌐 全部来源</option>
                    {uniqueSources.map(source => (
                      <option key={source} value={source.toLowerCase()}>{source}</option>
                    ))}
                  </select>
                </div>
              )}
            </div>

            {/* 题目列表 */}
            <div className="p-8">
              {/* 统计信息 */}
              {activeTab === 'all' && stats && (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                  {[
                    { label: '总题目数', value: stats.total || 0, icon: '📚', color: 'from-blue-500 to-cyan-400' },
                    { label: '简单题目', value: stats.easy || 0, icon: '🟢', color: 'from-green-500 to-emerald-400' },
                    { label: '中等题目', value: stats.medium || 0, icon: '🟡', color: 'from-yellow-500 to-amber-400' },
                    { label: '困难题目', value: stats.hard || 0, icon: '🔴', color: 'from-red-500 to-rose-400' }
                  ].map((stat, index) => (
                    <div key={index} className={`bg-gradient-to-r ${stat.color} rounded-3xl p-6 text-white`}>
                      <div className="flex items-center gap-4">
                        <div className="w-12 h-12 bg-white bg-opacity-20 backdrop-blur-sm rounded-xl flex items-center justify-center">
                          <span className="text-2xl">{stat.icon}</span>
                        </div>
                        <div>
                          <h3 className="text-lg font-medium text-white text-opacity-90">{stat.label}</h3>
                          <p className="text-3xl font-bold">{stat.value}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* 题目列表 */}
              <div className="space-y-6">
                {currentProblems.map((problem) => (
                  <div
                    key={problem.id}
                    className="bg-white rounded-2xl p-6 shadow-lg hover:shadow-xl transition-shadow duration-300"
                  >
                    <div className="flex flex-col lg:flex-row lg:items-center gap-4 lg:gap-6">
                      <div className="flex-1">
                        <h3 className="text-xl font-bold text-gray-900 mb-2">{problem.title}</h3>
                        {problem.description && (
                          <p className="text-gray-600 mb-4 line-clamp-2">{problem.description}</p>
                        )}
                        <div className="flex flex-wrap gap-2">
                          <span className={`px-3 py-1 rounded-full text-sm font-medium ${getDifficultyConfig(problem.difficulty).textColor}`}>
                            {getDifficultyConfig(problem.difficulty).text}
                          </span>
                          {problem.source && (
                            <span className="px-3 py-1 bg-gray-100 rounded-full text-sm font-medium text-gray-600">
                              {problem.source}
                            </span>
                          )}
                          {problem.tags && problem.tags.map((tag, index) => (
                            <span key={index} className="px-3 py-1 bg-purple-50 text-purple-600 rounded-full text-sm font-medium">
                              #{tag}
                            </span>
                          ))}
                        </div>
                      </div>
                      <div className="flex items-center gap-4">
                        <button
                          onClick={() => handleProblemClick(problem)}
                          className="px-6 py-2 bg-purple-600 text-white rounded-xl font-semibold hover:bg-purple-700 transition-colors duration-300"
                        >
                          开始
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* 分页 */}
              {activeTab === 'all' && filteredProblems.length > itemsPerPage && (
                <div className="mt-8 flex justify-center">
                  <div className="flex gap-2">
                    {Array.from({ length: Math.ceil(filteredProblems.length / itemsPerPage) }, (_, i) => (
                      <button
                        key={i}
                        onClick={() => setCurrentPage(i + 1)}
                        className={`w-10 h-10 rounded-xl font-semibold transition-all duration-300 ${
                          currentPage === i + 1
                            ? 'bg-purple-600 text-white'
                            : 'text-gray-600 hover:bg-purple-100'
                        }`}
                      >
                        {i + 1}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CodingDashboard 