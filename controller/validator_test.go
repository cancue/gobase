package controller

import (
	"errors"
	"reflect"
	"testing"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FieldError struct {
	mock.Mock
}

func (f *FieldError) Field() string                             { args := f.Called(); return args.String(0) }
func (f *FieldError) Tag() string                               { args := f.Called(); return args.String(0) }
func (f *FieldError) Param() string                             { args := f.Called(); return args.String(0) }
func (f *FieldError) ActualTag() (result string)                { return }
func (f *FieldError) Error() (result string)                    { return }
func (f *FieldError) Namespace() (result string)                { return }
func (f *FieldError) StructNamespace() (result string)          { return }
func (f *FieldError) StructField() (result string)              { return }
func (f *FieldError) Value() (result interface{})               { return }
func (f *FieldError) Kind() (result reflect.Kind)               { return }
func (f *FieldError) Type() (result reflect.Type)               { return }
func (f *FieldError) Translate(_ ut.Translator) (result string) { return }

func TestNewInvalidErrMap(t *testing.T) {
	t.Run("invalid error", func(t *testing.T) {
		expect := "invalid Field"
		result := newInvalidErrMap(errors.New(""))

		assert.Equal(t, expect, result.Message)
	})

	t.Run("single invalid field", func(t *testing.T) {
		mockFieldError := new(FieldError)
		mockFieldError.On("Field").Return("email")
		mockFieldError.On("Tag").Return("max")
		mockFieldError.On("Param").Return("64")
		verr := validator.ValidationErrors{
			mockFieldError,
		}
		expect := errMap{
			"email field is invalid",
			[]map[string]string{
				{"email": "max 64"},
			},
		}

		result := newInvalidErrMap(verr)

		assert.Equal(t, expect, result)
		mockFieldError.AssertExpectations(t)
	})

	t.Run("multiple invalid fields", func(t *testing.T) {
		emailError := new(FieldError)
		emailError.On("Field").Return("email")
		emailError.On("Tag").Return("max")
		emailError.On("Param").Return("64")
		passwordError := new(FieldError)
		passwordError.On("Field").Return("password")
		passwordError.On("Tag").Return("min")
		passwordError.On("Param").Return("8")
		vErr := validator.ValidationErrors{
			emailError,
			passwordError,
		}
		expect := errMap{
			"some fields are invalid",
			[]map[string]string{
				{"email": "max 64"},
				{"password": "min 8"},
			},
		}

		result := newInvalidErrMap(vErr)

		assert.Equal(t, expect, result)
		emailError.AssertExpectations(t)
		passwordError.AssertExpectations(t)
	})
}
