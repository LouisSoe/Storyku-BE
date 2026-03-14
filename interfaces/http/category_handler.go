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

type CategoryHandler struct {
	usecase usecase.CategoryUsecase
	logger  *logrus.Logger
}

func NewCategoryHandler(u usecase.CategoryUsecase, logger *logrus.Logger) *CategoryHandler {
	return &CategoryHandler{usecase: u, logger: logger}
}

type categoryRequest struct {
	Name string `json:"name"`
}

// List godoc
// GET /api/v1/categories
func (h *CategoryHandler) List(c echo.Context) error {
	categories, err := h.usecase.GetAll(c.Request().Context())
	if err != nil {
		h.logger.WithError(err).Error("failed to list categories")
		return utils.InternalError(c, "failed to fetch categories")
	}
	return utils.OK(c, "categories retrieved successfully", categories)
}

// GetByID godoc
// GET /api/v1/categories/:id
func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid category ID")
	}
	cat, err := h.usecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "category not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, "failed to fetch category")
	}
	return utils.OK(c, "category retrieved successfully", cat)
}

// Create godoc
// POST /api/v1/categories
func (h *CategoryHandler) Create(c echo.Context) error {
	var req categoryRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	cat := &domain.Category{Name: strings.TrimSpace(req.Name)}
	if err := h.usecase.Create(c.Request().Context(), cat); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "category created successfully", cat)
}

// Update godoc
// PUT /api/v1/categories/:id
func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid category ID")
	}
	var req categoryRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	cat := &domain.Category{Name: strings.TrimSpace(req.Name)}
	if err := h.usecase.Update(c.Request().Context(), id, cat); err != nil {
		if err.Error() == "category not found" {
			return utils.NotFound(c, err.Error())
		}
		return utils.BadRequest(c, err.Error())
	}
	cat.ID = id
	return utils.OK(c, "category updated successfully", cat)
}

// Delete godoc
// DELETE /api/v1/categories/:id
func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.BadRequest(c, "invalid category ID")
	}
	if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
		if err.Error() == "category not found" {
			return utils.NotFound(c, err.Error())
		}
		h.logger.WithError(err).Error("failed to delete category")
		return utils.InternalError(c, "failed to delete category")
	}
	return utils.OK(c, "category deleted successfully", nil)
}