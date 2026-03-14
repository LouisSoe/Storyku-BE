package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"storyku-be/core/usecase"
)

// ─── Mock Category Repository ─────────────────────────────────────────────────

type mockCategoryRepo struct {
	categories map[uuid.UUID]*domain.Category
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[uuid.UUID]*domain.Category)}
}

func (m *mockCategoryRepo) FindAll(_ context.Context) ([]domain.Category, error) {
	var list []domain.Category
	for _, c := range m.categories {
		list = append(list, *c)
	}
	return list, nil
}

func (m *mockCategoryRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.Category, error) {
	if c, ok := m.categories[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *mockCategoryRepo) FindBySlug(_ context.Context, slug string) (*domain.Category, error) {
	for _, c := range m.categories {
		if c.Slug == slug {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryRepo) Create(_ context.Context, c *domain.Category) error {
	m.categories[c.ID] = c
	return nil
}

func (m *mockCategoryRepo) Update(_ context.Context, c *domain.Category) error {
	m.categories[c.ID] = c
	return nil
}

func (m *mockCategoryRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.categories, id)
	return nil
}

// ─── Mock Tag Repository ──────────────────────────────────────────────────────

type mockTagRepo struct {
	tags map[uuid.UUID]*domain.Tag
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{tags: make(map[uuid.UUID]*domain.Tag)}
}

func (m *mockTagRepo) FindAll(_ context.Context) ([]domain.Tag, error) {
	var list []domain.Tag
	for _, t := range m.tags {
		list = append(list, *t)
	}
	return list, nil
}

func (m *mockTagRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.Tag, error) {
	if t, ok := m.tags[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *mockTagRepo) FindByIDs(_ context.Context, ids []uuid.UUID) ([]domain.Tag, error) {
	var list []domain.Tag
	for _, id := range ids {
		if t, ok := m.tags[id]; ok {
			list = append(list, *t)
		}
	}
	return list, nil
}

func (m *mockTagRepo) FindBySlug(_ context.Context, slug string) (*domain.Tag, error) {
	for _, t := range m.tags {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTagRepo) Create(_ context.Context, t *domain.Tag) error {
	m.tags[t.ID] = t
	return nil
}

func (m *mockTagRepo) Update(_ context.Context, t *domain.Tag) error {
	m.tags[t.ID] = t
	return nil
}

func (m *mockTagRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.tags, id)
	return nil
}

// ─── Mock Story Repository ────────────────────────────────────────────────────

type mockStoryRepo struct {
	stories   map[uuid.UUID]*domain.Story
	storyTags map[uuid.UUID][]uuid.UUID
}

func newMockStoryRepo() *mockStoryRepo {
	return &mockStoryRepo{
		stories:   make(map[uuid.UUID]*domain.Story),
		storyTags: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (m *mockStoryRepo) FindAll(_ context.Context, _ repository.StoryFilter) ([]domain.StoryDetail, int64, error) {
	var list []domain.StoryDetail
	for _, s := range m.stories {
		list = append(list, domain.StoryDetail{Story: *s})
	}
	return list, int64(len(list)), nil
}

func (m *mockStoryRepo) FindByID(_ context.Context, id uuid.UUID) (*domain.StoryDetail, error) {
	if s, ok := m.stories[id]; ok {
		return &domain.StoryDetail{Story: *s}, nil
	}
	return nil, nil
}

func (m *mockStoryRepo) Create(_ context.Context, story *domain.Story, tagIDs []uuid.UUID) error {
	m.stories[story.ID] = story
	m.storyTags[story.ID] = tagIDs
	return nil
}

func (m *mockStoryRepo) Update(_ context.Context, story *domain.Story, tagIDs []uuid.UUID) error {
	m.stories[story.ID] = story
	m.storyTags[story.ID] = tagIDs
	return nil
}

func (m *mockStoryRepo) Delete(_ context.Context, id uuid.UUID) error {
	delete(m.stories, id)
	delete(m.storyTags, id)
	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func buildUsecase() (usecase.StoryUsecase, *mockStoryRepo, *mockCategoryRepo, *mockTagRepo) {
	sr := newMockStoryRepo()
	cr := newMockCategoryRepo()
	tr := newMockTagRepo()
	uc := usecase.NewStoryUsecase(sr, cr, tr)
	return uc, sr, cr, tr
}

func seedCategory(cr *mockCategoryRepo) *domain.Category {
	cat := &domain.Category{
		ID:   uuid.New(),
		Name: "Technology",
		Slug: "technology",
	}
	cr.categories[cat.ID] = cat
	return cat
}

func seedTag(tr *mockTagRepo, name string) *domain.Tag {
	tag := &domain.Tag{
		ID:   uuid.New(),
		Name: name,
		Slug: name,
	}
	tr.tags[tag.ID] = tag
	return tag
}

func validStory(catID *uuid.UUID) *domain.Story {
	return &domain.Story{
		Title:      "Belajar Golang Clean Architecture",
		Author:     "Vito Hidayat",
		Synopsis:   "Panduan membangun REST API dengan Go",
		CategoryID: catID,
		Status:     domain.StatusDraft,
	}
}

// ─── Tests: Create ────────────────────────────────────────────────────────────

func TestStory_Create_Success(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)

	if err := uc.Create(context.Background(), story, nil); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if story.ID == uuid.Nil {
		t.Error("expected story.ID to be set after create")
	}
}

func TestStory_Create_WithTags_Success(t *testing.T) {
	uc, _, cr, tr := buildUsecase()
	cat := seedCategory(cr)
	tag1 := seedTag(tr, "golang")
	tag2 := seedTag(tr, "backend")

	story := validStory(&cat.ID)
	tagIDs := []uuid.UUID{tag1.ID, tag2.ID}

	if err := uc.Create(context.Background(), story, tagIDs); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestStory_Create_MissingTitle(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)
	story.Title = ""

	err := uc.Create(context.Background(), story, nil)
	if err == nil || err.Error() != "title is required" {
		t.Fatalf("expected 'title is required', got: %v", err)
	}
}

func TestStory_Create_MissingAuthor(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)
	story.Author = ""

	err := uc.Create(context.Background(), story, nil)
	if err == nil || err.Error() != "author is required" {
		t.Fatalf("expected 'author is required', got: %v", err)
	}
}

func TestStory_Create_CategoryNotFound(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	invalidCatID := uuid.New()
	story := validStory(&invalidCatID)

	err := uc.Create(context.Background(), story, nil)
	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected 'category not found', got: %v", err)
	}
}

func TestStory_Create_NilCategory_Success(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)

	if err := uc.Create(context.Background(), story, nil); err != nil {
		t.Fatalf("expected no error with nil category, got: %v", err)
	}
}

func TestStory_Create_InvalidTagID(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)

	fakeTagID := uuid.New()
	err := uc.Create(context.Background(), story, []uuid.UUID{fakeTagID})
	if err == nil || err.Error() != "one or more tag IDs are invalid" {
		t.Fatalf("expected 'one or more tag IDs are invalid', got: %v", err)
	}
}

func TestStory_Create_DefaultDraftStatus(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	story.Status = ""

	if err := uc.Create(context.Background(), story, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if story.Status != domain.StatusDraft {
		t.Errorf("expected status 'draft', got: %s", story.Status)
	}
}

func TestStory_Create_InvalidStatus(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	story.Status = "archived"

	err := uc.Create(context.Background(), story, nil)
	if err == nil || err.Error() != "status must be 'publish' or 'draft'" {
		t.Fatalf("expected status error, got: %v", err)
	}
}

// ─── Tests: GetByID ───────────────────────────────────────────────────────────

func TestStory_GetByID_Success(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	_ = uc.Create(context.Background(), story, nil)

	result, err := uc.GetByID(context.Background(), story.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != story.ID {
		t.Errorf("expected ID %s, got %s", story.ID, result.ID)
	}
}

func TestStory_GetByID_NotFound(t *testing.T) {
	uc, _, _, _ := buildUsecase()

	_, err := uc.GetByID(context.Background(), uuid.New())
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

// ─── Tests: Update ────────────────────────────────────────────────────────────

func TestStory_Update_Success(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)
	_ = uc.Create(context.Background(), story, nil)

	story.Title = "Updated Title"
	story.Status = domain.StatusPublish

	if err := uc.Update(context.Background(), story.ID, story, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStory_Update_WithNewTags(t *testing.T) {
	uc, _, cr, tr := buildUsecase()
	cat := seedCategory(cr)
	tag := seedTag(tr, "postgresql")
	story := validStory(&cat.ID)
	_ = uc.Create(context.Background(), story, nil)

	story.Status = domain.StatusPublish
	if err := uc.Update(context.Background(), story.ID, story, []uuid.UUID{tag.ID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStory_Update_NotFound(t *testing.T) {
	uc, _, cr, _ := buildUsecase()
	cat := seedCategory(cr)
	story := validStory(&cat.ID)

	err := uc.Update(context.Background(), uuid.New(), story, nil)
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

func TestStory_Update_ChangeCategoryToInvalid(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	_ = uc.Create(context.Background(), story, nil)

	invalidCatID := uuid.New()
	story.CategoryID = &invalidCatID

	err := uc.Update(context.Background(), story.ID, story, nil)
	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected 'category not found', got: %v", err)
	}
}

func TestStory_Update_InvalidTagID(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	_ = uc.Create(context.Background(), story, nil)

	err := uc.Update(context.Background(), story.ID, story, []uuid.UUID{uuid.New()})
	if err == nil || err.Error() != "one or more tag IDs are invalid" {
		t.Fatalf("expected 'one or more tag IDs are invalid', got: %v", err)
	}
}

// ─── Tests: Delete ────────────────────────────────────────────────────────────

func TestStory_Delete_Success(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	story := validStory(nil)
	_ = uc.Create(context.Background(), story, nil)

	if err := uc.Delete(context.Background(), story.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStory_Delete_NotFound(t *testing.T) {
	uc, _, _, _ := buildUsecase()

	err := uc.Delete(context.Background(), uuid.New())
	if err == nil || err.Error() != "story not found" {
		t.Fatalf("expected 'story not found', got: %v", err)
	}
}

// ─── Tests: GetAll ────────────────────────────────────────────────────────────

func TestStory_GetAll_DefaultPagination(t *testing.T) {
	uc, _, _, _ := buildUsecase()
	for i := 0; i < 3; i++ {
		s := validStory(nil)
		_ = uc.Create(context.Background(), s, nil)
	}

	_, total, err := uc.GetAll(context.Background(), repository.StoryFilter{Page: 0, Limit: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
}

func TestStory_GetAll_LimitCapped(t *testing.T) {
	uc, _, _, _ := buildUsecase()

	_, _, err := uc.GetAll(context.Background(), repository.StoryFilter{Page: 1, Limit: 999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}