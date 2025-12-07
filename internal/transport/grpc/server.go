package grpc

import (
	"context"
	"net"

	"github.com/Grishun/curate/internal/transport/grpc/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	options    *ServerOptions
	grpcServer *grpc.Server
}

func NewServer(opts ...ServerOption) *Server {
	options := NewServerOptions()

	for _, opt := range opts {
		opt(options)
	}

	options.logger.Info("create new grpc server", "host", options.host, "port", options.port)
	return &Server{
		options: options,
		grpcServer: grpc.NewServer(
			grpc.Creds(insecure.NewCredentials())),
	}
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.options.host, s.options.port))
	if err != nil {
		s.options.logger.Error("failed to listen", "error", err)
		return err
	}

	handler := NewHandler(
		WithHandlerLogger(s.options.logger),
		WithHandlerService(s.options.service),
	)
	generated.RegisterRatesServiceServer(s.grpcServer, handler)
	reflection.Register(s.grpcServer) // add reflection to the server

	s.options.logger.Info("registered handlers for grpc. Starting server")

	errCh := make(chan error)
	go func(chan error) {
		errCh <- s.grpcServer.Serve(lis)
	}(errCh)

	select {
	case <-ctx.Done():
		s.Stop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop() {
	s.options.logger.Info("calling for grpcServer.GracefulStop()")
	s.grpcServer.GracefulStop()
}

func (s *Server) Address() string {
	return net.JoinHostPort(s.options.host, s.options.port)
}
