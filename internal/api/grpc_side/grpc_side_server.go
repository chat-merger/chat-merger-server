package grpc_side

import (
	"chatmerger/internal/api/pb"
	"chatmerger/internal/usecase"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

var _ pb.BaseServiceServer = (*GrpcSideServer)(nil)

type GrpcSideServer struct {
	cfg        Config
	grpcServer *grpc.Server
	requiredUsecases
	pb.UnimplementedBaseServiceServer
}

type Config struct {
	Host string
	Port int
}

type requiredUsecases interface {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
}

func NewGrpcSideServer(cfg Config, usecases requiredUsecases) *GrpcSideServer {
	var server = &GrpcSideServer{
		cfg:              cfg,
		requiredUsecases: usecases,
	}
	var opts []grpc.ServerOption

	server.grpcServer = grpc.NewServer(opts...)
	pb.RegisterBaseServiceServer(server.grpcServer, server)

	return server
}

func (s *GrpcSideServer) Serve(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port))
	if err != nil {
		return fmt.Errorf("create listener: %v", err)
	}
	go ctxHandler(ctx, s.grpcServer.Stop)

	if err = s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to server grpc server: %v", err)
	}
	return nil
}

func ctxHandler(ctx context.Context, callback func()) {
	select {
	case <-ctx.Done():
		callback()
	}
}
