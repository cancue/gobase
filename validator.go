package gobase

import (
	"reflect"
	"strings"

	validation "github.com/go-playground/validator/v10"
)

type (
	// ValidationError describes the error information in the format for the end user.
	ValidationError struct {
		Message       string              `json:"message"`
		InvalidFields []map[string]string `json:"invalidFields"`
	}
)

var (
	// Validator is an alias for go-playground/validator
	Validator *validation.Validate
)

func init() {
	Validator = validation.New()

	// register json key as Field.
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) (name string) {
		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		return
	})
}

// NewInvalidErrorMap returns rich error message.
func NewInvalidErrorMap(err error) interface{} {
	validationErrors, ok := err.(validation.ValidationErrors)
	if !ok {
		return "invalid Field"
	}

	verr := new(ValidationError)

	if len(validationErrors) > 1 {
		verr.Message = "some fields are invalid"
	} else {
		verr.Message = validationErrors[0].Field() + " field is invalid"
	}

	for _, ve := range validationErrors {
		verr.InvalidFields = append(
			verr.InvalidFields,
			map[string]string{ve.Field(): strings.Trim(ve.Tag()+" "+ve.Param(), " ")},
		)
	}

	return verr
}
