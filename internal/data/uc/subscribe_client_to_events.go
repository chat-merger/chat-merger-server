package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"errors"
	"fmt"
	"slices"
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
	clients, err := c.cRepo.GetClients(model.ClientsFilterExceptStatus{Id: &id})
	if err != nil {
		return fmt.Errorf("get clients: %s", err)
	}

	// client not found
	if len(clients) == 0 {
		return ErrorClientWithGivenApiKeyNotFound
	}

	// check what client already is e.b. subject
	idsOfConn := c.bus.Subjects()
	isConnected := slices.ContainsFunc(idsOfConn, func(subject eventbus.Subject) bool {
		return clients[0].Id == subject
	})

	// client don't to be connected
	if isConnected {
		return ErrorClientAlreadyConnected
	}

	// subscribe to NewMessages
	c.bus.Subscribe(clients[0].Id, handler)

	return nil
}
