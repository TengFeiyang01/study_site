package web

import (
	"github.com/gin-gonic/gin"
)

// WebHandler 统一的Web处理器，简化为只有题目显示功能
type WebHandler struct {
	questHandler         *QuestHandler
	codingProblemHandler *CodingProblemHandler
}

func NewWebHandler(
	questHandler *QuestHandler,
	codingProblemHandler *CodingProblemHandler,
) *WebHandler {
	return &WebHandler{
		questHandler:         questHandler,
		codingProblemHandler: codingProblemHandler,
	}
}

func (h *WebHandler) RegisterRoutes(server *gin.Engine) {
	// 八股题目路由
	h.questHandler.RegisterRoutes(server)

	// 刷题模块路由 - 简化为只显示题目，支持跳转到LeetCode
	h.codingProblemHandler.RegisterRoutes(server)
}
