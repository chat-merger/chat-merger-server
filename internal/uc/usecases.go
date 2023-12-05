package uc

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
)

var _ usecase.SendMessageToEveryoneExceptUc = (*SendMessageToEveryoneExcept)(nil)

type SendMessageToEveryoneExcept struct{}

func (s *SendMessageToEveryoneExcept) SendMessageToEveryoneExcept(ids []model.ID) error {
	return nil
}

var _ usecase.CreateClientSessionUc = (*CreateClientSession)(nil)

type CreateClientSession struct{}

func (c *CreateClientSession) CreateClientSession(input model.CreateClientSession) (*model.ClientSession, error) {
	return nil, nil
}

var _ usecase.DropClientSessionUc = (*DropClientSession)(nil)

type DropClientSession struct{}

func (d *DropClientSession) DropClientSession(ids []model.ID) error {
	return nil
}

var _ usecase.ClientsListUc = (*ClientsList)(nil)

type ClientsList struct{}

func (c *ClientsList) ClientsList() ([]model.Client, error) {
	return nil, nil
}

var _ usecase.ClientsSessionsListUc = (*ClientsSessionsList)(nil)

type ClientsSessionsList struct{}

func (c *ClientsSessionsList) ClientsConnectionsList() ([]model.ClientSession, error) {
	return nil, nil
}

var _ usecase.CreateClientUc = (*CreateClient)(nil)

type CreateClient struct{}

func (c *CreateClient) CreateClient() error {
	return nil
}

var _ usecase.DeleteClientUc = (*DeleteClient)(nil)

type DeleteClient struct{}

func (d *DeleteClient) DeleteClient() error {
	return nil
}
