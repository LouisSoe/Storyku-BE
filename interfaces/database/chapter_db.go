package database

import (
	"context"
	"database/sql"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"time"
)

type chapterRepository struct{ db *sql.DB }

func NewChapterRepository(db *sql.DB) repository.ChapterRepository {
	return &chapterRepository{db: db}
}

func (r *chapterRepository) FindByStoryID(ctx context.Context, storyID string) ([]domain.Chapter, error) {
	query := `
		SELECT chapter_id, story_id, title, content, order_index, created_at, updated_at
		FROM chapters
		WHERE story_id = $1
		ORDER BY order_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, storyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []domain.Chapter
	for rows.Next() {
		var c domain.Chapter
		if err := rows.Scan(
			&c.ChapterID, &c.StoryID, &c.Title,
			&c.Content, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		chapters = append(chapters, c)
	}
	return chapters, nil
}

func (r *chapterRepository) FindByID(ctx context.Context, id string) (*domain.Chapter, error) {
	query := `
		SELECT chapter_id, story_id, title, content, order_index, created_at, updated_at
		FROM chapters
		WHERE chapter_id = $1
	`
	var c domain.Chapter
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ChapterID, &c.StoryID, &c.Title,
		&c.Content, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) Create(ctx context.Context, c *domain.Chapter) error {
	query := `
		INSERT INTO chapters (chapter_id, story_id, title, content, order_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		c.ChapterID, c.StoryID, c.Title, c.Content, c.OrderIndex, now, now,
	)
	return err
}

func (r *chapterRepository) Update(ctx context.Context, c *domain.Chapter) error {
	query := `
		UPDATE chapters
		SET title = $1, content = $2, updated_at = $3
		WHERE chapter_id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		c.Title, c.Content, time.Now(), c.ChapterID,
	)
	return err
}

func (r *chapterRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chapters WHERE chapter_id = $1`, id)
	return err
}

func (r *chapterRepository) CountByStoryID(ctx context.Context, storyID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM chapters WHERE story_id = $1`, storyID,
	).Scan(&count)
	return count, err
}