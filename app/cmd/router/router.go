package router

import (
	"github.com/Flussen/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/paginate"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/nurkenspashev92/bookit/docs"
	analyticssvc "github.com/nurkenspashev92/bookit/internal/analytics/service"
	bookingsvc "github.com/nurkenspashev92/bookit/internal/booking/service"
	contentsvc "github.com/nurkenspashev92/bookit/internal/content/service"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
	interactionsvc "github.com/nurkenspashev92/bookit/internal/interaction/service"
	locationsvc "github.com/nurkenspashev92/bookit/internal/location/service"
	"github.com/nurkenspashev92/bookit/internal/platform/healthcheck"
	propertysvc "github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/cache"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

var cachedResponseRoutes = map[string]string{
	"/api/v1/houses":  "houses",
	"/api/v1/houses/": "houses",
}

var quietLogPaths = map[string]struct{}{
	"/api/v1/healthcheck": {},
}

type Services struct {
	Cache     *cache.Cache
	User      *identitysvc.UserService
	JWT       *identitysvc.JWTService
	House     *propertysvc.HouseService
	HouseLike *interactionsvc.HouseLikeService
	Image     *propertysvc.ImageService
	Avatar    *identitysvc.AvatarService
	Category  *propertysvc.CategoryService
	Country   *locationsvc.CountryService
	City      *locationsvc.CityService
	Type      *propertysvc.TypeService
	FAQ       *contentsvc.FAQService
	Inquiry   *contentsvc.InquiryService
	Stats     *analyticssvc.StatsService
	Booking   *bookingsvc.BookingService
}

type guards struct {
	required fiber.Handler
	optional fiber.Handler
}

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, svc *Services) *fiber.App {
	useMiddleware(app, svc)

	access := guards{
		required: middleware.AuthRequired(svc.JWT),
		optional: middleware.AuthOptional(svc.JWT),
	}

	apiV1 := app.Group("/api/v1")
	apiV1.Get("/healthcheck", healthcheck.HealthCheck(db))

	registerIdentityRoutes(apiV1, svc, access)
	registerPropertyRoutes(apiV1, svc, access)
	registerLocationRoutes(apiV1, svc, access)
	registerContentRoutes(apiV1, svc, access)
	registerBookingRoutes(apiV1, svc, access)
	registerAnalyticsRoutes(apiV1, svc, access)

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}

func useMiddleware(app *fiber.App, svc *Services) {
	app.Use(middleware.Cors())
	app.Use(requestid.New())
	app.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
		QuietPaths: quietLogPaths,
	}))
	app.Use(middleware.RecoverPanic())
	app.Use(middleware.ResponseCache(middleware.ResponseCacheConfig{
		Cache:  svc.Cache,
		Routes: cachedResponseRoutes,
	}))
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(etag.New())
	app.Use(paginate.New(shared.PaginationConfig()))
}
