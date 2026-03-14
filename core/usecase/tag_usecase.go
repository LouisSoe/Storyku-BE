package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type TagUsecase interface {
	GetAll(ctx context.Context) ([]domain.Tag, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error)
	Create(ctx context.Context, tag *domain.Tag) error
	Update(ctx context.Context, id uuid.UUID, tag *domain.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type tagUsecase struct {
	repo repository.TagRepository
}

func NewTagUsecase(repo repository.TagRepository) TagUsecase {
	return &tagUsecase{repo: repo}
}

func (u *tagUsecase) GetAll(ctx context.Context) ([]domain.Tag, error) {
	return u.repo.FindAll(ctx)
}

func (u *tagUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	tag, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

func (u *tagUsecase) Create(ctx context.Context, tag *domain.Tag) error {
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	tag.ID = uuid.New()
	tag.Slug = toSlug(tag.Name)

	existing, _ := u.repo.FindBySlug(ctx, tag.Slug)
	if existing != nil {
		return errors.New("tag with this name already exists")
	}
	return u.repo.Create(ctx, tag)
}

func (u *tagUsecase) Update(ctx context.Context, id uuid.UUID, tag *domain.Tag) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("tag not found")
	}
	if tag.Name == "" {
		return errors.New("tag name is required")
	}
	tag.ID = id
	tag.Slug = toSlug(tag.Name)
	return u.repo.Update(ctx, tag)
}

func (u *tagUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("tag not found")
	}
	return u.repo.Delete(ctx, id)
}