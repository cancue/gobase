# gobase

[![godoc - documentation](https://godoc.org/github.com/cancue/gobase?status.svg)](https://pkg.go.dev/github.com/cancue/gobase@v0.2.1)
[![go report card](https://goreportcard.com/badge/github.com/cancue/gobase)](https://goreportcard.com/report/github.com/cancue/gobase)
[![github action - test](https://github.com/cancue/gobase/workflows/test/badge.svg)](https://github.com/cancue/gobase/actions)
[![codecov - code coverage](https://img.shields.io/codecov/c/github/cancue/gobase.svg?style=flat-square)](https://codecov.io/gh/cancue/gobase)

**gobase** is a web framework with basic settings and structure, wrapping [fiber](https://github.com/gofiber/fiber).

## Installation
```go
go get github.com/cancue/gobase
```

## Example
```go
package main

import (
	"github.com/cancue/gobase"
	"github.com/cancue/gobase/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := config.Config{
		Stage:        "local",
		Name:         "gobase",
		Domain:       "localhost",
		Port:         3000,
		AllowOrigins: []string{"http://localhost:3000"},
	}

	gobase.Start(&config, router)
}

func router(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})
}
```

## Demo
[gobase-demo](https://github.com/cancue/gobase-demo)

## License

[MIT](https://github.com/cancue/gobase/blob/master/LICENSE)
