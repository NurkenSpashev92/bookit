package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type HouseRepository struct {
	db     *pgxpool.Pool
	awsCfg *configs.AwsConfig
}

func NewHouseRepository(db *pgxpool.Pool, awsCfg *configs.AwsConfig) *HouseRepository {
	return &HouseRepository{db: db, awsCfg: awsCfg}
}

func (r *HouseRepository) GetAll(ctx context.Context) ([]schemas.HouseListItem, error) {
	baseURL := r.awsCfg.BaseURL()

	query := `
		SELECT
			h.id,
			h.name_en,
			h.name_kz,
			h.name_ru,
			h.slug,
			h.price,
			h.address_en,
			h.address_kz,
			h.address_ru,
			h.priority,
			h.guests_with_pets,
			h.best_house,
			h.promotion,
			CONCAT(c.name_kz, ', ', ct.name_kz) AS country_name_kz,
			CONCAT(c.name_ru, ', ', ct.name_ru) AS country_name_ru,
			CONCAT(c.name_en, ', ', ct.name_en) AS country_name_en,
			CONCAT(u.first_name, ' ', u.last_name) AS full_name,
			h.like_count,
			COALESCE(img.images, '[]') as images
		FROM houses h
		LEFT JOIN countries c ON c.id = h.country_id
		LEFT JOIN cities ct ON ct.id = h.city_id
		LEFT JOIN users u ON u.id = h.owner_id
		LEFT JOIN LATERAL (
			SELECT COALESCE(json_agg(
				json_build_object(
					'id', i.id,
					'original', $1 || i.original,
					'thumbnail', CASE WHEN i.thumbnail IS NOT NULL AND i.thumbnail <> '' THEN $1 || i.thumbnail ELSE '' END,
					'mime_type', i.mimetype,
					'size', i.size,
					'house_id', i.house_id
				)
			) FILTER (WHERE i.id IS NOT NULL), '[]') as images
			FROM (
				SELECT id, original, thumbnail, mimetype, size, house_id
				FROM images
				WHERE house_id = h.id
				ORDER BY id
				LIMIT 5
			) i
		) img ON true
		ORDER BY h.id DESC
	`

	rows, err := r.db.Query(ctx, query, baseURL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []schemas.HouseListItem

	for rows.Next() {
		var h schemas.HouseListItem
		var imagesJSON []byte

		err := rows.Scan(
			&h.ID,
			&h.NameEN, &h.NameKZ, &h.NameRU, &h.Slug,
			&h.Price,
			&h.AddressEN, &h.AddressKZ, &h.AddressRU,
			&h.Priority,
			&h.GuestsWithPets,
			&h.BestHouse,
			&h.Promotion,
			&h.CountryCityNameKZ, &h.CountryCityNameRU, &h.CountryCityNameEN,
			&h.OwnerFullName,
			&h.LikeCount,
			&imagesJSON,
		)

		if err != nil {
			return nil, err
		}

		if len(imagesJSON) == 0 {
			imagesJSON = []byte("[]")
		}

		err = json.Unmarshal(imagesJSON, &h.Images)
		if err != nil {
			return nil, err
		}

		houses = append(houses, h)
	}

	return houses, nil
}

func (r *HouseRepository) GetBySlug(ctx context.Context, slug string) (models.House, error) {
	var house models.House
	var imagesJSON []byte
	baseURL := r.awsCfg.BaseURL()

	query := `
		SELECT
			h.id, h.name_en, h.name_kz, h.name_ru, h.slug, h.price, h.rooms_qty, h.guest_qty, h.bedroom_qty, h.bath_qty,
			h.description_en, h.description_kz, h.description_ru,
			h.address_en, h.address_kz, h.address_ru,
			h.lng, h.lat, h.is_active, h.priority,
			h.comments_ru, h.comments_en, h.comments_kz,
			h.owner_id, h.type_id, h.city_id, h.country_id, h.guests_with_pets, h.best_house, h.promotion,
			h.district_en, h.district_kz, h.district_ru, h.phone_number, h.created_at, h.updated_at,
			h.like_count,
			COALESCE(img.images, '[]') as images
		FROM houses h
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
			FROM images i
			WHERE i.house_id = h.id
		) img ON true
		WHERE h.slug=$1
	`

	err := r.db.QueryRow(ctx, query, slug, baseURL).Scan(
		&house.ID, &house.NameEN, &house.NameKZ, &house.NameRU, &house.Slug,
		&house.Price, &house.RoomsQty, &house.GuestQty, &house.BedroomQty, &house.BathQty,
		&house.DescriptionEN, &house.DescriptionKZ, &house.DescriptionRU,
		&house.AddressEN, &house.AddressKZ, &house.AddressRU,
		&house.Lng, &house.Lat, &house.IsActive, &house.Priority,
		&house.CommentsRU, &house.CommentsEN, &house.CommentsKZ,
		&house.OwnerID, &house.TypeID, &house.CityID, &house.CountryID,
		&house.GuestsWithPets, &house.BestHouse, &house.Promotion,
		&house.DistrictEN, &house.DistrictKZ, &house.DistrictRU,
		&house.PhoneNumber, &house.CreatedAt, &house.UpdatedAt,
		&house.LikeCount,
		&imagesJSON,
	)
	if err != nil {
		return house, err
	}

	if len(imagesJSON) == 0 {
		imagesJSON = []byte("[]")
	}

	err = json.Unmarshal(imagesJSON, &house.Images)
	return house, err
}

func (r *HouseRepository) Create(ctx context.Context, h schemas.HouseCreateRequest) (models.House, error) {
	var house models.House

	slugValue := utils.GenerateSlug(h.Slug, h.NameEN, h.NameKZ, h.NameRU)

	exists, err := r.SlugExists(ctx, slugValue)
	if err != nil {
		return house, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return house, fmt.Errorf("slug '%s' already exists", slugValue)
	}

	query := `INSERT INTO houses (
			name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru, address_en, address_kz, address_ru,
			lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
			guests_with_pets, best_house, promotion, district_en, district_kz, district_ru, phone_number
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
		RETURNING id, name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru, address_en, address_kz, address_ru,
			lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
			guests_with_pets, best_house, promotion, district_en, district_kz, district_ru, phone_number,
			created_at, updated_at
	`
	if err := r.db.QueryRow(ctx,
		query,
		h.NameEN, h.NameKZ, h.NameRU, slugValue, h.Price.Int(), h.RoomsQty.Int(), h.GuestQty.Int(), h.BedroomQty.Int(), h.BathQty.IntPtr(),
		h.DescriptionEN, h.DescriptionKZ, h.DescriptionRU, h.AddressEN, h.AddressKZ, h.AddressRU,
		h.Lng.Float64Ptr(), h.Lat.Float64Ptr(), h.IsActive, h.Priority.Int(), h.OwnerID, h.TypeID.Int(), h.CityID.IntPtr(), h.CountryID.IntPtr(),
		h.GuestsWithPets, h.BestHouse, h.Promotion, h.DistrictEN, h.DistrictKZ, h.DistrictRU, h.PhoneNumber,
	).Scan(
		&house.ID, &house.NameEN, &house.NameKZ, &house.NameRU, &house.Slug, &house.Price, &house.RoomsQty, &house.GuestQty,
		&house.BedroomQty, &house.BathQty, &house.DescriptionEN, &house.DescriptionKZ, &house.DescriptionRU,
		&house.AddressEN, &house.AddressKZ, &house.AddressRU, &house.Lng, &house.Lat,
		&house.IsActive, &house.Priority, &house.OwnerID, &house.TypeID, &house.CityID, &house.CountryID,
		&house.GuestsWithPets, &house.BestHouse, &house.Promotion, &house.DistrictEN, &house.DistrictKZ, &house.DistrictRU, &house.PhoneNumber,
		&house.CreatedAt, &house.UpdatedAt,
	); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "houses_type_id_fkey":
				return house, fmt.Errorf("type_id %d does not exist", h.TypeID.Int())
			case "houses_city_id_fkey":
				return house, fmt.Errorf("city_id does not exist")
			case "houses_country_id_fkey":
				return house, fmt.Errorf("country_id does not exist")
			case "houses_owner_id_fkey":
				return house, fmt.Errorf("owner_id %d does not exist", h.OwnerID)
			}
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return house, fmt.Errorf("slug already exists")
		}
		return house, err
	}

	return house, err
}

