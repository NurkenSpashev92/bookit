package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
)

type Deps struct {
	FAQ     FAQService
	Inquiry InquiryService
}

func RegisterRoutes(api fiber.Router, deps Deps, guards shared.Guards) {
	shared.RegisterCRUD(api.Group("/faqs"), guards, NewFAQHandler(deps.FAQ))
	shared.RegisterCRUD(api.Group("/inquiry"), guards, NewInquiryHandler(deps.Inquiry))
}
