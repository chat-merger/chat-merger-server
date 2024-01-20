package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"fmt"
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
	clients, err := r.clientRepos.GetClients(model.ClientsFilter{})
	if err != nil {
		return fmt.Errorf("getting client list: %s", err)
	}
	err = r.clientRepos.SetClients(append(clients, newClient))
	if err != nil {
		return fmt.Errorf("setting clients list: %s", err)
	}
	return nil
}
