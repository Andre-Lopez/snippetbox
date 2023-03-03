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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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
	// Get session
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		app.clientError(c, fiber.StatusUnauthorized)
		return err
	}

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

	flash := sess.Get("flash")
	sess.Delete("flash")

	// Save session, still render template if cannot save
	if err := sess.Save(); err != nil {
		return c.Render("view", fiber.Map{"currentYear": time.Now().Year(), "snippet": snippet, "flash": flash})
	}

	return c.Render("view", fiber.Map{"currentYear": time.Now().Year(), "snippet": snippet, "flash": flash})
}

func (app *application) createSnippetPost(c *fiber.Ctx) error {
	// Get session
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		app.clientError(c, fiber.StatusUnauthorized)
		return err
	}

	var form createSnippetForm

	if err := c.BodyParser(&form); err != nil {
		app.clientError(c, fiber.StatusBadRequest)
		return err
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

	// Create entry in session to display successful toast
	sess.Set("flash", "Snippet successfully created!")

	// Save our session, still render template if cannot save
	if err := sess.Save(); err != nil {
		return c.Redirect(fmt.Sprintf("/snippet/view/%d", id), fiber.StatusSeeOther)
	}

	// Redirect user to new snippet view page
	return c.Redirect(fmt.Sprintf("/snippet/view/%d", id), fiber.StatusSeeOther)
}

func (app *application) createSnippet(c *fiber.Ctx) error {
	return c.Render("create", fiber.Map{"expires": 365})
}
