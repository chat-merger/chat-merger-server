package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/service/msgbus"
	"chatmerger/internal/usecase"
	"errors"
	"fmt"
)

var _ usecase.SubscribeClientToNewMsgsUc = (*SubscribeClientToNewMsgs)(nil)

type SubscribeClientToNewMsgs struct {
	cRepo domain.ClientsRepository
	bus   *msgbus.MessagesBus
}

func NewSubscribeClientToNewMsgs(
	cRepo domain.ClientsRepository,
	bus *msgbus.MessagesBus,
) *SubscribeClientToNewMsgs {
	return &SubscribeClientToNewMsgs{cRepo: cRepo, bus: bus}
}

var (
	ErrorClientWithGivenApiKeyNotFound = errors.New("client with given ApiKey not found")
	ErrorClientAlreadyConnected        = errors.New("client already connected")
)

func (c *SubscribeClientToNewMsgs) SubscribeClientToNewMsgs(id model.ID, handler func(newMsg model.Message) error) error {
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
	err = c.cRepo.UpdateClient(client.Id, client)
	if err != nil {
		return fmt.Errorf("update calient status: %s", err)
	}

	// subscribe to NewMessages
	c.bus.Subscribe(client.Id, handler)

	return nil
}
