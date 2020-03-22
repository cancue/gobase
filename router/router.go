package router

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/cancue/gobase/validator"
)

type (
	// Router public
	Router func(Server)

	// Server is interface for Router
	Server interface {
		Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
		GET(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		POST(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		PUT(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		DELETE(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		PATCH(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		OPTIONS(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		HEAD(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	}
)

// Context is for controller.
type Context echo.Context
type ctrlr interface {
	Exec(Context) (interface{}, error)
}

// Controller returns echo controller from a struct written in suggested form.
func Controller(obj interface{}) func(echo.Context) error {
	t := reflect.TypeOf(obj)

	return func(c echo.Context) (err error) {
		con := reflect.New(t.Elem()).Interface()

		if err := c.Bind(con); err != nil {
			return err
		}

		err = validator.Validate.Struct(con)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, validator.NewInvalidErrMap(err))
		}

		executer := con.(ctrlr)
		res, err := executer.Exec(c)
		if err != nil {
			return
		}

		return c.JSON(200, res)
	}
}

// Routes returns labstack/echo Routes result in marshalled json.
func Routes(echo *echo.Echo) (data []byte, err error) {
	routes := echo.Routes()
	data, err = json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return
	}

	return
}
