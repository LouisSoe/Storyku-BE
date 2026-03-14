package usecase

import (
	"context"
	"errors"
	"storyku-be/core/domain"
	"storyku-be/core/repository"

	"github.com/google/uuid"
)

type ChapterUsecase interface {
	AddChapter(ctx context.Context, storyID string, chapter *domain.Chapter) error
	UpdateChapter(ctx context.Context, storyID, chapterID string, chapter *domain.Chapter) error
	DeleteChapter(ctx context.Context, storyID, chapterID string) error
}

type chapterUsecase struct {
	storyRepo   repository.StoryRepository
	chapterRepo repository.ChapterRepository
}

func NewChapterUsecase(
	storyRepo repository.StoryRepository,
	chapterRepo repository.ChapterRepository,
) ChapterUsecase {
	return &chapterUsecase{
		storyRepo:   storyRepo,
		chapterRepo: chapterRepo,
	}
}

func (u *chapterUsecase) AddChapter(ctx context.Context, storyID string, chapter *domain.Chapter) error {
	story, err := u.storyRepo.FindByID(ctx, storyID)
	if err != nil || story == nil {
		return errors.New("story not found")
	}
	if chapter.Title == "" {
		return errors.New("chapter title is required")
	}

	count, err := u.chapterRepo.CountByStoryID(ctx, storyID)
	if err != nil {
		return err
	}

	chapter.ChapterID = uuid.New().String()
	chapter.StoryID = storyID
	chapter.OrderIndex = count
	return u.chapterRepo.Create(ctx, chapter)
}

func (u *chapterUsecase) UpdateChapter(ctx context.Context, storyID, chapterID string, chapter *domain.Chapter) error {
	story, err := u.storyRepo.FindByID(ctx, storyID)
	if err != nil || story == nil {
		return errors.New("story not found")
	}

	existing, err := u.chapterRepo.FindByID(ctx, chapterID)
	if err != nil || existing == nil {
		return errors.New("chapter not found")
	}
	if existing.StoryID != storyID {
		return errors.New("chapter does not belong to this story")
	}
	if chapter.Title == "" {
		return errors.New("chapter title is required")
	}

	chapter.ChapterID = chapterID
	chapter.StoryID = storyID
	chapter.OrderIndex = existing.OrderIndex
	return u.chapterRepo.Update(ctx, chapter)
}

func (u *chapterUsecase) DeleteChapter(ctx context.Context, storyID, chapterID string) error {
	story, err := u.storyRepo.FindByID(ctx, storyID)
	if err != nil || story == nil {
		return errors.New("story not found")
	}

	existing, err := u.chapterRepo.FindByID(ctx, chapterID)
	if err != nil || existing == nil {
		return errors.New("chapter not found")
	}
	if existing.StoryID != storyID {
		return errors.New("chapter does not belong to this story")
	}

	return u.chapterRepo.Delete(ctx, chapterID)
}