package domain

import "time"

type StoryCategory string
type StoryStatus string

const (
	CategoryFinancial  StoryCategory = "Financial"
	CategoryTechnology StoryCategory = "Technology"
	CategoryHealth     StoryCategory = "Health"
)

const (
	StatusPublish StoryStatus = "publish"
	StatusDraft   StoryStatus = "draft"
)

type Story struct {
	StoryID   string    `json:"story_id"   db:"story_id"`
	Title     string    `json:"title"      db:"title"`
	Author    string    `json:"author"     db:"author"`
	Synopsis  string    `json:"synopsis"   db:"synopsis"`
	Category  string    `json:"category"   db:"category"`
	CoverURL  string    `json:"cover_url"  db:"cover_url"`
	Tags      []string  `json:"tags"       db:"tags"`
	Status    string    `json:"status"     db:"status"`
	Chapters  []Chapter `json:"chapters,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Story) IsValidCategory() bool {
	return s.Category == string(CategoryFinancial) ||
		s.Category == string(CategoryTechnology) ||
		s.Category == string(CategoryHealth)
}

func (s *Story) IsValidStatus() bool {
	return s.Status == string(StatusPublish) || s.Status == string(StatusDraft)
}