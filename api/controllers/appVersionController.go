package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type VersionCheckResponse struct {
	UpdateRequired bool   `json:"update_required"`
	MinimumVersion string `json:"minimum_version,omitempty"`
	CurrentVersion string `json:"current_version,omitempty"`
}

// CheckAppVersion is a public endpoint that mobile apps call to verify their
// version meets the minimum required. Pass the app version via the X-App-Version header.
// When update_required is true, the app should prompt the user to update before continuing.
func CheckAppVersion(w http.ResponseWriter, r *http.Request) {
	currentVersion := r.Header.Get("X-App-Version")
	if currentVersion == "" {
		currentVersion = r.URL.Query().Get("version")
	}
	if currentVersion == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing X-App-Version header or version query parameter",
		})
		return
	}

	minimumVersion := "1.0.13" // Update this when you need to force users to update
	updateRequired := !versionAtLeast(currentVersion, minimumVersion)
	json.NewEncoder(w).Encode(VersionCheckResponse{
		UpdateRequired: updateRequired,
		MinimumVersion: minimumVersion,
		CurrentVersion: currentVersion,
	})
}

// versionAtLeast returns true if v >= minimum (v meets or exceeds the minimum required version).
func versionAtLeast(v string, minimum string) bool {
	vParts := parseVersion(v)
	minParts := parseVersion(minimum)

	for i := 0; i < 3; i++ {
		vPart := 0
		minPart := 0
		if i < len(vParts) {
			vPart = vParts[i]
		}
		if i < len(minParts) {
			minPart = minParts[i]
		}
		if vPart > minPart {
			return true
		}
		if vPart < minPart {
			return false
		}
	}
	return true
}

func parseVersion(s string) []int {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	parts := strings.Split(s, ".")
	result := make([]int, 0, 3)
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if idx := strings.IndexAny(p, "-+"); idx >= 0 {
			p = p[:idx]
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			n = 0
		}
		result = append(result, n)
	}
	return result
}
