package repository

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/repository/dao"
	"context"
	"time"
)

type CodingProblemRepository interface {
	Create(ctx context.Context, problem domain.CodingProblem) error
	FindAll(ctx context.Context) ([]domain.CodingProblem, error)
	FindById(ctx context.Context, id int64) (domain.CodingProblem, error)
	FindBySource(ctx context.Context, source string) ([]domain.CodingProblem, error)
	FindBySourceId(ctx context.Context, sourceId string) ([]domain.CodingProblem, error)
	FindByDifficulty(ctx context.Context, difficulty string) ([]domain.CodingProblem, error)
	Update(ctx context.Context, problem domain.CodingProblem) error
	Delete(ctx context.Context, id int64) error
	GetDailyProblem(ctx context.Context) (*domain.CodingProblem, error)
	SetDailyProblem(ctx context.Context, problemId int64) error
	GetDailyProblemHistory(ctx context.Context) ([]domain.CodingProblem, error)
	MarkAsDailyProblem(ctx context.Context, problemId int64, date time.Time) error
	UpdateStudyStatus(ctx context.Context, problemId int64, status string, lastStudied *time.Time) error
	SaveDailyProblem(dailyProblem *domain.DailyProblem) error
}

type CachedCodingProblemRepository struct {
	dao dao.CodingProblemDAO
}

func NewCachedCodingProblemRepository(dao dao.CodingProblemDAO) CodingProblemRepository {
	return &CachedCodingProblemRepository{dao: dao}
}

// CodingProblem Repository 实现
func (r *CachedCodingProblemRepository) Create(ctx context.Context, problem domain.CodingProblem) error {
	return r.dao.Insert(ctx, toEntity(problem))
}

func (r *CachedCodingProblemRepository) FindAll(ctx context.Context) ([]domain.CodingProblem, error) {
	problems, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CodingProblem, 0, len(problems))
	for _, p := range problems {
		result = append(result, r.toDomain(p))
	}
	return result, nil
}

func (r *CachedCodingProblemRepository) FindById(ctx context.Context, id int64) (domain.CodingProblem, error) {
	problem, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.CodingProblem{}, err
	}
	return r.toDomain(problem), nil
}

func (r *CachedCodingProblemRepository) FindBySource(ctx context.Context, source string) ([]domain.CodingProblem, error) {
	problems, err := r.dao.FindBySource(ctx, source)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CodingProblem, 0, len(problems))
	for _, p := range problems {
		result = append(result, r.toDomain(p))
	}
	return result, nil
}

func (r *CachedCodingProblemRepository) FindBySourceId(ctx context.Context, sourceId string) ([]domain.CodingProblem, error) {
	problems, err := r.dao.FindBySourceId(ctx, sourceId)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CodingProblem, 0, len(problems))
	for _, p := range problems {
		result = append(result, r.toDomain(p))
	}
	return result, nil
}

func (r *CachedCodingProblemRepository) FindByDifficulty(ctx context.Context, difficulty string) ([]domain.CodingProblem, error) {
	problems, err := r.dao.FindByDifficulty(ctx, difficulty)
	if err != nil {
		return nil, err
	}
	result := make([]domain.CodingProblem, 0, len(problems))
	for _, p := range problems {
		result = append(result, r.toDomain(p))
	}
	return result, nil
}

func (r *CachedCodingProblemRepository) Update(ctx context.Context, problem domain.CodingProblem) error {
	return r.dao.UpdateById(ctx, toEntity(problem))
}

func (r *CachedCodingProblemRepository) Delete(ctx context.Context, id int64) error {
	return r.dao.DeleteById(ctx, id)
}

func (c *CachedCodingProblemRepository) GetDailyProblem(ctx context.Context) (*domain.CodingProblem, error) {
	problem, err := c.dao.GetDailyProblem(ctx)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, nil
	}
	result := c.toDomain(*problem)
	return &result, nil
}

func (c *CachedCodingProblemRepository) SetDailyProblem(ctx context.Context, problemId int64) error {
	return c.dao.SetDailyProblem(ctx, problemId)
}

func (c *CachedCodingProblemRepository) GetDailyProblemHistory(ctx context.Context) ([]domain.CodingProblem, error) {
	problems, err := c.dao.GetDailyProblemHistory(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.CodingProblem, 0, len(problems))
	for _, p := range problems {
		result = append(result, c.toDomain(p))
	}
	return result, nil
}

func (c *CachedCodingProblemRepository) MarkAsDailyProblem(ctx context.Context, problemId int64, date time.Time) error {
	return c.dao.MarkAsDailyProblem(ctx, problemId, date)
}

func (c *CachedCodingProblemRepository) UpdateStudyStatus(ctx context.Context, problemId int64, status string, lastStudied *time.Time) error {
	return c.dao.UpdateStudyStatus(ctx, problemId, status, lastStudied)
}

func (c *CachedCodingProblemRepository) SaveDailyProblem(dailyProblem *domain.DailyProblem) error {
	ctx := context.Background()
	return c.dao.SaveDailyProblem(ctx, dailyProblem)
}

func toEntity(p domain.CodingProblem) dao.CodingProblem {
	return dao.CodingProblem{
		Id:             p.Id,
		Title:          p.Title,
		Difficulty:     p.Difficulty,
		Tags:           dao.StringSlice(p.Tags),
		Source:         p.Source,
		SourceId:       p.SourceId,
		SourceUrl:      p.SourceUrl,
		StudyStatus:    p.StudyStatus,
		LastStudied:    p.LastStudied,
		IsDailyProblem: p.IsDailyProblem,
		DailyDate:      p.DailyDate,
		IsHot100:       p.IsHot100,
		Ctime:          p.Ctime,
		Utime:          p.Utime,
	}
}

// 转换方法
func (r *CachedCodingProblemRepository) toDomain(p dao.CodingProblem) domain.CodingProblem {
	return domain.CodingProblem{
		Id:             p.Id,
		Title:          p.Title,
		Difficulty:     p.Difficulty,
		Tags:           []string(p.Tags),
		Source:         p.Source,
		SourceId:       p.SourceId,
		SourceUrl:      p.SourceUrl,
		StudyStatus:    p.StudyStatus,
		LastStudied:    p.LastStudied,
		IsDailyProblem: p.IsDailyProblem,
		DailyDate:      p.DailyDate,
		IsHot100:       p.IsHot100,
		Ctime:          p.Ctime,
		Utime:          p.Utime,
	}
}
