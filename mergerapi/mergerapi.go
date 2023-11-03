package mergerapi

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Server struct {
}

func (s Server) Connect(server BaseService_ConnectServer) error {
	for {
		counter = counter + 1

		newMsg := &NewMessageEvent{
			Id: fmt.Sprintf("%d", counter),
			CreatedAt: &timestamppb.Timestamp{
				Seconds: time.Now().UnixMilli() / 1000,
			},
			ClientName: "golang server",
			Author: &NewMessageEvent_AuthorValue{
				AuthorValue: &Author{
					AuthorId:   "42",
					AuthorName: "saime",
				},
			},
			Modifiers:   make([]*Modifier, 0),
			Action:      &NewMessageEvent_ActionValue{ActionValue: 1},
			Attachments: make([]*Attachment, 0),
		}

		err := server.Send(&Response{
			Event: &Response_OnCreateMsg{
				OnCreateMsg: newMsg,
			},
		})

		if err != nil {
			log.Printf("failed sending messages: %v", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s Server) mustEmbedUnimplementedBaseServiceServer() {
	//TODO implement me
	panic("implement me")
}

var counter = 0
