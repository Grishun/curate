package http

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/service"
)

type HandlerOption func(*HandlerOptions)

type HandlerOptions struct {
	historyLimit uint32
	currecies    []string
	service      *service.Service
	logger       domain.Logger
}

func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		historyLimit: 10,
		currecies:    []string{"BTC", "ETH", "TRX"},
		service:      service.New(),
		logger:       log.NewSlog(),
	}
}

func WithHistoryLimit(l uint32) HandlerOption {
	return func(options *HandlerOptions) {
		options.historyLimit = l
	}
}

func WithCurrencies(currencies []string) HandlerOption {
	return func(options *HandlerOptions) {
		options.currecies = currencies
	}
}

func WithService(s *service.Service) HandlerOption {
	return func(options *HandlerOptions) {
		options.service = s
	}
}

func WithLogger(l domain.Logger) HandlerOption {
	return func(options *HandlerOptions) {
		options.logger = l
	}
}
