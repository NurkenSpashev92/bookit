package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/cmd/router"
	"github.com/nurkenspashev92/bookit/configs"
	_ "github.com/nurkenspashev92/bookit/docs"
	analyticsrepo "github.com/nurkenspashev92/bookit/internal/analytics/repository"
	analyticssvc "github.com/nurkenspashev92/bookit/internal/analytics/service"
	bookingrepo "github.com/nurkenspashev92/bookit/internal/booking/repository"
	bookingsvc "github.com/nurkenspashev92/bookit/internal/booking/service"
	contentrepo "github.com/nurkenspashev92/bookit/internal/content/repository"
	contentsvc "github.com/nurkenspashev92/bookit/internal/content/service"
	identityrepo "github.com/nurkenspashev92/bookit/internal/identity/repository"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
	"github.com/nurkenspashev92/bookit/internal/initializers"
	interactionrepo "github.com/nurkenspashev92/bookit/internal/interaction/repository"
	interactionsvc "github.com/nurkenspashev92/bookit/internal/interaction/service"
	locationrepo "github.com/nurkenspashev92/bookit/internal/location/repository"
	locationsvc "github.com/nurkenspashev92/bookit/internal/location/service"
	propertyrepo "github.com/nurkenspashev92/bookit/internal/property/repository"
	propertysvc "github.com/nurkenspashev92/bookit/internal/property/service"
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

	userRepo := identityrepo.NewUserRepository(db)
	houseRepo := propertyrepo.NewHouseRepository(db, cfgAws)
	houseLikeRepo := interactionrepo.NewHouseLikeRepository(db, cfgAws)
	imageRepo := propertyrepo.NewHouseImageRepository(db)
	categoryRepo := propertyrepo.NewCategoryRepository(db)
	countryRepo := locationrepo.NewCountryRepository(db)
	cityRepo := locationrepo.NewCityRepository(db)
	typeRepo := propertyrepo.NewTypeRepository(db)
	faqRepo := contentrepo.NewFAQRepository(db)
	inquiryRepo := contentrepo.NewInquiryRepository(db)

	houseCache, err := initializers.NewCache(configs.NewCacheConfig(), configs.NewRedisConfig())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer houseCache.Close()

	statsRepo := analyticsrepo.NewStatsRepository(db)
	bookingRepo := bookingrepo.NewBookingRepository(db)

	jwtService := identitysvc.NewJWTService(cfgJwt)
	userService := identitysvc.NewUserService(userRepo, jwtService, cfgAws)
	houseService := propertysvc.NewHouseService(houseRepo, houseLikeRepo, bookingRepo, houseCache)
	houseLikeService := interactionsvc.NewHouseLikeService(houseLikeRepo)
	imageService := propertysvc.NewImageService(imageRepo, s3client, houseCache)
	avatarService := identitysvc.NewAvatarService(userRepo, s3client)
	categoryService := propertysvc.NewCategoryService(categoryRepo, s3client, cfgAws)
	countryService := locationsvc.NewCountryService(countryRepo)
	cityService := locationsvc.NewCityService(cityRepo)
	typeService := propertysvc.NewTypeService(typeRepo, s3client, cfgAws)
	statsService := analyticssvc.NewStatsService(statsRepo)
	bookingService := bookingsvc.NewBookingService(bookingRepo)
	faqService := contentsvc.NewFAQService(faqRepo)
	inquiryService := contentsvc.NewInquiryService(inquiryRepo)

	svc := &router.Services{
		Cache:     houseCache,
		User:      userService,
		JWT:       jwtService,
		House:     houseService,
		HouseLike: houseLikeService,
		Image:     imageService,
		Avatar:    avatarService,
		Category:  categoryService,
		Country:   countryService,
		City:      cityService,
		Type:      typeService,
		FAQ:       faqService,
		Inquiry:   inquiryService,
		Stats:     statsService,
		Booking:   bookingService,
	}

	app.App = router.RegisterRoutes(app.App, db, svc)

	if pprofPort := os.Getenv("PPROF_PORT"); pprofPort != "" {
		go func() {
			log.Printf("pprof listening on 0.0.0.0:%s", pprofPort)
			if err := http.ListenAndServe("0.0.0.0:"+pprofPort, nil); err != nil {
				log.Printf("pprof server error: %v", err)
			}
		}()
	}

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
