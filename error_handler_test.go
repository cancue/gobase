package gobase

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHTTPErrorHandler(t *testing.T) {
	existing := gobase
	buf := new(bytes.Buffer)
	gobase = new(Server)

	logger := logrus.New()
	logger.Out = buf
	gobase.Logger = logger

	defer func() {
		gobase = existing
	}()

	t.Run("HTTP Error", func(t *testing.T) {
		server := echo.New()
		res := httptest.NewRecorder()
		ctx := server.NewContext(httptest.NewRequest("get", "/", nil), res)
		err := echo.NewHTTPError(400, "raboof")

		httpErrorHandler(err, ctx)
		result := buf.String()

		assert.Zero(t, result)
	})

	t.Run("Server Error", func(t *testing.T) {
		res := httptest.NewRecorder()
		server := echo.New()
		ctx := server.NewContext(httptest.NewRequest("get", "/", nil), res)
		expect := "foobar"
		err := errors.New(expect)

		httpErrorHandler(err, ctx)
		result := buf.String()

		assert.Contains(t, result, expect)
	})
}
