package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/property/model"
)

type HouseImageRepository struct {
	db *pgxpool.Pool
}

func NewHouseImageRepository(db *pgxpool.Pool) *HouseImageRepository {
	return &HouseImageRepository{db: db}
}

func (r *HouseImageRepository) GetHouseIDBySlug(ctx context.Context, slug string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM houses WHERE slug=$1`, slug).Scan(&id)
	return id, err
}

func (r *HouseImageRepository) CountByHouse(ctx context.Context, houseID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM images WHERE house_id=$1`,
		houseID,
	).Scan(&count)
	return count, err
}

func (r *HouseImageRepository) Create(ctx context.Context, img *model.Image) error {
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

func (r *HouseImageRepository) CreateBatch(ctx context.Context, images []model.Image) error {
	if len(images) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("INSERT INTO images (original, thumbnail, mimetype, width, height, size, house_id) VALUES ")

	args := make([]any, 0, len(images)*7)
	for i, img := range images {
		if i > 0 {
			b.WriteString(",")
		}
		base := i * 7
		b.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			base+1, base+2, base+3, base+4, base+5, base+6, base+7))
		args = append(args, img.Original, img.Thumbnail, img.MimeType, img.Width, img.Height, img.Size, img.HouseID)
	}

	_, err := r.db.Exec(ctx, b.String(), args...)
	return err
}

type ImageKeys struct {
	Original  string
	Thumbnail string
}

func (r *HouseImageRepository) DeleteReturningKeys(ctx context.Context, imageID int) (*ImageKeys, error) {
	var keys ImageKeys
	err := r.db.QueryRow(ctx,
		`DELETE FROM images WHERE id=$1 RETURNING COALESCE(original,''), COALESCE(thumbnail,'')`,
		imageID,
	).Scan(&keys.Original, &keys.Thumbnail)
	if err != nil {
		return nil, err
	}
	return &keys, nil
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
