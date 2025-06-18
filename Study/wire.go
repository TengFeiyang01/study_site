//go:build wireinject

package main

import (
	"Training/Study/internal/repository"
	"Training/Study/internal/repository/dao"
	"Training/Study/internal/service"
	"Training/Study/internal/web"
	"Training/Study/ioc"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 数据库
		ioc.InitDB,

		// DAO层
		dao.NewQuestionDao,
		dao.NewGormCodingProblemDAO,

		// Repository层
		repository.NewQuestRepository,
		repository.NewCachedCodingProblemRepository,

		// Service层
		service.NewQuestService,
		service.NewLeetCodeCrawler,
		service.NewCodingProblemService,

		// Handler层
		web.NewQuestHandler,
		web.NewCodingProblemHandler,

		// Web服务器
		InitGinServer,
	)
	return &gin.Engine{}
}

func InitGinServer(
	questionHandler *web.QuestHandler,
	codingHandler *web.CodingProblemHandler,
	codingService service.CodingProblemService,
) *gin.Engine {
	server := gin.Default()

	// 配置CORS
	server.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 预热缓存并启动每日一题爬虫
	go func() {
		ctx := context.Background()
		if problems, err := codingService.GetAllProblems(ctx); err != nil {
			log.Printf("Failed to warm up cache: %v", err)
		} else {
			log.Printf("Cache warmed up successfully with %d problems", len(problems))
		}

		// 启动每日一题定时爬虫
		go func() {
			if err := codingService.StartDailyCrawler(ctx); err != nil {
				log.Printf("Daily crawler stopped: %v", err)
			}
		}()
	}()

	// 注册路由 - 八股复习 + 刷题模块
	questionHandler.RegisterRoutes(server)
	codingHandler.RegisterRoutes(server)

	return server
}
