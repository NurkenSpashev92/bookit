package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type citySeed struct {
	nameEN, nameKZ, nameRU, postalCode string
}

type cityRef struct {
	id        int
	countryID int
}

var defaultCities = []citySeed{
	{nameEN: "Astana", nameKZ: "Астана", nameRU: "Астана", postalCode: "010000"},
	{nameEN: "Almaty", nameKZ: "Алматы", nameRU: "Алматы", postalCode: "050000"},
	{nameEN: "Shymkent", nameKZ: "Шымкент", nameRU: "Шымкент", postalCode: "160000"},
	{nameEN: "Karaganda", nameKZ: "Қарағанды", nameRU: "Караганда", postalCode: "100000"},
	{nameEN: "Aktau", nameKZ: "Ақтау", nameRU: "Актау", postalCode: "130000"},
	{nameEN: "Turkestan", nameKZ: "Түркістан", nameRU: "Туркестан", postalCode: "161200"},
	{nameEN: "Aktobe", nameKZ: "Ақтөбе", nameRU: "Актобе", postalCode: "030000"},
	{nameEN: "Taraz", nameKZ: "Тараз", nameRU: "Тараз", postalCode: "080000"},
	{nameEN: "Oskemen", nameKZ: "Өскемен", nameRU: "Усть-Каменогорск", postalCode: "070000"},
	{nameEN: "Pavlodar", nameKZ: "Павлодар", nameRU: "Павлодар", postalCode: "140000"},
	{nameEN: "Semey", nameKZ: "Семей", nameRU: "Семей", postalCode: "071400"},
	{nameEN: "Oral", nameKZ: "Орал", nameRU: "Уральск", postalCode: "090000"},
	{nameEN: "Kostanay", nameKZ: "Қостанай", nameRU: "Костанай", postalCode: "110000"},
	{nameEN: "Kyzylorda", nameKZ: "Қызылорда", nameRU: "Кызылорда", postalCode: "120000"},
	{nameEN: "Atyrau", nameKZ: "Атырау", nameRU: "Атырау", postalCode: "060000"},
	{nameEN: "Petropavl", nameKZ: "Петропавл", nameRU: "Петропавловск", postalCode: "150000"},
	{nameEN: "Kokshetau", nameKZ: "Көкшетау", nameRU: "Кокшетау", postalCode: "020000"},
	{nameEN: "Taldykorgan", nameKZ: "Талдықорған", nameRU: "Талдыкорган", postalCode: "040000"},
	{nameEN: "Konayev", nameKZ: "Қонаев", nameRU: "Конаев", postalCode: "040800"},
	{nameEN: "Zhezkazgan", nameKZ: "Жезқазған", nameRU: "Жезказган", postalCode: "100600"},
	{nameEN: "Ekibastuz", nameKZ: "Екібастұз", nameRU: "Экибастуз", postalCode: "141200"},
	{nameEN: "Rudny", nameKZ: "Рудный", nameRU: "Рудный", postalCode: "111500"},
	{nameEN: "Temirtau", nameKZ: "Теміртау", nameRU: "Темиртау", postalCode: "101400"},
	{nameEN: "Zhanaozen", nameKZ: "Жаңаөзен", nameRU: "Жанаозен", postalCode: "130200"},
	{nameEN: "Balkhash", nameKZ: "Балқаш", nameRU: "Балхаш", postalCode: "100300"},
	{nameEN: "Kentau", nameKZ: "Кентау", nameRU: "Кентау", postalCode: "161000"},
	{nameEN: "Ridder", nameKZ: "Риддер", nameRU: "Риддер", postalCode: "070300"},
}

func ensureCities(ctx context.Context, conn *pgxpool.Pool, countryID int) ([]cityRef, error) {
	refs := make([]cityRef, 0, len(defaultCities))

	for _, city := range defaultCities {
		id, err := ensureRow(ctx, conn, "city", city.nameEN,
			`SELECT id FROM cities WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO cities (name_en, name_kz, name_ru, postall_code, country_id, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			 RETURNING id`,
			city.nameEN, city.nameKZ, city.nameRU, city.postalCode, countryID,
		)
		if err != nil {
			return nil, err
		}
		refs = append(refs, cityRef{id: id, countryID: countryID})
	}

	return refs, nil
}
