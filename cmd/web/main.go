package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// Obtain command line args
	PORT := flag.String("port", ":4000", "HTTP Network Address")
	STATIC_PATH := flag.String("static-path", "./ui/static", "Path of static conent to serve")
	flag.Parse()

	engine := html.New("./ui/html", ".html")
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Static("/static", *STATIC_PATH, fiber.Static{Browse: true})

	app.Get("/", home)
	app.Get("/snippet/view", viewSnippet)
	app.Post("/snippet/create", createSnippet)

	log.Println("Starting on server", *PORT)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusNotFound)
	})

	log.Fatal(app.Listen(*PORT))
}
