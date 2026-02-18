package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/internal/schemas"
	"github.com/nurkenspashev92/bookit/internal/services"
)

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
func Register(db *pgxpool.Pool, jwtService *services.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req schemas.UserCreateRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewUserRepository(db)
		user, err := repo.Create(c.Context(), req)
		if err != nil {
			if err.Error() == fmt.Sprintf("email %s already exists", req.Email) {
				return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
			}
			return c.Status(500).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		token, err := jwtService.GenerateToken(user)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: "failed to generate token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: "Lax",
			Path:     "/",
		})

		authUser := schemas.AuthUser{
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
			Avatar:     user.Avatar,
		}

		return c.Status(201).JSON(schemas.AuthResponse{
			User:  authUser,
			Token: token,
		})
	}
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
func Login(db *pgxpool.Pool, jwtService *services.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		if cookie := c.Cookies("jwt"); cookie != "" {
			if user, err := jwtService.ValidateToken(cookie); err == nil {
				authUser := schemas.AuthUser{
					Email:      user.Email,
					FirstName:  user.FirstName,
					LastName:   user.LastName,
					MiddleName: user.MiddleName,
					Avatar:     user.Avatar,
				}
				return c.Status(200).JSON(schemas.AuthResponse{
					User:  authUser,
					Token: cookie,
				})
			}
		}

		var req schemas.UserLoginRequest
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(400).JSON(schemas.ErrorResponse{Error: err.Error()})
		}

		repo := repositories.NewUserRepository(db)
		user, err := repo.GetByEmail(c.Context(), req.Email)
		if err != nil {
			return c.Status(401).JSON(schemas.ErrorResponse{Error: "invalid credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.Status(401).JSON(schemas.ErrorResponse{Error: "invalid credentials"})
		}

		token, err := jwtService.GenerateToken(user)
		if err != nil {
			return c.Status(500).JSON(schemas.ErrorResponse{Error: "failed to generate token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: "Lax",
			Path:     "/",
		})

		authUser := schemas.AuthUser{
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
			Avatar:     user.Avatar,
		}

		return c.JSON(schemas.AuthResponse{
			User:  authUser,
			Token: token,
		})
	}
}

// Logout godoc
// @Summary Logout user
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.MessageResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /auth/logout [post]
func Logout() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			HTTPOnly: true,
			Expires:  time.Now().Add(-time.Hour),
		})
		return c.JSON(schemas.MessageResponse{Message: "logged out"})
	}
}

// Me godoc
// @Summary Get current authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.AuthResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure      401   {object}  schemas.ErrorResponse "Unauthorized"
// @Security     ApiKeyAuth
// @Router /auth/me [get]
func Me(db *pgxpool.Pool, jwtService *services.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return c.Status(401).JSON(schemas.ErrorResponse{Error: "unauthenticated"})
		}

		tokenUser, err := jwtService.ValidateToken(cookie)
		if err != nil {
			return c.Status(401).JSON(schemas.ErrorResponse{Error: "invalid token"})
		}

		repo := repositories.NewUserRepository(db)
		user, err := repo.GetByEmail(c.Context(), tokenUser.Email)
		if err != nil {
			return c.Status(404).JSON(schemas.ErrorResponse{Error: "user not found"})
		}

		authUser := schemas.AuthUser{
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
			Avatar:     user.Avatar,
		}

		return c.JSON(schemas.AuthResponse{
			User:  authUser,
			Token: cookie,
		})
	}
}
