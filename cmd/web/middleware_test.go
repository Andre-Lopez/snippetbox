package main

import (
	"net/http"
	"testing"

	"github.com/Andre-Lopez/snippetbox/internal/assert"
	"github.com/gofiber/fiber/v2"
)

func TestSetSecureHeaders(t *testing.T) {
	// Create new app
	app := fiber.New()
	app.Use(SetSecureHeaders(Config{}))

	// Create route to use with middleware
	app.Get("/hello", func(c *fiber.Ctx) error {
		// Return simple string as response
		return c.SendString("Hello, World!")
	})

	// Create test request
	req, err := http.NewRequest(fiber.MethodGet, "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Send req
	res, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}

	expected := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expected)
}
