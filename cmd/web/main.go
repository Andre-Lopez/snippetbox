package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const PORT = ":4000"

func main() {
	app := fiber.New()
	app.Get("/", home)
	app.Get("/snippet/view", viewSnippet)
	app.Post("/snippet/create", createSnippet)

	log.Println("Starting on server", PORT)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})

	log.Fatal(app.Listen(PORT))
}
