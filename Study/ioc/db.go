package ioc

import (
	"Training/Study/internal/repository"
	"Training/Study/internal/repository/dao"
	"Training/Study/internal/service"
	"Training/Study/internal/web"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Application 应用程序依赖
type Application struct {
	DB                   *gorm.DB
	QuestHandler         *web.QuestHandler
	CodingProblemHandler *web.CodingProblemHandler
	Crawler              *service.LeetCodeCrawler
	CodingProblemRepo    repository.CodingProblemRepository
}

func InitDB() *gorm.DB {
	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // Log level
		},
	)

	// 连接MySQL数据库
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3306)/study?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect database")
	}

	return db
}

func InitWebServer() *gin.Engine {
	server := gin.Default()

	// 添加CORS中间件
	server.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	return server
}

// InitApplication 初始化应用程序依赖
func InitApplication(db *gorm.DB) *Application {
	// 初始化DAO
	questDAO := dao.NewQuestionDao(db)
	codingProblemDAO := dao.NewGormCodingProblemDAO(db)

	// 初始化Repository
	questRepo := repository.NewQuestRepository(questDAO)
	codingProblemRepo := repository.NewCachedCodingProblemRepository(codingProblemDAO)

	// 初始化Service
	leetcodeCrawler := service.NewLeetCodeCrawler(codingProblemRepo)
	questService := service.NewQuestService(questRepo)
	codingProblemService := service.NewCodingProblemService(codingProblemRepo, leetcodeCrawler)

	// 初始化Handler
	questHandler := web.NewQuestHandler(questService)
	codingProblemHandler := web.NewCodingProblemHandler(codingProblemService)

	return &Application{
		DB:                   db,
		QuestHandler:         questHandler,
		CodingProblemHandler: codingProblemHandler,
		Crawler:              leetcodeCrawler,
		CodingProblemRepo:    codingProblemRepo,
	}
}
