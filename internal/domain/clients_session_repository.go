package domain

import "chatmerger/internal/domain/model"

type ClientsSessionRepository interface {
	Connect(client model.Client) (*model.ClientSession, error)
	Connected() ([]model.Client, error)
	Disconnect(id model.ID) error
}
