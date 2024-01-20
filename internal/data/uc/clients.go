package uc

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/domain/model"
	"chatmerger/internal/usecase"
	"fmt"
)

var _ usecase.ClientsUc = (*Clients)(nil)

type Clients struct {
	cRepo domain.ClientsRepository
}

func NewClients(cRepo domain.ClientsRepository) *Clients {
	return &Clients{cRepo: cRepo}
}

func (c *Clients) Clients(filter model.ClientsFilter) ([]model.Client, error) {
	clients, err := c.cRepo.GetClients(filter)
	if err != nil {
		return nil, fmt.Errorf("clients from repo: %s", err)
	}

	return clients, nil
}
