package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
)

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
