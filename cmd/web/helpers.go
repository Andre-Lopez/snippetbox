package main

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func (app *application) serverError(c *fiber.Ctx, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
}

func (app *application) clientError(c *fiber.Ctx, status int) *fiber.Ctx {
	return c.Status(status)
}

func (app *application) notFound(c *fiber.Ctx) {
	app.clientError(c, fiber.StatusNotFound)
}
