package service

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/repository"
	"context"
	"time"
)

type CodingProblemService interface {
	CreateProblem(ctx context.Context, problem domain.CodingProblem) error
	GetAllProblems(ctx context.Context) ([]domain.CodingProblem, error)
	GetProblemById(ctx context.Context, id int64) (domain.CodingProblem, error)
	GetProblemsBySource(ctx context.Context, source string) ([]domain.CodingProblem, error)
	GetProblemsByDifficulty(ctx context.Context, difficulty string) ([]domain.CodingProblem, error)
	UpdateProblem(ctx context.Context, problem domain.CodingProblem) error
	DeleteProblem(ctx context.Context, id int64) error
	GetDailyProblem(ctx context.Context) (*domain.CodingProblem, error)
	GetDailyProblemHistory(ctx context.Context) ([]domain.CodingProblem, error)
	SetDailyProblem(ctx context.Context, problemId int64) error
	// 爬虫相关方法
	CrawlDailyProblem(ctx context.Context) error
	CrawlProblemBySlug(ctx context.Context, titleSlug string) (*domain.CodingProblem, error)
	StartDailyCrawler(ctx context.Context) error
}

type codingProblemService struct {
	repo    repository.CodingProblemRepository
	crawler *LeetCodeCrawler
}

func NewCodingProblemService(repo repository.CodingProblemRepository, crawler *LeetCodeCrawler) CodingProblemService {
	return &codingProblemService{
		repo:    repo,
		crawler: crawler,
	}
}

// CodingProblem Service 实现
func (svc *codingProblemService) CreateProblem(ctx context.Context, problem domain.CodingProblem) error {
	problem.Ctime = time.Now()
	problem.Utime = time.Now()
	return svc.repo.Create(ctx, problem)
}

func (svc *codingProblemService) GetAllProblems(ctx context.Context) ([]domain.CodingProblem, error) {
	return svc.repo.FindAll(ctx)
}

func (svc *codingProblemService) GetProblemById(ctx context.Context, id int64) (domain.CodingProblem, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *codingProblemService) GetProblemsBySource(ctx context.Context, source string) ([]domain.CodingProblem, error) {
	return svc.repo.FindBySource(ctx, source)
}

func (svc *codingProblemService) GetProblemsByDifficulty(ctx context.Context, difficulty string) ([]domain.CodingProblem, error) {
	return svc.repo.FindByDifficulty(ctx, difficulty)
}

func (svc *codingProblemService) UpdateProblem(ctx context.Context, problem domain.CodingProblem) error {
	problem.Utime = time.Now()
	return svc.repo.Update(ctx, problem)
}

func (svc *codingProblemService) DeleteProblem(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}

// 每日一题相关方法
func (svc *codingProblemService) GetDailyProblem(ctx context.Context) (*domain.CodingProblem, error) {
	return svc.repo.GetDailyProblem(ctx)
}

func (svc *codingProblemService) GetDailyProblemHistory(ctx context.Context) ([]domain.CodingProblem, error) {
	return svc.repo.GetDailyProblemHistory(ctx)
}

func (svc *codingProblemService) SetDailyProblem(ctx context.Context, problemId int64) error {
	return svc.repo.SetDailyProblem(ctx, problemId)
}

// 爬虫相关方法实现
func (svc *codingProblemService) CrawlDailyProblem(ctx context.Context) error {
	return svc.crawler.CrawlAndSaveDailyProblem(ctx)
}

func (svc *codingProblemService) CrawlProblemBySlug(ctx context.Context, titleSlug string) (*domain.CodingProblem, error) {
	return svc.crawler.CrawlProblemBySlug(ctx, titleSlug)
}

func (svc *codingProblemService) StartDailyCrawler(ctx context.Context) error {
	return svc.crawler.StartDailyCrawler(ctx)
}
