package handler

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

type UserService interface {
	Register(ctx context.Context, req schema.UserCreateRequest) (*schema.AuthResponse, error)
	Login(ctx context.Context, req schema.UserLoginRequest) (*schema.AuthResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*schema.AuthResponse, error)
	UpdateProfile(ctx context.Context, userID int, req schema.UserUpdateRequest) (*schema.AuthUser, error)
	ChangePassword(ctx context.Context, userID int, req schema.ChangePasswordRequest) error
	Me(ctx context.Context, accessToken string) (*schema.AuthResponse, error)
	ListUsers(ctx context.Context, search string) ([]schema.AdminUser, error)
	ListUsersPaginated(ctx context.Context, search string, limit, offset int) ([]schema.AdminUser, int, error)
	UpdateUserFlags(ctx context.Context, id int, req schema.UserAdminUpdateRequest) (*schema.AdminUser, error)
}

type AuthHandler struct {
	userService UserService
}

func NewAuthHandler(userService UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schema.UserCreateRequest true "User data"
// @Success 201 {object} schema.AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 409 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	request, err := shared.Bind[schema.UserCreateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	response, err := h.userService.Register(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	issueAuthCookies(c, response.AccessToken, response.RefreshToken)

	return shared.Created(c, response)
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schema.UserLoginRequest true "Login data"
// @Success 200 {object} schema.AuthResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	request, err := shared.Bind[schema.UserLoginRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	response, err := h.userService.Login(c.Context(), request)
	if err != nil {
		return shared.Fail(c, err)
	}

	issueAuthCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Uses refresh token from cookie or body to issue new token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body schema.RefreshRequest false "Refresh token (optional, can use cookie)"
// @Success 200 {object} schema.AuthResponse
// @Failure 401 {object} shared.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	refreshToken := h.refreshToken(c)
	if refreshToken == "" {
		return shared.Fail(c, shared.Unauthorized("refresh token required"))
	}

	response, err := h.userService.RefreshTokens(c.Context(), refreshToken)
	if err != nil {
		clearAuthCookies(c)
		return shared.Fail(c, err)
	}

	issueAuthCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(response)
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schema.UserUpdateRequest true "Fields to update"
// @Success 200 {object} schema.AuthUser
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me [patch]
func (h *AuthHandler) UpdateProfile(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.UserUpdateRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	authUser, err := h.userService.UpdateProfile(c.Context(), user.ID, request)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(authUser)
}

// ChangePassword godoc
// @Summary Change user password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body schema.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} shared.MessageResponse
// @Failure 400 {object} shared.ErrorResponse
// @Failure 401 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me/password [patch]
func (h *AuthHandler) ChangePassword(c fiber.Ctx) error {
	user, err := middleware.CurrentUser(c)
	if err != nil {
		return shared.Fail(c, err)
	}

	request, err := shared.Bind[schema.ChangePasswordRequest](c)
	if err != nil {
		return shared.Fail(c, err)
	}

	if err := h.userService.ChangePassword(c.Context(), user.ID, request); err != nil {
		return shared.Fail(c, err)
	}

	return shared.OK(c, "password changed")
}

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Produce json
// @Success 200 {object} shared.MessageResponse
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	clearAuthCookies(c)

	return shared.OK(c, "logged out")
}

// Me godoc
// @Summary Get current authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} schema.AuthResponse
// @Failure 401 {object} shared.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c fiber.Ctx) error {
	token := c.Cookies(accessCookieName)
	if token == "" {
		return shared.Fail(c, shared.Unauthorized("unauthenticated"))
	}

	response, err := h.userService.Me(c.Context(), token)
	if err != nil {
		return shared.Fail(c, err)
	}

	return c.JSON(response)
}

func (h *AuthHandler) refreshToken(c fiber.Ctx) string {
	if token := c.Cookies(refreshCookieName); token != "" {
		return token
	}

	request, err := shared.Bind[schema.RefreshRequest](c)
	if err != nil {
		return ""
	}

	return request.RefreshToken
}
