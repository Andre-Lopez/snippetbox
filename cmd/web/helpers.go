package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           interface{}
	IsAuthenticated bool
	CSRFToken       string
}

// Return true if request is from authenticated user
func (app *application) isAuthenticated(c *fiber.Ctx) bool {
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		return false
	}

	authUserId := sess.Get("authUserId")

	return authUserId != nil
}

// Obtains flash message and deletes if it exists
func (app *application) popFlashMessage(c *fiber.Ctx) interface{} {
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		return false
	}

	flash := sess.Get("flash")

	if flash != nil {
		sess.Delete("flash")
		if err = sess.Save(); err != nil {
			app.serverError(c, err)
			return nil
		}
	}

	return flash
}

// Returns template data struct
func (app *application) newTemplateData(c *fiber.Ctx) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.popFlashMessage(c),
		IsAuthenticated: app.isAuthenticated(c),
		CSRFToken:       c.Cookies("csrf"),
	}
}

func (app *application) clientError(c *fiber.Ctx, status int) {
	c.Status(status).SendString(http.StatusText(status))
}

func (app *application) notFound(c *fiber.Ctx) {
	app.clientError(c, fiber.StatusNotFound)
}

func (app *application) serverError(c *fiber.Ctx, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
}
