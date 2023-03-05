package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetSecureHeaders(config Config) fiber.Handler {
	cfg := SetConfigDefault(config)

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			fmt.Println("secure headers middleware skipped")
			return c.Next()
		}

		c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		c.Set("Referrer-Policy", "origin-when-cross-origin")
		c.Set("X-Content-Type-Options", "deny")
		c.Set("X-XSS-Protection", "0")

		return c.Next()
	}
}
