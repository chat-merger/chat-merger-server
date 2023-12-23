package domain

import "chatmerger/internal/domain/model"

type ClientsRepository interface {
	GetClients() ([]model.Client, error)
	SetClients(clients []model.Client) error
}
