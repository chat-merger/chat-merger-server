package grpc_controller

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/common/vals"
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *GrpcController) Updates(_ *emptypb.Empty, rpcCall pb.BaseService_UpdatesServer) error {
	var md = parseConnMetaData(rpcCall.Context())
	log.Printf("%s: %+v", msgs.ClientConnectedToServer, md)

	clients, err := s.Clients(model.ClientsFilter{ApiKey: &md.ApiKey})
	if err != nil {
		return fmt.Errorf("get clients: %s", err)
	}
	if len(clients) == 0 {
		return errors.New("invalid apikey")
	}

	client := clients[0]

	err = s.SubscribeClientToNewMsgs(client.Id, msgsEventHandler(rpcCall))
	if err != nil {
		return fmt.Errorf("failed subscribe to new msgs: %s", err)
	}
	log.Println(msgs.ClientSubscribedToNewMsgs)

	select {
	case <-rpcCall.Context().Done():
		s.DropClientSubscription(client.Id)
		return nil
	}
}

type metaData struct {
	ApiKey model.ApiKey
}

func parseConnMetaData(ctx context.Context) metaData {
	var md, _ = metadata.FromIncomingContext(ctx)
	var apiKeyRaw = md.Get(vals.AUTHENTICATE_HEADER)
	var apiKey model.ApiKey
	if len(apiKeyRaw) > 0 {
		apiKey = model.NewApiKey(apiKeyRaw[0])
	}
	return metaData{
		ApiKey: apiKey,
	}
}

func msgsEventHandler(rpcCall pb.BaseService_UpdatesServer) func(model.Message) error {
	return func(message model.Message) error {
		response, err := messageToResponse(message)
		if err != nil {
			return fmt.Errorf("failed convert msg to response: %s\n", err)
		}
		err = rpcCall.Send(response)
		if err != nil {
			return fmt.Errorf("failed send response: %s\n", err)
		}
		return nil
	}
}
