package helpers

import (
	"fmt"
	"net/http"
	"server/appError"
	"strconv"
)

// Tries to fetch value of "id" in parameter.
// On success it returns the value.
// On failure it returns 0 and err
func GetParamId(r *http.Request) (uint, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, appError.ErrMissingIdParam
	}

	if idStr == "0" {
		return 0, appError.ErrIdBelowOne
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return uint(id), nil
}

// Tries to fetch value of parameter with given name.
// On success it returns the value.
// On failure it returns 0 and err
func GetParamIdDynamic(r *http.Request, name string) (uint, error) {
	idStr := r.PathValue(name)
	if idStr == "" {
		return 0, fmt.Errorf("missing parameter %s", name)
	}

	if idStr == "0" {
		return 0, fmt.Errorf("parameter %s cannot be 0", name)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid parameter %s: %w", name, err)
	}

	return uint(id), nil
}
