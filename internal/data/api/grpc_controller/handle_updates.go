package grpc_controller

import (
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/common/vals"
	"chatmerger/internal/data/api/pb"
	"chatmerger/internal/domain/model"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (s *GrpcController) Updates(_ *emptypb.Empty, rpcCall pb.BaseService_UpdatesServer) error {
	var md = parseConnMetaData(rpcCall.Context())
	log.Printf("%s: %+v", msgs.ClientConnectedToServer, md)
	var input = model.CreateClientSession{ApiKey: md.ApiKey}
	clientSession, err := s.CreateClientSession(input)
	if err != nil {
		return fmt.Errorf("failed create session session: %s", err)
	}
	log.Printf("%s: %+v", msgs.ClientSessionCreated, clientSession)

	conn, onErr := newConnection(rpcCall, clientSession)

	//  when need send msg from some other clients to current connect
	go conn.handleSessionReceivedMsgs()

	select {
	case <-rpcCall.Context().Done():
		s.DropClientSession([]model.ID{conn.session.Id})
		return nil
	case err := <-onErr:
		s.DropClientSession([]model.ID{conn.session.Id})
		return err // todo replace with friendly error
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

func (c *connection) handleSessionReceivedMsgs() {
	for {
		select {
		case <-c.rpcCall.Context().Done():
			return
		case msg, ok := <-c.session.MsgCh:
			if !ok {
				c.errorf("failed to read channel of client session %#v\n", c.session)
				return
			}
			response, err := messageToResponse(msg)
			if err != nil {
				c.errorf("failed convert msg to respponse: %s\n", err)
				return
			}
			err = c.rpcCall.Send(response)
			if err != nil {
				c.errorf("failed send response to client session (%s): %s\n", c.session.Name, err)
				return
			}
		}
	}
}
