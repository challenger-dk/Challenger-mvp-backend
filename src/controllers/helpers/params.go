package helpers

import (
	"fmt"
	"net/http"
	"strconv"
)

// Tries to fetch value of "id" in parameter.
// On success it returns the value.
// On failure it returns 0, so make sure to check if it is 0 after calling
func GetParamId(w http.ResponseWriter, r *http.Request) uint {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing parameter id", http.StatusBadRequest)
		return 0
	}

	if idStr == "0" {
		http.Error(w, "Id cannot be 0", http.StatusBadRequest)
		return 0
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}

	return uint(id)
}

// Tries to fetch value of parameter with given name.
// On success it returns the value.
// On failure it returns 0, so make sure to check if it is 0 after calling
func GetParamIdDynamic(w http.ResponseWriter, r *http.Request, name string) uint {
	idStr := r.PathValue(name)
	if idStr == "" {
		http.Error(w, fmt.Sprintf("Missing parameter %s", name), http.StatusBadRequest)
		return 0
	}

	if idStr == "0" {
		http.Error(w, fmt.Sprintf("Parameter %s cannot be 0", name), http.StatusBadRequest)
		return 0
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}

	return uint(id)
}
