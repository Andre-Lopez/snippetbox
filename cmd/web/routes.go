package main

import (
	"errors"

	"github.com/Andre-Lopez/snippetbox/cmd/web/middleware"
	"github.com/Andre-Lopez/snippetbox/cmd/web/middleware/secureHeaders"
	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2"
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
		Views: html.New("./ui/html", ".html"),
	})

	// Logger Middleware
	mux.Use(logger.New(logger.Config{
		Format: "[${ip}] - ${port} ${method} ${path}\n",
	}))

	// Set Secure Headers Middleware
	mux.Use(secureHeaders.New(middleware.Config{}))

	// Panic Recovery Middleware
	mux.Use(recover.New())

	mux.Static("/static", app.staticPath, fiber.Static{Browse: true})

	mux.Get("/", app.home)
	mux.Get("/snippet/view/:id", app.viewSnippet)
	mux.Get("/snippet/create", app.createSnippet)
	mux.Post("/snippet/create", app.createSnippetPost)

	mux.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return mux
}
