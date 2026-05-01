package router

import (
	"time"

	"github.com/Flussen/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/etag"
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
	"github.com/nurkenspashev92/bookit/internal/initializers"
	"github.com/nurkenspashev92/bookit/internal/platform/healthcheck"
	interactionh "github.com/nurkenspashev92/bookit/internal/interaction/handler"
	interactionsvc "github.com/nurkenspashev92/bookit/internal/interaction/service"
	locationh "github.com/nurkenspashev92/bookit/internal/location/handler"
	locationsvc "github.com/nurkenspashev92/bookit/internal/location/service"
	propertyh "github.com/nurkenspashev92/bookit/internal/property/handler"
	propertysvc "github.com/nurkenspashev92/bookit/internal/property/service"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type Services struct {
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

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, svc *Services) *fiber.App {
	app.Use(middleware.CorsHandler)
	app.Use(initializers.NewLogger())
	// Gzip JSON responses — list endpoints shrink 5-10×, smaller wire size means lower latency.
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	// ETag turns repeat GETs into 304 Not Modified — empty body, fastest possible response.
	app.Use(etag.New())

	authHandler := identityh.NewAuthHandler(svc.User)
	houseHandler := propertyh.NewHouseHandler(svc.House)
	houseLikeHandler := interactionh.NewHouseLikeHandler(svc.HouseLike)
	imageHandler := propertyh.NewImageHandler(svc.Image)
	categoryHandler := propertyh.NewCategoryHandler(svc.Category)
	countryHandler := locationh.NewCountryHandler(svc.Country)
	cityHandler := locationh.NewCityHandler(svc.City)
	typeHandler := propertyh.NewTypeHandler(svc.Type)
	avatarHandler := identityh.NewAvatarHandler(svc.Avatar)
	statsHandler := analyticsh.NewStatsHandler(svc.Stats)
	bookingHandler := bookingh.NewBookingHandler(svc.Booking)
	faqHandler := contenth.NewFAQHandler(svc.FAQ, svc.Inquiry)

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", healthcheck.HealthCheck(db))

		auth := apiV1.Group("/auth")
		{
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
			auth.Post("/refresh", authHandler.Refresh)
			auth.Post("/logout", authHandler.Logout)
			auth.Get("/me", authHandler.Me)
			auth.Patch("/me", middleware.AuthRequired(svc.JWT), authHandler.UpdateProfile)
			auth.Patch("/me/password", middleware.AuthRequired(svc.JWT), authHandler.ChangePassword)
			auth.Post("/me/avatar",
				middleware.AuthRequired(svc.JWT),
				middleware.UploadLimits(10*1024*1024, 30*time.Second),
				avatarHandler.Upload,
			)
			auth.Delete("/me/avatar", middleware.AuthRequired(svc.JWT), avatarHandler.Delete)
		}

		category := apiV1.Group("/categories")
		{
			category.Get("/", categoryHandler.GetAll)
			category.Get("/:id", categoryHandler.GetByID)
			category.Post("", middleware.AuthRequired(svc.JWT), categoryHandler.Create)
			category.Patch("/:id", middleware.AuthRequired(svc.JWT), categoryHandler.Update)
			category.Delete("/:id", middleware.AuthRequired(svc.JWT), categoryHandler.Delete)
		}

		country := apiV1.Group("/countries")
		{
			country.Get("/", countryHandler.GetAll)
			country.Get("/:id", countryHandler.GetByID)
			country.Post("/", middleware.AuthRequired(svc.JWT), countryHandler.Create)
			country.Patch("/:id", middleware.AuthRequired(svc.JWT), countryHandler.Update)
			country.Delete("/:id", middleware.AuthRequired(svc.JWT), countryHandler.Delete)
		}

		city := apiV1.Group("/cities")
		{
			city.Get("/", cityHandler.GetAll)
			city.Get("/:id", cityHandler.GetByID)
			city.Post("/", middleware.AuthRequired(svc.JWT), cityHandler.Create)
			city.Patch("/:id", middleware.AuthRequired(svc.JWT), cityHandler.Update)
			city.Delete("/:id", middleware.AuthRequired(svc.JWT), cityHandler.Delete)
		}

		types := apiV1.Group("/types")
		{
			types.Get("/", typeHandler.GetAll)
			types.Get("/:id", typeHandler.GetByID)
			types.Post("/", middleware.AuthRequired(svc.JWT), typeHandler.Create)
			types.Patch("/:id", middleware.AuthRequired(svc.JWT), typeHandler.Update)
			types.Delete("/:id", middleware.AuthRequired(svc.JWT), typeHandler.Delete)
		}

		faq := apiV1.Group("/faqs")
		{
			faq.Get("/", faqHandler.GetAll)
			faq.Get("/:id", faqHandler.GetByID)
			faq.Post("/", middleware.AuthRequired(svc.JWT), faqHandler.Create)
			faq.Patch("/:id", middleware.AuthRequired(svc.JWT), faqHandler.Update)
			faq.Delete("/:id", middleware.AuthRequired(svc.JWT), faqHandler.Delete)
		}

		inquiry := apiV1.Group("/inquiry")
		{
			inquiry.Get("/", faqHandler.GetInquiries)
			inquiry.Get("/:id", faqHandler.GetInquiryByID)
			inquiry.Post("/", middleware.AuthRequired(svc.JWT), faqHandler.CreateInquiry)
			inquiry.Patch("/:id", middleware.AuthRequired(svc.JWT), faqHandler.UpdateInquiry)
			inquiry.Delete("/:id", middleware.AuthRequired(svc.JWT), faqHandler.DeleteInquiry)
		}

		apiV1.Get("/my-houses", middleware.AuthRequired(svc.JWT), houseHandler.MyHouses)

		bookings := apiV1.Group("/bookings")
		{
			bookings.Post("/", middleware.AuthRequired(svc.JWT), bookingHandler.Create)
			bookings.Get("/", middleware.AuthRequired(svc.JWT), bookingHandler.GetMyBookings)
			bookings.Get("/owner", middleware.AuthRequired(svc.JWT), bookingHandler.GetOwnerBookings)
			bookings.Get("/:id", middleware.AuthRequired(svc.JWT), bookingHandler.GetByID)
			bookings.Patch("/:id/status", middleware.AuthRequired(svc.JWT), bookingHandler.UpdateStatus)
		}

		stats := apiV1.Group("/stats")
		{
			stats.Get("/dashboard", middleware.AuthRequired(svc.JWT), statsHandler.Dashboard)
			stats.Get("/houses", middleware.AuthRequired(svc.JWT), statsHandler.HouseStats)
			stats.Get("/houses/:slug", middleware.AuthRequired(svc.JWT), statsHandler.HouseDetail)
			stats.Get("/charts", middleware.AuthRequired(svc.JWT), statsHandler.Charts)
		}

		houses := apiV1.Group("/houses")
		{
			houses.Get("/", middleware.AuthOptional(svc.JWT), houseHandler.GetAll)
			houses.Post("/", middleware.AuthRequired(svc.JWT), houseHandler.Create)
			houses.Patch("/:slug", middleware.AuthRequired(svc.JWT), houseHandler.Update)
			houses.Delete("/:slug", middleware.AuthRequired(svc.JWT), houseHandler.Delete)

			houses.Get("/check-slug", houseHandler.CheckSlug)
			houses.Get("/liked", middleware.AuthRequired(svc.JWT), houseLikeHandler.UserLikedHouses)
			houses.Delete("/images/:image_id", middleware.AuthRequired(svc.JWT), imageHandler.Delete)

			houses.Get("/:slug", middleware.AuthOptional(svc.JWT), houseHandler.GetBySlug)

			houses.Post("/:slug/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Like)
			houses.Delete("/:slug/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Unlike)
			houses.Get("/:slug/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Status)

			houses.Post("/:slug/images",
				middleware.AuthRequired(svc.JWT),
				middleware.UploadLimits(50*1024*1024, 2*time.Minute),
				imageHandler.Upload,
			)
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
