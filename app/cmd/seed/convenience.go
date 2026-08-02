package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/pkg/utils"
)

var defaultConveniences = []string{
	"Wi-Fi", "Air conditioning", "Heating", "Kitchen", "Washer",
	"Dryer", "Free parking", "Pool", "Hot tub", "Gym",
	"TV", "Workspace", "Elevator", "Breakfast", "BBQ grill",
	"Fireplace", "Balcony", "Garden", "Sea view", "Mountain view",
	"Pet friendly", "Smoke alarm", "First aid kit", "Security cameras", "Self check-in",
	"EV charger", "Crib", "Iron", "Hair dryer", "24/7 security",
}

func ensureConveniences(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultConveniences))

	for _, name := range defaultConveniences {
		id, err := ensureRow(ctx, conn, "convenience", name,
			`SELECT id FROM conveniences WHERE name = $1 LIMIT 1`,
			`INSERT INTO conveniences (name, slug, is_active, created_at, updated_at)
			 VALUES ($1, $2, TRUE, NOW(), NOW())
			 RETURNING id`,
			name, utils.GenerateSlug("", name, "", ""),
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
