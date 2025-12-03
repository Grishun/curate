package http

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/service"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	recover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func NewRouter(opts ...RouterOption) *fiber.App {
	options := NewRouterOptions()

	for _, opt := range opts {
		opt(options)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "currate-rest",
	})

	app.Use(
		recover.New(),
		requestid.New(),
		fiberlogger.New(),
	)

	api := app.Group("/")

	handlers := NewHandlers(
		WithService(options.service),
		WithHistoryLimit(options.historyLimit),
		WithCurrencies(options.currencies),
		WithLogger(options.logger),
	)

	RegisterHandlersWithOptions(api, handlers, FiberServerOptions{})

	return app
}

type RouterOption func(*RouterOptions)

type RouterOptions struct {
	service      *service.Service
	logger       domain.Logger
	historyLimit uint32
	currencies   []string
}

func NewRouterOptions() *RouterOptions {
	return &RouterOptions{
		service:      service.New(),
		logger:       log.NewSlog(),
		historyLimit: 100,
		currencies:   []string{"BTC", "ETH", "TRX"},
	}
}

func WithRouterService(s *service.Service) RouterOption {
	return func(o *RouterOptions) {
		o.service = s
	}
}

func WithRouterLogger(l domain.Logger) RouterOption {
	return func(o *RouterOptions) {
		o.logger = l
	}
}

func WithRouterHistoryLimit(l uint32) RouterOption {
	return func(o *RouterOptions) {
		o.historyLimit = l
	}
}

func WithRouterCurrencies(currencies []string) RouterOption {
	return func(o *RouterOptions) {
		o.currencies = currencies
	}
}
