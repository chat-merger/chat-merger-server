package client_sessions_repository

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
)

var _ domain.ClientsSessionRepository = (*ClientSessionsRepositoryBase)(nil)

type ClientSessionsRepositoryBase struct {
	conns []*connect
}

func (c *ClientSessionsRepositoryBase) Connect(client model.Client) (*model.ClientSession, error) {
	var newConn = &connect{
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
