package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	minCategoriesPerHouse = 1
	maxCategoriesPerHouse = 3
)

const linkHouseCategorySQL = `
	INSERT INTO house_category (house_id, category_id)
	VALUES ($1, $2)
	ON CONFLICT (house_id, category_id) DO NOTHING`

func linkHouseCategories(ctx context.Context, conn *pgxpool.Pool, houseIDs, categoryIDs []int) error {
	if len(houseIDs) == 0 || len(categoryIDs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	queued := 0

	for _, houseID := range houseIDs {
		for _, categoryID := range pickCategories(categoryIDs) {
			batch.Queue(linkHouseCategorySQL, houseID, categoryID)
			queued++
		}
	}

	if err := execBatch(ctx, conn, batch, queued, "link house category"); err != nil {
		return err
	}

	log.Printf("House-category links written: %d for %d houses", queued, len(houseIDs))
	return nil
}

func pickCategories(categoryIDs []int) []int {
	count := minCategoriesPerHouse + rand.Intn(maxCategoriesPerHouse-minCategoriesPerHouse+1)
	if count > len(categoryIDs) {
		count = len(categoryIDs)
	}

	shuffled := make([]int, len(categoryIDs))
	copy(shuffled, categoryIDs)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:count]
}

func backfillHouseLocations(ctx context.Context, conn *pgxpool.Pool, cities []cityRef) error {
	if len(cities) == 0 {
		return nil
	}

	houseIDs, err := selectHouseIDs(ctx, conn, `SELECT id FROM houses WHERE city_id IS NULL`)
	if err != nil {
		return fmt.Errorf("select houses without city: %w", err)
	}
	if len(houseIDs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, houseID := range houseIDs {
		city := cities[rand.Intn(len(cities))]
		batch.Queue(
			`UPDATE houses SET city_id = $1, country_id = $2, updated_at = NOW() WHERE id = $3`,
			city.id, city.countryID, houseID,
		)
	}

	if err := execBatch(ctx, conn, batch, len(houseIDs), "backfill house location"); err != nil {
		return err
	}

	log.Printf("Houses backfilled with a city: %d", len(houseIDs))
	return nil
}

func backfillHouseCategories(ctx context.Context, conn *pgxpool.Pool, categoryIDs []int) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	houseIDs, err := selectHouseIDs(ctx, conn, `
		SELECT h.id
		FROM houses h
		LEFT JOIN house_category hc ON hc.house_id = h.id
		WHERE hc.house_id IS NULL`)
	if err != nil {
		return fmt.Errorf("select houses without category: %w", err)
	}

	return linkHouseCategories(ctx, conn, houseIDs, categoryIDs)
}

func selectHouseIDs(ctx context.Context, conn *pgxpool.Pool, query string) ([]int, error) {
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func execBatch(ctx context.Context, conn *pgxpool.Pool, batch *pgx.Batch, queued int, action string) error {
	if queued == 0 {
		return nil
	}

	results := conn.SendBatch(ctx, batch)

	var execErr error
	for range queued {
		if _, err := results.Exec(); err != nil && execErr == nil {
			execErr = fmt.Errorf("%s: %w", action, err)
		}
	}

	closeErr := results.Close()
	if execErr != nil {
		return execErr
	}
	if closeErr != nil {
		return fmt.Errorf("close %s batch: %w", action, closeErr)
	}

	return nil
}
