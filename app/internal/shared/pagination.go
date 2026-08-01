package shared

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/paginate"
)

const (
	PageSizeKey     = "page_size"
	PageKey         = "page"
	DefaultPage     = 1
	DefaultPageSize = 10
	// ListDefaultPageSize is the fallback page size for the optional pagination
	// on reference list endpoints (categories, types, countries, ...).
	ListDefaultPageSize = 20
	MaxPageSize         = 100
	maxPageNumber       = 1_000_000
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

// PageParams holds the optional pagination request parsed from query params.
type PageParams struct {
	Page     int
	PageSize int
}

func (p PageParams) Limit() int  { return p.PageSize }
func (p PageParams) Offset() int { return (p.Page - 1) * p.PageSize }

// WantsPagination reports whether the request opted into the paginated
// envelope by supplying a `page` query param. When false, list handlers keep
// returning the plain JSON array for backward compatibility.
func WantsPagination(c fiber.Ctx) bool {
	return c.Query(PageKey) != ""
}

// ParsePageParams reads `page` (default 1) and `page_size` (default 20, capped
// at 100) from the query string. Out-of-range or invalid values fall back to
// the defaults, while a page_size above the cap is clamped to MaxPageSize.
func ParsePageParams(c fiber.Ctx) PageParams {
	page := QueryIntInRange(c, PageKey, DefaultPage, 1, maxPageNumber)

	pageSize, err := strconv.Atoi(c.Query(PageSizeKey))
	switch {
	case err != nil || pageSize < 1:
		pageSize = ListDefaultPageSize
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	}

	return PageParams{Page: page, PageSize: pageSize}
}

// ListMaybePaginated serves a list endpoint that returns a plain JSON array by
// default and the paginated envelope when the request opts in via `page`. The
// caller supplies the two service methods; everything else (param parsing,
// error handling, envelope building) is shared.
func ListMaybePaginated[T any](
	c fiber.Ctx,
	getAll func(context.Context) ([]T, error),
	getPage func(ctx context.Context, limit, offset int) ([]T, int, error),
) error {
	if WantsPagination(c) {
		p := ParsePageParams(c)
		items, total, err := getPage(c.Context(), p.Limit(), p.Offset())
		if err != nil {
			return Fail(c, err)
		}

		return c.JSON(PageEnvelope(Items(items), total, p))
	}

	items, err := getAll(c.Context())
	if err != nil {
		return Fail(c, err)
	}

	return List(c, items)
}

// PageEnvelope builds the PaginatedResponse envelope from parsed page params.
func PageEnvelope(data interface{}, total int, p PageParams) PaginatedResponse {
	return newPaginatedResponse(data, total, p.Page, p.PageSize)
}

func Paginated(data interface{}, total int, page *paginate.PageInfo) PaginatedResponse {
	return newPaginatedResponse(data, total, page.Page, page.Limit)
}

func newPaginatedResponse(data interface{}, total, page, pageSize int) PaginatedResponse {
	totalPages := 0
	if pageSize > 0 {
		totalPages = total / pageSize
		if total%pageSize > 0 {
			totalPages++
		}
	}

	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
