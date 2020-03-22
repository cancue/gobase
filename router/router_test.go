package router

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockController struct {
	Input string `json:"input" query:"input" param:"input" validate:"min=3"`
	Error bool   `json:"error" query:"error" param:"error"`
}

func (mc *mockController) Exec(_ echo.Context) (interface{}, error) {
	if mc.Error {
		err := errors.New("")
		return mc.Input, err
	}
	return mc.Input, nil
}

type keyval struct {
	Keys   []string
	Values []string
}

func TestController(t *testing.T) {
	server := echo.New()
	r := server.Router()
	// set max number of path params for echo.maxParams
	r.Add("GET", "/dummy/:one/:two", func(echo.Context) error { return nil })
	controller := Controller(new(mockController))
	tests := []struct {
		Type      string
		Body      interface{}
		ExecError bool // nothing but clarification
		Expect    bool
	}{
		{"json", "{\"input\": \"valid\"}", false, true},
		{"json", "{\"input\": \"valid\",}", false, false},
		{"json", "{\"input\": \"va\"}", false, false},
		{"json", "{\"input\": \"valid\", \"error\": true}", true, false},
		{"path", keyval{Keys: []string{"input"}, Values: []string{"valid"}}, false, true},
		{"path", keyval{Keys: []string{"input"}, Values: []string{"va"}}, false, false},
		{"path", keyval{Keys: []string{"input", "error"}, Values: []string{"valid", "true"}}, true, false},
		{"query", keyval{Keys: []string{"input"}, Values: []string{"valid"}}, false, true},
		{"query", keyval{Keys: []string{"input"}, Values: []string{"va"}}, false, false},
		{"query", keyval{Keys: []string{"input", "error"}, Values: []string{"valid", "true"}}, true, false},
	}

	for _, test := range tests {
		res := httptest.NewRecorder()
		var ctx echo.Context

		switch test.Type {
		case "json":
			body := test.Body.(string)
			req := httptest.NewRequest("get", "/", strings.NewReader(body))
			req.Header.Add("Content-Type", "application/json")
			ctx = server.NewContext(req, res)
		case "path":
			body := test.Body.(keyval)
			req := httptest.NewRequest("get", "/", strings.NewReader(""))
			ctx = server.NewContext(req, res)
			ctx.SetPath("/:input")
			ctx.SetParamNames(body.Keys...)
			ctx.SetParamValues(body.Values...)
		case "query":
			body := test.Body.(keyval)
			path := "/?"
			for i, k := range body.Keys {
				path = path + k + "=" + body.Values[i] + "&"
			}
			req := httptest.NewRequest("get", path, strings.NewReader(""))
			ctx = server.NewContext(req, res)
		default:
			t.Error("not implemented.")
		}
		err := controller(ctx)

		result := err == nil

		assert.Equal(t, test.Expect, result, err)
	}
}
