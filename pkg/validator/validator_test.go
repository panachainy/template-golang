package validator

import (
	"strings"
	"sync"
	"testing"

	pkgErrors "template-golang/pkg/errors"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=0,max=120"`
	Password string `json:"password" validate:"password_strength"`
	Phone    string `json:"phone" validate:"phone"`
	Username string `json:"username" validate:"username"`
	Slug     string `json:"slug" validate:"slug"`
	NoSpace  string `json:"no_space" validate:"no_spaces"`
}

func TestNew(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.NotNil(t, v.validate)
}

func TestValidator_Validate_Success(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	testData := TestStruct{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Password: "Password123!",
		Phone:    "+1234567890",
		Username: "john_doe",
		Slug:     "john-doe",
		NoSpace:  "nospaces",
	}

	err = v.Validate(testData)
	assert.NoError(t, err)
}

func TestValidator_Validate_Required(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	testData := TestStruct{
		// Name is missing (required)
		Email: "john@example.com",
		Age:   25,
	}

	err = v.Validate(testData)
	assert.Error(t, err)

	errorList, ok := err.(*pkgErrors.ErrorList)
	assert.True(t, ok)
	assert.True(t, errorList.HasErrors())

	// Should have error for missing name
	found := false
	for _, e := range errorList.Errors {
		if e.Context["field"] == "name" {
			found = true
			assert.Contains(t, e.Message, "required")
		}
	}
	assert.True(t, found, "Should have error for missing name field")
}

func TestValidator_Validate_Email(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{"valid email", "test@example.com", true},
		{"invalid email", "invalid-email", false},
		{"empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData := TestStruct{
				Name:  "John Doe",
				Email: tt.email,
				Age:   25,
			}

			err := v.Validate(testData)
			if tt.valid {
				// Should not have email validation error
				if err != nil {
					errorList := err.(*pkgErrors.ErrorList)
					for _, e := range errorList.Errors {
						assert.NotEqual(t, "email", e.Context["field"])
					}
				}
			} else {
				assert.Error(t, err)
				errorList := err.(*pkgErrors.ErrorList)
				found := false
				for _, e := range errorList.Errors {
					if e.Context["field"] == "email" {
						found = true
						break
					}
				}
				assert.True(t, found, "Should have email validation error")
			}
		})
	}
}

