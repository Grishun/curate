package http

import "github.com/Grishun/curate/internal/service"

type HandlerOption func(*HandlerOptions)

type HandlerOptions struct {
	historyLimit uint32
	currecies    []string
	service      *service.Service
}

func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{
		historyLimit: 10,
		currecies:    []string{"BTC", "ETH", "TRX"},
		service:      service.New(),
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
