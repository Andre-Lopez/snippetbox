package main

import (
	"io"
	"log"
	"testing"

	"github.com/Andre-Lopez/snippetbox/internal/models/mocks"
)

func newTestApp(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		snippets: &mocks.SnippetModel{},
		users:    &mocks.UserModel{},
	}
}
