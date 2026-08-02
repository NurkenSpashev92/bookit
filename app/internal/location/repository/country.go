package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/location/model"
	"github.com/nurkenspashev92/bookit/internal/location/schema"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type CountryRepository struct {
	db *pgxpool.Pool
}

func NewCountryRepository(db *pgxpool.Pool) *CountryRepository {
	return &CountryRepository{db: db}
}

const countrySearchWhere = " WHERE (name_kz ILIKE $1 OR name_en ILIKE $1 OR name_ru ILIKE $1 OR code ILIKE $1)"

func (r *CountryRepository) GetAll(ctx context.Context, search string) ([]model.Country, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = countrySearchWhere
		args = append(args, "%"+search+"%")
	}

	rows, err := r.db.Query(ctx, `SELECT id, name_kz, name_en, name_ru, code, slug, created_at, updated_at FROM countries`+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []model.Country
	for rows.Next() {
		var c model.Country
		err := rows.Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return countries, nil
}

func (r *CountryRepository) GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]model.Country, int, error) {
	var args []interface{}
	where := ""
	if search != "" {
		where = countrySearchWhere
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM countries`+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx,
		fmt.Sprintf(`SELECT id, name_kz, name_en, name_ru, code, slug, created_at, updated_at FROM countries%s ORDER BY id LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var countries []model.Country
	for rows.Next() {
		var c model.Country
		if err := rows.Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.Slug, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		countries = append(countries, c)
	}
	return countries, total, rows.Err()
}

func (r *CountryRepository) GetByID(ctx context.Context, id int) (model.Country, error) {
	var c model.Country
	err := r.db.QueryRow(ctx, `SELECT id, name_kz, name_en, name_ru, code, slug, created_at, updated_at FROM countries WHERE id=$1`, id).
		Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CountryRepository) Create(ctx context.Context, req schema.CountryCreateRequest) (model.Country, error) {
	var c model.Country
	slug := utils.GenerateSlug(req.Slug, req.NameEN, req.NameKZ, req.NameRU)
	err := r.db.QueryRow(ctx,
		`INSERT INTO countries (name_kz, name_en, name_ru, code, slug, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		 RETURNING id, name_kz, name_en, name_ru, code, slug, created_at, updated_at`,
		req.NameKZ, req.NameEN, req.NameRU, req.Code, slug,
	).Scan(&c.ID, &c.NameKZ, &c.NameEN, &c.NameRU, &c.Code, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
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
	if req.Slug != nil {
		c.Slug = *req.Slug
	}
	c.UpdatedAt = time.Now()

	_, err = r.db.Exec(ctx,
		`UPDATE countries SET name_kz=$1, name_en=$2, name_ru=$3, code=$4, slug=$5, updated_at=$6 WHERE id=$7`,
		c.NameKZ, c.NameEN, c.NameRU, c.Code, c.Slug, c.UpdatedAt, id,
	)
	return c, err
}

func (r *CountryRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM countries WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return model.ErrCountryNotFound
	}
	return nil
}
