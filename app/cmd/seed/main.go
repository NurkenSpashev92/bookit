package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/imageproc"
)

const (
	totalHouses = 100
	imagesDir   = "img"
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

func main() {
	cfgDb := configs.NewDBConfig()
	cfgAws := configs.NewAwsConfig()

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, cfgDb.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	s3client, err := aws.NewAwsS3Client(
		cfgAws.S3Region, cfgAws.S3AccessKey, cfgAws.S3SecretKey, cfgAws.S3Bucket,
	)
	if err != nil {
		log.Fatalf("failed to init S3: %v", err)
	}

	imageFiles, err := loadImageFiles(imagesDir)
	if err != nil {
		log.Fatalf("failed to load images: %v", err)
	}
	log.Printf("Loaded %d image files", len(imageFiles))

	var ownerID, typeID int
	if err := conn.QueryRow(ctx, `SELECT id FROM users LIMIT 1`).Scan(&ownerID); err != nil {
		log.Fatalf("no users found: %v", err)
	}
	if err := conn.QueryRow(ctx, `SELECT id FROM types LIMIT 1`).Scan(&typeID); err != nil {
		log.Fatalf("no types found: %v", err)
	}

	var cityID, countryID *int
	var cid, cnid int
	if err := conn.QueryRow(ctx, `SELECT id FROM cities LIMIT 1`).Scan(&cid); err == nil {
		cityID = &cid
	}
	if err := conn.QueryRow(ctx, `SELECT id FROM countries LIMIT 1`).Scan(&cnid); err == nil {
		countryID = &cnid
	}

	for i := 0; i < totalHouses; i++ {
		nameEN := fmt.Sprintf("%s %d", pick(namesEN), i+1)
		nameKZ := fmt.Sprintf("%s %d", pick(namesKZ), i+1)
		nameRU := fmt.Sprintf("%s %d", pick(namesRU), i+1)
		houseSlug := slug.Make(nameEN)

		price := 10000 + rand.Intn(490000)
		rooms := 1 + rand.Intn(6)
		guests := 1 + rand.Intn(10)
		bedrooms := 1 + rand.Intn(4)
		baths := 1 + rand.Intn(3)
		priority := rand.Intn(10)
		lng := 51.0 + rand.Float64()*2
		lat := 71.0 + rand.Float64()*2

		var houseID int
		err := conn.QueryRow(ctx, `
			INSERT INTO houses (
				name_en, name_kz, name_ru, slug, price, rooms_qty, guest_qty, bedroom_qty, bath_qty,
				description_en, description_kz, description_ru,
				address_en, address_kz, address_ru,
				lng, lat, is_active, priority, owner_id, type_id, city_id, country_id,
				guests_with_pets, best_house, promotion,
				district_en, district_kz, district_ru, phone_number
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30)
			RETURNING id`,
			nameEN, nameKZ, nameRU, houseSlug, price, rooms, guests, bedrooms, baths,
			pick(descriptionsEN), pick(descriptionsKZ), pick(descriptionsRU),
			pick(addressesEN), pick(addressesKZ), pick(addressesRU),
			lng, lat, true, priority, ownerID, typeID, cityID, countryID,
			rand.Intn(2) == 1, rand.Intn(5) == 0, rand.Intn(5) == 0,
			pick(districtsEN), pick(districtsKZ), pick(districtsRU), "+77001234567",
		).Scan(&houseID)
		if err != nil {
			if strings.Contains(err.Error(), "23505") {
				log.Printf("[%d/%d] slug '%s' exists, skipping", i+1, totalHouses, houseSlug)
				continue
			}
			log.Fatalf("[%d/%d] failed to insert house: %v", i+1, totalHouses, err)
		}

		for j := 0; j < 5; j++ {
			imgData := imageFiles[j%len(imageFiles)]

			result, err := imageproc.ProcessBytes(imgData.data, imgData.ext)
			if err != nil {
				log.Printf("[house %d] failed to process image %d: %v", houseID, j+1, err)
				continue
			}

			ts := time.Now().UnixNano() + int64(j)
			originalKey := fmt.Sprintf("houses/original/%d_%d.jpg", houseID, ts)
			thumbKey := fmt.Sprintf("houses/thumbnail/%d_%d.webp", houseID, ts)

			if _, err := s3client.UploadCompressed(ctx, originalKey, result.Original, "image/jpeg"); err != nil {
				log.Printf("[house %d] failed to upload original %d: %v", houseID, j+1, err)
				continue
			}
			if _, err := s3client.UploadCompressed(ctx, thumbKey, result.Thumbnail, "image/webp"); err != nil {
				log.Printf("[house %d] failed to upload thumbnail %d: %v", houseID, j+1, err)
				continue
			}

			_, err = conn.Exec(ctx, `
				INSERT INTO images (original, thumbnail, mimetype, width, height, size, house_id)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				originalKey, thumbKey, result.Mime, result.Width, result.Height, result.Size, houseID,
			)
			if err != nil {
				log.Printf("[house %d] failed to save image %d to db: %v", houseID, j+1, err)
			}
		}

		log.Printf("[%d/%d] House '%s' created with 5 images", i+1, totalHouses, nameEN)
	}

	log.Println("Seeding complete!")
}

type imageFile struct {
	data []byte
	ext  string
}

func loadImageFiles(dir string) ([]imageFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []imageFile
	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, imageFile{data: data, ext: ext})
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no image files found in %s", dir)
	}
	return files, nil
}

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}
