package controllers

import (
	"reflect"
	"server/common/config"
	"strings"

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
	// Add custom validators here
	err := validate.RegisterValidation("is-valid-sport", validateSport)
	if err != nil {
		panic("Failed to register custom validation: " + err.Error())
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
