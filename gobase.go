package gobase

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"

	"github.com/cancue/gobase/config"
	"github.com/cancue/gobase/errors"
)

// Router .
type Router interface {
	Set(app *fiber.App)
}

// Start .
func Start(conf *config.Config, router func(*fiber.App)) {
	errors.SetLogger(conf)
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.FiberErrorHandler,
	})

	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(conf.AllowOrigins, ", "),
		AllowCredentials: true,
	}))
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	router(app)

	initError := app.Listen(fmt.Sprintf(":%d", conf.Port))
	errors.LogError(initError)
}
