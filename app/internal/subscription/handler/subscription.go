package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/internal/subscription/schema"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type SubscriptionService interface {
	GetMy(ctx context.Context, userID int) (schema.SubscriptionResponse, error)
	Activate(ctx context.Context, userID int, req schema.SubscriptionActivateRequest) (schema.SubscriptionActivationResponse, error)
	Cancel(ctx context.Context, userID int) error
	GetByID(ctx context.Context, id int) (schema.SubscriptionResponse, error)
	GetByUserID(ctx context.Context, userID int) ([]schema.SubscriptionResponse, error)
	GetAll(ctx context.Context) ([]schema.SubscriptionResponse, error)
	GetAllPaginated(ctx context.Context, limit, offset int) ([]schema.SubscriptionResponse, int, error)
	Create(ctx context.Context, req schema.SubscriptionCreateRequest) (schema.SubscriptionResponse, error)
	Update(ctx context.Context, id int, req schema.SubscriptionUpdateRequest) (schema.SubscriptionResponse, error)
	Delete(ctx context.Context, id int) error
}

type SubscriptionHandler struct {
	service SubscriptionService
}

func NewSubscriptionHandler(service SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// GetMy godoc
// @Summary      Get my subscription
// @Tags         Subscriptions
// @Produce      json
// @Success      200  {object}  schema.SubscriptionResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/me [get]
func (h *SubscriptionHandler) GetMy(c fiber.Ctx) error {
	sub, err := h.service.GetMy(c.Context(), middleware.CurrentUserID(c))
	if err != nil {
		return shared.Fail(c, err)
	}
	return c.JSON(sub)
}

// Activate godoc
// @Summary      Activate or change my subscription plan
// @Description  Activates the given plan for the current user. If the same plan is already active, returns a message with the current end date and does not create a new record. Switching to a different plan deactivates the current one, creates a new active subscription, and updates the user's subscription_type.
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        request body schema.SubscriptionActivateRequest true "Plan to activate"
// @Success      200  {object}  schema.SubscriptionActivationResponse
// @Failure      400  {object}  shared.ErrorResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/activate [post]
func (h *SubscriptionHandler) Activate(c fiber.Ctx) error {
	req, err := shared.Bind[schema.SubscriptionActivateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	res, err := h.service.Activate(c.Context(), middleware.CurrentUserID(c), req)
	if err != nil {
		return shared.Fail(c, err)
	}
	return c.JSON(res)
}

// Cancel godoc
// @Summary      Cancel my subscription plan
// @Description  Deactivates the current user's active plan (if any) and resets their subscription_type to basic. Idempotent.
// @Tags         Subscriptions
// @Produce      json
// @Success      200  {object}  shared.MessageResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/cancel [post]
func (h *SubscriptionHandler) Cancel(c fiber.Ctx) error {
	if err := h.service.Cancel(c.Context(), middleware.CurrentUserID(c)); err != nil {
		return shared.Fail(c, err)
	}
	return shared.OK(c, "subscription plan cancelled")
}

// List godoc
// @Summary      List subscriptions
// @Description  Admin only. Without a `page` query param the response is a plain array. With `page` it is the paginated envelope.
// @Tags         Subscriptions
// @Produce      json
// @Param        page       query int false "Page number (enables the paginated envelope)"
// @Param        page_size  query int false "Items per page (default 20, max 100)"
// @Success      200  {array}   schema.SubscriptionResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      403  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(c fiber.Ctx) error {
	return shared.ListMaybePaginated(c, h.service.GetAll, h.service.GetAllPaginated)
}

// ByUser godoc
// @Summary      List a user's subscriptions
// @Description  Admin only.
// @Tags         Subscriptions
// @Produce      json
// @Param        userId  path  int  true  "User ID"
// @Success      200  {array}   schema.SubscriptionResponse
// @Failure      401  {object}  shared.ErrorResponse
// @Failure      403  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/user/{userId} [get]
func (h *SubscriptionHandler) ByUser(c fiber.Ctx) error {
	userID, err := shared.ParamInt(c, "userId")
	if err != nil {
		return shared.Fail(c, err)
	}

	subs, err := h.service.GetByUserID(c.Context(), userID)
	if err != nil {
		return shared.Fail(c, err)
	}
	return shared.List(c, subs)
}

// GetByID godoc
// @Summary      Get subscription by id
// @Tags         Subscriptions
// @Produce      json
// @Param        id   path  int  true  "Subscription ID"
// @Success      200  {object}  schema.SubscriptionResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	sub, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return shared.Fail(c, err)
	}
	return c.JSON(sub)
}

// Create godoc
// @Summary      Create subscription
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        request body schema.SubscriptionCreateRequest true "Subscription"
// @Success      201  {object}  schema.SubscriptionResponse
// @Failure      400  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(c fiber.Ctx) error {
	req, err := shared.Bind[schema.SubscriptionCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	sub, err := h.service.Create(c.Context(), req)
	if err != nil {
		return shared.Fail(c, err)
	}
	return shared.Created(c, sub)
}

// Update godoc
// @Summary      Update subscription
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id      path  int  true  "Subscription ID"
// @Param        request body  schema.SubscriptionUpdateRequest true "Fields to update"
// @Success      200  {object}  schema.SubscriptionResponse
// @Failure      400  {object}  shared.ErrorResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/{id} [patch]
func (h *SubscriptionHandler) Update(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	req, err := shared.Bind[schema.SubscriptionUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	sub, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return shared.Fail(c, err)
	}
	return c.JSON(sub)
}

// Delete godoc
// @Summary      Delete subscription
// @Tags         Subscriptions
// @Produce      json
// @Param        id   path  int  true  "Subscription ID"
// @Success      200  {object}  shared.MessageResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c fiber.Ctx) error {
	id, err := shared.ParamInt(c, "id")
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return shared.Fail(c, err)
	}
	return shared.OK(c, "subscription deleted")
}
