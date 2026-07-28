package router

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/timeout"

	identityh "github.com/nurkenspashev92/bookit/internal/identity/handler"
	"github.com/nurkenspashev92/bookit/pkg/middleware"
)

const (
	avatarUploadLimit   = 10 * 1024 * 1024
	avatarUploadTimeout = 30 * time.Second
)

func registerIdentityRoutes(api fiber.Router, svc *Services, access guards) {
	authHandler := identityh.NewAuthHandler(svc.User)
	avatarHandler := identityh.NewAvatarHandler(svc.Avatar)

	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", authHandler.Me)
	auth.Patch("/me", access.required, authHandler.UpdateProfile)
	auth.Patch("/me/password", access.required, authHandler.ChangePassword)
	auth.Post("/me/avatar",
		access.required,
		middleware.UploadLimits(avatarUploadLimit),
		timeout.New(avatarHandler.Upload, timeout.Config{Timeout: avatarUploadTimeout}),
	)
	auth.Delete("/me/avatar", access.required, avatarHandler.Delete)
}
