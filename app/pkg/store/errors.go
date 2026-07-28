package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func MapNoRows(err error, notFound error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return notFound
	}

	return err
}
