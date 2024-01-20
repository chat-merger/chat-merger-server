package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/service/msgbus"
	"chatmerger/internal/usecase"
	"fmt"
)

var _ usecase.DropClientSubscriptionUc = (*DropClientSubscription)(nil)

type DropClientSubscription struct {
	cRepo domain.ClientsRepository
	bus   *msgbus.MessagesBus
}

func NewDropClientSubscription(repo domain.ClientsRepository, bus *msgbus.MessagesBus) *DropClientSubscription {
	return &DropClientSubscription{cRepo: repo, bus: bus}
}

func (d *DropClientSubscription) DropClientSubscription(ids []model.ID) error {

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
		err = d.cRepo.UpdateClient(client.Id, client)
		if err != nil {
			return fmt.Errorf("update calient status: %s", err)
		}

		d.bus.Unsubscribe(id)
	}

	return nil
}
