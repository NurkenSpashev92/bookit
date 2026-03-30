package apiserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/cmd/router"
	"github.com/nurkenspashev92/bookit/configs"
	_ "github.com/nurkenspashev92/bookit/docs"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/services"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type ApiApp struct {
	*fiber.App
}

func (app *ApiApp) Run() {
	done := make(chan bool, 1)
	cfgDb := configs.NewDBConfig()
	cfgAws := configs.NewAwsConfig()
	cfgJwt := configs.NewAuthConfig()

	database, err := store.NewPostgresDb(cfgDb)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer database.Close()

	s3client, err := aws.NewAwsS3Client(
		cfgAws.S3Region,
		cfgAws.S3AccessKey,
		cfgAws.S3SecretKey,
		cfgAws.S3Bucket,
	)
	if err != nil {
		log.Fatalf("Failed to init S3: %v", err)
	}

	db := database.Conn

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	houseRepo := repositories.NewHouseRepository(db, cfgAws)
	houseLikeRepo := repositories.NewHouseLikeRepository(db)
	imageRepo := repositories.NewHouseImageRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	countryRepo := repositories.NewCountryRepository(db)
	cityRepo := repositories.NewCityRepository(db)
	typeRepo := repositories.NewTypeRepository(db)
	faqRepo := repositories.NewFAQRepository(db)
	inquiryRepo := repositories.NewInquiryRepository(db)

	// Services
	jwtService := services.NewJWTService(cfgJwt)
	userService := services.NewUserService(userRepo, jwtService)
	houseService := services.NewHouseService(houseRepo)
	houseLikeService := services.NewHouseLikeService(houseLikeRepo)
	imageService := services.NewImageService(imageRepo, s3client)
	categoryService := services.NewCategoryService(categoryRepo, s3client, cfgAws)
	countryService := services.NewCountryService(countryRepo)
	cityService := services.NewCityService(cityRepo)
	typeService := services.NewTypeService(typeRepo, s3client, cfgAws)
	faqService := services.NewFAQService(faqRepo)
	inquiryService := services.NewInquiryService(inquiryRepo)

	svc := &router.Services{
		User:      userService,
		JWT:       jwtService,
		House:     houseService,
		HouseLike: houseLikeService,
		Image:     imageService,
		Category: categoryService,
		Country:  countryService,
		City:     cityService,
		Type:     typeService,
		FAQ:      faqService,
		Inquiry:  inquiryService,
	}

	app.App = router.RegisterRoutes(app.App, db, svc)

	go func() {
		appPort := os.Getenv("APP_PORT")
		err := app.Listen("0.0.0.0:" + appPort)
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go app.Shutdown(done)
	<-done
	log.Println("Graceful shutdown complete.")
}

func (app *ApiApp) Shutdown(done chan<- bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}
