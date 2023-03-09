package main

import (
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/Andre-Lopez/snippetbox/internal/assert"
	"github.com/gofiber/fiber/v2"
)

func TestSnippetView(t *testing.T) {
	app := newTestApp(t)

	// TODO: user existance check from DB causes failure
	mux := app.routes()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: fiber.StatusOK,
			wantBody: "Test title",
		},
		{
			name:     "Invalid ID",
			urlPath:  "/snippet/view/2",
			wantCode: fiber.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-3",
			wantCode: fiber.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.75",
			wantCode: fiber.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/cool",
			wantCode: fiber.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(fiber.MethodGet, tt.urlPath, nil)
			if err != nil {
				t.Fatal(err)
			}

			res, _ := mux.Test(req)
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantCode)

			if tt.wantBody != "" {
				// Stringify response body
				body, err := httputil.DumpResponse(res, true)
				if err != nil {
					t.Fatal(err)
				}
				assert.StringContains(t, string(body), tt.wantBody)
			}
		})
	}
}
