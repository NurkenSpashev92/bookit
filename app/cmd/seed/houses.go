package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type houseRow struct {
	nameEN, nameKZ, nameRU string
	slug                   string

	price, rooms, guests, bedrooms, baths, priority int

	descriptionEN, descriptionKZ, descriptionRU string
	addressEN, addressKZ, addressRU             string
	districtEN, districtKZ, districtRU          string

	lng, lat float64

	ownerID, typeID   int
	cityID, countryID *int

	guestsWithPets, bestHouse, promotion bool
}

const insertHouseSQL = `
	INSERT INTO houses (
		name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
		description_en, description_kz, description_ru,
		address_en, address_kz, address_ru,
		lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
		guests_with_pets, best_house, promotion,
		district_en, district_kz, district_ru, phone_number
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
	ON CONFLICT (slug) DO NOTHING
	RETURNING id`

func buildHouseRows(ownerIDs, typeIDs []int, cityID, countryID *int) []houseRow {
	rows := make([]houseRow, 0, totalHouses)

	for i := range totalHouses {
		nameEN := fmt.Sprintf("%s %d", pick(namesEN), i+1)

		rows = append(rows, houseRow{
			nameEN: nameEN,
			nameKZ: fmt.Sprintf("%s %d", pick(namesKZ), i+1),
			nameRU: fmt.Sprintf("%s %d", pick(namesRU), i+1),
			slug:   slug.Make(nameEN),

			price:    10000 + rand.Intn(490000),
			rooms:    1 + rand.Intn(6),
			guests:   1 + rand.Intn(10),
			bedrooms: 1 + rand.Intn(4),
			baths:    1 + rand.Intn(3),
			priority: rand.Intn(10),

			descriptionEN: pick(descriptionsEN),
			descriptionKZ: pick(descriptionsKZ),
			descriptionRU: pick(descriptionsRU),
			addressEN:     pick(addressesEN),
			addressKZ:     pick(addressesKZ),
			addressRU:     pick(addressesRU),
			districtEN:    pick(districtsEN),
			districtKZ:    pick(districtsKZ),
			districtRU:    pick(districtsRU),

			lng: 51.0 + rand.Float64()*2,
			lat: 71.0 + rand.Float64()*2,

			ownerID:   ownerIDs[i/housesPerOwner],
			typeID:    typeIDs[rand.Intn(len(typeIDs))],
			cityID:    cityID,
			countryID: countryID,

			guestsWithPets: rand.Intn(2) == 1,
			bestHouse:      rand.Intn(5) == 0,
			promotion:      rand.Intn(5) == 0,
		})
	}

	return rows
}

func insertHouses(ctx context.Context, conn *pgxpool.Pool, rows []houseRow) ([]int, error) {
	batch := &pgx.Batch{}

	for _, r := range rows {
		batch.Queue(insertHouseSQL,
			r.nameEN, r.nameKZ, r.nameRU, r.slug, r.price, r.rooms, r.guests, r.bedrooms, r.baths,
			r.descriptionEN, r.descriptionKZ, r.descriptionRU,
			r.addressEN, r.addressKZ, r.addressRU,
			r.lng, r.lat, true, r.priority, r.ownerID, r.typeID, r.cityID, r.countryID,
			r.guestsWithPets, r.bestHouse, r.promotion,
			r.districtEN, r.districtKZ, r.districtRU, housePhoneNumber,
		)
	}

	results := conn.SendBatch(ctx, batch)

	ids, skipped, scanErr := scanHouseIDs(results, rows)

	closeErr := results.Close()
	if scanErr != nil {
		return nil, scanErr
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close houses batch: %w", closeErr)
	}

	log.Printf("Houses inserted: %d (skipped as existing: %d)", len(ids), skipped)
	return ids, nil
}

func scanHouseIDs(results pgx.BatchResults, rows []houseRow) (ids []int, skipped int, err error) {
	ids = make([]int, 0, len(rows))

	for _, r := range rows {
		var id int

		switch scanErr := results.QueryRow().Scan(&id); {
		case scanErr == nil:
			ids = append(ids, id)
		case errors.Is(scanErr, pgx.ErrNoRows):
			skipped++
		default:
			return nil, skipped, fmt.Errorf("insert house %q: %w", r.slug, scanErr)
		}
	}

	return ids, skipped, nil
}

func insertImages(ctx context.Context, conn *pgxpool.Pool, rows []imageRow) error {
	if len(rows) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	for _, r := range rows {
		batch.Queue(`
			INSERT INTO images (original, thumbnail, mimetype, width, height, size, house_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			r.original, r.thumbnail, r.mime, r.width, r.height, r.size, r.houseID,
		)
	}

	results := conn.SendBatch(ctx, batch)

	var execErr error
	for range rows {
		if _, err := results.Exec(); err != nil && execErr == nil {
			execErr = fmt.Errorf("insert image: %w", err)
		}
	}

	closeErr := results.Close()
	if execErr != nil {
		return execErr
	}
	if closeErr != nil {
		return fmt.Errorf("close images batch: %w", closeErr)
	}

	log.Printf("Images inserted: %d", len(rows))
	return nil
}
