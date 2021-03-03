package errors

import (
	"runtime"

	"github.com/pkg/errors"
)

// richError is rich error struct.
type richError struct {
	Raw  error
	Data interface{}
	*stack
}

// Error returns description to implement basic error.
func (err *richError) Error() string {
	if err.Raw == nil {
		return "debug\n"
	}

	return err.Raw.Error()
}

type stack []uintptr

// trace returns stacked traces for logger.
func (s *stack) trace() errors.StackTrace {
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
