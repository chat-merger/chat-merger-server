package api

import "log"

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
