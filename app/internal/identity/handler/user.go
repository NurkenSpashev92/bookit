package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type UserAdminService interface {
	ListUsers(ctx context.Context) ([]schema.AdminUser, error)
	ListUsersPaginated(ctx context.Context, limit, offset int) ([]schema.AdminUser, int, error)
	UpdateUserFlags(ctx context.Context, id int, req schema.UserAdminUpdateRequest) (*schema.AdminUser, error)
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
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success      200  {array}   schema.AdminUser
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      403  {object}  shared.ErrorResponse
// @Failure      500  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users [get]
func (h *UserHandler) List(c fiber.Ctx) error {
	return shared.ListMaybePaginated(c, h.userService.ListUsers, h.userService.ListUsersPaginated)
}

// Update godoc
// @Summary      Update user status flags
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

	user, err := h.userService.UpdateUserFlags(c.Context(), id, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(user)
}
