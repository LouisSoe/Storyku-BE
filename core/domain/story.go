package domain

import (
	"time"

	"github.com/google/uuid"
)

type StoryStatus string

const (
	StatusPublish StoryStatus = "publish"
	StatusDraft   StoryStatus = "draft"
)

type Story struct {
	ID         uuid.UUID   `json:"id"          db:"id"`
	CategoryID *uuid.UUID  `json:"category_id" db:"category_id"`
	Title      string      `json:"title"       db:"title"`
	Author     string      `json:"author"      db:"author"`
	Synopsis   string      `json:"synopsis"    db:"synopsis"`
	CoverURL   string      `json:"cover_url"   db:"cover_url"`
	Status     StoryStatus `json:"status"      db:"status"`
	CreatedAt  time.Time   `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"  db:"updated_at"`
}

type StoryDetail struct {
	Story
	Category *Category `json:"category"`
	Tags     []Tag     `json:"tags"`
}

func (s *Story) IsValidStatus() bool {
	return s.Status == StatusPublish || s.Status == StatusDraft
}