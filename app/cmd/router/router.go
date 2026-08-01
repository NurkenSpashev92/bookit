package router

import (
	"time"

	"github.com/Flussen/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/paginate"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/nurkenspashev92/bookit/docs"
	analyticsh "github.com/nurkenspashev92/bookit/internal/analytics/handler"
	analyticssvc "github.com/nurkenspashev92/bookit/internal/analytics/service"
	bookingh "github.com/nurkenspashev92/bookit/internal/booking/handler"
	bookingsvc "github.com/nurkenspashev92/bookit/internal/booking/service"
	contenth "github.com/nurkenspashev92/bookit/internal/content/handler"
	contentsvc "github.com/nurkenspashev92/bookit/internal/content/service"
	identityh "github.com/nurkenspashev92/bookit/internal/identity/handler"
	identitysvc "github.com/nurkenspashev92/bookit/internal/identity/service"
	interactionh "github.com/nurkenspashev92/bookit/internal/interaction/handler"
	interactionsvc "github.com/nurkenspashev92/bookit/internal/interaction/service"
	locationh "github.com/nurkenspashev92/bookit/internal/location/handler"
	locationsvc "github.com/nurkenspashev92/bookit/internal/location/service"
	"github.com/nurkenspashev92/bookit/internal/platform/healthcheck"
	propertyh "github.com/nurkenspashev92/bookit/internal/property/handler"
	propertysvc "github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/cache"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const referenceTTL = 15 * time.Second

var cachedResponseRoutes = map[string]middleware.CachedRoute{
	"/api/v1/houses":      {Namespace: "houses"},
	"/api/v1/houses/":     {Namespace: "houses"},
	"/api/v1/categories":  {Namespace: "categories", TTL: referenceTTL},
	"/api/v1/categories/": {Namespace: "categories", TTL: referenceTTL},
	"/api/v1/countries":   {Namespace: "countries", TTL: referenceTTL},
	"/api/v1/countries/":  {Namespace: "countries", TTL: referenceTTL},
	"/api/v1/cities":      {Namespace: "cities", TTL: referenceTTL},
	"/api/v1/cities/":     {Namespace: "cities", TTL: referenceTTL},
	"/api/v1/types":       {Namespace: "types", TTL: referenceTTL},
	"/api/v1/types/":      {Namespace: "types", TTL: referenceTTL},
}

var quietLogPaths = map[string]struct{}{
	"/api/v1/healthcheck": {},
}

type Services struct {
	Cache       *cache.Cache
	User        *identitysvc.UserService
	JWT         *identitysvc.JWTService
	House       *propertysvc.HouseService
	HouseLike   *interactionsvc.HouseLikeService
	Image       *propertysvc.ImageService
	Avatar      *identitysvc.AvatarService
	Category    *propertysvc.CategoryService
	Convenience *propertysvc.ConvenienceService
	Country     *locationsvc.CountryService
	City        *locationsvc.CityService
	Type        *propertysvc.TypeService
	FAQ         *contentsvc.FAQService
	Inquiry     *contentsvc.InquiryService
	Stats       *analyticssvc.StatsService
	Booking     *bookingsvc.BookingService
}

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, svc *Services) *fiber.App {
	useMiddleware(app, svc)

	guards := shared.Guards{
		Required: middleware.AuthRequired(svc.JWT),
		Optional: middleware.AuthOptional(svc.JWT),
		Admin:    middleware.AuthAdmin(svc.JWT),
	}

	apiV1 := app.Group("/api/v1")
	apiV1.Get("/healthcheck", healthcheck.HealthCheck(db))

	identityh.RegisterRoutes(apiV1, identityh.Deps{
		User:   svc.User,
		Avatar: svc.Avatar,
	}, guards)

	interactionh.RegisterRoutes(apiV1, svc.HouseLike, guards)

	propertyh.RegisterRoutes(apiV1, propertyh.Deps{
		House:       svc.House,
		Image:       svc.Image,
		Category:    svc.Category,
		Type:        svc.Type,
		Convenience: svc.Convenience,
	}, guards)

	locationh.RegisterRoutes(apiV1, locationh.Deps{
		Country: svc.Country,
		City:    svc.City,
	}, guards)

	contenth.RegisterRoutes(apiV1, contenth.Deps{
		FAQ:     svc.FAQ,
		Inquiry: svc.Inquiry,
	}, guards)

	bookingh.RegisterRoutes(apiV1, svc.Booking, guards)
	analyticsh.RegisterRoutes(apiV1, svc.Stats, guards)

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
