package http

import (
	"github.com/gofiber/fiber/v2"
)

type Server struct{}

// ensure that we've conformed to the `ServerInterface` with a compile-time check
var _ ServerInterface = (*Server)(nil)

func New() Server {
	return Server{}
}

func (s *Server) GetAllRates(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetRateByCode(c *fiber.Ctx, code string) error {
	//TODO implement me
	panic("implement me")
}
