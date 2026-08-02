package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type typeSeed struct {
	nameEN, nameKZ, nameRU string
}

var defaultTypes = []typeSeed{
	{nameEN: "Apartment", nameKZ: "Пәтер", nameRU: "Квартира"},
	{nameEN: "House", nameKZ: "Үй", nameRU: "Дом"},
	{nameEN: "Villa", nameKZ: "Вилла", nameRU: "Вилла"},
	{nameEN: "Cottage", nameKZ: "Коттедж", nameRU: "Коттедж"},
	{nameEN: "Hotel", nameKZ: "Қонақ үй", nameRU: "Отель"},
	{nameEN: "Studio", nameKZ: "Студия", nameRU: "Студия"},
	{nameEN: "Townhouse", nameKZ: "Таунхаус", nameRU: "Таунхаус"},
	{nameEN: "Guest house", nameKZ: "Қонақжай үй", nameRU: "Гостевой дом"},
	{nameEN: "Hostel", nameKZ: "Хостел", nameRU: "Хостел"},
	{nameEN: "Yurt", nameKZ: "Киіз үй", nameRU: "Юрта"},
}

func ensureTypes(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultTypes))

	for _, houseType := range defaultTypes {
		id, err := ensureRow(ctx, conn, "type", houseType.nameEN,
			`SELECT id FROM types WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO types (name_en, name_kz, name_ru, slug, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, TRUE, NOW(), NOW())
			 RETURNING id`,
			houseType.nameEN, houseType.nameKZ, houseType.nameRU,
			utils.GenerateSlug("", houseType.nameEN, houseType.nameKZ, houseType.nameRU),
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
