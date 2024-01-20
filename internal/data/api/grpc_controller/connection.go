package grpc_controller

import (
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"fmt"
)

type connection struct {
	rpcCall pb.BaseService_UpdatesServer
	session *model.ClientSession
	errCh   chan<- error
}

func newConnection(rpcCall pb.BaseService_UpdatesServer, session *model.ClientSession) (*connection, <-chan error) {
	errCh := make(chan error)
	return &connection{
		rpcCall: rpcCall,
		session: session,
		errCh:   errCh,
	}, errCh
}

func (c *connection) errorf(format string, a ...any) {
	select {
	case c.errCh <- fmt.Errorf(format, a):
	default:
	}
}
