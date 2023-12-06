package clients_repository

import (
	"chatmerger/internal/domain/model"
)

type ClientsRepositoryBase struct {
	clients []model.Client
}

func (c *ClientsRepositoryBase) GetClients() []model.Client {
	return c.clients
}

func (c *ClientsRepositoryBase) SetClients(clients []model.Client) {
	c.clients = clients
}
