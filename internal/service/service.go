package service

import (
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"github.com/Grishun/curate/internal/storage/memory"
)

type Service struct {
	logger  domain.Logger
	storage domain.Storage
	//TODO: add provider
}

type Option func(*Options)

type Options struct {
	logger  domain.Logger
	storage domain.Storage
}

func WithNewLogger(l domain.Logger) Option {
	return func(options *Options) {
		options.logger = l
	}
}

func WithNewStorage(s domain.Storage) Option {
	return func(options *Options) {
		options.storage = s
	}
}

func NewOptions() *Options {
	return &Options{
		logger:  log.NewSlog(),
		storage: memory.New(),
	}
}

func New(opts ...Option) *Service {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Service{
		logger:  options.logger,
		storage: options.storage,
	}
}
