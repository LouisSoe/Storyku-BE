package domain

import "time"

type Chapter struct {
	ChapterID  string    `json:"chapter_id"  db:"chapter_id"`
	StoryID    string    `json:"story_id"    db:"story_id"`
	Title      string    `json:"title"       db:"title"`
	Content    string    `json:"content"     db:"content"`
	OrderIndex int       `json:"order_index" db:"order_index"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
}