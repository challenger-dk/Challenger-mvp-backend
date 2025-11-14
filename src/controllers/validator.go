package controllers

import (
	"server/config"

	"github.com/go-playground/validator/v10"
)

// Initialize a validator instance
// This can be reused across the application
// to validate structs based on `validate` tags
var validate = validator.New()

// Initalize the validator with custom validations
// This function is called once when the package is initialized, so dont call this manually
// Pretty cool
func init() {
	validate.RegisterValidation("is-valid-sport", validateSport)
}

func validateSport(fl validator.FieldLevel) bool {
	// Get the string value from the field
	sportName := fl.Field().String()

	// Check if the name exists in the cache map
	// The 'ok' bool will be 'true' if it's found, 'false' otherwise
	_, ok := config.SportsCache[sportName]
	return ok
}
