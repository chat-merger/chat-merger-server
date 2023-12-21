package app

import (
	adm "chatmerger/internal/api/admin"
	"chatmerger/internal/api/pb"
	"chatmerger/internal/domain"
	. "chatmerger/internal/repositories/client_sessions_repository"
	. "chatmerger/internal/repositories/clients_repository"
	"chatmerger/internal/uc"
	"chatmerger/internal/usecase"
	"context"
	"log"
)

type State uint8

const (
	Working State = iota
	Running
	Stopping
	Stopped
	RunningFailure
	StoppingFailure
	InternalError
)

type application struct {
	clientRepository  domain.ClientsRepository
	sessionRepository domain.ClientsSessionRepository
	apiHandler        domain.Handler
	adminHandler      domain.Handler
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

// use for graceful shutdown
func (a *application) shutdown() {
	// todo: a.changeStatus(Stopping)
	cc, err := a.sessionRepository.Connected()
	if err == nil {
		for _, client := range cc {
			a.sessionRepository.Disconnect(client.Id)
		}
	}
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
	var clientsApiHandler = pb.NewClientsServer(
		pb.Config{
			Host: "localhost",
			Port: 8080,
		},
		pb.Usecases{
			CreateAndSendMsgToEveryoneExceptUc: usecases,
			CreateClientSessionUc:              usecases,
			DropClientSessionUc:                usecases,
		},
	)
	go serveHandler(clientsApiHandler, ctx)

	// crate and run admin panel api handler
	var adminPanelApiHandler = adm.NewAdminPanelServer(
		adm.Config{
			Host: "localhost",
			Port: 8081,
		},
		adm.Usecases{
			CreateClientUc:         usecases,
			DeleteClientUc:         usecases,
			ClientsListUc:          usecases,
			ConnectedClientsListUc: usecases,
		},
	)
	go serveHandler(adminPanelApiHandler, ctx)

	// todo:
	var _ = &application{
		clientRepository:  clientsRepo,
		sessionRepository: sessionsRepo,
		apiHandler:        clientsApiHandler,
		adminHandler:      adminPanelApiHandler,
	}

	return nil
}

func (a *application) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		a.shutdown()
	}
}

func serveHandler(h domain.Handler, ctx context.Context) {
	err := h.Serve(ctx)
	if err != nil {
		log.Fatalf("handler serve: %s", err)
	}
}
