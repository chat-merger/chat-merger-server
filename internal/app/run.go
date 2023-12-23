package app

import (
	"chatmerger/internal/api/grpc_side"
	"chatmerger/internal/api/http_side"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"chatmerger/internal/uc"
	"context"
	"fmt"
	"log"
)

func Run(ctx context.Context) error {

	repos, err := initRepositories()
	if err != nil {
		return fmt.Errorf("create repositories: %s", err)
	}

	var app = &application{
		commonDeps: commonDeps{
			usecases: newUsecases(*repos),
			ctx:      ctx,
		},
		httpSideCfg: http_side.Config{
			Host: "localhost",
			Port: 8081,
		},
		grpcSideCfg: grpc_side.Config{
			Host: "localhost",
			Port: 8080,
		},
	}

	// create and run clients api handler
	go app.runGrpcSideServer()
	// crate and run admin panel api handler
	go app.runHttpSideServer()

	<-ctx.Done()
	return nil
}

func (a *application) runHttpSideServer() {
	h := http_side.NewHttpSideServer(a.httpSideCfg, a.usecases)
	err := h.Serve(a.ctx)
	if err != nil {
		log.Fatalf("http side server serve: %s", err)
	}
}

func (a *application) runGrpcSideServer() {
	h := grpc_side.NewGrpcSideServer(a.grpcSideCfg, a.usecases)
	err := h.Serve(a.ctx)
	if err != nil {
		log.Fatalf("grpc side server serve: %s", err)
	}
}

func initRepositories() (*repositories, error) {

	sessionsRepo := NewClientSessionsRepositoryBase()
	clientsRepo, err := NewClientsRepositoryBase(Config{FilePath: "./clients.json"})
	if err != nil {
		return nil, fmt.Errorf("create clients repository: %s", err)
	}
	return &repositories{
		clientsRepo:  clientsRepo,
		sessionsRepo: sessionsRepo,
	}, nil

}

func newUsecases(repos repositories) *usecasesImpls {
	return &usecasesImpls{
		ConnectedClientsListUc: uc.NewConnectedClientsList(repos.sessionsRepo),
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(repos.sessionsRepo),
		CreateClientSessionUc:              uc.NewCreateClientSession(repos.clientsRepo, repos.sessionsRepo),
		DropClientSessionUc:                uc.NewDropClientSession(repos.sessionsRepo),
		// admin panel api server
		ClientsListUc:  uc.NewClientsList(repos.clientsRepo),
		CreateClientUc: uc.NewCreateClient(repos.clientsRepo),
		DeleteClientUc: uc.NewDeleteClient(repos.clientsRepo),
	}
}
