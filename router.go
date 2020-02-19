package gobase

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type (
	controller interface {
		Exec(Context) (interface{}, error)
	}

	genController func() controller
)

// Controller returns controller from a struct written in suggested form.
func Controller(obj interface{}) func(echo.Context) error {
	t := reflect.TypeOf(obj)

	return func(c echo.Context) (err error) {
		con := reflect.New(t.Elem()).Interface()

		if err := c.Bind(con); err != nil {
			return err
		}

		err = Validator.Struct(con)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, NewInvalidErrorMap(err))
		}

		executer := con.(controller)
		res, err := executer.Exec(c)
		if err != nil {
			return
		}

		return c.JSON(200, res)
	}
}
