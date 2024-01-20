package domain

import "chatmerger/internal/domain/model"

type ClientsRepository interface {
	GetClients(filter model.ClientsFilter) ([]model.Client, error)
	SetClients(clients []model.Client) error
	UpdateClient(id model.ID, new model.Client) error
	// todo: + SetStatus, working with specific client
}
