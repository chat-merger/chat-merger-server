package grpc_controller

import (
	"context"
	"fmt"
	"net"
)

func (s *GrpcController) Run(ctx context.Context) error {
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
