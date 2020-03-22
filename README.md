# gobase

[![godoc - documentation](https://godoc.org/github.com/cancue/gobase?status.svg)](https://pkg.go.dev/github.com/cancue/gobase@v0.2.0)
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
```

## Demo
[gobase-demo](https://github.com/cancue/gobase-demo)

## License

[MIT](https://github.com/cancue/gobase/blob/master/LICENSE)
