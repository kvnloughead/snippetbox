package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Returns true if there are no input field errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

/*
Add an error to the validator's FieldErrors struct, unless the field in question already has an error.

The struct will be initialized if it hasn't been already.
*/
func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = message
	}
}

/*
Adds an error to the validator's FieldErrors struct if the field isn't valid.

'ok' should be true if the field is valid, otherwise false. 'field' is the name of the input field. 'message' is the associated error message.
*/
func (v *Validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}

// Returns true if the string is not an empty string.
func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

// Returns true if a value contains no more than n characters.
func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) < n
}

// Returns true if the value matches one of the permittedValues.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
