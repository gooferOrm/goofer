package validation

import (
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
	"github.com/gooferOrm/goofer/pkg/schema"
)

// Validator is a wrapper around go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a struct using the "validate" tag
func (v *Validator) Validate(entity interface{}) error {
	return v.validate.Struct(entity)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidateEntity validates an entity and returns a list of validation errors
func (v *Validator) ValidateEntity(entity schema.Entity) ([]ValidationError, error) {
	err := v.validate.Struct(entity)
	if err == nil {
		return nil, nil
	}

	var validationErrors []ValidationError
	
	// Check if the error is a validator.ValidationErrors
	if errors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errors {
			validationErrors = append(validationErrors, ValidationError{
				Field:   e.Field(),
				Message: buildErrorMessage(e),
			})
		}
		return validationErrors, nil
	}
	
	return nil, err
}

// buildErrorMessage builds a human-readable error message
func buildErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		if e.Type().Kind() == reflect.String {
			return "Must be at least " + e.Param() + " characters long"
		}
		return "Must be at least " + e.Param()
	case "max":
		if e.Type().Kind() == reflect.String {
			return "Must be at most " + e.Param() + " characters long"
		}
		return "Must be at most " + e.Param()
	case "oneof":
		return "Must be one of: " + strings.Replace(e.Param(), " ", ", ", -1)
	}
	return "Invalid value"
}

// RegisterValidationHooks registers validation hooks with the repository
func RegisterValidationHooks(registry *schema.SchemaRegistry) error {
	// This would integrate validation with the repository hooks
	// For example, validating entities before saving them
	return nil
}

// ValidatableEntity is an interface for entities that can validate themselves
type ValidatableEntity interface {
	schema.Entity
	Validate() error
}

// ValidateHook is a hook that validates entities before saving them
type ValidateHook struct {
	validator *Validator
}

// NewValidateHook creates a new validate hook
func NewValidateHook() *ValidateHook {
	return &ValidateHook{
		validator: NewValidator(),
	}
}

// BeforeSave validates the entity before saving
func (h *ValidateHook) BeforeSave(entity interface{}) error {
	// Check if the entity implements ValidatableEntity
	if validatable, ok := entity.(ValidatableEntity); ok {
		return validatable.Validate()
	}
	
	// Otherwise, use the validator
	return h.validator.Validate(entity)
}