package main

import (
	"io"
	"log"
	"testing"

	"github.com/Andre-Lopez/snippetbox/internal/models/mocks"
	"github.com/gofiber/fiber/v2/middleware/session"
)

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