func TestValidator_Validate_MinMax(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"valid length", "John", true},
		{"too short", "J", false},
		{"too long", strings.Repeat("a", 51), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData := TestStruct{
				Name:  tt.value,
				Email: "test@example.com",
				Age:   25,
			}

			err := v.Validate(testData)
			if tt.valid {
				// Should not have name validation error
				if err != nil {
					errorList := err.(*pkgErrors.ErrorList)
					for _, e := range errorList.Errors {
						assert.NotEqual(t, "name", e.Context["field"])
					}
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidator_ValidateVar(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	// Test valid email
	err = v.ValidateVar("test@example.com", "email")
	assert.NoError(t, err)

	// Test invalid email
	err = v.ValidateVar("invalid-email", "email")
	assert.Error(t, err)
}

func TestCustomValidators(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	tests := []struct {
		name      string
		value     string
		tag       string
		shouldErr bool
	}{
		// Password strength tests
		{"strong password", "Password123!", "password_strength", false},
		{"weak password", "password", "password_strength", true},
		{"short password", "Pass1!", "password_strength", true},
		{"no uppercase", "password123!", "password_strength", true},
		{"no lowercase", "PASSWORD123!", "password_strength", true},
		{"no number", "Password!", "password_strength", true},
		{"no special", "Password123", "password_strength", true},

		// Phone tests
		{"valid phone", "+1234567890", "phone", false},
		{"valid phone without plus", "1234567890", "phone", false},
		{"invalid phone", "abc123", "phone", true},
		{"too short phone", "123", "phone", true},

		// Slug tests
		{"valid slug", "hello-world", "slug", false},
		{"invalid slug with spaces", "hello world", "slug", true},
		{"invalid slug with uppercase", "Hello-World", "slug", true},
		{"invalid slug with special chars", "hello_world!", "slug", true},

		// No spaces tests
		{"no spaces valid", "nospaces", "no_spaces", false},
		{"no spaces invalid", "has spaces", "no_spaces", true},

		// Username tests
		{"valid username", "john_doe123", "username", false},
		{"username too short", "jo", "username", true},
		{"username too long", "a_very_long_username_that_exceeds_limit", "username", true},
		{"username with spaces", "john doe", "username", true},
		{"username with special chars", "john@doe", "username", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateVar(tt.value, tt.tag)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGlobalValidator(t *testing.T) {
	// Save original global validator
	originalValidator := globalValidator

	// Reset to ensure clean state
	globalValidator = nil
	globalOnce = sync.Once{}

	// Test getting global validator
	v1 := GetGlobalValidator()
	v2 := GetGlobalValidator()
	assert.Equal(t, v1, v2) // Should be the same instance

	// Test setting global validator
	newValidator, err := New()
	assert.NoError(t, err)

	// Reset once to allow setting new validator
	globalValidator = nil
	globalOnce = sync.Once{}
	SetGlobalValidator(newValidator)
	v3 := GetGlobalValidator()
	assert.Equal(t, newValidator, v3)

	// Restore original state
	globalValidator = originalValidator
	globalOnce = sync.Once{}
}

func TestValidate_Global(t *testing.T) {
	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	_ = Validate(testData)
	// May have errors due to other required fields, but should not panic
	assert.NotPanics(t, func() {
		Validate(testData)
	})
}

func TestValidateVar_Global(t *testing.T) {
	err := ValidateVar("test@example.com", "email")
	assert.NoError(t, err)

	err = ValidateVar("invalid", "email")
	assert.Error(t, err)
}

func TestValidateStruct(t *testing.T) {
	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
	}

	errorList := ValidateStruct(testData)
	// May have errors, but should return ErrorList or nil
	if errorList != nil {
		assert.IsType(t, &pkgErrors.ErrorList{}, errorList)
	}
}

func TestMustValidate(t *testing.T) {
	validData := TestStruct{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      25,
		Password: "Password123!",
		Phone:    "+1234567890",
		Username: "john_doe",
		Slug:     "john-doe",
		NoSpace:  "nospaces",
	}

	assert.NotPanics(t, func() {
		MustValidate(validData)
	})

	invalidData := TestStruct{
		// Missing required fields
	}

	assert.Panics(t, func() {
		MustValidate(invalidData)
	})
}

func TestConvenienceFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() bool
		expected bool
	}{
		{"valid email", func() bool { return IsValidEmail("test@example.com") }, true},
		{"invalid email", func() bool { return IsValidEmail("invalid") }, false},
		{"valid URL", func() bool { return IsValidURL("https://example.com") }, true},
		{"invalid URL", func() bool { return IsValidURL("not-a-url") }, false},
		{"valid UUID", func() bool { return IsValidUUID("550e8400-e29b-41d4-a716-446655440000") }, true},
		{"invalid UUID", func() bool { return IsValidUUID("not-a-uuid") }, false},
		{"valid phone", func() bool { return IsValidPhone("+1234567890") }, true},
		{"invalid phone", func() bool { return IsValidPhone("abc") }, false},
		{"valid slug", func() bool { return IsValidSlug("hello-world") }, true},
		{"invalid slug", func() bool { return IsValidSlug("Hello World") }, false},
		{"valid username", func() bool { return IsValidUsername("john_doe") }, true},
		{"invalid username", func() bool { return IsValidUsername("j") }, false},
		{"strong password", func() bool { return IsStrongPassword("Password123!") }, true},
		{"weak password", func() bool { return IsStrongPassword("password") }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	v, err := New()
	assert.NoError(t, err)

	testData := struct {
		Email string `validate:"email"`
		Min   string `validate:"min=5"`
		Max   string `validate:"max=10"`
	}{
		Email: "invalid",
		Min:   "abc",
		Max:   "this is too long",
	}

	err = v.Validate(testData)
	assert.Error(t, err)

	errorList := err.(*pkgErrors.ErrorList)
	assert.True(t, errorList.HasErrors())

	// Check that error messages are human-readable
	for _, e := range errorList.Errors {
		assert.NotEmpty(t, e.Message)
		assert.NotContains(t, e.Message, "validation")
	}
}
