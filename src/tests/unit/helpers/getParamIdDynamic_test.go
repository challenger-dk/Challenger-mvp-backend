package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"server/controllers/helpers"
	"testing"
)

func newRequestWithDynamicID(paramName string, paramValue string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test/"+paramValue, nil)
	req.SetPathValue(paramName, paramValue)
	return req
}

// ------------------ TESTS FOR GetParamIdDynamic ------------------

func TestGetParamIdDynamic_MissingID(t *testing.T) {
	paramName := "userId"
	r := newRequestWithDynamicID(paramName, "")

	id, err := helpers.GetParamIdDynamic(r, paramName)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_ZeroID(t *testing.T) {
	paramName := "teamId"
	r := newRequestWithDynamicID(paramName, "0")

	id, err := helpers.GetParamIdDynamic(r, paramName)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_InvalidID(t *testing.T) {
	paramName := "challengeId"
	r := newRequestWithDynamicID(paramName, "abc")

	id, err := helpers.GetParamIdDynamic(r, paramName)

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if id != 0 {
		t.Fatalf("expected ID 0, got %d", id)
	}
}

func TestGetParamIdDynamic_ValidID(t *testing.T) {
	paramName := "id"
	r := newRequestWithDynamicID(paramName, "123")

	id, err := helpers.GetParamIdDynamic(r, paramName)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 123 {
		t.Fatalf("expected ID 123, got %d", id)
	}
}
