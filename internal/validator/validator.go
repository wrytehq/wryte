package validator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
	decoder  *form.Decoder
}

type ValidationErrors struct {
	errors map[string]string
}

func New() *Validator {
	validate := validator.New()
	decoder := form.NewDecoder()

	// Register custom validation tags if needed
	// validate.RegisterValidation("custom_tag", customValidationFunc)

	return &Validator{
		validate: validate,
		decoder:  decoder,
	}
}

func (v *Validator) DecodeAndValidate(r *http.Request, dst interface{}) (*ValidationErrors, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("error parsing form: %w", err)
	}

	// Decode form values into the struct
	if err := v.decoder.Decode(dst, r.PostForm); err != nil {
		return nil, fmt.Errorf("error decoding form: %w", err)
	}

	// Validate the struct
	if err := v.validate.Struct(dst); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.formatErrors(validationErrors), nil
		}
		return nil, fmt.Errorf("unexpected validation error: %w", err)
	}

	return &ValidationErrors{errors: make(map[string]string)}, nil
}

func (v *Validator) Validate(s interface{}) *ValidationErrors {
	err := v.validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.formatErrors(validationErrors)
		}
	}
	return &ValidationErrors{errors: make(map[string]string)}
}

func (v *Validator) formatErrors(errs validator.ValidationErrors) *ValidationErrors {
	errors := make(map[string]string)

	for _, err := range errs {
		fieldName := v.formatFieldName(err.Field())
		errors[fieldName] = v.getErrorMessage(err)
	}

	return &ValidationErrors{errors: errors}
}

func (v *Validator) formatFieldName(field string) string {
	return strings.ToLower(field)
}

func (v *Validator) getErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be no more than %s characters", field, err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, err.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", field)
	case "url":
		return "Must be a valid URL"
	case "uri":
		return "Must be a valid URI"
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.errors) > 0
}

func (ve *ValidationErrors) Get(field string) string {
	return ve.errors[field]
}

func (ve *ValidationErrors) Has(field string) bool {
	_, exists := ve.errors[field]
	return exists
}

// All returns all errors (for debugging)
func (ve *ValidationErrors) All() map[string]string {
	return ve.errors
}

// AddError adds a custom error for a field
func (ve *ValidationErrors) AddError(field, message string) {
	if ve.errors == nil {
		ve.errors = make(map[string]string)
	}
	ve.errors[field] = message
}
