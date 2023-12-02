package api

import (
	"chatmerger/pkg/mergerapi"
)

var _ mergerapi.BaseServiceServer = (*Server)(nil)

type Server struct {
}

func (s *Server) Connect(connService mergerapi.BaseService_ConnectServer) error {
	return nil
}

// func (s *Server) mustEmbedUnimplementedBaseServiceServer() {}
