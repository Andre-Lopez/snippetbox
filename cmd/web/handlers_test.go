package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func TestUserSignup(t *testing.T) {
	app := newTestApp(t)
	mux := app.routes()

	req, err := http.NewRequest(fiber.MethodGet, "/user/signup", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, _ := mux.Test(req, -1)

	// Stringify response body
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		t.Fatal(err)
	}

	validCsrfToken := extractCSRFToken(t, string(body))

	const (
		validName     = "Tom"
		validPassword = "validPassword"
		validEmail    = "tom@test.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name        string
		userName    string
		email       string
		password    string
		csrfToken   string
		wantCode    int
		wantFormTag string
	}{
		{
			name:      "Valid submission",
			userName:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: validCsrfToken,
			wantCode:  fiber.StatusSeeOther,
		},
		{
			name:      "Invalid CSRF Token",
			userName:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: "BAD TOKEN",
			wantCode:  fiber.StatusBadRequest,
		},
		{
			name:        "Empty Name",
			userName:    "",
			email:       validEmail,
			password:    validPassword,
			csrfToken:   validCsrfToken,
			wantCode:    fiber.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			req, err := http.NewRequest(fiber.MethodPost, "/user/signup", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			res, _ := mux.Test(req)

			// TODO: determine why we get 500 res
			t.Logf("%v", res)
		})
	}

}
