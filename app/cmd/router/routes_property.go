package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"

	interactionh "github.com/nurkenspashev92/bookit/internal/interaction/handler"
	propertyh "github.com/nurkenspashev92/bookit/internal/property/handler"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const (
	houseImagesUploadLimit   = 50 * 1024 * 1024
	houseImagesUploadTimeout = 2 * time.Minute
)

func registerPropertyRoutes(api fiber.Router, svc *Services, access guards) {
	houseHandler := propertyh.NewHouseHandler(svc.House)
	imageHandler := propertyh.NewImageHandler(svc.Image)
	likeHandler := interactionh.NewHouseLikeHandler(svc.HouseLike)
	categoryHandler := propertyh.NewCategoryHandler(svc.Category)
	typeHandler := propertyh.NewTypeHandler(svc.Type)

	api.Get("/my-houses", access.required, houseHandler.MyHouses)

	registerHouseRoutes(api, access, houseHandler, imageHandler, likeHandler)
	registerReferenceRoutes(api.Group("/categories"), access, categoryHandler)
	registerReferenceRoutes(api.Group("/types"), access, typeHandler)
}

func registerHouseRoutes(
	api fiber.Router,
	access guards,
	houseHandler *propertyh.HouseHandler,
	imageHandler *propertyh.ImageHandler,
	likeHandler *interactionh.HouseLikeHandler,
) {
	houses := api.Group("/houses")

	houses.Get("/", access.optional, houseHandler.GetAll)
	houses.Post("/", access.required, houseHandler.Create)

	houses.Get("/check-slug", houseHandler.CheckSlug)
	houses.Get("/liked", access.required, likeHandler.UserLikedHouses)
	houses.Delete("/images/:image_id", access.required, imageHandler.Delete)

	houses.Get("/:slug", access.optional, houseHandler.GetBySlug)
	houses.Patch("/:slug", access.required, houseHandler.Update)
	houses.Delete("/:slug", access.required, houseHandler.Delete)

	houses.Post("/:slug/like", access.required, likeHandler.Like)
	houses.Delete("/:slug/like", access.required, likeHandler.Unlike)
	houses.Get("/:slug/like", access.required, likeHandler.Status)

	houses.Post("/:slug/images",
		access.required,
		middleware.UploadLimits(houseImagesUploadLimit),
		timeout.New(imageHandler.Upload, timeout.Config{Timeout: houseImagesUploadTimeout}),
	)
}
