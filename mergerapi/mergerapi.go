package mergerapi

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Server struct {
}

var counter = 0

func (s Server) CreateMessage(server BaseService_CreateMessageServer) error {
	for {
		counter = counter + 1

		newMsg := MsgBody{
			Id: fmt.Sprintf("%d", counter),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().UnixMilli() / 1000,
			},
			ClientName: "golang server",
			Author: &MsgBody_AuthorValue{
				AuthorValue: &Author{
					AuthorId:   "42",
					AuthorName: "saime",
				},
			},
			Modifiers: make([]*Modifier, 0),
			Action: &MsgBody_ActionValue{
				ActionValue: 1,
			},
			Attachments: make([]*Attachment, 0),
		}
		err := server.Send(&newMsg)
		if err != nil {
			log.Fatalf("failed sending messages: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s Server) Edit(server BaseService_EditServer) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) Delete(server BaseService_DeleteServer) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) mustEmbedUnimplementedBaseServiceServer() {}
