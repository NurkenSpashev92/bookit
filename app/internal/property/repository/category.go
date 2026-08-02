package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

const categorySearchWhere = " AND (name_kz ILIKE $1 OR name_ru ILIKE $1 OR name_en ILIKE $1 OR slug ILIKE $1)"

func (r *CategoryRepository) GetCategories(ctx context.Context, search string) ([]schema.CategoryPaginate, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = categorySearchWhere
		args = append(args, "%"+search+"%")
	}

	query := `
		SELECT id, name_kz, name_ru, name_en, slug, is_active
		FROM categories
		WHERE is_active = TRUE` + where

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var categories []schema.CategoryPaginate

	for rows.Next() {
		var c schema.CategoryPaginate
		err := rows.Scan(&c.Id, &c.NameKz, &c.NameRu, &c.NameEn, &c.Slug, &c.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) GetCategoriesPaginated(ctx context.Context, search string, limit, offset int) ([]schema.CategoryPaginate, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = categorySearchWhere
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories WHERE is_active = TRUE`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, name_kz, name_ru, name_en, slug, is_active
		FROM categories
		WHERE is_active = TRUE%s
		ORDER BY id
		LIMIT $%d OFFSET $%d`,
		where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var categories []schema.CategoryPaginate

	for rows.Next() {
		var c schema.CategoryPaginate
		if err := rows.Scan(&c.Id, &c.NameKz, &c.NameRu, &c.NameEn, &c.Slug, &c.IsActive); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, total, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (model.Category, error) {
	query := `
		SELECT
			id,
			name_kz,
			name_ru,
			name_en,
			slug,
			is_active,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	var category model.Category

	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.Slug,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return category, store.MapNoRows(err, model.ErrCategoryNotFound)
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, req schema.CategoryCreateRequest) (model.Category, error) {
	query := `
		INSERT INTO categories (name_kz, name_ru, name_en, slug, is_active)
		VALUES ($1, $2, $3, $4, COALESCE($5, TRUE))
		RETURNING id, name_kz, name_ru, name_en, slug, is_active, created_at, updated_at
	`

	var category model.Category

	err := r.db.QueryRow(
		ctx,
		query,
		req.NameKz,
		req.NameRu,
		req.NameEn,
		req.Slug,
		req.IsActive,
	).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.Slug,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return category, store.MapUnique(err, model.ErrCategorySlugExists)
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int, req schema.CategoryUpdateRequest) (model.Category, error) {
	query := `
		UPDATE categories
		SET
			name_kz = COALESCE($1, name_kz),
			name_ru = COALESCE($2, name_ru),
			name_en = COALESCE($3, name_en),
			slug = COALESCE($4, slug),
			is_active = COALESCE($5, is_active),
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, name_kz, name_ru, name_en, slug, is_active, created_at, updated_at
	`

	var category model.Category

	err := r.db.QueryRow(
		ctx,
		query,
		req.NameKz,
		req.NameRu,
		req.NameEn,
		req.Slug,
		req.IsActive,
		id,
	).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.Slug,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return category, store.MapUnique(store.MapNoRows(err, model.ErrCategoryNotFound), model.ErrCategorySlugExists)
	}

	return category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM categories WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return model.ErrCategoryNotFound
	}

	return nil
}
