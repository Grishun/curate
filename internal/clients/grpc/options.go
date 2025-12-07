package grpc

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type ClientOptions struct {
	logger     domain.Logger
	serverAddr string
}

type ClientOption func(*ClientOptions)

func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		logger:     log.NewSlog(),
		serverAddr: "127.0.0.1:8081",
	}
}

func WithServerAddr(addr string) ClientOption {
	return func(o *ClientOptions) {
		o.serverAddr = addr
	}
}

func WithLogger(logger domain.Logger) ClientOption {
	return func(o *ClientOptions) {
		o.logger = logger
	}
}
