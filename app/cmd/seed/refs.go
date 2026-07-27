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

	for _, name := range defaultTypes {
		var id int

		err := conn.QueryRow(ctx, `SELECT id FROM types WHERE name = $1 LIMIT 1`, name).Scan(&id)
		switch {
		case err == nil:

		case errors.Is(err, pgx.ErrNoRows):
			if err := conn.QueryRow(ctx, `
				INSERT INTO types (name, is_active, created_at, updated_at)
				VALUES ($1, TRUE, NOW(), NOW())
				RETURNING id`, name,
			).Scan(&id); err != nil {
				return nil, fmt.Errorf("create type %q: %w", name, err)
			}
			log.Printf("Type '%s' created (id=%d)", name, id)

		default:
			return nil, fmt.Errorf("lookup type %q: %w", name, err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func firstIDOrNil(ctx context.Context, conn *pgxpool.Pool, table string) *int {
	var id int

	if err := conn.QueryRow(ctx, fmt.Sprintf(`SELECT id FROM %s LIMIT 1`, table)).Scan(&id); err != nil {
		return nil
	}

	return &id
}
