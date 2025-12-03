package coindesk

import (
	"github.com/Grishun/curate/internal/clients/rest"
	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
)

type Option func(*Options)

type Options struct {
	uri        string
	token      string
	httpClient *rest.Client
	quote      string
	currencies []string
	logger     domain.Logger
}

func NewOptions() *Options {
	return &Options{
		uri:        "https://min-api.cryptocompare.com",
		token:      "", // not required for min api
		httpClient: rest.NewClient(),
		quote:      "USD",
		currencies: []string{"BTC", "ETH", "TRX"}, // TODO: create custom type for currency codes
		logger:     log.NewSlog(),
	}
}

func WithURI(uri string) Option {
	return func(options *Options) {
		options.uri = uri
	}
}

func WithToken(token string) Option {
	return func(options *Options) {
		options.token = token
	}
}

func WithHTTPClient(httpClient *rest.Client) Option {
	return func(options *Options) {
		options.httpClient = httpClient
	}
}

func WithQuote(quote string) Option {
	return func(options *Options) {
		options.quote = quote
	}
}

func WithCurrencies(currencies ...string) Option {
	return func(options *Options) {
		options.currencies = currencies
	}
}

func WithLogger(l domain.Logger) Option {
	return func(options *Options) {
		options.logger = l
	}
}
