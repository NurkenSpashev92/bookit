package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ensureTypes(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultTypes))

	for _, houseType := range defaultTypes {
		id, err := ensureRow(ctx, conn, "type", houseType.nameEN,
			`SELECT id FROM types WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO types (name_en, name_kz, name_ru, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, TRUE, NOW(), NOW())
			 RETURNING id`,
			houseType.nameEN, houseType.nameKZ, houseType.nameRU,
		)
		if err != nil {
			return nil, err
		}

		if err := fillTypeTranslations(ctx, conn, id, houseType); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// fillTypeTranslations backfills kz/ru names for types created before the
// multilingual columns existed, where they were copied from the english name.
func fillTypeTranslations(ctx context.Context, conn *pgxpool.Pool, id int, houseType typeSeed) error {
	cmd, err := conn.Exec(ctx, `
		UPDATE types
		SET name_kz = $1, name_ru = $2, updated_at = NOW()
		WHERE id = $3 AND (name_kz <> $1 OR name_ru <> $2)`,
		houseType.nameKZ, houseType.nameRU, id,
	)
	if err != nil {
		return fmt.Errorf("translate type %q: %w", houseType.nameEN, err)
	}

	if cmd.RowsAffected() > 0 {
		log.Printf("Type '%s' translations updated (id=%d)", houseType.nameEN, id)
	}

	return nil
}

func ensureCountries(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultCountries))

	for _, country := range defaultCountries {
		id, err := ensureRow(ctx, conn, "country", country.nameEN,
			`SELECT id FROM countries WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO countries (name_en, name_kz, name_ru, code, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, NOW(), NOW())
			 RETURNING id`,
			country.nameEN, country.nameKZ, country.nameRU, country.code,
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

type cityRef struct {
	id        int
	countryID int
}

func ensureCities(ctx context.Context, conn *pgxpool.Pool, countryID int) ([]cityRef, error) {
	refs := make([]cityRef, 0, len(defaultCities))

	for _, city := range defaultCities {
		id, err := ensureRow(ctx, conn, "city", city.nameEN,
			`SELECT id FROM cities WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO cities (name_en, name_kz, name_ru, postall_code, country_id, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			 RETURNING id`,
			city.nameEN, city.nameKZ, city.nameRU, city.postalCode, countryID,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, cityRef{id: id, countryID: countryID})
	}

	return refs, nil
}

func ensureCategories(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultCategories))

	for _, category := range defaultCategories {
		id, err := ensureRow(ctx, conn, "category", category.nameEN,
			`SELECT id FROM categories WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO categories (name_en, name_kz, name_ru, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, TRUE, NOW(), NOW())
			 RETURNING id`,
			category.nameEN, category.nameKZ, category.nameRU,
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func ensureConveniences(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultConveniences))

	for _, name := range defaultConveniences {
		id, err := ensureRow(ctx, conn, "convenience", name,
			`SELECT id FROM conveniences WHERE name = $1 LIMIT 1`,
			`INSERT INTO conveniences (name, is_active, created_at, updated_at)
			 VALUES ($1, TRUE, NOW(), NOW())
			 RETURNING id`,
			name,
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func ensureRow(
	ctx context.Context,
	conn *pgxpool.Pool,
	kind, name, selectSQL, insertSQL string,
	insertArgs ...any,
) (int, error) {
	var id int

	err := conn.QueryRow(ctx, selectSQL, name).Scan(&id)
	switch {
	case err == nil:
		return id, nil

	case errors.Is(err, pgx.ErrNoRows):
		if err := conn.QueryRow(ctx, insertSQL, insertArgs...).Scan(&id); err != nil {
			return 0, fmt.Errorf("create %s %q: %w", kind, name, err)
		}
		log.Printf("%s '%s' created (id=%d)", kind, name, id)
		return id, nil

	default:
		return 0, fmt.Errorf("lookup %s %q: %w", kind, name, err)
	}
}
