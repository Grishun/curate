package coindesk

import (
	"context"
	http2 "net/http"
	"net/url"
	"strings"

	"github.com/Grishun/curate/internal/clients/rest"
)

const multSymbolsPrice = "/data/pricemulti"

type Coindesk struct {
	options *Options
}

func New(opts ...Option) *Coindesk {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	return &Coindesk{options: options}
}

func (c *Coindesk) Name() string {
	uri, _ := url.Parse(c.options.uri)

	return uri.Hostname()
}

func (c *Coindesk) Fetch(ctx context.Context) (map[string]float64, error) {
	result := make(map[string]map[string]float64)

	params := url.Values{}
	params.Add("fsyms", strings.Join(c.options.currencies, ","))
	params.Add("tsyms", c.options.quote)
	params.Add("api_key", c.options.token)

	headers := make(http2.Header)
	headers.Add("Content-type", "application/json")
	headers.Add("charset", "UTF-8")

	_, err := c.options.httpClient.Do(ctx,
		rest.WithMethod(http2.MethodGet),
		rest.WithQueryParams(params),
		rest.WithHeaders(headers),
		rest.WithURI(c.options.uri+multSymbolsPrice),
		rest.WithUnmarshallTo(&result),
	)

	return convert(result, c.options.quote), err
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
