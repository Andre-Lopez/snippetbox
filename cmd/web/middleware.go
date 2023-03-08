package main

// For middleware that require access to application struct

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Config for fiber middleware
type Config struct {
	Filter func(c *fiber.Ctx) bool
}

var ConfigDefault = Config{
	Filter: nil,
}

func SetConfigDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}

	return cfg
}

func (app *application) authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := app.sessionManager.Get(c)
		if err != nil {
			return c.Next()
		}

		var authUserId int
		id := sess.Get("authUserId")

		if id == nil {
			return c.Next()
		} else {
			authUserId = id.(int)
		}

		// Check if user exists
		exists, err := app.users.Exists(authUserId)
		if err != nil {
			app.serverError(c, err)
			return err
		}

		// Store bool in req context if exists
		if exists {
			c.Locals(isAuthenticatedContextKey, "true")
		}

		return c.Next()
	}
}

func (app *application) RequiresAuth(c *fiber.Ctx) error {

	// Redirect user to login page if not authenticated
	if !app.isAuthenticated(c) {
		return c.Redirect("/user/login")
	}

	return c.Next()
}

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
