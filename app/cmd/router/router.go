package router

import (
	"github.com/Flussen/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/configs"
	_ "github.com/nurkenspashev92/bookit/docs"
	"github.com/nurkenspashev92/bookit/internal/handlers"
	"github.com/nurkenspashev92/bookit/internal/initializers"
	"github.com/nurkenspashev92/bookit/internal/services"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, s3 *aws.AwsS3Client, cfg *configs.AwsConfig, jwtService *services.JWTService) *fiber.App {
	app.Use(middleware.CorsHandler)
	app.Use(initializers.NewLogger())

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handlers.HealthCheck(db))

		// Authentication
		auth := apiV1.Group("/auth")
		{
			auth.Post("/register", handlers.Register(db, jwtService))
			auth.Post("/login", handlers.Login(db, jwtService))
			auth.Post("/logout", middleware.AuthRequired(jwtService), handlers.Logout())
			auth.Get("/me", middleware.AuthRequired(jwtService), handlers.Me(db, jwtService))
		}

		// Categories
		category := apiV1.Group("/categories")
		{
			category.Get("/", handlers.GetCategories(db, cfg))
			category.Get("/:id", handlers.GetCategory(db, cfg))
			category.Post("", middleware.AuthRequired(jwtService), handlers.CreateCategory(db, s3, cfg))
			category.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateCategory(db, s3, cfg))
			category.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteCategory(db, s3))
		}

		// Country
		country := apiV1.Group("/countries")
		{
			country.Get("/", handlers.GetCountries(db))
			country.Get("/:id", handlers.GetCountry(db))
			country.Post("/", middleware.AuthRequired(jwtService), handlers.CreateCountry(db))
			country.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateCountry(db))
			country.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteCountry(db))
		}

		// City
		city := apiV1.Group("/cities")
		{
			city.Get("/", handlers.GetCities(db))
			city.Get("/:id", handlers.GetCity(db))
			city.Post("/", middleware.AuthRequired(jwtService), handlers.CreateCity(db))
			city.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateCity(db))
			city.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteCity(db))
		}

		// Types
		types := apiV1.Group("/types")
		{
			types.Get("/", handlers.GetTypes(db, s3, cfg))
			types.Get("/:id", handlers.GetTypeByID(db, s3, cfg))
			types.Post("/", middleware.AuthRequired(jwtService), handlers.CreateType(db, s3, cfg))
			types.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateType(db, s3, cfg))
			types.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteType(db, s3))
		}

		// FAQ
		faq := apiV1.Group("/faqs")
		{
			faq.Get("/", handlers.GetFAQs(db))
			faq.Get("/:id", handlers.GetFAQByID(db))
			faq.Post("/", middleware.AuthRequired(jwtService), handlers.CreateFAQ(db))
			faq.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateFAQ(db))
			faq.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteFAQ(db))
		}

		// Inquiry
		inquiry := apiV1.Group("/inquiry")
		{
			inquiry.Get("/", handlers.GetInquiries(db))
			inquiry.Get("/:id", handlers.GetInquiryByID(db))
			inquiry.Post("/", middleware.AuthRequired(jwtService), handlers.CreateInquiry(db))
			inquiry.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateInquiry(db))
			inquiry.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteInquiry(db))
		}

		// Houses
		houses := apiV1.Group("/houses")
		{
			houses.Get("/", handlers.GetHouses(db, cfg))
			houses.Get("/:id", handlers.GetHouseByID(db, cfg))
			houses.Post("/", middleware.AuthRequired(jwtService), handlers.CreateHouse(db))
			houses.Patch("/:id", middleware.AuthRequired(jwtService), handlers.UpdateHouse(db))
			houses.Delete("/:id", middleware.AuthRequired(jwtService), handlers.DeleteHouse(db))

			houses.Post("/:id/images", middleware.AuthRequired(jwtService), handlers.UploadHouseImages(db, s3, cfg))
			houses.Delete("/images/:image_id", middleware.AuthRequired(jwtService), handlers.DeleteHouseImage(db, s3))
		}
	}

	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
