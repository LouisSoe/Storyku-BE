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

type storyRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) repository.StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) FindAll(ctx context.Context, filter repository.StoryFilter) ([]domain.StoryDetail, int64, error) {
	args := []interface{}{}
	idx := 1
	conds := []string{"1=1"}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"
		conds = append(conds, fmt.Sprintf("(LOWER(s.title) LIKE $%d OR LOWER(s.author) LIKE $%d)", idx, idx+1))
		args = append(args, search, search)
		idx += 2
	}
	if filter.CategoryID != "" {
		conds = append(conds, fmt.Sprintf("s.category_id = $%d", idx))
		args = append(args, filter.CategoryID)
		idx++
	}
	if filter.Status != "" {
		conds = append(conds, fmt.Sprintf("s.status = $%d", idx))
		args = append(args, filter.Status)
		idx++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stories s "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	dataSQL := fmt.Sprintf(`
		SELECT s.id, s.category_id, s.title, s.author, s.synopsis, s.cover_url, s.status, s.created_at, s.updated_at,
		       c.id, c.name, c.slug
		FROM stories s
		LEFT JOIN categories c ON s.category_id = c.id
		%s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, idx, idx+1)

	rows, err := r.db.QueryContext(ctx, dataSQL, append(args, filter.Limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var details []domain.StoryDetail
	for rows.Next() {
		d := domain.StoryDetail{}
		var (
			catID   *uuid.UUID
			catName *string
			catSlug *string
		)
		if err := rows.Scan(
			&d.ID, &d.CategoryID, &d.Title, &d.Author, &d.Synopsis,
			&d.CoverURL, &d.Status, &d.CreatedAt, &d.UpdatedAt,
			&catID, &catName, &catSlug,
		); err != nil {
			return nil, 0, err
		}
		if catID != nil {
			d.Category = &domain.Category{ID: *catID, Name: *catName, Slug: *catSlug}
		}
		d.Tags = []domain.Tag{}
		details = append(details, d)
	}

	for i := range details {
		tags, err := r.fetchTagsByStoryID(ctx, details[i].ID)
		if err != nil {
			return nil, 0, err
		}
		details[i].Tags = tags
	}

	return details, total, nil
}

func (r *storyRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.StoryDetail, error) {
	query := `
		SELECT s.id, s.category_id, s.title, s.author, s.synopsis, s.cover_url, s.status, s.created_at, s.updated_at,
		       c.id, c.name, c.slug
		FROM stories s
		LEFT JOIN categories c ON s.category_id = c.id
		WHERE s.id = $1
	`
	d := &domain.StoryDetail{}
	var (
		catID   *uuid.UUID
		catName *string
		catSlug *string
	)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.CategoryID, &d.Title, &d.Author, &d.Synopsis,
		&d.CoverURL, &d.Status, &d.CreatedAt, &d.UpdatedAt,
		&catID, &catName, &catSlug,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if catID != nil {
		d.Category = &domain.Category{ID: *catID, Name: *catName, Slug: *catSlug}
	}

	tags, err := r.fetchTagsByStoryID(ctx, id)
	if err != nil {
		return nil, err
	}
	d.Tags = tags
	return d, nil
}

func (r *storyRepository) Create(ctx context.Context, s *domain.Story, tagIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	_, err = tx.ExecContext(ctx, `
		INSERT INTO stories (id, category_id, title, author, synopsis, cover_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, s.ID, s.CategoryID, s.Title, s.Author, s.Synopsis, s.CoverURL, s.Status, now, now)
	if err != nil {
		return err
	}

	if err := insertStoryTags(ctx, tx, s.ID, tagIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *storyRepository) Update(ctx context.Context, s *domain.Story, tagIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE stories
		SET category_id = $1, title = $2, author = $3, synopsis = $4,
		    cover_url = $5, status = $6, updated_at = $7
		WHERE id = $8
	`, s.CategoryID, s.Title, s.Author, s.Synopsis, s.CoverURL, s.Status, time.Now(), s.ID)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM story_tags WHERE story_id = $1`, s.ID); err != nil {
		return err
	}
	if err := insertStoryTags(ctx, tx, s.ID, tagIDs); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *storyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM stories WHERE id = $1`, id)
	return err
}

func (r *storyRepository) fetchTagsByStoryID(ctx context.Context, storyID uuid.UUID) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.name, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN story_tags st ON st.tag_id = t.id
		WHERE st.story_id = $1
		ORDER BY t.name ASC
	`, storyID)
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
	if tags == nil {
		return []domain.Tag{}, nil
	}
	return tags, nil
}

func insertStoryTags(ctx context.Context, tx *sql.Tx, storyID uuid.UUID, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO story_tags (story_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			storyID, tagID,
		); err != nil {
			return err
		}
	}
	return nil
}