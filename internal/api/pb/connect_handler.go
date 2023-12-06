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
	log.Printf("meta: %#v", md)
	var input = model.CreateClientSession{ApiKey: md.ApiKey}
	client, err := s.CreateClientSession(input)
	if err != nil {
		return fmt.Errorf("failed create client session: %s", err)
	}
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

		resp := Response{
			Author: r.Author,
		}
		connService.Send(&resp)
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
