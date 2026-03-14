package usecase_test

import (
	"context"
	"errors"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"storyku-be/core/usecase"
	"testing"
)

// ─── Mock StoryRepository ────────────────────────────────────────────────────

type mockStoryRepo struct {
	stories map[string]*domain.Story
}

func newMockStoryRepo() *mockStoryRepo {
	return &mockStoryRepo{stories: make(map[string]*domain.Story)}
}

func (m *mockStoryRepo) FindAll(_ context.Context, _ repository.StoryFilter) ([]domain.Story, int64, error) {
	var list []domain.Story
	for _, s := range m.stories {
		list = append(list, *s)
	}
	return list, int64(len(list)), nil
}

func (m *mockStoryRepo) FindByID(_ context.Context, id string) (*domain.Story, error) {
	if s, ok := m.stories[id]; ok {
		return s, nil
	}
	return nil, errors.New("not found")
}

func (m *mockStoryRepo) Create(_ context.Context, story *domain.Story) error {
	m.stories[story.StoryID] = story
	return nil
}

func (m *mockStoryRepo) Update(_ context.Context, story *domain.Story) error {
	if _, ok := m.stories[story.StoryID]; !ok {
		return errors.New("not found")
	}
	m.stories[story.StoryID] = story
	return nil
}

func (m *mockStoryRepo) Delete(_ context.Context, id string) error {
	delete(m.stories, id)
	return nil
}

// ─── Mock ChapterRepository ──────────────────────────────────────────────────

type mockChapterRepo struct {
	chapters map[string]*domain.Chapter
}

func newMockChapterRepo() *mockChapterRepo {
	return &mockChapterRepo{chapters: make(map[string]*domain.Chapter)}
}

