package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/cancue/gobase/logger"
)

// ErrorHandler is for logging internal server errors.
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// Set err
	e, ok := err.(*echo.HTTPError)
	if !ok {
		logger.Get().Error(err)
		e = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	if e.Internal != nil {
		if ee, ok := e.Internal.(*echo.HTTPError); ok {
			e = ee
		}
	}

	// TODO github.com/labstack/echo - Issue #1426
	code := e.Code
	message := e.Message

	if m, ok := message.(string); ok {
		message = echo.Map{"message": m}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // TODO github.com/labstack/echo - Issue #608
		err = c.NoContent(e.Code)
	} else {
		err = c.JSON(code, message)
	}

	if err != nil {
		logger.Get().Error(err)
	}
}
