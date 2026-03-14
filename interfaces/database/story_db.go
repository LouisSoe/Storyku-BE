package database

import (
	"context"
	"database/sql"
	"fmt"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
	"strings"
	"time"

	"github.com/lib/pq"
)

type storyRepository struct{ db *sql.DB }

func NewStoryRepository(db *sql.DB) repository.StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) FindAll(ctx context.Context, filter repository.StoryFilter) ([]domain.Story, int64, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if filter.Search != "" {
		like := "%" + strings.ToLower(filter.Search) + "%"
		conditions = append(conditions,
			fmt.Sprintf("(LOWER(title) LIKE $%d OR LOWER(author) LIKE $%d)", idx, idx+1),
		)
		args = append(args, like, like)
		idx += 2
	}

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", idx))
		args = append(args, filter.Category)
		idx++
	}

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", idx))
		args = append(args, filter.Status)
		idx++
	}

	where := strings.Join(conditions, " AND ")

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stories WHERE %s", where)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	dataQuery := fmt.Sprintf(`
		SELECT story_id, title, author, synopsis, category, cover_url, tags, status, created_at, updated_at
		FROM stories
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, idx, idx+1)
	args = append(args, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stories []domain.Story
	for rows.Next() {
		var s domain.Story
		var tags pq.StringArray
		if err := rows.Scan(
			&s.StoryID, &s.Title, &s.Author, &s.Synopsis,
			&s.Category, &s.CoverURL, &tags,
			&s.Status, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		s.Tags = []string(tags)
		stories = append(stories, s)
	}

	return stories, total, nil
}

func (r *storyRepository) FindByID(ctx context.Context, id string) (*domain.Story, error) {
	query := `
		SELECT story_id, title, author, synopsis, category, cover_url, tags, status, created_at, updated_at
		FROM stories
		WHERE story_id = $1
	`
	var s domain.Story
	var tags pq.StringArray
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.StoryID, &s.Title, &s.Author, &s.Synopsis,
		&s.Category, &s.CoverURL, &tags,
		&s.Status, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	s.Tags = []string(tags)
	return &s, nil
}

func (r *storyRepository) Create(ctx context.Context, s *domain.Story) error {
	query := `
		INSERT INTO stories (story_id, title, author, synopsis, category, cover_url, tags, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		s.StoryID, s.Title, s.Author, s.Synopsis,
		s.Category, s.CoverURL, pq.Array(s.Tags),
		s.Status, now, now,
	)
	return err
}

func (r *storyRepository) Update(ctx context.Context, s *domain.Story) error {
	query := `
		UPDATE stories
		SET title = $1, author = $2, synopsis = $3, category = $4,
		    cover_url = $5, tags = $6, status = $7, updated_at = $8
		WHERE story_id = $9
	`
	_, err := r.db.ExecContext(ctx, query,
		s.Title, s.Author, s.Synopsis, s.Category,
		s.CoverURL, pq.Array(s.Tags), s.Status, time.Now(), s.StoryID,
	)
	return err
}

func (r *storyRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM stories WHERE story_id = $1`, id)
	return err
}