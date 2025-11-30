package http

import (
	"context"
	"time"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type ClientOption func(opt *ClientOptions)

type ClientOptions struct {
	logger     domain.Logger
	ctx        context.Context
	baseURI    string
	token      string
	authScheme string
	timeout    time.Duration
}

func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		logger:     log.NewSlog(),
		ctx:        context.Background(),
		baseURI:    "",
		token:      "",
		authScheme: "Bearer",
		timeout:    time.Second * 10,
	}
}

func WithLogger(l domain.Logger) ClientOption {
	return func(opt *ClientOptions) {
		opt.logger = l
	}
}

func WithContext(ctx context.Context) ClientOption {
	return func(opt *ClientOptions) {
		opt.ctx = ctx
	}
}

func WithToken(token string) ClientOption {
	return func(opt *ClientOptions) {
		opt.token = token
	}
}

func WithBaseURI(uri string) ClientOption {
	return func(opt *ClientOptions) {
		opt.baseURI = uri
	}
}

func WithAuthScheme(scheme string) ClientOption {
	return func(opt *ClientOptions) {
		opt.authScheme = scheme
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(opt *ClientOptions) {
		opt.timeout = timeout
	}
}
