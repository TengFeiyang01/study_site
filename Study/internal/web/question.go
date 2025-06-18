package web

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestHandler struct {
	svc service.QuestService
}

func NewQuestHandler(svc service.QuestService) *QuestHandler {
	return &QuestHandler{
		svc: svc,
	}
}

func (q *QuestHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/question")
	g.GET("/", q.FindAll)
	g.GET("/categories", q.FindAllCategories)
	g.GET("/mastery-stats", q.GetMasteryStats)
	g.GET("/:category", q.FindByCategory)
	g.POST("/", q.Insert)
	g.PUT("/:id", q.UpdateById)
	g.PUT("/:id/mastery", q.UpdateMasteryLevel)
	g.DELETE("/:id", q.DeleteById)
}

func (q *QuestHandler) Insert(ctx *gin.Context) {
	type Request struct {
		Content  string `json:"content" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := q.svc.Insert(ctx, domain.Question{
		Content:  req.Content,
		Answer:   req.Answer,
		Category: req.Category,
	}); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})
}

func (q *QuestHandler) FindAll(ctx *gin.Context) {
	questions, err := q.svc.FindAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, questions)
}

func (q *QuestHandler) FindByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	if category == "" {
		ctx.JSON(400, gin.H{"error": "category is required"})
		return
	}
	questions, err := q.svc.FindByCategory(ctx, category)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, questions)
}

func (q *QuestHandler) UpdateById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	type BodyRequest struct {
		Content  string `json:"content" binding:"required"`
		Answer   string `json:"answer" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	var bodyReq BodyRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := q.svc.UpdateById(ctx, domain.Question{
		Id:       id,
		Category: bodyReq.Category,
		Content:  bodyReq.Content,
		Answer:   bodyReq.Answer,
	}); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})
}

func (q *QuestHandler) UpdateMasteryLevel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	type Request struct {
		MasteryLevel int `json:"mastery_level" binding:"required,min=0,max=2"`
	}
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := q.svc.UpdateMasteryLevel(ctx, id, req.MasteryLevel); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})
}

func (q *QuestHandler) DeleteById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	if err := q.svc.DeleteById(ctx, id); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "success"})
}

func (q *QuestHandler) FindAllCategories(ctx *gin.Context) {
	categories, err := q.svc.FindAllCategories(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, categories)
}

func (q *QuestHandler) GetMasteryStats(ctx *gin.Context) {
	stats, err := q.svc.GetMasteryStats(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, stats)
}
