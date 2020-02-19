# gobase

[![godoc - documentation](https://godoc.org/github.com/cancue/gobase?status.svg)](https://pkg.go.dev/github.com/cancue/gobase@v0.1.1)
[![go report card](https://goreportcard.com/badge/github.com/cancue/gobase)](https://goreportcard.com/report/github.com/cancue/gobase)
[![codecov - code coverage](https://img.shields.io/codecov/c/github/cancue/gobase.svg?style=flat-square)](https://codecov.io/gh/cancue/gobase)
[![github action - test](https://github.com/cancue/gobase/workflows/test/badge.svg)](https://github.com/cancue/gobase/actions)

**gobase** is a web framework with basic settings and structure, wrapping [echo](https://github.com/labstack/echo).

## Installation
```go
go get github.com/cancue/gobase
```

## Example
```go
package main

import (
	"net/http"

	"github.com/cancue/gobase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	// Gobase instance
	g := gobase.NewWithConfig(&gobase.Config{
		Stage:        "local",
		Name:         "gobase-example",
		Port:         ":65533",
		ReadTimeout:  0,
		WriteTimeout: 0,
	})
	// Middleware
	g.Use(middleware.Secure())
	// Routes
	g.GET("/", hello)
	// Start server
	g.Start()
}
```

## Demo
[gobase-demo](https://github.com/cancue/gobase-demo)

## License

[MIT](https://github.com/cancue/gobase/blob/master/LICENSE)
