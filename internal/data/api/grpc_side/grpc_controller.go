package grpc_side

import (
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/usecase"
	"google.golang.org/grpc"
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

func NewGrpcController(cfg Config, usecases requiredUsecases) *GrpcSideServer {
	var server = &GrpcSideServer{
		cfg:              cfg,
		requiredUsecases: usecases,
	}
	var opts []grpc.ServerOption

	server.grpcServer = grpc.NewServer(opts...)
	pb.RegisterBaseServiceServer(server.grpcServer, server)

	return server
}
