package main

import (
	"flag"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Obtain command line args
	PORT := flag.String("port", ":4000", "HTTP Network Address")
	STATIC_PATH := flag.String("static-path", "./ui/static", "Path of static conent to serve")
	flag.Parse()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// init our app struct
	application := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	engine := html.New("./ui/html", ".html")
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			application.serverError(ctx, err)
			return nil
		},
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Static("/static", *STATIC_PATH, fiber.Static{Browse: true})

	app.Get("/", application.home)
	app.Get("/snippet/view", application.viewSnippet)
	app.Post("/snippet/create", application.createSnippet)

	infoLog.Println("Starting on server", *PORT)
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	errorLog.Fatal(app.Listen(*PORT))
}
