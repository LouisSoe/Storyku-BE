package domain

import "time"

type DashboardSummary struct {
	TotalStories    int64 `json:"total_stories"`
	TotalPublished  int64 `json:"total_published"`
	TotalDraft      int64 `json:"total_draft"`
	TotalChapters   int64 `json:"total_chapters"`
	TotalCategories int64 `json:"total_categories"`
	TotalTags       int64 `json:"total_tags"`
}

type StoriesPerCategory struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategorySlug string `json:"category_slug"`
	Total        int64  `json:"total"`
}

type StoriesPerStatus struct {
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type TopTag struct {
	TagID   string `json:"tag_id"`
	TagName string `json:"tag_name"`
	TagSlug string `json:"tag_slug"`
	Total   int64  `json:"total"`
}

type RecentStory struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Status       string    `json:"status"`
	CoverURL     string    `json:"cover_url"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type RecentChapter struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	StoryID    string    `json:"story_id"`
	StoryTitle string    `json:"story_title"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DashboardData struct {
	Summary            DashboardSummary     `json:"summary"`
	StoriesPerCategory []StoriesPerCategory `json:"stories_per_category"`
	StoriesPerStatus   []StoriesPerStatus   `json:"stories_per_status"`
	TopTags            []TopTag             `json:"top_tags"`
	RecentStories      []RecentStory        `json:"recent_stories"`
	RecentChapters     []RecentChapter      `json:"recent_chapters"`
}