package main

import (
	"errors"
	"time"

	"github.com/Andre-Lopez/snippetbox/cmd/web/middleware"
	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func (app *application) routes() *fiber.App {
	mux := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(ctx)
			} else {
				app.serverError(ctx, err)
			}
			return nil
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Views:        html.New("./ui/html", ".html"),
	})

	// Logger Middleware
	mux.Use(logger.New(logger.Config{
		Format: "[${ip}] - ${port} ${method} ${path}\n",
	}))

	// CSRF middleware
	mux.Use(csrf.New(csrf.Config{
		KeyLookup:      "form:csrf_token",
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieName:     "csrf",
	}))

	// Set Secure Headers Middleware
	mux.Use(middleware.SetSecureHeaders(middleware.Config{}))

	// Panic Recovery Middleware
	mux.Use(recover.New())

	mux.Static("/static", app.staticPath, fiber.Static{Browse: true})

	mux.Get("/", app.viewHome)
	mux.Get("/snippet/view/:id", app.viewSnippet)
	mux.Get("/snippet/create", app.RequiresAuth(), app.createSnippet)
	mux.Post("/snippet/create", app.RequiresAuth(), app.createSnippetPost)

	mux.Get("/user/signup", app.userSignup)
	mux.Post("/user/signup", app.userSignupPost)
	mux.Get("/user/login", app.userLogin)
	mux.Post("/user/login", app.userLoginPost)
	mux.Post("/user/logout", app.userLogout)

	mux.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return mux
}
