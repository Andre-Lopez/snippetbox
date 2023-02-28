package secureHeaders

import (
	"fmt"

	"github.com/Andre-Lopez/snippetbox/cmd/web/middleware"
	"github.com/gofiber/fiber/v2"
)

func configDefault(config ...middleware.Config) middleware.Config {
	if len(config) < 1 {
		return middleware.ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = middleware.ConfigDefault.Filter
	}

	return cfg
}

func New(config middleware.Config) fiber.Handler {
	cfg := configDefault(config)

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
