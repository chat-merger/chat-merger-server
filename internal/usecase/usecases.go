package usecase

import (
	"chatmerger/internal/domain"
)

var _ domain.Usecases = new(Usecases)

type Usecases struct {
	ClientsConnections
	DropAllClientConnections
	DropClientConnection
	ConnectClient
	ClientsList
	DeleteClient
}

type ClientsConnections struct{}

func (u *ClientsConnections) ClientsConnections() {}

type DropAllClientConnections struct{}

func (u *DropAllClientConnections) DropAllClientConnections() {}

type DropClientConnection struct{}

func (u *DropClientConnection) DropClientConnection(id int) {}

type ConnectClient struct{}

func (u *ConnectClient) ConnectClient() {}

type ClientsList struct{}

func (u *ClientsList) ClientsList() {}

type DeleteClient struct{}

func (u *DeleteClient) DeleteClient() {}
