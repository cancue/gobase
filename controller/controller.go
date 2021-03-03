package controller

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// Controller is interface to handle routines in router.
type Controller interface {
	Exec(ctx *fiber.Ctx) (interface{}, error)
}

// RouteJSON handles json request.
func RouteJSON(ctrl interface{}) func(ctx *fiber.Ctx) error {
	t := reflect.TypeOf(ctrl)

	return func(ctx *fiber.Ctx) error {
		controller := reflect.New(t.Elem()).Interface()
		_ = ctx.BodyParser(&controller)

		if controller == nil {
			controller = reflect.New(t.Elem()).Interface()
		}

		if err := ValidateAndExecute(controller.(Controller), ctx); err != nil {
			return err
		}
		return nil
	}
}

// ValidateAndExecute handles controllers set by router.
func ValidateAndExecute(ctrl Controller, ctx *fiber.Ctx) error {
	if err := validate.Struct(ctrl); err != nil {
		RespondInvalidBody(err, ctx)
		return nil
	}

	res, err := ctrl.Exec(ctx)
	if err != nil {
		return err
	}

	ctx.JSON(res)
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return nil
}
