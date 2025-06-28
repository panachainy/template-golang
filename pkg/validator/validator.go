package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	pkgErrors "template-golang/pkg/errors"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground/validator with additional functionality
type Validator struct {
	validate *validator.Validate
}

// ValidationError represents a single validation error with context
type ValidationError struct {
	Field   string      `json:"field"`
	Tag     string      `json:"tag"`
	Value   interface{} `json:"value,omitempty"`
	Message string      `json:"message"`
	Param   string      `json:"param,omitempty"`
}

// New creates a new validator instance
func New() (*Validator, error) {
	validate := validator.New()

	// Register custom tag name function to use JSON tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	if err := registerCustomValidators(validate); err != nil {
		return nil, fmt.Errorf("failed to register custom validators: %w", err)
	}

	return &Validator{
		validate: validate,
	}, nil
}

// Validate validates a struct and returns detailed error information
func (v *Validator) Validate(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		var validationErrors []ValidationError

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: getErrorMessage(err),
				Param:   err.Param(),
			})
		}

		// Create error list
		errorList := pkgErrors.NewErrorList()
		for _, vErr := range validationErrors {
			errorList.AddValidation(vErr.Field, vErr.Message)
		}

		return errorList
	}

	return nil
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// RegisterValidation registers a custom validation function
func (v *Validator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return v.validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

// RegisterAlias registers an alias for a validation tag
func (v *Validator) RegisterAlias(alias, tags string) {
	v.validate.RegisterAlias(alias, tags)
}

// getErrorMessage returns a human-readable error message for validation errors
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Must be at least %s characters long", fe.Param())
		}
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Must be at most %s characters long", fe.Param())
		}
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "len":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Must be exactly %s characters long", fe.Param())
		}
		return fmt.Sprintf("Must be exactly %s", fe.Param())
	case "alpha":
		return "Must contain only alphabetic characters"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must be a valid number"
	case "url":
		return "Must be a valid URL"
	case "uuid":
		return "Must be a valid UUID"
	case "uuid4":
		return "Must be a valid UUID v4"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "eqfield":
		return fmt.Sprintf("Must be equal to %s", fe.Param())
	case "nefield":
		return fmt.Sprintf("Must not be equal to %s", fe.Param())
	case "password_strength":
		return "Password must contain at least 8 characters with uppercase, lowercase, number and special character"
	case "phone":
		return "Must be a valid phone number"
	case "slug":
		return "Must be a valid slug (alphanumeric and hyphens only)"
	case "no_spaces":
		return "Must not contain spaces"
	case "username":
		return "Must be a valid username (3-30 characters, alphanumeric and underscores only)"
	default:
		return fmt.Sprintf("Validation failed for '%s'", fe.Tag())
	}
}

// registerCustomValidators registers custom validation functions
func registerCustomValidators(validate *validator.Validate) error {
	// Password strength validator
	if err := validate.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`\d`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

		return hasUpper && hasLower && hasNumber && hasSpecial
	}); err != nil {
		return fmt.Errorf("failed to register password_strength validator: %w", err)
	}

	// Phone number validator (basic)
	if err := validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
		return phoneRegex.MatchString(phone)
	}); err != nil {
		return fmt.Errorf("failed to register phone validator: %w", err)
	}

	// Slug validator
	if err := validate.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		slug := fl.Field().String()
		slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		return slugRegex.MatchString(slug)
	}); err != nil {
		return fmt.Errorf("failed to register slug validator: %w", err)
	}

	// No spaces validator
	if err := validate.RegisterValidation("no_spaces", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return !strings.Contains(value, " ")
	}); err != nil {
		return fmt.Errorf("failed to register no_spaces validator: %w", err)
	}

	// Username validator
	if err := validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		if len(username) < 3 || len(username) > 30 {
			return false
		}
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return usernameRegex.MatchString(username)
	}); err != nil {
		return fmt.Errorf("failed to register username validator: %w", err)
	}

	return nil
}

// Global validator instance
var globalValidator *Validator

// SetGlobalValidator sets the global validator instance
func SetGlobalValidator(v *Validator) {
	globalValidator = v
}

// GetGlobalValidator returns the global validator instance
func GetGlobalValidator() *Validator {
	if globalValidator == nil {
		globalValidator = New()
	}
	return globalValidator
}

// Validate validates using the global validator
func Validate(s interface{}) error {
	return GetGlobalValidator().Validate(s)
}

// ValidateVar validates a variable using the global validator
func ValidateVar(field interface{}, tag string) error {
	return GetGlobalValidator().ValidateVar(field, tag)
}

// ValidateStruct validates a struct and returns a structured error
func ValidateStruct(s interface{}) *pkgErrors.ErrorList {
	if err := Validate(s); err != nil {
		if errorList, ok := err.(*pkgErrors.ErrorList); ok {
			return errorList
		}

		// If not an ErrorList, create one
		errorList := pkgErrors.NewErrorList()
		errorList.Add(pkgErrors.Validation(err.Error()))
		return errorList
	}

	return nil
}

// MustValidate validates and panics on validation error
func MustValidate(s interface{}) {
	if err := Validate(s); err != nil {
		panic(fmt.Sprintf("validation failed: %v", err))
	}
}

// IsValidEmail checks if a string is a valid email
func IsValidEmail(email string) bool {
	return ValidateVar(email, "email") == nil
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(url string) bool {
	return ValidateVar(url, "url") == nil
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(uuid string) bool {
	return ValidateVar(uuid, "uuid") == nil
}

// IsValidPhone checks if a string is a valid phone number
func IsValidPhone(phone string) bool {
	return ValidateVar(phone, "phone") == nil
}

// IsValidSlug checks if a string is a valid slug
func IsValidSlug(slug string) bool {
	return ValidateVar(slug, "slug") == nil
}

// IsValidUsername checks if a string is a valid username
func IsValidUsername(username string) bool {
	return ValidateVar(username, "username") == nil
}

// IsStrongPassword checks if a password meets strength requirements
func IsStrongPassword(password string) bool {
	return ValidateVar(password, "password_strength") == nil
}

// Common validation tags as constants
const (
	Required         = "required"
	Email            = "email"
	Min              = "min"
	Max              = "max"
	Len              = "len"
	Alpha            = "alpha"
	Alphanum         = "alphanum"
	Numeric          = "numeric"
	URL              = "url"
	UUID             = "uuid"
	UUID4            = "uuid4"
	OneOf            = "oneof"
	GreaterThan      = "gt"
	GreaterThanEq    = "gte"
	LessThan         = "lt"
	LessThanEq       = "lte"
	EqualField       = "eqfield"
	NotEqualField    = "nefield"
	PasswordStrength = "password_strength"
	Phone            = "phone"
	Slug             = "slug"
	NoSpaces         = "no_spaces"
	Username         = "username"
)
