package apiserver

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"

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
		fiberlog.Fatalw("failed to initialize the database", "error", err)
	}
	defer database.Close()

	s3client, err := aws.NewAwsS3Client(
		cfgAws.S3Region,
		cfgAws.S3AccessKey,
		cfgAws.S3SecretKey,
		cfgAws.S3Bucket,
	)
	if err != nil {
		fiberlog.Fatalw("failed to initialize S3 client", "error", err)
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
		fiberlog.Fatalw("failed to connect to Redis", "error", err)
	}
	defer houseCache.Close()

	statsRepo := analyticsrepo.NewStatsRepository(db)
	bookingRepo := bookingrepo.NewBookingRepository(db)

	jwtService := identitysvc.NewJWTService(cfgJwt)
	userService := identitysvc.NewUserService(userRepo, jwtService, cfgAws)
	houseService := propertysvc.NewHouseService(houseRepo, houseLikeRepo, bookingRepo, houseCache)
	houseLikeService := interactionsvc.NewHouseLikeService(houseLikeRepo)
	imageService := propertysvc.NewImageService(imageRepo, s3client, houseCache)
	avatarService := identitysvc.NewAvatarService(userRepo, s3client, cfgAws)
	categoryService := propertysvc.NewCategoryService(categoryRepo)
	countryService := locationsvc.NewCountryService(countryRepo)
	cityService := locationsvc.NewCityService(cityRepo)
	typeService := propertysvc.NewTypeService(typeRepo)
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
			fiberlog.Infow("pprof listening", "addr", "0.0.0.0:"+pprofPort)
			if err := http.ListenAndServe("0.0.0.0:"+pprofPort, nil); err != nil {
				fiberlog.Errorw("pprof server stopped", "error", err)
			}
		}()
	}

	go func() {
		appPort := os.Getenv("APP_PORT")
		fiberlog.Infow("http server starting", "addr", "0.0.0.0:"+appPort)

		err := app.Listen("0.0.0.0:" + appPort)
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	go app.Shutdown(done)
	<-done
	fiberlog.Info("graceful shutdown complete")
}

func (app *ApiApp) Shutdown(done chan<- bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	fiberlog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		fiberlog.Errorw("server forced to shutdown", "error", err)
	}

	fiberlog.Info("server exiting")

	done <- true
}
