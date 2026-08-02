package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/utils"
)

const pgUniqueViolation = "23505"

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

func (r *HouseRepository) OwnerIDBySlug(ctx context.Context, slug string) (int, error) {
	var ownerID int
	err := r.db.QueryRow(ctx, `SELECT owner_id FROM houses WHERE slug=$1`, slug).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrHouseNotFound
		}
		return 0, err
	}
	return ownerID, nil
}

func (r *HouseRepository) queryHousesPaginated(ctx context.Context, filter schema.HouseFilter, limit, offset int) ([]schema.HouseListItem, int, error) {
	baseURL := r.awsCfg.BaseURL()

	wb := newWhereBuilder()
	if filter.OwnerID != nil {
		wb.add("h.owner_id", "=", *filter.OwnerID)
	} else if filter.Moderation {
		if filter.IsActive != nil {
			wb.add("h.is_active", "=", *filter.IsActive)
		}
	} else {
		wb.add("h.is_active", "=", true)
	}
	if filter.Name != nil {
		wb.addILike(*filter.Name)
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
	if filter.GuestsWithBabies != nil && *filter.GuestsWithBabies {
		wb.add("h.guests_with_babies", "=", true)
	}
	if filter.TypeSlug != nil {
		wb.addExpr("h.type_id = (SELECT id FROM types WHERE slug = $%d)", *filter.TypeSlug)
	}
	if filter.CountryID != nil {
		wb.add("h.country_id", "=", *filter.CountryID)
	}
	if filter.CityID != nil {
		wb.add("h.city_id", "=", *filter.CityID)
	}
	if filter.CategorySlug != nil {
		wb.addArg(*filter.CategorySlug)
	}

	g, gctx := errgroup.WithContext(ctx)

	var total int
	g.Go(func() error {
		countWhere, countArgs := wb.build(0)
		countCatJoin := ""
		if filter.CategorySlug != nil {
			countCatJoin = fmt.Sprintf("INNER JOIN house_category hc ON hc.house_id = h.id AND hc.category_id = (SELECT id FROM categories WHERE slug = $%d)", len(countArgs))
		}
		q := fmt.Sprintf(`SELECT COUNT(*) FROM houses h %s %s`, countCatJoin, countWhere)
		return r.db.QueryRow(gctx, q, countArgs...).Scan(&total)
	})

	var houses []schema.HouseListItem
	g.Go(func() error {
		selectWhere, selectArgs := wb.build(1)
		allArgs := append([]interface{}{baseURL}, selectArgs...)

		selectCatJoin := ""
		if filter.CategorySlug != nil {
			selectCatJoin = fmt.Sprintf("INNER JOIN house_category hc ON hc.house_id = h.id AND hc.category_id = (SELECT id FROM categories WHERE slug = $%d)", len(allArgs))
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
				h.best_house, h.promotion,
				h.is_verified, h.is_sale, h.is_newest, h.is_hot, h.is_featured, h.is_discount,
				h.is_active,
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
						'thumbnail', CASE WHEN i.thumbnail IS NOT NULL AND i.thumbnail <> '' THEN $1 || i.thumbnail ELSE '' END,
						'mime_type', i.mimetype,
						'size', i.size,
						'is_label', i.is_label,
						'house_id', i.house_id
					)
				) FILTER (WHERE i.id IS NOT NULL), '[]') as images
				FROM (
					SELECT id, thumbnail, mimetype, size, is_label, house_id
					FROM images WHERE house_id = h.id ORDER BY (is_label IS TRUE) DESC, id LIMIT 5
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
				&h.BestHouse, &h.Promotion,
				&h.IsVerified, &h.IsSale, &h.IsNewest, &h.IsHot, &h.IsFeatured, &h.IsDiscount,
				&h.IsActive,
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

type condition struct {
	format    string
	valueIdxs []int
}

type whereBuilder struct {
	conditions []condition
	values     []interface{}
}

func newWhereBuilder() *whereBuilder {
	return &whereBuilder{}
}

func (w *whereBuilder) add(col, op string, val interface{}) {
	w.values = append(w.values, val)
	w.conditions = append(w.conditions, condition{
		format:    col + " " + op + " $%d",
		valueIdxs: []int{len(w.values) - 1},
	})
}

func (w *whereBuilder) addArg(val interface{}) {
	w.values = append(w.values, val)
}

func (w *whereBuilder) addExpr(format string, val interface{}) {
	w.values = append(w.values, val)
	w.conditions = append(w.conditions, condition{
		format:    format,
		valueIdxs: []int{len(w.values) - 1},
	})
}

func (w *whereBuilder) addILike(val string) {
	w.values = append(w.values, "%"+val+"%")

	idx := len(w.values) - 1
	w.conditions = append(w.conditions, condition{
		format:    "(h.name_en ILIKE $%d OR h.name_kz ILIKE $%d OR h.name_ru ILIKE $%d)",
		valueIdxs: []int{idx, idx, idx},
	})
}

func (w *whereBuilder) build(argOffset int) (string, []interface{}) {
	if len(w.conditions) == 0 {
		return "", w.values
	}

	parts := make([]string, len(w.conditions))
	for i, cond := range w.conditions {
		numbers := make([]interface{}, len(cond.valueIdxs))
		for j, idx := range cond.valueIdxs {
			numbers[j] = idx + 1 + argOffset
		}
		parts[i] = fmt.Sprintf(cond.format, numbers...)
	}

	return "WHERE " + strings.Join(parts, " AND "), w.values
}

func (r *HouseRepository) GetBySlug(ctx context.Context, slug string) (schema.HouseDetailResponse, error) {
	var h schema.HouseDetailResponse
	var imagesJSON []byte
	baseURL := r.awsCfg.BaseURL()

	query := `
		SELECT
			h.id, h.name_en, h.name_kz, h.name_ru, h.slug, h.price, h.rooms_qty, h.guest_qty, h.bedroom_qty, h.bath_qty,
			h.description_en, h.description_kz, h.description_ru,
			h.address_en, h.address_kz, h.address_ru,
			h.lng, h.lat, h.is_active,
			h.comments_ru, h.comments_en, h.comments_kz,
			h.type_id, h.city_id, h.country_id, h.guests_with_pets, h.guests_with_babies, h.best_house, h.promotion,
			h.is_verified, h.is_sale, h.is_newest, h.is_hot, h.is_featured, h.is_discount,
			h.district_en, h.district_kz, h.district_ru, h.phone_number,
			h.like_count,
			CONCAT(u.first_name, ' ', u.last_name),
			CASE WHEN u.payment_qr IS NOT NULL AND u.payment_qr <> '' THEN $2 || u.payment_qr ELSE '' END,
			COALESCE(u.payment_phone, ''),
			COALESCE(img.images, '[]')
		FROM houses h
		LEFT JOIN users u ON u.id = h.owner_id
		LEFT JOIN LATERAL (
			SELECT COALESCE(json_agg(
				json_build_object(
					'id', i.id,
					'original', $2 || i.original,
					'thumbnail', CASE WHEN i.thumbnail IS NOT NULL AND i.thumbnail <> '' THEN $2 || i.thumbnail ELSE '' END,
					'mime_type', i.mimetype,
					'size', i.size,
					'is_label', i.is_label,
					'house_id', i.house_id
				)
				ORDER BY (i.is_label IS TRUE) DESC, i.id
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
		&h.Lng, &h.Lat, &h.IsActive,
		&h.CommentsRU, &h.CommentsEN, &h.CommentsKZ,
		&h.TypeID, &h.CityID, &h.CountryID,
		&h.GuestsWithPets, &h.GuestsWithBabies, &h.BestHouse, &h.Promotion,
		&h.IsVerified, &h.IsSale, &h.IsNewest, &h.IsHot, &h.IsFeatured, &h.IsDiscount,
		&h.DistrictEN, &h.DistrictKZ, &h.DistrictRU,
		&h.PhoneNumber,
		&h.LikeCount,
		&h.OwnerFullName,
		&h.OwnerPaymentQR,
		&h.OwnerPaymentPhone,
		&imagesJSON,
	)
	if err != nil {
		return h, err
	}

	if len(imagesJSON) == 0 {
		imagesJSON = []byte("[]")
	}
	err = json.Unmarshal(imagesJSON, &h.Images)
	return h, err
}

func (r *HouseRepository) RecordView(ctx context.Context, slug string, userID *int, ip string) {
	const query = `
		WITH target AS (
			SELECT id FROM houses WHERE slug = $1
		), logged AS (
			INSERT INTO house_views (house_id, user_id, ip_address)
			SELECT id, $2, $3 FROM target
		)
		UPDATE houses SET view_count = view_count + 1
		WHERE id IN (SELECT id FROM target)`

	if _, err := r.db.Exec(ctx, query, slug, userID, ip); err != nil {
		log.WithContext(ctx).Errorw("failed to record house view",
			"slug", slug,
			"error", err,
		)
	}
}

func (r *HouseRepository) Create(ctx context.Context, h schema.HouseCreateRequest) (model.House, error) {
	var house model.House

	slugValue := utils.GenerateSlug(h.Slug, h.NameEN, h.NameKZ, h.NameRU)

	exists, err := r.SlugExists(ctx, slugValue)
	if err != nil {
		return house, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return house, fmt.Errorf("%w: %s", model.ErrSlugExists, slugValue)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return house, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `INSERT INTO houses (
			name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru, address_en, address_kz, address_ru,
			lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
			guests_with_pets, guests_with_babies, best_house, promotion, district_en, district_kz, district_ru, phone_number,
			is_verified, is_sale, is_hot, is_featured, is_discount
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)
		RETURNING id, name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru, address_en, address_kz, address_ru,
			lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
			guests_with_pets, guests_with_babies, best_house, promotion, district_en, district_kz, district_ru, phone_number,
			is_verified, is_sale, is_newest, is_hot, is_featured, is_discount,
			created_at, updated_at
	`
	if err := tx.QueryRow(ctx,
		query,
		h.NameEN, h.NameKZ, h.NameRU, slugValue, h.Price.Int(), h.RoomsQty.Int(), h.GuestQty.Int(), h.BedroomQty.Int(), h.BathQty.IntPtr(),
		h.DescriptionEN, h.DescriptionKZ, h.DescriptionRU, h.AddressEN, h.AddressKZ, h.AddressRU,
		h.Lng.Float64Ptr(), h.Lat.Float64Ptr(), false, h.Priority.Int(), h.OwnerID, h.TypeID.Int(), h.CityID.IntPtr(), h.CountryID.IntPtr(),
		h.GuestsWithPets, h.GuestsWithBabies, h.BestHouse, h.Promotion, h.DistrictEN, h.DistrictKZ, h.DistrictRU, h.PhoneNumber,
		h.IsVerified, h.IsSale, h.IsHot, h.IsFeatured, h.IsDiscount,
	).Scan(
		&house.ID, &house.NameEN, &house.NameKZ, &house.NameRU, &house.Slug, &house.Price, &house.RoomsQty, &house.GuestQty,
		&house.BedroomQty, &house.BathQty, &house.DescriptionEN, &house.DescriptionKZ, &house.DescriptionRU,
		&house.AddressEN, &house.AddressKZ, &house.AddressRU, &house.Lng, &house.Lat,
		&house.IsActive, &house.Priority, &house.OwnerID, &house.TypeID, &house.CityID, &house.CountryID,
		&house.GuestsWithPets, &house.GuestsWithBabies, &house.BestHouse, &house.Promotion, &house.DistrictEN, &house.DistrictKZ, &house.DistrictRU, &house.PhoneNumber,
		&house.IsVerified, &house.IsSale, &house.IsNewest, &house.IsHot, &house.IsFeatured, &house.IsDiscount,
		&house.CreatedAt, &house.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.ConstraintName == "houses_type_id_fkey":
				return house, shared.Invalid(fmt.Sprintf("type_id %d does not exist", h.TypeID.Int()))
			case pgErr.ConstraintName == "houses_city_id_fkey":
				return house, shared.Invalid("city_id does not exist")
			case pgErr.ConstraintName == "houses_country_id_fkey":
				return house, shared.Invalid("country_id does not exist")
			case pgErr.ConstraintName == "houses_owner_id_fkey":
				return house, shared.Invalid(fmt.Sprintf("owner_id %d does not exist", h.OwnerID))
			case pgErr.Code == pgUniqueViolation:
				return house, model.ErrSlugExists
			}
		}
		return house, err
	}

	if err := linkHouseCategories(ctx, tx, house.ID, h.CategoryIDs); err != nil {
		return model.House{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.House{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	house.CategoryIDs = uniqueInts(h.CategoryIDs)

	return house, nil
}

func linkHouseCategories(ctx context.Context, tx pgx.Tx, houseID int, categoryIDs []int) error {
	if _, err := tx.Exec(ctx, `DELETE FROM house_category WHERE house_id = $1`, houseID); err != nil {
		return fmt.Errorf("failed to reset categories: %w", err)
	}

	ids := uniqueInts(categoryIDs)
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO house_category (house_id, category_id)
		SELECT $1, unnest($2::int[])
		ON CONFLICT (house_id, category_id) DO NOTHING`,
		houseID, ids,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "house_category_category_id_fkey" {
			return model.ErrCategoryRefInvalid
		}
		return fmt.Errorf("failed to link categories: %w", err)
	}

	return nil
}

func uniqueInts(values []int) []int {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[int]struct{}, len(values))
	unique := make([]int, 0, len(values))

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}

	return unique
}

