package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

const (
	totalHouses      = 100
	housesPerOwner   = 2
	imagesPerHouse   = 5
	housePhoneNumber = "+77001234567"

	imagesDir = "img"

	uploadConcurrency = 16

	defaultPassword = "1q2w3e4r"

	adminEmail     = "admin@bookit.kz"
	adminFirstName = "admin"

	ownerEmailFormat = "owner%d@bookit.kz"
	ownerNameFormat  = "owner%d"
)

func main() {
	started := time.Now()

	if err := run(); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Printf("Seeding complete in %s", time.Since(started).Round(time.Millisecond))
}

func run() error {
	ctx := context.Background()

	cfgDb := configs.NewDBConfig()
	cfgAws := configs.NewAwsConfig()

	conn, err := pgxpool.New(ctx, cfgDb.DatabaseURL())
	if err != nil {
		return err
	}
	defer conn.Close()

	s3client, err := aws.NewAwsS3Client(
		cfgAws.S3Region, cfgAws.S3AccessKey, cfgAws.S3SecretKey, cfgAws.S3Bucket,
	)
	if err != nil {
		return err
	}

	images, err := prepareImages(imagesDir)
	if err != nil {
		return err
	}
	log.Printf("Prepared %d images", len(images))

	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordHash := string(hashed)

	if _, err := ensureAdmin(ctx, conn, passwordHash); err != nil {
		return err
	}

	ownerIDs, err := ensureOwners(ctx, conn, ownersNeeded(), passwordHash)
	if err != nil {
		return err
	}

	typeIDs, err := ensureTypes(ctx, conn)
	if err != nil {
		return err
	}

	countryIDs, err := ensureCountries(ctx, conn)
	if err != nil {
		return err
	}

	cities, err := ensureCities(ctx, conn, countryIDs[0])
	if err != nil {
		return err
	}

	categoryIDs, err := ensureCategories(ctx, conn)
	if err != nil {
		return err
	}

	convenienceIDs, err := ensureConveniences(ctx, conn)
	if err != nil {
		return err
	}

	if err := ensureFAQs(ctx, conn); err != nil {
		return err
	}

	houseIDs, err := insertHouses(ctx, conn, buildHouseRows(ownerIDs, typeIDs, cities))
	if err != nil {
		return err
	}

	if err := linkHouseCategories(ctx, conn, houseIDs, categoryIDs); err != nil {
		return err
	}

	if err := linkHouseConveniences(ctx, conn, houseIDs, convenienceIDs); err != nil {
		return err
	}

	if err := backfillHouseLocations(ctx, conn, cities); err != nil {
		return err
	}
	if err := backfillHouseCategories(ctx, conn, categoryIDs); err != nil {
		return err
	}
	if err := backfillHouseConveniences(ctx, conn, convenienceIDs); err != nil {
		return err
	}

	if len(houseIDs) == 0 {
		log.Println("No new houses created — nothing to upload")
		return nil
	}

	return insertImages(ctx, conn, uploadHouseImages(ctx, s3client, houseIDs, images))
}

func ownersNeeded() int {
	return (totalHouses + housesPerOwner - 1) / housesPerOwner
}
