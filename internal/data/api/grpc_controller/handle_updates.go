package grpc_controller

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/component/eventbus"
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

	ctx, cancel := context.WithCancel(rpcCall.Context())

	err = s.SubscribeClientToEvents(client.Id, eventHandler(rpcCall, cancel))
	if err != nil {
		return fmt.Errorf("failed subscribe to events: %s", err)
	}
	log.Println(msgs.ClientSubscribedToNewMsgs)

	select {
	case <-ctx.Done():
		s.DropClientSubscription(client.Id)
		return nil
	}
}

type metaData struct {
	ApiKey model.ApiKey
}

const (
	authenticateHeader = "X-Api-Key"
)

func parseConnMetaData(ctx context.Context) metaData {
	var md, _ = metadata.FromIncomingContext(ctx)
	var apiKeyRaw = md.Get(authenticateHeader)
	var apiKey model.ApiKey
	if len(apiKeyRaw) > 0 {
		apiKey = model.ApiKey(apiKeyRaw[0])
	}
	return metaData{
		ApiKey: apiKey,
	}
}

func eventHandler(rpcCall pb.BaseService_UpdatesServer, cancel context.CancelFunc) eventbus.Handler {
	return func(event eventbus.Event) error {
		switch {

		case event.Message != nil:
			response, err := messageToResponse(*event.Message)
			if err != nil {
				return fmt.Errorf("failed convert msg to response: %s\n", err)
			}
			err = rpcCall.Send(response)
			if err != nil {
				return fmt.Errorf("failed send response: %s\n", err)
			}

		case event.DropSubscription != nil:
			cancel()
		}
		return nil
	}
}
