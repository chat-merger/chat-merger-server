package grpc_side

import (
	"chatmerger/internal/api/pb"
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/common/vals"
	"chatmerger/internal/domain/model"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
)

func (s *GrpcSideServer) Connect(rpcCall pb.BaseService_ConnectServer) error {
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
	// read clients input msg and fanout to other clients
	go conn.handleMsgsFromRpcCall(s.requiredUsecases)

	select {
	case <-rpcCall.Context().Done():
		return nil
	case err := <-onErr:
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
