package http

import (
	"github.com/Grishun/curate/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	service *service.Service
}

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Handlers)(nil)

func NewHandlers(srv *service.Service) *Handlers {
	return &Handlers{service: srv}
}

func (s *Handlers) GetAllRates(c *fiber.Ctx, params GetAllRatesParams) error {
	if params.Limit == nil {
		return fiber.NewError(fiber.StatusBadRequest, "limit is required")
	}
	ratesMap, err := s.service.GetRates(c.Context(), uint(*params.Limit))
	if err != nil {
		return err
	}

	return c.JSON(ratesMap)
}

func (s *Handlers) GetRateByCurrency(c *fiber.Ctx, currency string, params GetRateByCurrencyParams) error {
	if params.Limit == nil {
		return fiber.NewError(fiber.StatusBadRequest, "limit is required")
	}

	ratesMap, err := s.service.GetRate(c.Context(), currency, uint(*params.Limit))
	if err != nil {
		return err
	}

	return c.JSON(ratesMap)
}
