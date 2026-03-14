package usecase

import (
	"context"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type DashboardUsecase interface {
	GetDashboard(ctx context.Context) (*domain.DashboardData, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

func (u *dashboardUsecase) GetDashboard(ctx context.Context) (*domain.DashboardData, error) {
	summary, err := u.repo.GetSummary(ctx)
	if err != nil {
		return nil, err
	}

	perCategory, err := u.repo.GetStoriesPerCategory(ctx)
	if err != nil {
		return nil, err
	}

	perStatus, err := u.repo.GetStoriesPerStatus(ctx)
	if err != nil {
		return nil, err
	}

	topTags, err := u.repo.GetTopTags(ctx, 5)
	if err != nil {
		return nil, err
	}

	recentStories, err := u.repo.GetRecentStories(ctx, 5)
	if err != nil {
		return nil, err
	}

	recentChapters, err := u.repo.GetRecentChapters(ctx, 5)
	if err != nil {
		return nil, err
	}

	// Pastikan slice tidak nil — frontend lebih mudah handle array kosong
	if perCategory == nil {
		perCategory = []domain.StoriesPerCategory{}
	}
	if perStatus == nil {
		perStatus = []domain.StoriesPerStatus{}
	}
	if topTags == nil {
		topTags = []domain.TopTag{}
	}
	if recentStories == nil {
		recentStories = []domain.RecentStory{}
	}
	if recentChapters == nil {
		recentChapters = []domain.RecentChapter{}
	}

	return &domain.DashboardData{
		Summary:            summary,
		StoriesPerCategory: perCategory,
		StoriesPerStatus:   perStatus,
		TopTags:            topTags,
		RecentStories:      recentStories,
		RecentChapters:     recentChapters,
	}, nil
}