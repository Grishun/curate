package grpc

import (
	"github.com/Grishun/curate/internal/service"
	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ generated.RatesServiceServer = (*Handler)(nil)

func (h *Handler) SubscribeRate(request *generated.SubscribeRateRequest,
	stream grpc.ServerStreamingServer[generated.SubscribeRateResponse]) error {
	userID, err := uuid.Parse(request.UserId)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid user id")
	}

	subscription := service.Subscription{
		Currency: request.Currency,
		Provider: request.Provider,
		ID:       userID,
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
				UserId:     request.UserId,
				ReceivedAt: rate.Timestamp.Unix(),
			})

			if err != nil {
				h.options.logger.Error("failed to send rate to client", "error", err, "clientId", request.UserId)
			}
		}
	}
}
