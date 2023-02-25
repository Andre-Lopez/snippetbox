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
	c.SendString("Create route from snippetbox")
	return nil
}
