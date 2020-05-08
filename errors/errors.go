package errors

import (
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// NewHTTPError returns error with http status code.
var NewHTTPError = echo.NewHTTPError

type stack []uintptr

func (s *stack) Trace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}

func callers() *stack {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// Err is rich error struct.
type Err struct {
	Raw  error
	Data interface{}
	*stack
}

func (err *Err) Error() string {
	if err.Raw == nil {
		return "debug\n"
	}

	return err.Raw.Error()
}

// Wrap wraps error with additional data and stack info.
func Wrap(err error, data interface{}) error {
	if _, ok := err.(*Err); ok {
		return err
	}

	if _, ok := err.(*echo.HTTPError); ok {
		return err
	}

	return &Err{
		err,
		data,
		callers(),
	}
}

// Debug returns anyway error for debugging.
func Debug(data interface{}) error {
	return &Err{
		nil,
		data,
		callers(),
	}
}
