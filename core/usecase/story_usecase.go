package usecase

import (
	"context"
	"errors"
	"storyku-be/core/domain"
	"storyku-be/core/repository"

	"github.com/google/uuid"
)

type StoryUsecase interface {
	GetStories(ctx context.Context, filter repository.StoryFilter) ([]domain.Story, int64, error)
	GetStoryByID(ctx context.Context, id string) (*domain.Story, error)
	CreateStory(ctx context.Context, story *domain.Story) error
	UpdateStory(ctx context.Context, id string, story *domain.Story) error
	DeleteStory(ctx context.Context, id string) error
}

type storyUsecase struct {
	repo repository.StoryRepository
}

func NewStoryUsecase(repo repository.StoryRepository) StoryUsecase {
	return &storyUsecase{repo: repo}
}

func (u *storyUsecase) GetStories(ctx context.Context, filter repository.StoryFilter) ([]domain.Story, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}
	return u.repo.FindAll(ctx, filter)
}

func (u *storyUsecase) GetStoryByID(ctx context.Context, id string) (*domain.Story, error) {
	if id == "" {
		return nil, errors.New("story_id is required")
	}
	story, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("story not found")
	}
	return story, nil
}

func (u *storyUsecase) CreateStory(ctx context.Context, story *domain.Story) error {
	if err := validateStory(story); err != nil {
		return err
	}
	if story.Status == "" {
		story.Status = string(domain.StatusDraft)
	}
	if !story.IsValidStatus() {
		return errors.New("status must be 'publish' or 'draft'")
	}
	story.StoryID = uuid.New().String()
	return u.repo.Create(ctx, story)
}

func (u *storyUsecase) UpdateStory(ctx context.Context, id string, story *domain.Story) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil || existing == nil {
		return errors.New("story not found")
	}
	if err := validateStory(story); err != nil {
		return err
	}
	if !story.IsValidStatus() {
		return errors.New("status must be 'publish' or 'draft'")
	}
	if story.CoverURL == "" {
		story.CoverURL = existing.CoverURL
	}
	story.StoryID = id
	return u.repo.Update(ctx, story)
}

func (u *storyUsecase) DeleteStory(ctx context.Context, id string) error {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil || existing == nil {
		return errors.New("story not found")
	}
	return u.repo.Delete(ctx, id)
}

func validateStory(story *domain.Story) error {
	if story.Title == "" {
		return errors.New("title is required")
	}
	if story.Author == "" {
		return errors.New("author is required")
	}
	if !story.IsValidCategory() {
		return errors.New("category must be 'Financial', 'Technology', or 'Health'")
	}
	return nil
}