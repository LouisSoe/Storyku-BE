package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM tags ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tag, error) {
	t := &domain.Tag{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM tags WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *tagRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Tag, error) {
	if len(ids) == 0 {
		return []domain.Tag{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(
		`SELECT id, name, slug, created_at, updated_at FROM tags WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	t := &domain.Tag{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM tags WHERE slug = $1`, slug,
	).Scan(&t.ID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *tagRepository) Create(ctx context.Context, t *domain.Tag) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tags (id, name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.Name, t.Slug, now, now,
	)
	return err
}

func (r *tagRepository) Update(ctx context.Context, t *domain.Tag) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tags SET name = $1, slug = $2, updated_at = $3 WHERE id = $4`,
		t.Name, t.Slug, time.Now(), t.ID,
	)
	return err
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tags WHERE id = $1`, id)
	return err
}