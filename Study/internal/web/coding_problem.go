package web

import (
	"Training/Study/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CodingProblemHandler struct {
	service service.CodingProblemService
}

func NewCodingProblemHandler(service service.CodingProblemService) *CodingProblemHandler {
	return &CodingProblemHandler{
		service: service,
	}
}

func (h *CodingProblemHandler) RegisterRoutes(server *gin.Engine) {
	codingGroup := server.Group("/api/coding")

	// 核心功能：显示题目列表，支持直接跳转到LeetCode
	codingGroup.GET("/problems", h.GetAllProblems)
	codingGroup.GET("/problems/:id", h.GetProblemById)
	codingGroup.GET("/problems/source/:source", h.GetProblemsBySource)
	codingGroup.GET("/problems/difficulty/:difficulty", h.GetProblemsByDifficulty)
	codingGroup.GET("/daily", h.GetDailyProblem)
	codingGroup.GET("/daily/history", h.GetDailyProblemHistory)
	codingGroup.GET("/random", h.GetRandomProblem)
	codingGroup.GET("/stats", h.GetStats)

	// 学习状态管理
	codingGroup.PUT("/problems/:id/study-status", h.UpdateStudyStatus)

	// 管理功能
	codingGroup.POST("/crawl-hot100", h.CrawlHot100Problems)
}

func (h *CodingProblemHandler) GetAllProblems(c *gin.Context) {
	// 支持分页参数
	page := 1
	limit := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	problems, err := h.service.GetAllProblems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 简单分页处理
	total := len(problems)
	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		c.JSON(http.StatusOK, gin.H{
			"problems": []interface{}{},
			"total":    total,
			"page":     page,
			"limit":    limit,
		})
		return
	}

	if end > total {
		end = total
	}

	pagedProblems := problems[start:end]

	c.JSON(http.StatusOK, gin.H{
		"problems": pagedProblems,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *CodingProblemHandler) GetProblemById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := h.service.GetProblemById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func (h *CodingProblemHandler) GetProblemsBySource(c *gin.Context) {
	source := c.Param("source")
	problems, err := h.service.GetProblemsBySource(c.Request.Context(), source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"total":    len(problems),
	})
}

func (h *CodingProblemHandler) GetProblemsByDifficulty(c *gin.Context) {
	difficulty := c.Param("difficulty")
	problems, err := h.service.GetProblemsByDifficulty(c.Request.Context(), difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"total":    len(problems),
	})
}

func (h *CodingProblemHandler) GetDailyProblem(c *gin.Context) {
	problem, err := h.service.GetDailyProblem(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if problem == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No daily problem found"})
		return
	}
	c.JSON(http.StatusOK, problem)
}

func (h *CodingProblemHandler) UpdateStudyStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	var req struct {
		StudyStatus string `json:"study_status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证学习状态值
	validStatuses := map[string]bool{
		"not_started": true,
		"in_progress": true,
		"completed":   true,
	}

	if !validStatuses[req.StudyStatus] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study status"})
		return
	}

	// 获取当前题目
	problem, err := h.service.GetProblemById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新学习状态
	problem.StudyStatus = req.StudyStatus
	if req.StudyStatus != "not_started" {
		now := time.Now()
		problem.LastStudied = &now
	}

	err = h.service.UpdateProblem(c.Request.Context(), problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study status updated successfully"})
}

// GetDailyProblemHistory 获取每日一题历史
func (h *CodingProblemHandler) GetDailyProblemHistory(c *gin.Context) {
	problems, err := h.service.GetDailyProblemHistory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"total":    len(problems),
	})
}

// GetRandomProblem 获取随机题目
func (h *CodingProblemHandler) GetRandomProblem(c *gin.Context) {
	problems, err := h.service.GetAllProblems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(problems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No problems found"})
		return
	}

	// 简单随机选择
	randomIndex := time.Now().UnixNano() % int64(len(problems))
	randomProblem := problems[randomIndex]

	c.JSON(http.StatusOK, randomProblem)
}

// GetStats 获取统计信息
func (h *CodingProblemHandler) GetStats(c *gin.Context) {
	problems, err := h.service.GetAllProblems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 初始化统计数据
	easyCount := 0
	mediumCount := 0
	hardCount := 0
	sourceMap := make(map[string]int)
	statusMap := map[string]int{
		"not_started": 0,
		"in_progress": 0,
		"completed":   0,
	}

	for _, problem := range problems {
		// 按难度统计 (统一格式)
		switch problem.Difficulty {
		case "Easy", "简单":
			easyCount++
		case "Medium", "中等":
			mediumCount++
		case "Hard", "困难":
			hardCount++
		}

		// 按来源统计
		sourceMap[problem.Source]++

		// 按状态统计
		statusMap[problem.StudyStatus]++
	}

	stats := map[string]interface{}{
		"total":  len(problems),
		"easy":   easyCount,
		"medium": mediumCount,
		"hard":   hardCount,
		"by_difficulty": map[string]int{
			"Easy":   easyCount,
			"Medium": mediumCount,
			"Hard":   hardCount,
		},
		"by_source": sourceMap,
		"by_status": statusMap,
	}

	c.JSON(http.StatusOK, stats)
}

// CrawlHot100Problems 爬取Hot100题目
func (h *CodingProblemHandler) CrawlHot100Problems(c *gin.Context) {
	go func() {
		if err := h.service.CrawlAndSaveHot100Problems(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Hot100题目爬取已开始，请稍候查看结果"})
}
