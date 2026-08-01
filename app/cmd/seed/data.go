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

	defaultTypes = []typeSeed{
		{nameEN: "Apartment", nameKZ: "Пәтер", nameRU: "Квартира"},
		{nameEN: "House", nameKZ: "Үй", nameRU: "Дом"},
		{nameEN: "Villa", nameKZ: "Вилла", nameRU: "Вилла"},
		{nameEN: "Cottage", nameKZ: "Коттедж", nameRU: "Коттедж"},
		{nameEN: "Hotel", nameKZ: "Қонақ үй", nameRU: "Отель"},
		{nameEN: "Studio", nameKZ: "Студия", nameRU: "Студия"},
		{nameEN: "Townhouse", nameKZ: "Таунхаус", nameRU: "Таунхаус"},
		{nameEN: "Guest house", nameKZ: "Қонақжай үй", nameRU: "Гостевой дом"},
		{nameEN: "Hostel", nameKZ: "Хостел", nameRU: "Хостел"},
		{nameEN: "Yurt", nameKZ: "Киіз үй", nameRU: "Юрта"},
	}

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

	defaultConveniences = []string{
		"Wi-Fi", "Air conditioning", "Heating", "Kitchen", "Washer",
		"Dryer", "Free parking", "Pool", "Hot tub", "Gym",
		"TV", "Workspace", "Elevator", "Breakfast", "BBQ grill",
		"Fireplace", "Balcony", "Garden", "Sea view", "Mountain view",
		"Pet friendly", "Smoke alarm", "First aid kit", "Security cameras", "Self check-in",
		"EV charger", "Crib", "Iron", "Hair dryer", "24/7 security",
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

type typeSeed struct {
	nameEN, nameKZ, nameRU string
}

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}
