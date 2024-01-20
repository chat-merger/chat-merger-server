package client_sessions_repository

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"errors"
	"sync"
)

var _ domain.ClientSessionsRepository = (*ClientSessionsRepositoryBase)(nil)

type ClientSessionsRepositoryBase struct {
	conns []connect
	mu    sync.RWMutex
}

func NewClientSessionsRepositoryBase() *ClientSessionsRepositoryBase {
	return &ClientSessionsRepositoryBase{
		conns: make([]connect, 0),
		mu:    sync.RWMutex{},
	}
}

func (c *ClientSessionsRepositoryBase) Send(msg model.Message, clientId model.ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.Lock()
	defer c.mu.Unlock()
	var newConn = connect{
		Client: client,
		ch:     make(chan model.Message),
	}
	c.conns = append(c.conns, newConn)
	return newConn.toDomain(), nil
}

func (c *ClientSessionsRepositoryBase) Connected() ([]model.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var clients []model.Client
	for _, conn := range c.conns {
		clients = append(clients, conn.Client)
	}
	return clients, nil
}

func (c *ClientSessionsRepositoryBase) Disconnect(id model.ID) error {
	c.mu.Lock()
	defer c.mu.Unlock()
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
