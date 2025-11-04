package helpers

import (
	"net/http"
	"strconv"
)

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
