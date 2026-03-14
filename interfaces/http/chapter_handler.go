package http

import (
	"storyku-be/core/domain"
	"storyku-be/core/usecase"
	"storyku-be/pkg/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ChapterHandler struct {
	uc     usecase.ChapterUsecase
	logger *logrus.Logger
}

func NewChapterHandler(uc usecase.ChapterUsecase, logger *logrus.Logger) *ChapterHandler {
	return &ChapterHandler{uc: uc, logger: logger}
}

type chapterRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Add godoc
// POST /api/v1/stories/:id/chapters
func (h *ChapterHandler) Create(c echo.Context) error {
	var req chapterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	chapter := &domain.Chapter{
		Title:   strings.TrimSpace(req.Title),
		Content: req.Content,
	}

	if err := h.uc.AddChapter(c.Request().Context(), c.Param("id"), chapter); err != nil {
		if err.Error() == "story not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}

	return utils.Created(c, "chapter added successfully", chapter)
}

// Update godoc
// PUT /api/v1/stories/:id/chapters/:cid
func (h *ChapterHandler) Update(c echo.Context) error {
	var req chapterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	chapter := &domain.Chapter{
		Title:   strings.TrimSpace(req.Title),
		Content: req.Content,
	}

	if err := h.uc.UpdateChapter(c.Request().Context(), c.Param("id"), c.Param("cid"), chapter); err != nil {
		if err.Error() == "story not found" || err.Error() == "chapter not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}

	return utils.OK(c, "chapter updated successfully", chapter)
}

// Delete godoc
// DELETE /api/v1/stories/:id/chapters/:cid
func (h *ChapterHandler) Delete(c echo.Context) error {
	if err := h.uc.DeleteChapter(c.Request().Context(), c.Param("id"), c.Param("cid")); err != nil {
		if err.Error() == "story not found" || err.Error() == "chapter not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("DeleteChapter failed")
		return utils.InternalError(c, "failed to delete chapter")
	}
	return utils.OK(c, "chapter deleted successfully", nil)
}