package repository

import "chatmerger/internal/domain/model"

type ClientsRepository interface {
	Create(client model.Client) error
	GetClients(filter model.ClientsFilterExceptStatus) ([]model.Client, error)
	Update(id model.ID, new model.Client) error
	Delete(id model.ID) error
}
