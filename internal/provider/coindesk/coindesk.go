package coindesk

import (
	"context"
	http2 "net/http"
	"net/url"
	"strings"

	"github.com/Grishun/curate/internal/clients/http"
	"github.com/Grishun/curate/internal/domain"
)

const multSymbolsPrice = "/data/pricemulti"

type Coindesk struct {
	uri        string
	token      string
	httpClient http.Client
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
		logger:     options.logger,
	}
}

func (c *Coindesk) Name() string {
	uri, _ := url.Parse(c.uri)

	return uri.Hostname()
}

func (c *Coindesk) Fetch(ctx context.Context) (map[string]float64, error) {
	result := make(map[string]map[string]float64)

	params := url.Values{}
	params.Add("fsyms", strings.Join(c.currencies, ","))
	params.Add("tsyms", c.quote)
	params.Add("api_key", c.token)

	headers := make(http2.Header)
	headers.Add("Content-type", "application/json")
	headers.Add("charset", "UTF-8")

	_, err := c.httpClient.NewRequest(
		http.WithMethod(http2.MethodGet),
		http.WithQueryParams(params),
		http.WithHeaders(headers),
		http.WithURI(c.uri+multSymbolsPrice),
		http.WithRequestContext(ctx),
		http.WithUnmarshallTo(&result),
	)

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
