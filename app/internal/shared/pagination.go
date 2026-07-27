package shared

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

const (
	PageSizeKey     = "page_size"
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total" example:"100"`
	Page       int         `json:"page" example:"1"`
	PageSize   int         `json:"page_size" example:"10"`
	TotalPages int         `json:"total_pages" example:"10"`
}

func PaginationConfig() paginate.Config {
	return paginate.Config{
		LimitKey:     PageSizeKey,
		DefaultPage:  DefaultPage,
		DefaultLimit: DefaultPageSize,
		MaxLimit:     MaxPageSize,
	}
}

func Page(c fiber.Ctx) *paginate.PageInfo {
	if page, ok := paginate.FromContext(c); ok {
		return page
	}

	return paginate.NewPageInfo(DefaultPage, DefaultPageSize, 0, nil)
}

func Paginated(data interface{}, total int, page *paginate.PageInfo) PaginatedResponse {
	totalPages := 0
	if page.Limit > 0 {
		totalPages = total / page.Limit
		if total%page.Limit > 0 {
			totalPages++
		}
	}

	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page.Page,
		PageSize:   page.Limit,
		TotalPages: totalPages,
	}
}
