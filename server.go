/*
Package gobase is a web framework with basic settings and structure, wrapping labstack/echo.

Example:
	package main
	import (
		"net/http"
		"github.com/cancue/gobase"
		"github.com/labstack/echo/v4/middleware"
	)
	// Handler
	func hello(c gobase.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
	func main() {
		// Gobase instance
		g := gobase.NewWithConfig(gobase.Config{
			Stage: "local",
			Name: "gobase-example",
			Port: 65535,
			ReadTimeout: 0,
			WriteTimeout: 0,
		})
		// Middleware
		g.Use(middleware.Secure())
		// Routes
		g.GET("/", hello)
		// Start server
		g.Start()
	}

You may want to check out https://github.com/cancue/gobase-demo
*/
package gobase

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/cancue/gobase/util"
)

type (
	// Server is the top-level framework instance.
	Server struct {
		Config     *Config
		Logger     *logrus.Logger
		echo       *echo.Echo
		httpConfig http.Server
	}

	// Config defines the config for Server.
	Config struct {
		Stage             string
		Name              string
		Port              string
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		HTTPRequestLogger bool
	}

	// Context for controller is an alias for labstack/echo
	Context echo.Context
)

var (
	// NewHTTPError is an alias for labstack/echo.
	NewHTTPError = echo.NewHTTPError

	gobase *Server
)

// New returns server with Config by config/{stage}.yaml written in suggested form.
func New(stage string) *Server {
	path := filepath.Join(util.CallerDir(), "config", stage+".yaml")
	yaml := util.ReadYAMLFile(path)

	server := yaml["server"].(map[string]interface{})
	timeout := server["timeout"].(map[string]interface{})

	config := Config{
		Stage:        stage,
		Name:         server["name"].(string),
		Port:         ":" + strconv.Itoa(server["port"].(int)),
		ReadTimeout:  time.Duration(timeout["read"].(int)) * time.Second,
		WriteTimeout: time.Duration(timeout["read"].(int)) * time.Second,
	}

	return NewWithConfig(&config)
}

// NewWithConfig returns server with CustomConfig.
func NewWithConfig(config *Config) *Server {
	gobase = new(Server)

	gobase.Config = config
	gobase.Logger = logrus.New()
	gobase.httpConfig = http.Server{
		Addr:         config.Port,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	gobase.echo = echo.New()
	gobase.echo.HideBanner = true
	gobase.echo.HTTPErrorHandler = httpErrorHandler

	return gobase
}

// Start starts server.
func (s *Server) Start() {
	s.Logger.Fatal(s.echo.StartServer(&s.httpConfig))
}

// Use for middleware is an alias for labstack/echo.
func (s *Server) Use(middleware ...echo.MiddlewareFunc) {
	s.echo.Use(middleware...)
}

// Routes returns labstack/echo Routes result in marshalled json.
func (s *Server) Routes() (data []byte, err error) {
	routes := s.echo.Routes()
	data, err = json.MarshalIndent(routes, "", "  ")
	if err != nil {
		return
	}

	return
}

// GET for router is an alias for labstack/echo.
func (s *Server) GET(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	s.echo.GET(path, controller, m...)
}

// POST for router is an alias for labstack/echo.
func (s *Server) POST(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	s.echo.POST(path, controller, m...)
}

// PUT for router is an alias for labstack/echo.
func (s *Server) PUT(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	s.echo.PUT(path, controller, m...)
}

// DELETE for router is an alias for labstack/echo.
func (s *Server) DELETE(path string, controller echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	s.echo.DELETE(path, controller, m...)
}
