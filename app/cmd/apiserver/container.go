package apiserver

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/cmd/router"
	"github.com/nurkenspashev92/bookit/configs"
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
	"github.com/nurkenspashev92/bookit/pkg/cache"
	"github.com/nurkenspashev92/bookit/pkg/store"
)

type infrastructure struct {
	db     *pgxpool.Pool
	s3     *aws.AwsS3Client
	cache  *cache.Cache
	awsCfg *configs.AwsConfig
}

type container struct {
	db       *pgxpool.Pool
	services *router.Services
	release  []func()
}

func newContainer() (*container, error) {
	awsCfg := configs.NewAwsConfig()

	database, err := store.NewPostgresDb(configs.NewDBConfig())
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	s3client, err := aws.NewAwsS3Client(awsCfg.S3Region, awsCfg.S3AccessKey, awsCfg.S3SecretKey, awsCfg.S3Bucket)
	if err != nil {
		database.Conn.Close()
		return nil, fmt.Errorf("s3 client: %w", err)
	}

	appCache, err := initializers.NewCache(configs.NewCacheConfig(), configs.NewRedisConfig())
	if err != nil {
		database.Conn.Close()
		return nil, fmt.Errorf("cache: %w", err)
	}

	infra := infrastructure{
		db:     database.Conn,
		s3:     s3client,
		cache:  appCache,
		awsCfg: awsCfg,
	}

	return &container{
		db:       infra.db,
		services: buildServices(infra),
		release: []func(){
			func() { _ = appCache.Close() },
			database.Conn.Close,
		},
	}, nil
}

func (c *container) Close() {
	for _, release := range c.release {
		release()
	}
}

func buildServices(infra infrastructure) *router.Services {
	likeRepo := interactionrepo.NewHouseLikeRepository(infra.db, infra.awsCfg)
	bookingRepo := bookingrepo.NewBookingRepository(infra.db)

	jwtService := identitysvc.NewJWTService(configs.NewAuthConfig())
	userRepo := identityrepo.NewUserRepository(infra.db)

	return &router.Services{
		Cache:     infra.cache,
		JWT:       jwtService,
		User:      identitysvc.NewUserService(userRepo, jwtService, infra.awsCfg),
		Avatar:    identitysvc.NewAvatarService(userRepo, infra.s3, infra.awsCfg),
		House:     propertysvc.NewHouseService(propertyrepo.NewHouseRepository(infra.db, infra.awsCfg), likeRepo, bookingRepo, infra.cache),
		Image:     propertysvc.NewImageService(propertyrepo.NewHouseImageRepository(infra.db), infra.s3, infra.cache),
		Category:  propertysvc.NewCategoryService(propertyrepo.NewCategoryRepository(infra.db)),
		Type:      propertysvc.NewTypeService(propertyrepo.NewTypeRepository(infra.db)),
		HouseLike: interactionsvc.NewHouseLikeService(likeRepo),
		Booking:   bookingsvc.NewBookingService(bookingRepo),
		Country:   locationsvc.NewCountryService(locationrepo.NewCountryRepository(infra.db)),
		City:      locationsvc.NewCityService(locationrepo.NewCityRepository(infra.db)),
		FAQ:       contentsvc.NewFAQService(contentrepo.NewFAQRepository(infra.db)),
		Inquiry:   contentsvc.NewInquiryService(contentrepo.NewInquiryRepository(infra.db)),
		Stats:     analyticssvc.NewStatsService(analyticsrepo.NewStatsRepository(infra.db)),
	}
}
