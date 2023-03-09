package main

import (
	"html"
	"io"
	"log"
	"regexp"
	"testing"

	"github.com/Andre-Lopez/snippetbox/internal/models/mocks"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Creates a new app struct for testing
func newTestApp(t *testing.T) *application {

	store := session.New(session.Config{
		CookieSecure: true,
	})

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		sessionManager: store,
	}
}

// Extracts the CSRF token from a stringified response body
func extractCSRFToken(t *testing.T, body string) string {
	left := "Set-Cookie: csrf="
	right := ";"
	rx := regexp.MustCompile(regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))
	matches := rx.FindAllStringSubmatch(body, -1)

	if len(matches) < 1 {
		t.Fatal("No CSRF token found in body")
	}

	return html.UnescapeString(string(matches[0][1]))
}

// func PostForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {

// }
