package main

import (
	"errors"
	"fmt"
	"strconv"

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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Serves the home page displaying last 10 created snippets
func (app *application) viewHome(c *fiber.Ctx) error {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(c, err)
		return err
	}

	c.Context()

	data := app.newTemplateData(c)
	data.Snippets = snippets

	return c.Render("home", *data)
}

// Handles snippet ID and serves details page of according snippet
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

	data := app.newTemplateData(c)
	data.Snippet = snippet

	return c.Render("view", *data)
}

// Handles snippet data and creates a new snippet if valid
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must have a value of 1, 7, or 365")

	// Return form with data and errors if needed
	if !form.Valid() {
		data := app.newTemplateData(c)
		data.Form = form

		return c.Render("create", *data)
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

// Serves the snippet creation form
func (app *application) createSnippet(c *fiber.Ctx) error {
	data := app.newTemplateData(c)
	data.Form = createSnippetForm{
		Expires: 365,
	}

	return c.Render("create", *data)
}

// Serves the user signup form
func (app *application) userSignup(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{})
}

// Handles signup data and registers user if valid
func (app *application) userSignupPost(c *fiber.Ctx) error {
	// Get session
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		app.clientError(c, fiber.StatusUnauthorized)
		return err
	}

	// Parse form data
	var form userSignupForm
	if err := c.BodyParser(&form); err != nil {
		app.clientError(c, fiber.StatusBadRequest)
		return err
	}

	// Run validations for each field
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "Must enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 7), "password", "Password must be at least 7 characters")

	// Return form with errors if needed
	if !form.Valid() {
		data := app.newTemplateData(c)
		data.Form = form
		return c.Render("signup", *data)
	}

	// Create new user
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			return c.Render("signup", fiber.Map{"name": form.Name, "email": form.Email, "errors": form.FieldErrors})
		} else {
			app.serverError(c, err)
			return err
		}
	}

	// Set flash message notifying success
	sess.Set("flash", "Signup successful. Please login...")

	// Save our session, still render template if cannot save
	if err := sess.Save(); err != nil {
		return c.Redirect("/user/login", fiber.StatusSeeOther)
	}

	return c.Redirect("/user/login", fiber.StatusSeeOther)
}

// Serves the user login form
func (app *application) userLogin(c *fiber.Ctx) error {
	data := app.newTemplateData(c)
	data.Form = userLoginForm{}
	return c.Render("login", *data)
}

// Handles login data and sets auth token if credentials valid
func (app *application) userLoginPost(c *fiber.Ctx) error {
	// Get session
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		app.clientError(c, fiber.StatusUnauthorized)
		return err
	}

	// Parse form data
	var form userLoginForm
	if err := c.BodyParser(&form); err != nil {
		app.clientError(c, fiber.StatusBadRequest)
		return err
	}

	// Run validations for each field
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "Must enter a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(c)
		data.Form = form
		return c.Render("login", *data)
	}

	// Check validation of credentials
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(c)
			data.Form = form
			return c.Render("login", *data)
		} else {
			app.serverError(c, err)
			return err
		}
	}

	// Set auth cookie and save session
	sess.Set("authUserId", id)
	if err := sess.Save(); err != nil {
		app.serverError(c, err)
		return err
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}

// Handles login data and sets auth token if credentials valid
func (app *application) userLogout(c *fiber.Ctx) error {
	// Get session
	sess, err := app.sessionManager.Get(c)
	if err != nil {
		app.clientError(c, fiber.StatusUnauthorized)
		return err
	}

	// Delete auth cookie
	sess.Delete("authUserId")

	// Create entry in session to display successful toast
	sess.Set("flash", "You have been successfully logged out")

	if err := sess.Save(); err != nil {
		app.serverError(c, err)
		return err
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}
