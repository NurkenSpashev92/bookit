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

var (
	namesEN = []string{
		"Cozy Apartment", "Beach Villa", "Mountain Lodge", "City Penthouse", "Lake House",
		"Forest Cabin", "Seaside Cottage", "Luxury Suite", "Garden Flat", "Rooftop Loft",
		"Country House", "Modern Studio", "Royal Palace", "Ocean View", "Sunset Villa",
		"River House", "Snow Lodge", "Desert Oasis", "Harbor View", "Hilltop Estate",
	}
	namesKZ = []string{
		"Жайлы пәтер", "Жағажай виллаcы", "Тау лоджы", "Қалалық пентхаус", "Көл үйі",
		"Орман үйі", "Теңіз коттеджі", "Люкс номер", "Бау пәтер", "Шатыр лофт",
		"Ауыл үйі", "Заманауи студия", "Ханшайым сарайы", "Мұхит көрінісі", "Күн батысы",
		"Өзен үйі", "Қар лоджы", "Шөл оазисі", "Айлақ көрінісі", "Тау шыңы",
	}
	namesRU = []string{
		"Уютная квартира", "Пляжная вилла", "Горный лодж", "Городской пентхаус", "Дом у озера",
		"Лесной домик", "Морской коттедж", "Люкс-сьют", "Садовая квартира", "Лофт на крыше",
		"Загородный дом", "Современная студия", "Королевский дворец", "Вид на океан", "Вилла заката",
		"Речной дом", "Снежный лодж", "Оазис в пустыне", "Вид на гавань", "Усадьба на холме",
	}

	descriptionsEN = []string{
		"A wonderful place to stay with your family and friends.",
		"Perfect getaway for a relaxing vacation.",
		"Enjoy the breathtaking views and modern amenities.",
		"Spacious and comfortable accommodation in a prime location.",
		"Experience luxury living at its finest.",
	}
	descriptionsKZ = []string{
		"Отбасыңызбен және достарыңызбен тұруға тамаша орын.",
		"Демалыс үшін тамаша орын.",
		"Тамаша көріністер мен заманауи ыңғайлылықтарды пайдаланыңыз.",
		"Бірінші дәрежелі орналасқан кең және жайлы тұрғын үй.",
		"Ең жоғары деңгейдегі сәнді өмірді сезініңіз.",
	}
	descriptionsRU = []string{
		"Прекрасное место для проживания с семьей и друзьями.",
		"Идеальное место для расслабляющего отдыха.",
		"Наслаждайтесь потрясающими видами и современными удобствами.",
		"Просторное и комфортное жилье в отличном месте.",
		"Испытайте роскошную жизнь на высшем уровне.",
	}

	addressesEN = []string{
		"123 Main Street", "456 Oak Avenue", "789 Pine Road", "321 Elm Boulevard", "654 Maple Lane",
	}
	addressesKZ = []string{
		"Абай көшесі 123", "Тоқтар көшесі 456", "Назарбаев даңғылы 789", "Бейбітшілік көшесі 321", "Республика көшесі 654",
	}
	addressesRU = []string{
		"ул. Абая 123", "ул. Токтара 456", "пр. Назарбаева 789", "ул. Мира 321", "ул. Республики 654",
	}

	districtsEN = []string{"Downtown", "Uptown", "Midtown", "Suburbs", "Old Town"}
	districtsKZ = []string{"Орталық", "Жоғары қала", "Орта қала", "Іргетас", "Ескі қала"}
	districtsRU = []string{"Центр", "Верхний город", "Средний город", "Пригород", "Старый город"}
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

	guestsWithPets, bestHouse, promotion                        bool
	isVerified, isSale, isNewest, isHot, isFeatured, isDiscount bool
}

const insertHouseSQL = `
	INSERT INTO houses (
		name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
		description_en, description_kz, description_ru,
		address_en, address_kz, address_ru,
		lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
		guests_with_pets, best_house, promotion,
		district_en, district_kz, district_ru, phone_number,
		is_verified, is_sale, is_newest, is_hot, is_featured, is_discount
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36)
	ON CONFLICT (slug) DO NOTHING
	RETURNING id`

func buildHouseRows(ownerIDs, typeIDs []int, cities []cityRef) []houseRow {
	rows := make([]houseRow, 0, totalHouses)

	for i := range totalHouses {
		nameEN := fmt.Sprintf("%s %d", pick(namesEN), i+1)
		cityID, countryID := pickCity(cities)

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
			isVerified:     rand.Intn(3) == 0,
			isSale:         rand.Intn(5) == 0,
			isNewest:       rand.Intn(4) == 0,
			isHot:          rand.Intn(6) == 0,
			isFeatured:     rand.Intn(6) == 0,
			isDiscount:     rand.Intn(5) == 0,
		})
	}

	return rows
}

func pickCity(cities []cityRef) (cityID, countryID *int) {
	if len(cities) == 0 {
		return nil, nil
	}

	city := cities[rand.Intn(len(cities))]
	return &city.id, &city.countryID
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
			r.isVerified, r.isSale, r.isNewest, r.isHot, r.isFeatured, r.isDiscount,
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
