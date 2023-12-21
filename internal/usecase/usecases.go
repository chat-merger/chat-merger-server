package usecase

import "chatmerger/internal/domain/model"

// grpc api usecases

type CreateAndSendMsgToEveryoneExceptUc interface {
	CreateAndSendMsgToEveryoneExcept(msg model.CreateMessage, ids []model.ID) error
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
type ConnectedClientsListUc interface {
	ConnectedClientsList() ([]model.Client, error)
}
type CreateClientUc interface {
	CreateClient(input model.CreateClient) error
}
type DeleteClientUc interface {
	DeleteClients(ids []model.ID) error
}
