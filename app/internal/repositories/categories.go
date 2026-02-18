package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]schemas.CategoryPaginate, error) {
	query := `
		SELECT id, name_kz, name_ru, name_en, icon
		FROM categories
		WHERE is_active = TRUE
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var categories []schemas.CategoryPaginate

	for rows.Next() {
		var c schemas.CategoryPaginate
		var icon *string
		err := rows.Scan(&c.Id, &c.NameKz, &c.NameRu, &c.NameEn, &icon)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		c.Icon = icon
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (models.Category, error) {
	query := `
		SELECT 
			id,
			name_kz,
			name_ru, 
			name_en, 
			is_active, 
			icon, 
			created_at, 
			updated_at
		FROM categories
		WHERE id = $1
	`

	var category models.Category

	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.IsActive,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return category, err
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, req schemas.CategoryCreateRequest) (models.Category, error) {
	query := `
		INSERT INTO categories (name_kz, name_ru, name_en, icon, is_active)
		VALUES ($1, $2, $3, $4, COALESCE($5, TRUE))
		RETURNING id, name_kz, name_ru, name_en, is_active, icon, created_at, updated_at
	`

	var category models.Category

	err := r.db.QueryRow(
		ctx,
		query,
		req.NameKz,
		req.NameRu,
		req.NameEn,
		req.Icon,
		req.IsActive,
	).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.IsActive,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return category, fmt.Errorf("faild to insert category %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id int, req schemas.CategoryUpdateRequest, icon *string) (models.Category, *string, error) {
	old, err := r.GetByID(ctx, id)
	if err != nil {
		return models.Category{}, nil, err
	}

	query := `
		UPDATE categories
		SET
			name_kz = COALESCE($1, name_kz),
			name_ru = COALESCE($2, name_ru),
			name_en = COALESCE($3, name_en),
			is_active = COALESCE($4, is_active),
			icon = COALESCE($5, icon),
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, name_kz, name_ru, name_en, is_active, icon, created_at, updated_at
	`

	var category models.Category

	err = r.db.QueryRow(
		ctx,
		query,
		req.NameKz,
		req.NameRu,
		req.NameEn,
		req.IsActive,
		icon,
		id,
	).Scan(
		&category.ID,
		&category.NameKz,
		&category.NameRu,
		&category.NameEn,
		&category.IsActive,
		&category.Icon,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	return category, old.Icon, err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) (*string, error) {
	var icon *string

	err := r.db.QueryRow(ctx,
		`SELECT icon FROM categories WHERE id = $1`,
		id,
	).Scan(&icon)

	if err != nil {
		return nil, err
	}

	cmd, err := r.db.Exec(ctx,
		`DELETE FROM categories WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	if cmd.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	return icon, nil
}
