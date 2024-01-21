package uc

import (
	"chatmerger/internal/domain/model"
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"fmt"
)

var _ usecase.ClientsUc = (*Clients)(nil)

type Clients struct {
	cRepo repository.ClientsRepository
}

func NewClients(cRepo repository.ClientsRepository) *Clients {
	return &Clients{cRepo: cRepo}
}

func (c *Clients) Clients(filter model.ClientsFilter) ([]model.Client, error) {
	clients, err := c.cRepo.GetClients(filter)
	if err != nil {
		return nil, fmt.Errorf("clients from repo: %s", err)
	}

	return clients, nil
}
