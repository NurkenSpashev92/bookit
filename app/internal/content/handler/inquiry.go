package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type InquiryService interface {
	GetAll(ctx context.Context, search string) ([]schema.Inquiry, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.Inquiry, int, error)
	GetByID(ctx context.Context, id int) (schema.Inquiry, error)
	Create(ctx context.Context, req schema.InquiryCreateRequest) (schema.Inquiry, error)
	Update(ctx context.Context, id int, req schema.InquiryUpdateRequest) (schema.Inquiry, error)
	Delete(ctx context.Context, id int) error
}

type InquiryHandler struct {
	inquiryService InquiryService
}

func NewInquiryHandler(inquiryService InquiryService) *InquiryHandler {
	return &InquiryHandler{inquiryService: inquiryService}
}

// GetInquiries godoc
// @Summary Get all inquiries
// @Description Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags Inquiry
// @Produce json
// @Param search query string false "Search by email/phone/text"
// @Param page query int false "Page number (enables the paginated envelope)"
// @Param page_size query int false "Items per page (default 20, max 100)"
// @Success 200 {array} schema.Inquiry
// @Failure 500 {object} shared.ErrorResponse
// @Router /inquiries [get]
func (h *InquiryHandler) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.Inquiry, error) {
			return h.inquiryService.GetAll(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.Inquiry, int, error) {
			return h.inquiryService.GetAllPaginated(ctx, search, limit, offset)
		},
	)
}

// GetInquiryByID godoc
// @Summary Get inquiry by ID
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} schema.Inquiry
// @Failure 404 {object} shared.ErrorResponse
// @Router /inquiries/{id} [get]
func (h *InquiryHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	inquiry, err := h.inquiryService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(inquiry)
}

// CreateInquiry godoc
// @Summary Create an inquiry
// @Tags Inquiry
// @Accept json
// @Produce json
// @Param inquiry body schema.InquiryCreateRequest true "Inquiry data"
// @Success 201 {object} schema.Inquiry
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries [post]
func (h *InquiryHandler) Create(c fiber.Ctx) error {
	request, err := shared.Bind[schema.InquiryCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	inquiry, err := h.inquiryService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, inquiry)
}

// UpdateInquiry godoc
// @Summary Update an inquiry
// @Tags Inquiry
// @Accept json
// @Produce json
// @Param id path int true "Inquiry ID"
// @Param inquiry body schema.InquiryUpdateRequest true "Inquiry data"
// @Success 200 {object} schema.Inquiry
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [patch]
func (h *InquiryHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.InquiryUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	inquiry, err := h.inquiryService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(inquiry)
}

// DeleteInquiry godoc
// @Summary Delete an inquiry
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [delete]
func (h *InquiryHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.inquiryService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "Inquiry deleted")
}
