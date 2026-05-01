package shared

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// PaginatedResponse generic paginated response wrapper
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total" example:"100"`
	Page       int         `json:"page" example:"1"`
	PageSize   int         `json:"page_size" example:"10"`
	TotalPages int         `json:"total_pages" example:"10"`
}

// PaginationParams parsed pagination query params
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

func ParsePagination(c fiber.Ctx) PaginationParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}
