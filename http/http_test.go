package http

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/cancue/gobase/config"
	"github.com/cancue/gobase/logger"
)

func init() {
	conf := config.Config{
		YAML: map[string]interface{}{
			"name": "gobase",
		},
		Stage: "test",
	}
	logger.Set(&conf)
}

func TestHTTPErrorHandler(t *testing.T) {
	buf := new(bytes.Buffer)
	lgr := logger.Get().(*logger.DefaultLogger)
	lgr.Logger.Out = buf

	t.Run("HTTP Error", func(t *testing.T) {
		server := echo.New()
		res := httptest.NewRecorder()
		ctx := server.NewContext(httptest.NewRequest("get", "/", nil), res)
		err := echo.NewHTTPError(400, "raboof")

		ErrorHandler(err, ctx)
		result := buf.String()

		assert.Zero(t, result)
	})

	t.Run("Server Error", func(t *testing.T) {
		res := httptest.NewRecorder()
		server := echo.New()
		ctx := server.NewContext(httptest.NewRequest("get", "/", nil), res)
		expect := "foobar"
		err := errors.New(expect)

		ErrorHandler(err, ctx)
		result := buf.String()

		assert.Contains(t, result, expect)
	})
}
