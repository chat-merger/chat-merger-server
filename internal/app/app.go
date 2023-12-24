package app

import (
	"chatmerger/internal/api/grpc_side"
	"chatmerger/internal/api/http_side"
	"chatmerger/internal/config"
	"chatmerger/internal/domain"
	"chatmerger/internal/usecase"
	"context"
	"sync"
)

type application struct {
	commonDeps
	httpSideCfg http_side.Config
	grpcSideCfg grpc_side.Config
	errCh       chan<- error
	wg          *sync.WaitGroup // for indicate all things (servers, handlers...) will stopped
	cancelFunc  context.CancelFunc
}

func newApplication(commonDeps commonDeps, cfg *config.Config) (*application, <-chan error) {
	errCh := make(chan error)
	return &application{
		errCh:      errCh,
		commonDeps: commonDeps,
		httpSideCfg: http_side.Config{
			Host: "localhost",
			Port: cfg.HttpServerPort,
		},
		grpcSideCfg: grpc_side.Config{
			Host: "localhost",
			Port: cfg.GrpcServerPort,
		},
		wg: new(sync.WaitGroup),
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
