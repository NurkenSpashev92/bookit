package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/configs"
	propertyschema "github.com/nurkenspashev92/bookit/internal/property/schema"
)

type HouseLikeRepository struct {
	db     *pgxpool.Pool
	awsCfg *configs.AwsConfig
}

func NewHouseLikeRepository(db *pgxpool.Pool, awsCfg *configs.AwsConfig) *HouseLikeRepository {
	return &HouseLikeRepository{db: db, awsCfg: awsCfg}
}

func (r *HouseLikeRepository) getHouseIDBySlug(ctx context.Context, slug string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM houses WHERE slug=$1`, slug).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, errors.New("house not found")
	}
	return id, err
}

func (r *HouseLikeRepository) LikeReturningCount(ctx context.Context, userID int, slug string) (int, error) {
	houseID, err := r.getHouseIDBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}

	var count int
	query := `
		WITH ins AS (
			INSERT INTO house_likes (user_id, house_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			RETURNING house_id
		)
		SELECT like_count FROM houses WHERE id = $2
	`
	err = r.db.QueryRow(ctx, query, userID, houseID).Scan(&count)
	return count, err
}

func (r *HouseLikeRepository) UnlikeReturningCount(ctx context.Context, userID int, slug string) (int, error) {
	houseID, err := r.getHouseIDBySlug(ctx, slug)
	if err != nil {
		return 0, err
	}

	var count int
	query := `
		WITH del AS (
			DELETE FROM house_likes WHERE user_id=$1 AND house_id=$2
			RETURNING house_id
		)
		SELECT like_count FROM houses WHERE id = $2
	`
	err = r.db.QueryRow(ctx, query, userID, houseID).Scan(&count)
	return count, err
}

func (r *HouseLikeRepository) StatusWithCount(ctx context.Context, userID int, slug string) (bool, int, error) {
	houseID, err := r.getHouseIDBySlug(ctx, slug)
	if err != nil {
		return false, 0, err
	}

	var liked bool
	var count int
	query := `
		SELECT
			EXISTS(SELECT 1 FROM house_likes WHERE user_id=$1 AND house_id=$2),
			(SELECT like_count FROM houses WHERE id=$2)
	`
	err = r.db.QueryRow(ctx, query, userID, houseID).Scan(&liked, &count)
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

func (r *HouseLikeRepository) GetUserLikedHouses(ctx context.Context, userID int) ([]propertyschema.HouseListItem, error) {
	result, _, err := r.GetUserLikedHousesPaginated(ctx, userID, 0, 0)
	return result, err
}

func (r *HouseLikeRepository) GetUserLikedHousesPaginated(ctx context.Context, userID, limit, offset int) ([]propertyschema.HouseListItem, int, error) {
	baseURL := r.awsCfg.BaseURL()

	g, gctx := errgroup.WithContext(ctx)

	var total int
	g.Go(func() error {
		return r.db.QueryRow(gctx,
			`SELECT COUNT(*) FROM house_likes WHERE user_id=$1`, userID,
		).Scan(&total)
	})

	var houses []propertyschema.HouseListItem
	g.Go(func() error {
		pagination := ""
		args := []interface{}{userID, baseURL}
		if limit > 0 {
			pagination = fmt.Sprintf("LIMIT $3 OFFSET $4")
			args = append(args, limit, offset)
		}

		query := fmt.Sprintf(`
			SELECT
				h.id, h.name_en, h.name_kz, h.name_ru, h.slug, h.price,
				h.address_en, h.address_kz, h.address_ru,
				h.priority, h.guests_with_pets, h.best_house, h.promotion,
				CONCAT(c.name_kz, ', ', ct.name_kz),
				CONCAT(c.name_ru, ', ', ct.name_ru),
				CONCAT(c.name_en, ', ', ct.name_en),
				CONCAT(u.first_name, ' ', u.last_name),
				h.like_count,
				COALESCE(img.images, '[]')
			FROM house_likes hl
			INNER JOIN houses h ON h.id = hl.house_id
			LEFT JOIN countries c ON c.id = h.country_id
			LEFT JOIN cities ct ON ct.id = h.city_id
			LEFT JOIN users u ON u.id = h.owner_id
			LEFT JOIN LATERAL (
				SELECT COALESCE(json_agg(
					json_build_object(
						'id', i.id,
						'original', $2 || i.original,
						'thumbnail', CASE WHEN i.thumbnail IS NOT NULL AND i.thumbnail <> '' THEN $2 || i.thumbnail ELSE '' END,
						'mime_type', i.mimetype,
						'size', i.size,
						'house_id', i.house_id
					)
				) FILTER (WHERE i.id IS NOT NULL), '[]') as images
				FROM (
					SELECT id, original, thumbnail, mimetype, size, house_id
					FROM images WHERE house_id = h.id ORDER BY id LIMIT 5
				) i
			) img ON true
			WHERE hl.user_id = $1
			ORDER BY hl.created_at DESC
			%s`, pagination)

		rows, err := r.db.Query(gctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var h propertyschema.HouseListItem
			var imagesJSON []byte

			if err := rows.Scan(
				&h.ID, &h.NameEN, &h.NameKZ, &h.NameRU, &h.Slug, &h.Price,
				&h.AddressEN, &h.AddressKZ, &h.AddressRU,
				&h.Priority, &h.GuestsWithPets, &h.BestHouse, &h.Promotion,
				&h.CountryCityNameKZ, &h.CountryCityNameRU, &h.CountryCityNameEN,
				&h.OwnerFullName, &h.LikeCount, &imagesJSON,
			); err != nil {
				return err
			}

			h.IsLiked = true

			if len(imagesJSON) == 0 {
				imagesJSON = []byte("[]")
			}
			if err := json.Unmarshal(imagesJSON, &h.Images); err != nil {
				return err
			}

			houses = append(houses, h)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return houses, total, nil
}
