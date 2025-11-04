package helpers

import (
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
	}

	return uint(id)
}
