package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"fmt"
)

var _ usecase.DropClientSubscriptionUc = (*DropClientSubscription)(nil)

type DropClientSubscription struct {
	cRepo repository.ClientsRepository
	bus   *eventbus.EventBus
}

func NewDropClientSubscription(repo repository.ClientsRepository, bus *eventbus.EventBus) *DropClientSubscription {
	return &DropClientSubscription{cRepo: repo, bus: bus}
}

func (d *DropClientSubscription) DropClientSubscription(ids ...model.ID) error {

	for _, id := range ids {
		clients, err := d.cRepo.GetClients(model.ClientsFilter{Id: &id})
		if err != nil {
			return fmt.Errorf("get clients: %s", err)
		}

		// client not found
		if len(clients) == 0 {
			return ErrorClientWithGivenApiKeyNotFound
		}

		// take first
		client := clients[0]

		//// client should be connected
		//if client.Status == model.ConnStatusInactive {
		//	return ErrorClientAlreadyConnected
		//}

		client.Status = model.ConnStatusInactive
		err = d.cRepo.Update(client.Id, client)
		if err != nil {
			return fmt.Errorf("update calient status: %s", err)
		}

		d.bus.Unsubscribe(id)
	}

	return nil
}
