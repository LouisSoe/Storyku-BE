package utils

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaginationParams struct {
	Page  int
	Limit int
}

func GetPagination(c echo.Context) PaginationParams {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{Page: page, Limit: limit}
}

func BuildMeta(page, limit int, total int64) Meta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}