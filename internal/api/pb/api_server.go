package pb

import (
	"chatmerger/internal/usecase"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

var _ BaseServiceServer = (*Server)(nil)

type Server struct {
	cfg        Config
	grpcServer *grpc.Server
	Usecases
}

type Config struct {
	Host string
	Port int
}

type Usecases struct {
	usecase.SendMessageToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
}

func NewApiServer(cfg Config, usecases Usecases) *Server {
	var server = &Server{
		cfg:      cfg,
		Usecases: usecases,
	}
	var opts []grpc.ServerOption

	server.grpcServer = grpc.NewServer(opts...)
	RegisterBaseServiceServer(server.grpcServer, server)

	return server
}

func (s *Server) Connect(connService BaseService_ConnectServer) error {
	var data = connService.Context()
	log.Printf("connected %v", data)

	for {
		r, err := connService.Recv()
		if err != nil {
			log.Fatalf("recv op err: %v", err)
		}
		resp := Response{
			Author: r.Author,
		}
		connService.Send(&resp)
	}
	return nil

}

func (s *Server) mustEmbedUnimplementedBaseServiceServer() {}

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