func (r *HouseRepository) getForUpdate(ctx context.Context, slug string) (models.House, error) {
	var house models.House
	query := `
		SELECT id, name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru,
			address_en, address_kz, address_ru,
			lng, lat, is_active, priority,
			comments_ru, comments_en, comments_kz,
			owner_id, type_id, city_id, country_id, guests_with_pets, best_house, promotion,
			district_en, district_kz, district_ru, phone_number, created_at, updated_at
		FROM houses WHERE slug=$1
	`
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&house.ID, &house.NameEN, &house.NameKZ, &house.NameRU, &house.Slug,
		&house.Price, &house.RoomsQty, &house.GuestQty, &house.BedroomQty, &house.BathQty,
		&house.DescriptionEN, &house.DescriptionKZ, &house.DescriptionRU,
		&house.AddressEN, &house.AddressKZ, &house.AddressRU,
		&house.Lng, &house.Lat, &house.IsActive, &house.Priority,
		&house.CommentsRU, &house.CommentsEN, &house.CommentsKZ,
		&house.OwnerID, &house.TypeID, &house.CityID, &house.CountryID,
		&house.GuestsWithPets, &house.BestHouse, &house.Promotion,
		&house.DistrictEN, &house.DistrictKZ, &house.DistrictRU,
		&house.PhoneNumber, &house.CreatedAt, &house.UpdatedAt,
	)
	return house, err
}

