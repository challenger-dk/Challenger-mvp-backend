package validator

import (
	"html"
	"reflect"
	"server/common/config"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Initialize a validator instance
// This can be reused across the application
// to validate structs based on `validate` tags
var V = validator.New() // Struct is the global validator instance

// Initalize the validator with custom validations
// This function is called once when the package is initialized, so dont call this manually
// Pretty cool
func init() {
	// Add custom validators here
	err := V.RegisterValidation("is-valid-sport", validateSport)
	if err != nil {
		panic("Failed to register custom validation: " + err.Error())
	}

	// Register sanitize validation that automatically cleans strings
	err = V.RegisterValidation("sanitize", sanitizeAndValidate)
	if err != nil {
		panic("Failed to register sanitize validation: " + err.Error())
	}

	V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func validateSport(fl validator.FieldLevel) bool {
	// Get the string value from the field
	sportName := fl.Field().String()

	// Check if the name exists in the cache map
	// The 'ok' bool will be 'true' if it's found, 'false' otherwise
	_, ok := config.SportsCache[sportName]
	return ok
}

func sanitizeAndValidate(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Only process string fields
	if field.Kind() != reflect.String {
		return true
	}

	// Get the original string
	original := field.String()

	// Sanitize it
	sanitized := sanitizeStrings(original)

	// Set the sanitized value back to the field
	if field.CanSet() {
		field.SetString(sanitized)
	}

	return true // Always return true since we're sanitizing, not validating
}

func sanitizeStrings(s string) string {
	trimmed := strings.TrimSpace(s)
	return html.EscapeString(trimmed)
}
