package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/property/model"
)

type TypeRepository struct {
	db *pgxpool.Pool
}

func NewTypeRepository(db *pgxpool.Pool) *TypeRepository {
	return &TypeRepository{db: db}
}

func (r *TypeRepository) GetAll(ctx context.Context) ([]model.Type, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, icon, is_active FROM types`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Type
	for rows.Next() {
		var t model.Type
		if err := rows.Scan(&t.ID, &t.Name, &t.Icon, &t.IsActive); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TypeRepository) GetByID(ctx context.Context, id int) (model.Type, error) {
	var t model.Type
	err := r.db.QueryRow(ctx, `SELECT id, name, icon, is_active FROM types WHERE id=$1`, id).
		Scan(&t.ID, &t.Name, &t.Icon, &t.IsActive)
	return t, err
}

func (r *TypeRepository) Create(ctx context.Context, t model.Type) (model.Type, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO types (name, icon, is_active, created_at, updated_at)
		 VALUES ($1,$2,$3,NOW(),NOW())
		 RETURNING id, name, icon, is_active`,
		t.Name, t.Icon, t.IsActive,
	).Scan(&t.ID, &t.Name, &t.Icon, &t.IsActive)
	return t, err
}

func (r *TypeRepository) Update(ctx context.Context, id int, t model.Type) (model.Type, error) {
	t.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE types SET name=$1, icon=$2, is_active=$3, updated_at=NOW() WHERE id=$4`,
		t.Name, t.Icon, t.IsActive, id,
	)
	t.ID = id
	return t, err
}

func (r *TypeRepository) Delete(ctx context.Context, id int) (string, error) {
	var icon string
	err := r.db.QueryRow(ctx, `SELECT icon FROM types WHERE id=$1`, id).Scan(&icon)
	if err != nil {
		return "", err
	}

	cmd, err := r.db.Exec(ctx, `DELETE FROM types WHERE id=$1`, id)
	if err != nil {
		return "", err
	}
	if cmd.RowsAffected() == 0 {
		return "", fmt.Errorf("type not found")
	}

	return icon, nil
}
