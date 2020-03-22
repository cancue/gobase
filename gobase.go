/*
Package gobase is a web framework with basic settings and structure, wrapping labstack/echo.

Example:
	package main

	import (
		"net/http"

		"github.com/labstack/echo/v4"

		"github.com/cancue/gobase"
		"github.com/cancue/gobase/config"
		"github.com/cancue/gobase/router"
	)

	func main() {
		gb := gobase.Server{
			Config: &config.Config{
				Stage: "local",
				YAML: map[string]interface{}{
					"name": "gobase-demo",
					"server": map[string]interface{}{
						"port": 8888,
						"timeout": map[string]interface{}{
							"read":  600,
							"write": 600,
						},
					},
				},
			},
			Router: func(s router.Server) {
				s.GET("/", func(ctx echo.Context) error {
					return ctx.String(http.StatusOK, "Hello, World!")
				})
			},
		}

		gb.Start()
	}

You may want to check out https://github.com/cancue/gobase-demo
*/
package gobase

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/cancue/gobase/config"
	"github.com/cancue/gobase/db/pg"
	gobasehttp "github.com/cancue/gobase/http"
	"github.com/cancue/gobase/logger"
	"github.com/cancue/gobase/router"
)

// Server is main dish.
type Server struct {
	Config       *config.Config
	Router       router.Router
	Logger       logger.Logger
	ErrorHandler func(err error, c echo.Context)
	Middlewares  []echo.MiddlewareFunc
	Echo         *echo.Echo
}

// Start starts server.
func (gb *Server) Start() {
	gb.Set(1, ".")

	logger.Get().Fatal(gb.Echo.StartServer(httpConfigFrom(gb.Config)))
}

// Set sets server with config/${stage}.yaml.
func (gb *Server) Set(lenCaller int, rpath string) {
	if gb.Config == nil {
		gb.Config = config.SetWithCallerLength(0+lenCaller, rpath)
	}

	gb.Echo = echo.New()

	if gb.Logger != nil {
		logger.SetCustomLogger(gb.Logger)
	} else {
		gb.Logger = logger.Set(gb.Config)
	}

	gb.Echo.Use(gb.Middlewares...)

	if gb.ErrorHandler != nil {
		gb.Echo.HTTPErrorHandler = gb.ErrorHandler
	} else {
		gb.ErrorHandler = gobasehttp.ErrorHandler
		gb.Echo.HTTPErrorHandler = gb.ErrorHandler
	}

	if gb.Router != nil {
		gb.Router(gb.Echo)
	}

	if _, ok := gb.Config.YAML["db"].(interface{}); ok {
		pg.Set(gb.Config)
	}

	gb.Echo.HideBanner = true
}

func httpConfigFrom(conf *config.Config) *http.Server {
	server := conf.YAML["server"].(map[string]interface{})
	timeout := server["timeout"].(map[string]interface{})

	return &http.Server{
		Addr:         ":" + strconv.Itoa(server["port"].(int)),
		ReadTimeout:  time.Duration(timeout["read"].(int)) * time.Second,
		WriteTimeout: time.Duration(timeout["write"].(int)) * time.Second,
	}
}
