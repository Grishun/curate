package grpc

import (
	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient generated.RatesServiceClient
	options    *ClientOptions
	conn       *grpc.ClientConn
}

func NewClient(opts ...ClientOption) (*Client, error) {
	options := NewClientOptions()

	for _, opt := range opts {
		opt(options)
	}

	conn, err := grpc.NewClient(options.serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		options.logger.Error("failed to create gRPC client", err, err.Error())
		return nil, err
	}

	options.logger.Info("created gRPC client", "serverAddr", options.serverAddr)

	return &Client{
		grpcClient: generated.NewRatesServiceClient(conn),
		options:    options,
		conn:       conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
