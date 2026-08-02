package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type faqSeed struct {
	QRU, ARU string
	QKZ, AKZ string
	QEN, AEN string
}

var defaultFAQs = []faqSeed{
	{
		QRU: "Как забронировать жильё?",
		ARU: "Откройте страницу дома, выберите даты и число гостей, затем нажмите «Забронировать». Заявка уйдёт хозяину на подтверждение.",
		QKZ: "Тұрғын үйді қалай брондауға болады?",
		AKZ: "Үй бетін ашып, күндер мен қонақ санын таңдап, «Брондау» түймесін басыңыз. Өтінім иесіне растауға жіберіледі.",
		QEN: "How do I book a place?",
		AEN: "Open a home’s page, choose dates and the number of guests, then press “Book”. The request goes to the host for confirmation.",
	},
	{
		QRU: "Как происходит оплата?",
		ARU: "Оплата проходит напрямую хозяину по его реквизитам (номер телефона или QR), которые отображаются в брони после подтверждения.",
		QKZ: "Төлем қалай жүреді?",
		AKZ: "Төлем расталғаннан кейін бронда көрсетілетін иесінің деректемелері (телефон нөмірі не QR) бойынша тікелей өтеді.",
		QEN: "How does payment work?",
		AEN: "Payment goes directly to the host using their details (phone number or QR), shown in the booking once it is confirmed.",
	},
	{
		QRU: "Нужно ли регистрироваться?",
		ARU: "Для бронирования и добавления жилья — да. Просмотр объявлений доступен без регистрации.",
		QKZ: "Тіркелу қажет пе?",
		AKZ: "Брондау мен үй қосу үшін — иә. Хабарландыруларды тіркелмей-ақ қарауға болады.",
		QEN: "Do I need to register?",
		AEN: "To book or list a home, yes. Browsing listings works without registration.",
	},
	{
		QRU: "Как стать хозяином?",
		ARU: "Перейдите в раздел «Сдать жильё», заполните информацию о доме и загрузите фотографии. После модерации объявление появится в поиске.",
		QKZ: "Иесі болу үшін не істеу керек?",
		AKZ: "«Тұрғын үй тапсыру» бөліміне өтіп, үй туралы ақпаратты толтырып, фотосуреттерді жүктеңіз. Модерациядан кейін хабарландыру іздеуде пайда болады.",
		QEN: "How do I become a host?",
		AEN: "Go to “List your place”, fill in your home’s details and upload photos. After moderation the listing appears in search.",
	},
	{
		QRU: "Почему моего объявления нет в поиске?",
		ARU: "Новые объявления публикуются после проверки администратором. Неактивные дома видны только вам в разделе «Мои дома».",
		QKZ: "Хабарландыруым неге іздеуде жоқ?",
		AKZ: "Жаңа хабарландырулар әкімші тексергеннен кейін жарияланады. Белсенді емес үйлер тек сізге «Менің үйлерім» бөлімінде көрінеді.",
		QEN: "Why isn’t my listing in search?",
		AEN: "New listings are published after an admin review. Inactive homes are visible only to you under “My homes”.",
	},
	{
		QRU: "Можно ли отменить бронирование?",
		ARU: "Да. Откройте «Мои бронирования», выберите бронь и нажмите «Отменить». Пока хозяин не подтвердил заявку, отмена мгновенная.",
		QKZ: "Бронды болдырмауға бола ма?",
		AKZ: "Иә. «Менің брондарым» бөлімін ашып, бронды таңдап, «Болдырмау» түймесін басыңыз. Иесі растамайынша, болдырмау бірден орындалады.",
		QEN: "Can I cancel a booking?",
		AEN: "Yes. Open “My bookings”, select the booking and press “Cancel”. While the host hasn’t confirmed, cancellation is instant.",
	},
	{
		QRU: "Как связаться с хозяином?",
		ARU: "После подтверждения брони телефон хозяина отображается в деталях бронирования и на странице дома.",
		QKZ: "Иесімен қалай байланысуға болады?",
		AKZ: "Бронь расталғаннан кейін иесінің телефоны бронь мәліметтерінде және үй бетінде көрсетіледі.",
		QEN: "How do I contact the host?",
		AEN: "Once the booking is confirmed, the host’s phone appears in the booking details and on the home’s page.",
	},
	{
		QRU: "На каких языках доступен сайт?",
		ARU: "Bookit работает на русском, казахском и английском — переключайте язык в шапке сайта.",
		QKZ: "Сайт қандай тілдерде қолжетімді?",
		AKZ: "Bookit орыс, қазақ және ағылшын тілдерінде жұмыс істейді — тілді сайттың жоғарғы жағында ауыстырыңыз.",
		QEN: "What languages does the site support?",
		AEN: "Bookit works in Russian, Kazakh and English — switch the language in the site header.",
	},
	{
		QRU: "Как редактировать моё объявление?",
		ARU: "В разделе «Мои дома» выберите дом и нажмите «Редактировать» — можно менять описание, цены, фотографии и обложку.",
		QKZ: "Хабарландыруымды қалай өңдеймін?",
		AKZ: "«Менің үйлерім» бөлімінде үйді таңдап, «Өңдеу» түймесін басыңыз — сипаттаманы, бағаларды, фотосуреттерді және мұқабаны өзгертуге болады.",
		QEN: "How do I edit my listing?",
		AEN: "Under “My homes”, select the home and press “Edit” — you can change the description, prices, photos and cover.",
	},
	{
		QRU: "Что даёт подписка хозяину?",
		ARU: "Тариф определяет доступное число объявлений и дополнительные возможности. Подробности — на странице «Тарифы».",
		QKZ: "Жазылым иесіне не береді?",
		AKZ: "Тариф қолжетімді хабарландыру санын және қосымша мүмкіндіктерді анықтайды. Толығырақ — «Тарифтер» бетінде.",
		QEN: "What does a subscription give a host?",
		AEN: "The plan sets the number of listings available and extra features. Details are on the “Pricing” page.",
	},
}

func ensureFAQs(ctx context.Context, conn *pgxpool.Pool) error {
	for _, f := range defaultFAQs {
		if _, err := ensureRow(ctx, conn, "faq", f.QRU,
			`SELECT id FROM faq WHERE question_ru = $1 LIMIT 1`,
			`INSERT INTO faq (question_ru, answer_ru, question_kz, answer_kz, question_en, answer_en, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
			 RETURNING id`,
			f.QRU, f.ARU, f.QKZ, f.AKZ, f.QEN, f.AEN,
		); err != nil {
			return err
		}
	}
	return nil
}
