package router

import (
	"github.com/gofiber/fiber/v3"

	contenth "github.com/nurkenspashev92/bookit/internal/content/handler"
)

func registerContentRoutes(api fiber.Router, svc *Services, access guards) {
	registerReferenceRoutes(api.Group("/faqs"), access, contenth.NewFAQHandler(svc.FAQ))
	registerReferenceRoutes(api.Group("/inquiry"), access, contenth.NewInquiryHandler(svc.Inquiry))
}
