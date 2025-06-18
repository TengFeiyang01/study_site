package service

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/repository"
	"context"
)

type QuestService interface {
	Insert(ctx context.Context, quest domain.Question) error
	FindByCategory(ctx context.Context, category string) ([]domain.Question, error)
	FindAll(ctx context.Context) ([]domain.Question, error)
	UpdateById(ctx context.Context, quest domain.Question) error
	UpdateMasteryLevel(ctx context.Context, id int64, masteryLevel int) error
	DeleteById(ctx context.Context, id int64) error
	DeleteByCategory(ctx context.Context, category string) error
	FindAllCategories(ctx context.Context) ([]string, error)
	GetMasteryStats(ctx context.Context) (map[string]int, error)
}

type questService struct {
	repo repository.QuestRepository
}

func NewQuestService(repo repository.QuestRepository) QuestService {
	return &questService{repo: repo}
}

func (svc *questService) Insert(ctx context.Context, quest domain.Question) error {
	return svc.repo.Insert(ctx, quest)
}

func (svc *questService) FindByCategory(ctx context.Context, category string) ([]domain.Question, error) {
	return svc.repo.FindByCategory(ctx, category)
}

func (svc *questService) FindAll(ctx context.Context) ([]domain.Question, error) {
	return svc.repo.FindAll(ctx)
}

func (svc *questService) UpdateById(ctx context.Context, quest domain.Question) error {
	return svc.repo.UpdateById(ctx, quest)
}

func (svc *questService) UpdateMasteryLevel(ctx context.Context, id int64, masteryLevel int) error {
	return svc.repo.UpdateMasteryLevel(ctx, id, masteryLevel)
}

func (svc *questService) DeleteById(ctx context.Context, id int64) error {
	return svc.repo.DeleteById(ctx, id)
}

func (svc *questService) DeleteByCategory(ctx context.Context, category string) error {
	return svc.repo.DeleteByCategory(ctx, category)
}

func (svc *questService) FindAllCategories(ctx context.Context) ([]string, error) {
	return svc.repo.FindAllCategories(ctx)
}

func (svc *questService) GetMasteryStats(ctx context.Context) (map[string]int, error) {
	return svc.repo.GetMasteryStats(ctx)
}
