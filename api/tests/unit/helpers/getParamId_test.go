package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"server/api/controllers/helpers"
	"testing"
)

func newRequestWithID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test/"+id, nil)
	req.SetPathValue("id", id)
	return req
}

// ------------------ TESTS ------------------
func TestGetParamId_MissingID(t *testing.T) {
	r := newRequestWithID("")

	id, err := helpers.GetParamId(r)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_ZeroID(t *testing.T) {
	r := newRequestWithID("0")

	id, err := helpers.GetParamId(r)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_InvalidID(t *testing.T) {
	r := newRequestWithID("abc")

	id, err := helpers.GetParamId(r)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamId_ValidID(t *testing.T) {
	r := newRequestWithID("123")

	id, err := helpers.GetParamId(r)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 123 {
		t.Fatalf("expected ID 123, got %d", id)
	}
}
