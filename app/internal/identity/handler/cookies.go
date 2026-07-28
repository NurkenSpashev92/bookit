package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	accessCookieName  = "access_token"
	refreshCookieName = "refresh_token"
	accessCookiePath  = "/"
	refreshCookiePath = "/api/v1/auth"
	accessCookieTTL   = 15 * time.Minute
	refreshCookieTTL  = 7 * 24 * time.Hour
)

func issueAuthCookies(c fiber.Ctx, accessToken, refreshToken string) {
	now := time.Now()
	setAuthCookie(c, accessCookieName, accessToken, accessCookiePath, now.Add(accessCookieTTL))
	setAuthCookie(c, refreshCookieName, refreshToken, refreshCookiePath, now.Add(refreshCookieTTL))
}

func clearAuthCookies(c fiber.Ctx) {
	expired := time.Now().Add(-time.Hour)
	setAuthCookie(c, accessCookieName, "", accessCookiePath, expired)
	setAuthCookie(c, refreshCookieName, "", refreshCookiePath, expired)
}

func setAuthCookie(c fiber.Ctx, name, value, path string, expires time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     path,
		Expires:  expires,
	})
}
