package domain

import "chatmerger/internal/domain/model"

type ClientSessionsRepository interface {
	Connect(client model.Client) (*model.ClientSession, error)
	Connected() ([]model.Client, error)
	Disconnect(id model.ID) error
	Send(msg model.Message, clientId model.ID) error
}
