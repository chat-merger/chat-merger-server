package grpc_side

import (
	"chatmerger/internal/api/pb"
	"chatmerger/internal/usecase"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var _ pb.BaseServiceServer = (*Server)(nil)

type Server struct {
	cfg        Config
	grpcServer *grpc.Server
	Usecases
	pb.UnimplementedBaseServiceServer
}

type Config struct {
	Host string
	Port int
}

type Usecases struct {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
}

func NewClientsServer(cfg Config, usecases Usecases) *Server {
	var server = &Server{
		cfg:      cfg,
		Usecases: usecases,
	}
	var opts []grpc.ServerOption

	server.grpcServer = grpc.NewServer(opts...)
	pb.RegisterBaseServiceServer(server.grpcServer, server)

	return server
}

func (s *Server) Serve(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go s.contextCancelHandler(ctx)

	if err = s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server grpc serfver: %v", err)
	}
	return nil
}

func (s *Server) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
	}
}
