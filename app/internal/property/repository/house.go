package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type HouseRepository struct {
	db     *pgxpool.Pool
	awsCfg *configs.AwsConfig
}

func NewHouseRepository(db *pgxpool.Pool, awsCfg *configs.AwsConfig) *HouseRepository {
	return &HouseRepository{db: db, awsCfg: awsCfg}
}

func (r *HouseRepository) GetAll(ctx context.Context) ([]schema.HouseListItem, error) {
	houses, _, err := r.GetAllPaginated(ctx, schema.HouseFilter{}, 0, 0)
	return houses, err
}

func (r *HouseRepository) GetAllPaginated(ctx context.Context, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error) {
	return r.queryHousesPaginated(ctx, filter, limit, offset)
}

func (r *HouseRepository) GetByOwnerPaginated(ctx context.Context, ownerID, limit, offset int) ([]schema.HouseListItem, int, error) {
	f := schema.HouseFilter{}
	f.OwnerID = &ownerID
	return r.queryHousesPaginated(ctx, f, limit, offset)
}

func (r *HouseRepository) queryHousesPaginated(ctx context.Context, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error) {
	baseURL := r.awsCfg.BaseURL()

	wb := newWhereBuilder()
	if filter.OwnerID != nil {
		wb.add("h.owner_id", "=", *filter.OwnerID)
	}
	if filter.Name != nil {
		wb.addILike("h.name_en", *filter.Name)
	}
	if filter.MinPrice != nil {
		wb.add("h.price", ">=", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		wb.add("h.price", "<=", *filter.MaxPrice)
	}
	if filter.GuestCount != nil {
		wb.add("h.guest_qty", ">=", *filter.GuestCount)
	}
	if filter.RoomsQty != nil {
		wb.add("h.rooms_qty", ">=", *filter.RoomsQty)
	}
	if filter.BedroomQty != nil {
		wb.add("h.bedroom_qty", ">=", *filter.BedroomQty)
	}
	if filter.BedQty != nil {
		wb.add("h.bedroom_qty", ">=", *filter.BedQty)
	}
	if filter.BathQty != nil {
		wb.add("h.bath_qty", ">=", *filter.BathQty)
	}
	if filter.GuestsWithPets != nil && *filter.GuestsWithPets {
		wb.add("h.guests_with_pets", "=", true)
	}
	if filter.TypeID != nil {
		wb.add("h.type_id", "=", *filter.TypeID)
	}
	if filter.CountryID != nil {
		wb.add("h.country_id", "=", *filter.CountryID)
	}
	if filter.CityID != nil {
		wb.add("h.city_id", "=", *filter.CityID)
	}
	if filter.CategoryID != nil {
		wb.add("", "", *filter.CategoryID)
	}

	g, gctx := errgroup.WithContext(ctx)

	var total int
	g.Go(func() error {
		countWhere, countArgs := wb.build(0)
		countCatJoin := ""
		if filter.CategoryID != nil {
			countCatJoin = fmt.Sprintf("INNER JOIN house_category hc ON hc.house_id = h.id AND hc.category_id = $%d", len(countArgs))
		}
		q := fmt.Sprintf(`SELECT COUNT(*) FROM houses h %s %s`, countCatJoin, countWhere)
		return r.db.QueryRow(gctx, q, countArgs...).Scan(&total)
	})

	var houses []schema.HouseListItem
	g.Go(func() error {
		selectWhere, selectArgs := wb.build(1)
		allArgs := append([]interface{}{baseURL}, selectArgs...)

		selectCatJoin := ""
		if filter.CategoryID != nil {
			selectCatJoin = fmt.Sprintf("INNER JOIN house_category hc ON hc.house_id = h.id AND hc.category_id = $%d", len(allArgs))
		}

		pagination := ""
		if limit > 0 {
			pagination = fmt.Sprintf("LIMIT $%d OFFSET $%d", len(allArgs)+1, len(allArgs)+2)
			allArgs = append(allArgs, limit, offset)
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
			FROM houses h
			%s
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
					FROM images WHERE house_id = h.id ORDER BY id LIMIT 5
				) i
			) img ON true
			%s
			ORDER BY h.id DESC
			%s`, selectCatJoin, selectWhere, pagination)

		rows, err := r.db.Query(gctx, query, allArgs...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var h schema.HouseListItem
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

type whereBuilder struct {
	conditions []string
	values     []interface{}
}

func newWhereBuilder() *whereBuilder {
	return &whereBuilder{}
}

func (w *whereBuilder) add(col, op string, val interface{}) {
	w.values = append(w.values, val)
	if col != "" {
		w.conditions = append(w.conditions, fmt.Sprintf("%s %s {%d}", col, op, len(w.values)-1))
	}
}

func (w *whereBuilder) addILike(col string, val string) {
	pattern := "%" + val + "%"
	w.values = append(w.values, pattern)
	idx := len(w.values) - 1
	w.conditions = append(w.conditions, fmt.Sprintf(
		"(h.name_en ILIKE {%d} OR h.name_kz ILIKE {%d} OR h.name_ru ILIKE {%d})",
		idx, idx, idx,
	))
}

func (w *whereBuilder) build(argOffset int) (string, []interface{}) {
	if len(w.conditions) == 0 && len(w.values) == len(w.conditions) {
		return "", w.values
	}

	result := ""
	if len(w.conditions) > 0 {
		parts := make([]string, len(w.conditions))
		for i, cond := range w.conditions {
			parts[i] = cond
		}
		for i := range parts {
			for j := range w.values {
				placeholder := fmt.Sprintf("{%d}", j)
				replacement := fmt.Sprintf("$%d", j+1+argOffset)
				parts[i] = strings.ReplaceAll(parts[i], placeholder, replacement)
			}
		}
		result = "WHERE " + strings.Join(parts, " AND ")
	}

	return result, w.values
}

func (r *HouseRepository) GetBySlug(ctx context.Context, slug string) (schema.HouseDetailResponse, error) {
	var h schema.HouseDetailResponse
	var imagesJSON []byte
	var createdAt, updatedAt time.Time
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
			CONCAT(c.name_kz, ', ', ct.name_kz),
			CONCAT(c.name_ru, ', ', ct.name_ru),
			CONCAT(c.name_en, ', ', ct.name_en),
			CONCAT(u.first_name, ' ', u.last_name),
			COALESCE(img.images, '[]')
		FROM houses h
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
			FROM images i
			WHERE i.house_id = h.id
		) img ON true
		WHERE h.slug=$1
	`

	err := r.db.QueryRow(ctx, query, slug, baseURL).Scan(
		&h.ID, &h.NameEN, &h.NameKZ, &h.NameRU, &h.Slug,
		&h.Price, &h.RoomsQty, &h.GuestQty, &h.BedroomQty, &h.BathQty,
		&h.DescriptionEN, &h.DescriptionKZ, &h.DescriptionRU,
		&h.AddressEN, &h.AddressKZ, &h.AddressRU,
		&h.Lng, &h.Lat, &h.IsActive, &h.Priority,
		&h.CommentsRU, &h.CommentsEN, &h.CommentsKZ,
		&h.OwnerID, &h.TypeID, &h.CityID, &h.CountryID,
		&h.GuestsWithPets, &h.BestHouse, &h.Promotion,
		&h.DistrictEN, &h.DistrictKZ, &h.DistrictRU,
		&h.PhoneNumber, &createdAt, &updatedAt,
		&h.LikeCount,
		&h.CountryCityNameKZ, &h.CountryCityNameRU, &h.CountryCityNameEN,
		&h.OwnerFullName,
		&imagesJSON,
	)
	if err != nil {
		return h, err
	}

	h.CreatedAt = createdAt.Format(time.RFC3339)
	h.UpdatedAt = updatedAt.Format(time.RFC3339)

	if len(imagesJSON) == 0 {
		imagesJSON = []byte("[]")
	}
	err = json.Unmarshal(imagesJSON, &h.Images)
	return h, err
}

func (r *HouseRepository) RecordView(ctx context.Context, slug string, userID *int, ip string) {
	var houseID int
	err := r.db.QueryRow(ctx, `SELECT id FROM houses WHERE slug=$1`, slug).Scan(&houseID)
	if err != nil {
		return
	}
	_, _ = r.db.Exec(ctx,
		`INSERT INTO house_views (house_id, user_id, ip_address) VALUES ($1, $2, $3)`,
		houseID, userID, ip,
	)
	_, _ = r.db.Exec(ctx, `UPDATE houses SET view_count = view_count + 1 WHERE id = $1`, houseID)
}

func (r *HouseRepository) Create(ctx context.Context, h schema.HouseCreateRequest) (model.House, error) {
	var house model.House

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
		h.Lng.Float64Ptr(), h.Lat.Float64Ptr(), true, h.Priority.Int(), h.OwnerID, h.TypeID.Int(), h.CityID.IntPtr(), h.CountryID.IntPtr(),
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

func (r *HouseRepository) getForUpdate(ctx context.Context, slug string) (model.House, error) {
	var house model.House
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

func (r *HouseRepository) Update(ctx context.Context, slug string, h schema.HouseUpdateRequest) (model.House, error) {
	house, err := r.getForUpdate(ctx, slug)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return model.House{}, fmt.Errorf("house with slug '%s' not found", slug)
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
