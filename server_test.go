package gobase

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewWithConfig(t *testing.T) {
	config := Config{
		Stage:             "foo",
		Name:              "bar",
		Port:              "baz",
		ReadTimeout:       time.Second * 1234,
		WriteTimeout:      time.Second * 4321,
		HTTPRequestLogger: false,
	}

	server := NewWithConfig(&config)

	assert.Equal(t, &config, server.Config)
	assert.IsType(t, new(logrus.Logger), server.Logger)
	assert.Equal(
		t,
		http.Server{
			Addr:         config.Port,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
		server.httpConfig,
	)
	assert.True(t, server.echo.HideBanner)
	assert.Equal(
		t,
		reflect.ValueOf(httpErrorHandler).Pointer(),
		reflect.ValueOf(server.echo.HTTPErrorHandler).Pointer(),
	)
}
