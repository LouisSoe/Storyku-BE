package http

import (
	"storyku-be/core/usecase"
	"storyku-be/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type DashboardHandler struct {
	usecase usecase.DashboardUsecase
	logger  *logrus.Logger
}

func NewDashboardHandler(u usecase.DashboardUsecase, logger *logrus.Logger) *DashboardHandler {
	return &DashboardHandler{usecase: u, logger: logger}
}

// Get godoc
// GET /api/v1/dashboard
func (h *DashboardHandler) Get(c echo.Context) error {
	data, err := h.usecase.GetDashboard(c.Request().Context())
	if err != nil {
		h.logger.WithError(err).Error("failed to get dashboard data")
		return utils.InternalError(c, "failed to fetch dashboard data")
	}
	return utils.OK(c, "dashboard data retrieved successfully", data)
}