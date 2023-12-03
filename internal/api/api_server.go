package api

import (
	"chatmerger/internal/domain/usecase"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var _ BaseServiceServer = (*Server)(nil)

type Server struct {
	usecase.Usecases
	cfg        Config
	grpcServer *grpc.Server
}

type Config struct {
	Host string
	Port int
}

func NewApiServer(usecases usecase.Usecases, cfg Config) *Server {
	var server = &Server{
		Usecases: usecases,
		cfg:      cfg,
	}
	var opts []grpc.ServerOption
	server.grpcServer = grpc.NewServer(opts...)
	RegisterBaseServiceServer(server.grpcServer, server)

	return server
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err = s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server grpc serfver: %v", err)
	}
	return nil
}
