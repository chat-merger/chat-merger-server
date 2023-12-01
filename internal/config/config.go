package config

import (
	"chatmerger/internal/domain/model"
)

type Config struct {
	clients []model.Client
	host    string
	port    int
}

func (c *Config) GetClients() []model.Client {
	return c.clients
}

func (c *Config) SetClients(clients []model.Client) {
	c.clients = clients
}

func (c *Config) GetHost() string {
	return c.host
}

func (c *Config) SetHost(host string) {
	c.host = host
}

func (c *Config) GetPort() int {
	return c.port
}

func (c *Config) SetPort(port int) {
	c.port = port
}
