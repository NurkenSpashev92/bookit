package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/pkg/utils"
)

type categorySeed struct {
	nameEN, nameKZ, nameRU string
}

var defaultCategories = []categorySeed{
	{nameEN: "Beachfront", nameKZ: "Жағажай маңы", nameRU: "У пляжа"},
	{nameEN: "Mountains", nameKZ: "Таулар", nameRU: "Горы"},
	{nameEN: "City center", nameKZ: "Қала орталығы", nameRU: "Центр города"},
	{nameEN: "Countryside", nameKZ: "Ауыл", nameRU: "За городом"},
	{nameEN: "Lakeside", nameKZ: "Көл жағасы", nameRU: "У озера"},
	{nameEN: "Ski-in", nameKZ: "Шаңғы курорты", nameRU: "Горнолыжные"},
	{nameEN: "Family friendly", nameKZ: "Отбасыларға", nameRU: "Для семьи"},
	{nameEN: "Business travel", nameKZ: "Іссапар", nameRU: "Командировка"},
	{nameEN: "Luxury", nameKZ: "Люкс", nameRU: "Люкс"},
	{nameEN: "Budget", nameKZ: "Үнемді", nameRU: "Бюджетные"},
	{nameEN: "Pool", nameKZ: "Бассейнмен", nameRU: "С бассейном"},
	{nameEN: "Sauna", nameKZ: "Сауна", nameRU: "Сауна и баня"},
	{nameEN: "Pet friendly", nameKZ: "Үй жануарларымен", nameRU: "С питомцами"},
	{nameEN: "Amazing views", nameKZ: "Керемет көріністер", nameRU: "Потрясающие виды"},
	{nameEN: "Camping", nameKZ: "Кемпинг", nameRU: "Кемпинг"},
	{nameEN: "Cabins", nameKZ: "Шағын үйлер", nameRU: "Домики"},
	{nameEN: "Design", nameKZ: "Дизайнерлік", nameRU: "Дизайнерские"},
	{nameEN: "Historic", nameKZ: "Тарихи", nameRU: "Исторические"},
	{nameEN: "Long stay", nameKZ: "Ұзақ мерзімге", nameRU: "Долгосрочно"},
	{nameEN: "Trending", nameKZ: "Танымал", nameRU: "Популярные"},
	{nameEN: "Islands", nameKZ: "Аралдар", nameRU: "Острова"},
	{nameEN: "Desert", nameKZ: "Шөл дала", nameRU: "Пустыня"},
	{nameEN: "Farm stay", nameKZ: "Фермада", nameRU: "На ферме"},
	{nameEN: "Vineyards", nameKZ: "Жүзімдіктер", nameRU: "Виноградники"},
	{nameEN: "Tiny homes", nameKZ: "Кішкентай үйлер", nameRU: "Мини-дома"},
	{nameEN: "Treehouses", nameKZ: "Ағаш үйлер", nameRU: "Дома на деревьях"},
	{nameEN: "Castles", nameKZ: "Сарайлар", nameRU: "Замки"},
	{nameEN: "Hot springs", nameKZ: "Ыстық бұлақтар", nameRU: "Горячие источники"},
	{nameEN: "National parks", nameKZ: "Ұлттық парктер", nameRU: "Нацпарки"},
	{nameEN: "Golf", nameKZ: "Гольф", nameRU: "Гольф"},
}

func ensureCategories(ctx context.Context, conn *pgxpool.Pool) ([]int, error) {
	ids := make([]int, 0, len(defaultCategories))

	for _, category := range defaultCategories {
		id, err := ensureRow(ctx, conn, "category", category.nameEN,
			`SELECT id FROM categories WHERE name_en = $1 LIMIT 1`,
			`INSERT INTO categories (name_en, name_kz, name_ru, slug, is_active, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, TRUE, NOW(), NOW())
			 RETURNING id`,
			category.nameEN, category.nameKZ, category.nameRU,
			utils.GenerateSlug("", category.nameEN, category.nameKZ, category.nameRU),
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
