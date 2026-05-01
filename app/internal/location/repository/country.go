package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
)

type CountryRepository struct {
	db *pgxpool.Pool
}

func NewCountryRepository(db *pgxpool.Pool) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetAll(ctx context.Context) ([]model.Country, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name_kz, name_en, name_ru, code, created_at, updated_at FROM countries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country
	for rows.Next() {
		var c model.Country
		err := rows.Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return countries, nil
}

func (r *CountryRepository) GetByID(ctx context.Context, id int) (model.Country, error) {
	var c model.Country
	err := r.db.QueryRow(ctx, `SELECT id, name_kz, name_en, name_ru, code, created_at, updated_at FROM countries WHERE id=$1`, id).
		Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CountryRepository) Create(ctx context.Context, req schema.CountryCreateRequest) (model.Country, error) {
	var c model.Country
	err := r.db.QueryRow(ctx,
		`INSERT INTO countries (name_kz, name_en, name_ru, code, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,NOW(),NOW())
		 RETURNING id, name_kz, name_en, name_ru, code, created_at, updated_at`,
		req.NameKZ, req.NameEN, req.NameRU, req.Code,
	).Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CountryRepository) Update(ctx context.Context, id int, req schema.CountryUpdateRequest) (model.Country, error) {
	c, err := r.GetByID(ctx, id)
	if err != nil {
		return c, err
	}

	if req.NameKZ != nil {
		c.NameKZ = *req.NameKZ
	}
	if req.NameEN != nil {
		c.NameEN = *req.NameEN
	}
	if req.NameRU != nil {
		c.NameRU = *req.NameRU
	}
	if req.Code != nil {
		c.Code = *req.Code
	}
	c.UpdatedAt = time.Now()

	_, err = r.db.Exec(ctx,
		`UPDATE countries SET name_kz=$1, name_en=$2, name_ru=$3, code=$4, updated_at=$5 WHERE id=$6`,
		c.NameKZ, c.NameEN, c.NameRU, c.Code, c.UpdatedAt, id,
	)
	return c, err
}

func (r *CountryRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM countries WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("country not found")
	}
	return nil
}
