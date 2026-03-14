package http

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"storyku-be/core/domain"
	"storyku-be/core/usecase"
	"storyku-be/pkg/utils"
)

type TagHandler struct {
	usecase usecase.TagUsecase
	logger  *logrus.Logger
}

func NewTagHandler(u usecase.TagUsecase, logger *logrus.Logger) *TagHandler {
	return &TagHandler{usecase: u, logger: logger}
}

type tagRequest struct {
	Name string `json:"name"`
}

// List godoc
// GET /api/v1/tags
func (h *TagHandler) List(c echo.Context) error {
	tags, err := h.usecase.GetAll(c.Request().Context())
	if err != nil {
		h.logger.WithError(err).Error("failed to list tags")
		return utils.InternalError(c, "failed to fetch tags")
	}
	return utils.OK(c, "tags retrieved successfully", tags)
}

// GetByID godoc
// GET /api/v1/tags/:id
func (h *TagHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid tag ID")
	}
	tag, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "tag not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "failed to fetch tag")
	}
	return utils.OK(c, "tag retrieved successfully", tag)
}

// Create godoc
// POST /api/v1/tags
func (h *TagHandler) Create(c echo.Context) error {
	var req tagRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	tag := &domain.Tag{Name: strings.TrimSpace(req.Name)}
	if err := h.usecase.Create(c.Request().Context(), tag); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "tag created successfully", tag)
}

// Update godoc
// PUT /api/v1/tags/:id
func (h *TagHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid tag ID")
	}
	var req tagRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	tag := &domain.Tag{Name: strings.TrimSpace(req.Name)}
	if err := h.usecase.Update(c.Request().Context(), id, tag); err != nil {
		if err.Error() == "tag not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}
	tag.ID = id
	return utils.OK(c, "tag updated successfully", tag)
}

// Delete godoc
// DELETE /api/v1/tags/:id
func (h *TagHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid tag ID")
	}
	if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "tag not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("failed to delete tag")
		return utils.InternalError(c, "failed to delete tag")
	}
	return utils.OK(c, "tag deleted successfully", nil)
}