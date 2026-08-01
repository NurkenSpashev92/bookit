package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type UserAdminService interface {
	ListUsers(ctx context.Context, search string) ([]schema.AdminUser, error)
	ListUsersPaginated(ctx context.Context, search string, limit, offset int) ([]schema.AdminUser, int, error)
	UpdateUser(ctx context.Context, id int, req schema.UserAdminUpdateRequest) (*schema.AdminUser, error)
}

type UserHandler struct {
	userService UserAdminService
}

func NewUserHandler(userService UserAdminService) *UserHandler {
	return &UserHandler{userService: userService}
}

// List godoc
// @Summary      List all users
// @Description  Admin only. Without a `page` query param the response is a plain array. With `page` it is the paginated envelope (shared.PaginatedResponse).
// @Tags         Users
// @Produce      json
// @Param        search     query string false "Search by name/email/phone"
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success      200  {array}   schema.AdminUser
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      403  {object}  shared.ErrorResponse
// @Failure      500  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	search := strings.TrimSpace(c.Query("search"))

	return shared.ListMaybePaginated(c,
		func(ctx context.Context) ([]schema.AdminUser, error) {
			return h.userService.ListUsers(ctx, search)
		},
		func(ctx context.Context, limit, offset int) ([]schema.AdminUser, int, error) {
			return h.userService.ListUsersPaginated(ctx, search, limit, offset)
		},
	)
}

// Update godoc
// @Summary      Update a user
// @Description  Admin only. Partial update of a user's profile fields (first_name, last_name, middle_name, phone_number, email) and status flags (is_active, is_superuser). All fields optional; only provided fields are changed. Changing email to one already used by another user returns 409.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id      path  int                            true  "User ID"
// @Param        request body  schema.UserAdminUpdateRequest  true  "Fields to update"
// @Success      200  {object}  schema.AdminUser
// @Failure      400  {object}  shared.ErrorResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      403  {object}  shared.ErrorResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Failure      409  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users/{id} [patch]
func (h *UserHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.UserAdminUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	user, err := h.userService.UpdateUser(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(user)
}
