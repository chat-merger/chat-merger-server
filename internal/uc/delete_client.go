package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"slices"
)

var _ usecase.DeleteClientUc = (*DeleteClient)(nil)

type DeleteClient struct {
	clientRepos domain.ClientsRepository
}

func NewDeleteClient(clientRepos domain.ClientsRepository) *DeleteClient {
	return &DeleteClient{clientRepos: clientRepos}
}

func (r *DeleteClient) DeleteClients(ids []model.ID) error {
	var clients = r.clientRepos.GetClients()
	// remove elements..
	var newClientsList = slices.DeleteFunc(clients, func(client model.Client) bool {
		// witch in the list
		return slices.ContainsFunc(ids, func(id model.ID) bool {
			return id == client.Id
		})
	})

	r.clientRepos.SetClients(newClientsList)

	return nil
}
