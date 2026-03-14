package database

import (
	"context"
	"database/sql"
	"storyku-be/core/domain"
	"storyku-be/core/repository"
)

type dashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) repository.DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetSummary mengambil ringkasan angka utama dalam satu query efisien
func (r *dashboardRepository) GetSummary(ctx context.Context) (domain.DashboardSummary, error) {
	query := `
		SELECT
			COUNT(*)                                          AS total_stories,
			COUNT(*) FILTER (WHERE status = 'publish')       AS total_published,
			COUNT(*) FILTER (WHERE status = 'draft')         AS total_draft,
			(SELECT COUNT(*) FROM chapters)                   AS total_chapters,
			(SELECT COUNT(*) FROM categories)                 AS total_categories,
			(SELECT COUNT(*) FROM tags)                       AS total_tags
		FROM stories
	`
	var s domain.DashboardSummary
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.TotalStories,
		&s.TotalPublished,
		&s.TotalDraft,
		&s.TotalChapters,
		&s.TotalCategories,
		&s.TotalTags,
	)
	return s, err
}

// GetStoriesPerCategory distribusi story per category
func (r *dashboardRepository) GetStoriesPerCategory(ctx context.Context) ([]domain.StoriesPerCategory, error) {
	query := `
		SELECT
			c.id,
			c.name,
			c.slug,
			COUNT(s.id) AS total
		FROM categories c
		LEFT JOIN stories s ON s.category_id = c.id
		GROUP BY c.id, c.name, c.slug
		ORDER BY total DESC, c.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.StoriesPerCategory
	for rows.Next() {
		var item domain.StoriesPerCategory
		if err := rows.Scan(
			&item.CategoryID,
			&item.CategoryName,
			&item.CategorySlug,
			&item.Total,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// GetStoriesPerStatus distribusi publish vs draft
func (r *dashboardRepository) GetStoriesPerStatus(ctx context.Context) ([]domain.StoriesPerStatus, error) {
	query := `
		SELECT status, COUNT(*) AS total
		FROM stories
		GROUP BY status
		ORDER BY status ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.StoriesPerStatus
	for rows.Next() {
		var item domain.StoriesPerStatus
		if err := rows.Scan(&item.Status, &item.Total); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// GetTopTags tag yang paling banyak dipakai di story
func (r *dashboardRepository) GetTopTags(ctx context.Context, limit int) ([]domain.TopTag, error) {
	query := `
		SELECT
			t.id,
			t.name,
			t.slug,
			COUNT(st.story_id) AS total
		FROM tags t
		LEFT JOIN story_tags st ON st.tag_id = t.id
		GROUP BY t.id, t.name, t.slug
		ORDER BY total DESC, t.name ASC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.TopTag
	for rows.Next() {
		var item domain.TopTag
		if err := rows.Scan(
			&item.TagID,
			&item.TagName,
			&item.TagSlug,
			&item.Total,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// GetRecentStories 5 story terbaru dengan info category
func (r *dashboardRepository) GetRecentStories(ctx context.Context, limit int) ([]domain.RecentStory, error) {
	query := `
		SELECT
			s.id,
			s.title,
			s.author,
			s.status,
			COALESCE(s.cover_url, ''),
			COALESCE(c.name, ''),
			s.created_at
		FROM stories s
		LEFT JOIN categories c ON c.id = s.category_id
		ORDER BY s.created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.RecentStory
	for rows.Next() {
		var item domain.RecentStory
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Author,
			&item.Status,
			&item.CoverURL,
			&item.CategoryName,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// GetRecentChapters 5 chapter terbaru dengan judul story-nya
func (r *dashboardRepository) GetRecentChapters(ctx context.Context, limit int) ([]domain.RecentChapter, error) {
	query := `
		SELECT
			ch.id,
			ch.title,
			ch.story_id,
			s.title,
			ch.updated_at
		FROM chapters ch
		JOIN stories s ON s.id = ch.story_id
		ORDER BY ch.updated_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.RecentChapter
	for rows.Next() {
		var item domain.RecentChapter
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.StoryID,
			&item.StoryTitle,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}