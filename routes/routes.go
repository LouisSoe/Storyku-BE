package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	httpHandler "storyku-be/interfaces/http"
)

func Register(
	e *echo.Echo,
	storyHandler *httpHandler.StoryHandler,
	chapterHandler *httpHandler.ChapterHandler,
	categoryHandler *httpHandler.CategoryHandler,
	tagHandler *httpHandler.TagHandler,
) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/uploads", "uploads")

	api := e.Group("/api/v1")

	// Master: Categories
	categories := api.Group("/categories")
	categories.GET("", categoryHandler.List)
	categories.POST("", categoryHandler.Create)
	categories.GET("/:id", categoryHandler.GetByID)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	// Master: Tags
	tags := api.Group("/tags")
	tags.GET("", tagHandler.List)
	tags.POST("", tagHandler.Create)
	tags.GET("/:id", tagHandler.GetByID)
	tags.PUT("/:id", tagHandler.Update)
	tags.DELETE("/:id", tagHandler.Delete)

	// Stories
	stories := api.Group("/stories")
	stories.GET("", storyHandler.List)
	stories.POST("", storyHandler.Create)
	stories.GET("/:id", storyHandler.GetByID)
	stories.PUT("/:id", storyHandler.Update)
	stories.DELETE("/:id", storyHandler.Delete)

	// Chapters
	chapters := stories.Group("/:id/chapters")
	chapters.POST("", chapterHandler.Create)
	chapters.PUT("/:cid", chapterHandler.Update)
	chapters.DELETE("/:cid", chapterHandler.Delete)
}