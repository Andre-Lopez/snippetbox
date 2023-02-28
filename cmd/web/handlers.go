package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
	id := c.Query("id")
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

func (app *application) createSnippet(c *fiber.Ctx) error {
	body := struct {
		title   string `json:"title"`
		content string `json:"content"`
		expires int    `json:"expires"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	id, err := app.snippets.Insert(body.title, body.content, body.expires)
	if err != nil {
		app.serverError(c, err)
		return err
	}

	c.Redirect(fmt.Sprintf("/snippet/view?id=%d", id), fiber.StatusSeeOther)
	return nil
}
