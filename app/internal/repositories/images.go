package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/models"
)

type HouseImageRepository struct {
	db *pgxpool.Pool
}

func NewHouseImageRepository(db *pgxpool.Pool) *HouseImageRepository {
	return &HouseImageRepository{db: db}
}

func (r *HouseImageRepository) CountByHouse(ctx context.Context, houseID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM images WHERE house_id=$1`,
		houseID,
	).Scan(&count)

	return count, err
}

func (r *HouseImageRepository) Create(ctx context.Context, img *models.Image) error {
	query := `
	INSERT INTO images
	(original, mimetype, size, house_id)
	VALUES ($1,$2,$3,$4)
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		img.Original,
		img.MimeType,
		img.Size,
		img.HouseID,
	).Scan(&img.ID, &img.CreatedAt, &img.UpdatedAt)
}

func (r *HouseImageRepository) Delete(ctx context.Context, imageID int) (*string, error) {
	var key *string

	err := r.db.QueryRow(ctx,
		`SELECT original FROM images WHERE id=$1`,
		imageID,
	).Scan(&key)
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(ctx, `DELETE FROM images WHERE id=$1`, imageID)
	return key, err
}
