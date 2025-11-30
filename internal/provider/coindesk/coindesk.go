package coindesk

import (
	"context"
	"net/url"
	"strings"

	"github.com/Grishun/curate/internal/domain"
	"github.com/Grishun/curate/internal/log"
	"resty.dev/v3"
)

const multSymbolsPrice = "/data/pricemulti"

type Coindesk struct {
	uri        string
	token      string
	httpClient *resty.Client //FIXME: use interface instead (implement in domain)
	quote      string
	currencies []string
	logger     domain.Logger
}

func New(opts ...Option) *Coindesk {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Coindesk{
		uri:        options.uri,
		token:      options.token,
		httpClient: options.httpClient,
		quote:      options.quote,
		currencies: options.currencies,
		logger:     log.NewSlog(),
	}
}

func (c *Coindesk) Name() string {
	uri, _ := url.Parse(c.uri)

	return uri.Hostname()
}

func (c *Coindesk) Fetch(ctx context.Context) (map[string]float64, error) {
	result := make(map[string]map[string]float64)

	params := map[string]string{
		"fsyms":   strings.Join(c.currencies, ","),
		"tsyms":   c.quote,
		"api_key": c.token,
	}

	_, err := c.httpClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetHeader("Content-type", "application/json").
		SetHeader("charset", "UTF-8").
		SetAuthScheme("Bearer").
		SetAuthToken(c.token).
		SetResult(&result).
		SetLogger(c.logger).
		Get(c.uri + multSymbolsPrice)

	return convert(result, c.quote), err
}

func convert(rates map[string]map[string]float64, quote string) map[string]float64 {
	if len(rates) == 0 || len(quote) == 0 {
		return nil
	}

	res := make(map[string]float64, len(rates))
	for currency, rate := range rates {
		if val, ok := rate[quote]; ok {
			res[currency] = val
		}
	}

	return res
}
