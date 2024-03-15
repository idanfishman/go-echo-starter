package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	_RequestIDHeader = regexp.MustCompile(`^X-(?:[a-zA-Z0-9-]+-)?Request-ID$`)
)

func isRequestIDHeader(fl validator.FieldLevel) bool {
	return _RequestIDHeader.MatchString(fl.Field().String())
}

// Validator wraps the go-playground validator to provide custom validation functionalities.
type Validator struct {
	V *validator.Validate
}

// NewValidator initializes a new Validator instance with custom settings if needed.
// It returns a pointer to a Validator.
func NewValidator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterValidation("request_id_header", isRequestIDHeader)

	return &Validator{V: v}
}
