import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080', // 明确指定后端地址
  timeout: 15000, // 增加到15秒
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    console.log('Request:', config.method?.toUpperCase(), config.url, config.data)
    return config
  },
  error => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    console.log('Response:', response.status, response.data)
    // 统一处理响应数据格式
    const data = response.data
    if (data.error) {
      return Promise.reject(new Error(data.error))
    }
    return {
      data: data.problems || data.stats || data,
      total: data.total,
      page: data.page,
      limit: data.limit
    }
  },
  error => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

// Question API
export const questionAPI = {
  // 获取所有问题
  getAll: () => api.get('/question/'),
  
  // 获取所有分类
  getCategories: () => api.get('/question/categories'),
  
  // 获取掌握度统计
  getMasteryStats: () => api.get('/question/mastery-stats'),
  
  // 按分类获取问题
  getByCategory: (category) => api.get(`/question/${category}`),
  
  // 创建问题
  create: (data) => api.post('/question/', data),
  
  // 更新问题
  update: (id, data) => api.put(`/question/${id}`, data),
  
  // 更新掌握度
  updateMasteryLevel: (id, masteryLevel) => api.put(`/question/${id}/mastery`, { mastery_level: masteryLevel }),
  
  // 删除问题
  delete: (id) => api.delete(`/question/${id}`)
}

// 简化的刷题问题API - 只保留查看和跳转功能  
export const codingProblemAPI = {
  // 核心查看功能
  getAll: (params) => api.get('/api/coding/problems', { params }),
  getBySource: (source) => api.get(`/api/coding/problems/source/${source}`),
  getByDifficulty: (difficulty) => api.get(`/api/coding/problems/difficulty/${difficulty}`),
  getDaily: () => api.get('/api/coding/daily'),
  getDailyHistory: () => api.get('/api/coding/daily/history'),
  getRandom: () => api.get('/api/coding/random'),
  getStats: () => api.get('/api/coding/stats'),
  
  // 学习状态管理
  updateStudyStatus: (problemId, status) => api.put(`/api/coding/problems/${problemId}/study-status`, { study_status: status }),
  
  // 管理功能
  refreshCache: () => api.post('/api/coding/refresh'),
  importProblems: (problems) => api.post('/api/coding/import', problems)
}

// 移除了所有刷题记录相关API，现在专注于题目展示和跳转LeetCode

export default api
