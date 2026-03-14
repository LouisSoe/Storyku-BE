package repository

import (
	"context"

	"github.com/google/uuid"
	"storyku-be/core/domain"
)

type TagRepository interface {
	FindAll(ctx context.Context) ([]domain.Tag, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Tag, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Tag, error)
	Create(ctx context.Context, tag *domain.Tag) error
	Update(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}