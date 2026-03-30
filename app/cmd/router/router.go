package router

import (
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
	Category  *services.CategoryService
	Country   *services.CountryService
	City      *services.CityService
	Type      *services.TypeService
	FAQ       *services.FAQService
	Inquiry   *services.InquiryService
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
	faqHandler := handlers.NewFAQHandler(svc.FAQ, svc.Inquiry)

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handlers.HealthCheck(db))

		auth := apiV1.Group("/auth")
		{
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
			auth.Post("/logout", middleware.AuthRequired(svc.JWT), authHandler.Logout)
			auth.Get("/me", middleware.AuthRequired(svc.JWT), authHandler.Me)
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

		houses := apiV1.Group("/houses")
		{
			houses.Get("/", houseHandler.GetAll)
			houses.Post("/", middleware.AuthRequired(svc.JWT), houseHandler.Create)
			houses.Patch("/:id", middleware.AuthRequired(svc.JWT), houseHandler.Update)
			houses.Delete("/:id", middleware.AuthRequired(svc.JWT), houseHandler.Delete)

			houses.Get("/check-slug", houseHandler.CheckSlug)
			houses.Get("/liked", middleware.AuthRequired(svc.JWT), houseLikeHandler.UserLikedHouses)
			houses.Delete("/images/:image_id", middleware.AuthRequired(svc.JWT), imageHandler.Delete)

			houses.Get("/:slug", houseHandler.GetBySlug)

			houses.Post("/:id/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Like)
			houses.Delete("/:id/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Unlike)
			houses.Get("/:id/like", middleware.AuthRequired(svc.JWT), houseLikeHandler.Status)

			houses.Post("/:id/images", middleware.AuthRequired(svc.JWT), imageHandler.Upload)
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
