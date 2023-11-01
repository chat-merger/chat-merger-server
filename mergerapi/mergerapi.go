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

// создается на каждого клиента
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
			Modifiers:   make([]*Modifier, 0),
			Action:      &MsgBody_ActionValue{ActionValue: 1},
			Attachments: make([]*Attachment, 0),
		}
		err := server.Send(&newMsg)
		if err != nil {
			log.Printf("failed sending messages: %v", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func clientInputReader(server BaseService_CreateMessageServer) {

}

func newMsgCreator(server BaseService_CreateMessageServer) {

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
