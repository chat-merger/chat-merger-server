package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"errors"
	"fmt"
)

var _ usecase.SubscribeClientToEventsUc = (*SubscribeClientToEvents)(nil)

type SubscribeClientToEvents struct {
	cRepo repository.ClientsRepository
	bus   *eventbus.EventBus
}

func NewSubscribeClientToEvents(
	cRepo repository.ClientsRepository,
	bus *eventbus.EventBus,
) *SubscribeClientToEvents {
	return &SubscribeClientToEvents{cRepo: cRepo, bus: bus}
}

var (
	ErrorClientWithGivenApiKeyNotFound = errors.New("client with given ApiKey not found")
	ErrorClientAlreadyConnected        = errors.New("client already connected")
)

func (c *SubscribeClientToEvents) SubscribeClientToEvents(id model.ID, handler eventbus.Handler) error {
	clients, err := c.cRepo.GetClients(model.ClientsFilter{Id: &id})
	if err != nil {
		return fmt.Errorf("get clients: %s", err)
	}

	// client not found
	if len(clients) == 0 {
		return ErrorClientWithGivenApiKeyNotFound
	}

	// take first
	client := clients[0]

	// client don't to be connected
	if client.Status == model.ConnStatusActive {
		return ErrorClientAlreadyConnected
	}

	client.Status = model.ConnStatusActive
	err = c.cRepo.Update(client.Id, client)
	if err != nil {
		return fmt.Errorf("update calient status: %s", err)
	}

	// subscribe to NewMessages
	c.bus.Subscribe(client.Id, handler)

	return nil
}
