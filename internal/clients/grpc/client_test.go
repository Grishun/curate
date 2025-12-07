package grpc

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/transport/grpc"
	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"github.com/stretchr/testify/require"
)

func TestGRPCClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	svc := service.New(service.WithPollingInterval(time.Second))
	go svc.Start(ctx)
	defer svc.Stop(ctx)

	server := grpc.NewServer(
		grpc.WithHost("localhost"),
		grpc.WithPort("8081"),
		grpc.WithService(svc),
	)

	go server.Run(ctx)
	defer server.Stop()

	subscribeRequests := []struct {
		request          *generated.SubscribeRateRequest
		expectedMessages int
		expectedError    bool
		name             string
	}{
		{
			request: &generated.SubscribeRateRequest{
				Currency: "BTC",
			},
			expectedMessages: 3,
			expectedError:    false,
			name:             "valid-BTC",
		},

		{
			request: &generated.SubscribeRateRequest{
				Currency: "ETH",
			},
			expectedMessages: 3,
			expectedError:    false,
			name:             "valid-ETH",
		},

		{
			request: &generated.SubscribeRateRequest{
				Currency: "USD",
			},
			expectedMessages: 0,
			expectedError:    false,
			name:             "invalid-currency",
		},

		{
			request: &generated.SubscribeRateRequest{
				Currency: "BTC",
				UserId:   "cecwzxa",
			},
			expectedMessages: 0,
			expectedError:    true,
			name:             "invlid-user-id",
		},
	}

	time.Sleep(time.Second * 2)

	var counter atomic.Int32
	for _, req := range subscribeRequests {
		req := req
		t.Run(req.name, func(t *testing.T) {
			defer func() {
				if int(counter.Add(1)) == len(subscribeRequests) {
					cancel()
				}
			}()

			client, err := NewClient(WithServerAddr(server.Address()))
			require.NoError(t, err)

			if req.request.UserId == "" {
				req.request.UserId = client.uuid.String()
			}

			stream, err := client.grpcClient.SubscribeRate(ctx, req.request)
			require.NoError(t, err)

			for i := 0; i < req.expectedMessages; i++ {
				response, err := stream.Recv()
				t.Log("received new data from stream", "response", response, "error", err)
				if req.expectedError {
					require.Error(t, err)
					return
				}
				require.NotNil(t, response)
				require.NotZero(t, response.Rate)
				require.NotZero(t, response.ReceivedAt)
			}
		})
	}

	t.Log("stopping server")
}
