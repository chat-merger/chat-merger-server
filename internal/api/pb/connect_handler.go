package pb

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/rule"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
)

func (s *Server) Connect(connService BaseService_ConnectServer) error {
	var md = parseConnMetaData(connService.Context())
	var input = model.CreateClientSession{ApiKey: md.ApiKey}
	client, err := s.CreateClientSession(input)
	if err != nil {
		return fmt.Errorf("failed create client session: %s", err)
	}
	//  when need send msg from some other clients to current connect
	go func() {
		select {
		case msg, ok := <-client.MsgCh:
			if !ok {
				log.Printf("failed to read channel of client %#v", client)
			} else {
				response, err := MessageToResponse(msg)
				if err != nil {
					log.Printf("failed convert msg to respponse: %s", err)
				}
				err = connService.Send(response)
				if err != nil {
					log.Printf("failed send response to client (%s): %s", client.Name, err)
				}
			}
		}
	}()
	// read clients input msg and fanout to other clients
	for {
		r, err := connService.Recv()
		if err != nil {
			s.DropClientSession([]model.ID{client.Id})
			if err == io.EOF {
				log.Println("err == io.EOF")
				return nil
			}
			log.Printf("recv op err: %v", err)
			return err
		}
		//resp, err := transform(r, client.Name)
		msg, err := RequestToCreateMessage(r, client.Name)
		if err != nil {
			log.Printf("transform request to response: %v", err)
			continue
		}
		err = s.CreateAndSendMsgToEveryoneExcept(*msg, []model.ID{client.Id})
		if err != nil {
			log.Printf("send msg to clients: %v", err)
		}
	}
}

func (s *Server) createSession(input model.CreateClientSession) {
	s.CreateClientSession(input)
}

type metaData struct {
	ApiKey model.ApiKey
}

func parseConnMetaData(ctx context.Context) metaData {
	var md, _ = metadata.FromIncomingContext(ctx)
	var apiKeyRaw = md.Get(rule.AUTHENTICATE_HEADER)
	var apiKey model.ApiKey
	if len(apiKeyRaw) > 0 {
		apiKey = model.NewApiKey(apiKeyRaw[0])
	}
	return metaData{
		ApiKey: apiKey,
	}
}

//
//func transform(request *Request, client string) (*Response, error) {
//	var body isResponse_Body
//	switch request.Body.(type) {
//	case *Request_Text:
//		body = &Response_Text{
//			Text: request.Body.(*Request_Text).Text,
//		}
//	case *Request_Media:
//		body = &Response_Media{
//			Media: request.Body.(*Request_Media).Media,
//		}
//	default:
//		return nil, errors.New("request body not match with RequestBody interface")
//	}
//	var resp = &Response{
//		Id:         uuid.NewString(),
//		ReplyMsgId: request.ReplyMsgId,
//		CreatedAt:  request.CreatedAt,
//		Author:     request.Author,
//		Client:     client,
//		IsSilent:   request.IsSilent,
//		Body:       body,
//	}
//	return resp, nil
//}
