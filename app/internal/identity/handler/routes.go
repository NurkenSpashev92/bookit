package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"

	"github.com/nurkenspashev92/bookit/internal/shared"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const (
	avatarUploadLimit   = 10 * 1024 * 1024
	avatarUploadTimeout = 30 * time.Second
)

type Deps struct {
	User   UserService
	Avatar AvatarService
}

func RegisterRoutes(api fiber.Router, deps Deps, guards shared.Guards) {
	authHandler := NewAuthHandler(deps.User)
	avatarHandler := NewAvatarHandler(deps.Avatar)

	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", authHandler.Me)
	auth.Patch("/me", guards.Required, authHandler.UpdateProfile)
	auth.Patch("/me/password", guards.Required, authHandler.ChangePassword)
	auth.Post("/me/avatar",
		guards.Required,
		middleware.UploadLimits(avatarUploadLimit),
		timeout.New(avatarHandler.Upload, timeout.Config{Timeout: avatarUploadTimeout}),
	)
	auth.Delete("/me/avatar", guards.Required, avatarHandler.Delete)
}
