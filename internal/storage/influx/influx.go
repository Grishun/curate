package influx

import (
	"context"
	"fmt"
	http2 "net/http"
	"strings"

	"github.com/Grishun/curate/internal/clients/rest"
	"github.com/Grishun/curate/internal/domain"
	influx "github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

type Client struct {
	client     *influx.Client
	options    *Options
	httpClient domain.Client
}

func NewClient(opts ...Option) (*Client, error) {
	options := NewOptions(opts...)

	for _, opt := range opts {
		opt(options)
	}

	httpClient := rest.NewClient(
		rest.WithToken(options.token),
		rest.WithContext(options.ctx),
		rest.WithAuthScheme("Bearer"),
	)

	influxClient, err := influx.New(influx.ClientConfig{
		Host:     options.hostURI,
		Token:    options.token,
		Database: options.database,
	})

	if err != nil {
		options.logger.Error("failed to create influx client",
			"error", err, "host", options.hostURI, "database", options.database, "token", options.token[:5]+"...")
	}

	options.logger.Debug("created influx client",
		"host", options.hostURI, "database", options.database, "token", options.token[:5]+"...")

	return &Client{
		client:     influxClient,
		httpClient: httpClient,
		options:    options}, nil
}

func (c *Client) Insert(ctx context.Context, rates ...domain.Rate) error {
	ratesToWrite := make([]any, len(rates))
	for i, rate := range rates {
		ratesToWrite[i] = rate
	}

	return c.client.WriteData(ctx, ratesToWrite)
}

func (c *Client) HealthCheck(ctx context.Context) error {
	resp, err := c.httpClient.Do(
		rest.WithURI(c.options.hostURI+"/health"),
		rest.WithMethod(http2.MethodGet),
		rest.WithRequestContext(ctx),
	)

	if resp.StatusCode != http2.StatusOK {
		return fmt.Errorf("non statusok response from influx (%d)", resp.StatusCode)
	}

	if err != nil {
		c.options.logger.Error("influx is not responding", "error", err)
		return err
	}

	c.options.logger.Debug("influx is up and ready")

	return nil
}

func (c *Client) Get(ctx context.Context, currecny string, limit uint) ([]domain.Rate, error) {
	query := fmt.Sprintf(`SELECT * FROM %s ORDER BY time DESC LIMIT %d`, strings.ToUpper(currecny), limit)

	response, err := c.client.Query(ctx, query, influx.WithQueryType(influx.InfluxQL))
	if err != nil {
		return nil, err
	}

	response.Value()

	return nil, nil
}

func (c *Client) GetAll(ctx context.Context, limit uint) (map[string][]domain.Rate, error) {
	return nil, nil
}

func (c *Client) GetHistoryLimit() uint {
	return 0
}
