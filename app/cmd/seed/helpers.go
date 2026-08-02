package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}
