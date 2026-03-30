package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/schemas"
)

type HouseLikeRepository struct {
	db *pgxpool.Pool
}

func NewHouseLikeRepository(db *pgxpool.Pool) *HouseLikeRepository {
	return &HouseLikeRepository{db: db}
}

func (r *HouseLikeRepository) LikeReturningCount(ctx context.Context, userID, houseID int) (int, error) {
	var count int
	query := `
		WITH ins AS (
			INSERT INTO house_likes (user_id, house_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			RETURNING house_id
		)
		SELECT like_count FROM houses WHERE id = $2
	`
	err := r.db.QueryRow(ctx, query, userID, houseID).Scan(&count)
	return count, err
}

func (r *HouseLikeRepository) UnlikeReturningCount(ctx context.Context, userID, houseID int) (int, error) {
	var count int
	query := `
		WITH del AS (
			DELETE FROM house_likes WHERE user_id=$1 AND house_id=$2
			RETURNING house_id
		)
		SELECT like_count FROM houses WHERE id = $2
	`
	err := r.db.QueryRow(ctx, query, userID, houseID).Scan(&count)
	return count, err
}

func (r *HouseLikeRepository) StatusWithCount(ctx context.Context, userID, houseID int) (bool, int, error) {
	var liked bool
	var count int
	query := `
		SELECT
			EXISTS(SELECT 1 FROM house_likes WHERE user_id=$1 AND house_id=$2),
			(SELECT like_count FROM houses WHERE id=$2)
	`
	err := r.db.QueryRow(ctx, query, userID, houseID).Scan(&liked, &count)
	return liked, count, err
}

func (r *HouseLikeRepository) CountByHouse(ctx context.Context, houseID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT like_count FROM houses WHERE id=$1`,
		houseID,
	).Scan(&count)
	return count, err
}

func (r *HouseLikeRepository) GetUserLikedHouseIDs(ctx context.Context, userID int) ([]int, error) {
	rows, err := r.db.Query(ctx,
		`SELECT house_id FROM house_likes WHERE user_id=$1 ORDER BY created_at DESC`,
		userID,
	)
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
	return ids, nil
}

func (r *HouseLikeRepository) GetUserLikedHouses(ctx context.Context, userID int) ([]schemas.HouseLikeItem, error) {
	query := `
		SELECT
			h.id, h.name_en, h.name_kz, h.name_ru, h.slug, h.price,
			h.address_en, h.address_kz, h.address_ru,
			hl.created_at AS liked_at
		FROM house_likes hl
		INNER JOIN houses h ON h.id = hl.house_id
		WHERE hl.user_id = $1
		ORDER BY hl.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []schemas.HouseLikeItem
	for rows.Next() {
		var item schemas.HouseLikeItem
		if err := rows.Scan(
			&item.ID, &item.NameEN, &item.NameKZ, &item.NameRU,
			&item.Slug, &item.Price,
			&item.AddressEN, &item.AddressKZ, &item.AddressRU,
			&item.LikedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
