package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type CategoryUsecase interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, id uuid.UUID, category *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

func NewCategoryUsecase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (u *categoryUsecase) GetAll(ctx context.Context) ([]domain.Category, error) {
	return u.repo.FindAll(ctx)
}

func (u *categoryUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	cat, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (u *categoryUsecase) Create(ctx context.Context, category *domain.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}
	category.ID = uuid.New()
	category.Slug = toSlug(category.Name)

	existing, _ := u.repo.FindBySlug(ctx, category.Slug)
	if existing != nil {
		return errors.New("category with this name already exists")
	}
	return u.repo.Create(ctx, category)
}

func (u *categoryUsecase) Update(ctx context.Context, id uuid.UUID, category *domain.Category) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}
	if category.Name == "" {
		return errors.New("category name is required")
	}
	category.ID = id
	category.Slug = toSlug(category.Name)
	return u.repo.Update(ctx, category)
}

func (u *categoryUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("category not found")
	}
	return u.repo.Delete(ctx, id)
}

func toSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}