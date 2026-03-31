package router

import (
	"time"

	"github.com/Flussen/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/nurkenspashev92/bookit/docs"
	"github.com/nurkenspashev92/bookit/internal/handlers"
	"github.com/nurkenspashev92/bookit/internal/initializers"
	"github.com/nurkenspashev92/bookit/internal/services"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type Services struct {
	User      *services.UserService
	JWT       *services.JWTService
	House     *services.HouseService
	HouseLike *services.HouseLikeService
	Image     *services.ImageService
	Avatar    *services.AvatarService
	Category  *services.CategoryService
	Country   *services.CountryService
	City      *services.CityService
	Type      *services.TypeService
	FAQ       *services.FAQService
	Inquiry   *services.InquiryService
	Stats     *services.StatsService
}

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, svc *Services) *fiber.App {
	app.Use(middleware.CorsHandler)
	app.Use(initializers.NewLogger())

	authHandler := handlers.NewAuthHandler(svc.User)
	houseHandler := handlers.NewHouseHandler(svc.House)
	houseLikeHandler := handlers.NewHouseLikeHandler(svc.HouseLike)
	imageHandler := handlers.NewImageHandler(svc.Image)
	categoryHandler := handlers.NewCategoryHandler(svc.Category)
	countryHandler := handlers.NewCountryHandler(svc.Country)
	cityHandler := handlers.NewCityHandler(svc.City)
	typeHandler := handlers.NewTypeHandler(svc.Type)
	avatarHandler := handlers.NewAvatarHandler(svc.Avatar)
	statsHandler := handlers.NewStatsHandler(svc.Stats)
	faqHandler := handlers.NewFAQHandler(svc.FAQ, svc.Inquiry)

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handlers.HealthCheck(db))

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
