package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
)

var _ usecase.ClientsUc = (*ClientSession)(nil)

type ClientSession struct {
	sRepo domain.ClientSessionsRepository
}

func NewClientSession(sRepo domain.ClientSessionsRepository) *ClientSession {
	return &ClientSession{sRepo: sRepo}
}

func (c *ClientSession) Clients(filter usecase.ClientsFilter) ([]model.Client, error) {
	connected, err := c.sRepo.Connected()
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Client, len(connected))

	for _, client := range connected {
		if filter.Id != nil && client.Id == *filter.Id &&
			filter.Name != nil && client.Name == *filter.Name &&
			filter.ApiKey != nil && client.ApiKey == *filter.ApiKey {
			filtered = append(filtered, client)
		}
	}

	return filtered, nil
}
