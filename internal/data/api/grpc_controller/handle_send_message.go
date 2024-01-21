package grpc_controller

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"context"
	"errors"
	"fmt"
	"log"
)

func (s *GrpcController) SendMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	var md = parseConnMetaData(ctx)
	log.Printf("%s: %+v", msgs.NewMessageFromClient, md)
	clients, err := s.Clients(model.ClientsFilter{ApiKey: &md.ApiKey})
	if err != nil {
		return nil, fmt.Errorf("failed get clients: %s", err)
	}
	if len(clients) == 0 {
		return nil, errors.New("not found client")
	}
	client := clients[0]
	msg, err := requestToCreateMessage(req, client.Name)
	if err != nil {
		return nil, fmt.Errorf("transform request to response: %v\n", err)
	}
	newMsg, err := s.CreateAndSendMsgToEveryoneExcept(*msg, client.Id)
	if err != nil {
		return nil, fmt.Errorf("send msg to clients: %v\n", err)
	}
	resp, err := messageToResponse(*newMsg)
	if err != nil {
		return nil, fmt.Errorf("failed convert msg to response: %v\n", err)
	}
	return resp, nil
}
