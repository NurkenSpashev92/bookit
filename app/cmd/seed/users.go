package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func upsertUser(
	ctx context.Context, conn *pgxpool.Pool, email, firstName, passwordHash string, isSuperuser bool,
) (int, bool, error) {
	var id int

	err := conn.QueryRow(ctx, `SELECT id FROM users WHERE email = $1 LIMIT 1`, email).Scan(&id)
	switch {
	case err == nil:
		return id, false, nil
	case errors.Is(err, pgx.ErrNoRows):

	default:
		return 0, false, fmt.Errorf("lookup user %q: %w", email, err)
	}

	if err := conn.QueryRow(ctx, `
		INSERT INTO users (
			email, first_name, last_name, middle_name, password,
			is_superuser, is_active, date_joined, created_at, updated_at
		) VALUES ($1, $2, '', '', $3, $4, TRUE, NOW(), NOW(), NOW())
		RETURNING id`,
		email, firstName, passwordHash, isSuperuser,
	).Scan(&id); err != nil {
		return 0, false, fmt.Errorf("create user %q: %w", email, err)
	}

	return id, true, nil
}

func ensureAdmin(ctx context.Context, conn *pgxpool.Pool, passwordHash string) (int, error) {
	id, created, err := upsertUser(ctx, conn, adminEmail, adminFirstName, passwordHash, true)
	if err != nil {
		return 0, err
	}

	if created {
		log.Printf("Admin '%s' created (id=%d, password=%s)", adminEmail, id, defaultPassword)
		return id, nil
	}

	tag, err := conn.Exec(ctx, `
		UPDATE users SET is_superuser = TRUE, is_active = TRUE, updated_at = NOW()
		WHERE id = $1 AND (is_superuser = FALSE OR is_active = FALSE)`, id)
	if err != nil {
		return 0, fmt.Errorf("promote admin to superuser: %w", err)
	}

	if tag.RowsAffected() > 0 {
		log.Printf("Admin '%s' promoted to superuser (id=%d)", adminEmail, id)
	} else {
		log.Printf("Admin '%s' already exists (id=%d)", adminEmail, id)
	}

	return id, nil
}

func ensureOwners(ctx context.Context, conn *pgxpool.Pool, count int, passwordHash string) ([]int, error) {
	ids := make([]int, 0, count)
	created := 0

	for i := 1; i <= count; i++ {
		email := fmt.Sprintf(ownerEmailFormat, i)

		id, isNew, err := upsertUser(ctx, conn, email, fmt.Sprintf(ownerNameFormat, i), passwordHash, false)
		if err != nil {
			return nil, err
		}
		if isNew {
			created++
		}

		ids = append(ids, id)
	}

	log.Printf("Owners ready: %d (created %d, existing %d)", len(ids), created, len(ids)-created)
	return ids, nil
}
