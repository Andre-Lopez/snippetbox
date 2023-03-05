package main

// For middleware that require access to application struct

import (
	"github.com/gofiber/fiber/v2"
)

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

func (app *application) RequiresAuth() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// Redirect user to login page if not authenticated
		if !app.isAuthenticated(c) {
			return c.Redirect("/user/login")
		}

		return c.Next()
	}
}
