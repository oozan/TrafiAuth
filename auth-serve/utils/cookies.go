// utils/cookies.go

package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetAuthCookies(c *fiber.Ctx, accessToken, refreshToken string) {
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Expires:  time.Now().Add(15 * time.Minute),
        HTTPOnly: true,
        Secure:   false, 
    })

    c.Cookie(&fiber.Cookie{
        Name:     "refresh_token",
        Value:    refreshToken,
        Expires:  time.Now().Add(24 * time.Hour),
        HTTPOnly: true,
        Secure:   false, 
    })
}

func SetAccessTokenCookie(c *fiber.Ctx, accessToken string) {
    c.Cookie(&fiber.Cookie{
        Name:     "access_token",
        Value:    accessToken,
        Expires:  time.Now().Add(15 * time.Minute),
        HTTPOnly: true,
        Secure:   false, 
    })
}
