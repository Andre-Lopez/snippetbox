package main

// For middleware that require access to application struct

import (
	"github.com/gofiber/fiber/v2"
)

func (app *application) RequiresAuth() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// Redirect user to login page if not authenticated
		if !app.isAuthenticated(c) {
			return c.Redirect("/user/login")
		}

		return c.Next()
	}
}
