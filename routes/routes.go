package routes

import (
	"net/http"
	httpHandler "storyku-be/interfaces/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Register(e *echo.Echo, story *httpHandler.StoryHandler, chapter *httpHandler.ChapterHandler) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/uploads", "uploads")

	api := e.Group("/api/v1")

	s := api.Group("/stories")
	s.GET("", story.List)
	s.POST("", story.Create)
	s.GET("/:id", story.GetByID)
	s.PUT("/:id", story.Update)
	s.DELETE("/:id", story.Delete)

	ch := s.Group("/:id/chapters")
	ch.POST("", chapter.Add)
	ch.PUT("/:cid", chapter.Update)
	ch.DELETE("/:cid", chapter.Delete)
}