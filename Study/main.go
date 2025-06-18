package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Training/Study/config"
	"Training/Study/internal/repository/dao"
	"Training/Study/ioc"
)

func main() {
	log.Printf("🚀 开始启动应用...")

	// 初始化依赖注入
	app := initApplication()

	// 初始化数据
	log.Printf("📊 开始初始化数据...")
	initializeData(app)

	// 启动web服务器
	startWebServer(app)
}

// 初始化应用
func initApplication() *ioc.Application {
	log.Printf("🔧 初始化数据库连接...")
	// 初始化数据库
	db := ioc.InitDB()

	// 自动创建表
	log.Printf("📋 创建数据库表...")
	err := dao.InitTables(db)
	if err != nil {
		panic("Failed to create tables: " + err.Error())
	}

	log.Printf("🏗️  初始化应用组件...")
	return ioc.InitApplication(db)
}

// 初始化数据
func initializeData(app *ioc.Application) {
	ctx := context.Background()

	// 1. 插入Hot100题目数据
	log.Printf("🔥 开始插入Hot100题目数据...")
	config.InsertHot100Problems(ctx, app)

	// 2. 获取每日一题（启动时爬取一次）
	log.Printf("📅 开始获取每日一题...")
	err := app.Crawler.CrawlAndSaveDailyProblem(ctx)
	if err != nil {
		log.Printf("❌ 获取每日一题失败: %v", err)
	} else {
		log.Printf("✅ 每日一题获取完成")
	}

	// 3. 启动定时任务（每天0点获取新的每日一题）
	log.Printf("⏰ 启动定时任务...")
	go func() {
		if err := app.Crawler.StartDailyCrawler(ctx); err != nil {
			log.Printf("❌ 定时爬虫停止: %v", err)
		}
	}()
}

// 启动web服务器
func startWebServer(app *ioc.Application) {
	log.Printf("🌐 启动Web服务器...")

	// 启动服务器
	server := ioc.InitWebServer()

	// 注册路由
	app.QuestHandler.RegisterRoutes(server)
	app.CodingProblemHandler.RegisterRoutes(server)

	// 启动服务器
	port := ":8080"
	log.Printf("🚀 服务器启动在 http://localhost%s", port)

	// 创建HTTP服务器
	httpServer := &http.Server{
		Addr:    port,
		Handler: server,
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("❌ 服务器强制关闭:", err)
	}

	log.Println("✅ 服务器已退出")
}
