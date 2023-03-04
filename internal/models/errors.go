package models

import (
	"errors"
)

var ErrDuplicateEmail = errors.New("models: Duplicate email")
var ErrInvalidCredentials = errors.New("models: Invalid credentials")
var ErrNoRecord = errors.New("models: No matching record found")
