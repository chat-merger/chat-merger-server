package uc

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"fmt"
	"slices"
)

var _ usecase.ClientsUc = (*Clients)(nil)

type Clients struct {
	cRepo repository.ClientsRepository
	bus   *eventbus.EventBus
}

func NewClients(cRepo repository.ClientsRepository, bus *eventbus.EventBus) *Clients {
	return &Clients{cRepo: cRepo, bus: bus}
}

func (c *Clients) Clients(filter model.ClientsFilter) ([]model.ClientWithStatus, error) {
	clients, err := c.cRepo.GetClients(filter.ExceptStatus())
	if err != nil {
		return nil, fmt.Errorf("clients from repo: %s", err)
	}

	withStatus := addStatus(clients, c.bus.Subjects())

	if filter.Status != model.ConnStatusUndefined {
		withStatus = filterByStatus(filter.Status, withStatus)
	}

	return withStatus, nil
}

func filterByStatus(
	status model.ConnStatus,
	clients []model.ClientWithStatus,
) []model.ClientWithStatus {
	newClients := make([]model.ClientWithStatus, 0, len(clients))
	for _, client := range clients {
		if status == client.Status {
			newClients = append(newClients, client)
		}
	}
	return newClients
}

func addStatus(
	clients []model.Client,
	subjects []eventbus.Subject,
) []model.ClientWithStatus {

	newClients := make([]model.ClientWithStatus, 0, len(clients))
	for _, client := range clients {
		isConnected := slices.ContainsFunc(subjects, func(subject eventbus.Subject) bool {
			return client.Id == subject
		})
		withStatus := model.ClientWithStatus{
			Id:     client.Id,
			Name:   client.Name,
			ApiKey: client.ApiKey,
			Status: model.ConnStatusUndefined,
		}
		if isConnected {
			withStatus.Status = model.ConnStatusActive
		} else {
			withStatus.Status = model.ConnStatusInactive
		}
		newClients = append(newClients, withStatus)
	}
	return newClients
}
