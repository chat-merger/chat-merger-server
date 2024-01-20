package grpc_controller

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"context"
	"errors"
	"fmt"
	"log"
)

func (s *GrpcController) SendMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	var md = parseConnMetaData(ctx)
	log.Printf("%s: %+v", msgs.NewMessageFromClient, md)
	sessions, err := s.Clients(usecase.ClientsFilter{ApiKey: &md.ApiKey})
	if err != nil {
		return nil, fmt.Errorf("failed create session session: %s", err)
	}
	if len(sessions) == 0 {
		return nil, errors.New("not found client")
	}
	session := sessions[0]
	msg, err := requestToCreateMessage(req, session.Name)
	if err != nil {
		return nil, fmt.Errorf("transform request to response: %v\n", err)
	}
	newMsg, err := s.CreateAndSendMsgToEveryoneExcept(*msg, []model.ID{session.Id})
	if err != nil {
		return nil, fmt.Errorf("send msg to clients: %v\n", err)
	}
	resp, err := messageToResponse(*newMsg)
	if err != nil {
		return nil, fmt.Errorf("failed convert msg to response: %v\n", err)
	}
	return resp, nil
}