func (m *mockChapterRepo) FindByStoryID(_ context.Context, storyID string) ([]domain.Chapter, error) {
	var list []domain.Chapter
	for _, c := range m.chapters {
		if c.StoryID == storyID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *mockChapterRepo) FindByID(_ context.Context, id string) (*domain.Chapter, error) {
	if c, ok := m.chapters[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

func (m *mockChapterRepo) Create(_ context.Context, chapter *domain.Chapter) error {
	m.chapters[chapter.ChapterID] = chapter
	return nil
}

func (m *mockChapterRepo) Update(_ context.Context, chapter *domain.Chapter) error {
	m.chapters[chapter.ChapterID] = chapter
	return nil
}

func (m *mockChapterRepo) Delete(_ context.Context, id string) error {
	delete(m.chapters, id)
	return nil
}

func (m *mockChapterRepo) CountByStoryID(_ context.Context, storyID string) (int, error) {
	count := 0
	for _, c := range m.chapters {
		if c.StoryID == storyID {
			count++
		}
	}
	return count, nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newStoryUsecase() (usecase.StoryUsecase, *mockStoryRepo) {
	sr := newMockStoryRepo()
	return usecase.NewStoryUsecase(sr), sr
}

func newChapterUsecase(sr *mockStoryRepo) (usecase.ChapterUsecase, *mockChapterRepo) {
	cr := newMockChapterRepo()
	return usecase.NewChapterUsecase(sr, cr), cr
}

func newValidStory() *domain.Story {
	return &domain.Story{
		Title:    "Test Story",
		Author:   "Author A",
		Synopsis: "Synopsis here",
		Category: string(domain.CategoryTechnology),
		Status:   string(domain.StatusDraft),
		Tags:     []string{"golang", "backend"},
	}
}

// ─── Story Tests ─────────────────────────────────────────────────────────────

func TestCreateStory_Success(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()

	if err := uc.CreateStory(context.Background(), story); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if story.StoryID == "" {
		t.Error("expected StoryID to be set")
	}
}

func TestCreateStory_MissingTitle(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	story.Title = ""

	err := uc.CreateStory(context.Background(), story)
	if err == nil || err.Error() != "title is required" {
		t.Fatalf("expected 'title is required', got: %v", err)
	}
}

func TestCreateStory_MissingAuthor(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	story.Author = ""

	err := uc.CreateStory(context.Background(), story)
	if err == nil || err.Error() != "author is required" {
		t.Fatalf("expected 'author is required', got: %v", err)
	}
}

func TestCreateStory_InvalidCategory(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	story.Category = "Sports"

	if err := uc.CreateStory(context.Background(), story); err == nil {
		t.Fatal("expected error for invalid category")
	}
}

func TestCreateStory_DefaultStatusDraft(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	story.Status = ""

	if err := uc.CreateStory(context.Background(), story); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if story.Status != string(domain.StatusDraft) {
		t.Errorf("expected status 'draft', got: %s", story.Status)
	}
}

func TestGetStoryByID_NotFound(t *testing.T) {
	uc, _ := newStoryUsecase()
	_, err := uc.GetStoryByID(context.Background(), "non-existent-id")
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

func TestUpdateStory_Success(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	_ = uc.CreateStory(context.Background(), story)

	story.Title = "Updated Title"
	story.Status = string(domain.StatusPublish)
	if err := uc.UpdateStory(context.Background(), story.StoryID, story); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateStory_NotFound(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	err := uc.UpdateStory(context.Background(), "bad-id", story)
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

func TestDeleteStory_Success(t *testing.T) {
	uc, _ := newStoryUsecase()
	story := newValidStory()
	_ = uc.CreateStory(context.Background(), story)

	if err := uc.DeleteStory(context.Background(), story.StoryID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteStory_NotFound(t *testing.T) {
	uc, _ := newStoryUsecase()
	err := uc.DeleteStory(context.Background(), "bad-id")
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

// ─── Chapter Tests ────────────────────────────────────────────────────────────

func TestAddChapter_Success(t *testing.T) {
	suc, sr := newStoryUsecase()
	cuc, _ := newChapterUsecase(sr)

	story := newValidStory()
	_ = suc.CreateStory(context.Background(), story)

	chapter := &domain.Chapter{Title: "Chapter 1", Content: "<p>Content</p>"}
	if err := cuc.AddChapter(context.Background(), story.StoryID, chapter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chapter.ChapterID == "" {
		t.Error("expected ChapterID to be set")
	}
}

func TestAddChapter_StoryNotFound(t *testing.T) {
	_, sr := newStoryUsecase()
	cuc, _ := newChapterUsecase(sr)

	chapter := &domain.Chapter{Title: "Chapter 1"}
	err := cuc.AddChapter(context.Background(), "bad-id", chapter)
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

func TestAddChapter_MissingTitle(t *testing.T) {
	suc, sr := newStoryUsecase()
	cuc, _ := newChapterUsecase(sr)

	story := newValidStory()
	_ = suc.CreateStory(context.Background(), story)

	chapter := &domain.Chapter{Title: ""}
	err := cuc.AddChapter(context.Background(), story.StoryID, chapter)
	if err == nil || err.Error() != "chapter title is required" {
		t.Fatalf("expected 'chapter title is required', got: %v", err)
	}
}

func TestDeleteChapter_WrongStory(t *testing.T) {
	suc, sr := newStoryUsecase()
	cuc, _ := newChapterUsecase(sr)

	story1 := newValidStory()
	_ = suc.CreateStory(context.Background(), story1)

	story2 := newValidStory()
	story2.Title = "Story 2"
	_ = suc.CreateStory(context.Background(), story2)

	chapter := &domain.Chapter{Title: "Chapter 1"}
	_ = cuc.AddChapter(context.Background(), story1.StoryID, chapter)

	err := cuc.DeleteChapter(context.Background(), story2.StoryID, chapter.ChapterID)
	if err == nil || err.Error() != "chapter does not belong to this story" {
		t.Fatalf("expected 'chapter does not belong to this story', got: %v", err)
	}
}

func TestAddChapter_OrderIndex(t *testing.T) {
	suc, sr := newStoryUsecase()
	cuc, _ := newChapterUsecase(sr)

	story := newValidStory()
	_ = suc.CreateStory(context.Background(), story)

	ch1 := &domain.Chapter{Title: "Chapter 1"}
	ch2 := &domain.Chapter{Title: "Chapter 2"}
	_ = cuc.AddChapter(context.Background(), story.StoryID, ch1)
	_ = cuc.AddChapter(context.Background(), story.StoryID, ch2)

	if ch1.OrderIndex != 0 {
		t.Errorf("expected ch1 order_index 0, got %d", ch1.OrderIndex)
	}
	if ch2.OrderIndex != 1 {
		t.Errorf("expected ch2 order_index 1, got %d", ch2.OrderIndex)
	}
}