package http

import (
	"github.com/Grishun/curate/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewRouter(service *service.Service, historyLimit uint32) *fiber.App { // TODO: build with options
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		AppName:               "currate-rest",
	})

	app.Use(
		recover.New(),
		requestid.New(),
		logger.New(),
	) // TODO: add our logger

	api := app.Group("/")

	handlers := NewHandlers(
		WithService(service),
		WithHistoryLimit(historyLimit),
	)

	RegisterHandlersWithOptions(api, handlers, FiberServerOptions{})

	return app
}
