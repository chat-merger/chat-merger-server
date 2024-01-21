package uc

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"log"
)

var _ usecase.DeleteClientUc = (*DeleteClient)(nil)

type DeleteClient struct {
	clientRepos repository.ClientsRepository
}

func NewDeleteClient(clientRepos repository.ClientsRepository) *DeleteClient {
	return &DeleteClient{clientRepos: clientRepos}
}

func (r *DeleteClient) DeleteClients(ids ...model.ID) error {
	// delete each client
	for _, id := range ids {
		err := r.clientRepos.Delete(id)
		if err != nil {
			log.Printf("[ERROR] delete client: %s", err)
		}
	}
	return nil
}
