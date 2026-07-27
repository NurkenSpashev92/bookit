package main

import "math/rand"

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

	defaultTypes = []string{"Apartment", "House", "Villa", "Cottage", "Hotel"}

	defaultCountries = []countrySeed{
		{nameEN: "Kazakhstan", nameKZ: "Қазақстан", nameRU: "Казахстан", code: "KZ"},
	}

	defaultCities = []citySeed{
		{nameEN: "Astana", nameKZ: "Астана", nameRU: "Астана", postalCode: "010000"},
		{nameEN: "Almaty", nameKZ: "Алматы", nameRU: "Алматы", postalCode: "050000"},
		{nameEN: "Shymkent", nameKZ: "Шымкент", nameRU: "Шымкент", postalCode: "160000"},
		{nameEN: "Karaganda", nameKZ: "Қарағанды", nameRU: "Караганда", postalCode: "100000"},
		{nameEN: "Aktau", nameKZ: "Ақтау", nameRU: "Актау", postalCode: "130000"},
		{nameEN: "Turkestan", nameKZ: "Түркістан", nameRU: "Туркестан", postalCode: "161200"},
	}

	defaultCategories = []categorySeed{
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
	}
)

type countrySeed struct {
	nameEN, nameKZ, nameRU, code string
}

type citySeed struct {
	nameEN, nameKZ, nameRU, postalCode string
}

type categorySeed struct {
	nameEN, nameKZ, nameRU string
}

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}
