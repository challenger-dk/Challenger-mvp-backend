package validator

import (
	"testing"

	"server/common/config"
	"server/common/validator"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize sports cache for unit tests
	config.SportsCache = map[string]bool{
		"Tennis":     true,
		"Football":   true,
		"Basketball": true,
	}
}

type TestStruct struct {
	Name        string `validate:"sanitize"`
	Email       string `validate:"sanitize,required,email"`
	Description string `validate:"sanitize"`
	Age         int    `validate:"min=1"`
	Sport       string `validate:"is-valid-sport"`
}

func TestSanitizeValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    TestStruct
		expected TestStruct
	}{
		{
			name: "Trims whitespace",
			input: TestStruct{
				Name:        "  John Doe  ",
				Email:       "  test@example.com  ",
				Description: "  Some description  ",
				Age:         25,
				Sport:       "Tennis",
			},
			expected: TestStruct{
				Name:        "John Doe",
				Email:       "test@example.com",
				Description: "Some description",
				Age:         25,
				Sport:       "Tennis",
			},
		},
		{
			name: "Escapes HTML",
			input: TestStruct{
				Name:        "<script>alert('xss')</script>",
				Email:       "test@example.com",
				Description: "<b>Bold</b> & <i>italic</i>",
				Age:         30,
				Sport:       "Football",
			},
			expected: TestStruct{
				Name:        "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
				Email:       "test@example.com",
				Description: "&lt;b&gt;Bold&lt;/b&gt; &amp; &lt;i&gt;italic&lt;/i&gt;",
				Age:         30,
				Sport:       "Football",
			},
		},
		{
			name: "Handles empty strings",
			input: TestStruct{
				Name:        "",
				Email:       "empty@test.com",
				Description: "   ",
				Age:         18,
				Sport:       "Basketball",
			},
			expected: TestStruct{
				Name:        "",
				Email:       "empty@test.com",
				Description: "",
				Age:         18,
				Sport:       "Basketball",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate (which should sanitize)
			err := validator.V.Struct(&tt.input)
			assert.NoError(t, err)

			// Check if sanitization occurred
			assert.Equal(t, tt.expected.Name, tt.input.Name)
			assert.Equal(t, tt.expected.Email, tt.input.Email)
			assert.Equal(t, tt.expected.Description, tt.input.Description)
			assert.Equal(t, tt.expected.Age, tt.input.Age)
			assert.Equal(t, tt.expected.Sport, tt.input.Sport)
		})
	}
}

func TestIsValidSportValidation(t *testing.T) {
	tests := []struct {
		name       string
		sport      string
		shouldPass bool
	}{
		{
			name:       "Valid sport - Tennis",
			sport:      "Tennis",
			shouldPass: true,
		},
		{
			name:       "Valid sport - Football",
			sport:      "Football",
			shouldPass: true,
		},
		{
			name:       "Valid sport - Basketball",
			sport:      "Basketball",
			shouldPass: true,
		},
		{
			name:       "Invalid sport",
			sport:      "InvalidSport",
			shouldPass: false,
		},
		{
			name:       "Empty sport",
			sport:      "",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := TestStruct{
				Name:        "Test",
				Email:       "test@example.com",
				Description: "Test desc",
				Age:         25,
				Sport:       tt.sport,
			}

			err := validator.V.Struct(&testStruct)

			if tt.shouldPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "is-valid-sport")
			}
		})
	}
}

func TestSanitizeWithUserDto(t *testing.T) {
	// Test with actual DTO from your project
	userDto := struct {
		FirstName string `json:"first_name" validate:"sanitize,required,min=3"`
		LastName  string `json:"last_name" validate:"sanitize"`
		Bio       string `json:"bio" validate:"sanitize"`
	}{
		FirstName: "  <script>John</script>  ",
		LastName:  "  Doe & Co  ",
		Bio:       "<p>Hello world!</p>",
	}

	err := validator.V.Struct(&userDto)
	assert.NoError(t, err)

	assert.Equal(t, "&lt;script&gt;John&lt;/script&gt;", userDto.FirstName)
	assert.Equal(t, "Doe &amp; Co", userDto.LastName)
	assert.Equal(t, "&lt;p&gt;Hello world!&lt;/p&gt;", userDto.Bio)
}

func TestCombinedValidations(t *testing.T) {
	// Test struct with multiple validation tags including custom ones
	testStruct := struct {
		Name  string `validate:"sanitize,required,min=3"`
		Sport string `validate:"required,is-valid-sport"`
		Email string `validate:"sanitize,required,email"`
	}{
		Name:  "  <b>John</b>  ",
		Sport: "Tennis",
		Email: "  test@example.com  ",
	}

	err := validator.V.Struct(&testStruct)
	assert.NoError(t, err)

	// Check sanitization occurred
	assert.Equal(t, "&lt;b&gt;John&lt;/b&gt;", testStruct.Name)
	assert.Equal(t, "test@example.com", testStruct.Email)
	// Sport should remain unchanged (no sanitize tag)
	assert.Equal(t, "Tennis", testStruct.Sport)
}
