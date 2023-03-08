package main

import (
	"context"
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

func TestRequiresAuth(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want int
	}{{
		name: "Authenticated",
		val:  "true",
		want: 200,
	},
		{
			name: "Not Authenticated Empty String",
			val:  "",
			want: 302,
		},
		{
			name: "Not Authenticated, False string",
			val:  "false",
			want: 302,
		},
	}

	// Create new mux and app instance
	mux := fiber.New()
	app := &application{}

	// Create route to use with middleware
	mux.Get("/hello", app.RequiresAuth, func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Set the context value
			ctx := context.Background()
			ctx = context.WithValue(ctx, isAuthenticatedContextKey, tt.val)

			// Create request with context
			req, err := http.NewRequestWithContext(ctx, fiber.MethodGet, "/hello", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Send req
			res, err := mux.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}

			// 302 redirect if not autheticated, 200 OK if auth
			assert.Equal(t, res.StatusCode, tt.want)
		})
	}
}
