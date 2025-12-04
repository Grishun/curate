package influx

import (
	"context"
	"fmt"
	http2 "net/http"
	"strings"
	"time"

	"github.com/Grishun/curate/internal/clients/rest"
	"github.com/Grishun/curate/internal/domain"
	influx "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type Client struct {
	client     *influx.Client
	options    *Options
	httpClient domain.HTTPClient
}

var (
	ErrFailedToParseData = errors.New("error parsing data from influx map")
	ErrInfluxNotReady    = errors.New("influx is not ready to receive quiries yet")
)

func NewClient(opts ...Option) (*Client, error) {
	options := NewOptions()

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
		options.logger.Error("failed to create influx client", "error", err)
		return nil, err
	}

	options.logger.Info("created influx client",
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
	if err != nil {
		c.options.logger.Error("failed to write rates to influx", "count", len(ratesToWrite), "error", err)
		return err
	}

	c.options.logger.Info("inserted rates to influx", "count", len(ratesToWrite))

	return nil
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
	c.options.logger.Debug("querying influx for currency", "currency", currecny, "limit", limit)

	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC LIMIT %d`, strings.ToUpper(currecny), limit)

	response, err := c.client.Query(ctx, query, influx.WithQueryType(influx.InfluxQL))
	if err != nil {
		c.options.logger.Error("influx query failed", "currency", currecny, "limit", limit, "error", err)
		return nil, err
	}

	rates := make([]domain.Rate, 0, limit)
	for i := 0; response.Next(); i++ {
		v := response.Value()

		rate, err := c.parseMapToRate(v)
		if err != nil {
			c.options.logger.Error("failed to parse influx row", "currency", currecny, "error", err)
			return nil, err
		}

		rates = append(rates, *rate)
	}

	rates = slices.Clip(rates)
	slices.Reverse(rates)

	c.options.logger.Info("fetched rates from influx", "currency", currecny, "count", len(rates))

	return rates, nil
}

func (c *Client) GetAll(ctx context.Context, limit uint32) (map[string][]domain.Rate, error) {
	c.options.logger.Debug("querying influx for all currencies", "limit", limit)

	query := fmt.Sprintf(`SELECT * FROM /.*/ ORDER BY time DESC LIMIT %d`, limit)

	response, err := c.client.Query(ctx, query, influx.WithQueryType(influx.InfluxQL))
	if err != nil {
		c.options.logger.Error("influx multi-currency query failed", "limit", limit, "error", err)
		return nil, err
	}

	ratesMap := make(map[string][]domain.Rate)
	for i := 0; response.Next(); i++ {
		v := response.Value()

		rate, err := c.parseMapToRate(v)
		if err != nil {
			c.options.logger.Error("failed to parse influx row", "error", err)
			return nil, err
		}

		rates, ok := ratesMap[rate.Currency]
		if !ok {
			rates = make([]domain.Rate, 0, limit)
		}

		rates = append(rates, *rate)
		ratesMap[rate.Currency] = rates
	}

	for currency, rates := range ratesMap {
		slices.Reverse(rates)
		ratesMap[currency] = slices.Clip(rates)
	}

	c.options.logger.Info("fetched rates for all currencies", "limit", limit, "currency_count", len(ratesMap))

	return ratesMap, nil
}

func (c *Client) parseMapToRate(m map[string]any) (*domain.Rate, error) {
	currency, ok := m["iox::measurement"].(string)
	if !ok {
		err := errors.Wrap(ErrFailedToParseData, "failed to parse currency")
		c.options.logger.Error("failed to parse influx currency", "error", err)
		return nil, err
	}

	timestamp, ok := m["time"].(time.Time)
	if !ok {
		err := errors.Wrap(ErrFailedToParseData, "failed to parse timestamp")
		c.options.logger.Error("failed to parse influx timestamp", "currency", currency, "error", err)
		return nil, err
	}

	value, ok := m["value"].(float64)
	if !ok {
		err := errors.Wrap(ErrFailedToParseData, "failed to parse value")
		c.options.logger.Error("failed to parse influx value", "currency", currency, "error", err)
		return nil, err
	}

	quote, ok := m["quote"].(string)
	if !ok {
		err := errors.Wrap(ErrFailedToParseData, "failed to parse quote")
		c.options.logger.Error("failed to parse influx quote", "currency", currency, "error", err)
		return nil, err
	}

	provider, ok := m["provider"].(string)
	if !ok {
		err := errors.Wrap(ErrFailedToParseData, "failed to parse provider")
		c.options.logger.Error("failed to parse influx provider", "currency", currency, "error", err)
		return nil, err
	}

	return &domain.Rate{
		Currency:  currency,
		Value:     value,
		Quote:     quote,
		Timestamp: timestamp,
		Provider:  provider,
	}, nil
}
