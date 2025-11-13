package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"server/controllers/helpers"
	"testing"
)

func newRequestWithID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test/"+id, nil)
	req.SetPathValue("id", id)
	return req
}

// ------------------ TESTS ------------------
func TestGetParamId_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	r := newRequestWithID("")

	id := helpers.GetParamId(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_ZeroID(t *testing.T) {
	w := httptest.NewRecorder()
	r := newRequestWithID("0")

	id := helpers.GetParamId(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_InvalidID(t *testing.T) {
	w := httptest.NewRecorder()
	r := newRequestWithID("abc")

	id := helpers.GetParamId(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_ValidID(t *testing.T) {
	w := httptest.NewRecorder()
	r := newRequestWithID("123")

	id := helpers.GetParamId(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 && resp.StatusCode != 0 {
		t.Fatalf("expected no error, got HTTP %d", resp.StatusCode)
	}
	if id != 123 {
		t.Fatalf("expected ID 123, got %d", id)
	}
}
