package domain

import "chatmerger/internal/domain/model"

type ClientsRepository interface {
	GetClients() []model.Client
	SetClients(clients []model.Client)
}
