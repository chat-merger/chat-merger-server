package usecase

import (
	"chatmerger/internal/component/eventbus"
	"chatmerger/internal/domain/model"
)

// grpc api usecases

type CreateAndSendMsgToEveryoneExceptUc interface {
	CreateAndSendMsgToEveryoneExcept(msg model.CreateMessage, ids ...model.ID) (*model.Message, error)
}
type SubscribeClientToEventsUc interface {
	SubscribeClientToEvents(id model.ID, handler eventbus.Handler) error
}
type DropClientSubscriptionUc interface {
	DropClientSubscription(ids ...model.ID) error
}

// admin site usecases

type ClientsUc interface {
	Clients(filter model.ClientsFilter) ([]model.Client, error)
}
type CreateClientUc interface {
	CreateClient(input model.CreateClient) error
}
type DeleteClientUc interface {
	DeleteClients(ids ...model.ID) error
}
