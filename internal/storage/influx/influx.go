package influx

import (
	"context"
	"fmt"
	http2 "net/http"
	"slices"
	"strings"
	"time"

	"github.com/Grishun/curate/internal/clients/rest"
	"github.com/Grishun/curate/internal/domain"
	influx "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/pkg/errors"
)

type Client struct {
	client     *influx.Client
	options    *Options
	httpClient domain.HTTPClient
}

var (
	ErrUnavailableCurrency = errors.New("don't know this currency")
	ErrFailedToParseData   = errors.New("error parsing data from influx map")
)

func NewClient(opts ...Option) (*Client, error) {
	options := NewOptions(opts...)

	for _, opt := range opts {
		opt(options)
	}

	httpClient := rest.NewClient(
		rest.WithToken(options.token),
		rest.WithAuthScheme("Bearer"),
	)

	influxClient, err := influx.New(influx.ClientConfig{
		Host:     options.hostURI,
		Token:    options.token,
		Database: options.database,
	})

	if err != nil {
		return nil, err
	}

	options.logger.Debug("created influx client",
		"host", options.hostURI, "database", options.database)

	return &Client{
		client:     influxClient,
		httpClient: httpClient,
		options:    options,
	}, nil
}

func (c *Client) Insert(ctx context.Context, rates ...domain.Rate) error {
	ratesToWrite := make([]any, len(rates))
	for i, rate := range rates {
		ratesToWrite[i] = rate
	}

	err := c.client.WriteData(ctx, ratesToWrite)

	return err
}

func (c *Client) HealthCheck(ctx context.Context) error {
	resp, err := c.httpClient.Do(ctx,
		rest.WithURI(c.options.hostURI+"/health"),
		rest.WithMethod(http2.MethodGet),
	)

	if err != nil {
		return err
	}
	if resp.StatusCode != http2.StatusOK {
		return fmt.Errorf("unexpected status code (%d) from %s", resp.StatusCode, resp.Request.URL)
	}

	resp, err = c.httpClient.Do(ctx,
		rest.WithURI(c.options.hostURI+"/ping"),
		rest.WithMethod(http2.MethodGet),
	)

	if err != nil {
		return err
	}
	if resp.StatusCode != http2.StatusOK {
		return fmt.Errorf("unexpected status code (%d) from %s", resp.StatusCode, resp.Request.URL)
	}

	c.options.logger.Debug("influx is up and ready")

	return nil
}

func (c *Client) Get(ctx context.Context, currecny string, limit uint32) ([]domain.Rate, error) {
	if !slices.Contains(c.options.currencies, currecny) {
		return nil, ErrUnavailableCurrency
	}

	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC LIMIT %d`, strings.ToUpper(currecny), limit)

	response, err := c.client.Query(ctx, query, influx.WithQueryType(influx.InfluxQL))
	if err != nil {
		return nil, err
	}

	rates := make([]domain.Rate, 0, limit)
	for i := 0; response.Next(); i++ {
		v := response.Value()

		rate, err := parseMapToRate(v)
		if err != nil {
			return nil, err
		}

		rates = append(rates, *rate)
	}

	return rates, nil
}

func (c *Client) GetAll(ctx context.Context, limit uint32) (map[string][]domain.Rate, error) {
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC LIMIT %d`,
		strings.Join(c.options.currencies, ","), limit)

	response, err := c.client.Query(ctx, query, influx.WithQueryType(influx.InfluxQL))
	if err != nil {
		return nil, err
	}

	ratesMap := make(map[string][]domain.Rate, len(c.options.currencies))
	for i := 0; response.Next(); i++ {
		v := response.Value()

		rate, err := parseMapToRate(v)
		if err != nil {
			return nil, err
		}

		rates, ok := ratesMap[rate.Currency]
		if !ok {
			rates = make([]domain.Rate, 0, limit)
		}

		rates = append(rates, *rate)
		ratesMap[rate.Currency] = rates
	}

	return ratesMap, nil
}

func parseMapToRate(m map[string]any) (*domain.Rate, error) {
	currency, ok := m["iox::measurement"].(string)
	if !ok {
		return nil, errors.Wrap(ErrFailedToParseData, "failed to parse currency")
	}

	timestamp, ok := m["time"].(time.Time)
	if !ok {
		return nil, errors.Wrap(ErrFailedToParseData, "failed to parse timestamp")
	}

	value, ok := m["value"].(float64)
	if !ok {
		return nil, errors.Wrap(ErrFailedToParseData, "failed to parse value")
	}

	quote, ok := m["quote"].(string)
	if !ok {
		return nil, errors.Wrap(ErrFailedToParseData, "failed to parse quote")
	}

	provider, ok := m["provider"].(string)
	if !ok {
		return nil, errors.Wrap(ErrFailedToParseData, "failed to parse provider")
	}

	return &domain.Rate{
		Currency:  currency,
		Value:     value,
		Quote:     quote,
		Timestamp: timestamp,
		Provider:  provider,
	}, nil
}
