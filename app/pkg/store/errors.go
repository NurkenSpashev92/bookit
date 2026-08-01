package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

func MapNoRows(err error, notFound error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return notFound
	}

	return err
}

func MapUnique(err error, conflict error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return conflict
	}

	return err
}
