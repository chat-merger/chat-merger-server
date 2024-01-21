package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"log"
)

var _ usecase.DeleteClientUc = (*DeleteClient)(nil)

type DeleteClient struct {
	clientRepos repository.ClientsRepository
	bus         *eventbus.EventBus
}

func NewDeleteClient(clientRepos repository.ClientsRepository, bus *eventbus.EventBus) *DeleteClient {
	return &DeleteClient{clientRepos: clientRepos, bus: bus}
}

func (r *DeleteClient) DeleteClients(ids ...model.ID) error {
	// delete each client
	for _, id := range ids {
		err := r.clientRepos.Delete(id)
		if err != nil {
			log.Printf("[ERROR] delete client: %s", err)
		}
		r.bus.Unsubscribe(id)
	}
	return nil
}
