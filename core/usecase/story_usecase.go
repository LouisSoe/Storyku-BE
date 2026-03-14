package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type StoryUsecase interface {
	GetAll(ctx context.Context, filter repository.StoryFilter) ([]domain.StoryDetail, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.StoryDetail, error)
	Create(ctx context.Context, story *domain.Story, tagIDs []uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, story *domain.Story, tagIDs []uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type storyUsecase struct {
	storyRepo    repository.StoryRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

func NewStoryUsecase(
	storyRepo repository.StoryRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) StoryUsecase {
	return &storyUsecase{
		storyRepo:    storyRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

func (u *storyUsecase) GetAll(ctx context.Context, filter repository.StoryFilter) ([]domain.StoryDetail, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}
	return u.storyRepo.FindAll(ctx, filter)
}

func (u *storyUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.StoryDetail, error) {
	story, err := u.storyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}
	return story, nil
}

func (u *storyUsecase) Create(ctx context.Context, story *domain.Story, tagIDs []uuid.UUID) error {
	if err := u.validateStory(ctx, story, tagIDs); err != nil {
		return err
	}
	if story.Status == "" {
		story.Status = domain.StatusDraft
	}
	if !story.IsValidStatus() {
		return errors.New("status must be 'publish' or 'draft'")
	}
	story.ID = uuid.New()
	return u.storyRepo.Create(ctx, story, tagIDs)
}

func (u *storyUsecase) Update(ctx context.Context, id uuid.UUID, story *domain.Story, tagIDs []uuid.UUID) error {
	existing, err := u.storyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("story not found")
	}
	if err := u.validateStory(ctx, story, tagIDs); err != nil {
		return err
	}
	if !story.IsValidStatus() {
		return errors.New("status must be 'publish' or 'draft'")
	}
	story.ID = id
	return u.storyRepo.Update(ctx, story, tagIDs)
}

func (u *storyUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := u.storyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("story not found")
	}
	return u.storyRepo.Delete(ctx, id)
}

func (u *storyUsecase) validateStory(ctx context.Context, story *domain.Story, tagIDs []uuid.UUID) error {
	if story.Title == "" {
		return errors.New("title is required")
	}
	if story.Author == "" {
		return errors.New("author is required")
	}

	if story.CategoryID != nil {
		cat, err := u.categoryRepo.FindByID(ctx, *story.CategoryID)
		if err != nil {
			return err
		}
		if cat == nil {
			return errors.New("category not found")
		}
	}

	if len(tagIDs) > 0 {
		tags, err := u.tagRepo.FindByIDs(ctx, tagIDs)
		if err != nil {
			return err
		}
		if len(tags) != len(tagIDs) {
			return errors.New("one or more tag IDs are invalid")
		}
	}
	return nil
}