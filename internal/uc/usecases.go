package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
)

var _ usecase.SendMessageToEveryoneExceptUc = (*SendMessageToEveryoneExcept)(nil)

type SendMessageToEveryoneExcept struct{}

func (r *SendMessageToEveryoneExcept) SendMessageToEveryoneExcept(ids []model.ID) error {
	return nil
}

var _ usecase.ClientsListUc = (*ClientsList)(nil)

type ClientsList struct {
	clientsRepos domain.ClientsRepository
}

func NewClientsList(clientsRepos domain.ClientsRepository) *ClientsList {
	return &ClientsList{clientsRepos: clientsRepos}
}

func (r *ClientsList) ClientsList() ([]model.Client, error) {
	return r.clientsRepos.GetClients(), nil
}

var _ usecase.ConnectedClientsListUc = (*ConnectedClientsList)(nil)

type ConnectedClientsList struct {
	sessionsRepo domain.ClientsSessionRepository
}

func NewConnectedClientsList(sessionsRepo domain.ClientsSessionRepository) *ConnectedClientsList {
	return &ConnectedClientsList{sessionsRepo: sessionsRepo}
}

func (r *ConnectedClientsList) ConnectedClientsList() ([]model.Client, error) {
	return r.sessionsRepo.Connected()
}
