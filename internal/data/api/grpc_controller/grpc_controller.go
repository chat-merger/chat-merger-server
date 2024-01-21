package grpc_controller

import (
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/usecase"
	"google.golang.org/grpc"
)

var _ pb.BaseServiceServer = (*GrpcController)(nil)

type GrpcController struct {
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
	usecase.SubscribeClientToEventsUc
	usecase.DropClientSubscriptionUc
	usecase.ClientsUc
}

func NewGrpcController(cfg Config, usecases requiredUsecases) *GrpcController {
	var server = &GrpcController{
		cfg:              cfg,
		requiredUsecases: usecases,
	}
	var opts []grpc.ServerOption

	server.grpcServer = grpc.NewServer(opts...)
	pb.RegisterBaseServiceServer(server.grpcServer, server)

	return server
}
