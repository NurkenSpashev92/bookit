package shared

import "github.com/gofiber/fiber/v3"

type Guards struct {
	Required fiber.Handler
	Optional fiber.Handler
}

type CRUDHandler interface {
	GetAll(c fiber.Ctx) error
	GetByID(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
}

func RegisterCRUD(group fiber.Router, guards Guards, handler CRUDHandler) {
	group.Get("/", handler.GetAll)
	group.Get("/:id", handler.GetByID)
	group.Post("/", guards.Required, handler.Create)
	group.Patch("/:id", guards.Required, handler.Update)
	group.Delete("/:id", guards.Required, handler.Delete)
}
