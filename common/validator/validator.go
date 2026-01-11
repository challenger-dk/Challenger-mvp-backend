package validator

import (
	"html"
	"log/slog"
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

	// Check if the name exists in the cache map (case-insensitive)
	// First try exact match
	if _, ok := config.SportsCache[sportName]; ok {
		slog.Info("Sport validation passed (exact match)", "sport", sportName)
		return true
	}

	// If not found, try case-insensitive match
	for cachedSport := range config.SportsCache {
		if strings.EqualFold(sportName, cachedSport) {
			slog.Info("Sport validation passed (case-insensitive match)",
				"input", sportName,
				"matched", cachedSport,
			)
			return true
		}
	}

	slog.Warn("Sport validation failed",
		"sport", sportName,
		"available_sports", getAvailableSports(),
	)
	return false
}

// getAvailableSports returns a slice of all available sports in the cache
func getAvailableSports() []string {
	sports := make([]string, 0, len(config.SportsCache))
	for sport := range config.SportsCache {
		sports = append(sports, sport)
	}
	return sports
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
