package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM categories ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	c := &domain.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	c := &domain.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, slug, created_at, updated_at FROM categories WHERE slug = $1`, slug,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *categoryRepository) Create(ctx context.Context, c *domain.Category) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO categories (id, name, slug, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		c.ID, c.Name, c.Slug, now, now,
	)
	return err
}

func (r *categoryRepository) Update(ctx context.Context, c *domain.Category) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = $1, slug = $2, updated_at = $3 WHERE id = $4`,
		c.Name, c.Slug, time.Now(), c.ID,
	)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}