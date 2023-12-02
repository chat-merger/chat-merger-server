package domain

import "chatmerger/internal/domain/model"

type ClientRepository interface {
	GetClients() []model.Client
	SetClients(clients []model.Client)
}