func (r *HouseRepository) Update(ctx context.Context, slug string, h schemas.HouseUpdateRequest) (models.House, error) {
	house, err := r.getForUpdate(ctx, slug)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.House{}, fmt.Errorf("house with slug '%s' not found", slug)
		}
		return house, err
	}

	if h.NameEN != nil {
		house.NameEN = *h.NameEN
	}
	if h.NameKZ != nil {
		house.NameKZ = *h.NameKZ
	}
	if h.NameRU != nil {
		house.NameRU = *h.NameRU
	}
	if h.Slug != nil {
		house.Slug = *h.Slug
	}
	if h.Price != nil {
		house.Price = h.Price.Int()
	}
	if h.RoomsQty != nil {
		house.RoomsQty = h.RoomsQty.Int()
	}
	if h.GuestQty != nil {
		house.GuestQty = h.GuestQty.Int()
	}
	if h.BedroomQty != nil {
		house.BedroomQty = h.BedroomQty.Int()
	}
	if h.BathQty != nil {
		house.BathQty = h.BathQty.IntPtr()
	}
	if h.DescriptionEN != nil {
		house.DescriptionEN = *h.DescriptionEN
	}
	if h.DescriptionKZ != nil {
		house.DescriptionKZ = *h.DescriptionKZ
	}
	if h.DescriptionRU != nil {
		house.DescriptionRU = *h.DescriptionRU
	}
	if h.AddressEN != nil {
		house.AddressEN = *h.AddressEN
	}
	if h.AddressKZ != nil {
		house.AddressKZ = *h.AddressKZ
	}
	if h.AddressRU != nil {
		house.AddressRU = *h.AddressRU
	}
	if h.Lng != nil {
		house.Lng = h.Lng.Float64Ptr()
	}
	if h.Lat != nil {
		house.Lat = h.Lat.Float64Ptr()
	}
	if h.IsActive != nil {
		house.IsActive = *h.IsActive
	}
	if h.Priority != nil {
		house.Priority = h.Priority.Int()
	}
	if h.TypeID != nil {
		house.TypeID = h.TypeID.Int()
	}
	if h.CityID != nil {
		house.CityID = h.CityID.IntPtr()
	}
	if h.CountryID != nil {
		house.CountryID = h.CountryID.IntPtr()
	}
	if h.GuestsWithPets != nil {
		house.GuestsWithPets = *h.GuestsWithPets
	}
	if h.BestHouse != nil {
		house.BestHouse = *h.BestHouse
	}
	if h.Promotion != nil {
		house.Promotion = *h.Promotion
	}
	if h.DistrictEN != nil {
		house.DistrictEN = *h.DistrictEN
	}
	if h.DistrictKZ != nil {
		house.DistrictKZ = *h.DistrictKZ
	}
	if h.DistrictRU != nil {
		house.DistrictRU = *h.DistrictRU
	}
	if h.PhoneNumber != nil {
		house.PhoneNumber = *h.PhoneNumber
	}

	house.UpdatedAt = time.Now()

	slugValue := utils.GenerateSlug(house.Slug, house.NameEN, house.NameKZ, house.NameRU)

	exists, err := r.SlugExistsExceptID(ctx, slugValue, house.ID)
	if err != nil {
		return house, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return house, fmt.Errorf("slug '%s' already exists", slugValue)
	}

	query := `
		UPDATE houses SET
			name_en=$1, name_kz=$2, name_ru=$3, slug=$4, price=$5, rooms_qty=$6, guest_qty=$7, bedroom_qty=$8, bath_qty=$9,
			description_en=$10, description_kz=$11, description_ru=$12,
			address_en=$13, address_kz=$14, address_ru=$15,
			lng=$16, lat=$17, is_active=$18, type_id=$19, city_id=$20, country_id=$21, guests_with_pets=$22, best_house=$23,
			promotion=$24, district_en=$25, district_kz=$26, district_ru=$27, phone_number=$28, updated_at=$29
		WHERE id=$30
	`

	_, err = r.db.Exec(ctx, query,
		house.NameEN, house.NameKZ, house.NameRU, slugValue, house.Price,
		house.RoomsQty, house.GuestQty, house.BedroomQty, house.BathQty,
		house.DescriptionEN, house.DescriptionKZ, house.DescriptionRU,
		house.AddressEN, house.AddressKZ, house.AddressRU,
		house.Lng, house.Lat, house.IsActive, house.TypeID,
		house.CityID, house.CountryID, house.GuestsWithPets, house.BestHouse,
		house.Promotion, house.DistrictEN, house.DistrictKZ, house.DistrictRU,
		house.PhoneNumber, house.UpdatedAt, house.ID,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "houses_type_id_fkey":
				return house, fmt.Errorf("type_id does not exist")
			case "houses_city_id_fkey":
				return house, fmt.Errorf("city_id does not exist")
			case "houses_country_id_fkey":
				return house, fmt.Errorf("country_id does not exist")
			}
		}
		return house, err
	}

	return house, nil
}

func (r *HouseRepository) Delete(ctx context.Context, slug string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM houses WHERE slug=$1", slug)
	return err
}

func (r *HouseRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM houses WHERE slug=$1)`,
		slug,
	).Scan(&exists)

	return exists, err
}

func (r *HouseRepository) SlugExistsExceptID(ctx context.Context, slug string, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM houses WHERE slug=$1 AND id<>$2)", slug, id).Scan(&exists)
	return exists, err
}
