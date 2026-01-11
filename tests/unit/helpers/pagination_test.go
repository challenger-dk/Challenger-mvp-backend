package helpers_test

import (
	"server/api/controllers/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeCursor(t *testing.T) {
	t.Run("Encode string cursor", func(t *testing.T) {
		cursor := "test-cursor-value"
		encoded, err := helpers.EncodeCursor(cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, encoded)
	})

	t.Run("Encode integer cursor", func(t *testing.T) {
		cursor := 12345
		encoded, err := helpers.EncodeCursor(cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, encoded)
	})

	t.Run("Encode struct cursor", func(t *testing.T) {
		type Cursor struct {
			ID        uint   `json:"id"`
			CreatedAt string `json:"created_at"`
		}
		cursor := Cursor{ID: 42, CreatedAt: "2024-01-01"}
		encoded, err := helpers.EncodeCursor(cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, encoded)
	})

	t.Run("Encode nil value", func(t *testing.T) {
		var cursor *string = nil
		encoded, err := helpers.EncodeCursor(cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, encoded)
	})
}

func TestDecodeCursor(t *testing.T) {
	t.Run("Decode string cursor", func(t *testing.T) {
		original := "test-cursor-value"
		encoded, err := helpers.EncodeCursor(original)
		assert.NoError(t, err)

		var decoded string
		err = helpers.DecodeCursor(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("Decode integer cursor", func(t *testing.T) {
		original := 12345
		encoded, err := helpers.EncodeCursor(original)
		assert.NoError(t, err)

		var decoded int
		err = helpers.DecodeCursor(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("Decode struct cursor", func(t *testing.T) {
		type Cursor struct {
			ID        uint   `json:"id"`
			CreatedAt string `json:"created_at"`
		}
		original := Cursor{ID: 42, CreatedAt: "2024-01-01"}
		encoded, err := helpers.EncodeCursor(original)
		assert.NoError(t, err)

		var decoded Cursor
		err = helpers.DecodeCursor(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original.ID, decoded.ID)
		assert.Equal(t, original.CreatedAt, decoded.CreatedAt)
	})

	t.Run("Decode invalid base64", func(t *testing.T) {
		var decoded string
		err := helpers.DecodeCursor("invalid-base64!!!", &decoded)
		assert.Error(t, err)
	})

	t.Run("Decode invalid JSON", func(t *testing.T) {
		// Create valid base64 but invalid JSON
		invalidJSON := "dGVzdA==" // base64 of "test" but not valid JSON for our struct
		var decoded struct {
			ID uint `json:"id"`
		}
		err := helpers.DecodeCursor(invalidJSON, &decoded)
		assert.Error(t, err)
	})
}

func TestEncodeDecodeCursor_RoundTrip(t *testing.T) {
	t.Run("Round trip with uint cursor", func(t *testing.T) {
		original := uint(999)
		encoded, err := helpers.EncodeCursor(original)
		assert.NoError(t, err)

		var decoded uint
		err = helpers.DecodeCursor(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("Round trip with complex struct", func(t *testing.T) {
		type Cursor struct {
			ID        uint   `json:"id"`
			Timestamp int64  `json:"timestamp"`
			Type      string `json:"type"`
		}
		original := Cursor{ID: 123, Timestamp: 1704067200, Type: "challenge"}
		encoded, err := helpers.EncodeCursor(original)
		assert.NoError(t, err)

		var decoded Cursor
		err = helpers.DecodeCursor(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})
}
