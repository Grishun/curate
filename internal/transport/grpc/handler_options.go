package grpc

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/transport/grpc/generated"
)

// ------------ grpc handler options ------------

type Handler struct {
	generated.UnimplementedRatesServiceServer
	options *HandlerOptions
}

type HandlerOptions struct {
	service *service.Service
	logger  domain.Logger
}

type HandlerOption func(*HandlerOptions)

func WithHandlerService(s *service.Service) HandlerOption {
	return func(options *HandlerOptions) {
		options.service = s
	}
}

func WithHandlerLogger(l domain.Logger) HandlerOption {
	return func(options *HandlerOptions) {
		options.logger = l
	}
}

func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		service: service.New(),
		logger:  log.NewSlog(),
	}
}

func NewHandler(opts ...HandlerOption) *Handler {
	options := NewHandlerOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Handler{
		UnimplementedRatesServiceServer: generated.UnimplementedRatesServiceServer{},
		options:                         options,
	}
}
