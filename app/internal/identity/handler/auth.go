package handler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/identity/model"
	"github.com/nurkenspashev92/bookit/internal/identity/schema"
	"github.com/nurkenspashev92/bookit/internal/identity/service"
	"github.com/nurkenspashev92/bookit/internal/shared"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
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
	var req schema.UserCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.userService.Register(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) || errors.Is(err, service.ErrPhoneAlreadyExists) {
			return c.Status(409).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	setAccessCookie(c, resp.AccessToken)
	setRefreshCookie(c, resp.RefreshToken)

	return c.Status(201).JSON(resp)
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
	var req schema.UserLoginRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.userService.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrAccountDisabled) {
			return c.Status(401).JSON(shared.ErrorResponse{Error: err.Error()})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	setAccessCookie(c, resp.AccessToken)
	setRefreshCookie(c, resp.RefreshToken)

	return c.JSON(resp)
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
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		var req schema.RefreshRequest
		if err := json.Unmarshal(c.Body(), &req); err == nil && req.RefreshToken != "" {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		return c.Status(401).JSON(shared.ErrorResponse{Error: "refresh token required"})
	}

	resp, err := h.userService.RefreshTokens(c.Context(), refreshToken)
	if err != nil {
		clearAuthCookies(c)
		return c.Status(401).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	setAccessCookie(c, resp.AccessToken)
	setRefreshCookie(c, resp.RefreshToken)

	return c.JSON(resp)
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
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{Error: "unauthenticated"})
	}

	var req schema.UserUpdateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	authUser, err := h.userService.UpdateProfile(c.Context(), user.ID, req)
	if err != nil {
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
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
	user, ok := c.Locals("user").(model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{Error: "unauthenticated"})
	}

	var req schema.ChangePasswordRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	if err := h.userService.ChangePassword(c.Context(), user.ID, req); err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(401).JSON(shared.ErrorResponse{Error: "old password is incorrect"})
		}
		return c.Status(500).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(shared.MessageResponse{Message: "password changed"})
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
	return c.JSON(shared.MessageResponse{Message: "logged out"})
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
	token := c.Cookies("access_token")
	if token == "" {
		return c.Status(401).JSON(shared.ErrorResponse{Error: "unauthenticated"})
	}

	resp, err := h.userService.Me(c.Context(), token)
	if err != nil {
		return c.Status(401).JSON(shared.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
}

func setAccessCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
	})
}

func setRefreshCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/api/v1/auth",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
}

func clearAuthCookies(c fiber.Ctx) {
	expired := time.Now().Add(-time.Hour)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		Path:     "/api/v1/auth",
		Expires:  expired,
	})
}
