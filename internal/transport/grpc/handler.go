package grpc

import (
	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ generated.RatesServiceServer = (*Handler)(nil)

func (h *Handler) SubscribeRate(request *generated.SubscribeRateRequest,
	stream grpc.ServerStreamingServer[generated.SubscribeRateResponse]) error {

	subscription := service.Subscription{
		Currency: request.Currency,
		Provider: request.Provider,
	}

	ratesCh := h.options.service.SubscribeRate(stream.Context(), subscription)

	for {
		select {
		case <-stream.Context().Done():
			h.options.logger.Info("stream context is done, closing channel")
			return status.Error(codes.Canceled, "stream context is done")
		case rate := <-ratesCh:
			err := stream.Send(&generated.SubscribeRateResponse{
				Rate:       rate.Value,
				ReceivedAt: rate.Timestamp.Unix(),
			})

			if err != nil {
				h.options.logger.Error("failed to send rate to client", "error", err)
			}
		}
	}
}
