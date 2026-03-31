package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type FAQHandler struct {
	faqService     *services.FAQService
	inquiryService *services.InquiryService
}

func NewFAQHandler(faqService *services.FAQService, inquiryService *services.InquiryService) *FAQHandler {
	return &FAQHandler{
		faqService:     faqService,
		inquiryService: inquiryService,
	}
}

// GetFAQs godoc
// @Summary Get all FAQs
// @Tags FAQ
// @Produce json
// @Success 200 {array} schemas.FAQ
// @Failure 500 {object} schemas.ErrorResponse
// @Router /faqs [get]
func (h *FAQHandler) GetAll(c fiber.Ctx) error {
	faqs, err := h.faqService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if faqs == nil {
		faqs = []schemas.FAQ{}
	}
	return c.JSON(faqs)
}

// GetFAQByID godoc
// @Summary Get FAQ by ID
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} schemas.FAQ
// @Failure 404 {object} schemas.ErrorResponse
// @Router /faqs/{id} [get]
func (h *FAQHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	faq, err := h.faqService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: "FAQ not found"})
	}

	return c.JSON(faq)
}

// CreateFAQ godoc
// @Summary Create a FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param faq body schemas.FAQCreateRequest true "FAQ data"
// @Success 201 {object} schemas.FAQ
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs [post]
func (h *FAQHandler) Create(c fiber.Ctx) error {
	var req schemas.FAQCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	faq, err := h.faqService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(faq)
}

// UpdateFAQ godoc
// @Summary Update a FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path int true "FAQ ID"
// @Param faq body schemas.FAQUpdateRequest true "FAQ data"
// @Success 200 {object} schemas.FAQ
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs/{id} [patch]
func (h *FAQHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	var req schemas.FAQUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	faq, err := h.faqService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(faq)
}

// DeleteFAQ godoc
// @Summary Delete a FAQ
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs/{id} [delete]
func (h *FAQHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	if err := h.faqService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "FAQ deleted"})
}

// GetInquiries godoc
// @Summary Get all inquiries
// @Tags Inquiry
// @Produce json
// @Success 200 {array} schemas.Inquiry
// @Failure 500 {object} schemas.ErrorResponse
// @Router /inquiries [get]
func (h *FAQHandler) GetInquiries(c fiber.Ctx) error {
	list, err := h.inquiryService.GetAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if list == nil {
		list = []schemas.Inquiry{}
	}
	return c.JSON(list)
}

// GetInquiryByID godoc
// @Summary Get inquiry by ID
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} schemas.Inquiry
// @Failure 404 {object} schemas.ErrorResponse
// @Router /inquiries/{id} [get]
func (h *FAQHandler) GetInquiryByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	inquiry, err := h.inquiryService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: "Inquiry not found"})
	}

	return c.JSON(inquiry)
}

// CreateInquiry godoc
// @Summary Create an inquiry
// @Tags Inquiry
// @Accept json
// @Produce json
// @Param inquiry body schemas.InquiryCreateRequest true "Inquiry data"
// @Success 201 {object} schemas.Inquiry
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries [post]
func (h *FAQHandler) CreateInquiry(c fiber.Ctx) error {
	var req schemas.InquiryCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	inquiry, err := h.inquiryService.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.Status(201).JSON(inquiry)
}

// UpdateInquiry godoc
// @Summary Update an inquiry
// @Tags Inquiry
// @Accept json
// @Produce json
// @Param id path int true "Inquiry ID"
// @Param inquiry body schemas.InquiryUpdateRequest true "Inquiry data"
// @Success 200 {object} schemas.Inquiry
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [patch]
func (h *FAQHandler) UpdateInquiry(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	var req schemas.InquiryUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	inquiry, err := h.inquiryService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(inquiry)
}

// DeleteInquiry godoc
// @Summary Delete an inquiry
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [delete]
func (h *FAQHandler) DeleteInquiry(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
	}

	if err := h.inquiryService.Delete(c.Context(), id); err != nil {
		return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(schemas.MessageResponse{Message: "Inquiry deleted"})
}
