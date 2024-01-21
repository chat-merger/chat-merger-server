package app

import (
	"chatmerger/internal/domain/repository"
	"chatmerger/internal/usecase"
	"context"
	"sync"
)

type application struct {
	usecases usecasesImpls

	errCh      chan<- error
	wg         *sync.WaitGroup // for indicate all things (servers, handlers...) will stopped
	cancelFunc context.CancelFunc
	ctx        context.Context
}

func newApplication(ctx context.Context, usecases usecasesImpls) (*application, <-chan error) {
	errCh := make(chan error)
	ctx, cancelFunc := context.WithCancel(ctx)
	return &application{
		usecases:   usecases,
		errCh:      errCh,
		wg:         new(sync.WaitGroup),
		cancelFunc: cancelFunc,
		ctx:        ctx,
	}, errCh
}

type usecasesImpls struct {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.SubscribeClientToEventsUc
	usecase.DropClientSubscriptionUc
	usecase.ClientsUc
	usecase.CreateClientUc
	usecase.DeleteClientUc
}

type repositories struct {
	cRepo repository.ClientsRepository
}
