package http

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	options *HandlerOptions
}

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Handlers)(nil)

func NewHandlers(opts ...HandlerOption) *Handlers {
	options := NewHandlerOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Handlers{
		options: options,
	}
}

func (s *Handlers) GetAllRates(c *fiber.Ctx, params GetAllRatesParams) error {
	limit := s.options.historyLimit
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > s.options.historyLimit {
		limit = s.options.historyLimit
	}

	ratesMap, err := s.options.service.GetRates(c.Context(), limit)
	if err != nil {
		return c.JSON(goErrorTocCustom(err))
	}

	return c.JSON(mapDomainRatesToOpenAPIRates(ratesMap))
}

func (s *Handlers) GetRateByCurrency(c *fiber.Ctx, currency string, params GetRateByCurrencyParams) error {
	limit := s.options.historyLimit
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > s.options.historyLimit {
		limit = s.options.historyLimit
	}

	rates, err := s.options.service.GetRate(c.Context(), currency, limit)
	if err != nil {
		return c.JSON(goErrorTocCustom(err))
	}

	return c.JSON(mapDomainRateToOpenAPIRate(rates))
}

func (s *Handlers) HealthCheck(c *fiber.Ctx) error {
	if err := s.options.service.HealthCheck(c.Context()); err != nil {
		return c.JSON(goErrorTocCustom(err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (s *Handlers) GetAllCurrencies(c *fiber.Ctx) error {
	return c.JSON(s.options.currecies)
}

func goErrorTocCustom(err error) Error {
	return Error{
		Message: err.Error(),
	}
}
