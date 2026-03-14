package http

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"storyku-be/core/usecase"
	"storyku-be/pkg/utils"
)

type StoryHandler struct {
	storyUsecase   usecase.StoryUsecase
	chapterUsecase usecase.ChapterUsecase
	logger         *logrus.Logger
}

func NewStoryHandler(su usecase.StoryUsecase, cu usecase.ChapterUsecase, logger *logrus.Logger) *StoryHandler {
	return &StoryHandler{storyUsecase: su, chapterUsecase: cu, logger: logger}
}

type storyRequest struct {
	CategoryID string   `json:"category_id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Synopsis   string   `json:"synopsis"`
	TagIDs     []string `json:"tag_ids"`
	Status     string   `json:"status"`
}

type storyDetailResponse struct {
	domain.StoryDetail
	Chapters []domain.Chapter `json:"chapters"`
}

// List godoc
// GET /api/v1/stories
func (h *StoryHandler) List(c echo.Context) error {
	pagination := utils.GetPagination(c)
	filter := repository.StoryFilter{
		Search:     strings.TrimSpace(c.QueryParam("search")),
		CategoryID: c.QueryParam("category_id"),
		Status:     c.QueryParam("status"),
		Page:       pagination.Page,
		Limit:      pagination.Limit,
	}

	stories, total, err := h.storyUsecase.GetAll(c.Request().Context(), filter)
	if err != nil {
		h.logger.WithError(err).Error("failed to list stories")
		return utils.InternalError(c, "failed to fetch stories")
	}

	meta := utils.BuildMeta(pagination.Page, pagination.Limit, total)
	return utils.OKWithMeta(c, "stories retrieved successfully", stories, meta)
}

// GetByID godoc
// GET /api/v1/stories/:id
func (h *StoryHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid story ID")
	}

	story, err := h.storyUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("failed to get story")
		return utils.InternalError(c, "failed to fetch story")
	}

	chapters, err := h.chapterUsecase.GetByStoryID(c.Request().Context(), id)
	if err != nil {
		chapters = []domain.Chapter{}
	}

	return utils.OK(c, "story retrieved successfully", storyDetailResponse{
		StoryDetail: *story,
		Chapters:    chapters,
	})
}

// Create godoc
// POST /api/v1/stories
func (h *StoryHandler) Create(c echo.Context) error {
	req, tagIDs, err := h.bindStoryRequest(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	coverURL, err := h.handleFileUpload(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	story := h.toStoreDomain(req, coverURL)
	if err := h.storyUsecase.Create(c.Request().Context(), story, tagIDs); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "story created successfully", story)
}

// Update godoc
// PUT /api/v1/stories/:id
func (h *StoryHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid story ID")
	}

	req, tagIDs, err := h.bindStoryRequest(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	existing, err := h.storyUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "failed to fetch story")
	}

	coverURL, err := h.handleFileUpload(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	if coverURL == "" {
		coverURL = existing.CoverURL
	}

	story := h.toStoreDomain(req, coverURL)
	if err := h.storyUsecase.Update(c.Request().Context(), id, story, tagIDs); err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}

	story.ID = id
	return utils.OK(c, "story updated successfully", story)
}

// Delete godoc
// DELETE /api/v1/stories/:id
func (h *StoryHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid story ID")
	}
	if err := h.storyUsecase.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("failed to delete story")
		return utils.InternalError(c, "failed to delete story")
	}
	return utils.OK(c, "story deleted successfully", nil)
}

func (h *StoryHandler) bindStoryRequest(c echo.Context) (*storyRequest, []uuid.UUID, error) {
	var req storyRequest
	if err := c.Bind(&req); err != nil {
		return nil, nil, errors.New("invalid request body")
	}

	var tagIDs []uuid.UUID
	for _, raw := range req.TagIDs {
		id, err := uuid.Parse(strings.TrimSpace(raw))
		if err != nil {
			return nil, nil, errors.New("invalid tag_id: " + raw)
		}
		tagIDs = append(tagIDs, id)
	}
	return &req, tagIDs, nil
}

func (h *StoryHandler) toStoreDomain(req *storyRequest, coverURL string) *domain.Story {
	story := &domain.Story{
		Title:    strings.TrimSpace(req.Title),
		Author:   strings.TrimSpace(req.Author),
		Synopsis: strings.TrimSpace(req.Synopsis),
		CoverURL: coverURL,
		Status:   domain.StoryStatus(req.Status),
	}
	if req.CategoryID != "" {
		if id, err := uuid.Parse(req.CategoryID); err == nil {
			story.CategoryID = &id
		}
	}
	return story
}

func (h *StoryHandler) handleFileUpload(c echo.Context) (string, error) {
	file, err := c.FormFile("cover")
	if err != nil {
		return "", nil
	}

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", errors.New("only jpg, jpeg, png, webp are allowed")
	}
	if file.Size > 5*1024*1024 {
		return "", errors.New("file size must not exceed 5MB")
	}

	uploadDir := "uploads/covers"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	filename := time.Now().Format("20060102150405") + "_" + uuid.New().String()[:8] + ext
	dst := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}
	return "/" + dst, nil
}