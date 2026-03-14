package repository

import (
	"context"

	"github.com/google/uuid"
	"storyku-be/core/domain"
)

type StoryFilter struct {
	Search     string
	CategoryID string
	Status     string
	Page       int
	Limit      int
}

type StoryRepository interface {
	FindAll(ctx context.Context, filter StoryFilter) ([]domain.StoryDetail, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.StoryDetail, error)
	Create(ctx context.Context, story *domain.Story, tagIDs []uuid.UUID) error
	Update(ctx context.Context, story *domain.Story, tagIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}