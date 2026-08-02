package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type ConvenienceRepository struct {
	db *pgxpool.Pool
}

func NewConvenienceRepository(db *pgxpool.Pool) *ConvenienceRepository {
	return &ConvenienceRepository{db: db}
}

const convenienceSearchWhere = " AND (name ILIKE $1 OR slug ILIKE $1)"

func (r *ConvenienceRepository) GetConveniences(ctx context.Context, search string) ([]schema.ConveniencePaginate, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = convenienceSearchWhere
		args = append(args, "%"+search+"%")
	}

	query := `
		SELECT id, name, slug, is_active
		FROM conveniences
		WHERE is_active = TRUE` + where

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var conveniences []schema.ConveniencePaginate

	for rows.Next() {
		var c schema.ConveniencePaginate
		if err := rows.Scan(&c.Id, &c.Name, &c.Slug, &c.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		conveniences = append(conveniences, c)
	}

	return conveniences, rows.Err()
}

func (r *ConvenienceRepository) GetConveniencesPaginated(ctx context.Context, search string, limit, offset int) ([]schema.ConveniencePaginate, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = convenienceSearchWhere
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM conveniences WHERE is_active = TRUE`+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count conveniences: %w", err)
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT id, name, slug, is_active
		FROM conveniences
		WHERE is_active = TRUE%s
		ORDER BY id
		LIMIT $%d OFFSET $%d`,
		where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var conveniences []schema.ConveniencePaginate

	for rows.Next() {
		var c schema.ConveniencePaginate
		if err := rows.Scan(&c.Id, &c.Name, &c.Slug, &c.IsActive); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		conveniences = append(conveniences, c)
	}

	return conveniences, total, rows.Err()
}

func (r *ConvenienceRepository) GetByID(ctx context.Context, id int) (schema.Convenience, error) {
	query := `
		SELECT id, name, slug, is_active
		FROM conveniences
		WHERE id = $1
	`

	var c schema.Convenience

	err := r.db.QueryRow(ctx, query, id).Scan(&c.Id, &c.Name, &c.Slug, &c.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c, model.ErrConvenienceNotFound
		}
		return c, err
	}

	return c, nil
}

func (r *ConvenienceRepository) CreateConvenience(ctx context.Context, req schema.ConvenienceCreateRequest) (schema.Convenience, error) {
	query := `
		INSERT INTO conveniences (name, slug, is_active)
		VALUES ($1, $2, COALESCE($3, TRUE))
		RETURNING id, name, slug, is_active
	`

	var c schema.Convenience

	err := r.db.QueryRow(ctx, query, req.Name, req.Slug, req.IsActive).Scan(&c.Id, &c.Name, &c.Slug, &c.IsActive)
	if err != nil {
		return c, store.MapUnique(err, model.ErrConvenienceSlugExists)
	}

	return c, nil
}

func (r *ConvenienceRepository) Update(ctx context.Context, id int, req schema.ConvenienceUpdateRequest) (schema.Convenience, error) {
	query := `
		UPDATE conveniences
		SET
			name = COALESCE($1, name),
			slug = COALESCE($2, slug),
			is_active = COALESCE($3, is_active),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, slug, is_active
	`

	var c schema.Convenience

	err := r.db.QueryRow(ctx, query, req.Name, req.Slug, req.IsActive, id).Scan(&c.Id, &c.Name, &c.Slug, &c.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c, model.ErrConvenienceNotFound
		}
		return c, store.MapUnique(err, model.ErrConvenienceSlugExists)
	}

	return c, nil
}

func (r *ConvenienceRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM conveniences WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return model.ErrConvenienceNotFound
	}

	return nil
}
