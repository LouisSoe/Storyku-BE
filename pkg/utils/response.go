package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func OK(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data})
}

func OKWithMeta(c echo.Context, message string, data interface{}, meta Meta) error {
	return c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data, Meta: &meta})
}

func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{Success: true, Message: message, Data: data})
}

func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, Response{Success: false, Message: message})
}

func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, Response{Success: false, Message: message})
}

func InternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, Response{Success: false, Message: message})
}