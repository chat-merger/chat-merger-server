package usecase

import "chatmerger/internal/domain/model"

// grpc api usecases

type SendMessageToEveryoneExceptUc interface {
	SendMessageToEveryoneExcept(ids []model.ID) error
}

type CreateClientSessionUc interface {
	CreateClientSession(input model.CreateClientSession) (*model.ClientSession, error)
}
type DropClientSessionUc interface {
	DropClientSession(ids []model.ID) error
}

// admin site usecases

type ClientsListUc interface {
	ClientsList() ([]model.Client, error)
}
type ClientsSessionsListUc interface {
	ClientsConnectionsList() ([]model.ClientSession, error)
}
type CreateClientUc interface {
	CreateClient(input model.CreateClient) error
}
type DeleteClientUc interface {
	DeleteClient() error
}
