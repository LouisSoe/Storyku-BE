package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"storyku-be/config"
	"storyku-be/core/usecase"
	dbImpl "storyku-be/interfaces/database"
	httpHandler "storyku-be/interfaces/http"
	"storyku-be/routes"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.Load()

	storyRepo := dbImpl.NewStoryRepository(cfg.DB)
	chapterRepo := dbImpl.NewChapterRepository(cfg.DB)

	storyUC := usecase.NewStoryUsecase(storyRepo)
	chapterUC := usecase.NewChapterUsecase(storyRepo, chapterRepo)

	storyHandler := httpHandler.NewStoryHandler(storyUC, cfg.Logger)
	chapterHandler := httpHandler.NewChapterHandler(chapterUC, cfg.Logger)

	e := echo.New()
	e.HideBanner = true

	routes.Register(e, storyHandler, chapterHandler)

	go func() {
		cfg.Logger.WithField("port", cfg.AppPort).Info("server starting")
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			cfg.Logger.WithError(err).Fatal("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cfg.Logger.Info("shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		cfg.Logger.WithError(err).Fatal("server forced shutdown")
	}

	cfg.Logger.Info("server stopped")
}