package grpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Grishun/curate/internal/provider/coindesk"
	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/transport/grpc"
	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"github.com/stretchr/testify/require"
)

func TestGRPCClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	mockedCoinDesk := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"BTC":{"USD":91194.7},"ETH":{"USD":3050.22},"TRX":{"USD":0.2813}}`))
	}))
	defer mockedCoinDesk.Close()
	svc := service.New(
		service.WithPollingInterval(time.Second),
		service.WithProviders(coindesk.New(coindesk.WithURI(mockedCoinDesk.URL))),
	)

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
			name:             "invalid-USD",
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
