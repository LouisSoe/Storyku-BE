package repository

import (
	"context"
	"storyku-be/core/domain"
)

type ChapterRepository interface {
	FindByStoryID(ctx context.Context, storyID string) ([]domain.Chapter, error)
	FindByID(ctx context.Context, id string) (*domain.Chapter, error)
	Create(ctx context.Context, chapter *domain.Chapter) error
	Update(ctx context.Context, chapter *domain.Chapter) error
	Delete(ctx context.Context, id string) error
	CountByStoryID(ctx context.Context, storyID string) (int, error)
}