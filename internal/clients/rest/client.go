package rest

import (
	"context"
	"net/http"

	"github.com/Grishun/curate/internal/domain"
	"resty.dev/v3"
)

type Client struct {
	*resty.Client
	logger domain.Logger
}

func NewClient(opts ...ClientOption) *Client {
	options := NewClientOptions()

	for _, opt := range opts {
		opt(options)
	}
	restyClient := resty.New()

	restyClient.
		SetBaseURL(options.baseURI).
		SetAuthToken(options.token).
		SetLogger(options.logger).
		SetAuthScheme(options.authScheme).
		SetTimeout(options.timeout)

	return &Client{Client: restyClient, logger: options.logger}
}

func (c *Client) Do(ctx context.Context, opts ...domain.RequestOption) (*http.Response, error) {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	c.logger.Debug("sending rest request", "method", options.Method, "uri", options.URI)

	resp, err := c.R().SetContext(ctx).
		SetHeaderMultiValues(options.Headers).
		SetQueryParamsFromValues(options.QueryParams).
		SetBody(options.Body).
		SetResult(options.UnmarshallTo).
		Execute(options.Method, options.URI)

	if err != nil {
		c.logger.Error("rest request failed", "method", options.Method, "uri", options.URI, "error", err)
		return nil, err
	}

	c.logger.Debug("rest request completed", "method", options.Method, "uri", options.URI, "status", resp.StatusCode())

	return resp.RawResponse, nil
}
