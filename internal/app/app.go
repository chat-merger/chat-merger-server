package app

import (
	"chatmerger/internal/domain"
	"chatmerger/internal/usecase"
	"context"
	"sync"
)

type application struct {
	commonDeps
	errCh      chan<- error
	wg         *sync.WaitGroup // for indicate all things (servers, handlers...) will stopped
	cancelFunc context.CancelFunc
	ctx        context.Context
}

func newApplication(ctx context.Context, commonDeps commonDeps) (*application, <-chan error) {
	errCh := make(chan error)
	ctx, cancelFunc := context.WithCancel(ctx)
	return &application{
		commonDeps: commonDeps,
		errCh:      errCh,
		wg:         new(sync.WaitGroup),
		cancelFunc: cancelFunc,
		ctx:        ctx,
	}, errCh
}

type commonDeps struct {
	usecases *usecasesImpls
	ctx      context.Context
}

type usecasesImpls struct {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
	usecase.ClientsListUc
	usecase.ConnectedClientsListUc
	usecase.CreateClientUc
	usecase.DeleteClientUc
}

type repositories struct {
	clientsRepo  domain.ClientsRepository
	sessionsRepo domain.ClientSessionsRepository
}
