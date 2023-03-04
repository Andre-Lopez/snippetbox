package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=? ^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Returns true if not errors have been stored in FieldErrors map
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Adds error message
func (v *Validator) AddFieldError(key, message string) {
	// Init new map if needed
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Adds non-field error message (ex: invalid credentials)
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Adds error message if field is invalid
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Returns true if a string is GE a min value
func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

// Returns true if a string is LE a max value
func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// Returns true if a string matches the provided regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Returns true if a string is not empty
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Returns true if a value is in a list of ints
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