func (r *HouseRepository) getForUpdate(ctx context.Context, slug string) (model.House, error) {
	var house model.House
	query := `
		SELECT id, name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
			description_en, description_kz, description_ru,
			address_en, address_kz, address_ru,
			lng, lat, is_active, priority,
			comments_ru, comments_en, comments_kz,
			owner_id, type_id, city_id, country_id, guests_with_pets, guests_with_babies, best_house, promotion,
			is_verified, is_sale, is_newest, is_hot, is_featured, is_discount,
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
		&house.GuestsWithPets, &house.GuestsWithBabies, &house.BestHouse, &house.Promotion,
		&house.IsVerified, &house.IsSale, &house.IsNewest, &house.IsHot, &house.IsFeatured, &house.IsDiscount,
		&house.DistrictEN, &house.DistrictKZ, &house.DistrictRU,
		&house.PhoneNumber, &house.CreatedAt, &house.UpdatedAt,
	)
	return house, err
}

func (r *HouseRepository) Update(ctx context.Context, slug string, h schema.HouseUpdateRequest) (model.House, error) {
	house, err := r.getForUpdate(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.House{}, fmt.Errorf("%w: slug %q", model.ErrHouseNotFound, slug)
		}
		return house, err
	}

	applyHouseUpdate(&house, h)
	house.UpdatedAt = time.Now()

	slugValue := utils.GenerateSlug(house.Slug, house.NameEN, house.NameKZ, house.NameRU)

	exists, err := r.SlugExistsExceptID(ctx, slugValue, house.ID)
	if err != nil {
		return house, fmt.Errorf("failed to check slug: %w", err)
	}
	if exists {
		return house, fmt.Errorf("%w: %s", model.ErrSlugExists, slugValue)
	}

	query := `
		UPDATE houses SET
			name_en=$1, name_kz=$2, name_ru=$3, slug=$4, price=$5, rooms_qty=$6, guest_qty=$7, bedroom_qty=$8, bath_qty=$9,
			description_en=$10, description_kz=$11, description_ru=$12,
			address_en=$13, address_kz=$14, address_ru=$15,
			lng=$16, lat=$17, is_active=$18, type_id=$19, city_id=$20, country_id=$21, guests_with_pets=$22, guests_with_babies=$23, best_house=$24,
			promotion=$25, district_en=$26, district_kz=$27, district_ru=$28, phone_number=$29,
			is_verified=$30, is_sale=$31, is_newest=$32, is_hot=$33, is_featured=$34, is_discount=$35, updated_at=$36
		WHERE id=$37
	`

	_, err = r.db.Exec(ctx, query,
		house.NameEN, house.NameKZ, house.NameRU, slugValue, house.Price,
		house.RoomsQty, house.GuestQty, house.BedroomQty, house.BathQty,
		house.DescriptionEN, house.DescriptionKZ, house.DescriptionRU,
		house.AddressEN, house.AddressKZ, house.AddressRU,
		house.Lng, house.Lat, house.IsActive, house.TypeID,
		house.CityID, house.CountryID, house.GuestsWithPets, house.GuestsWithBabies, house.BestHouse,
		house.Promotion, house.DistrictEN, house.DistrictKZ, house.DistrictRU,
		house.PhoneNumber,
		house.IsVerified, house.IsSale, house.IsNewest, house.IsHot, house.IsFeatured, house.IsDiscount,
		house.UpdatedAt, house.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "houses_type_id_fkey":
				return house, shared.Invalid("type_id does not exist")
			case "houses_city_id_fkey":
				return house, shared.Invalid("city_id does not exist")
			case "houses_country_id_fkey":
				return house, shared.Invalid("country_id does not exist")
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
