package repository

import (
	"context"
	"storyku-be/core/domain"
)

type DashboardRepository interface {
	GetSummary(ctx context.Context) (domain.DashboardSummary, error)
	GetStoriesPerCategory(ctx context.Context) ([]domain.StoriesPerCategory, error)
	GetStoriesPerStatus(ctx context.Context) ([]domain.StoriesPerStatus, error)
	GetTopTags(ctx context.Context, limit int) ([]domain.TopTag, error)
	GetRecentStories(ctx context.Context, limit int) ([]domain.RecentStory, error)
	GetRecentChapters(ctx context.Context, limit int) ([]domain.RecentChapter, error)
}