package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/content/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type FAQService interface {
	GetAll(ctx context.Context, search string) ([]schema.FAQ, error)
	GetAllPaginated(ctx context.Context, search string, limit, offset int) ([]schema.FAQ, int, error)
	GetByID(ctx context.Context, id int) (schema.FAQ, error)
	Create(ctx context.Context, req schema.FAQCreateRequest) (schema.FAQ, error)
	Update(ctx context.Context, id int, req schema.FAQUpdateRequest) (schema.FAQ, error)
	Delete(ctx context.Context, id int) error
}

type FAQHandler struct {
	faqService FAQService
}

func NewFAQHandler(faqService FAQService) *FAQHandler {
	return &FAQHandler{faqService: faqService}
}

// GetFAQs godoc
// @Summary Get all FAQs
// @Description Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags FAQ
// @Produce json
// @Param search query string false "Search by question/answer"
// @Param page query int false "Page number (enables the paginated envelope)"
// @Param page_size query int false "Items per page (default 20, max 100)"
// @Success 200 {array} schema.FAQ
// @Failure 500 {object} shared.ErrorResponse
// @Router /faqs [get]
func (h *FAQHandler) GetAll(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.FAQ, error) {
			return h.faqService.GetAll(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.FAQ, int, error) {
			return h.faqService.GetAllPaginated(ctx, search, limit, offset)
		},
	)
}

// GetFAQByID godoc
// @Summary Get FAQ by ID
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} schema.FAQ
// @Failure 404 {object} shared.ErrorResponse
// @Router /faqs/{id} [get]
func (h *FAQHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	faq, err := h.faqService.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(faq)
}

// CreateFAQ godoc
// @Summary Create a FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param faq body schema.FAQCreateRequest true "FAQ data"
// @Success 201 {object} schema.FAQ
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs [post]
func (h *FAQHandler) Create(c fiber.Ctx) error {
	request, err := shared.Bind[schema.FAQCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	faq, err := h.faqService.Create(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return shared.Created(c, faq)
}

// UpdateFAQ godoc
// @Summary Update a FAQ
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path int true "FAQ ID"
// @Param faq body schema.FAQUpdateRequest true "FAQ data"
// @Success 200 {object} schema.FAQ
// @Failure 400 {object} shared.ErrorResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs/{id} [patch]
func (h *FAQHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.FAQUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	faq, err := h.faqService.Update(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(faq)
}

// DeleteFAQ godoc
// @Summary Delete a FAQ
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} shared.MessageResponse
// @Failure 404 {object} shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router /faqs/{id} [delete]
func (h *FAQHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.faqService.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "FAQ deleted")
}
