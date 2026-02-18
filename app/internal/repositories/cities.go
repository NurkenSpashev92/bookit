package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type CityRepository struct {
	db *pgxpool.Pool
}

func NewCityRepository(db *pgxpool.Pool) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) GetAllWithCountry(ctx context.Context) ([]schemas.City, error) {
	query := `
		SELECT 
			c.id, c.name_ru, c.name_en, c.name_kz, c.postall_code,
			ct.id, ct.name_kz, ct.name_en, ct.name_ru, ct.code
		FROM cities c
		INNER JOIN countries ct ON c.country_id = ct.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query cities: %w", err)
	}
	defer rows.Close()

	var result []schemas.City
	for rows.Next() {
		var c schemas.City
		var ct schemas.Country
		if err := rows.Scan(
			&c.ID, &c.NameRU, &c.NameEN, &c.NameKZ, &c.PostallCode,
			&ct.ID, &ct.NameKZ, &ct.NameEN, &ct.NameRU, &ct.Code,
		); err != nil {
			return nil, err
		}
		c.Country = ct
		result = append(result, c)
	}
	return result, nil
}

func (r *CityRepository) GetByIDWithCountry(ctx context.Context, id int) (schemas.City, error) {
	query := `
		SELECT 
			c.id, c.name_ru, c.name_en, c.name_kz, c.postall_code,
			ct.id, ct.name_kz, ct.name_en, ct.name_ru, ct.code
		FROM cities c
		INNER JOIN countries ct ON c.country_id = ct.id
		WHERE c.id = $1
	`

	var c schemas.City
	var ct schemas.Country
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.NameRU, &c.NameEN, &c.NameKZ, &c.PostallCode,
		&ct.ID, &ct.NameKZ, &ct.NameEN, &ct.NameRU, &ct.Code,
	)
	if err != nil {
		return c, fmt.Errorf("city not found: %w", err)
	}
	c.Country = ct
	return c, nil
}

func (r *CityRepository) GetByID(ctx context.Context, id int) (models.City, error) {
	var c models.City
	err := r.db.QueryRow(ctx, `SELECT id, name_ru, name_en, name_kz, postall_code, country_id, created_at, updated_at FROM cities WHERE id=$1`, id).
		Scan(&c.ID, &c.NameRU, &c.NameEN, &c.NameKZ, &c.PostallCode, &c.CountryID, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CityRepository) Create(ctx context.Context, req schemas.CityCreateRequest) (models.City, error) {
	var c models.City
	err := r.db.QueryRow(ctx,
		`INSERT INTO cities (name_ru, name_en, name_kz, postall_code, country_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		 RETURNING id, name_ru, name_en, name_kz, postall_code, country_id, created_at, updated_at`,
		req.NameRU, req.NameEN, req.NameKZ, req.PostallCode, req.CountryID,
	).Scan(&c.ID, &c.NameRU, &c.NameEN, &c.NameKZ, &c.PostallCode, &c.CountryID, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *CityRepository) Update(ctx context.Context, id int, req schemas.CityUpdateRequest) (models.City, error) {
	c, err := r.GetByID(ctx, id)
	if err != nil {
		return c, err
	}

	if req.NameRU != nil {
		c.NameRU = *req.NameRU
	}
	if req.NameEN != nil {
		c.NameEN = *req.NameEN
	}
	if req.NameKZ != nil {
		c.NameKZ = *req.NameKZ
	}
	if req.PostallCode != nil {
		c.PostallCode = *req.PostallCode
	}
	if req.CountryID != nil {
		c.CountryID = *req.CountryID
	}
	c.UpdatedAt = time.Now()

	_, err = r.db.Exec(ctx,
		`UPDATE cities SET name_ru=$1, name_en=$2, name_kz=$3, postall_code=$4, country_id=$5, updated_at=$6 WHERE id=$7`,
		c.NameRU, c.NameEN, c.NameKZ, c.PostallCode, c.CountryID, c.UpdatedAt, id,
	)
	return c, err
}

func (r *CityRepository) Delete(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM cities WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("city not found")
	}
	return nil
}
