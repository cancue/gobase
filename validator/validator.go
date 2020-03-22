package validator

import (
	"reflect"
	"strings"

	validation "github.com/go-playground/validator/v10"
)

// ErrMap describes the error information in the format for the end user.
type ErrMap struct {
	Message       string              `json:"message"`
	InvalidFields []map[string]string `json:"invalidFields"`
}

var (
	// Validate is an alias for go-playground/validator
	Validate *validation.Validate
)

func init() {
	Validate = validation.New()

	// register json key as Field.
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) (name string) {
		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		return
	})
}

// NewInvalidErrMap returns rich error message.
func NewInvalidErrMap(err error) *ErrMap {
	validationErrors, ok := err.(validation.ValidationErrors)
	if !ok {
		return &ErrMap{Message: "invalid Field"}
	}

	verr := new(ErrMap)

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
