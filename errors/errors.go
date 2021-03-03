package errors

import (
	"github.com/gofiber/fiber/v2"
)

// NewHTTPError is fiber.NewError wrapper.
var NewHTTPError = fiber.NewError

// FiberErrorHandler is a fiber.Config parameter.
var FiberErrorHandler = func(ctx *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError
	msg := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if code == 500 {
		msg = "Internal Server Error"
		LogError(err)
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return ctx.Status(code).SendString(msg)
}

// Wrap wraps error with additional data and stack info.
func Wrap(err error, data interface{}) error {
	if _, ok := err.(*richError); ok {
		return err
	}

	if _, ok := err.(*fiber.Error); ok {
		return err
	}

	return &richError{
		err,
		data,
		callers(),
	}
}

// Debug returns any data as error for debugging.
func Debug(data interface{}) error {
	return &richError{
		nil,
		data,
		callers(),
	}
}
