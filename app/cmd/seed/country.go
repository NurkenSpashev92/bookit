package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type countrySeed struct {
	nameEN, nameKZ, nameRU, code string
}

var defaultCountries = []countrySeed{
	{nameEN: "Kazakhstan", nameKZ: "Қазақстан", nameRU: "Казахстан", code: "KZ"},
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
