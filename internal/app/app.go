package app

import (
	"chatmerger/internal/api/grpc_side"
	"chatmerger/internal/api/http_side"
	"chatmerger/internal/domain"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"chatmerger/internal/uc"
	"chatmerger/internal/usecase"
	"context"
	"log"
)

type application struct {
	clientRepository  domain.ClientsRepository
	sessionRepository domain.ClientsSessionRepository
	grpcHandler       domain.Handler
	httpHandler       domain.Handler
	usecases          *Usecases
}

type Usecases struct {
	usecase.CreateAndSendMsgToEveryoneExceptUc
	usecase.CreateClientSessionUc
	usecase.DropClientSessionUc
	usecase.ClientsListUc
	usecase.ConnectedClientsListUc
	usecase.CreateClientUc
	usecase.DeleteClientUc
}

func Run(ctx context.Context) error {
	var sessionsRepo = &ClientSessionsRepositoryBase{}
	clientsRepo, err := NewClientsRepositoryBase(Config{FilePath: "./clients.json"})
	if err != nil {
		log.Fatalf("create clients repository: %s", err)
	}
	var usecases = &Usecases{
		ConnectedClientsListUc: uc.NewConnectedClientsList(sessionsRepo),
		// clients api server
		CreateAndSendMsgToEveryoneExceptUc: uc.NewCreateAndSendMsgToEveryoneExcept(sessionsRepo),
		CreateClientSessionUc:              uc.NewCreateClientSession(clientsRepo, sessionsRepo),
		DropClientSessionUc:                uc.NewDropClientSession(sessionsRepo),
		// admin panel api server
		ClientsListUc:  uc.NewClientsList(clientsRepo),
		CreateClientUc: uc.NewCreateClient(clientsRepo),
		DeleteClientUc: uc.NewDeleteClient(clientsRepo),
	}

	// create and run clients api handler
	var grpcHandler = grpc_side.NewClientsServer(
		grpc_side.Config{
			Host: "localhost",
			Port: 8080,
		},
		grpc_side.Usecases{
			CreateAndSendMsgToEveryoneExceptUc: usecases,
			CreateClientSessionUc:              usecases,
			DropClientSessionUc:                usecases,
		},
	)
	go serveHandler(grpcHandler, ctx)

	// crate and run admin panel api handler
	var httpHandler = http_side.NewAdminPanelServer(
		http_side.Config{
			Host: "localhost",
			Port: 8081,
		},
		http_side.Usecases{
			CreateClientUc:         usecases,
			DeleteClientUc:         usecases,
			ClientsListUc:          usecases,
			ConnectedClientsListUc: usecases,
		},
	)
	go serveHandler(httpHandler, ctx)

	var app = &application{
		clientRepository:  clientsRepo,
		sessionRepository: sessionsRepo,
		grpcHandler:       grpcHandler,
		httpHandler:       httpHandler,
	}

	<-ctx.Done()
	app.shutdown()
	return nil
}

func serveHandler(h domain.Handler, ctx context.Context) {
	err := h.Serve(ctx)
	if err != nil {
		log.Fatalf("handler serve: %s", err)
	}
}

// use for graceful shutdown
func (a *application) shutdown() {
	cc, err := a.sessionRepository.Connected()
	if err == nil {
		for _, client := range cc {
			err := a.sessionRepository.Disconnect(client.Id)
			if err != nil {
				log.Printf("session disconnect: %s\n", err)
			}
		}
	}
}
