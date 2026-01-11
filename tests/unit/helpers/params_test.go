package helpers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"server/api/controllers/helpers"
	"server/common/appError"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryParam(t *testing.T) {
	t.Run("Get existing query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=test&value=123",
			},
		}

		value, err := helpers.GetQueryParam(req, "name")
		assert.NoError(t, err)
		assert.Equal(t, "test", value)

		value, err = helpers.GetQueryParam(req, "value")
		assert.NoError(t, err)
		assert.Equal(t, "123", value)
	})

	t.Run("Get missing query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=test",
			},
		}

		value, err := helpers.GetQueryParam(req, "missing")
		assert.Error(t, err)
		assert.Empty(t, value)
		assert.Contains(t, err.Error(), "missing query parameter")
	})

	t.Run("Get empty query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=",
			},
		}

		value, err := helpers.GetQueryParam(req, "name")
		assert.Error(t, err)
		assert.Empty(t, value)
	})
}

func TestGetQueryParamOptional(t *testing.T) {
	t.Run("Get existing query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=test&value=123",
			},
		}

		value := helpers.GetQueryParamOptional(req, "name")
		assert.Equal(t, "test", value)

		value = helpers.GetQueryParamOptional(req, "value")
		assert.Equal(t, "123", value)
	})

	t.Run("Get missing query parameter returns empty string", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=test",
			},
		}

		value := helpers.GetQueryParamOptional(req, "missing")
		assert.Empty(t, value)
	})

	t.Run("Get empty query parameter returns empty string", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "name=",
			},
		}

		value := helpers.GetQueryParamOptional(req, "name")
		assert.Empty(t, value)
	})
}

func TestGetQueryInt(t *testing.T) {
	t.Run("Get valid integer query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "limit=10&offset=20",
			},
		}

		limit := helpers.GetQueryInt(req, "limit", 5)
		assert.Equal(t, 10, limit)

		offset := helpers.GetQueryInt(req, "offset", 0)
		assert.Equal(t, 20, offset)
	})

	t.Run("Get missing query parameter returns default", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "other=value",
			},
		}

		value := helpers.GetQueryInt(req, "limit", 25)
		assert.Equal(t, 25, value)
	})

	t.Run("Get invalid integer query parameter returns default", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "limit=invalid",
			},
		}

		value := helpers.GetQueryInt(req, "limit", 10)
		assert.Equal(t, 10, value)
	})

	t.Run("Get negative integer query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "limit=-5",
			},
		}

		value := helpers.GetQueryInt(req, "limit", 10)
		assert.Equal(t, -5, value)
	})

	t.Run("Get zero integer query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "limit=0",
			},
		}

		value := helpers.GetQueryInt(req, "limit", 10)
		assert.Equal(t, 0, value)
	})
}

func TestGetQueryUint(t *testing.T) {
	t.Run("Get valid uint query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=5&size=50",
			},
		}

		page := helpers.GetQueryUint(req, "page", 1)
		assert.Equal(t, uint(5), page)

		size := helpers.GetQueryUint(req, "size", 10)
		assert.Equal(t, uint(50), size)
	})

	t.Run("Get missing query parameter returns default", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "other=value",
			},
		}

		value := helpers.GetQueryUint(req, "page", 1)
		assert.Equal(t, uint(1), value)
	})

	t.Run("Get invalid uint query parameter returns default", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=invalid",
			},
		}

		value := helpers.GetQueryUint(req, "page", 1)
		assert.Equal(t, uint(1), value)
	})

	t.Run("Get zero uint query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=0",
			},
		}

		value := helpers.GetQueryUint(req, "page", 1)
		assert.Equal(t, uint(0), value)
	})

	t.Run("Get negative uint query parameter returns default (invalid)", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=-5",
			},
		}

		value := helpers.GetQueryUint(req, "page", 1)
		assert.Equal(t, uint(1), value)
	})

	t.Run("Get large uint query parameter", func(t *testing.T) {
		req := &http.Request{
			URL: &url.URL{
				RawQuery: "page=4294967295",
			},
		}

		value := helpers.GetQueryUint(req, "page", 1)
		// Should handle max uint32 value
		assert.Equal(t, uint(4294967295), value)
	})
}

func newRequestWithPathValue(key, value string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if value != "" {
		req.SetPathValue(key, value)
	}
	return req
}

func TestGetParamId(t *testing.T) {
	t.Run("Get valid ID parameter", func(t *testing.T) {
		req := newRequestWithPathValue("id", "123")
		id, err := helpers.GetParamId(req)
		assert.NoError(t, err)
		assert.Equal(t, uint(123), id)
	})

	t.Run("Get missing ID parameter", func(t *testing.T) {
		req := newRequestWithPathValue("id", "")
		id, err := helpers.GetParamId(req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, appError.ErrMissingIdParam)
		assert.Equal(t, uint(0), id)
	})

	t.Run("Get zero ID parameter", func(t *testing.T) {
		req := newRequestWithPathValue("id", "0")
		id, err := helpers.GetParamId(req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, appError.ErrIdBelowOne)
		assert.Equal(t, uint(0), id)
	})

	t.Run("Get invalid ID parameter", func(t *testing.T) {
		req := newRequestWithPathValue("id", "invalid")
		id, err := helpers.GetParamId(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid id")
		assert.Equal(t, uint(0), id)
	})
}

func TestGetParamIdDynamic(t *testing.T) {
	t.Run("Get valid dynamic parameter", func(t *testing.T) {
		req := newRequestWithPathValue("team_id", "456")
		id, err := helpers.GetParamIdDynamic(req, "team_id")
		assert.NoError(t, err)
		assert.Equal(t, uint(456), id)
	})

	t.Run("Get missing dynamic parameter", func(t *testing.T) {
		req := newRequestWithPathValue("team_id", "")
		id, err := helpers.GetParamIdDynamic(req, "team_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing parameter")
		assert.Equal(t, uint(0), id)
	})

	t.Run("Get zero dynamic parameter", func(t *testing.T) {
		req := newRequestWithPathValue("team_id", "0")
		id, err := helpers.GetParamIdDynamic(req, "team_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be 0")
		assert.Equal(t, uint(0), id)
	})

	t.Run("Get invalid dynamic parameter", func(t *testing.T) {
		req := newRequestWithPathValue("team_id", "abc")
		id, err := helpers.GetParamIdDynamic(req, "team_id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parameter")
		assert.Equal(t, uint(0), id)
	})
}
