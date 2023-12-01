package domain

import "chatmerger/internal/domain/model"

type Config interface {
	GetClients() []model.Client
	SetClients(clients []model.Client)

	GetHost() string
	SetHost(host string)

	GetPort() int
	SetPort(port int)
}
