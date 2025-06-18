package repository

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/repository/dao"

	"github.com/ecodeclub/ekit/slice"

	"context"
	"time"
)

type QuestRepository interface {
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

type questRepository struct {
	dao dao.QuestDao
}

func NewQuestRepository(dao dao.QuestDao) QuestRepository {
	return &questRepository{dao: dao}
}

func (r *questRepository) Insert(ctx context.Context, quest domain.Question) error {
	return r.dao.Insert(ctx, r.toEntity(quest))
}

func (r *questRepository) FindByCategory(ctx context.Context, category string) ([]domain.Question, error) {
	res, err := r.dao.FindByCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	data := slice.Map(res, func(idx int, src dao.Question) domain.Question {
		return r.toDomain(src)
	})
	return data, err
}

func (r *questRepository) FindAll(ctx context.Context) ([]domain.Question, error) {
	res, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	data := slice.Map(res, func(idx int, src dao.Question) domain.Question {
		return r.toDomain(src)
	})
	return data, err
}

func (r *questRepository) UpdateById(ctx context.Context, quest domain.Question) error {
	return r.dao.UpdateById(ctx, r.toEntity(quest))
}

func (r *questRepository) UpdateMasteryLevel(ctx context.Context, id int64, masteryLevel int) error {
	return r.dao.UpdateMasteryLevel(ctx, id, masteryLevel)
}

func (r *questRepository) DeleteById(ctx context.Context, id int64) error {
	return r.dao.DeleteById(ctx, id)
}

func (r *questRepository) DeleteByCategory(ctx context.Context, category string) error {
	return r.dao.DeleteByCategory(ctx, category)
}

func (r *questRepository) FindAllCategories(ctx context.Context) ([]string, error) {
	return r.dao.FindAllCategories(ctx)
}

func (r *questRepository) GetMasteryStats(ctx context.Context) (map[string]int, error) {
	return r.dao.GetMasteryStats(ctx)
}

func (r *questRepository) toDomain(quest dao.Question) domain.Question {
	return domain.Question{
		Id:           quest.Id,
		Category:     quest.Category,
		Content:      quest.Content,
		Answer:       quest.Answer,
		MasteryLevel: quest.MasteryLevel,
		Ctime:        time.UnixMilli(quest.Ctime),
		Utime:        time.UnixMilli(quest.Utime),
	}
}

func (r *questRepository) toEntity(quest domain.Question) dao.Question {
	return dao.Question{
		Id:           quest.Id,
		Category:     quest.Category,
		Content:      quest.Content,
		Answer:       quest.Answer,
		MasteryLevel: quest.MasteryLevel,
		Ctime:        quest.Ctime.UnixMilli(),
		Utime:        quest.Utime.UnixMilli(),
	}
}
