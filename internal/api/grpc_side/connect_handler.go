package grpc_side

import (
	"chatmerger/internal/api/pb"
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/common/vals"
	"chatmerger/internal/domain/model"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

func (s *GrpcSideServer) Connect(connService pb.BaseService_ConnectServer) error {
	//return errors.ErrUnsupported
	var md = parseConnMetaData(connService.Context())
	log.Printf("%s: %+v", msgs.ClientConnectedToServer, md)
	var input = model.CreateClientSession{ApiKey: md.ApiKey}
	clientSession, err := s.CreateClientSession(input)
	if err != nil {
		return fmt.Errorf("failed create clientSession session: %s", err)
	}
	errCh := make(chan error)
	log.Printf("%s: %+v", msgs.ClientSessionCreated, clientSession)
	//  when need send msg from some other clients to current connect
	go handleSessionReceivedMsgs(connService, clientSession, errCh)
	// read clients input msg and fanout to other clients
	go handleMessagesFromClient(connService, clientSession, errCh, s.requiredUsecases)

	select {
	case <-connService.Context().Done():
		return nil
	case err := <-errCh:
		return err
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

func handleMessagesFromClient(connService pb.BaseService_ConnectServer, clientSession *model.ClientSession, errCh chan<- error, usecases requiredUsecases) {
	for {
		r, err := connService.Recv()
		if err != nil {
			usecases.DropClientSession([]model.ID{clientSession.Id})
			if err == io.EOF {
				log.Println(msgs.ClientSessionCloseConnection)
				return
			}

			errCh <- fmt.Errorf("recv op err: %v\n", err)
			return
		}
		log.Println(msgs.NewMessageFromClient)
		msg, err := requestToCreateMessage(r, clientSession.Name)
		if err != nil {
			errCh <- fmt.Errorf("transform request to response: %v\n", err)
			continue
		}
		err = usecases.CreateAndSendMsgToEveryoneExcept(*msg, []model.ID{clientSession.Id})
		if err != nil {
			errCh <- fmt.Errorf("send msg to clients: %v\n", err)
		}
	}
}

func handleSessionReceivedMsgs(connService pb.BaseService_ConnectServer, session *model.ClientSession, errCh chan<- error) {
	for {
		select {
		case <-connService.Context().Done():
			return
		case msg, ok := <-session.MsgCh:
			if !ok {
				errCh <- fmt.Errorf("failed to read channel of client session %#v\n", session)
				return
			} else {
				response, err := messageToResponse(msg)
				if err != nil {
					errCh <- fmt.Errorf("failed convert msg to respponse: %s\n", err)
					return
				}
				err = connService.Send(response)
				if err != nil {
					errCh <- fmt.Errorf("failed send response to client session (%s): %s\n", session.Name, err)
					return
				}
			}
		}
	}
}
