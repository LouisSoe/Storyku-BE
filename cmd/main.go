package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"storyku-be/config"
	dbImpl "storyku-be/interfaces/database"
	httpHandler "storyku-be/interfaces/http"
	"storyku-be/core/usecase"
	"storyku-be/routes"
)

func main() {
	cfg := config.Load()

	// Repositories
	storyRepo     := dbImpl.NewStoryRepository(cfg.DB)
	chapterRepo   := dbImpl.NewChapterRepository(cfg.DB)
	categoryRepo  := dbImpl.NewCategoryRepository(cfg.DB)
	tagRepo       := dbImpl.NewTagRepository(cfg.DB)
	dashboardRepo := dbImpl.NewDashboardRepository(cfg.DB)

	// Use Cases
	categoryUsecase  := usecase.NewCategoryUsecase(categoryRepo)
	tagUsecase       := usecase.NewTagUsecase(tagRepo)
	storyUsecase     := usecase.NewStoryUsecase(storyRepo, categoryRepo, tagRepo)
	chapterUsecase   := usecase.NewChapterUsecase(storyRepo, chapterRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo)

	// Handlers
	storyHandler     := httpHandler.NewStoryHandler(storyUsecase, chapterUsecase, cfg.Logger)
	chapterHandler   := httpHandler.NewChapterHandler(chapterUsecase, cfg.Logger)
	categoryHandler  := httpHandler.NewCategoryHandler(categoryUsecase, cfg.Logger)
	tagHandler       := httpHandler.NewTagHandler(tagUsecase, cfg.Logger)
	dashboardHandler := httpHandler.NewDashboardHandler(dashboardUsecase, cfg.Logger)

	e := echo.New()
	e.HideBanner = true

	routes.Register(e, cfg.Logger, storyHandler, chapterHandler, categoryHandler, tagHandler, dashboardHandler)

	go func() {
		cfg.Logger.Infof("server starting on port %s", cfg.AppPort)
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cfg.Logger.Info("shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		cfg.Logger.Fatalf("server forced shutdown: %v", err)
	}
	cfg.Logger.Info("server stopped")
}