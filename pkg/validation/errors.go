package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom error types for specific error scenarios encountered during validation error mapping.
type ErrNonStructType struct {
	Type reflect.Kind
}

func (e *ErrNonStructType) Error() string {
	return fmt.Sprintf("expected a struct but got a %s", e.Type)
}

type ErrFieldNotFound struct {
	FieldName string
}

func (e *ErrFieldNotFound) Error() string {
	return fmt.Sprintf("field '%s' not found", e.FieldName)
}

// StructFieldError represents a detailed error for a specific field within a struct.
type StructFieldError struct {
	FieldName    string `json:"field"`
	ErrorMessage string `json:"message"`
}

// ApiError represents a detailed error for a specific field within a struct.
type ApiError struct {
	Message string             `json:"message"`
	Details []StructFieldError `json:"details"`
}

// FormatTagErrorMessage formats validation error messages based on the validation tag and parameter.
func FormatTagErrorMessage(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "max":
		return fmt.Sprintf("The value must be less than or equal to %s.", param)
	case "min":
		return fmt.Sprintf("The value must be greater than or equal to %s.", param)
	case "fqdn":
		return "The value must be a valid fully-qualified domain name."
	case "uuid":
		return "The value must be a valid UUID."
	default:
		return "The value provided is invalid."
	}
}

// MapStructValidationErrors maps validation errors to a human-readable format, supporting custom struct tags.
func MapStructValidationErrors(ve validator.ValidationErrors, i interface{}, customTag ...string) ([]StructFieldError, error) {
	// Default tag for identifying struct fields.
	tag := "json"
	if len(customTag) > 0 && customTag[0] != "" {
		tag = customTag[0]
	}

	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Ptr {
		return nil, &ErrNonStructType{Type: val.Kind()}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer to work with the actual struct.
	}

	errors := []StructFieldError{}
	for _, e := range ve {
		fieldPath, err := getFieldPath(e, val, tag)
		if err != nil {
			return nil, err
		}

		if fieldPath != "" {
			message := FormatTagErrorMessage(e.Tag(), e.Param())
			errors = append(errors, StructFieldError{FieldName: fieldPath, ErrorMessage: message})
		}
	}

	return errors, nil
}

// getFieldPath builds the path to a field using the provided tag.
func getFieldPath(e validator.FieldError, val reflect.Value, tag string) (string, error) {
	tagPath := []string{}
	parts := strings.Split(e.StructNamespace(), ".")[1:] // Exclude root struct name.

	for _, part := range parts {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if val.Kind() != reflect.Struct {
			return "", &ErrNonStructType{Type: val.Kind()}
		}

		field, found := val.Type().FieldByName(part)
		if !found {
			return "", &ErrFieldNotFound{FieldName: part}
		}

		fieldTag := getFieldTag(field, tag)
		if fieldTag != "" {
			tagPath = append(tagPath, fieldTag)
		}

		val = val.FieldByName(part) // Move to next nested field, if any.
	}

	return strings.Join(tagPath, "."), nil
}

// getFieldTag retrieves the specified tag for a field, defaulting to the field name if the tag is absent.
func getFieldTag(field reflect.StructField, tag string) string {
	tagValue, ok := field.Tag.Lookup(tag)
	if !ok || tagValue == "" {
		return field.Name // Default to field name if tag is missing or empty.
	}
	return strings.Split(tagValue, ",")[0] // Return first part of the tag, ignoring options.
}
