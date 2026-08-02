package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type TypeRepository struct {
	db *pgxpool.Pool
}

func NewTypeRepository(db *pgxpool.Pool) *TypeRepository {
	return &TypeRepository{db: db}
}

const typeSearchWhere = " WHERE (name_kz ILIKE $1 OR name_ru ILIKE $1 OR name_en ILIKE $1 OR slug ILIKE $1)"

func (r *TypeRepository) GetAll(ctx context.Context, search string) ([]model.Type, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = typeSearchWhere
		args = append(args, "%"+search+"%")
	}

	rows, err := r.db.Query(ctx, `SELECT id, name_kz, name_ru, name_en, slug, is_active FROM types`+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Type
	for rows.Next() {
		var t model.Type
		if err := rows.Scan(&t.ID, &t.NameKz, &t.NameRu, &t.NameEn, &t.Slug, &t.IsActive); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *TypeRepository) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]model.Type, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = typeSearchWhere
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM types`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx,
		fmt.Sprintf(`SELECT id, name_kz, name_ru, name_en, slug, is_active FROM types%s ORDER BY id LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []model.Type
	for rows.Next() {
		var t model.Type
		if err := rows.Scan(&t.ID, &t.NameKz, &t.NameRu, &t.NameEn, &t.Slug, &t.IsActive); err != nil {
			return nil, 0, err
		}
		result = append(result, t)
	}
	return result, total, rows.Err()
}

func (r *TypeRepository) GetByID(ctx context.Context, id int) (model.Type, error) {
	var t model.Type
	err := r.db.QueryRow(ctx, `SELECT id, name_kz, name_ru, name_en, slug, is_active FROM types WHERE id=$1`, id).
		Scan(&t.ID, &t.NameKz, &t.NameRu, &t.NameEn, &t.Slug, &t.IsActive)
	return t, store.MapNoRows(err, model.ErrTypeNotFound)
}

func (r *TypeRepository) Create(ctx context.Context, t model.Type) (model.Type, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO types (name_kz, name_ru, name_en, slug, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		 RETURNING id, name_kz, name_ru, name_en, slug, is_active`,
		t.NameKz, t.NameRu, t.NameEn, t.Slug, t.IsActive,
	).Scan(&t.ID, &t.NameKz, &t.NameRu, &t.NameEn, &t.Slug, &t.IsActive)
	return t, store.MapUnique(err, model.ErrTypeSlugExists)
}

func (r *TypeRepository) Update(ctx context.Context, id int, t model.Type) (model.Type, error) {
	t.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE types SET name_kz=$1, name_ru=$2, name_en=$3, slug=$4, is_active=$5, updated_at=NOW() WHERE id=$6`,
		t.NameKz, t.NameRu, t.NameEn, t.Slug, t.IsActive, id,
	)
	t.ID = id
	return t, store.MapUnique(err, model.ErrTypeSlugExists)
}

func (r *TypeRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM types WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return model.ErrTypeNotFound
	}

	return nil
}
