package grpc_side

import (
	"chatmerger/internal/api/pb"
	"chatmerger/internal/common/msgs"
	"chatmerger/internal/domain/model"
	"fmt"
	"io"
	"log"
)

type connection struct {
	rpcCall pb.BaseService_ConnectServer
	session *model.ClientSession
	errCh   chan<- error
}

func newConnection(rpcCall pb.BaseService_ConnectServer, session *model.ClientSession) (*connection, <-chan error) {
	errCh := make(chan error)
	return &connection{
		rpcCall: rpcCall,
		session: session,
		errCh:   errCh,
	}, errCh
}

func (c *connection) errorf(format string, a ...any) {
	select {
	case c.errCh <- fmt.Errorf(format, a):
	default:
	}
}

func (c *connection) handleMsgsFromRpcCall(usecases requiredUsecases) {
	for {
		r, err := c.rpcCall.Recv()
		if err != nil {
			usecases.DropClientSession([]model.ID{c.session.Id})
			if err == io.EOF {
				log.Println(msgs.ClientSessionCloseConnection)
				return
			}

			c.errCh <- fmt.Errorf("recv op errCh: %v\n", err)
			return
		}
		log.Println(msgs.NewMessageFromClient)
		msg, err := requestToCreateMessage(r, c.session.Name)
		if err != nil {
			c.errCh <- fmt.Errorf("transform request to response: %v\n", err)
			continue
		}
		err = usecases.CreateAndSendMsgToEveryoneExcept(*msg, []model.ID{c.session.Id})
		if err != nil {
			c.errCh <- fmt.Errorf("send msg to clients: %v\n", err)
		}
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
