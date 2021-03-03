package controller

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// register json key as Field.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) (name string) {
		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		return
	})
}

// RespondInvalidBody returns ErrMap from validation.
func RespondInvalidBody(err error, ctx *fiber.Ctx) {
	ctx.Status(fiber.StatusBadRequest).JSON(newInvalidErrMap(err))
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
}

/* private */

// newInvalidErrMap returns rich error message.
func newInvalidErrMap(err error) (result errMap) {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		result.Message = "invalid Field"
		return
	}

	if len(validationErrors) > 1 {
		result.Message = "some fields are invalid"
	} else {
		result.Message = validationErrors[0].Field() + " field is invalid"
	}

	for _, ve := range validationErrors {
		result.InvalidFields = append(
			result.InvalidFields,
			map[string]string{ve.Field(): strings.Trim(ve.Tag()+" "+ve.Param(), " ")},
		)
	}
	return
}

// errMap describes the error information in the format for the end user.
type errMap struct {
	Message       string              `json:"message"`
	InvalidFields []map[string]string `json:"invalidFields"`
}
