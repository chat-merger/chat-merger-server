package uc

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"fmt"
	"github.com/google/uuid"
)

var _ usecase.CreateClientUc = (*CreateClient)(nil)

type CreateClient struct {
	clientRepos repository.ClientsRepository
}

func NewCreateClient(clientRepos repository.ClientsRepository) *CreateClient {
	return &CreateClient{clientRepos: clientRepos}
}

func (r *CreateClient) CreateClient(input model.CreateClient) error {
	var newClient = model.Client{
		Id:     model.ID(uuid.New().String()),
		Name:   input.Name,
		ApiKey: model.ApiKey(uuid.New().String()),
	}
	err := r.clientRepos.Create(newClient)
	if err != nil {
		return fmt.Errorf("create client: %s", err)
	}
	return nil
}
