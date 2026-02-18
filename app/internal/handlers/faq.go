package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

// GetFAQs godoc
// @Summary Get all FAQs
// @Tags FAQ
// @Produce json
// @Success 200 {array} schemas.FAQ
// @Failure 500 {object} schemas.ErrorResponse
// @Router /faqs [get]
func GetFAQs(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewFAQRepository(db)
		faqs, err := repo.GetAll(c.Context())
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(faqs)
	}
}

// GetFAQByID godoc
// @Summary Get FAQ by ID
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} schemas.FAQ
// @Failure 404 {object} schemas.ErrorResponse
// @Router /faqs/{id} [get]
func GetFAQByID(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewFAQRepository(db)
		faq, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "FAQ not found"})
		}

		return c.JSON(faq)
	}
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
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /faqs [post]
func CreateFAQ(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req schemas.FAQCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewFAQRepository(db)
		faq, err := repo.Create(c.Context(), req)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(201).JSON(faq)
	}
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
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /faqs/{id} [patch]
func UpdateFAQ(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		var req schemas.FAQUpdateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewFAQRepository(db)
		faq, err := repo.Update(c.Context(), id, req)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(faq)
	}
}

// DeleteFAQ godoc
// @Summary Delete a FAQ
// @Tags FAQ
// @Produce json
// @Param id path int true "FAQ ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /faqs/{id} [delete]
func DeleteFAQ(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewFAQRepository(db)
		if err := repo.Delete(c.Context(), id); err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(schemas.MessageResponse{Message: "FAQ deleted"})
	}
}

// GetInquiries godoc
// @Summary Get all inquiries
// @Tags Inquiry
// @Produce json
// @Success 200 {array} schemas.Inquiry
// @Failure 500 {object} schemas.ErrorResponse
// @Router /inquiries [get]
func GetInquiries(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		repo := repositories.NewInquiryRepository(db)
		list, err := repo.GetAll(c.Context())
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(list)
	}
}

// GetInquiryByID godoc
// @Summary Get inquiry by ID
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} schemas.Inquiry
// @Failure 404 {object} schemas.ErrorResponse
// @Router /inquiries/{id} [get]
func GetInquiryByID(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewInquiryRepository(db)
		inquiry, err := repo.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "Inquiry not found"})
		}

		return c.JSON(inquiry)
	}
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
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /inquiries [post]
func CreateInquiry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req schemas.InquiryCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewInquiryRepository(db)
		inquiry, err := repo.Create(c.Context(), req)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.Status(201).JSON(inquiry)
	}
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
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [patch]
func UpdateInquiry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		var req schemas.InquiryUpdateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewInquiryRepository(db)
		inquiry, err := repo.Update(c.Context(), id, req)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(inquiry)
	}
}

// DeleteInquiry godoc
// @Summary Delete an inquiry
// @Tags Inquiry
// @Produce json
// @Param id path int true "Inquiry ID"
// @Success 200 {object} schemas.MessageResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /inquiries/{id} [delete]
func DeleteInquiry(db *pgxpool.Pool) fiber.Handler {
	return func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: "invalid id"})
		}

		repo := repositories.NewInquiryRepository(db)
		if err := repo.Delete(c.Context(), id); err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(schemas.MessageResponse{Message: "Inquiry deleted"})
	}
}
