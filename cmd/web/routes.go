package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func (app *application) routes() *fiber.App {
	engine := html.New("./ui/html", ".html")
	mux := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			app.serverError(ctx, err)
			return nil
		},

		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	mux.Static("/static", app.staticPath, fiber.Static{Browse: true})

	mux.Get("/", app.home)
	mux.Get("/snippet/view", app.viewSnippet)
	mux.Post("/snippet/create", app.createSnippet)

	mux.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return mux
}
