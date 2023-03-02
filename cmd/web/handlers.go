package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Andre-Lopez/snippetbox/internal/models"
	"github.com/Andre-Lopez/snippetbox/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type createSnippetForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

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
	expires, err := strconv.Atoi(c.FormValue("expires"))
	if err != nil {
		app.clientError(c, fiber.StatusBadRequest)
		return err
	}

	// Init form values
	form := createSnippetForm{
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
		Expires: expires,
	}

	// Run validations for each field
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be empty")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be empty")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must have a value of 1, 7, or 365")

	// Return form with data and errors if needed
	if !form.Valid() {
		return c.Render("create", fiber.Map{"title": form.Title, "content": form.Content, "expires": form.Expires, "errors": form.FieldErrors})
	}

	// Create new snippet if data is valid
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(c, err)
		return err
	}

	// Redirect user to new snippet view page
	return c.Redirect(fmt.Sprintf("/snippet/view/%d", id), fiber.StatusSeeOther)
}

func (app *application) createSnippet(c *fiber.Ctx) error {
	return c.Render("create", fiber.Map{"expires": 365})
}
