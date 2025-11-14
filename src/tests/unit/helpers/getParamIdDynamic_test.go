package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"server/controllers/helpers"
	"testing"
)

// Helper function to create a new request with a dynamic parameter
func newRequestWithDynamicID(paramName string, paramValue string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test/"+paramValue, nil)
	// Set the path value using the dynamic name
	req.SetPathValue(paramName, paramValue)
	return req
}

// ------------------ TESTS FOR GetParamIdDynamic ------------------

func TestGetParamIdDynamic_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	paramName := "userId" // Use a dynamic name
	r := newRequestWithDynamicID(paramName, "")

	id := helpers.GetParamIdDynamic(w, r, paramName)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_ZeroID(t *testing.T) {
	w := httptest.NewRecorder()
	paramName := "teamId" // Use a different dynamic name
	r := newRequestWithDynamicID(paramName, "0")

	id := helpers.GetParamIdDynamic(w, r, paramName)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	paramName := "challengeId" // Use a different dynamic name
	r := newRequestWithDynamicID(paramName, "abc")

	id := helpers.GetParamIdDynamic(w, r, paramName)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_ValidID(t *testing.T) {
	w := httptest.NewRecorder()
	paramName := "id" // Also test with the name "id"
	r := newRequestWithDynamicID(paramName, "123")

	id := helpers.GetParamIdDynamic(w, r, paramName)

	resp := w.Result()
	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		t.Fatalf("expected no error, got HTTP %d", resp.StatusCode)
	}
	if id != 123 {
		t.Fatalf("expected ID 123, got %d", id)
	}
}
