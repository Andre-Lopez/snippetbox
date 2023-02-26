package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (app *application) home(c *fiber.Ctx) error {
	if c.Path() != "/" {
		app.notFound(c)
		return nil
	}

	return c.Render("home", fiber.Map{})
}

func (app *application) viewSnippet(c *fiber.Ctx) error {
	id := c.Query("id")
	intId, err := strconv.Atoi(id)

	if err != nil || intId < 1 {
		app.notFound(c)
		return nil
	}

	c.SendString(fmt.Sprintf("Viewing id: %d", intId))
	return nil
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
