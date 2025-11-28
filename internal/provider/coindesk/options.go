package coindesk

import (
	"time"

	"resty.dev/v3"
)

type Option func(*Options)

type Options struct {
	uri             string
	token           string
	pollingInterval time.Duration
	httpClient      *resty.Client //FIXME: use interface instead (implement in domain)
	quote           string
	currencies      []string
}

func NewOptions() *Options {
	return &Options{
		uri:             "https://min-api.cryptocompare.com",
		token:           "", // not required for min api
		pollingInterval: 1 * time.Minute,
		httpClient:      resty.New(),
		quote:           "USD",
		currencies:      []string{"BTC", "ETH", "TRX"}, //TODO: create custom type for codes
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

func WithHTTPClient(httpClient *resty.Client) Option {
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

func WithPollingInterval(pollingInterval time.Duration) Option {
	return func(options *Options) {
		options.pollingInterval = pollingInterval
	}
}
