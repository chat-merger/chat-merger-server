package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"github.com/google/uuid"
)

var _ usecase.CreateClientUc = (*CreateClient)(nil)

type CreateClient struct {
	clientRepos domain.ClientsRepository
}

func NewCreateClient(clientRepos domain.ClientsRepository) *CreateClient {
	return &CreateClient{clientRepos: clientRepos}
}

func (r *CreateClient) CreateClient(input model.CreateClient) error {
	var newClient = model.Client{
		Id:     model.NewID(uuid.New().String()),
		Name:   input.Name,
		ApiKey: model.NewApiKey(uuid.New().String()),
	}
	var clients = r.clientRepos.GetClients()
	r.clientRepos.SetClients(append(clients, newClient))

	return nil
}
