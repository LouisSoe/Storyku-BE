package http

import (
	"io"
	"os"
	"path/filepath"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"storyku-be/core/usecase"
	"storyku-be/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type StoryHandler struct {
	uc     usecase.StoryUsecase
	logger *logrus.Logger
}

func NewStoryHandler(uc usecase.StoryUsecase, logger *logrus.Logger) *StoryHandler {
	return &StoryHandler{uc: uc, logger: logger}
}

type storyRequest struct {
	Title    string   `json:"title"    form:"title"`
	Author   string   `json:"author"   form:"author"`
	Synopsis string   `json:"synopsis" form:"synopsis"`
	Category string   `json:"category" form:"category"`
	Tags     []string `json:"tags"     form:"tags"`
	Status   string   `json:"status"   form:"status"`
}

// List godoc
// GET /api/v1/stories
func (h *StoryHandler) List(c echo.Context) error {
	pg := utils.GetPagination(c)

	filter := repository.StoryFilter{
		Search:   strings.TrimSpace(c.QueryParam("search")),
		Category: c.QueryParam("category"),
		Status:   c.QueryParam("status"),
		Page:     pg.Page,
		Limit:    pg.Limit,
	}

	stories, total, err := h.uc.GetStories(c.Request().Context(), filter)
	if err != nil {
		h.logger.WithError(err).Error("List stories failed")
		return utils.InternalError(c, "failed to fetch stories")
	}

	meta := utils.BuildMeta(pg.Page, pg.Limit, total)
	return utils.OKWithMeta(c, "stories retrieved successfully", stories, meta)
}

// GetByID godoc
// GET /api/v1/stories/:id
func (h *StoryHandler) GetByID(c echo.Context) error {
	story, err := h.uc.GetStoryByID(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.NotFound(c, err.Error())
	}
	return utils.OK(c, "story retrieved successfully", story)
}

// Create godoc
// POST /api/v1/stories
func (h *StoryHandler) Create(c echo.Context) error {
	req, err := bindStoryRequest(c)
	if err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	coverURL, err := handleCoverUpload(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	story := &domain.Story{
		Title:    strings.TrimSpace(req.Title),
		Author:   strings.TrimSpace(req.Author),
		Synopsis: strings.TrimSpace(req.Synopsis),
		Category: req.Category,
		CoverURL: coverURL,
		Tags:     sanitizeTags(req.Tags),
		Status:   req.Status,
	}

	if err := h.uc.CreateStory(c.Request().Context(), story); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "story created successfully", story)
}

// Update godoc
// PUT /api/v1/stories/:id
func (h *StoryHandler) Update(c echo.Context) error {
	req, err := bindStoryRequest(c)
	if err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	coverURL, err := handleCoverUpload(c)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	story := &domain.Story{
		Title:    strings.TrimSpace(req.Title),
		Author:   strings.TrimSpace(req.Author),
		Synopsis: strings.TrimSpace(req.Synopsis),
		Category: req.Category,
		CoverURL: coverURL,
		Tags:     sanitizeTags(req.Tags),
		Status:   req.Status,
	}

	if err := h.uc.UpdateStory(c.Request().Context(), c.Param("id"), story); err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "story updated successfully", story)
}

// Delete godoc
// DELETE /api/v1/stories/:id
func (h *StoryHandler) Delete(c echo.Context) error {
	if err := h.uc.DeleteStory(c.Request().Context(), c.Param("id")); err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("DeleteStory failed")
		return utils.InternalError(c, "failed to delete story")
	}
	return utils.OK(c, "story deleted successfully", nil)
}

// ─── Shared helpers (dipakai chapter_handler juga via package) ───────────────

func bindStoryRequest(c echo.Context) (*storyRequest, error) {
	var req storyRequest
	if err := c.Bind(&req); err != nil {
		return nil, err
	}
	if strings.Contains(c.Request().Header.Get("Content-Type"), "multipart/form-data") {
		if raw := c.FormValue("tags"); raw != "" {
			req.Tags = strings.Split(raw, ",")
		}
	}
	return &req, nil
}

func handleCoverUpload(c echo.Context) (string, error) {
	file, err := c.FormFile("cover")
	if err != nil {
		return "", nil
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		return "", echo.NewHTTPError(400, "only jpg, jpeg, png, webp allowed")
	}
	if file.Size > 5*1024*1024 {
		return "", echo.NewHTTPError(400, "cover image must not exceed 5MB")
	}

	dir := "uploads/covers"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	filename := time.Now().Format("20060102150405") + "_" + uuid.New().String()[:8] + ext
	dst := filepath.Join(dir, filename)

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

func sanitizeTags(tags []string) []string {
	var result []string
	for _, t := range tags {
		if t = strings.TrimSpace(t); t != "" {
			result = append(result, t)
		}
	}
	return result
}