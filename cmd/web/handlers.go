package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (app *application) home(c *fiber.Ctx) error {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(c, err)
		return err
	}

	return c.Render("home", fiber.Map{"currentYear": time.Now().Year(), "snippets": snippets})
}

func (app *application) viewSnippet(c *fiber.Ctx) error {
	id := c.Params("id")
	intId, err := strconv.Atoi(id)

	if err != nil || intId < 1 {
		app.notFound(c)
		return nil
	}

	snippet, err := app.snippets.Get(intId)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(c)
		} else {
			app.serverError(c, err)
		}
		return err
	}

	return c.Render("view", fiber.Map{"currentYear": time.Now().Year(), "snippet": snippet})
}

func (app *application) createSnippetPost(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")
	expires, err := strconv.Atoi(c.FormValue("expires"))
	if err != nil {
		app.clientError(c, fiber.StatusBadRequest)
		return err
	}

	fieldErrors := make(map[string]string)

	// Validate title
	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be empty"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be longer than 100 characters"
	}

	// Validate content
	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be empty"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must have a value of 1, 7, or 365"
	}

	// Re-serve the creation form with displayed errors
	if len(fieldErrors) > 0 {
		return c.Render("create", fiber.Map{"title": title, "content": content, "expires": expires, "errors": fieldErrors})
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(c, err)
		return err
	}

	c.Redirect(fmt.Sprintf("/snippet/view/%d", id), fiber.StatusSeeOther)
	return nil
}

func (app *application) createSnippet(c *fiber.Ctx) error {
	return c.Render("create", fiber.Map{"expires": 365})
}
