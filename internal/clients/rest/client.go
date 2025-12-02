package rest

import (
	"net/http"

	"github.com/Grishun/curate/internal/domain"
	"resty.dev/v3"
)

type Client struct {
	*resty.Client
}

func NewClient(opts ...ClientOption) *Client {
	options := NewClientOptions()

	for _, opt := range opts {
		opt(options)
	}
	restyClient := resty.New()

	restyClient.
		SetContext(options.ctx).
		SetBaseURL(options.baseURI).
		SetAuthToken(options.token).
		SetLogger(options.logger).
		SetAuthScheme(options.authScheme).
		SetTimeout(options.timeout)

	return &Client{restyClient}
}

func (c *Client) Do(opts ...domain.RequestOption) (*http.Response, error) {
	options := NewOptions()

	for _, opt := range opts {
		opt(options)
	}

	resp, err := c.R().SetContext(options.Ctx).
		SetHeaderMultiValues(options.Headers).
		SetQueryParamsFromValues(options.QueryParams).
		SetBody(options.Body).
		SetResult(options.UnmarshallTo).
		Execute(options.Method, options.URI)

	if err != nil {
		return nil, err
	}

	return resp.RawResponse, nil
}
