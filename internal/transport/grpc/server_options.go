package grpc

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/service"
)

// ServerOption grpc server option
type ServerOption func(*ServerOptions)

type ServerOptions struct {
	port    string
	host    string
	logger  domain.Logger
	service *service.Service
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		port:    "8081",
		host:    "127.0.0.1",
		logger:  log.NewSlog(),
		service: service.New(),
	}
}

func WithPort(port string) ServerOption {
	return func(options *ServerOptions) {
		options.port = port
	}
}

func WithHost(host string) ServerOption {
	return func(options *ServerOptions) {
		options.host = host
	}
}

func WithLogger(logger domain.Logger) ServerOption {
	return func(options *ServerOptions) {
		options.logger = logger
	}
}

func WithService(srv *service.Service) ServerOption {
	return func(options *ServerOptions) {
		options.service = srv
	}
}
