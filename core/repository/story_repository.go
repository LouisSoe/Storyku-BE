package repository

import (
	"context"
	"storyku-be/core/domain"
)

type StoryFilter struct {
	Search   string
	Category string
	Status   string
	Page     int
	Limit    int
}

type StoryRepository interface {
	FindAll(ctx context.Context, filter StoryFilter) ([]domain.Story, int64, error)
	FindByID(ctx context.Context, id string) (*domain.Story, error)
	Create(ctx context.Context, story *domain.Story) error
	Update(ctx context.Context, story *domain.Story) error
	Delete(ctx context.Context, id string) error
}