package domain

import "chatmerger/internal/domain/model"

type ClientConnectionRepository interface {
	Connect(client model.Client) (*model.ClientConnection, error)
	Connected() ([]model.Client, error)
	Disconnect(id int) error
}
