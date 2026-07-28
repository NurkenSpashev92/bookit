package router

import (
	"github.com/gofiber/fiber/v3"
)

type referenceHandler interface {
	GetAll(c fiber.Ctx) error
	GetByID(c fiber.Ctx) error
	Create(c fiber.Ctx) error
	Update(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
}

func registerReferenceRoutes(group fiber.Router, access guards, handler referenceHandler) {
	group.Get("/", handler.GetAll)
	group.Get("/:id", handler.GetByID)
	group.Post("/", access.required, handler.Create)
	group.Patch("/:id", access.required, handler.Update)
	group.Delete("/:id", access.required, handler.Delete)
}
