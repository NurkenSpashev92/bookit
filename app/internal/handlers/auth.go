package handlers

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schemas.UserCreateRequest true "User data"
// @Success 201 {object} schemas.AuthResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req schemas.UserCreateRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.userService.Register(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	setCookie(c, resp.Token)

	return c.Status(201).JSON(resp)
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body schemas.UserLoginRequest true "Login data"
// @Success 200 {object} schemas.AuthResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	if cookie := c.Cookies("jwt"); cookie != "" {
		resp, err := h.userService.ValidateTokenAndGetUser(cookie)
		if err == nil {
			return c.JSON(resp)
		}
	}

	var req schemas.UserLoginRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	resp, err := h.userService.Login(c.Context(), req)
	if err != nil {
		return c.Status(401).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	setCookie(c, resp.Token)

	return c.JSON(resp)
}

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.MessageResponse
// @Failure 401 {object} schemas.ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	})
	return c.JSON(schemas.MessageResponse{Message: "logged out"})
}

// Me godoc
// @Summary Get current authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.AuthResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Security ApiKeyAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(401).JSON(schemas.ErrorResponse{Error: "unauthenticated"})
	}

	resp, err := h.userService.Me(c.Context(), cookie)
	if err != nil {
		return c.Status(401).JSON(schemas.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(resp)
}

func setCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: "Lax",
		Path:     "/",
	})
}
