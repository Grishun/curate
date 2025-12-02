package http

import (
	"github.com/Grishun/curate/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	service     *service.Service
	histryLimit uint32
}

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Handlers)(nil)

func NewHandlers(srv *service.Service, historyLimit uint32) *Handlers { //TODO: add options
	return &Handlers{
		service:     srv,
		histryLimit: historyLimit,
	}
}

func (s *Handlers) GetAllRates(c *fiber.Ctx, params GetAllRatesParams) error {
	limit := s.histryLimit
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > s.histryLimit {
		limit = s.histryLimit
	}

	ratesMap, err := s.service.GetRates(c.Context(), limit)
	if err != nil {
		return err
	}

	return c.JSON(ratesMap)
}

func (s *Handlers) GetRateByCurrency(c *fiber.Ctx, currency string, params GetRateByCurrencyParams) error {
	limit := s.histryLimit
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > s.histryLimit {
		limit = s.histryLimit
	}

	ratesMap, err := s.service.GetRate(c.Context(), currency, limit)
	if err != nil {
		return err
	}

	return c.JSON(ratesMap)
}

func (s *Handlers) HealthCheck(c *fiber.Ctx) error {
	if err := s.service.HealthCheck(c.Context()); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
