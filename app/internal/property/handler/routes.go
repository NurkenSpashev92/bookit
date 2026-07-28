package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"

	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const (
	houseImagesUploadLimit   = 50 * 1024 * 1024
	houseImagesUploadTimeout = 2 * time.Minute
)

type Deps struct {
	House    HouseService
	Image    ImageService
	Category CategoryService
	Type     TypeService
}

func RegisterRoutes(api fiber.Router, deps Deps, guards shared.Guards) {
	houseHandler := NewHouseHandler(deps.House)
	imageHandler := NewImageHandler(deps.Image)

	api.Get("/my-houses", guards.Required, houseHandler.MyHouses)

	houses := api.Group("/houses")

	houses.Get("/", guards.Optional, houseHandler.GetAll)
	houses.Post("/", guards.Required, houseHandler.Create)

	houses.Get("/check-slug", houseHandler.CheckSlug)
	houses.Delete("/images/:image_id", guards.Required, imageHandler.Delete)

	houses.Get("/:slug", guards.Optional, houseHandler.GetBySlug)
	houses.Patch("/:slug", guards.Required, houseHandler.Update)
	houses.Delete("/:slug", guards.Required, houseHandler.Delete)

	houses.Post("/:slug/images",
		guards.Required,
		middleware.UploadLimits(houseImagesUploadLimit),
		timeout.New(imageHandler.Upload, timeout.Config{Timeout: houseImagesUploadTimeout}),
	)

	shared.RegisterCRUD(api.Group("/categories"), guards, NewCategoryHandler(deps.Category))
	shared.RegisterCRUD(api.Group("/types"), guards, NewTypeHandler(deps.Type))
}
