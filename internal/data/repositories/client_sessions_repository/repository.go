package client_sessions_repository

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"errors"
)

var _ domain.ClientSessionsRepository = (*ClientSessionsRepositoryBase)(nil)

type ClientSessionsRepositoryBase struct {
	conns []connect
}

func NewClientSessionsRepositoryBase() *ClientSessionsRepositoryBase {
	return &ClientSessionsRepositoryBase{}
}

func (c *ClientSessionsRepositoryBase) Send(msg model.Message, clientId model.ID) error {
	var expectConn *connect
	for _, conn := range c.conns {
		if conn.Id == clientId {
			expectConn = &conn
			break
		}
	}
	if expectConn == nil {
		return errors.New("client not connected")
	}

	return expectConn.sendMsg(msg)
}

func (c *ClientSessionsRepositoryBase) Connect(client model.Client) (*model.ClientSession, error) {
	var newConn = connect{
		Client: client,
		ch:     make(chan model.Message),
	}
	c.conns = append(c.conns, newConn)
	return newConn.toDomain(), nil
}

func (c *ClientSessionsRepositoryBase) Connected() ([]model.Client, error) {
	var clients []model.Client
	for _, conn := range c.conns {
		clients = append(clients, conn.Client)
	}
	return clients, nil
}

func (c *ClientSessionsRepositoryBase) Disconnect(id model.ID) error {
	for i, conn := range c.conns {
		if conn.Id == id {
			// remove from conns list
			c.conns = append(c.conns[:i], c.conns[i+1:]...)
			// close channel
			conn.closeChan()
		}
	}
	return nil
}

type connect struct {
	model.Client
	ch chan model.Message
}

func (c *connect) sendMsg(msg model.Message) error {
	select {
	case c.ch <- msg:
		return nil
	default:
		return errors.New("channel do not listing")
	}
}

func (c *connect) closeChan() {
	select {
	case _, ok := <-c.ch:
		if ok {
			close(c.ch)
		}
	default:
		close(c.ch)
	}
}

func (c *connect) toDomain() *model.ClientSession {
	return &model.ClientSession{
		Client: c.Client,
		MsgCh:  c.ch,
	}
}
