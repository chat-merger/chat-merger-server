package config

import (
	"chatmerger/internal/domain/model"
)

type Config struct {
	clients []model.Client
}

func (c *Config) GetClients() []model.Client {
	return c.clients
}

func (c *Config) SetClients(clients []model.Client) {
	c.clients = clients
}
